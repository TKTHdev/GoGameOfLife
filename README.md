# 🧬 Go Game of Life

A simple implementation of Conway's Game of Life in Go using the [faiface/pixel](https://github.com/faiface/pixel) graphics library.

## What is this
- Conway's Game of Life(https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) developed in Golang
- Interactive speed adjustment with up and down arrow key
- Interactive cell spwaning by clicking

## GIF

![gif](https://github.com/TKTHdev/GoGameOfLife/blob/master/gameoflife.gif)

## 🧰 Requirements

- Go 1.18 or newer
- OpenGL (included by default on most OSes)
- Git

---

## 🛠 Installation & Running

### 1. Clone this repository

```bash
git clone https://github.com/TKTHdev/GoGameOfLife.git
cd GoGameOfLife
```

### 2. Install dependencies
This project uses Pixel, a 2D game library for Go. To install dependencies:
```bash
go mod tidy
```

## 3. Run
```bash
go run game.go
```
