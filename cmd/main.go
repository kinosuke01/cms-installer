package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	wpi "github.com/kinosuke01/cms-installer/wordpress"
)

// Receive the password value from the environment variable.
type PwParams struct {
	ftp  string
	db   string
	site string
}

func (pp *PwParams) init() {
	pp.ftp = os.Getenv("CMSI_FTP_PASSOWRD")
	pp.db = os.Getenv("CMSI_DB_PASSOWRD")
	pp.site = os.Getenv("CMSI_SITE_PASSWORD")
}

// A simple logger with standard output.
type Logger struct {
	withTS bool
}

func (l *Logger) write(msg string, level string) {
	if l.withTS {
		ts := time.Now().Local().Format("2006-01-02T15:04:05-07:00")
		fmt.Printf("%+v %+v %+v\n", ts, level, msg)
	} else {
		fmt.Printf("%+v %+v\n", level, msg)
	}
}

func (l *Logger) Info(msg string) {
	l.write(msg, "INFO")
}

func (l *Logger) Error(msg string) {
	l.write(msg, "ERROR")
}

func main() {
	optHelp := flag.Bool("h", false, "Show usage")
	flag.Parse()

	if *optHelp {
		showUsage()
		return
	}

	args := flag.Args()

	if len(args) != 2 {
		showUsage()
		return
	}

	cms := args[0]
	paramStr := args[1]
	pp := &PwParams{}
	pp.init()
	var err error

	switch cms {
	case "wp", "wordpress":
		err = wpInstall(paramStr, pp)
	// TODO other cms installation
	default:
		err = errors.New("cms name is invalid")
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Installation completed")
	}
}

func wpInstall(paramStr string, pp *PwParams) error {
	var config wpi.Config
	err := json.Unmarshal([]byte(paramStr), &config)
	if err != nil {
		return fmt.Errorf("not json ( %+v )", err.Error())
	}

	if config.FtpPassword == "" {
		config.FtpPassword = pp.ftp
	}
	if config.DBPassword == "" {
		config.DBPassword = pp.db
	}
	if config.SitePassword == "" {
		config.SitePassword = pp.site
	}

	wp, err := wpi.New(&config)
	if err != nil {
		return err
	}

	wp.Logger = &Logger{
		withTS: true,
	}

	err = wp.Install()
	if err != nil {
		return err
	}

	return nil
}

func showUsage() {
	toText := func(list [][2]string) string {
		var text string
		for _, kv := range list {
			k := kv[0]
			v := kv[1]
			text = text + fmt.Sprintf("%20s", k) + "\t" + v + "\n"
		}
		return text
	}

	text := usageTemplate

	cmstext := toText([][2]string{
		{"wp, wordpress", "WordPress installation"},
	})
	text = strings.Replace(text, "CMS_NAME_PLACEHOLDER", cmstext, 1)

	wpcnf := wpi.Config{}
	wptext := toText(wpcnf.DescList())
	text = strings.Replace(text, "WORDPRESS_PLACEHOLDER", wptext, 1)

	fmt.Println(text)
}

const usageTemplate = `Usage: cmsi <cms_name> <json_params>

<cms_name>:
CMS_NAME_PLACEHOLDER

<json_params> (for wordpress):
WORDPRESS_PLACEHOLDER
`
