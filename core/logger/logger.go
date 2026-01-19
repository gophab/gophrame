package logger

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gophab/gophrame/core/global"

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
func Info(args ...any) {
	Logger.SetPrefix("[INFO]\t")
	Logger.Output(2, fmt.Sprintln(args...))
}

// Danger 错误 为什么不命名为 error？避免和 error 类型重名
func Danger(args ...any) {
	Logger.SetPrefix("[ERROR]\t")
	Logger.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}

// Warn 警告
func Warn(args ...any) {
	Logger.SetPrefix("[WARNING]\t")
	Logger.Output(2, fmt.Sprintln(args...))
}

// Debug debug
func Debug(args ...any) {
	Logger.SetPrefix("[DEBUG]\t")
	Logger.Output(2, fmt.Sprintln(args...))
}

func Error(args ...any) {
	Logger.SetPrefix("[ERROR]\t")
	Logger.Output(2, fmt.Sprintln(args...))
}

func Fatal(args ...any) {
	Logger.SetPrefix("[FATAL]\t")
	Logger.Output(2, fmt.Sprintln(args...))
	os.Exit(1)
}
