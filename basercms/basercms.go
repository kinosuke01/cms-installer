package basercms

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
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

type BaserCMS struct {
	ftpc  FtpClient
	httpc HttpClient

	ftpDir string

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

	return &BaserCMS{
		ftpc:  fc,
		httpc: hc,

		ftpDir: cnf.FtpDir,

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

	return nil
}

func (cms *BaserCMS) httpContext(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}
