package log

import (
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func GetLog(ignorePath ...bool) *log.Entry {
	needPath := true
	if len(ignorePath) > 0 {
		needPath = !ignorePath[0]
	}
	var fileName, funcName string
	pc, file, _, ok := runtime.Caller(1)

	if ok {
		path := strings.Split(file, "/")
		fileName = path[len(path)-1]
	}
	fn := runtime.FuncForPC(pc).Name()
	path := strings.Split(fn, "/")
	funcName = path[len(path)-1]

	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)

	if needPath {
		return logger.WithFields(log.Fields{
			"file": fileName,
			"func": funcName,
		})
	}
	return logger.WithFields(log.Fields{})

}
