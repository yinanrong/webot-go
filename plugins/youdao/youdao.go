package youdao

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"webot-go/service"

	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs" // 导入日志包
)

// Register plugin
// 必须有的插件注册函数
func Register(session *service.Session) {
	// 将插件注册到session
	session.HandlerRegister.Add(service.MSG_TEXT, service.Handler(youdao), "youdao")
}

// 消息处理函数
func youdao(session *service.Session, msg *service.ReceivedMessage) {
	uri := "http://fanyi.youdao.com/openapi.do?keyfrom=go-aida&key=145986666&type=data&doctype=json&version=1.1&q=" + url.QueryEscape(msg.Content)
	response, err := http.Get(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		logs.Error(err)
		return
	}
	errorCode, err := jc.GetInt("errorCode")
	if err != nil {
		logs.Error(err)
		return
	}
	if errorCode != 0 {
		logs.Error("youdao API", errorCode)
		return
	}
	trans, err := jc.GetSliceString("translation")
	if err != nil {
		logs.Error(err)
		return
	}
	if len(trans) < 1 {
		return
	}

	session.SendText(msg.At+trans[0], session.Bot.UserName, service.RealTargetUserName(session, msg))
}
