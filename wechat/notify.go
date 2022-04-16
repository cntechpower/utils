package wechat

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	MsgTypeText     = "text"
	MsgTypeMarkdown = "markdown"
)

func buildStringList(s []string) (res string) {
	builder := strings.Builder{}
	builder.WriteString("[")
	l := len(s)
	for idx, str := range s {
		builder.WriteString(fmt.Sprintf("%q", str))
		if idx != l-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("]")
	return builder.String()
}

func SendMsgByGroupRobot(ctx context.Context, webhook, typ, content string, mentionedUserList, mentionedMobileList []string) (err error) {
	var body io.Reader
	if typ == MsgTypeText {
		body = strings.NewReader(fmt.Sprintf(`
{
    "msgtype": "%s",
    "text": {
        "content": "%s",
		"mentioned_list":%v,
		"mentioned_mobile_list":%v
    }
}`,
			typ, content, buildStringList(mentionedUserList), buildStringList(mentionedMobileList)))
	} else {
		body = strings.NewReader(fmt.Sprintf(`
{
    "msgtype": "%s",
    "markdown": {
        "content": "%s"
    }
}`,
			typ, content))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", webhook, body)
	if err != nil {
		return
	}
	defer func() {
		_ = req.Body.Close()
	}()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != 200 {
		bs, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("webbook call fail, code: %v, body: %v", resp.StatusCode, string(bs))
		return
	}
	//markdown not support mention user, so send another text message.
	if typ == MsgTypeMarkdown &&
		(len(mentionedUserList) != 0 || len(mentionedMobileList) != 0) {
		err = SendMsgByGroupRobot(ctx, webhook, MsgTypeText, "", mentionedUserList, mentionedMobileList)
	}
	return
}
