package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// IsValidFn validates the input string and returns true if valid otherwise false.
type IsValidFn = func(string) bool

var defaultIsValidFn = func(s string) bool { return true }

// requestInput prompts for input, reads, validates it and returns it.
// If 'required' is specified, the function will keep asking for inputs until a valid input is entered.
func requestInput(prompt string, required bool, isValidFn IsValidFn, ctx IOContext) string {
	var str string
	reader := bufio.NewReader(ctx.StdinReader)

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

func questionYN(prompt string, ctx IOContext) bool {
	var str string
	reader := bufio.NewReader(ctx.StdinReader)

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

func requestString(prompt string, required bool, ctx IOContext) string {
	return requestInput(prompt, required, defaultIsValidFn, ctx)
}

func requestOptionalExistingDirectory(prompt string, defaultVal string, ctx IOContext) string {
	return requestInput(
		formatOptionalPrompt(prompt, defaultVal),
		false,
		func(path string) bool {
			if path == "" {
				return true
			}
			expandedPath := expandPath(path)
			_, err := os.Stat(expandedPath)
			exists := !os.IsNotExist(err)
			if !exists {
				if questionYN(fmt.Sprintf("the directory '%s' does not exist. create it now?", path), ctx) {
					exists = os.MkdirAll(expandedPath, 0755) == nil
				} else {
					_, _ = printfRed("the directory '%s' does not exist\r\n", path)
				}
			}

			return exists
		},
		ctx,
	)
}

func requestUint(prompt string, required bool, ctx IOContext) uint {
	var str string
	for {
		str = requestInput(prompt, required, defaultIsValidFn, ctx)
		if str == "" {
			return 0
		}
		if v, err := strconv.ParseUint(str, 10, 32); err == nil {
			return uint(v)
		}

		_, _ = printRed("please enter an unsigned integer value\r\n")
	}
}

func requestOptionalBool(prompt string, defaultVal bool, ctx IOContext) bool {
	var str string
	for {
		str = requestInput(formatOptionalPrompt(prompt, defaultVal), false, defaultIsValidFn, ctx)
		if str == "" {
			return false
		}
		if v, err := strconv.ParseBool(str); err == nil {
			return v
		}

		_, _ = printRed("please enter 'true', 'false', '1' or '0'\r\n")
	}
}

func requestCommandLine(prompt string, required bool, ctx IOContext) []string {
	var final []string
	str := requestInput(prompt, required, defaultIsValidFn, ctx)
	if str != "" {
		final = parseCommand(str)
	}

	return final
}

func formatOptionalPrompt(prompt string, defaultVal interface{}) string {
	return fmt.Sprintf("%s (%v)", prompt, defaultVal)
}
