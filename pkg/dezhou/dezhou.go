package dezhou

import (
	"sync"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	cards      *Cards
	groupMutex sync.Map
)

func init() {
	cards = NewCards()
	zero.OnCommandGroup([]string{"德州", "111"}).Handle(func(ctx *zero.Ctx) {
		if !zero.OnlyGroup(ctx) {
			ctx.Send(message.Text("该功能只有群组可以使用"))
			return
		}

		switch ctx.Event.Message.String() {
		case ".德州":
			if _, ok := groupMutex.Load(ctx.Event.GroupID); ok {
				ctx.Send("该群组已有游戏在进行")
				return
			}
			groupMutex.Store(ctx.Event.GroupID, map[int64]struct{}{})
			ctx.Send("30s内输入.111报名...")
			go func() {
				<-time.After(time.Second * 30)
				if v, ok := groupMutex.Load(ctx.Event.GroupID); ok {
					userSet := v.(map[int64]struct{})
					if len(userSet) <= 3 {
						ctx.Send("30s少于4人加入，游戏关闭")
						groupMutex.Delete(ctx.Event.GroupID)
					} else {
						retMsg := message.Message{}
						retMsg = append(retMsg, message.Text("游戏开始！参与者为"))
						userList := make([]int64, 0, len(userSet))
						for user := range userSet {
							retMsg = append(retMsg, message.At(user))
							userList = append(userList, user)
						}
						ctx.Send(retMsg)
						groupMutex.Delete(ctx.Event.GroupID)
						groupMutex.Store(ctx.Event.GroupID, userList)
						game := NewGame(ctx)
						groupMutex.Delete(ctx.Event.GroupID)
						groupMutex.Store(ctx.Event.GroupID, game)
						game.Start()
					}
				}
			}()
		case ".111":
			if v, ok := groupMutex.Load(ctx.Event.GroupID); !ok {
				ctx.Send(message.Text("请先输入.德州开始游戏"))
				return
			} else {
				switch v.(type) {
				case []int64:
					ctx.Send(message.Text("当前并非等待状态"))
					return
				case map[int64]struct{}:
					// set sender wait
					v.(map[int64]struct{})[ctx.Event.UserID] = struct{}{}
					ctx.Send(
						message.Message{
							message.At(ctx.Event.UserID),
							message.Text("加入成功"),
						})
				}
			}
		}
	})
	zero.OnCommandGroup([]string{"跟", "加20", "溜"}).Handle(func(ctx *zero.Ctx) {
		if !zero.OnlyGroup(ctx) {
			ctx.Send(message.Text("该功能只有群组可以使用"))
			return
		}
		gameInterface, ok := groupMutex.Load(ctx.Event.GroupID)
		if !ok {
			ctx.Send("请先输入.德州开始游戏")
			return
		}
		switch gameInterface.(type) {
		case []int64:
		case map[int64]struct{}:
			ctx.Send("游戏尚未初始化")
			return
		}

		game := gameInterface.(Game)
		if ctx.Event.UserID != game.WaitedPlayer {
			return
		}
		game.OrderChan <- Order{
			msg:    ctx.Event.Message.String()[1:],
			fromQQ: ctx.Event.UserID,
		}
	})
	zero.OnCommandGroup([]string{"当前筹码", "当前明牌"}).Handle(func(ctx *zero.Ctx) {
		if !zero.OnlyGroup(ctx) {
			ctx.Send(message.Text("该功能只有群组可以使用"))
			return
		}
		gameInterface, ok := groupMutex.Load(ctx.Event.GroupID)
		if !ok {
			ctx.Send("请先输入.德州开始游戏")
			return
		}
		switch gameInterface.(type) {
		case []int64:
		case map[int64]struct{}:
			ctx.Send("游戏尚未初始化")
			return
		}
		// show info
	})
}
