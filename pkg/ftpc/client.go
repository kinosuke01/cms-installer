package ftpc

import (
	"bytes"
	"crypto/tls"
	"path"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type Config struct {
	LoginID  string
	Password string
	Host     string
	Port     string
}

type Client struct {
	Config
	conn *ftp.ServerConn
}

func New(cnf *Config) *Client {
	return &Client{
		Config: *cnf,
	}
}

func (c *Client) login() error {
	id := c.LoginID
	password := c.Password
	host := c.Host
	port := c.Port

	conn, err := ftp.Dial(
		(host + ":" + port),
		ftp.DialWithTimeout(5*time.Second),
		ftp.DialWithExplicitTLS(
			&tls.Config{
				InsecureSkipVerify: true,
				ServerName:         host,
			},
		),
	)
	if err != nil {
		return err
	}
	c.conn = conn

	err = conn.Login(id, password)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) quit() error {
	if err := c.conn.Quit(); err != nil {
		return err
	}
	return nil
}

func (c *Client) exists(entryPath string) (bool, error) {
	parentPath := path.Join(entryPath, "../")
	entryName := path.Base(entryPath)

	entries, err := c.conn.List(parentPath)
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entryName == entry.Name {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) mkdirP(dirPath string) error {
	if dirPath == "." {
		return nil
	}

	dirNames := strings.Split(dirPath, "/")

	d := "./"
	for _, dirName := range dirNames {
		d = path.Join(d, dirName)
		exists, err := c.exists(d)
		if err != nil {
			return err
		}
		if !exists {
			err := c.conn.MakeDir(d)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) Upload(filePath string, pContent *string) error {
	if err := c.login(); err != nil {
		return err
	}
	defer c.quit()

	if err := c.mkdirP(path.Dir(filePath)); err != nil {
		return err
	}

	data := bytes.NewBufferString(*pContent)
	err := c.conn.Stor(filePath, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(filePath string, onlyExists bool) error {
	if err := c.login(); err != nil {
		return err
	}
	defer c.quit()

	if onlyExists {
		exists, err := c.exists(filePath)
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}
	}

	// No error if the file doesn't exist.
	// `onlyExists` may be unnecessary.
	err := c.conn.Delete(filePath)
	if err != nil {
		return err
	}

	return nil
}
