package basercms

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/kinosuke01/cms-installer/pkg/thelper"
)

func TestNew(t *testing.T) {
	genValuedConfig := func(f func(cnf *Config)) *Config {
		cnf := &Config{
			FtpLoginID:     "ftp_user",
			FtpPassword:    "ftp_password",
			FtpHost:        "ftp_host",
			FtpDir:         "mysite",
			DBType:         "mysql",
			DBName:         "db_name",
			DBUser:         "db_user",
			DBPassword:     "db_password",
			DBHost:         "db_host",
			SiteURL:        "https://example.com/mysite",
			SiteUser:       "site_admin",
			SitePassword:   "site_password",
			SiteEmail:      "site-admin@example.com",
			InitArchiveURL: "https://basercms.net/packages/download_exec/basercms-x.x.x.zip",
		}
		f(cnf)
		return cnf
	}

	tt := []struct {
		name string
		cnf  *Config

		expectedFtpcType  string
		expectedHttpcType string

		expectedFtpDir string

		expectedDbType     string
		expectedDbName     string
		expectedDbUser     string
		expectedDbPassword string
		expectedDbHost     string
		expectedDbPort     string
		expectedDbPrefix   string

		expectedSiteUser     string
		expectedSitePassword string
		expectedSiteEmail    string

		expectedInitScript     string
		expectedInitTokenLen   int
		expectedInitArchiveURL string
		expectedInitArchiveDir string

		expectedBcInstallScript   string
		expectedBcInstallTokenLen int
		expectedBcInstallPHPPath  string

		expectedErrorKeywords []string
	}{
		{
			name: "invalid_config",
			cnf:  &Config{},

			expectedErrorKeywords: []string{
				"required",
			},
		},
		{
			name: "valid_config",
			cnf:  genValuedConfig(func(cnf *Config) {}),

			expectedFtpcType:     "*ftpc.Client",
			expectedHttpcType:    "*httpc.Client",
			expectedFtpDir:       "mysite",
			expectedDbType:       "mysql",
			expectedDbName:       "db_name",
			expectedDbUser:       "db_user",
			expectedDbPassword:   "db_password",
			expectedDbHost:       "db_host",
			expectedDbPort:       "3306",
			expectedDbPrefix:     "mysite_",
			expectedSiteUser:     "site_admin",
			expectedSitePassword: "site_password",
			expectedSiteEmail:    "site-admin@example.com",

			expectedInitScript:     initScript,
			expectedInitTokenLen:   tokenLen,
			expectedInitArchiveURL: "https://basercms.net/packages/download_exec/basercms-x.x.x.zip",
			expectedInitArchiveDir: initArchiveDir,

			expectedBcInstallScript:   bcInstallScript,
			expectedBcInstallTokenLen: tokenLen,
			expectedBcInstallPHPPath:  "php",

			expectedErrorKeywords: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cms, err := New(tc.cnf)

			msg := thelper.CheckError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
			if err != nil {
				return
			}

			kvvs := [][3]string{
				{"ftpcType", tc.expectedFtpcType, reflect.TypeOf(cms.ftpc).String()},
				{"httpcType", tc.expectedHttpcType, reflect.TypeOf(cms.httpc).String()},
				{"ftpDir", tc.expectedFtpDir, cms.ftpDir},
				{"dbType", tc.expectedDbType, cms.dbType},
				{"dbName", tc.expectedDbName, cms.dbName},
				{"dbUser", tc.expectedDbUser, cms.dbUser},
				{"dbPassword", tc.expectedDbPassword, cms.dbPassword},
				{"dbHost", tc.expectedDbHost, cms.dbHost},
				{"dbPort", tc.expectedDbPort, cms.dbPort},
				{"dbPrefix", tc.expectedDbPrefix, cms.dbPrefix},
				{"siteUser", tc.expectedSiteUser, cms.siteUser},
				{"sitePassword", tc.expectedSitePassword, cms.sitePassword},
				{"siteEmail", tc.expectedSiteEmail, cms.siteEmail},
				{"initScript", tc.expectedInitScript, cms.initScript},
				{"initArchiveURL", tc.expectedInitArchiveURL, cms.initArchiveURL},
				{"bcInstallScript", tc.expectedBcInstallScript, cms.bcInstallScript},
				{"bcInstallPHPPath", tc.expectedBcInstallPHPPath, cms.bcInstallPHPPath},
			}
			for _, kvv := range kvvs {
				if kvv[1] != kvv[2] {
					t.Fatalf("%+v wrong. want=%+v, got=%+v", kvv[0], kvv[1], kvv[2])
				}
			}

			if tc.expectedInitTokenLen != len(cms.initToken) {
				t.Fatalf("initTokenLength wrong. want=%+v, got=%+v", tc.expectedInitTokenLen, len(cms.initToken))
			}
			if tc.expectedBcInstallTokenLen != len(cms.bcInstallToken) {
				t.Fatalf("bcInstallTokenLength wrong. want=%+v, got=%+v", tc.expectedBcInstallTokenLen, len(cms.bcInstallToken))
			}
		})
	}
}

func TestBaserCMS_BcInstallScript(t *testing.T) {
	tt := []struct {
		name    string
		token   string
		phpPath string
		now     string

		expectedToken     string
		expectedPHPPath   string
		expectedExpiredAt string
	}{
		{
			name:    "execute",
			token:   "asdfghjkklqwertyuiopzxcvbnm1234567890",
			phpPath: "/usr/local/php",
			now:     "2021-10-19T20:30:00+09:00",

			expectedToken:     "asdfghjkklqwertyuiopzxcvbnm1234567890",
			expectedPHPPath:   "/usr/local/php",
			expectedExpiredAt: "1634643060", // 2021-10-19T20:31:00+09:00
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cms := &BaserCMS{
				bcInstallToken:   tc.token,
				bcInstallPHPPath: tc.phpPath,
			}

			iso8601 := "2006-01-02T15:04:05-07:00"
			tm, err := time.Parse(iso8601, tc.now)
			if err != nil {
				t.Fatalf(err.Error())
			}
			tm = tm.Local()

			result, _ := cms.BcInstallScript(tm)

			if !strings.Contains(*result, tc.expectedToken) {
				t.Fatalf("bcInstallScript wrong. want_keyword=%+v", tc.expectedToken)
			}
			if !strings.Contains(*result, tc.expectedPHPPath) {
				t.Fatalf("bcInstallScript wrong. want_keyword=%+v", tc.expectedPHPPath)
			}
			if !strings.Contains(*result, tc.expectedExpiredAt) {
				t.Fatalf("bcInstallScript wrong. want_keyword=%+v", tc.expectedExpiredAt)
			}
		})
	}
}

/*
func TestWordpress_WpAdminInstall(t *testing.T) {
	tt := []struct {
		name             string
		statusCode       int
		errorMessage     string
		responseBodyFile string

		expectedErrorKeywords []string
	}{
		{
			name:         "error_exists",
			errorMessage: "timeout error",

			expectedErrorKeywords: []string{"timeout error"},
		},
		{
			name:       "status_code_500",
			statusCode: 500,

			expectedErrorKeywords: []string{"status code is 500"},
		},
		{
			name:             "validation_error",
			statusCode:       200,
			responseBodyFile: "testdata/install/empty-error-body",

			expectedErrorKeywords: []string{
				"invalid post parameter",
			},
		},
		{
			name:             "status_code_200",
			statusCode:       200,
			responseBodyFile: "testdata/install/success-body",

			expectedErrorKeywords: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			httpClient := &TestHttpClient{
				statusCode:       tc.statusCode,
				errorMessage:     tc.errorMessage,
				responseBodyFile: tc.responseBodyFile,
			}
			cms := &WordPress{
				httpc: httpClient,
			}
			err := cms.WpAdminInstall()
			msg := testError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
		})
	}
}
*/
