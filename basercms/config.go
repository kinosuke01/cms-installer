package basercms

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"strings"

	"github.com/kinosuke01/cms-installer/pkg/randstr"
)

type Config struct {
	FtpLoginID  string `json:"ftp_login_id" desc:"[required] FTP loginID"`
	FtpPassword string `json:"ftp_password" desc:"[required] FTP password"`
	FtpHost     string `json:"ftp_host" desc:"[required] FTP hostname"`
	FtpPort     string `json:"ftp_port" desc:"[optional] FTP port number (default=21)"`
	FtpDir      string `json:"ftp_dir" desc:"[optional] FTP directory to be installed"`

	DBType     string `json:"db_type" desc:"[required] DB type"`
	DBHost     string `json:"db_host" desc:"[optional] DB hostname"`
	DBName     string `json:"db_name" desc:"[optional] DB name"`
	DBUser     string `json:"db_user" desc:"[optional] DB username"`
	DBPassword string `json:"db_password" desc:"[optional] DB password"`
	DBPrefix   string `json:"db_prefix" desc:"[optional] Prefix of table name on WordPress (default=mysite_)"`
	DBPort     string `json:"db_port" desc:"[optional] DB port number (default = 3306)"`

	SiteURL      string `json:"site_url" desc:"[required] Site URL"`
	SiteUser     string `json:"site_user" desc:"[required] LoginID of BaserCMS"`
	SitePassword string `json:"site_password" desc:"[required] LoginPassword of BaserCMS"`
	SiteEmail    string `json:"site_email" desc:"[required] Email address to register with BaserCMS"`

	InitArchiveURL string `json:"init_archive_url"`
	InitArchiveDir string `json:"init_archive_dir"`
	InitToken      string `json:"init_token"`
}

func (cnf *Config) init() error {
	if cnf.FtpPort == "" {
		cnf.FtpPort = "21"
	}

	if cnf.DBPrefix == "" {
		cnf.DBPrefix = "mysite_"
	}

	if cnf.InitArchiveURL == "" {
		cnf.InitArchiveURL = initArchiveURL
	}

	if cnf.InitArchiveDir == "" {
		cnf.InitArchiveDir = initArchiveDir
	}

	if cnf.InitToken == "" {
		token, err := randstr.Generate(64)
		if err != nil {
			return err
		}
		cnf.InitToken = token
	}

	return nil
}

func (cnf *Config) validate() error {
	msgs := []string{}

	kvs := map[string]string{
		"ftp_login_id":  cnf.FtpLoginID,
		"ftp_password":  cnf.FtpPassword,
		"ftp_host":      cnf.FtpHost,
		"db_type":       cnf.DBType,
		"site_url":      cnf.SiteURL,
		"site_user":     cnf.SiteUser,
		"site_password": cnf.SitePassword,
		"site_email":    cnf.SiteEmail,
	}
	for k, v := range kvs {
		if v == "" {
			msgs = append(msgs, fmt.Sprintf("%+v is required.", k))
		}
	}

	k := "db_type"
	v := cnf.DBType
	if !(v == "mysql" || v == "postgres" || v == "sqlite" || v == "csv") {
		msgs = append(msgs, fmt.Sprintf("%+v is invalid.", k))
	}

	if cnf.DBType == "mysql" || cnf.DBType == "postgres" {
		kvs = map[string]string{
			"db_name":     cnf.DBName,
			"db_user":     cnf.DBUser,
			"db_password": cnf.DBPassword,
			"db_host":     cnf.DBHost,
			"db_port":     cnf.DBPort,
		}
		for k, v := range kvs {
			if v == "" {
				msgs = append(msgs, fmt.Sprintf("%+v is required.", k))
			}
		}
	}

	k = "site_email"
	v = cnf.SiteEmail
	_, err := mail.ParseAddress(v)
	if err != nil {
		msgs = append(msgs, fmt.Sprintf("%+v is invalid format.", k))
	}

	kvs = map[string]string{
		"site_url":    cnf.SiteURL,
		"archive_url": cnf.InitArchiveURL,
	}
	for k, v := range kvs {
		if v == "" {
			break
		}
		parsedURL, err := url.Parse(v)
		if err != nil {
			msgs = append(msgs, fmt.Sprintf("%+v is invalid format.", k))
		} else if parsedURL.Scheme == "" || parsedURL.Host == "" {
			msgs = append(msgs, fmt.Sprintf("%+v is invalid format.", k))
		}
	}

	if len(msgs) > 0 {
		msg := strings.Join(msgs, " / ")
		return errors.New(msg)
	}

	return nil
}

func (cnf *Config) DescList() [][2]string {
	list := [][2]string{}
	rt := reflect.TypeOf(*cnf)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		k := f.Tag.Get("json")
		v := f.Tag.Get("desc")
		if k != "" && v != "" {
			list = append(list, [2]string{k, v})
		}
	}
	return list
}
