package wx

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/httplib"
	"github.com/cdle/sillyGirl/core"
	"github.com/gin-gonic/gin"
)

var wx = core.NewBucket("wx")

var api_url = wx.Get("api_url")
var robot_wxid = wx.Get("robot_wxid")

func sendMsg(pmsg *PushMsg) {
	if pmsg.Msg == "" {
		return
	}
	if pmsg.Type == 0 {
		pmsg.Type = 1
	}
	pmsg.RobotWxid = robot_wxid
	req := httplib.Post(api_url)
	pmsg.Msg = url.QueryEscape(pmsg.Msg)
	data, _ := json.Marshal(pmsg)
	data, _ = json.Marshal(map[string]string{
		"data": string(data),
	})
	req.Header("Content-Type", "application/json")
	req.Body(data)
	req.Response()
}

func init() {
	core.Pushs["wx"] = func(i interface{}, s string) {
		if robot_wxid != "" {
			pmsg := PushMsg{
				Msg:      s,
				FromWxid: fmt.Sprint(i),
			}
			sendMsg(&pmsg)
		}
	}
	core.GroupPushs["wx"] = func(i, _ interface{}, s string) {
		pmsg := PushMsg{
			Type:      1,
			FromWxid:  fmt.Sprint(i) + "@chatroom",
			RobotWxid: robot_wxid,
		}
		for _, v := range regexp.MustCompile(`\[CQ:image,file=([^\[\]]+)\]`).FindAllStringSubmatch(s, -1) {
			s = strings.Replace(s, fmt.Sprintf(`[CQ:image,file=%s]`, v[1]), "", -1)
			data, err := os.ReadFile(core.ExecPath + "/data/images/" + v[1])
			if err == nil {
				add := regexp.MustCompile("(https.*)").FindString(string(data))
				if add != "" {
					pmsg := pmsg
					pmsg.Type = 3
					pmsg.Msg = add
					sendMsg(&pmsg)
				}
			}
		}
		s = regexp.MustCompile(`\[CQ:([^\[\]]+)\]`).ReplaceAllString(s, "")
		// core.NotifyMasters(pmsg.Msg)
		// {
		// 	t := []string{}
		// 	for _, v := range strings.Split(s, "\n") {
		// 		if v != "" {
		// 			t = append(t, v)
		// 		}
		// 	}
		// 	s = strings.Join(t, "\n")
		// }
		// core.NotifyMasters(pmsg.Msg)
		s = strings.Replace(s, "\n\n", "\n", -1)
		s = strings.Trim(s, "\n")
		pmsg.Msg = s
		sendMsg(&pmsg)
	}
	core.Server.POST("/yawx", func(c *gin.Context) {
		data, _ := c.GetRawData()
		s, err := url.QueryUnescape(string(data))
		if err != nil {
			return
		}
		args, err := url.ParseQuery(s)
		if err != nil {
			return
		}
		if args.Get("type") == "" {
			return
		}
		if args.Get("robot_wxid") != robot_wxid {
			robot_wxid = args.Get("robot_wxid")
			wx.Set("robot_wxid", robot_wxid)
		}
		core.Senders <- &Sender{
			value: args,
		}
	})
}

type Sender struct {
	leixing int
	mtype   int
	matches [][]string
	deleted bool
	goon    bool
	value   url.Values
}

type JsonMsg struct {
	Content         string `json:"content"`
	FinalFromName   string `json:"final_from_name"`
	FinalFromWxid   string `json:"final_from_wxid"`
	FromName        string `json:"from_name"`
	FromWxid        string `json:"from_wxid"`
	MsgType         int    `json:"msg_type"`
	Msgid           int    `json:"msgid"`
	OriginalContent string `json:"original_content"`
	SendOutType     int    `json:"send_out_type"`
	Timestamp       int    `json:"timestamp"`
	ToName          string `json:"to_name"`
	ToWxid          string `json:"to_wxid"`
}

type PushMsg struct {
	Type      int    `json:"type"`
	Msg       string `json:"msg"`
	FromWxid  string `json:"from_wxid"`
	RobotWxid string `json:"robot_wxid"`
}

func (sender *Sender) GetContent() string {
	return sender.value.Get("msg")
}

func (sender *Sender) GetUserID() interface{} {
	if uid := sender.value.Get("final_from_wxid"); uid != "" {
		return uid
	} else {
		return sender.value.Get("from_wxid")
	}
}

func (sender *Sender) GetChatID() interface{} {
	if uid := sender.value.Get("final_from_wxid"); uid != "" {
		return strings.Replace(sender.value.Get("from_wxid"), "@chatroom", "", -1)
	} else {
		return nil
	}
}

func (sender *Sender) GetImType() string {
	return "wx"
}

func (sender *Sender) GetMessageID() int {
	return core.Int(sender.value.Get("msgid"))
}

func (sender *Sender) GetUsername() string {
	if uid := sender.value.Get("final_from_wxid"); uid != "" {
		return sender.value.Get("final_from_name")
	} else {
		return sender.value.Get("from_name")
	}
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
	return strings.Contains(wx.Get("masters"), fmt.Sprint(sender.GetUserID()))
}

func (sender *Sender) IsMedia() bool {
	return false
}

func (sender *Sender) Reply(msgs ...interface{}) (int, error) {
	pmsg := PushMsg{
		FromWxid:  sender.value.Get("from_wxid"),
		RobotWxid: robot_wxid,
	}
	for _, item := range msgs {
		switch item.(type) {
		case string:
			pmsg.Msg = item.(string)
		case []byte:
			pmsg.Msg = string(item.([]byte))
		case core.ImageUrl:
			pmsg.Type = 3
			pmsg.Msg = string(item.(core.ImageUrl))
		}
	}
	sendMsg(&pmsg)
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
