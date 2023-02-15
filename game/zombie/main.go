package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const GameFrameWidth = 30
const GameFrameHeight = 15
const GameFrameSymbol = '|'

type Point struct {
	row, col int
}

type GameObject struct {
	points         []*Point
	velRow, velCol int
}

var screen tcell.Screen
var isGameOver bool
var isGamePaused bool
var debugLog string
var score int

func main() {
	rand.Seed(time.Now().UnixNano())
	InitScreen()
	InitGameState()
	screen.HideCursor()
	inputChan := InitUserInput()

	for !isGameOver {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()
		time.Sleep(100 * time.Millisecond)
	}

	screenWidth, screenHeight := screen.Size()
	PrintStringCentered(screenHeight/2-1, screenWidth/2, "Game Over!!!")
	PrintStringCentered(screenHeight/2, screenWidth/2, fmt.Sprintf("Your score is %d...", score))
	screen.Show()

	time.Sleep(3 * time.Second)
	screen.Fini()
}

func UpdateState() {
	if isGamePaused {
		return
	}
}

func DrawState() {
	if isGamePaused {
		return
	}

	screen.Clear()
	PrintString(0, 0, debugLog)
	PrintGameFrame()

	screen.Show()
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func InitGameState() {
}

func HandleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" {
	} else if key == "Rune[a]" {
	} else if key == "Rune[s]" {
	} else if key == "Rune[d]" {
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	}
}

func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}

	return key
}

func PrintGameFrame() {
	gameFrameTopLeftRow, gameFrameTopLeftCol := GetGameFrameTopLeft()
	row, col := gameFrameTopLeftRow-1, gameFrameTopLeftCol-1
	width, height := GameFrameWidth+2, GameFrameHeight+2

	PrintUnfilledRect(row, col, width, height, GameFrameSymbol)
}

func PrintStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		PrintFilledRect(row, col, 1, 1, c)
		col += 1
	}
}

func PrintFilledRectInGameFrame(row, col, width, height int, ch rune) {
	r, c := GetGameFrameTopLeft()
	PrintFilledRect(row+r, col+c, width, height, ch)
}

func PrintFilledRect(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func PrintUnfilledRect(row, col, width, height int, ch rune) {
	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row, ch, nil, tcell.StyleDefault)
	}

	for r := 1; r < height-1; r++ {
		screen.SetContent(col, row+r, ch, nil, tcell.StyleDefault)
		screen.SetContent(col+width-1, row+r, ch, nil, tcell.StyleDefault)
	}

	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row+height-1, ch, nil, tcell.StyleDefault)
	}
}

func GetGameFrameTopLeft() (int, int) {
	screenWidth, screenHeight := screen.Size()
	return screenHeight/2 - GameFrameHeight/2, screenWidth/2 - GameFrameWidth/2
}

