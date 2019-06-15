package menu

import (
	"log"
	"os/exec"
	"strings"
)

type output struct {
	code int
	body string
}

func execCommand(command string) (string, error) {
	if command == "todo" {
		return "", nil
	}

	log.Println("Command", command)
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
