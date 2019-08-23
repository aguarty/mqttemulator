package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	//"log/syslog"
)

type severity int

// Severity levels.
const (
	sDebug severity = iota
	sInfo
	sWarning
	sError
	sFatal
)

//Logger logger
type Logger struct {
	logDebug *log.Logger
	logInfo  *log.Logger
	logWarn  *log.Logger
	logError *log.Logger
	logFatal *log.Logger
	loglevel severity
}

//Info print with info level
func (l *Logger) Info(v ...interface{}) {
	if l.loglevel <= sInfo {
		l.output(sInfo, 0, fmt.Sprint(v...))
	} else {
		return
	}
}

//Infof printf with info level
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.loglevel <= sInfo {
		l.output(sInfo, 0, fmt.Sprintf(format, v...))
	} else {
		return
	}
}

//Debug print with debug level
func (l *Logger) Debug(v ...interface{}) {
	if l.loglevel <= sDebug {
		l.output(sDebug, 0, fmt.Sprint(v...))
	} else {
		return
	}
}

//Debugf printf with debug level
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.loglevel <= sDebug {
		l.output(sDebug, 0, fmt.Sprintf(format, v...))
	} else {
		return
	}
}

//Error print with error level
func (l *Logger) Error(v ...interface{}) {
	if l.loglevel <= sError {
		l.output(sError, 0, fmt.Sprint(v...))
	} else {
		return
	}
}

//Errorf printf with error level
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.loglevel <= sError {
		l.output(sError, 0, fmt.Sprintf(format, v...))
	} else {
		return
	}
}

//Warn print with warn level
func (l *Logger) Warn(v ...interface{}) {
	if l.loglevel <= sWarning {
		l.output(sWarning, 0, fmt.Sprint(v...))
	} else {
		return
	}
}

//Warnf print with warnf level
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.loglevel <= sWarning {
		l.output(sWarning, 0, fmt.Sprintf(format, v...))
	} else {
		return
	}
}

//Fatal print with fatal level
func (l *Logger) Fatal(v ...interface{}) {
	l.output(sFatal, 0, fmt.Sprint(v...))
	os.Exit(1)
}

//Fatalf printf with fatal level
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.output(sFatal, 0, fmt.Sprintf(format, v...))
	os.Exit(1)
}

//output print result
func (l *Logger) output(s severity, depth int, txt string) {
	switch s {
	case sDebug:
		l.logDebug.Output(3, txt)
	case sInfo:
		l.logInfo.Output(3, txt)
	case sWarning:
		l.logWarn.Output(3, txt)
	case sError:
		l.logError.Output(3, txt)
	case sFatal:
		l.logFatal.Output(3, txt)
	default:
		panic(fmt.Sprintln("unrecognized severity:", s))
	}
}

//initLogger initialize logger
func Init(logLevel string, logFilePath string) *Logger {

	lg := &Logger{}

	switch logLevel {
	case "debug":
		lg.loglevel = sDebug
	case "info":
		lg.loglevel = sInfo
	case "warning":
		lg.loglevel = sWarning
	case "error":
		lg.loglevel = sError
	case "fatal":
		lg.loglevel = sFatal

	}

	var logFile io.Writer

	if logFilePath == "" {
		logFile = ioutil.Discard
	} else {
		var err error
		dir, file := filepath.Split(logFilePath)
		if _, err = os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		}
		logFile, err = os.OpenFile(dir+file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
	}

	var out io.Writer
	out = os.Stderr

	iLogs := []io.Writer{out, logFile}
	wLogs := []io.Writer{out, logFile}
	eLogs := []io.Writer{out, logFile}
	dLogs := []io.Writer{out, logFile}
	fLogs := []io.Writer{out, logFile}

	lg.logDebug = log.New(io.MultiWriter(dLogs...), "DEBUG: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	lg.logInfo = log.New(io.MultiWriter(iLogs...), "INFO: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	lg.logWarn = log.New(io.MultiWriter(wLogs...), "WARN: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	lg.logError = log.New(io.MultiWriter(eLogs...), "ERROR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	lg.logFatal = log.New(io.MultiWriter(fLogs...), "FATAL: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	return lg

}

// For output to syslog
// func setup(src string) (*syslog.Writer, *syslog.Writer, *syslog.Writer, error) {
// 	const facility = syslog.LOG_USER
// 	il, err := syslog.New(facility|syslog.LOG_NOTICE, src)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	wl, err := syslog.New(facility|syslog.LOG_WARNING, src)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	el, err := syslog.New(facility|syslog.LOG_ERR, src)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}
// 	return il, wl, el, nil
// }
