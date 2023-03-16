package world

import (
	"log"
	"snake_game_l1/internel/common"
	"time"
)

type player interface {
	GetShape() []common.Pixel
	SetDirection(d common.Direction)
	SetMoveInterval(moveInterval time.Duration)
	GetMoveInterval() time.Duration
	Eat()
	Run()
	Kill()
	IsDead() bool
}

type npc interface {
	GetShape() []common.Pixel
	SetMoveInterval(moveInterval time.Duration)
	GetMoveInterval() time.Duration
	BeEated()
	Run()
	IsDead() bool
}

type canvas interface {
	IsOut(box []common.Pixel) bool
	SetBox(name string, box []common.Pixel, color common.Color)
	SetCenterStr(str string)
	GetKeyboardEvent() chan common.Direction
	GetTermEvent() chan struct{}
	Draw()
	Init() error
	Destory()
}

const (
	GAME_STATUS_INIT = iota
	GAME_STATUS_RUNING
	GAME_STATUS_FAILED
	GAME_STATUS_WIN
)

type World struct {
	player     player
	npc        npc
	canvas     canvas
	winLen     int
	gameStatus int
}

func Create(p player, n npc, c canvas, winLen int) *World {
	return &World{
		player:     p,
		npc:        n,
		canvas:     c,
		winLen:     winLen,
		gameStatus: GAME_STATUS_INIT,
	}
}

func (w *World) getRefreshInterval() time.Duration {
	refreshInterval := w.player.GetMoveInterval()
	if w.player.GetMoveInterval() > w.npc.GetMoveInterval() {
		refreshInterval = w.npc.GetMoveInterval()
	}

	return refreshInterval / 2
}

func (w *World) isCrashPlayerNpc() bool {
	for _, ps := range w.player.GetShape() {
		for _, ns := range w.npc.GetShape() {
			if ps.X == ns.X && ps.Y == ns.Y {
				return true
			}
		}
	}
	return false
}

func (w *World) refresh() {
	w.player.Run()
	w.npc.Run()

	if w.isCrashPlayerNpc() {
		log.Println("player eat npc")
		w.npc.BeEated()
		w.player.Eat()
	}

	w.canvas.SetBox("player", w.player.GetShape(), common.COLOR_WHITE)
	w.canvas.SetBox("npc", w.npc.GetShape(), common.COLOR_GREEN)

	if w.canvas.IsOut(w.player.GetShape()) {
		w.player.Kill()
	}
	//判断游戏有没有结束
	if w.player.IsDead() {
		w.gameStatus = GAME_STATUS_FAILED
	} else if w.npc.IsDead() {
		w.gameStatus = GAME_STATUS_WIN
	} else if len(w.player.GetShape()) >= w.winLen {
		w.gameStatus = GAME_STATUS_WIN
	}

	if w.gameStatus == GAME_STATUS_FAILED {
		w.canvas.SetCenterStr("Game Over")
	} else if w.gameStatus == GAME_STATUS_WIN {
		w.canvas.SetCenterStr("You Win")
	}

	w.canvas.Draw()
}

func (w *World) Run() {

	w.gameStatus = GAME_STATUS_RUNING

	rInterval := w.getRefreshInterval()
	keyChan := w.canvas.GetKeyboardEvent()
	termChan := w.canvas.GetTermEvent()

	w.refresh()

	rTicker := time.NewTicker(rInterval)
	for {
		select {
		case <-rTicker.C:
			if w.gameStatus == GAME_STATUS_RUNING {
				w.refresh()
			}
		case dir := <-keyChan:
			w.player.SetDirection(dir)
		case <-termChan:
			w.canvas.Destory()
			rTicker.Stop()
			return
		}
	}
}
