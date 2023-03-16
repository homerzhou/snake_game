package food

import (
	"snake_game_l1/internel/common"
	"time"
)

type RandomInterface interface {
	Random(int) ([]common.Pixel, error)
}

type Food struct {
	p             []common.Pixel
	status        common.LiveStatus
	moveInterval  time.Duration
	lastMoveTime  time.Time
	randomHandler RandomInterface
}

func Create(initPixel common.Pixel, moveInterval time.Duration, randomHandler RandomInterface) *Food {
	f := Food{
		p:             []common.Pixel{initPixel},
		moveInterval:  moveInterval,
		status:        common.LIVE_STATUS_LIVING,
		lastMoveTime:  time.Now(),
		randomHandler: randomHandler,
	}
	return &f
}

func (f *Food) Run() {
	if f.status == common.LIVE_STATUS_DEAD {
		return
	}

	for {
		if time.Now().After(f.lastMoveTime.Add(f.moveInterval)) {
			f.BeEated()
			f.lastMoveTime = f.lastMoveTime.Add(f.moveInterval)
		} else {
			break
		}
	}
}

func (f *Food) BeEated() {
	p, err := f.randomHandler.Random(1)
	if err != nil || len(p) <= 0 {
		f.p = nil
		f.status = common.LIVE_STATUS_DEAD
		return
	}
	f.p = []common.Pixel{p[0]}
}

func (f *Food) IsDead() bool {
	return f.status == common.LIVE_STATUS_DEAD
}

func (f *Food) GetShape() []common.Pixel {
	return f.p
}

func (f *Food) SetMoveInterval(moveInterval time.Duration) {
	f.moveInterval = moveInterval
}
func (f *Food) GetMoveInterval() time.Duration {
	return f.moveInterval
}
