package testhelper

import (
	"fmt"
	"strings"
)

func CheckError(err error, keywords *[]string) string {
	exists := (err != nil)
	expectedExists := (len(*keywords) > 0)
	if expectedExists && !exists {
		return fmt.Sprintf("error exists wrong. want=%+v, got=%+v", expectedExists, exists)
	}
	if !expectedExists && exists {
		return fmt.Sprintf("error wrong. want=%+v, got=%+v", expectedExists, err.Error())
	}
	if err != nil {
		for _, keyword := range *keywords {
			if !strings.Contains(err.Error(), keyword) {
				return fmt.Sprintf("error messages wrong. want_keywords=%+v, got=%+v", keyword, err.Error())
			}
		}
	}
	return ""
}
