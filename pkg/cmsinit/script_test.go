package cmsinit

import (
	"strings"
	"testing"
	"time"
)

func TestScript(t *testing.T) {
	tt := []struct {
		name string
		cnf  *Config
		now  string

		expectedErrorExists bool
		expectedContain     []string
		expectedExclusion   []string
	}{
		{
			name: "invalid_config",
			cnf:  &Config{},
			now:  "2021-10-19T19:30:00+09:00",

			expectedErrorExists: true,
			expectedContain:     []string{},
			expectedExclusion:   []string{},
		},
		{
			name: "valid_config",
			cnf: &Config{
				URL:    "http://example.com",
				Dir:    "cms_directory",
				Token:  "dummy_token_text",
				Period: 3600,
			},
			now: "2021-10-19T19:30:00+09:00",

			expectedErrorExists: false,
			expectedContain: []string{
				"http://example.com",
				"cms_directory",
				"dummy_token_text",
				"1634643000", // 2021-10-19T20:30:00+09:00
			},
			expectedExclusion: []string{
				"ARCHIVE_URL_PLACEHOLDER",
				"EXTRACTED_DIR_PLACEHOLDER",
				"TOKEN_PLACEHOLDER",
				"EXPIRED_AT_PLACEHOLDER",
			},
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

			str, err := Script(tc.cnf)

			if tc.expectedErrorExists != (err != nil) {
				t.Fatalf("error exists wrong. want=%+v, got=%+v", tc.expectedErrorExists, (err != nil))
			}

			for _, keyword := range tc.expectedContain {
				if !strings.Contains(str, keyword) {
					t.Fatalf("script result wrong. want_keyword=%+v", keyword)
				}
			}
			for _, keyword := range tc.expectedExclusion {
				if strings.Contains(str, keyword) {
					t.Fatalf("script result wrong. want_exclusion_keyword=%+v", keyword)
				}
			}
		})
	}
}
