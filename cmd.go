package main

import (
	"errors"
	"fmt"
	"strings"
)

type redisCommand struct{
	command string
	argv 	[]string
}

var cmd redisCommand

func execCommand(cmdParams any) (string, error) {
	arr, ok := cmdParams.([]any)
	if !ok {
		return "", errors.New("expected array command")
	}

	if len(arr) == 0 {
		return "", errors.New("empty command")
	}

	strParams := make([]string, len(arr))
	for i, v := range arr {
		s, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("non-string element in command array: %v", v)
		}
		strParams[i] = s
	}

	command := strings.ToUpper(strParams[0])
	argv := strParams[1:]

	switch command {
	case "PING":
		return execPing(argv)
	default:
		return "", fmt.Errorf("unknown command '%s'", command)
	}
}

func execPing(argv []string)(string, error){
	if(len(argv) == 0){
		return ("PONG"), nil
	}
	if(len(argv) == 1){
		return argv[0], nil
	}
	return "", errors.New("ERR Number of arguments is more than expected")
}