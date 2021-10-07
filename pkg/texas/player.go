package texas

import zero "github.com/wdvxdr1123/ZeroBot"

type Player struct {
	Counter      int
	TableCounter int
	QQid         int64
	Droped       bool
	PCards       Cards
	OwnToCtx     *zero.Ctx
}

func (p *Player) Drop() {
	p.Droped = true
}

func (p *Player) AddBonus(n int) {
	if n != 20 && n != 50 && n != 100 {
		p.OwnToCtx.Send("只能加20/50/100")
		return
	}
	p.Counter -= n
	p.TableCounter += n
}

func (p *Player) SmallBlind() {
	p.Counter -= 10
	p.TableCounter += 10
}

func (p *Player) BigBlind() {
	p.Counter -= 20
	p.TableCounter += 20
}

func (p *Player) Follow(target int) {
	if target < p.TableCounter {
		p.OwnToCtx.Send("不用跟")
		return
	}
	p.Counter -= (target - p.TableCounter)
	p.TableCounter = target
}
