package hhsh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	baseUrl     = "https://lab.magiconch.com/api/nbnhhsh/guess"
	transFormat = "%s : %s \n"
)

type Tran struct {
	Name  string
	Trans []string
}
type Resp []*Tran

func init() {
	zero.OnCommand("好好说话").Handle(func(ctx *zero.Ctx) {
		cmd := extension.CommandModel{}
		err := ctx.Parse(&cmd)
		if err != nil {
			log.Errorf("fail parse ctx: %s", ctx.Event.RawMessage)
			return
		}
		if cmd.Args == "" {
			ctx.Send("无参数")
			return
		}

		bodyLoad := fmt.Sprintf(`{"text":"%s"}`, cmd.Args)
		req, _ := http.NewRequest(http.MethodPost, baseUrl, bytes.NewBufferString(bodyLoad))
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Errorf("fail getting hhsh url %s: %v", baseUrl, err)
			ctx.Send(message.Text("抓取结果失败"))
			return
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		respRes := Resp{}
		err = json.Unmarshal(respBody, &respRes)
		if err != nil {
			log.Errorf("fail parse hhsh json: %v", err)
			ctx.Send("解析失败，请联系管理员")
			return
		}
		var retSb strings.Builder
		for _, v := range respRes {
			retSb.WriteString(fmt.Sprintf(transFormat, v.Name, strings.Join(v.Trans, ",")))
		}
		ctx.Send(message.Text(retSb.String()))
	})
}
