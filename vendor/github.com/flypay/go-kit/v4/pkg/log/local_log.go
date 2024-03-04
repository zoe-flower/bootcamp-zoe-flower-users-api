package log

import (
	"fmt"
	"strings"

	"github.com/aws/aws-xray-sdk-go/xraylog"

	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

type localFormatter struct {
	logrus.TextFormatter
}

const timeFormat = "2006-01-02 15:04:05.000"

type XrayLocalLogger struct {
	Logger
}

func (l *XrayLocalLogger) Log(level xraylog.LogLevel, msg fmt.Stringer) {
	// debugs are pretty noisy, so let's remove them
	if level < xraylog.LogLevelInfo {
		return
	}

	switch level {
	case xraylog.LogLevelWarn:
		l.Warn(msg)
	case xraylog.LogLevelError:
		l.Error(msg)
	default:
		l.Info(msg)
	}
}

func (f *localFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	lvl := strings.ToUpper(entry.Level.String())
	if len(lvl) > 5 {
		lvl = lvl[:5]
	}

	var lvlOut aurora.Value
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		lvlOut = aurora.Blue(lvl)
	case logrus.WarnLevel:
		lvlOut = aurora.Yellow(lvl)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		lvlOut = aurora.Red(lvl)
	default:
		lvlOut = aurora.Gray(15, lvl)
	}

	time := entry.Time.Format(timeFormat)
	msg := entry.Message

	o := fmt.Sprintf("%s\t[%s]\t%s\n", time, lvlOut, msg)

	return []byte(o), nil
}
