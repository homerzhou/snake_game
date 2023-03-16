package snake

import (
	"snake_game_l1/internel/common"
	"time"
)

type Snake struct {
	moveInterval time.Duration
	p            []common.Pixel
	headDir      common.Direction
	tailDir      common.Direction
	lastMoveTime time.Time
	status       common.LiveStatus
}

func Create(initPixel []common.Pixel, moveInterval time.Duration) *Snake {
	if len(initPixel) <= 2 {
		return nil
	}

	f := Snake{
		moveInterval: moveInterval,
		p:            initPixel,
		headDir:      common.DIRECTION_RIGHT,
		status:       common.LIVE_STATUS_LIVING,
		lastMoveTime: time.Now(),
	}
	f.updateTailDir()
	return &f
}

func (s *Snake) updateTailDir() {

	tailBox := s.p[len(s.p)-1]
	tailBox2 := s.p[len(s.p)-2]
	if tailBox2.X == tailBox.X {
		if tailBox2.Y > tailBox.Y {
			s.tailDir = common.DIRECTION_UP
		} else {
			s.tailDir = common.DIRECTION_DOWN
		}
	} else {
		if tailBox2.X > tailBox.X {
			s.tailDir = common.DIRECTION_RIGHT
		} else {
			s.tailDir = common.DIRECTION_LEFT
		}
	}

}

func (s *Snake) crawl() {
	for i := len(s.p) - 1; i > 0; i-- {
		s.p[i].X = s.p[i-1].X
		s.p[i].Y = s.p[i-1].Y
	}

	switch s.headDir {
	case common.DIRECTION_UP:
		s.p[0].Y++
	case common.DIRECTION_DOWN:
		s.p[0].Y--
	case common.DIRECTION_LEFT:
		s.p[0].X--
	case common.DIRECTION_RIGHT:
		s.p[0].X++
	}

	s.updateTailDir()
}

func (s *Snake) selfCrash() bool {
	selfCrash := false
	for i := range s.p {
		for j := i + 1; j < len(s.p); j++ {
			if s.p[i].X == s.p[j].X &&
				s.p[i].Y == s.p[j].Y {
				selfCrash = true
			}
		}
	}
	return selfCrash
}

func (s *Snake) Run() {
	if s.status != common.LIVE_STATUS_LIVING {
		return
	}

	for {
		if time.Now().After(s.lastMoveTime.Add(s.moveInterval)) {
			s.crawl()
			if s.selfCrash() {
				s.status = common.LIVE_STATUS_DEAD
				break
			}
			s.lastMoveTime = s.lastMoveTime.Add(s.moveInterval)
		} else {
			break
		}
	}
}

func (s *Snake) Eat() {
	if s.status != common.LIVE_STATUS_LIVING {
		return
	}

	tailP := s.p[len(s.p)-1]
	switch s.tailDir {
	case common.DIRECTION_UP:
		s.p = append(s.p, common.Pixel{X: tailP.X, Y: tailP.Y - 1})
	case common.DIRECTION_DOWN:
		s.p = append(s.p, common.Pixel{X: tailP.X, Y: tailP.Y + 1})
	case common.DIRECTION_LEFT:
		s.p = append(s.p, common.Pixel{X: tailP.X + 1, Y: tailP.Y})
	case common.DIRECTION_RIGHT:
		s.p = append(s.p, common.Pixel{X: tailP.X - 1, Y: tailP.Y})
	}
}

func (s *Snake) Kill() {
	s.status = common.LIVE_STATUS_DEAD
}

func (s *Snake) IsDead() bool {
	return s.status == common.LIVE_STATUS_DEAD
}

func (s *Snake) GetShape() []common.Pixel {
	return s.p
}

func (s *Snake) SetDirection(d common.Direction) {
	if s.status != common.LIVE_STATUS_LIVING {
		return
	}

	if s.headDir == common.DIRECTION_DOWN ||
		s.headDir == common.DIRECTION_UP {
		if d == common.DIRECTION_LEFT ||
			d == common.DIRECTION_RIGHT {

			s.headDir = d
		}
	} else {
		if d == common.DIRECTION_DOWN ||
			d == common.DIRECTION_UP {

			s.headDir = d
		}
	}

}

func (s *Snake) SetMoveInterval(moveInterval time.Duration) {
	if s.status != common.LIVE_STATUS_LIVING {
		return
	}
	s.moveInterval = moveInterval
}

func (s *Snake) GetMoveInterval() time.Duration {
	return s.moveInterval
}
