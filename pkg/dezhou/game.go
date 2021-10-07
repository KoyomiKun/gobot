package dezhou

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	RoundFormat = "第%d轮开始"
)

type Order struct {
	msg    string
	fromQQ int64
}

type Game struct {
	BankerIndex   int
	WaitedPlayer  int64
	ShuffledCards Cards
	ShowCards     []Card
	JoinedPlayer  []*Player
	OrderChan     chan Order
	ctx           *zero.Ctx
}

func NewGame(ctx *zero.Ctx) *Game {

	groupId := ctx.Event.GroupID

	// shuffle cards
	var cardsContent []Card
	copy(cardsContent, cards.Content)
	copyCards := Cards{
		TopIndex: cards.TopIndex,
		Content:  cardsContent,
	}

	copyCards.Shuffle()

	userIds, ok := groupMutex.Load(groupId)
	if !ok {
		log.Errorf("fail create game in %d: group has not registried", groupId)
		ctx.Send("创建游戏失败，请联系管理员处理")
		return nil
	}
	userIdArr := userIds.([]int64)

	// init player
	joinedPlayers := make([]*Player, 0, len(userIdArr))
	for _, v := range userIdArr {
		joinedPlayers = append(joinedPlayers, &Player{
			Counter:      1000,
			TableCounter: 0,
			QQid:         v,
			Droped:       false,
			PCards: Cards{
				TopIndex: 2,
				Content:  make([]Card, 0, 3),
			},
			OwnToCtx: ctx,
		})
	}

	ctx.Send("初始化游戏完毕，游戏即将开始...")

	return &Game{
		BankerIndex:   0,
		WaitedPlayer:  0,
		ShuffledCards: copyCards,
		ShowCards:     make([]Card, 0, 4),
		JoinedPlayer:  joinedPlayers,
		OrderChan:     make(chan Order),
		ctx:           ctx,
	}
}

func (g *Game) Start() {
	gctx := g.ctx
	gctx.Send("游戏开始！")
	for i := 0; len(g.JoinedPlayer) != 0; i++ {
		g.Round(i)
	}
	// delete group mutex
}

func (g *Game) Round(index int) {
	playerNum := len(g.JoinedPlayer)

	banker := g.JoinedPlayer[g.BankerIndex%playerNum].QQid
	smallBlind := g.JoinedPlayer[(g.BankerIndex+1)%playerNum].QQid
	bigBlind := g.JoinedPlayer[(g.BankerIndex+2)%playerNum].QQid

	g.ctx.Send(message.Message{
		message.Text(fmt.Sprintf(RoundFormat, index)),
		message.At(banker),
		message.Text("是庄家"),
		message.At(smallBlind),
		message.Text("是小盲"),
		message.At(bigBlind),
		message.Text("是大盲"),
	})
	g.BankerIndex++

	// send cards
	for _, player := range g.JoinedPlayer {
		var sb strings.Builder
		for i := 0; i < 3; i++ {
			card := g.ShuffledCards.Pop()
			sb.WriteString(card.String())
		}
		g.ctx.SendPrivateMessage(player.QQid, message.Text(sb.String()))
	}

	for !(len(g.ShowCards) == 4) {
		// round start
		var i int = 3
		for g.WaitedPlayer != bigBlind {
			g.WaitedPlayer = g.JoinedPlayer[(g.BankerIndex+i)%playerNum].QQid
			g.ctx.Send(message.Message{
				message.At(g.WaitedPlayer),
				message.Text("说话"),
			})
			select {
			// order chan must rec target player's msg
			// TODO: strict speed <06-10-21, komikun> //
			case order := <-g.OrderChan:
				switch order.msg {
				case "跟":
					g.ctx.Send("不准跟")
				case "加20":
					g.ctx.Send("不准加")
				case "溜":
					g.ctx.Send("不准溜")
				}
			}
			i++
		}
		// show card
		showCard := *g.ShuffledCards.Pop()
		g.ShowCards = append(g.ShowCards, showCard)
		var sb strings.Builder
		sb.WriteString("当前展示的卡牌为：\n")
		for _, card := range g.ShowCards {
			sb.WriteString(card.String())
		}
		g.ctx.Send(sb.String())
	}
}

func (g *Game) Close() {
	g.ctx.Send("游戏关闭...")
	groupMutex.Delete(g.ctx.Event.GroupID)
}
