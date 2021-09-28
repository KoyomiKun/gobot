package setu

import zero "github.com/wdvxdr1123/ZeroBot"

func init() {
	zero.OnCommand("色图").Handle(func(ctx *zero.Ctx) {
		ctx.Send("滚")
	})
}
