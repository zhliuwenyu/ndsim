package ndsim

import (
	"errors"
	"fmt"
	"log"
	"os"
)

//Log 对log包进行封装，进行初始化和wf日志的拆分
type Log struct {
	Path     string
	FileName string
	Level    int
	LogDN    *log.Logger
	LogWF    *log.Logger
}

const (
	logLevelFatal = 1 << iota
	logLevelWarning
	logLevelNotice
	logLevelDebug
)

//Debug for log Debug level
func (l *Log) Debug(v ...interface{}) {
	if l.Level > logLevelDebug {
		l.LogDN.Output(2, fmt.Sprintln(" Debug ", v))
	}
}

//Notice for log Notice level
func (l *Log) Notice(v ...interface{}) {
	if l.Level > logLevelNotice {
		l.LogDN.Output(2, fmt.Sprintln(" Notice ", v))
	}
}

//Warning for log Warning level
func (l *Log) Warning(v ...interface{}) {
	if l.Level > logLevelWarning {
		l.LogWF.Output(2, fmt.Sprintln(" Warning ", v))
	}
}

//Fatal for log Fatal level
func (l *Log) Fatal(v ...interface{}) {
	if l.Level > logLevelFatal {
		l.LogWF.Output(2, fmt.Sprintln(" Fatal ", v))
	}
	os.Exit(255)
}

//Info for notice log has changed and no log level to set
func (l *Log) Info(v ...interface{}) {
	l.LogDN.Println(v)
	l.LogWF.Println(v)
}

//gLog is global Log object
var gLog Log

const logPREFIX = "[ndsim]"

func initLog() error {
	if gLog.LogDN != nil {
		return nil
	}
	logPath, fileName, level := gConfig.LogPath, gConfig.LogFileName, gConfig.LogLevel
	dir, err := os.Lstat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(logPath, 0664)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			fmt.Println(err)
			return err
		}
	} else if !dir.IsDir() {
		fmt.Println(logPath, "is already exist and not path")
		return errors.New("log path is not dir")
	}
	logIO1, err := os.OpenFile(logPath+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	logIO2, err := os.OpenFile(logPath+"/"+fileName+".wf", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log1 := log.New(logIO1, logPREFIX, log.Llongfile|log.Ldate|log.Ldate|log.Ltime)
	log2 := log.New(logIO2, logPREFIX, log.Llongfile|log.Ldate|log.Ldate|log.Ltime)
	gLog = Log{Path: logPath, FileName: fileName, Level: level, LogDN: log1, LogWF: log2}

	return nil
}
