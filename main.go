package main

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var board [8][8]string

func main() {
	fmt.Println("Starting")

	boardSetup()
	fmt.Println(board)
	printBoard()

	pixelgl.Run(run)

}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Sharpy's Othello Game",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Skyblue)

	board := imdraw.New(nil)
	board.Color = colornames.Lightgreen
	board.Push(pixel.V(50, 50))
	board.Push(pixel.V(650, 650))
	board.Rectangle(0)

	grid := imdraw.New(nil)
	grid.Color = colornames.Black
	for x := 0.0; x < 600; x = x + 75 {
		for y := 0.0; y < 600; y = y + 75 {
			grid.Push(pixel.V(50+x, 50+y))
			grid.Push(pixel.V(125+x, 125+y))
			grid.Rectangle(1)
		}
	}

	circle := drawCircle(4, 4, colornames.Black)
	circle2 := drawCircle(4, 5, colornames.White)

	// ... later in the code
	for !win.Closed() {
		// ...
		win.Clear(colornames.Aliceblue)
		board.Draw(win)
		grid.Draw(win)
		circle.Draw(win)
		circle2.Draw(win)

		win.Update()
		// ...
	}
}

func drawCircle(x, y float64, c color.Color) *imdraw.IMDraw {
	circle := imdraw.New(nil)
	circle.Color = c
	circle.Push(pixel.V((x*75)+12, (y*75)+12))
	circle.Circle(30, 0)
	return circle
}

func boardSetup() {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			board[x][y] = " "
		}
	}
	board[3][3] = "B"
	board[3][4] = "W"
	board[4][3] = "W"
	board[4][4] = "B"
}

func printBoard() {
	for x := 0; x < 8; x++ {
		fmt.Printf("\n")
		for y := 0; y < 8; y++ {
			switch board[x][y] {
			case " ":
				fmt.Printf(".")
			case "B":
				fmt.Printf("B")
			case "W":
				fmt.Printf("W")
			}
		}
	}
	fmt.Printf("\n")
}
