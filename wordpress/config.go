package wordpress

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"strings"
)

type Config struct {
	FtpLoginID  string `json:"ftp_login_id" desc:"[required] FTP loginID"`
	FtpPassword string `json:"ftp_password" desc:"[required] FTP password"`
	FtpHost     string `json:"ftp_host" desc:"[required] FTP hostname"`
	FtpPort     string `json:"ftp_port" desc:"[optional] FTP port number (default=21)"`
	FtpDir      string `json:"ftp_dir" desc:"[optional] FTP directory to be installed"`

	DBName     string `json:"db_name" desc:"[required] DB name"`
	DBUser     string `json:"db_user" desc:"[required] DB username"`
	DBPassword string `json:"db_password" desc:"[required] DB password"`
	DBHost     string `json:"db_host" desc:"[required] DB hostname"`
	DBPrefix   string `json:"db_prefix" desc:"[optional] Prefix of table name on WordPress (default=wp_)"`

	SiteURL      string `json:"site_url" desc:"[required] Site URL"`
	SiteTitle    string `json:"site_title" desc:"[required] Site name"`
	SiteUser     string `json:"site_user" desc:"[required] LoginID of WordPress"`
	SitePassword string `json:"site_password" desc:"[required] LoginPassword of WordPress"`
	SiteEmail    string `json:"site_email" desc:"[required] Email address to register with WordPress"`
	SiteUnpublic bool   `json:"site_unpublic" desc:"[optional] Prevent search engines from indexing your site (default=false)"`

	InitArchiveURL string `json:"init_archive_url"`
	InitArchiveDir string `json:"init_archive_dir"`
}

func (cnf *Config) init() error {
	if cnf.FtpPort == "" {
		cnf.FtpPort = "21"
	}

	if cnf.DBPrefix == "" {
		cnf.DBPrefix = "wp_"
	}

	if cnf.InitArchiveURL == "" {
		cnf.InitArchiveURL = initArchiveURL
	}

	if cnf.InitArchiveDir == "" {
		cnf.InitArchiveDir = initArchiveDir
	}

	return nil
}

func (cnf *Config) validate() error {
	msgs := []string{}

	kvs := map[string]string{
		"ftp_login_id":  cnf.FtpLoginID,
		"ftp_password":  cnf.FtpPassword,
		"ftp_host":      cnf.FtpHost,
		"db_name":       cnf.DBName,
		"db_user":       cnf.DBUser,
		"db_password":   cnf.DBPassword,
		"db_host":       cnf.DBHost,
		"site_url":      cnf.SiteURL,
		"site_title":    cnf.SiteTitle,
		"site_user":     cnf.SiteUser,
		"site_password": cnf.SitePassword,
		"site_email":    cnf.SiteEmail,
	}
	for k, v := range kvs {
		if v == "" {
			msgs = append(msgs, fmt.Sprintf("%+v is required.", k))
		}
	}

	k := "site_email"
	v := cnf.SiteEmail
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
