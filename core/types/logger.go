package types

import (
	"fmt"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"runtime"
)

type CustomLogger struct {
	tmLog.Logger
}

func (l *CustomLogger) Info(msg string, keyvals ...interface{}) {
	file, line := getCallerInfo()
	keyvals = append(keyvals, "file", file, "line", line)
	l.Logger.Info(msg, keyvals...)
}

func (l *CustomLogger) Error(msg string, keyvals ...interface{}) {
	file, line := getCallerInfo()
	keyvals = append(keyvals, "file", file, "line", line)
	l.Logger.Error(msg, keyvals...)
}

func (l *CustomLogger) Debug(msg string, keyvals ...interface{}) {
	file, line := getCallerInfo()
	keyvals = append(keyvals, "file", file, "line", line)
	l.Logger.Debug(msg, keyvals...)
}

///////////////////////////////////////////////////////////////////////

func (l *CustomLogger) Errorf(s string, i ...interface{}) {
	errMsg := fmt.Sprintf(s, i...)
	keyvals := []interface{}{"error", errMsg}
	l.Logger.Error("Badger", keyvals...)
}

func (l *CustomLogger) Warningf(s string, i ...interface{}) {
	errMsg := fmt.Sprintf(s, i...)
	keyvals := []interface{}{"error", errMsg}
	l.Logger.Error("Badger", keyvals...)
}

func (l *CustomLogger) Infof(s string, i ...interface{}) {
	infoMsg := fmt.Sprintf(s, i...)
	keyvals := []interface{}{"info", infoMsg}
	l.Logger.Info("Badger", keyvals...)
}

func (l *CustomLogger) Debugf(s string, i ...interface{}) {
	debugMsg := fmt.Sprintf(s, i...)
	keyvals := []interface{}{"debug", debugMsg}
	l.Logger.Debug("Badger", keyvals...)
}

// getCallerInfo returns the file name and line number of the caller's parent.
// It uses runtime.Caller to retrieve this information.
// If the caller information cannot be retrieved, it returns "unknown" as the file name and 0 as the line number.
func getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown", 0
	}
	return file, line
}
