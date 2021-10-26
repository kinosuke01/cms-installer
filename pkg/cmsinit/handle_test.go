package cmsinit

import (
	"fmt"
	"strings"
	"testing"
)

func testError(err error, errorKeywords *[]string) string {
	exists := (err != nil)
	expectedExists := (len(*errorKeywords) > 0)
	if expectedExists != exists {
		return fmt.Sprintf("error exists wrong. want=%+v, got=%+v", expectedExists, exists)
	}
	if err != nil {
		for _, keyword := range *errorKeywords {
			if !strings.Contains(err.Error(), keyword) {
				return fmt.Sprintf("error messages wrong. want_keywords=%+v, got=%+v", keyword, err.Error())
			}
		}
	}
	return ""
}

func TestHandle(t *testing.T) {
	tt := []struct {
		name       string
		statusCode int
		bodyString string

		expectedErrorKeywords []string
	}{
		{
			name:       "status_code_500",
			statusCode: 500,
			bodyString: "",

			expectedErrorKeywords: []string{"status code is 500"},
		},
		{
			name:       "empty_body",
			statusCode: 200,
			bodyString: "",

			expectedErrorKeywords: []string{"unexpected end of JSON input"},
		},
		{
			name:       "error_message_exists",
			statusCode: 200,
			bodyString: `{"result":false,"error_message":"AUTH_ERROR"}`,

			expectedErrorKeywords: []string{"AUTH_ERROR"},
		},
		{
			name:       "result_false",
			statusCode: 200,
			bodyString: `{"result":false,"error_message":""}`,

			expectedErrorKeywords: []string{"result is false"},
		},
		{
			name:       "result_true",
			statusCode: 200,
			bodyString: `{"result":true,"error_message":""}`,

			expectedErrorKeywords: []string{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := Handle(tc.statusCode, []byte(tc.bodyString))
			msg := testError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
		})
	}
}
