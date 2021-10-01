package hhsh

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
)

const (
	baseUrl = "https://lab.magiconch.com/api/nbnhhsh/guess"
)

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
		req, _ := http.NewRequest(http.MethodPost, baseUrl, nil)

	})
}
