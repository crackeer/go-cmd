package util

import (
	"strings"
)

// SplitArgs
//  @param args
//  @return []string
func SplitArgs(args []string) []string {
	retData := []string{}
	for _, value := range args {
		retData = append(retData, strings.Split(strings.TrimSpace(value), ",")...)
	}
	return retData
}
