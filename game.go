package main

import (
	"image/color"
	"math/rand"
	"time"
	"sync"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/imdraw"
)

const (
	width     = 1024
	height    = 768
	cellSize  = 10
	rows      = height / cellSize
	cols      = width / cellSize
)

var grid [rows][cols]bool

// 初期化：ランダムにセルの状態を決定
func initGrid() {
	rand.Seed(time.Now().UnixNano())
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			grid[y][x] = rand.Float64() < 0.2 // 20%の確率で生きてる
		}
	}
}

// 隣接する生きたセルの数を数える
func liveNeighbors(y, x int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dy == 0 && dx == 0 {
				continue
			}
			ny := (y + dy + rows) % rows
			nx := (x + dx + cols) % cols
			if grid[ny][nx] {
				count++
			}
		}
	}
	return count
}

func randomGenerateCell()  {
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if !grid[y][x] {
				if rand.Float64() < 0.01 { // 1%の確率で新しいセルを生成
					grid[y][x] = rand.Float64() < 0.01 // 1%の確率で生きてる
				}
			}
		}
	}
}


func updateGrid() {
	var next [rows][cols]bool
	var wg sync.WaitGroup
	numWorkers := 8 // コア数に応じて調整

	randomGenerateCell()

	chunkSize := rows / numWorkers
	for i := 0; i < numWorkers; i++ {
		startY := i * chunkSize
		endY := startY + chunkSize
		if i == numWorkers-1 {
			endY = rows // 最後のチャンクは余りも含める
		}

		wg.Add(1)
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				for x := 0; x < cols; x++ {
					n := liveNeighbors(y, x)
					if grid[y][x] {
						next[y][x] = n == 2 || n == 3
					} else {
						next[y][x] = n == 3
					}
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	grid = next
}


func drawGrid(win *pixelgl.Window) {
	win.Clear(colornames.Black)
	imd := imdraw.New(nil)
	imd.Color = color.White

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if grid[y][x] {
				imd.Push(pixel.V(float64(x*cellSize), float64(y*cellSize)))
				imd.Push(pixel.V(float64((x+1)*cellSize), float64((y+1)*cellSize)))
				imd.Rectangle(0)
			}
		}
	}

	imd.Draw(win)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game of Life",
		Bounds: pixel.R(0, 0, width, height),
		//VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	initGrid()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for !win.Closed() {
		select {
		case <-ticker.C:
			updateGrid()
		default:
		}

		drawGrid(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
