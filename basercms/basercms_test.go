package basercms

import (
	"reflect"
	"testing"

	"github.com/kinosuke01/cms-installer/pkg/thelper"
)

func TestNew(t *testing.T) {
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
		expectedDbPrefix   string

		expectedSiteTitle    string
		expectedSiteUser     string
		expectedSitePassword string
		expectedSiteEmail    string

		expectedInitScript string
		expectedArchiveURL string

		expectedErrorExists   bool
		expectedErrorKeywords []string
	}{
		{
			name: "invalid_config",
			cnf:  &Config{},

			expectedErrorExists: true,
			expectedErrorKeywords: []string{
				"required",
			},
		},
		{
			name: "valid_config",
			cnf: &Config{
				FtpLoginID:   "ftpuser",
				FtpPassword:  "ftppassword",
				FtpHost:      "ftp.example.com",
				FtpDir:       "example.com/blog",
				DBType:       "mysql",
				DBName:       "site_production",
				DBUser:       "dbuser",
				DBPassword:   "dbpassword",
				DBHost:       "db.example.com",
				DBPrefix:     "mysite_",
				SiteURL:      "https://example.com/blog",
				SiteUser:     "siteadmin",
				SitePassword: "sitepassword",
				SiteEmail:    "site@example.com",
			},

			expectedFtpcType:     "*ftpc.Client",
			expectedHttpcType:    "*httpc.Client",
			expectedFtpDir:       "example.com/blog",
			expectedDbType:       "mysql",
			expectedDbName:       "site_production",
			expectedDbUser:       "dbuser",
			expectedDbPassword:   "dbpassword",
			expectedDbHost:       "db.example.com",
			expectedDbPrefix:     "mysite_",
			expectedSiteTitle:    "MyBlog",
			expectedSiteUser:     "siteadmin",
			expectedSitePassword: "sitepassword",
			expectedSiteEmail:    "site@example.com",
			expectedInitScript:   "init.php",
			expectedArchiveURL:   "https://ja.wordpress.org/latest-ja.zip",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cms, err := New(tc.cnf)

			msg := thelper.CheckInclusion(err, &tc.expectedErrorKeywords)
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
				{"dbName", tc.expectedDbName, cms.dbName},
				{"dbUser", tc.expectedDbUser, cms.dbUser},
				{"dbPassword", tc.expectedDbPassword, cms.dbPassword},
				{"dbHost", tc.expectedDbHost, cms.dbHost},
				{"dbPrefix", tc.expectedDbPrefix, cms.dbPrefix},
				{"siteUser", tc.expectedSiteUser, cms.siteUser},
				{"sitePassword", tc.expectedSitePassword, cms.sitePassword},
				{"siteEmail", tc.expectedSiteEmail, cms.siteEmail},
				{"initScript", tc.expectedInitScript, cms.initScript},
				{"archiveURL", tc.expectedArchiveURL, cms.initArchiveURL},
			}
			for _, kvv := range kvvs {
				if kvv[1] != kvv[2] {
					t.Fatalf("%+v wrong. want=%+v, got=%+v", kvv[0], kvv[1], kvv[2])
				}
			}
		})
	}
}

func TestWordpress_InjectInitScript(t *testing.T) {
	tt := []struct {
		name string

		ftpDir       string
		initScript   string
		initToken    string
		errorMessage string

		expectedFilePath      string
		expectedErrorKeywords []string
	}{
		{
			name:                  "success",
			ftpDir:                "example.com/blog",
			initScript:            "init.php",
			initToken:             "xxxxx",
			errorMessage:          "",
			expectedFilePath:      "example.com/blog/init.php",
			expectedErrorKeywords: []string{},
		},
		{
			name:                  "error",
			ftpDir:                "example.com/blog",
			initScript:            "init.php",
			initToken:             "xxxxx",
			errorMessage:          "error",
			expectedFilePath:      "",
			expectedErrorKeywords: []string{"error"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ftpc := &TestFtpClient{
				errorMessage: tc.errorMessage,
			}
			cms := &WordPress{
				ftpc:       ftpc,
				ftpDir:     tc.ftpDir,
				initScript: tc.initScript,
				initToken:  tc.initToken,
			}
			err := cms.InjectInitScript()

			msg := testError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
			if err != nil {
				return
			}

			if tc.expectedFilePath != ftpc.filePath {
				t.Fatalf("ftpc.filePath wrong. want=%+v, got=%+v", tc.expectedFilePath, ftpc.filePath)
			}

			content, err := cms.InitScript()
			if err != nil {
				t.Fatalf(err.Error())
			}
			if *content != ftpc.content {
				t.Fatalf("ftpc.content wrong.")
			}
		})
	}
}

func TestWordpress_ExecInit(t *testing.T) {
	tt := []struct {
		name         string
		statusCode   int
		bodyString   string
		errorMessage string

		expectedErrorKeywords []string
	}{
		{
			name:         "error_exists",
			statusCode:   200,
			bodyString:   "",
			errorMessage: "timeout error",

			expectedErrorKeywords: []string{"timeout error"},
		},
		{
			name:         "invalid_status_code",
			statusCode:   500,
			bodyString:   "",
			errorMessage: "",

			expectedErrorKeywords: []string{"status code"},
		},
		{
			name:         "valid_response",
			statusCode:   200,
			bodyString:   `{"result":true,"error_message":""}`,
			errorMessage: "",

			expectedErrorKeywords: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			httpClient := &TestHttpClient{
				statusCode:   tc.statusCode,
				bodyString:   tc.bodyString,
				errorMessage: tc.errorMessage,
			}
			cms := &WordPress{
				httpc: httpClient,
			}

			err := cms.ExecInit()
			msg := testError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
		})
	}
}

func TestWordpress_DeleteInitScript(t *testing.T) {
	tt := []struct {
		name string

		ftpDir       string
		initScript   string
		errorMessage string

		expectedFilePath      string
		expectedErrorKeywords []string
	}{
		{
			name:                  "success",
			ftpDir:                "example.com/blog",
			initScript:            "init.php",
			errorMessage:          "",
			expectedFilePath:      "example.com/blog/init.php",
			expectedErrorKeywords: []string{},
		},
		{
			name:                  "error",
			ftpDir:                "example.com/blog",
			initScript:            "init.php",
			errorMessage:          "error",
			expectedFilePath:      "",
			expectedErrorKeywords: []string{"error"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ftpc := &TestFtpClient{
				errorMessage: tc.errorMessage,
			}
			cms := &WordPress{
				ftpc:       ftpc,
				ftpDir:     tc.ftpDir,
				initScript: tc.initScript,
			}
			err := cms.DeleteInitScript()

			msg := testError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
			if err != nil {
				return
			}

			if tc.expectedFilePath != ftpc.filePath {
				t.Fatalf("ftpc.filePath wrong. want=%+v, got=%+v", tc.expectedFilePath, ftpc.filePath)
			}
		})
	}
}

func TestWordpress_WpAdminSetupConfig(t *testing.T) {
	tt := []struct {
		name         string
		statusCode   int
		errorMessage string

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
			name:       "status_code_200",
			statusCode: 200,

			expectedErrorKeywords: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			httpClient := &TestHttpClient{
				statusCode:   tc.statusCode,
				errorMessage: tc.errorMessage,
			}
			cms := &WordPress{
				httpc: httpClient,
			}
			err := cms.WpAdminSetupConfig()
			msg := testError(err, &tc.expectedErrorKeywords)
			if msg != "" {
				t.Fatalf(msg)
			}
		})
	}
}

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
