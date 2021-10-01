package searcher

import (
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
)

func init() {
	zero.OnCommand("查图").Handle(func(ctx *zero.Ctx) {
		var cmd extension.CommandModel
		err := ctx.Parse(&cmd)
		if err != nil {
			log.Errorf("fail parse ctx: %s", ctx.Event.RawMessage)
			return
		}
	})
}
