package snake

import (
	"log"
	"snake_game_l1/internel/common"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestInputSnakeAction struct {
	dir     common.Direction
	isEat   bool
	isCrawl bool
}

func TestEatCrawl(t *testing.T) {
	tests := []struct {
		inputInitPixels []common.Pixel
		inputActions    []TestInputSnakeAction
		expectPixels    []common.Pixel
	}{
		{
			inputInitPixels: []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			inputActions: []TestInputSnakeAction{
				{common.DIRECTION_UP, false, true},
			},
			expectPixels: []common.Pixel{{X: 12, Y: 11}, {X: 12, Y: 10}, {X: 11, Y: 10}},
		},
		{
			inputInitPixels: []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			inputActions: []TestInputSnakeAction{
				{common.DIRECTION_RIGHT, false, true},
				{common.DIRECTION_UP, false, true},
				{common.DIRECTION_LEFT, false, true},
				{common.DIRECTION_DOWN, false, true},
			},
			expectPixels: []common.Pixel{{X: 12, Y: 10}, {X: 12, Y: 11}, {X: 13, Y: 11}},
		},
		{
			//测试一些无效的方向变化
			inputInitPixels: []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			inputActions: []TestInputSnakeAction{
				{common.DIRECTION_UP, false, true},
				{common.DIRECTION_UP, false, true},
				{common.DIRECTION_DOWN, false, true},
				{common.DIRECTION_LEFT, false, true},
				{common.DIRECTION_LEFT, false, true},
				{common.DIRECTION_RIGHT, false, true},
				{common.DIRECTION_DOWN, false, true},
			},
			expectPixels: []common.Pixel{{X: 9, Y: 12}, {X: 9, Y: 13}, {X: 10, Y: 13}},
		},
		{
			inputInitPixels: []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			inputActions: []TestInputSnakeAction{
				{common.DIRECTION_RIGHT, false, true},
				{common.DIRECTION_UP, false, true},
				{common.DIRECTION_UP, true, false},
				{common.DIRECTION_UP, false, true},
				{common.DIRECTION_UP, false, true},
				{common.DIRECTION_UP, true, false},

				{common.DIRECTION_LEFT, false, true},
				{common.DIRECTION_LEFT, false, true},
				{common.DIRECTION_LEFT, false, true},

				{common.DIRECTION_DOWN, false, true},
				{common.DIRECTION_DOWN, false, true},

				{common.DIRECTION_DOWN, true, false},

				{common.DIRECTION_RIGHT, false, true},
				{common.DIRECTION_RIGHT, false, true},
				{common.DIRECTION_RIGHT, false, true},

				{common.DIRECTION_RIGHT, true, false},
			},
			expectPixels: []common.Pixel{{X: 13, Y: 11}, {X: 12, Y: 11}, {X: 11, Y: 11}, {X: 10, Y: 11}, {X: 10, Y: 12}, {X: 10, Y: 13}, {X: 10, Y: 14}},
		},
	}

	for idx, test := range tests {
		s := Create(test.inputInitPixels, 1*time.Second)
		for _, action := range test.inputActions {
			s.SetDirection(action.dir)
			if action.isCrawl {
				s.crawl()
			}
			if action.isEat {
				s.Eat()
			}
			log.Println(s.GetShape())
		}
		assert.Equal(t, test.expectPixels, s.GetShape(), strconv.Itoa(idx))
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		inputInitPixels       []common.Pixel
		inputInitMoveInterval time.Duration
		inputNewMoveInterval  time.Duration
		inputLoopDirection    []common.Direction
		inputLoopEat          []bool
		inputLoopTimes        int
		inputLoopInterval     time.Duration
		inputIsKill           bool
		expectIsDeal          bool
		expectMoveInterval    time.Duration
		expectShape           []common.Pixel
	}{
		{
			inputInitPixels:       []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			inputInitMoveInterval: 100 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputLoopTimes:        1,
			inputLoopInterval:     110 * time.Millisecond,
			inputIsKill:           false,
			expectIsDeal:          false,
			expectMoveInterval:    100 * time.Millisecond,
			expectShape:           []common.Pixel{{X: 13, Y: 10}, {X: 12, Y: 10}, {X: 11, Y: 10}},
		},
		{
			inputInitPixels:       []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			inputInitMoveInterval: 100 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputLoopTimes:        1,
			inputLoopInterval:     220 * time.Millisecond,
			inputIsKill:           false,
			expectIsDeal:          false,
			expectMoveInterval:    100 * time.Millisecond,
			expectShape:           []common.Pixel{{X: 14, Y: 10}, {X: 13, Y: 10}, {X: 12, Y: 10}},
		},
		{
			inputInitPixels:       []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			inputInitMoveInterval: 10 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputLoopTimes:        4,
			inputLoopInterval:     60 * time.Millisecond,
			inputIsKill:           true,
			expectIsDeal:          true,
			expectMoveInterval:    100 * time.Millisecond,
			expectShape:           []common.Pixel{{X: 14, Y: 10}, {X: 13, Y: 10}, {X: 12, Y: 10}},
		},
		{
			//移动的过程中， 吃一些东西
			inputInitPixels:       []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}},
			inputInitMoveInterval: 10 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputLoopEat:          []bool{false, true, false},
			inputLoopTimes:        3,
			inputLoopInterval:     110 * time.Millisecond,
			inputIsKill:           false,
			expectIsDeal:          false,
			expectMoveInterval:    100 * time.Millisecond,
			expectShape:           []common.Pixel{{X: 15, Y: 10}, {X: 14, Y: 10}, {X: 13, Y: 10}, {X: 12, Y: 10}},
		},
		{
			//自己吃掉了自己
			inputInitPixels:       []common.Pixel{{X: 12, Y: 10}, {X: 11, Y: 10}, {X: 10, Y: 10}, {X: 9, Y: 10}, {X: 8, Y: 10}},
			inputInitMoveInterval: 10 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputLoopDirection:    []common.Direction{common.DIRECTION_UP, common.DIRECTION_LEFT, common.DIRECTION_DOWN, common.DIRECTION_DOWN},
			inputLoopTimes:        4,
			inputLoopInterval:     110 * time.Millisecond,
			inputIsKill:           false,
			expectIsDeal:          true,
			expectMoveInterval:    100 * time.Millisecond,
			expectShape:           []common.Pixel{{X: 11, Y: 10}, {X: 11, Y: 11}, {X: 12, Y: 11}, {X: 12, Y: 10}, {X: 11, Y: 10}},
		},
	}

	for idx, test := range tests {
		s := Create(test.inputInitPixels, test.inputInitMoveInterval)
		for i := 0; i < test.inputLoopTimes; i++ {
			time.Sleep(test.inputLoopInterval)
			s.SetMoveInterval(test.inputNewMoveInterval)
			if test.inputLoopDirection != nil {
				s.SetDirection(test.inputLoopDirection[i])
			}
			if test.inputLoopEat != nil && test.inputLoopEat[i] == true {
				s.Eat()
			}
			s.Run()
		}
		if test.inputIsKill {
			s.Kill()
		}
		assert.Equal(t, test.expectIsDeal, s.IsDead(), strconv.Itoa(idx))
		assert.Equal(t, test.expectMoveInterval, s.GetMoveInterval(), strconv.Itoa(idx))
		assert.Equal(t, test.expectShape, s.GetShape(), strconv.Itoa(idx))
	}
}
