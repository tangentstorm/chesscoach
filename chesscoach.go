// Unnamed Chess Coach Program

package main

import (
	eb "github.com/hajimehoshi/ebiten"
	eu "github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
)

// images of the white/black chess pieces
var wp, wr, wn, wb, wq, wk, bp, br, bn, bb, bq, bk *eb.Image

func sprite(path string) *eb.Image {
	im, _, err := eu.NewImageFromFile(path, eb.FilterDefault)
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
}

const cellSize = 48 // size of grid cells

func blit(screen, im *eb.Image, gx, gy int) {
	op := &eb.DrawImageOptions{}
	op.GeoM.Translate(float64(gx*cellSize), float64(gy*cellSize))
	screen.DrawImage(im, op)
}

func update(screen *eb.Image) error {
	if eb.IsDrawingSkipped() {
		return nil
	}

	// draw the board
	light := color.RGBA{0, 0, 255, 255}
	dark := color.RGBA{0, 0, 127, 255}

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
	return nil
}

func main() {
	if err := eb.Run(update, 640, cellSize*8, 2, "chesscoach"); err != nil {
		log.Fatal(err)
	}
}
