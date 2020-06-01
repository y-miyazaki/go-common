package command

import (
	"fmt"
	"os/exec"
)

// Output execute exec.Command and output command.
func CombinedOutput(command string, output bool, options ...string) (string, error) {
	out, err := exec.Command(command, options...).CombinedOutput()
	outputCommand := command + " "
	for _, s := range options {
		outputCommand = s + " "
	}
	if output {
		fmt.Println(outputCommand)
		fmt.Println(string(out))
	}
	return string(out), err
}

// OutputStr execute exec.Command and output command.
func CombinedOutputStr(command string, output bool) (string, error) {
	out, err := exec.Command("sh", "-c", command).CombinedOutput()
	if output {
		fmt.Println(command)
		fmt.Println(string(out))
	}
	return string(out), err
}
