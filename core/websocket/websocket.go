package websocket

import (
	"github.com/wjshen/gophrame/core/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/**
websocket 想要了解更多具体细节请参见以下文档
文档地址：https://github.com/gorilla/websocket/tree/master/examples
*/

type Websocket struct {
	Client *Client
}

// onOpen 事件函数
func (w *Websocket) OnOpen(context *gin.Context) (*Websocket, bool) {
	if client, ok := (&Client{}).OnOpen(context); ok {

		token := context.GetString("token")
		logger.Info("获取到的客户端上线时携带的唯一标记值：", "token", token)

		// 成功上线以后，开发者可以基于客户端上线时携带的唯一参数(这里用token键表示)
		// 在数据库查询更多的其他字段信息，直接追加在 Client 结构体上，方便后续使用
		//client.ClientMoreParams.UserParams1 = "123"
		//client.ClientMoreParams.UserParams2 = "456"
		//fmt.Printf("最终每一个客户端(client) 已有的参数：%+v\n", client)

		w.Client = client

		go w.Client.Heartbeat() // 一旦握手+协议升级成功，就为每一个连接开启一个自动化的隐式心跳检测包

		return w, true
	}

	return nil, false
}

// OnMessage 处理业务消息
func (w *Websocket) OnMessage(context *gin.Context) {
	go w.Client.ReadPump(func(messageType int, receivedData []byte) {
		//参数说明
		//messageType 消息类型，1=文本
		//receivedData 服务器接收到客户端（例如js客户端）发来的的数据，[]byte 格式
		// 实际项目中，消息的传递请统一按照json格式传递
		tempMsg := "服务器已经收到了你的消息：" + string(receivedData)
		// 回复客户端已经收到消息;
		if err := w.Client.SendMessage(messageType, tempMsg); err != nil {
			logger.Error("消息发送出现错误", err.Error())
		}

	}, w.OnError, w.OnClose)
}

// OnError 客户端与服务端在消息交互过程中发生错误回调函数
func (w *Websocket) OnError(err error) {
	logger.Error("远端掉线、卡死、刷新浏览器等会触发该错误:", err.Error())
	//fmt.Printf("远端掉线、卡死、刷新浏览器等会触发该错误: %v\n", err.Error())
}

// OnClose 客户端关闭回调，发生onError回调以后会继续回调该函数
func (w *Websocket) OnClose() {
	w.Client.State = 0
	w.Client.Hub.UnRegister <- w.Client // 向hub管道投递一条注销消息，由hub中心负责关闭连接、删除在线数据
}

// 获取在线的全部客户端
func (w *Websocket) GetOnlineClients() int {
	//fmt.Printf("在线客户端数量：%d\n", len(w.WsClient.Hub.Clients))
	return len(w.Client.Hub.Clients)
}

// 向全部在线客户端广播消息
// 广播函数可能被不同的逻辑同时调用，由于操作的都是 Conn ，因此为了保证并发安全，加互斥锁
func (w *Websocket) BroadcastMsg(sendMsg string) {
	for onlineClient := range w.Client.Hub.Clients {
		//获取每一个在线的客户端，向远端发送消息
		if err := onlineClient.SendMessage(websocket.TextMessage, sendMsg); err != nil {
			logger.Error(ErrorsWebsocketWriteMgsFail, err.Error())
		}
	}
}
