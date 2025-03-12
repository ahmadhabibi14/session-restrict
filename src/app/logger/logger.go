package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var BaseDir string = `web-api-docs`

func init() {
	fullPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	BaseDir = filepath.Base(fullPath)
}

const (
	logFieldMsg  = `message`
	logFieldStck = `stack`
)

type Logger struct {
	logger *logrus.Logger
}

func (l *Logger) GetLogger() *logrus.Logger {
	return l.logger
}

func (l *Logger) InfoWithFields(info any, fields logrus.Fields) {
	l.logger.WithFields(fields).Info(info)
}

func (l *Logger) Trace(trace any, msg ...string) {
	_, file, line, _ := runtime.Caller(1)

	idxBase := strings.Index(file, BaseDir)
	str := color.CyanString(`▣ ` + file[idxBase:] + `:` + strconv.Itoa(line) + ` ▶`)
	now := time.Now()
	formattedTime := now.Format("2006/01/02 03:04 PM")
	stack := color.YellowString("STAC")
	fmt.Println(stack + "[" + formattedTime + "] " + strings.Replace(str, `%`, `%%`, -1))

	if len(msg) > 0 {
		l.logger.WithField(logFieldMsg, msg).Trace(trace)
	} else {
		l.logger.Trace(trace)
	}

	stack = file + `:` + strconv.Itoa(line)

	if len(msg) > 0 {
		l.logger.
			WithField(logFieldMsg, msg).
			WithField(logFieldStck, stack).
			Trace(trace)
	} else {
		l.logger.Trace(trace)
	}
}

func (l *Logger) Debug(debug any, msg ...string) {
	_, file, line, _ := runtime.Caller(1)

	idxBase := strings.Index(file, BaseDir)
	str := color.CyanString(`▣ ` + file[idxBase:] + `:` + strconv.Itoa(line) + ` ▶`)
	now := time.Now()
	formattedTime := now.Format("2006/01/02 03:04 PM")
	stack := color.YellowString("STAC")
	fmt.Println(stack + "[" + formattedTime + "] " + strings.Replace(str, `%`, `%%`, -1))

	stack = file + `:` + strconv.Itoa(line)

	if len(msg) > 0 {
		l.logger.
			WithField(logFieldMsg, msg).
			WithField(logFieldStck, stack).
			Debug(debug)
	} else {
		l.logger.Debug(debug)
	}
}

func (l *Logger) Info(info any, msg ...string) {
	_, file, line, _ := runtime.Caller(1)

	idxBase := strings.Index(file, BaseDir)
	str := color.CyanString(`▣ ` + file[idxBase:] + `:` + strconv.Itoa(line) + ` ▶`)
	now := time.Now()
	formattedTime := now.Format("2006/01/02 03:04 PM")
	stack := color.YellowString("STAC")
	fmt.Println(stack + "[" + formattedTime + "] " + strings.Replace(str, `%`, `%%`, -1))

	stack = file + `:` + strconv.Itoa(line)

	if len(msg) > 0 {
		l.logger.
			WithField(logFieldMsg, msg).
			WithField(logFieldStck, stack).
			Info(info)
	} else {
		l.logger.Info(info)
	}
}

func (l *Logger) Warn(warn any, msg ...string) {
	_, file, line, _ := runtime.Caller(1)

	idxBase := strings.Index(file, BaseDir)
	str := color.CyanString(`▣ ` + file[idxBase:] + `:` + strconv.Itoa(line) + ` ▶`)
	now := time.Now()
	formattedTime := now.Format("2006/01/02 03:04 PM")
	stack := color.YellowString("STAC")
	fmt.Println(stack + "[" + formattedTime + "] " + strings.Replace(str, `%`, `%%`, -1))

	stack = file + `:` + strconv.Itoa(line)

	if len(msg) > 0 {
		l.logger.
			WithField(logFieldMsg, msg).
			WithField(logFieldStck, stack).
			Warn(warn)
	} else {
		l.logger.Warn(warn)
	}
}

func (l *Logger) Error(err any, msg ...string) {
	_, file, line, _ := runtime.Caller(1)
	idxBase := strings.Index(file, BaseDir)
	str := color.CyanString(`▣ ` + file[idxBase:] + `:` + strconv.Itoa(line) + ` ▶`)
	now := time.Now()
	formattedTime := now.Format("2006/01/02 03:04 PM")
	stack := color.YellowString("STAC")
	fmt.Println(stack + "[" + formattedTime + "] " + strings.Replace(str, `%`, `%%`, -1))

	stack = file + `:` + strconv.Itoa(line)

	if len(msg) > 0 {
		l.logger.
			WithField(logFieldMsg, msg).
			WithField(logFieldStck, stack).
			Error(err)
	} else {
		l.logger.Error(err)
	}
}

func (l *Logger) Fatal(fatal any, msg ...string) {
	_, file, line, _ := runtime.Caller(1)

	idxBase := strings.Index(file, BaseDir)
	str := color.CyanString(`▣ ` + file[idxBase:] + `:` + strconv.Itoa(line) + ` ▶`)
	now := time.Now()
	formattedTime := now.Format("2006/01/02 03:04 PM")
	stack := color.YellowString("STAC")
	fmt.Println(stack + "[" + formattedTime + "] " + strings.Replace(str, `%`, `%%`, -1))

	stack = file + `:` + strconv.Itoa(line)

	if len(msg) > 0 {
		l.logger.WithField(logFieldMsg, msg).
			WithField(logFieldStck, stack).
			Fatal(fatal)
	} else {
		l.logger.WithField(logFieldStck, stack).Fatal(fatal)
	}
}

func (l *Logger) Panic(panic any, msg string) {
	_, file, line, _ := runtime.Caller(1)

	idxBase := strings.Index(file, BaseDir)
	str := color.CyanString(`▣ ` + file[idxBase:] + `:` + strconv.Itoa(line) + ` ▶`)
	now := time.Now()
	formattedTime := now.Format("2006/01/02 03:04 PM")
	stack := color.YellowString("STAC")
	fmt.Println(stack + "[" + formattedTime + "] " + strings.Replace(str, `%`, `%%`, -1))

	stack = file + `:` + strconv.Itoa(line)

	if len(msg) > 0 {
		l.logger.WithField(logFieldMsg, msg).
			WithField(logFieldStck, stack).
			Panic(panic)
	} else {
		l.logger.WithField(logFieldStck, stack).Panic(panic)
	}
}

var Log *Logger

func InitLogger() {
	logger := logrus.New()

	defer func() {
		Log = &Logger{
			logger: logger,
		}
	}()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: `2006/01/02 03:04 PM`,
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
}
