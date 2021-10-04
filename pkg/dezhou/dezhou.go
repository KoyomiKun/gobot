package dezhou

import (
	"sync"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	// value: nil 0 1
	// 0: wait for players
	// 1: start game
	groupMutex sync.Map
)

func init() {
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
					if len(userSet) == 0 {
						ctx.Send("30s无人加入，游戏关闭")
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

}
