// Unnamed Chess Coach Program

package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	eb "github.com/hajimehoshi/ebiten"
	ebu "github.com/hajimehoshi/ebiten/ebitenutil"
	ebf "github.com/hajimehoshi/ebiten/examples/resources/fonts"
	ebi "github.com/hajimehoshi/ebiten/inpututil"
	ebt "github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
)

const cellSize = 48 // size of grid cells

type Square struct{ x, y int }

type marker int

const (
	p0 marker = iota // player's starting square
	p1               // ... ending square
	o0               // opponent's starting square
	o1               // ... ending square
)

// images of the white/black chess pieces
var wp, wr, wn, wb, wq, wk, bp, br, bn, bb, bq, bk *eb.Image
var mainFont font.Face
var files = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}
var nowhere = Square{-1, -1}
var markers = map[marker]Square{
	p0: nowhere,
	p1: nowhere,
	o0: nowhere,
	o1: nowhere,
}

func sprite(path string) *eb.Image {
	im, _, err := ebu.NewImageFromFile(path, eb.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return im
}

func init() {
	wp = sprite("sprites/wp.png")
	wr = sprite("sprites/wr.png")
	wn = sprite("sprites/wn.png")
	wb = sprite("sprites/wb.png")
	wq = sprite("sprites/wq.png")
	wk = sprite("sprites/wk.png")
	bp = sprite("sprites/bp.png")
	br = sprite("sprites/br.png")
	bn = sprite("sprites/bn.png")
	bb = sprite("sprites/bb.png")
	bq = sprite("sprites/bq.png")
	bk = sprite("sprites/bk.png")

	// set up the font
	tt, err := truetype.Parse(ebf.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mainFont = truetype.NewFace(tt, &truetype.Options{
		Size:    8,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func blit(screen, im *eb.Image, gx, gy int) {
	op := &eb.DrawImageOptions{}
	op.GeoM.Translate(float64(gx*cellSize), float64(gy*cellSize))
	screen.DrawImage(im, op)
}

func mouseSquare() Square {
	x, y := eb.CursorPosition()
	x, y = x/cellSize, y/cellSize
	if x >= 0 && y >= 0 && x <= 7 && y <= 7 {
		return Square{x: x, y: y}
	}
	return nowhere
}

func update(screen *eb.Image) error {

	if ebi.IsMouseButtonJustPressed(eb.MouseButtonLeft) {
		sq := mouseSquare()
		switch {
		case markers[p0] == nowhere:
			markers[p0] = sq
		case markers[p0] == sq:
			markers[p0] = nowhere
		default:
			// TODO: actually make the move
		}
	}

	if eb.IsDrawingSkipped() {
		return nil
	}

	// draw the board
	light := color.RGBA{0, 0, 255, 255}
	dark := color.RGBA{0, 0, 127, 255}
	highlight := color.RGBA{224, 164, 0, 63}

	c := &light
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			x0, y0 := x*cellSize, y*cellSize
			sq := image.Rect(x0, y0, x0+cellSize-1, y0+cellSize-1)
			if (y & 1) == (x & 1) {
				c = &light
			} else {
				c = &dark
			}
			if x == markers[p0].x && y == markers[p0].y {
				c = &highlight
			}
			draw.Draw(screen, sq, &image.Uniform{c}, image.ZP, draw.Src)
		}
	}

	// draw opening board
	blit(screen, br, 0, 0)
	blit(screen, bn, 1, 0)
	blit(screen, bb, 2, 0)
	blit(screen, bq, 3, 0)
	blit(screen, bk, 4, 0)
	blit(screen, bb, 5, 0)
	blit(screen, bn, 6, 0)
	blit(screen, br, 7, 0)
	blit(screen, wr, 0, 7)
	blit(screen, wn, 1, 7)
	blit(screen, wb, 2, 7)
	blit(screen, wq, 3, 7)
	blit(screen, wk, 4, 7)
	blit(screen, wb, 5, 7)
	blit(screen, wn, 6, 7)
	blit(screen, wr, 7, 7)
	for i := 0; i < 8; i++ {
		blit(screen, wp, i, 6)
		blit(screen, bp, i, 1)
	}

	const textX = cellSize*8 + 16
	ebt.Draw(screen, "Hello. I am chess coach.", mainFont, textX, 16, color.White)
	ebt.Draw(screen, "You are playing white today.", mainFont, textX, 32, color.White)
	ebt.Draw(screen, "Use mouse to select move.", mainFont, textX, 48, color.White)

	ms := mouseSquare()
	if ms != nowhere {
		f := files[ms.x]
		r := 8 - ms.y
		ebt.Draw(screen, fmt.Sprintf("Mouse is over %c%d", f, r), mainFont, textX, 96, color.White)
	}

	return nil
}

func main() {
	if err := eb.Run(update, 640, cellSize*8, 2, "chesscoach"); err != nil {
		log.Fatal(err)
	}
}
