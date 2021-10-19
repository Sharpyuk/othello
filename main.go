package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var board [8][8]string
var player string = "B"
var allValidMoves = []affectedSquares{}

type affectedSquares struct {
	X, Y int
}

func main() {
	fmt.Println("Starting")

	boardSetup()
	fmt.Println(board)
	printBoard(board)

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

	drawBoard := drawBoard()
	grid := drawGrid()
	currentPlayerText, currentPlayer := drawCurrentPlayer()
	currentPlayerText.Draw(win, pixel.IM.Scaled(currentPlayerText.Orig, 4))
	currentPlayer.Draw(win)

	for !win.Closed() {
		win.Clear(colornames.Navy)
		drawBoard.Draw(win)
		grid.Draw(win)
		drawAvailableMoves(allValidMoves, win)

		currentPlayerText, currentPlayer := drawCurrentPlayer()
		currentPlayerText.Draw(win, pixel.IM.Scaled(currentPlayerText.Orig, 4))
		currentPlayer.Draw(win)

		refreshBoard(win)
		blackPlayerScore, whitePlayerScore := updateScores()
		blackPlayerScore.Draw(win, pixel.IM.Scaled(blackPlayerScore.Orig, 4))
		whitePlayerScore.Draw(win, pixel.IM.Scaled(whitePlayerScore.Orig, 4))

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			// Which square was clicked on?
			x, y := getClickedBox(win.MousePosition())

			// Check if player made a valid move
			affected := findAffectedSquares(player, x-1, y-1)
			if len(affected) > 0 {

				// Update board to remove possible move grey squares
				drawBoard.Draw(win)
				grid.Draw(win)
				refreshBoard(win)
				win.Update()

				drawFlipAnimation(float64(x-1), float64(y-1), player, win)
				board[x-1][y-1] = player
				flipAffectedSquares(affected, win)

				player = updatePlayer(player)
				//currentPlayerText, currentPlayer = drawCurrentPlayer()

				for player == "W" {
					allValidMoves = validMoves(player)
					if len(allValidMoves) == 0 {
						player = updatePlayer(player)
						break
					}
					time.Sleep(1 * time.Second)
					playComputerMove(allValidMoves, win)
					player = updatePlayer(player)
				}

				allValidMoves = validMoves(player)
				if len(allValidMoves) == 0 {
					player = updatePlayer(player)
					if len(validMoves(player)) == 0 {
						gameOver(win)
					}
				}
			}

		}

		win.Update()
		// ...
	}
}

func gameOver(win *pixelgl.Window) {
	fmt.Println("Game Over!")
	var w, b int
	for x := 0; x < 7; x++ {
		for y := 0; y < 7; y++ {
			switch board[x][y] {
			case "W":
				w++
			case "B":
				b++
			}
		}
	}
	fmt.Println("White: %d, Black: %d", w, b)
	if w > b {
		fmt.Println("White Wins..")
	} else {
		fmt.Println("Black Wins!  Well done!!")
	}
}

func playComputerMove(possibleMoves []affectedSquares, win *pixelgl.Window) {
	var foundMove affectedSquares
	var maxAffectedSquares int
	for _, p := range possibleMoves {
		fmt.Println("-Playing Computer Move")
		fmt.Printf(" --Checking X:%d  Y:%d\n", p.X, p.Y)
		a := findAffectedSquares("W", p.X, p.Y)

		if len(a) > maxAffectedSquares {
			fmt.Printf("  ---Highest: %d", len(a))
			foundMove.X = p.X
			foundMove.Y = p.Y
			maxAffectedSquares = len(a)
		}
	}

	affected := findAffectedSquares(player, foundMove.X, foundMove.Y)

	drawFlipAnimation(float64(foundMove.X), float64(foundMove.Y), player, win)
	board[foundMove.X][foundMove.Y] = "W"

	flipAffectedSquares(affected, win)

	fmt.Printf("Computer Move: %d, %d\n", foundMove.X, foundMove.Y)
}

func validMoves(p string) []affectedSquares {
	possibleMoves := []affectedSquares{}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if board[x][y] == " " {
				chkSquare := findAffectedSquares(p, x, y)
				if len(chkSquare) > 0 {
					//possibleMoves = append(possibleMoves, chkSquare...)

					possibleMoves = append(possibleMoves, affectedSquares{X: x, Y: y})

				}
			}
		}
	}
	return possibleMoves
}

func flipAffectedSquares(affected []affectedSquares, win *pixelgl.Window) {
	for _, square := range affected {
		switch board[square.X][square.Y] {
		case "B":
			board[square.X][square.Y] = "W"

		case "W":
			board[square.X][square.Y] = "B"
		}

		drawFlipAnimation(float64(square.X), float64(square.Y), board[square.X][square.Y], win)

	}
}
func findAffectedSquares(p string, x, y int) []affectedSquares {
	if board[x][y] != " " {
		return []affectedSquares{}
	}
	AllSquaresAffected := []affectedSquares{}

	// Check left
	affected := []affectedSquares{}
chkLeft:
	for rowIndex := x + 1; rowIndex < 8; rowIndex++ {
		fmt.Printf("Checking Left X: %d, Y: %d = %s (p: %s)\n", x, y, board[rowIndex][y], p)
		switch board[rowIndex][y] {
		case " ":
			break chkLeft
		case p:
			AllSquaresAffected = append(AllSquaresAffected, affected...)
			break chkLeft
		default:
			affected = append(affected, affectedSquares{X: rowIndex, Y: y})
		}
	}

	// Check Right
	affected = []affectedSquares{}
chkRight:
	for rowIndex := x - 1; rowIndex > -1; rowIndex-- {
		fmt.Printf("Checking Right X: %d, Y: %d = %s (p: %s)\n", x, y, board[rowIndex][y], p)
		switch board[rowIndex][y] {
		case " ":
			break chkRight
		case p:
			AllSquaresAffected = append(AllSquaresAffected, affected...)
			break chkRight
		default:
			affected = append(affected, affectedSquares{X: rowIndex, Y: y})
		}
	}

	// Check Up
	affected = []affectedSquares{}
chkUp:
	for colIndex := y + 1; colIndex < 8; colIndex++ {
		fmt.Printf("Checking Up X: %d, Y: %d = %s (p: %s)\n", x, y, board[x][colIndex], p)
		switch board[x][colIndex] {
		case " ":
			break chkUp
		case p:
			AllSquaresAffected = append(AllSquaresAffected, affected...)
			break chkUp
		default:
			affected = append(affected, affectedSquares{X: x, Y: colIndex})
		}
	}

	// Check Dowm
	affected = []affectedSquares{}
chkDown:
	for colIndex := y - 1; colIndex > -1; colIndex-- {
		fmt.Printf("Checking Down X: %d, Y: %d = %s (p: %s)\n", x, y, board[x][colIndex], p)
		switch board[x][colIndex] {
		case " ":
			break chkDown
		case p:
			AllSquaresAffected = append(AllSquaresAffected, affected...)
			break chkDown
		default:
			affected = append(affected, affectedSquares{X: x, Y: colIndex})
		}
	}

	// Check DiagUpRight
	affected = []affectedSquares{}

	colIndex := y
	rowIndex := x
DiagUpRight:
	for colIndex < 7 && rowIndex < 7 {
		colIndex++
		rowIndex++
		fmt.Printf("Checking DiagUpRight X: %d, Y: %d = %s (p: %s)\n", rowIndex, colIndex, board[rowIndex][colIndex], p)
		switch board[rowIndex][colIndex] {
		case " ":
			break DiagUpRight
		case p:
			AllSquaresAffected = append(AllSquaresAffected, affected...)
			break DiagUpRight
		default:
			affected = append(affected, affectedSquares{X: rowIndex, Y: colIndex})
		}
	}

	// Check DiagDownLeft
	affected = []affectedSquares{}

	colIndex = y
	rowIndex = x
DiagDownLeft:
	for colIndex > 0 && rowIndex > 0 {
		colIndex--
		rowIndex--
		fmt.Printf("Checking DiagDownLeft X: %d, Y: %d = %s (p: %s)\n", rowIndex, colIndex, board[rowIndex][colIndex], p)
		switch board[rowIndex][colIndex] {
		case " ":
			break DiagDownLeft
		case p:
			AllSquaresAffected = append(AllSquaresAffected, affected...)
			break DiagDownLeft
		default:
			affected = append(affected, affectedSquares{X: rowIndex, Y: colIndex})
		}
	}

	// Check DiagUpLeft
	affected = []affectedSquares{}

	colIndex = y
	rowIndex = x
DiagUpLeft:
	for colIndex < 7 && rowIndex > 0 {
		colIndex++
		rowIndex--
		fmt.Printf("Checking DiagUpLeft X: %d, Y: %d = %s (p: %s)\n", rowIndex, colIndex, board[rowIndex][colIndex], p)
		switch board[rowIndex][colIndex] {
		case " ":
			break DiagUpLeft
		case p:
			AllSquaresAffected = append(AllSquaresAffected, affected...)
			break DiagUpLeft
		default:
			affected = append(affected, affectedSquares{X: rowIndex, Y: colIndex})
		}
	}

	// Check DiagDownRight
	affected = []affectedSquares{}

	colIndex = y
	rowIndex = x
DiagDownRight:
	for colIndex > 0 && rowIndex < 7 {
		colIndex--
		rowIndex++
		fmt.Printf("Checking DiagDownRight X: %d, Y: %d = %s (p: %s)\n", rowIndex, colIndex, board[rowIndex][colIndex], p)
		switch board[rowIndex][colIndex] {
		case " ":
			break DiagDownRight
		case p:
			AllSquaresAffected = append(AllSquaresAffected, affected...)
			break DiagDownRight
		default:
			affected = append(affected, affectedSquares{X: rowIndex, Y: colIndex})
		}
	}

	fmt.Printf("Affected: %+v\n", AllSquaresAffected)
	return AllSquaresAffected
}

func updatePlayer(p string) string {
	// Check if next player can move
	nextPlayer := ""
	switch p {
	case "B":
		nextPlayer = "W"
	case "W":
		nextPlayer = "B"
	}
	valid := validMoves(nextPlayer)
	if len(valid) > 0 {
		p = nextPlayer
	}
	return p
}

func getClickedBox(vect pixel.Vec) (int, int) {
	x := ((vect.X - 50) / 75)
	y := ((vect.Y - 50) / 75)
	fmt.Println(x, y)
	if x > 8 || x < 0 {
		return 0, 0
	}
	if y > 8 || y < 0 {
		return 0, 0
	}
	return int(x + 1), int(y + 1)
}

func refreshBoard(win *pixelgl.Window) {

	for x := 0.0; x < 8; x++ {
		for y := 0.0; y < 8; y++ {

			switch board[int(x)][int(y)] {
			case "W":
				drawCircle(x+1, y+1, colornames.White).Draw(win)
			case "B":
				drawCircle(x+1, y+1, colornames.Black).Draw(win)
			}

		}
	}
}

func updateScores() (*text.Text, *text.Text) {

	var black, white int
	for x := 0; x < 7; x++ {
		for y := 0; y < 7; y++ {
			switch board[x][y] {
			case "B":
				black++
			case "W":
				white++
			}
		}
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(700, 300), basicAtlas)

	fmt.Fprintf(basicTxt, "WHITE: %d", white)

	basicAtlas2 := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt2 := text.New(pixel.V(700, 200), basicAtlas2)

	fmt.Fprintf(basicTxt2, "BLACK: %d", black)
	return basicTxt, basicTxt2

}

func drawBoard() *imdraw.IMDraw {
	board := imdraw.New(nil)
	board.Color = colornames.Darkgreen
	board.Push(pixel.V(50, 50))
	board.Push(pixel.V(650, 650))
	board.Rectangle(0)
	return board
}

func drawAvailableMoves(a []affectedSquares, w *pixelgl.Window) {
	for _, av := range a {
		circle := imdraw.New(nil)
		circle.Color = colornames.Gray
		circle.Push(pixel.V(((float64(av.X)+1)*75)+12, ((float64(av.Y)+1)*75)+12))
		circle.Circle(30, 5)
		circle.Draw(w)
	}
}

func drawGrid() *imdraw.IMDraw {
	grid := imdraw.New(nil)
	grid.Color = colornames.Black
	for x := 0.0; x < 600; x = x + 75 {
		for y := 0.0; y < 600; y = y + 75 {
			grid.Push(pixel.V(50+x, 50+y))
			grid.Push(pixel.V(125+x, 125+y))
			grid.Rectangle(1)
		}
	}
	return grid
}

func drawFlipAnimation(x, y float64, c string, win *pixelgl.Window) {
	x += 1
	y += 1
	playerColor := colornames.White
	if player == "B" {
		playerColor = colornames.Black
	}
	circleArc := imdraw.New(nil)
	circleArc.Color = playerColor
	for i := 0.0; i < 9.0; i++ {
		time.Sleep(30 * time.Millisecond)
		circleArc.Push(pixel.V((x*75)+12, (y*75)+12))
		circleArc.CircleArc(30, 0, i*math.Pi, 0)
		circleArc.Draw(win)
		win.Update()
	}
}

func drawCircle(x, y float64, c color.Color) *imdraw.IMDraw {
	circle := imdraw.New(nil)
	circle.Color = c
	circle.Push(pixel.V((x*75)+12, (y*75)+12))
	circle.Circle(30, 0)
	return circle
}

func drawCurrentPlayer() (*text.Text, *imdraw.IMDraw) {
	playerColor := colornames.White
	if player == "B" {
		playerColor = colornames.Black
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(700, 700), basicAtlas)

	fmt.Fprintf(basicTxt, "Player:")

	currentPlayer := imdraw.New(nil)
	currentPlayer.Color = playerColor
	currentPlayer.Push(pixel.V(950, 720))
	currentPlayer.Circle(30, 0)

	return basicTxt, currentPlayer
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

func printBoard(b [8][8]string) {
	for x := 0; x < 8; x++ {
		fmt.Printf("\n")
		for y := 0; y < 8; y++ {
			switch b[x][y] {
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
