package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	width    = 1024
	height   = 768
	cellSize = 10
	rows     = height / cellSize
	cols     = width / cellSize
)

var (
	grid       [rows][cols]bool
	frameDelay = 50 * time.Millisecond
)

func initGrid() {
	rand.Seed(time.Now().UnixNano())
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if rand.Float64() < 0.2 {
				grid[y][x] = true
			} else {
				grid[y][x] = false
			}
		}
	}
}

func liveNeighbors(y, x int) int {
	count := 0
	// Pre-calculate wrapped coordinates
	ny := (y + rows) % rows
	nx := (x + cols) % cols

	// Check all 8 neighbors
	if grid[(ny-1+rows)%rows][(nx-1+cols)%cols] {
		count++
	}
	if grid[(ny-1+rows)%rows][nx] {
		count++
	}
	if grid[(ny-1+rows)%rows][(nx+1)%cols] {
		count++
	}
	if grid[ny][(nx-1+cols)%cols] {
		count++
	}
	if grid[ny][(nx+1)%cols] {
		count++
	}
	if grid[(ny+1)%rows][(nx-1+cols)%cols] {
		count++
	}
	if grid[(ny+1)%rows][nx] {
		count++
	}
	if grid[(ny+1)%rows][(nx+1)%cols] {
		count++
	}

	return count
}

func updateGrid() {
	var next [rows][cols]bool
	var wg sync.WaitGroup
	numWorkers := 8
	chunkSize := rows / numWorkers
	for i := 0; i < numWorkers; i++ {
		startY := i * chunkSize
		endY := startY + chunkSize
		if i == numWorkers-1 {
			endY = rows
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

	var wg sync.WaitGroup
	numWorkers := 8
	chunkSize := rows / numWorkers

	// Create a channel to collect all imdraw objects
	imdChan := make(chan *imdraw.IMDraw, numWorkers)

	for i := 0; i < numWorkers; i++ {
		startY := i * chunkSize
		endY := startY + chunkSize
		if i == numWorkers-1 {
			endY = rows
		}
		wg.Add(1)
		go func(startY, endY int) {
			defer wg.Done()
			// Create a separate imdraw object for each worker
			workerImd := imdraw.New(nil)
			workerImd.Color = color.White

			// Batch the drawing operations
			for y := startY; y < endY; y++ {
				for x := 0; x < cols; x++ {
					if grid[y][x] {
						workerImd.Push(pixel.V(float64(x*cellSize), float64(y*cellSize)))
						workerImd.Push(pixel.V(float64((x+1)*cellSize), float64((y+1)*cellSize)))
						workerImd.Rectangle(0)
					}
				}
			}
			imdChan <- workerImd
		}(startY, endY)
	}

	// Wait for all workers to complete
	wg.Wait()
	close(imdChan)

	// Combine all imdraw objects
	finalImd := imdraw.New(nil)
	finalImd.Color = color.White
	for workerImd := range imdChan {
		// Draw each worker's imdraw object to the final one
		workerImd.Draw(finalImd)
	}

	// Draw the final result
	finalImd.Draw(win)
}

func drawFrameDelay(txt *text.Text, win *pixelgl.Window) {
	txt.Clear()
	fmt.Fprintf(txt, "Frame Delay: %d ms", frameDelay.Milliseconds())
	txt.Draw(win, pixel.IM)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game of Life",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	initGrid()
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(10, height-20), atlas)
	for !win.Closed() {
		if win.Pressed(pixelgl.MouseButtonLeft) {
			mouse := win.MousePosition()
			x := int(mouse.X) / cellSize
			y := int(mouse.Y) / cellSize
			if x >= 0 && x < cols && y >= 0 && y < rows {
				grid[y][x] = true
			}
		}
		if win.Pressed(pixelgl.KeyUp) {
			frameDelay -= 10 * time.Millisecond
			if frameDelay < 10*time.Millisecond {
				frameDelay = 10 * time.Millisecond
			}
		}
		if win.Pressed(pixelgl.KeyDown) {
			frameDelay += 10 * time.Millisecond
			if frameDelay > 1000*time.Millisecond {
				frameDelay = 1000 * time.Millisecond
			}
		}
		updateGrid()
		drawGrid(win)
		drawFrameDelay(txt, win)
		win.Update()
		time.Sleep(frameDelay)
	}
}

func main() {
	pixelgl.Run(run)
}
