package main

import (
	_ "embed"
	
	"github.com/hajimehoshi/ebiten/v2"
	
	"github.com/RileySun/TextWrapper"
)

const GAMEWIDTH = 960
const GAMEHEIGHT = 640

//go:embed Font.ttf
var FONTDATA []byte

type Game struct {
	wrapper *TextWrapper
}

func main() {
	ebiten.SetWindowTitle("TextWrapper")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := &Game{
		wrapper:textwrapper.NewTextWrapper(FONTDATA),
	}
	
	game.wrapper.W = GAMEWIDTH
	output := []string{
		"This is a short line",
		"This is a much longer line, it should wrap until the edge of the screen and then get turned into a new line. You can even go a little further.",
		"This line is to show what happens after the aforementioned line wrapping.",
	}
	game.wrapper.SetText(output)
	game.wrapper.SetSize(20, 28)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func (g *Game) Update() error {
	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	g.wrapper.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GAMEWIDTH, GAMEHEIGHT
}