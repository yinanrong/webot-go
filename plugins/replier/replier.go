package replier

import (
	"net/http"
	"net/url"
	"strings"

	"io/ioutil"

	"fmt"

	"webot-go/service"

	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs"
)

// Register plugin
func Register(session *service.Session) {
	session.HandlerRegister.Add(service.MSG_TEXT, service.Handler(autoReply), "text-replier")
	if err := session.HandlerRegister.Add(service.MSG_IMG, service.Handler(autoReply), "img-replier"); err != nil {
		logs.Error(err)
	}

}
func autoReply(session *service.Session, msg *service.ReceivedMessage) {
	if msg.IsBot(session) {
		return
	}
	uri := fmt.Sprintf("http://api.qingyunke.com/api.php?key=free&appid=0&msg=%s", url.QueryEscape(msg.Content))
	resp, err := http.Get(uri)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	jc, err := rrconfig.LoadJsonConfigFromBytes(body)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	content, err := jc.GetString("content")
	if err != nil {
		logs.Error(err.Error())
		return
	}
	session.SendText(strings.Replace(content, "{br}", "\n", -1), session.Bot.UserName, msg.FromUserName)
}
