package cli

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type IsValidFn = func(string) bool

var defaultIsValidFn = func(s string) bool { return true }

func RequestInput(prompt string, required bool, isValidFn IsValidFn) string {
	var str string
	reader := bufio.NewReader(os.Stdin)

	displayPrompt := func() {
		if required {
			printfBold("%s %s: ", prompt, sprintRed("*"))
		} else {
			printfBold("%s %s: ", prompt, sprintGreen("?"))
		}
	}

	for {
		displayPrompt()
		str, _ = reader.ReadString('\n')
		str = strings.TrimSpace(str)
		if (str == "" && required) || !isValidFn(str) {
			continue
		} else {
			return str
		}
	}
}

func QuestionYN(prompt string) bool {
	var str string
	reader := bufio.NewReader(os.Stdin)

	displayPrompt := func() {
		printfBold("%s (y/n): ", prompt)
	}

	for {
		displayPrompt()
		str, _ = reader.ReadString('\n')
		str = strings.TrimSpace(str)
		if str == "y" {
			return true
		} else if str == "n" || str == "" {
			return false
		}
	}
}

func RequestString(prompt string, required bool) string {
	return RequestInput(prompt, required, defaultIsValidFn)
}

func RequestExistingDirectory(prompt string, required bool) string {
	return RequestInput(
		prompt,
		required,
		func(path string) bool {
			if path == "" {
				return true
			}
			_, err := os.Stat(path)
			return !os.IsNotExist(err)
		},
	)
}

func RequestUint(prompt string, required bool) uint {
	var str string
	for {
		str = RequestInput(prompt, required, defaultIsValidFn)
		if str == "" {
			return 0
		}
		if v, err := strconv.ParseUint(str, 10, 32); err == nil {
			return uint(v)
		}
	}
}

func RequestBool(prompt string, required bool) bool {
	var str string
	for {
		str = RequestInput(prompt, required, defaultIsValidFn)
		if str == "" {
			return false
		}
		if v, err := strconv.ParseBool(str); err == nil {
			return v
		}
	}
}

func RequestCommandLine(prompt string, required bool) []string {
	var final []string
	str := RequestInput(prompt, required, defaultIsValidFn)
	if str != "" {
		final = parseCommand(str)
	}

	return final
}
