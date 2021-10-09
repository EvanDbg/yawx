package wx

import (
	"time"

	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

var wx = core.NewBucket("wx")

var api_url = wx.Get("api_url")

func init() {
	core.Server.POST("/yawx", func(c *gin.Context) {
		data, _ := c.GetRawData()
		core.NotifyMasters(string(data))
	})
}

type Message struct {
	Type      int    `json:"type"`
	Msg       string `json:"msg"`
	FromWxid  string `json:"from_wxid"`
	RobotWxid string `json:"robot_wxid"`
}

type Sender struct {
	matches [][]string
	deleted bool
	goon    bool
}

func (sender *Sender) GetContent() string {
	return ""
}

func (sender *Sender) GetUserID() int {
	return 0
}

func (sender *Sender) GetChatID() int {
	return 0
}

func (sender *Sender) GetImType() string {
	return "wx"
}

func (sender *Sender) GetMessageID() int {
	return 0
}

func (sender *Sender) GetUsername() string {
	return ""
}

func (sender *Sender) IsReply() bool {
	return false
}

func (sender *Sender) GetReplySenderUserID() int {
	if !sender.IsReply() {
		return 0
	}
	return 0
}

func (sender *Sender) GetRawMessage() interface{} {
	return nil
}

func (sender *Sender) SetMatch(ss []string) {
	sender.matches = [][]string{ss}
}
func (sender *Sender) SetAllMatch(ss [][]string) {
	sender.matches = ss
}

func (sender *Sender) GetMatch() []string {
	return sender.matches[0]
}

func (sender *Sender) GetAllMatch() [][]string {
	return sender.matches
}

func (sender *Sender) Get(index ...int) string {

	i := 0
	if len(index) != 0 {
		i = index[0]
	}
	if len(sender.matches) == 0 {
		return ""
	}
	if len(sender.matches[0]) < i+1 {
		return ""
	}
	return sender.matches[0][i]
}

func (sender *Sender) IsAdmin() bool {

	return true
}

func (sender *Sender) IsMedia() bool {
	return false
}

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {

	return 0, nil
}

func (sender *Sender) Delete() error {
	return nil
}

func (sender *Sender) Disappear(lifetime ...time.Duration) {

}

func (sender *Sender) Finish() {

}

func (sender *Sender) Continue() {
	sender.goon = true
}

func (sender *Sender) IsContinue() bool {
	return sender.goon
}
