package randstr

import (
	"fmt"
	"sort"
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

func TestGenerate(t *testing.T) {
	tt := []struct {
		name         string
		size         int
		count        int
		expectedSize int

		expectedErrorKeywords []string
	}{
		{
			name:         "exec_10_times",
			size:         64,
			count:        10,
			expectedSize: 64,

			expectedErrorKeywords: []string{},
		},
		{
			name:         "size_is_0",
			size:         0,
			count:        1,
			expectedSize: 0,

			expectedErrorKeywords: []string{"invalid size"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			strs := []string{}
			for i := 0; i < tc.count; i++ {
				str, err := Generate(tc.size)
				msg := testError(err, &tc.expectedErrorKeywords)
				if msg != "" {
					t.Fatalf(msg)
				}
				if tc.expectedSize != len(str) {
					t.Fatalf("size wrong. want=%+v, got=%+v", tc.expectedSize, len(str))
				}
				strs = append(strs, str)
			}

			sort.Slice(strs, func(i, j int) bool {
				return strs[i] < strs[j]
			})
			for i := 0; i < (len(strs) - 1); i++ {
				if strs[i] == strs[i+1] {
					t.Fatalf("same token generated. got=%+v", strs[i])
				}
			}
		})
	}
}
