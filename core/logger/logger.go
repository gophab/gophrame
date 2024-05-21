package logger

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/wjshen/gophrame/core/global"

	"github.com/astaxie/beego/validation"
)

var Logger *log.Logger
var once sync.Once

func init() {
	once.Do(func() {
		if global.Debug {
			Logger = log.New(os.Stdout, "[DEBUG]\t", log.LstdFlags|log.Llongfile)
		} else {
			Logger = log.New(os.Stdout, "[INFO]\t", log.LstdFlags|log.Lshortfile)
		}
	})
}

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		Info(err.Key, err.Message)
	}
}

// Info 详情
func Info(args ...interface{}) {
	Logger.SetPrefix("[INFO]\t")
	Logger.Output(2, fmt.Sprintln(args...))
}

// Danger 错误 为什么不命名为 error？避免和 error 类型重名
func Danger(args ...interface{}) {
	Logger.SetPrefix("[ERROR]\t")
	Logger.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

// Warn 警告
func Warn(args ...interface{}) {
	Logger.SetPrefix("[WARNING]\t")
	Logger.Output(2, fmt.Sprintln(args...))
}

// Debug debug
func Debug(args ...interface{}) {
	Logger.SetPrefix("[DEBUG]\t")
	Logger.Output(2, fmt.Sprintln(args...))
}

func Error(args ...interface{}) {
	Logger.SetPrefix("[ERROR]\t")
	Logger.Output(2, fmt.Sprintln(args...))
}

func Fatal(args ...interface{}) {
	Logger.SetPrefix("[FATAL]\t")
	Logger.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}
