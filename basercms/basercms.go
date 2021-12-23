package basercms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/kinosuke01/cms-installer/pkg/cmsinit"
	"github.com/kinosuke01/cms-installer/pkg/ftpc"
	"github.com/kinosuke01/cms-installer/pkg/httpc"
	"github.com/kinosuke01/cms-installer/pkg/randstr"
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

type BaserCMS struct {
	ftpc  FtpClient
	httpc HttpClient

	ftpDir string

	dbType     string
	dbName     string
	dbUser     string
	dbPassword string
	dbHost     string
	dbPort     string
	dbPrefix   string

	siteURL      string
	siteUser     string
	sitePassword string
	siteEmail    string

	initScript     string
	initToken      string
	initArchiveURL string
	initArchiveDir string

	bcInstallScript  string
	bcInstallToken   string
	bcInstallPHPPath string

	Logger Logger
}

func New(cnf *Config) (*BaserCMS, error) {
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

	bcInstallToken, err := randstr.Generate(64)
	if err != nil {
		return nil, err
	}

	return &BaserCMS{
		ftpc:  fc,
		httpc: hc,

		ftpDir: cnf.FtpDir,

		dbType:     cnf.DBType,
		dbName:     cnf.DBName,
		dbUser:     cnf.DBUser,
		dbPassword: cnf.DBPassword,
		dbHost:     cnf.DBHost,
		dbPort:     cnf.DBPort,
		dbPrefix:   cnf.DBPrefix,

		siteURL:      cnf.SiteURL,
		siteUser:     cnf.SiteUser,
		sitePassword: cnf.SitePassword,
		siteEmail:    cnf.SiteEmail,

		initScript:     initScript,
		initToken:      cnf.InitToken,
		initArchiveURL: cnf.InitArchiveURL,
		initArchiveDir: cnf.InitArchiveDir,

		bcInstallScript:  bcInstallScript,
		bcInstallToken:   bcInstallToken,
		bcInstallPHPPath: cnf.PHPPath,
	}, nil
}

func (cms *BaserCMS) InitScript() (*string, error) {
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

func (cms *BaserCMS) InjectInitScript() error {
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

func (cms *BaserCMS) ExecInit() error {
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

func (cms *BaserCMS) DeleteInitScript() error {
	filePath := path.Join(cms.ftpDir, cms.initScript)

	err := cms.ftpc.Delete(filePath, false)
	if err != nil {
		return err
	}

	return nil
}

func (cms *BaserCMS) BcInstallScript(now time.Time) (*string, error) {
	str := bcInstallScriptTemplate
	str = strings.Replace(str, "TOKEN_PLACEHOLDER", cms.bcInstallToken, 1)
	str = strings.Replace(str, "PHP_PATH_PLACEHOLDER", cms.bcInstallPHPPath, 1)

	t := now.Add(time.Duration(60) * time.Second).Local()
	timestr := strconv.FormatInt(t.Unix(), 10)
	str = strings.Replace(str, "EXPIRED_AT_PLACEHOLDER", timestr, 1)

	return &str, nil
}

func (cms *BaserCMS) InjectBcInstallScript() error {
	pContent, err := cms.BcInstallScript(time.Now().Local())
	if err != nil {
		return err
	}
	filePath := path.Join(cms.ftpDir, cms.bcInstallScript)

	err = cms.ftpc.Upload(filePath, pContent)
	if err != nil {
		return err
	}
	return nil
}

func (cms *BaserCMS) BcIntall() error {
	ctx, cancel := cms.httpContext(bcHttpTimeout)
	defer cancel()

	parsedURL, err := url.Parse(cms.siteURL)
	if err != nil {
		return err
	}
	siteURL := parsedURL.Scheme + "://" + parsedURL.Host
	baseURL := parsedURL.Path

	res, err := cms.httpc.DoRequest(
		ctx,
		&httpc.RequestOptions{
			Path:   cms.bcInstallScript,
			Method: http.MethodPost,
			BodyValues: url.Values{
				"token":        []string{cms.bcInstallToken},
				"siteurl":      []string{siteURL},
				"dbtype":       []string{cms.dbType},
				"siteuser":     []string{cms.siteUser},
				"sitepassword": []string{cms.sitePassword},
				"email":        []string{cms.siteEmail},
				"host":         []string{cms.dbHost},
				"database":     []string{cms.dbName},
				"login":        []string{cms.dbUser},
				"password":     []string{cms.dbPassword},
				"prefix":       []string{cms.dbPrefix},
				"port":         []string{cms.dbPort},
				"baseurl":      []string{baseURL},
			},
		},
	)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("status code is %+v", res.StatusCode)
	}

	bodyData := struct {
		ExitCode string   `json:"exit_code"`
		Messages []string `json:"messages"`
	}{}
	err = json.Unmarshal(res.BodyBytes, &bodyData)

	if err != nil {
		return err
	}

	exitCode := bodyData.ExitCode
	msgLine := strings.Join(bodyData.Messages, "\t")

	if exitCode != "0" {
		return fmt.Errorf("exit_code is %+v, messages is %+v", bodyData.ExitCode, msgLine)
	}

	failedKeywords := []string{
		"AUTH_ERROR",
		"EXEC_ERROR",
		"既にインストール済です",
		"baserCMSのインストールを行うには",
		"baserCMSのインストールに失敗しました",
	}
	for _, keyword := range failedKeywords {
		if !strings.Contains(msgLine, keyword) {
			return fmt.Errorf("messages is %+v", msgLine)
		}
	}

	return nil
}

func (cms *BaserCMS) DeleteBcInstallScript() error {
	filePath := path.Join(cms.ftpDir, cms.bcInstallScript)

	err := cms.ftpc.Delete(filePath, false)
	if err != nil {
		return err
	}

	return nil
}

func (cms *BaserCMS) Install() error {
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

	err = withlog.Exec(cms.InjectBcInstallScript, desc("InjectBcInstallScript"), cms.Logger)
	if err != nil {
		return err
	}

	err = withlog.Exec(cms.BcIntall, desc("BcInstall"), cms.Logger)
	if err != nil {
		return err
	}

	err = withlog.Exec(cms.DeleteBcInstallScript, desc("DeleteBcInstallScript"), cms.Logger)
	if err != nil {
		return err
	}

	return nil
}

func (cms *BaserCMS) httpContext(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}
