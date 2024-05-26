package starter

import (
	"sync"

	_ "github.com/gophab/gophrame/core/social"
	_ "github.com/gophab/gophrame/core/social/dingtalk"
	_ "github.com/gophab/gophrame/core/social/feishu"
	_ "github.com/gophab/gophrame/core/social/wxcp"
	_ "github.com/gophab/gophrame/core/social/wxma"
	_ "github.com/gophab/gophrame/core/social/wxmp"

	"github.com/gophab/gophrame/core/social/config"

	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/starter"
)

var (
	once sync.Once
)

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	logger.Debug("Starting Social: ...", config.Setting.Enabled)
	if config.Setting.Enabled {
		once.Do(func() {
			// ...
		})

	}
}
