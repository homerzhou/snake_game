package food

import (
	"fmt"
	"log"
	"snake_game_l1/internel/common"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RandomStubElement struct {
	pixels  []common.Pixel
	isError bool
}

type RandomStub struct {
	randoms []RandomStubElement
	count   int
}

func (r *RandomStub) Random(len int) (rp []common.Pixel, err error) {
	rp = r.randoms[r.count].pixels
	if r.randoms[r.count].isError {
		err = fmt.Errorf("random fail")
	}
	r.count++
	return
}

func TestFood1(t *testing.T) {
	tests := []struct {
		inputInitPixels       common.Pixel
		inputInitMoveInterval time.Duration
		inputNewMoveInterval  time.Duration
		inputRandom           []RandomStubElement
		inputLoopTimes        int
		inputLoopInterval     time.Duration
		expectDead            bool
		expectRandomCount     int
		expectShape           []common.Pixel
	}{
		{
			inputInitPixels:       common.Pixel{X: 11, Y: 11},
			inputInitMoveInterval: 100 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputRandom: []RandomStubElement{
				{pixels: []common.Pixel{{X: 1, Y: 1}}, isError: false},
				{pixels: []common.Pixel{{X: 2, Y: 2}}, isError: false},
				{pixels: nil, isError: true},
			},
			inputLoopTimes:    1,
			inputLoopInterval: 50 * time.Millisecond,
			expectDead:        false,
			expectRandomCount: 0,
			expectShape:       []common.Pixel{{X: 11, Y: 11}},
		},
		{
			inputInitPixels:       common.Pixel{X: 11, Y: 11},
			inputInitMoveInterval: 100 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputRandom: []RandomStubElement{
				{pixels: []common.Pixel{{X: 1, Y: 1}}, isError: false},
				{pixels: []common.Pixel{{X: 2, Y: 2}}, isError: false},
				{pixels: nil, isError: true},
			},
			inputLoopTimes:    2,
			inputLoopInterval: 80 * time.Millisecond,
			expectDead:        false,
			expectRandomCount: 1,
			expectShape:       []common.Pixel{{X: 1, Y: 1}},
		},
		{
			inputInitPixels:       common.Pixel{X: 11, Y: 11},
			inputInitMoveInterval: 100 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputRandom: []RandomStubElement{
				{pixels: []common.Pixel{{X: 1, Y: 1}}, isError: false},
				{pixels: []common.Pixel{{X: 2, Y: 2}}, isError: false},
				{pixels: []common.Pixel{{X: 3, Y: 3}}, isError: false},
				{pixels: nil, isError: true},
			},
			inputLoopTimes:    4,
			inputLoopInterval: 80 * time.Millisecond,
			expectDead:        false,
			expectRandomCount: 3,
			expectShape:       []common.Pixel{{X: 3, Y: 3}},
		},
		{
			inputInitPixels:       common.Pixel{X: 11, Y: 11},
			inputInitMoveInterval: 100 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputRandom: []RandomStubElement{
				{pixels: []common.Pixel{{X: 1, Y: 1}}, isError: false},
				{pixels: []common.Pixel{{X: 2, Y: 2}}, isError: false},
				{pixels: nil, isError: true},
			},
			inputLoopTimes:    3,
			inputLoopInterval: 110 * time.Millisecond,
			expectDead:        true,
			expectRandomCount: 3,
			expectShape:       nil,
		},
		{
			inputInitPixels:       common.Pixel{X: 11, Y: 11},
			inputInitMoveInterval: 100 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputRandom: []RandomStubElement{
				{pixels: []common.Pixel{{X: 1, Y: 1}}, isError: false},
				{pixels: []common.Pixel{{X: 2, Y: 2}}, isError: false},
				{pixels: nil, isError: true},
			},
			inputLoopTimes:    1,
			inputLoopInterval: 330 * time.Millisecond,
			expectDead:        true,
			expectRandomCount: 3,
			expectShape:       nil,
		},
		{
			inputInitPixels:       common.Pixel{X: 11, Y: 11},
			inputInitMoveInterval: 10 * time.Millisecond,
			inputNewMoveInterval:  100 * time.Millisecond,
			inputRandom: []RandomStubElement{
				{pixels: []common.Pixel{{X: 1, Y: 1}}, isError: false},
				{pixels: []common.Pixel{{X: 2, Y: 2}}, isError: false},
				{pixels: nil, isError: true},
			},
			inputLoopTimes:    3,
			inputLoopInterval: 110 * time.Millisecond,
			expectDead:        true,
			expectRandomCount: 3,
			expectShape:       nil,
		},
	}

	for idx, test := range tests {
		rStub := RandomStub{test.inputRandom, 0}
		f := Create(test.inputInitPixels, test.inputInitMoveInterval, &rStub)
		assert.NotEqual(t, f, nil, strconv.Itoa(idx))

		for i := 0; i < test.inputLoopTimes; i++ {
			log.Println("debug sleep, ", time.Now().UnixMilli(), "    ", f.lastMoveTime.UnixMilli())
			time.Sleep(test.inputLoopInterval)
			log.Println("debug sleep, ", time.Now().UnixMilli(), "    ", f.lastMoveTime.UnixMilli())
			f.SetMoveInterval(test.inputNewMoveInterval)
			f.Run()
		}
		assert.Equal(t, test.expectDead, f.IsDead(), strconv.Itoa(idx))
		assert.Equal(t, test.expectRandomCount, rStub.count, strconv.Itoa(idx))
		assert.Equal(t, test.expectShape, f.GetShape(), strconv.Itoa(idx))
		assert.Equal(t, test.inputNewMoveInterval, f.GetMoveInterval(), strconv.Itoa(idx))
	}
}
