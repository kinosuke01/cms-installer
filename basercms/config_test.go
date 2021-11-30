package basercms

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

func TestConfig_validate(t *testing.T) {
	genValuedConfig := func(f func(cnf *Config)) *Config {
		cnf := &Config{
			FtpLoginID:   "ftp_user",
			FtpPassword:  "ftp_password",
			FtpHost:      "ftp_host",
			DBType:       "mysql",
			DBName:       "db_name",
			DBUser:       "db_user",
			DBPassword:   "db_password",
			DBHost:       "db_host",
			DBPort:       "3306",
			SiteURL:      "https://example.com/blog",
			SiteUser:     "site_admin",
			SitePassword: "site_password",
			SiteEmail:    "site-admin@example.com",
		}
		f(cnf)
		return cnf
	}

	tt := []struct {
		name string
		cnf  *Config

		expectedErrorKeywords []string
	}{
		{
			name: "field_blank",
			cnf:  &Config{},

			expectedErrorKeywords: []string{
				"required",
				"ftp_login_id",
				"ftp_password",
				"ftp_host",
				"db_type",
				"site_url",
				"site_user",
				"site_password",
				"site_email",
			},
		},
		{
			name: "invalid_email",
			cnf: genValuedConfig(func(cnf *Config) {
				cnf.SiteEmail = "invalid-email-address"
			}),

			expectedErrorKeywords: []string{
				"invalid format",
				"site_email",
			},
		},
		{
			name: "invalid_url",
			cnf: genValuedConfig(func(cnf *Config) {
				cnf.SiteURL = "invalid-site-url"
			}),

			expectedErrorKeywords: []string{
				"invalid format",
				"site_url",
			},
		},
		{
			name: "invalid_db_type",
			cnf: genValuedConfig(func(cnf *Config) {
				cnf.DBType = "invalid-db-type"
			}),

			expectedErrorKeywords: []string{
				"invalid format",
				"db_type",
			},
		},
		{
			name: "empty-mysql-info",
			cnf: genValuedConfig(func(cnf *Config) {
				cnf.DBType = "mysql"
				cnf.DBHost = ""
				cnf.DBPort = ""
				cnf.DBName = ""
				cnf.DBUser = ""
				cnf.DBPassword = ""
			}),
			expectedErrorKeywords: []string{
				"required",
				"db_host",
				"db_port",
				"db_name",
				"db_user",
				"db_password",
			},
		},
		{
			name: "valid_config_when_using_mysql",
			cnf:  genValuedConfig(func(cnf *Config) {}),

			expectedErrorKeywords: []string{},
		},
		{
			name: "valid_config_when_using_sqlite",
			cnf: genValuedConfig(func(cnf *Config) {
				cnf.DBType = "sqlite"
				cnf.DBHost = ""
				cnf.DBPort = ""
				cnf.DBName = ""
				cnf.DBUser = ""
				cnf.DBPassword = ""
			}),

			expectedErrorKeywords: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cnf.validate()

			msg := testError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
		})
	}
}