package main

import (
	"io/ioutil"
	"log"
	"os"
	"snake_game_l1/internel/common"
	"snake_game_l1/internel/food"
	"snake_game_l1/internel/snake"
	"snake_game_l1/internel/terminal"
	"snake_game_l1/internel/world"
	"time"
)

func main() {
	f, err := os.OpenFile("run.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	// log.SetOutput(f)
	log.SetOutput(ioutil.Discard)

	terminal := terminal.Create()

	tWidth, tHeight := terminal.GetSize()
	foodInitPixel := common.Pixel{X: tWidth / 2, Y: tHeight / 2}
	snakeInitPixel := []common.Pixel{
		{X: tWidth / 3, Y: tHeight / 3},
		{X: tWidth/3 - 1, Y: tHeight / 3},
		{X: tWidth/3 - 2, Y: tHeight / 3},
	}

	log.Println("Init pixel: ", foodInitPixel, snakeInitPixel)
	food := food.Create(foodInitPixel, 100*time.Second, terminal)
	snake := snake.Create(snakeInitPixel, 200*time.Millisecond)

	world := world.Create(snake, food, terminal, 10)
	world.Run()
}
