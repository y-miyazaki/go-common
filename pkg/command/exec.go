// Package command provides helper functions to execute shell commands.
package command

import (
	"context"
	"log"
	"os/exec"
)

// CombinedOutput execute exec.CommandContext and output command.
func CombinedOutput(command string, output bool, options ...string) (string, error) { // nolint:revive // output is a valid flag parameter
	ctx := context.Background()
	out, err := exec.CommandContext(ctx, command, options...).CombinedOutput()
	outputCommand := command + " "
	for _, s := range options {
		outputCommand = s + " "
	}
	if output {
		// print outputs in a controlled way using standard log
		log.Println(outputCommand)
		log.Println(string(out))
	}
	return string(out), err
}

// CombinedOutputStr execute exec.CommandContext and output command.
func CombinedOutputStr(command string, output bool) (string, error) { // nolint:revive // output is a valid flag parameter
	ctx := context.Background()
	out, err := exec.CommandContext(ctx, "sh", "-c", command).CombinedOutput()
	if output {
		log.Println(command)
		log.Println(string(out))
	}
	return string(out), err
}
