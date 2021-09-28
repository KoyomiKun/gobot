package searcher

import zero "github.com/wdvxdr1123/ZeroBot"

func init() {
	zero.OnCommand("看图").Handle(func(ctx *zero.Ctx) {
		ctx.Send("不行")
	})
}
