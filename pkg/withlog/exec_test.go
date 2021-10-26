package withlog

import (
	"errors"
	"sort"
	"strings"
	"testing"
	"time"
)

type TestStruct struct {
	executedFlag   bool
	errorMessage   string
	processingTime int
}

func (s *TestStruct) exec() error {
	s.executedFlag = true
	time.Sleep(time.Duration(s.processingTime) * time.Second)

	if s.errorMessage != "" {
		return errors.New(s.errorMessage)
	}
	return nil
}

type TestLogger struct {
	infoList  []string
	errorList []string
}

func (l *TestLogger) Info(str string) {
	if l.infoList == nil {
		l.infoList = []string{}
	}
	l.infoList = append(l.infoList, str)
}

func (l *TestLogger) Error(str string) {
	if l.errorList == nil {
		l.errorList = []string{}
	}
	l.errorList = append(l.errorList, str)
}

func TestExec(t *testing.T) {
	tt := []struct {
		name           string
		funcName       string
		processingTime int
		errorMessage   string
		useLogger      bool

		expectedExecutedFlag  bool
		expectedInfoKeywords  []string
		expectedErrorKeywords []string
	}{
		{
			name:           "logger_is_nil",
			funcName:       "exec",
			processingTime: 1,
			useLogger:      false,

			expectedExecutedFlag:  true,
			expectedInfoKeywords:  []string{},
			expectedErrorKeywords: []string{},
		},
		{
			name:           "result_is_error",
			funcName:       "exec",
			processingTime: 0,
			useLogger:      true,

			expectedExecutedFlag: true,
			expectedInfoKeywords: []string{
				"exec",
				"Running",
				"Finished",
				"0",
				"successful",
			},
			expectedErrorKeywords: []string{},
		},
		{
			name:           "result_is_success",
			funcName:       "exec",
			processingTime: 0,
			errorMessage:   "test_error",
			useLogger:      true,

			expectedExecutedFlag: true,
			expectedInfoKeywords: []string{
				"exec",
				"Running",
				"Finished",
				"0",
				"failed",
			},
			expectedErrorKeywords: []string{
				"exec",
				"test_error",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := &TestStruct{
				errorMessage:   tc.errorMessage,
				processingTime: tc.processingTime,
			}
			var logger *TestLogger
			if tc.useLogger {
				logger = &TestLogger{
					infoList:  []string{},
					errorList: []string{},
				}
			}

			Exec(ts.exec, tc.funcName, logger)

			if tc.expectedExecutedFlag != ts.executedFlag {
				t.Fatalf(
					"func executed wrong. want=%+v got=%+v",
					tc.expectedExecutedFlag,
					ts.executedFlag,
				)
			}
			if logger != nil {
				infoText := strings.Join(logger.infoList, "\t")
				for _, keyword := range tc.expectedInfoKeywords {
					if !strings.Contains(infoText, keyword) {
						t.Fatalf("logger info wrong. want_keyword=%+v, got=%+v", keyword, infoText)
					}
				}

				errorText := strings.Join(logger.errorList, "\t")
				for _, keyword := range tc.expectedErrorKeywords {
					if !strings.Contains(errorText, keyword) {
						t.Fatalf("logger error wrong. want_keyword=%+v, got=%+v", keyword, errorText)
					}
				}
			}
		})
	}
}

func TestGenHash(t *testing.T) {
	tt := []struct {
		name      string
		funcNames []string
	}{
		{
			name: "exec",
			funcNames: []string{
				"step01",
				"step02",
				"step03",
				"step04",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hashList := []string{}
			for _, fn := range tc.funcNames {
				hashList = append(hashList, genHash(fn))
			}
			sort.Slice(hashList, func(i, j int) bool {
				return hashList[i] < hashList[j]
			})
			for i := 0; i < (len(hashList) - 1); i++ {
				if hashList[i] == hashList[i+1] {
					t.Fatalf("same hash generated. got=%+v", hashList[i])
				}
			}
		})
	}
}

func TestPeriod(t *testing.T) {
	tt := []struct {
		name   string
		before int64
		after  int64

		expectedResult float32
	}{
		{
			name:   "exec",
			before: 1603592280, // 2020-10-25T11:18:00+09:00
			after:  1603592282,

			expectedResult: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			startAt := time.Unix(tc.before, 0)
			endAt := time.Unix(tc.after, 0)
			result := period(&startAt, &endAt)

			if tc.expectedResult != result {
				t.Fatalf("result wrong. want=%+v, got=%+v", tc.expectedResult, result)
			}
		})
	}
}
