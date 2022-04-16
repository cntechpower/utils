package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var tokenMap = make(map[string]*appAccessToken)

func getAppAccessToken(corpId, corpSecret string) (accessToken string, err error) {
	token, ok := tokenMap[fmt.Sprintf("%v:%v", corpId, corpSecret)]
	if ok && token.ExpiresAt.Unix() > time.Now().Unix() {
		accessToken = token.AccessToken
		return
	}
	resp, err := http.Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpId, corpSecret))
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	tokenResp := &genAccessTokenResp{}
	err = json.Unmarshal(bs, tokenResp)
	if err != nil {
		return
	}
	if tokenResp.ErrCode != 0 {
		err = fmt.Errorf("call token api error: code:%v, msg:%v", tokenResp.ErrCode, tokenResp.ErrMsg)
		return
	}
	token = &appAccessToken{
		AccessToken: tokenResp.AccessToken,
		ExpiresAt:   time.Now().Add(time.Second * time.Duration(tokenResp.ExpiresIn)),
	}
	accessToken = token.AccessToken
	tokenMap[fmt.Sprintf("%v:%v", corpId, corpSecret)] = token
	return
}

func SendMsgByApp(ctx context.Context, corpId, corpSecret string,
	toUser []string, toParty []string, toTag []string,
	agentId int64, content string) (err error) {
	body := strings.NewReader(fmt.Sprintf(`
{
   "touser" : "%s",
   "toparty" : "%s",
   "totag" : "%s",
   "msgtype" : "text",
   "agentid" : %d,
   "text" : {
       "content" : "%s"
   },
   "safe":0,
   "enable_id_trans": 0,
   "enable_duplicate_check": 0,
   "duplicate_check_interval": 1800
}`, strings.Join(toUser, "|"), strings.Join(toParty, "|"), strings.Join(toTag, "|"), agentId, content))

	token, err := getAppAccessToken(corpId, corpSecret)
	if err != nil {
		return
	}
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token), body)
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
		err = fmt.Errorf("send msg call fail, code: %v, body: %v", resp.StatusCode, string(bs))
		return
	}
	return
}
