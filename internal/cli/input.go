package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


type IsValidFn = func(string) bool

var defaultIsValidFn = func(s string) bool { return true }

func RequestInput(prompt string, required bool, isValidFn IsValidFn) string {
	var str string
	reader := bufio.NewReader(StdinReader)

	displayPrompt := func() {
		if required {
			_, _ = fmt.Printf("%s %s: ", prompt, sprintRed("*"))
		} else {
			_, _ = fmt.Printf("%s %s: ", prompt, sprintGreen("?"))
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
	reader := bufio.NewReader(StdinReader)

	displayPrompt := func() {
		_, _ = fmt.Printf("%s (y/n|enter): ", prompt)
	}

	for {
		displayPrompt()
		str, _ = reader.ReadString('\n')
		str = strings.TrimSpace(strings.ToLower(str))
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

func RequestOptionalExistingDirectory(prompt string, defaultVal string) string {
	return RequestInput(
		formatOptionalPrompt(prompt, defaultVal),
		false,
		func(path string) bool {
			if path == "" {
				return true
			}
			_, err := os.Stat(expandPath(path))
			exists := !os.IsNotExist(err)
			if !exists {
				_, _ = printfRed("the directory '%s' does not exist\r\n", path)
			}

			return exists
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
		} else {
			_, _ = printRed("please enter an unsigned integer value\r\n")
		}
	}
}

func RequestOptionalBool(prompt string, defaultVal bool) bool {
	var str string
	for {
		str = RequestInput(formatOptionalPrompt(prompt, defaultVal), false, defaultIsValidFn)
		if str == "" {
			return false
		}
		if v, err := strconv.ParseBool(str); err == nil {
			return v
		} else {
			_, _ = printRed("please enter 'true', 'false', '1' or '0'\r\n")
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

func formatOptionalPrompt(prompt string, defaultVal interface{}) string {
	return fmt.Sprintf("%s (%v)", prompt, defaultVal)
}
