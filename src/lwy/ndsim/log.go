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
		l.LogDN.Println(" DEBUG ", v)
	}
}

//Notice for log Notice level
func (l *Log) Notice(v ...interface{}) {
	if l.Level > logLevelNotice {
		l.LogDN.Println(" Notice ", v)
	}
}

//Warning for log Warning level
func (l *Log) Warning(v ...interface{}) {
	if l.Level > logLevelWarning {
		l.LogWF.Println(" Warning ", v)
	}
}

//Fatal for log Fatal level
func (l *Log) Fatal(v ...interface{}) {
	if l.Level > logLevelFatal {
		l.LogWF.Println(" Fatal ", v)
	}
}

//Info for notice log has changed and no log level to set
func (l *Log) Info(v ...interface{}) {
	l.LogDN.Println(v)
	l.LogWF.Println(v)
}

//GLog is global Log object
var GLog Log

const logPREFIX = "[ndsim]"

func initLog() error {
	logPath, fileName, level := GConfig.LogPath, GConfig.LogFileName, GConfig.LogLevel
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
	GLog = Log{Path: logPath, FileName: fileName, Level: level, LogDN: log1, LogWF: log2}

	return nil
}
