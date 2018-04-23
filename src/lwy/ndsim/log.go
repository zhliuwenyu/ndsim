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
	logLevelDebug = 1 << iota
	logLevelNotice
	logLevelWarning
	logLevelFatal
)

//Debug for log Debug level
func (l *Log) Debug(v ...interface{}) {
	if l.Level&logLevelDebug != 0 {
		l.LogDN.Println(" DEBUG ", v)
	}
}

//Notice for log Notice level
func (l *Log) Notice(v ...interface{}) {
	if l.Level&logLevelDebug != 0 {
		l.LogDN.Println(" Notice ", v)
	}
}

//Warning for log Warning level
func (l *Log) Warning(v ...interface{}) {
	if l.Level&logLevelDebug != 0 {
		l.LogWF.Println(" Warning ", v)
	}
}

//Fatal for log Fatal level
func (l *Log) Fatal(v ...interface{}) {
	if l.Level&logLevelDebug != 0 {
		l.LogWF.Println(" Fatal ", v)
	}
}

//GLog is global Log object
var GLog Log

const logPREFIX = "[ndsim]"

func initLog(path string, fileName string, level int) error {
	dir, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0664)
			if err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			fmt.Println(err)
			return err
		}
	}
	if !dir.IsDir() {
		fmt.Println(path, "is already exist and not path")
		return errors.New("log path is not dir")
	}
	logIO1, err := os.OpenFile(path+"/"+fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	logIO2, err := os.OpenFile(path+"/"+fileName+".wf", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log1 := log.New(logIO1, logPREFIX, log.Llongfile|log.Ldate|log.Ldate|log.Ltime)
	log2 := log.New(logIO2, logPREFIX, log.Llongfile|log.Ldate|log.Ldate|log.Ltime)
	GLog = Log{Path: path, FileName: fileName, Level: level, LogDN: log1, LogWF: log2}

	return nil
}
