package wordpress

import (
	"testing"

	"github.com/kinosuke01/cms-installer/pkg/testhelper"
)

func TestConfig_validate(t *testing.T) {
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
				"db_name",
				"db_user",
				"db_password",
				"db_host",
				"site_url",
				"site_title",
				"site_user",
				"site_password",
				"site_email",
			},
		},
		{
			name: "invalid_email",
			cnf: &Config{
				FtpLoginID:   "ftp_user",
				FtpPassword:  "ftp_password",
				FtpHost:      "ftp_host",
				DBName:       "db_name",
				DBUser:       "db_user",
				DBPassword:   "db_password",
				DBHost:       "db_host",
				SiteURL:      "https://example.com/blog",
				SiteTitle:    "Hello WordPress!",
				SiteUser:     "wp_admin",
				SitePassword: "wp_password",
				SiteEmail:    "invalid-email-address",
			},

			expectedErrorKeywords: []string{
				"invalid format",
				"site_email",
			},
		},
		{
			name: "invalid_url",
			cnf: &Config{
				FtpLoginID:   "ftp_user",
				FtpPassword:  "ftp_password",
				FtpHost:      "ftp_host",
				DBName:       "db_name",
				DBUser:       "db_user",
				DBPassword:   "db_password",
				DBHost:       "db_host",
				SiteURL:      "invaid-site-url",
				SiteTitle:    "Hello WordPress!",
				SiteUser:     "wp_admin",
				SitePassword: "wp_password",
				SiteEmail:    "wp-admin@example.com",
			},

			expectedErrorKeywords: []string{
				"invalid format",
				"site_url",
			},
		},
		{
			name: "valid_config",
			cnf: &Config{
				FtpLoginID:   "ftp_user",
				FtpPassword:  "ftp_password",
				FtpHost:      "ftp_host",
				DBName:       "db_name",
				DBUser:       "db_user",
				DBPassword:   "db_password",
				DBHost:       "db_host",
				SiteURL:      "https://example.com/blog",
				SiteTitle:    "Hello WordPress!",
				SiteUser:     "wp_admin",
				SitePassword: "wp_password",
				SiteEmail:    "wp-admin@example.com",
			},

			expectedErrorKeywords: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cnf.validate()

			msg := testhelper.CheckError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
		})
	}
}
