// Unnamed Chess Coach Program

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"

	"github.com/golang/freetype/truetype"
	eb "github.com/hajimehoshi/ebiten"
	ebu "github.com/hajimehoshi/ebiten/ebitenutil"
	ebf "github.com/hajimehoshi/ebiten/examples/resources/fonts"
	ebi "github.com/hajimehoshi/ebiten/inpututil"
	ebt "github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	lib "github.com/tangentstorm/chesscoach/lib"
)

const cellSize = 48 // size of grid cells

// Marker is an enum used for annotating the board in the UI.
type marker int

const (
	p0 marker = iota // player's starting square
	p1               // ... ending square
	o0               // opponent's starting square
	o1               // ... ending square
)

var board lib.Board
var markers = map[marker]lib.Square{
	p0: lib.Nowhere,
	p1: lib.Nowhere,
	o0: lib.Nowhere,
	o1: lib.Nowhere,
}

var mainFont font.Face

// images of the white/black chess pieces
var icons [13]*eb.Image // 13 = card(piece)

func sprite(path string) *eb.Image {
	im, _, err := ebu.NewImageFromFile(path, eb.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return im
}

func init() {
	board = lib.StartPos()
	icons[lib.WP] = sprite("sprites/wp.png")
	icons[lib.WR] = sprite("sprites/wr.png")
	icons[lib.WN] = sprite("sprites/wn.png")
	icons[lib.WB] = sprite("sprites/wb.png")
	icons[lib.WQ] = sprite("sprites/wq.png")
	icons[lib.WK] = sprite("sprites/wk.png")
	icons[lib.BP] = sprite("sprites/bp.png")
	icons[lib.BR] = sprite("sprites/br.png")
	icons[lib.BN] = sprite("sprites/bn.png")
	icons[lib.BB] = sprite("sprites/bb.png")
	icons[lib.BQ] = sprite("sprites/bq.png")
	icons[lib.BK] = sprite("sprites/bk.png")

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

func mouseSquare() lib.Square {
	x, y := eb.CursorPosition()
	x, y = x/cellSize, y/cellSize
	return lib.SquareAt(x, y)
}

func update(screen *eb.Image) error {

	if ebi.IsMouseButtonJustPressed(eb.MouseButtonLeft) {
		sq := mouseSquare()
		switch {
		case markers[p0] == lib.Nowhere:
			markers[p0] = sq
		case markers[p0] == sq:
			markers[p0] = lib.Nowhere
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

	var c color.Color
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			x0, y0 := x*cellSize, y*cellSize
			sq := image.Rect(x0, y0, x0+cellSize-1, y0+cellSize-1)
			if (y & 1) == (x & 1) {
				c = &light
			} else {
				c = &dark
			}
			if markers[p0] == lib.SquareAt(x,y) {
				c = &highlight
			}
			draw.Draw(screen, sq, &image.Uniform{c}, image.ZP, draw.Src)
		}
	}

	// draw opening board
	for y, file := range board {
		for x, p := range file {
			if p > lib.NO {
				blit(screen, icons[p], x, y)
			}
		}
	}

	const textX = cellSize*8 + 16
	ebt.Draw(screen, "Hello. I am chess coach.", mainFont, textX, 16, color.White)
	ebt.Draw(screen, "You are playing white today.", mainFont, textX, 32, color.White)
	ebt.Draw(screen, "Use mouse to select move.", mainFont, textX, 48, color.White)

	ms := mouseSquare()
	if ms != lib.Nowhere {
		ebt.Draw(screen, fmt.Sprintf("Mouse is over %s", ms.Name()), mainFont, textX, 96, color.White)
	}
	return nil
}

func main() {
	if err := eb.Run(update, 640, cellSize*8, 2, "chesscoach"); err != nil {
		log.Fatal(err)
	}
}
