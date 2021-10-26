package withlog

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Logger interface {
	Info(string)
	Error(string)
}

// Log output before and after function execution
func Exec(f func() error, fn string, logger Logger) error {
	if (logger == nil) || reflect.ValueOf(logger).IsNil() {
		return f()
	}

	hash := genHash(fn)

	logger.Info(fmt.Sprintf("[%+v] Running %+v", hash, fn))
	startAt := time.Now()

	err := f()

	endAt := time.Now()
	sec := period(&startAt, &endAt)
	if err == nil {
		logger.Info(fmt.Sprintf("[%+v] Finished in %+v seconds (successful).", hash, sec))
	} else {
		logger.Info(fmt.Sprintf("[%+v] Finished in %+v seconds (failed).", hash, sec))
		logger.Error(fmt.Sprintf("[%+v] %+v %+v", hash, fn, err.Error()))
	}

	return err
}

func genHash(fn string) string {
	pstr := strconv.FormatInt(time.Now().UnixNano(), 10) + "__" + fn

	hash := md5.New()
	defer hash.Reset()
	hash.Write([]byte(pstr))
	hstr := hex.EncodeToString(hash.Sum(nil))

	return string([]rune(hstr)[:8])
}

func period(startAt, endAt *time.Time) float32 {
	sec := (float32)(endAt.UnixMilli()-startAt.UnixMilli()) * 0.001
	return sec
}
