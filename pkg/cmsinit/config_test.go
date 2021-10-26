package cmsinit

import (
	"strings"
	"testing"
	"time"
)

func TestConfig_init(t *testing.T) {
	tt := []struct {
		name string
		cnf  *Config

		expectedURL       string
		expectedDir       string
		expectedPeriod    int
		expectedNowExists bool
	}{
		{
			name: "manual_setting",
			cnf: &Config{
				URL:    "http://example.com",
				Dir:    "cms",
				Period: 86400,
			},
			expectedURL:       "http://example.com",
			expectedDir:       "cms",
			expectedPeriod:    86400,
			expectedNowExists: true,
		},
		{
			name:              "default_setting",
			cnf:               &Config{},
			expectedURL:       defaultURL,
			expectedDir:       defaultDir,
			expectedPeriod:    defaultPeriod,
			expectedNowExists: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.cnf.init()
			if tc.expectedURL != tc.cnf.URL {
				t.Fatalf("URL wrong. want=%+v, got=%+v", tc.expectedURL, tc.cnf.URL)
			}
			if tc.expectedDir != tc.cnf.Dir {
				t.Fatalf("Dir wrong. want=%+v, got=%+v", tc.expectedDir, tc.cnf.Dir)
			}
			if tc.expectedPeriod != tc.cnf.Period {
				t.Fatalf("Period wrong. want=%+v, got=%+v", tc.expectedPeriod, tc.cnf.Period)
			}
			if tc.expectedNowExists != (tc.cnf.Now != nil) {
				t.Fatalf(
					"ExpiredAtExists wrong. want=%+v, got=%+v",
					tc.expectedNowExists,
					(tc.cnf.Now != nil),
				)
			}
		})
	}
}

func TestConfig_validate(t *testing.T) {
	tt := []struct {
		name                  string
		cnf                   *Config
		expectedErrorExists   bool
		expectedErrorMessages []string
	}{
		{
			name: "valid_token",
			cnf: &Config{
				Token: "xxxxx",
			},
			expectedErrorExists: false,
			expectedErrorMessages: []string{
				"token",
				"required",
			},
		},
		{
			name:                "empty_token",
			cnf:                 &Config{},
			expectedErrorExists: true,
			expectedErrorMessages: []string{
				"token",
				"required",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cnf.validate()

			if tc.expectedErrorExists != (err != nil) {
				t.Fatalf("error exists wrong. want=%+v, got=%+v", tc.expectedErrorExists, (err != nil))
			}
			if err != nil {
				for _, keyword := range tc.expectedErrorMessages {
					if !strings.Contains(err.Error(), keyword) {
						t.Fatalf("error messages wrong. want_keywords=%+v, got=%+v", tc.expectedErrorMessages, err.Error())
					}
				}
			}
		})
	}
}

func TestConfig_expiredAt(t *testing.T) {
	tt := []struct {
		name string
		cnf  *Config
		now  string

		expectedExpiredAtUt int64
	}{
		{
			name: "get_expired_at",
			cnf: &Config{
				Period: 3620,
			},
			now: "2021-10-19T19:30:00+09:00",

			expectedExpiredAtUt: 1634643020, // 2021-10-19T20:30:20+09:00
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.cnf.init()

			iso8601 := "2006-01-02T15:04:05-07:00"
			tm, err := time.Parse(iso8601, tc.now)
			if err != nil {
				t.Fatalf(err.Error())
			}
			tm = tm.Local()
			tc.cnf.Now = &tm

			expiredAt := tc.cnf.expiredAt()
			expiredAtUt := expiredAt.Unix()

			if tc.expectedExpiredAtUt != expiredAtUt {
				t.Fatalf("expiredAt wrong. want=%+v, got=%+v", tc.expectedExpiredAtUt, expiredAtUt)
			}
		})
	}
}
