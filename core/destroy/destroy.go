package destroy

import (
	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/logger"

	"os"
	"os/signal"
	"syscall"
)

const (
	// 进程被结束
	ProcessKilled string = "收到信号，进程被结束"
)

func init() {
	//  用于系统信号的监听
	go func() {
		c := make(chan os.Signal, 1000)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM) // 监听可能的退出信号
		received := <-c                                                                           //接收信号管道中的值
		logger.Warn(ProcessKilled, "信号值", received.String())
		eventbus.FuzzyPublishEvent(global.EventDestroyPrefix)
		close(c)
		os.Exit(1)
	}()

}
