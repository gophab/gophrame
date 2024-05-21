package websocket

const (
	ErrorsWebsocketOnOpenFail                 string = "websocket onopen 发生阶段错误"
	ErrorsWebsocketUpgradeFail                string = "websocket Upgrade 协议升级, 发生错误"
	ErrorsWebsocketReadMessageFail            string = "websocket ReadPump(实时读取消息)协程出错"
	ErrorsWebsocketBeatHeartFail              string = "websocket BeatHeart心跳协程出错"
	ErrorsWebsocketBeatHeartsMoreThanMaxTimes string = "websocket BeatHeart 失败次数超过最大值"
	ErrorsWebsocketSetWriteDeadlineFail       string = "websocket  设置消息写入截止时间出错"
	ErrorsWebsocketWriteMgsFail               string = "websocket  Write Msg(send msg) 失败"
	ErrorsWebsocketStateInvalid               string = "websocket  state 状态已经不可用(掉线、卡死等愿意，造成双方无法进行数据交互)"
)
