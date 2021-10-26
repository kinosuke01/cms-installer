package cmsinit

import (
	"strconv"
	"strings"
)

func Script(cnf *Config) (string, error) {
	cnf.init()
	err := cnf.validate()

	if err != nil {
		return "", err
	}

	str := php
	str = strings.Replace(str, "ARCHIVE_URL_PLACEHOLDER", cnf.URL, 1)
	str = strings.Replace(str, "EXTRACTED_DIR_PLACEHOLDER", cnf.Dir, 1)
	str = strings.Replace(str, "TOKEN_PLACEHOLDER", cnf.Token, 1)

	t := *cnf.expiredAt()
	timestr := strconv.FormatInt(t.Unix(), 10)
	str = strings.Replace(str, "EXPIRED_AT_PLACEHOLDER", timestr, 1)

	return str, nil
}
