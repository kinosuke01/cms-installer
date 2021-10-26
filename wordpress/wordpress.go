package wordpress

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/kinosuke01/cms-installer/pkg/cmsinit"
	"github.com/kinosuke01/cms-installer/pkg/ftpc"
	"github.com/kinosuke01/cms-installer/pkg/httpc"
	"github.com/kinosuke01/cms-installer/pkg/withlog"
)

type FtpClient interface {
	Upload(string, *string) error
	Delete(string, bool) error
}

type HttpClient interface {
	DoRequest(context.Context, *httpc.RequestOptions) (*httpc.Response, error)
}

type Logger interface {
	Info(string)
	Error(string)
}

type WordPress struct {
	ftpc  FtpClient
	httpc HttpClient

	ftpDir string

	dbName     string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPrefix   string

	siteURL      string
	siteTitle    string
	siteUser     string
	sitePassword string
	siteEmail    string
	sitePublic   bool

	initScript     string
	initToken      string
	initArchiveURL string
	initArchiveDir string

	Logger Logger
}

func New(cnf *Config) (*WordPress, error) {
	err := cnf.init()
	if err != nil {
		return nil, err
	}

	err = cnf.validate()
	if err != nil {
		return nil, err
	}

	fc := ftpc.New(&ftpc.Config{
		LoginID:  cnf.FtpLoginID,
		Password: cnf.FtpPassword,
		Host:     cnf.FtpHost,
		Port:     cnf.FtpPort,
	})

	baseURL, err := url.Parse(cnf.SiteURL)
	if err != nil {
		return nil, err
	}
	hc := httpc.New(&httpc.Config{
		Scheme:   baseURL.Scheme,
		Host:     baseURL.Host,
		BasePath: baseURL.Path,
	})

	return &WordPress{
		ftpc:  fc,
		httpc: hc,

		ftpDir: cnf.FtpDir,

		dbName:     cnf.DBName,
		dbUser:     cnf.DBUser,
		dbPassword: cnf.DBPassword,
		dbHost:     cnf.DBHost,
		dbPrefix:   cnf.DBPrefix,

		siteURL:      cnf.SiteURL,
		siteTitle:    cnf.SiteTitle,
		siteUser:     cnf.SiteUser,
		sitePassword: cnf.SitePassword,
		siteEmail:    cnf.SiteEmail,
		sitePublic:   !cnf.SiteUnpublic,

		initScript:     initScript,
		initToken:      cnf.InitToken,
		initArchiveURL: cnf.InitArchiveURL,
		initArchiveDir: cnf.InitArchiveDir,
	}, nil
}

func (cms *WordPress) InitScript() (*string, error) {
	str, err := cmsinit.Script(&cmsinit.Config{
		URL:   cms.initArchiveURL,
		Dir:   cms.initArchiveDir,
		Token: cms.initToken,
	})

	if err != nil {
		return nil, err
	}

	return &str, nil
}

func (cms *WordPress) InjectInitScript() error {
	pContent, err := cms.InitScript()
	if err != nil {
		return err
	}
	filePath := path.Join(cms.ftpDir, cms.initScript)

	// TODO Add injection methods other than FTP.
	err = cms.ftpc.Upload(filePath, pContent)
	if err != nil {
		return err
	}
	return nil
}

func (cms *WordPress) ExecInit() error {
	ctx, cancel := cms.httpContext(initHttpTimeout)
	defer cancel()

	res, err := cms.httpc.DoRequest(
		ctx,
		&httpc.RequestOptions{
			Path:   cms.initScript,
			Method: http.MethodPost,
			BodyValues: url.Values{
				"token": []string{cms.initToken},
			},
		},
	)

	if err != nil {
		return err
	}

	err = cmsinit.Handle(res.StatusCode, res.BodyBytes)
	if err != nil {
		return err
	}

	return nil
}

func (cms *WordPress) DeleteInitScript() error {
	filePath := path.Join(cms.ftpDir, cms.initScript)

	err := cms.ftpc.Delete(filePath, false)
	if err != nil {
		return err
	}

	return nil
}

func (cms *WordPress) DeleteWpConfig() error {
	filePath := path.Join(cms.ftpDir, "wp-config.php")

	err := cms.ftpc.Delete(filePath, true)
	if err != nil {
		return err
	}

	return nil
}

func (cms *WordPress) WpAdminSetupConfig() error {
	ctx, cancel := cms.httpContext(wpHttpTimeout)
	defer cancel()

	res, err := cms.httpc.DoRequest(
		ctx,
		&httpc.RequestOptions{
			Method: http.MethodPost,
			Path:   "wp-admin/setup-config.php",
			Queries: map[string]string{
				"step": "2",
			},
			BodyValues: url.Values{
				"dbname": []string{cms.dbName},
				"uname":  []string{cms.dbUser},
				"pwd":    []string{cms.dbPassword},
				"dbhost": []string{cms.dbHost},
				"prefix": []string{cms.dbPrefix},
			},
		},
	)

	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("status code is %+v", res.StatusCode)
	}

	return nil
}

func (cms *WordPress) WpAdminInstall() error {
	pubval := "0"
	if cms.sitePublic {
		pubval = "1"
	}

	ctx, cancel := cms.httpContext(wpHttpTimeout)
	defer cancel()

	res, err := cms.httpc.DoRequest(
		ctx,
		&httpc.RequestOptions{
			Method: http.MethodPost,
			Path:   "wp-admin/install.php",
			Queries: map[string]string{
				"step": "2",
			},
			BodyValues: url.Values{
				"weblog_title":    []string{cms.siteTitle},
				"user_name":       []string{cms.siteUser},
				"admin_password":  []string{cms.sitePassword},
				"admin_password2": []string{cms.sitePassword},
				"admin_email":     []string{cms.siteEmail},
				"blog_public":     []string{pubval},
			},
		},
	)

	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("status code is %+v", res.StatusCode)
	}

	// validation error handling
	// Cannot determine from status code.
	// So scraping required.
	bodyStr := string(res.BodyBytes)
	ngList := []string{
		`name="weblog_title"`,
		`name="user_name"`,
		`name="admin_password"`,
		`name="admin_email"`,
	}
	isValid := false
	for _, ngStr := range ngList {
		if !strings.Contains(bodyStr, ngStr) {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid post parameter")
	}

	return nil
}

func (cms *WordPress) Install() error {
	var err error

	desc := func(fn string) string {
		return fmt.Sprintf("%+v - %+v", cms.siteURL, fn)
	}

	err = withlog.Exec(cms.InjectInitScript, desc("InjectInitScript"), cms.Logger)
	if err != nil {
		return err
	}

	err = withlog.Exec(cms.ExecInit, desc("ExecInit"), cms.Logger)
	if err != nil {
		return err
	}

	err = withlog.Exec(cms.DeleteInitScript, desc("DeleteInitScript"), cms.Logger)
	if err != nil {
		return err
	}

	err = withlog.Exec(cms.DeleteWpConfig, desc("DeleteWpConfig"), cms.Logger)
	if err != nil {
		return err
	}

	err = withlog.Exec(cms.WpAdminSetupConfig, desc("WpAdminSetupConfig"), cms.Logger)
	if err != nil {
		return err
	}

	err = withlog.Exec(cms.WpAdminInstall, desc("WpAdminInstall"), cms.Logger)
	if err != nil {
		return err
	}

	return nil
}

func (cms *WordPress) httpContext(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}
