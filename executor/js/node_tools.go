package js

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func formatArgs(args []string) string {
	argsWithQuotes := make([]string, 0, len(args))
	for _, arg := range args {
		argsWithQuotes = append(argsWithQuotes, fmt.Sprintf("'%s'", arg))
	}
	return strings.Join(argsWithQuotes, ",")
}

func prepareNodeCommand(path string, command string, args []string) string {
	return fmt.Sprintf("console.log(require('%s').%s(%s))", path, command, formatArgs(args))
}

//TODO: rewrite this method
func parseNodeArrayResult(src string) ([]string, error) {

	src = strings.Trim(src, "[]")
	result := make([]string, 0)
	buf := ""
	var state = 0

	for _, nextRune := range src {
		switch state {
		case 0:
			if nextRune == ' ' {
			} else if nextRune == '\'' {
				state = 1
			} else if nextRune == ',' {
				result = append(result, buf)
				buf = ""
			} else {
				state = 2
				buf = buf + string(nextRune)
			}
		case 1:
			if nextRune == '\'' {
				state = 3
			} else {
				buf = buf + string(nextRune)
			}
		case 2:
			if nextRune == ',' {
				state = 0
				result = append(result, buf)
				buf = ""
			} else if nextRune == ' ' {
				state = 3
			} else {
				buf = buf + string(nextRune)
			}
		case 3:
			if nextRune == ',' {
				state = 0
				result = append(result, buf)
				buf = ""
			} else if nextRune == ' ' {
			} else {
				return nil, errors.New("parse error 2")
			}
		}
	}

	if state == 2 || state == 3 {
		result = append(result, buf)
		buf = ""
		state = 0
	}

	if state == 0 {
		return result, nil
	} else {
		return nil, errors.New("parse error: " + strconv.Itoa(state))
	}
}
