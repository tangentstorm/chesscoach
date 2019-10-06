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
	"golang.org/x/image/font"
	eb "github.com/hajimehoshi/ebiten"
	ebu "github.com/hajimehoshi/ebiten/ebitenutil"
	ebf "github.com/hajimehoshi/ebiten/examples/resources/fonts"
	ebi "github.com/hajimehoshi/ebiten/inpututil"
	ebt "github.com/hajimehoshi/ebiten/text"
	"github.com/notnil/chess"
)

const cellSize = 48 // size of grid cells

var game *chess.Game
var p0 = chess.NoSquare  // player's selected square

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
	game = chess.NewGame()
	icons[chess.WhitePawn] = sprite("sprites/wp.png")
	icons[chess.WhiteRook] = sprite("sprites/wr.png")
	icons[chess.WhiteKnight] = sprite("sprites/wn.png")
	icons[chess.WhiteBishop] = sprite("sprites/wb.png")
	icons[chess.WhiteQueen] = sprite("sprites/wq.png")
	icons[chess.WhiteKing] = sprite("sprites/wk.png")
	icons[chess.BlackPawn] = sprite("sprites/bp.png")
	icons[chess.BlackRook] = sprite("sprites/br.png")
	icons[chess.BlackKnight] = sprite("sprites/bn.png")
	icons[chess.BlackBishop] = sprite("sprites/bb.png")
	icons[chess.BlackQueen] = sprite("sprites/bq.png")
	icons[chess.BlackKing] = sprite("sprites/bk.png")

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

func squareAt(x, y int) chess.Square {
	if x < 0 || x > 7 || y < 0 || y > 7 {
		return chess.NoSquare
	}
	return chess.Square(int(8*(7-y))+x)
}

func mouseSquare() chess.Square {
	x, y := eb.CursorPosition()
	if x < 0 || y < 0 {
		return chess.NoSquare
	}
	x, y = x/cellSize, y/cellSize
	return squareAt(x, y)
}

func validSquares() (result []chess.Square) {
	counts := make(map[chess.Square] int)
	for _, mv := range game.ValidMoves() {
		if p0 == chess.NoSquare {
			counts[mv.S1()]++
		} else if p0 == mv.S1() {
			counts[mv.S2()]++
		}
	}
	for sq := range counts {
		result = append(result, sq)
	}
	return result
}

func watchMouse() {
	if ebi.IsMouseButtonJustPressed(eb.MouseButtonLeft) {
		sq := mouseSquare()
		switch {
		case p0 == chess.NoSquare:
			p0 = sq
		case p0 == sq:
			p0 = chess.NoSquare
		default:
			// TODO: actually make the move
		}
	}
}

func drawSquareXY(screen *eb.Image, x, y int, c color.Color) {
	x0, y0 := x*cellSize, (7-y)*cellSize
	rect := image.Rect(x0, y0, x0+cellSize-1, y0+cellSize-1)
	draw.Draw(screen, rect, &image.Uniform{c}, image.ZP, draw.Src)
}

func drawSquare(screen *eb.Image, sq chess.Square, c color.Color) {
	drawSquareXY(screen, int(sq.File()), int(sq.Rank()), c)
}

func drawBoard(screen *eb.Image) {
	light := color.RGBA{0, 0, 255, 255}
	dark := color.RGBA{0, 0, 127, 255}
	var c color.Color
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if (y & 1) == (x & 1) {
				c = &light
			} else {
				c = &dark
			}
			drawSquareXY(screen, x, 7-y, c)
		}
	}
}

func drawMarks(screen *eb.Image) {
	highlight := color.RGBA{224, 164, 0, 63}
	validLight := color.RGBA{63, 63, 127, 63}
	validDark := color.RGBA{32, 32, 64, 63}
	for _, sq := range validSquares() {
		if int(sq.Rank())&1 == int(sq.File())&1  {
			drawSquare(screen, sq, validDark)
		} else {
			drawSquare(screen, sq, validLight)
		}
	}
	if p0 != chess.NoSquare {
		drawSquare(screen, p0, highlight)
	}
}

func drawPieces(screen *eb.Image) {
	for sq, p := range game.Position().Board().SquareMap() {
		blit(screen, icons[p], int(sq.File()), 7-int(sq.Rank()))
	}
}

func drawText(screen *eb.Image) {
	const textX = cellSize*8 + 16
	ebt.Draw(screen, "Hello. I am chess coach.", mainFont, textX, 16, color.White)
	ebt.Draw(screen, "You are playing white today.", mainFont, textX, 32, color.White)
	ebt.Draw(screen, "Use mouse to select move.", mainFont, textX, 48, color.White)
	if ms := mouseSquare(); ms != chess.NoSquare {
		ebt.Draw(screen, fmt.Sprintf("Mouse is over %s", ms), mainFont, textX, 96, color.White)
	}
}

func update(screen *eb.Image) error {
	watchMouse()
	if ! eb.IsDrawingSkipped() {
		drawBoard(screen)
		drawMarks(screen)
		drawPieces(screen)
		drawText(screen)
	}
	return nil
}

func main() {
	if err := eb.Run(update, 640, cellSize*8, 2, "chesscoach"); err != nil {
		log.Fatal(err)
	}
}
