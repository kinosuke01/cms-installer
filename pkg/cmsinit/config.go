package cmsinit

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const defaultURL = "https://ja.wordpress.org/latest-ja.zip"
const defaultDir = "wordpress"
const defaultPeriod = 60 // 1min

type Config struct {
	URL    string
	Dir    string
	Token  string
	Period int
	Now    *time.Time
}

func (cnf *Config) init() {
	if cnf.URL == "" {
		cnf.URL = defaultURL
	}
	if cnf.Dir == "" {
		cnf.Dir = defaultDir
	}
	if cnf.Period == 0 {
		cnf.Period = defaultPeriod
	}
	if cnf.Now == nil {
		t := time.Now().Local()
		cnf.Now = &t
	}
}

func (cnf *Config) validate() error {
	msgs := []string{}

	k := "token"
	v := cnf.Token
	if v == "" {
		msgs = append(msgs, fmt.Sprintf("%+v is required.", k))
	}

	if len(msgs) > 0 {
		msg := strings.Join(msgs, " / ")
		return errors.New(msg)
	}

	return nil
}

func (cnf *Config) expiredAt() *time.Time {
	t := cnf.Now.Add(time.Duration(cnf.Period) * time.Second).Local()
	return &t
}
