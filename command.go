package menu

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
)

var errTimeout = errors.New("command timed out")

type command struct {
	commandDetails

	busy           bool
	updateInterval time.Duration
	timeout        time.Duration
	ticker         *time.Ticker

	out   *string
	error error
}

type commandDetails struct {
	ShellCommand   string `json:"cmd"`
	UpdateInterval string `json:"update_interval"`
	Timeout        string `json:"timeout"`
}

// Can be both a structure or a command string.
func (c *command) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '{' {
		err := json.Unmarshal(b, &c.commandDetails)
		if err != nil {
			return err
		}
		if c.UpdateInterval != "" {
			c.updateInterval, err = time.ParseDuration(c.UpdateInterval)
			if err != nil {
				return err
			}
		}
		if c.Timeout != "" {
			c.timeout, err = time.ParseDuration(c.Timeout)
			if err != nil {
				return err
			}
		}
		return nil
	}

	str, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	c.ShellCommand = str

	return nil
}

func (c *command) exec() {
	c.busy = true
	defer func() { c.busy = false }()

	switch {
	case c.ShellCommand != "":
		c.out, c.error = execShellCommand(c.ShellCommand, c.timeout)
	}
}

func (c *command) keepUpdated() {
	c.exec()
	if c.updateInterval == 0 {
		return
		// c.UpdateInterval = 3 * time.Second
	}

	c.ticker = time.NewTicker(c.updateInterval)
	go func() {
		for range c.ticker.C {
			c.exec()
			if c.error != nil {
				log.Printf("Command failed: %v", c.error)
			}
		}
	}()
}

func (c *command) resetTimer() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.keepUpdated()
}

func execShellCommand(shellCommand string, timeout time.Duration) (*string, error) {
	log.Println("Command:", shellCommand)
	var out bytes.Buffer
	cmd := exec.Command("bash", "-c", shellCommand)
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if err := waitWithTimeout(cmd, timeout); err != nil {
		return nil, err
	}

	strOut := strings.TrimSpace(out.String())
	// log.Println("Output:", strOut)
	return &strOut, nil
}

func waitWithTimeout(cmd *exec.Cmd, timeout time.Duration) error {
	if timeout == 0 {
		return cmd.Wait()
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Failed to kill command after timeout: %v", err)
		}
		return errTimeout
	case err := <-errCh:
		return err
	}
}
