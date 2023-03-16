package world

import (
	"fmt"
	"snake_game_l1/internel/common"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type StubPlayer struct {
	shapeListIdx int
	shapeList    [][]common.Pixel
	moveInterval time.Duration
	lastMoveTime time.Time
	dir          common.Direction
	dead         bool
	eatCount     int
}

func (s *StubPlayer) GetShape() []common.Pixel { return s.shapeList[s.shapeListIdx] }
func (s *StubPlayer) SetDirection(d common.Direction) {
	s.dir = d
}
func (s *StubPlayer) SetMoveInterval(moveInterval time.Duration) {}
func (s *StubPlayer) GetMoveInterval() time.Duration             { return s.moveInterval }
func (s *StubPlayer) Eat()                                       { s.eatCount++ }
func (s *StubPlayer) Run() {
	for {
		if time.Now().After(s.lastMoveTime.Add(s.moveInterval)) {
			fmt.Println("player run")
			s.shapeListIdx++
			s.lastMoveTime = s.lastMoveTime.Add(s.moveInterval)
		} else {
			break
		}
	}
}
func (s *StubPlayer) Kill()        { s.dead = true }
func (s *StubPlayer) IsDead() bool { return s.dead }

type StubNpc struct {
	shapeListIdx int
	shapeList    [][]common.Pixel
	moveInterval time.Duration
	dead         bool
	beEatedCount int
	lastMoveTime time.Time
}

func (s *StubNpc) GetShape() []common.Pixel                   { return s.shapeList[s.shapeListIdx] }
func (s *StubNpc) SetMoveInterval(moveInterval time.Duration) {}
func (s *StubNpc) GetMoveInterval() time.Duration             { return s.moveInterval }
func (s *StubNpc) BeEated()                                   { s.beEatedCount++ }
func (s *StubNpc) Run() {
	for {
		if time.Now().After(s.lastMoveTime.Add(s.moveInterval)) {
			fmt.Println("npc run")
			s.shapeListIdx++
			s.lastMoveTime = s.lastMoveTime.Add(s.moveInterval)
		} else {
			break
		}
	}
}
func (s *StubNpc) IsDead() bool { return s.dead }

type StubCanvas struct {
	maxWidth     int
	maxHeight    int
	boxMap       map[string][]common.Pixel
	keyChan      chan common.Direction
	termChan     chan struct{}
	drawCount    int
	initCount    int
	destoryCount int
}

func (s *StubCanvas) IsOut(box []common.Pixel) bool {
	for _, b := range box {
		if b.X > s.maxWidth || b.Y > s.maxHeight {
			return true
		}
	}
	return false
}
func (s *StubCanvas) SetBox(name string, box []common.Pixel, color common.Color) {
	s.boxMap[name] = box
}
func (s *StubCanvas) SetCenterStr(str string)                 {}
func (s *StubCanvas) GetKeyboardEvent() chan common.Direction { return s.keyChan }
func (s *StubCanvas) GetTermEvent() chan struct{}             { return s.termChan }
func (s *StubCanvas) Draw() {
	s.drawCount++
}
func (s *StubCanvas) Init() error {
	s.initCount++
	return nil
}
func (s *StubCanvas) Destory() {
	s.destoryCount++
}

func TestIntervalAndCrash(t *testing.T) {
	tests := []struct {
		inputPlayer    *StubPlayer
		inputNpc       *StubNpc
		expectInterval time.Duration
		expectIsCash   bool
	}{
		{
			inputPlayer: &StubPlayer{
				shapeListIdx: 0,
				shapeList:    [][]common.Pixel{{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}}},
				moveInterval: 100 * time.Millisecond,
				lastMoveTime: time.Now(),
				dead:         false,
			},
			inputNpc: &StubNpc{
				shapeListIdx: 0,
				shapeList:    [][]common.Pixel{{{X: 12, Y: 10}}},
				moveInterval: 50 * time.Millisecond,
				dead:         false,
				lastMoveTime: time.Now(),
			},
			expectInterval: 25 * time.Millisecond,
			expectIsCash:   true,
		},
		{
			inputPlayer: &StubPlayer{
				shapeListIdx: 0,
				shapeList:    [][]common.Pixel{{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}}},
				moveInterval: 100 * time.Millisecond,
				lastMoveTime: time.Now(),
				dead:         false,
			},
			inputNpc: &StubNpc{
				shapeListIdx: 0,
				shapeList:    [][]common.Pixel{{{X: 13, Y: 10}}},
				moveInterval: 50 * time.Millisecond,
				dead:         false,
				lastMoveTime: time.Now(),
			},
			expectInterval: 25 * time.Millisecond,
			expectIsCash:   false,
		},
	}

	for idx, test := range tests {
		w := Create(test.inputPlayer, test.inputNpc, nil, 0)
		assert.Equal(t, test.expectInterval, w.getRefreshInterval(), strconv.Itoa(idx))
		assert.Equal(t, test.expectIsCash, w.isCrashPlayerNpc(), strconv.Itoa(idx))
	}
}

func TestRefresh(t *testing.T) {
	tests := []struct {
		inputPlayer           *StubPlayer
		inputNpc              *StubNpc
		inputCanvas           *StubCanvas
		inputWinLen           int
		inputLoopTimes        int
		inputInterval         time.Duration
		expectPlayerIdx       int
		expectPlayerShape     []common.Pixel
		expectPlayerIsDead    bool
		expectPlayerEatCount  int
		expectNpcIdx          int
		expectNpcShape        []common.Pixel
		expectNpcIsDead       bool
		expectNpcBeEatedCount int
		expectCanvasBox       map[string][]common.Pixel
		expectCanvasDrawCount int
		expectGameStatus      int
	}{
		{
			//player正常的走一步， npc正常的走一次，
			//npc和player没有发生碰撞
			//player没有出界
			//player长度没有达到赢的长度
			inputPlayer: &StubPlayer{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
					{{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
				},
				moveInterval: 100 * time.Millisecond,
				lastMoveTime: time.Now(),
				dead:         false,
			},
			inputNpc: &StubNpc{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 15, Y: 10}},
					{{X: 23, Y: 10}},
				},
				moveInterval: 100 * time.Millisecond,
				dead:         false,
				lastMoveTime: time.Now(),
			},
			inputCanvas: &StubCanvas{
				maxWidth:  30,
				maxHeight: 30,
				boxMap:    make(map[string][]common.Pixel),
				drawCount: 0,
			},
			inputWinLen:           10,
			inputLoopTimes:        2,
			inputInterval:         80 * time.Millisecond,
			expectPlayerIdx:       1,
			expectPlayerShape:     []common.Pixel{{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
			expectPlayerIsDead:    false,
			expectPlayerEatCount:  0,
			expectNpcIdx:          1,
			expectNpcShape:        []common.Pixel{{X: 23, Y: 10}},
			expectNpcIsDead:       false,
			expectNpcBeEatedCount: 0,
			expectCanvasBox: map[string][]common.Pixel{
				"player": {{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
				"npc":    {{X: 23, Y: 10}},
			},
			expectCanvasDrawCount: 2,
			expectGameStatus:      GAME_STATUS_RUNING,
		},
		{
			//player正常的走一步， npc正常的走一次，
			//npc和player发生碰撞
			//player没有出界
			//player长度没有达到赢的长度
			inputPlayer: &StubPlayer{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
					{{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
				},
				moveInterval: 100 * time.Millisecond,
				lastMoveTime: time.Now(),
				dead:         false,
				eatCount:     0,
			},
			inputNpc: &StubNpc{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 15, Y: 10}},
					{{X: 13, Y: 10}},
				},
				moveInterval: 100 * time.Millisecond,
				dead:         false,
				lastMoveTime: time.Now(),
				beEatedCount: 0,
			},
			inputCanvas: &StubCanvas{
				maxWidth:  30,
				maxHeight: 30,
				boxMap:    make(map[string][]common.Pixel),
				drawCount: 0,
			},
			inputWinLen:           10,
			inputLoopTimes:        2,
			inputInterval:         80 * time.Millisecond,
			expectPlayerIdx:       1,
			expectPlayerShape:     []common.Pixel{{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
			expectPlayerIsDead:    false,
			expectPlayerEatCount:  1,
			expectNpcIdx:          1,
			expectNpcShape:        []common.Pixel{{X: 13, Y: 10}},
			expectNpcIsDead:       false,
			expectNpcBeEatedCount: 1,
			expectCanvasBox: map[string][]common.Pixel{
				"player": {{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
				"npc":    {{X: 13, Y: 10}},
			},
			expectCanvasDrawCount: 2,
			expectGameStatus:      GAME_STATUS_RUNING,
		},
		{
			//player正常的走一步， npc正常的走一次，
			//npc和player没有发生碰撞
			//player走出界
			//player长度没有达到赢的长度
			inputPlayer: &StubPlayer{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
					{{X: 33, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
				},
				moveInterval: 100 * time.Millisecond,
				lastMoveTime: time.Now(),
				dead:         false,
				eatCount:     0,
			},
			inputNpc: &StubNpc{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 15, Y: 10}},
					{{X: 23, Y: 10}},
				},
				moveInterval: 100 * time.Millisecond,
				dead:         false,
				lastMoveTime: time.Now(),
				beEatedCount: 0,
			},
			inputCanvas: &StubCanvas{
				maxWidth:  30,
				maxHeight: 30,
				boxMap:    make(map[string][]common.Pixel),
				drawCount: 0,
			},
			inputWinLen:           10,
			inputLoopTimes:        2,
			inputInterval:         80 * time.Millisecond,
			expectPlayerIdx:       1,
			expectPlayerShape:     []common.Pixel{{X: 33, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
			expectPlayerIsDead:    true,
			expectPlayerEatCount:  0,
			expectNpcIdx:          1,
			expectNpcShape:        []common.Pixel{{X: 23, Y: 10}},
			expectNpcIsDead:       false,
			expectNpcBeEatedCount: 0,
			expectCanvasBox: map[string][]common.Pixel{
				"player": {{X: 33, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
				"npc":    {{X: 23, Y: 10}},
			},
			expectCanvasDrawCount: 2,
			expectGameStatus:      GAME_STATUS_FAILED,
		},
		{
			//player正常的走一步， npc正常的走一次，
			//npc和player没有发生碰撞
			//player没有出界
			//player长度达到赢的长度
			inputPlayer: &StubPlayer{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
					{{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
				},
				moveInterval: 100 * time.Millisecond,
				lastMoveTime: time.Now(),
				dead:         false,
				eatCount:     0,
			},
			inputNpc: &StubNpc{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 15, Y: 10}},
					{{X: 23, Y: 10}},
				},
				moveInterval: 100 * time.Millisecond,
				dead:         false,
				lastMoveTime: time.Now(),
				beEatedCount: 0,
			},
			inputCanvas: &StubCanvas{
				maxWidth:  30,
				maxHeight: 30,
				boxMap:    make(map[string][]common.Pixel),
				drawCount: 0,
			},
			inputWinLen:           4,
			inputLoopTimes:        2,
			inputInterval:         80 * time.Millisecond,
			expectPlayerIdx:       1,
			expectPlayerShape:     []common.Pixel{{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			expectPlayerIsDead:    false,
			expectPlayerEatCount:  0,
			expectNpcIdx:          1,
			expectNpcShape:        []common.Pixel{{X: 23, Y: 10}},
			expectNpcIsDead:       false,
			expectNpcBeEatedCount: 0,
			expectCanvasBox: map[string][]common.Pixel{
				"player": {{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
				"npc":    {{X: 23, Y: 10}},
			},
			expectCanvasDrawCount: 2,
			expectGameStatus:      GAME_STATUS_WIN,
		},
	}

	for idx, test := range tests {
		w := Create(test.inputPlayer,
			test.inputNpc,
			test.inputCanvas,
			test.inputWinLen)

		w.gameStatus = GAME_STATUS_RUNING
		test.inputPlayer.lastMoveTime = time.Now()
		test.inputNpc.lastMoveTime = time.Now()

		fmt.Println("debug")
		for i := 0; i < test.inputLoopTimes; i++ {
			time.Sleep(test.inputInterval)
			w.refresh()
		}
		assert.Equal(t,
			test.expectPlayerIdx, test.inputPlayer.shapeListIdx, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectPlayerShape,
			test.inputPlayer.shapeList[test.inputPlayer.shapeListIdx],
			strconv.Itoa(idx))
		assert.Equal(t,
			test.expectPlayerIsDead, test.inputPlayer.dead, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectPlayerEatCount, test.inputPlayer.eatCount, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectNpcIdx, test.inputNpc.shapeListIdx, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectNpcShape,
			test.inputNpc.shapeList[test.inputNpc.shapeListIdx],
			strconv.Itoa(idx))
		assert.Equal(t,
			test.expectNpcIsDead, test.inputNpc.dead, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectNpcBeEatedCount, test.inputNpc.beEatedCount, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectCanvasBox, test.inputCanvas.boxMap, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectCanvasDrawCount, test.inputCanvas.drawCount, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectGameStatus, w.gameStatus, strconv.Itoa(idx))
	}
}

func TestRun(t *testing.T) {
	//游戏没有结束
	//refresh的间隔确实刷新了
	//键位的输入，确实可以
	//输入stop确实会结束

	tests := []struct {
		inputPlayer              *StubPlayer
		inputNpc                 *StubNpc
		inputCanvas              *StubCanvas
		inputCanvasKeyDuration   time.Duration
		inputCanvasTermDuration  time.Duration
		expectPlayerIdx          int
		expectPlayerShape        []common.Pixel
		expectPlayDirection      common.Direction
		expectNpcIdx             int
		expectNpcShape           []common.Pixel
		expectCanvasBox          map[string][]common.Pixel
		expectCanvasInitCount    int
		expectCanvasDestoryCount int
	}{
		{
			//测试refresh是不是正常
			//中间设置下方向
			//然后去stop
			inputPlayer: &StubPlayer{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
					{{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
					{{X: 14, Y: 10}, {X: 13, Y: 10}, {X: 12, Y: 10}},
				},
				moveInterval: 50 * time.Millisecond,
				lastMoveTime: time.Now(),
				dead:         false,
				dir:          common.DIRECTION_RIGHT,
			},
			inputNpc: &StubNpc{
				shapeListIdx: 0,
				shapeList: [][]common.Pixel{
					{{X: 15, Y: 10}},
					{{X: 23, Y: 10}},
					{{X: 24, Y: 10}},
				},
				moveInterval: 50 * time.Millisecond,
				dead:         false,
				lastMoveTime: time.Now(),
			},
			inputCanvas: &StubCanvas{
				maxWidth:     30,
				maxHeight:    30,
				boxMap:       make(map[string][]common.Pixel),
				keyChan:      make(chan common.Direction, 1),
				termChan:     make(chan struct{}, 1),
				drawCount:    0,
				initCount:    0,
				destoryCount: 0,
			},
			inputCanvasKeyDuration:  60 * time.Millisecond,
			inputCanvasTermDuration: 130 * time.Millisecond,
			expectPlayerIdx:         2,
			expectPlayerShape:       []common.Pixel{{X: 14, Y: 10}, {X: 13, Y: 10}, {X: 12, Y: 10}},
			expectPlayDirection:     common.DIRECTION_UP,
			expectNpcIdx:            2,
			expectNpcShape:          []common.Pixel{{X: 24, Y: 10}},
			expectCanvasBox: map[string][]common.Pixel{
				"player": {{X: 14, Y: 10}, {X: 13, Y: 10}, {X: 12, Y: 10}},
				"npc":    {{X: 24, Y: 10}},
			},
			expectCanvasDestoryCount: 1,
		},
	}

	for idx, test := range tests {
		w := Create(test.inputPlayer,
			test.inputNpc,
			test.inputCanvas,
			10)

		go func() {
			<-time.After(test.inputCanvasKeyDuration)
			test.inputCanvas.keyChan <- common.DIRECTION_UP
		}()
		go func() {
			<-time.After(test.inputCanvasTermDuration)
			test.inputCanvas.termChan <- struct{}{}
		}()
		w.Run()
		assert.Equal(t,
			test.expectPlayerIdx, test.inputPlayer.shapeListIdx, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectPlayerShape,
			test.inputPlayer.shapeList[test.inputPlayer.shapeListIdx],
			strconv.Itoa(idx))
		assert.Equal(t,
			test.expectPlayDirection,
			test.inputPlayer.dir,
			strconv.Itoa(idx))
		assert.Equal(t,
			test.expectNpcIdx, test.inputNpc.shapeListIdx, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectNpcShape,
			test.inputNpc.shapeList[test.inputNpc.shapeListIdx],
			strconv.Itoa(idx))
		assert.Equal(t,
			test.expectCanvasBox, test.inputCanvas.boxMap, strconv.Itoa(idx))
		assert.Equal(t,
			test.expectCanvasDestoryCount, test.inputCanvas.destoryCount, strconv.Itoa(idx))
	}
}
