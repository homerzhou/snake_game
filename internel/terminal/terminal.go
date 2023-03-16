package terminal

import (
	"log"
	"math/rand"
	"snake_game_l1/internel/common"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Terminal struct {
	boxs map[string]struct {
		shape []common.Pixel
		color tcell.Style
	}
	screen     tcell.Screen
	centerStr  string
	keyChan    chan common.Direction
	termChan   chan struct{}
	maxWidth   int
	maxHeigh   int
	colorWhite tcell.Style
	colorGreen tcell.Style
}

func Create() *Terminal {
	screen, e := tcell.NewScreen()
	if e != nil {
		return nil
	}

	t := Terminal{
		boxs: make(map[string]struct {
			shape []common.Pixel
			color tcell.Style
		}),
		screen:    screen,
		centerStr: "",
		keyChan:   make(chan common.Direction, 1),
		termChan:  make(chan struct{}, 1),
	}

	if t.Init() != nil {
		return nil
	}

	t.colorWhite = tcell.StyleDefault.Foreground(tcell.ColorCadetBlue.TrueColor()).Background(tcell.ColorWhite)
	t.colorGreen = tcell.StyleDefault.Foreground(tcell.ColorCadetBlue.TrueColor()).Background(tcell.ColorGreen)

	return &t
}

func (t *Terminal) Init() error {
	if e := t.screen.Init(); e != nil {
		return e
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	t.screen.SetStyle(defStyle)

	t.screen.Clear()

	go func() {
		for {
			switch ev := t.screen.PollEvent().(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyUp:
					t.pushKeyChanNonBlock(common.DIRECTION_UP)
				case tcell.KeyDown:
					t.pushKeyChanNonBlock(common.DIRECTION_DOWN)
				case tcell.KeyLeft:
					t.pushKeyChanNonBlock(common.DIRECTION_LEFT)
				case tcell.KeyRight:
					t.pushKeyChanNonBlock(common.DIRECTION_RIGHT)
				case tcell.KeyCtrlC:
					t.pushTermChanNonBlock()
				default:
				}
			}
		}
	}()

	t.maxWidth, t.maxHeigh = t.screen.Size()

	return nil
}

func (t *Terminal) pushKeyChanNonBlock(d common.Direction) {
	select {
	case t.keyChan <- d:
	default:
	}
}

func (t *Terminal) pushTermChanNonBlock() {
	select {
	case t.termChan <- struct{}{}:
	default:
	}
}

func (t *Terminal) IsOut(box []common.Pixel) bool {
	for _, b := range box {
		if b.X >= t.maxWidth || b.Y >= t.maxHeigh || b.X <= 0 || b.Y <= 0 {
			return true
		}
	}
	return false
}

func (t *Terminal) Random(len int) ([]common.Pixel, error) {
	//todo conflict with box
	w := rand.Intn(t.maxWidth-1) + 1
	h := rand.Intn(t.maxHeigh-1) + 1
	return []common.Pixel{{X: w, Y: h}}, nil
}

func (t *Terminal) SetBox(name string, box []common.Pixel, color common.Color) {
	if color == common.COLOR_WHITE {
		t.boxs[name] = struct {
			shape []common.Pixel
			color tcell.Style
		}{
			shape: box,
			color: t.colorWhite,
		}
	} else {
		t.boxs[name] = struct {
			shape []common.Pixel
			color tcell.Style
		}{
			shape: box,
			color: t.colorGreen,
		}
	}
}

func (t *Terminal) SetCenterStr(str string) {
	t.centerStr = str
}

func (t *Terminal) drawCenterStr() {

	w, h := t.maxWidth, t.maxHeigh
	x := w/2 - 7
	y := h / 2

	style := tcell.StyleDefault.Foreground(tcell.ColorCadetBlue.TrueColor()).Background(tcell.ColorWhite)
	for _, c := range t.centerStr {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		t.screen.SetContent(x, y, c, comb, style)
		x += w
	}
}

func (t *Terminal) drawBoxes() {
	log.Println("start to drawbox: ", t.boxs)
	for _, box := range t.boxs {
		for _, p := range box.shape {
			t.screen.SetContent(p.X, t.maxHeigh-p.Y, ' ', nil, box.color)
		}
	}
}

func (t *Terminal) Draw() {
	t.screen.Clear()
	t.drawBoxes()
	t.drawCenterStr()
	t.screen.Show()
}

func (t *Terminal) GetKeyboardEvent() chan common.Direction {
	return t.keyChan
}

func (t *Terminal) GetTermEvent() chan struct{} {
	return t.termChan
}

func (t *Terminal) Destory() {
	t.screen.Fini()
}

func (t *Terminal) GetSize() (int, int) {
	return t.maxWidth, t.maxHeigh
}
