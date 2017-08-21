package verify

import (
	"fmt"

	"webot-go/service"

	"github.com/songtianyi/rrframework/logs"
)

// Register plugin
// 必须有的插件注册函数
// 指定session, 可以对不同用户注册不同插件
func Register(session *service.Session) {
	// 将插件注册到session
	// 第一个参数: 指定消息类型, 所有该类型的消息都会被转发到此插件
	// 第二个参数: 指定消息处理函数, 消息会进入此函数
	// 第三个参数: 自定义插件名，不能重名，switcher插件会用到此名称
	session.HandlerRegister.Add(service.MSG_FV, service.Handler(verify), "verify")
}

// 消息处理函数
func verify(session *service.Session, msg *service.ReceivedMessage) {

	logs.Info(msg.Content)

	master := session.Cm.GetContactByPYQuanPin("yinanrong")

	if err := session.AcceptFriend("", []*service.VerifyUser{{Value: msg.RecommendInfo.UserName, VerifyUserTicket: msg.RecommendInfo.Ticket}}); err != nil {
		errMsg := fmt.Sprintf("accept %s's friend request error, %s", msg.RecommendInfo.NickName, err.Error())
		logs.Error(errMsg)
		if master != nil {
			session.SendText(errMsg, session.Bot.UserName, master.UserName)
		}
		return
	}

	// 回复消息
	// 第一个参数: 回复的内容
	// 第二个参数: 机器人ID
	// 第三个参数: 联系人/群组/特殊账号ID
	if master != nil {
		session.SendText(fmt.Sprintf("%s accepted", msg.RecommendInfo.NickName), session.Bot.UserName, master.UserName)
	}

}
