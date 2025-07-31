package main

import (
	_ "embed"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	
	"github.com/RileySun/TextWrapper"
)

const GAMEWIDTH = 960
const GAMEHEIGHT = 640

//go:embed Font.ttf
var FONTDATA []byte
var keyTimer float64

type Game struct {
	wrapper *textwrapper.TextWrapper
	scroll *textwrapper.TextWrapper
}

func main() {
	ebiten.SetWindowTitle("TextWrapper")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := &Game{
		wrapper:textwrapper.NewTextWrapper(GAMEWIDTH, -1, FONTDATA),
		scroll:textwrapper.NewTextWrapper(GAMEWIDTH, 240, FONTDATA),
	}
	
	//Regular TextWrapper
	output := []string{
		"This is a short line\n",
		"This is a much longer line, it should wrap until the edge of the screen and then get turned into a new line. You can even go a little further.\n",
		"This line is to show what happens after the aforementioned line wrapping.",
	}
	game.wrapper.SetText(output)
	game.wrapper.SetSize(30, 35)
	

	//Scroll Wrapper
	game.scroll.Y = 320
	scrollOutput := []string {
		"This is a scroll text.\n",
		"Scroll texts are overly long and usually will go off the page.",
		"Using the arrow keys on your keyboard you should be able to scroll this text up and down.",
		"Scrolling text can be useful for long peices of text like stories or lore.\n",
		"Hopefully you will be able to make use of this new scrolling and wrapping text utility.\n",
		"Language shapes the way we think, and determines what we can think about.",
		"-Benjamin Lee Whorf",
	}
	game.scroll.SetText(scrollOutput)
	game.scroll.SetSize(30, 35)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func (g *Game) Update() error {
	keyTimer += 0.1
	if key(ebiten.KeyUp) && keyTimer > 1 {
		g.scroll.ScrollUp()
		keyTimer = 0
	}
	
	if key(ebiten.KeyDown) && keyTimer > 1 {
		g.scroll.ScrollDown()
		keyTimer = 0
	}

	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	g.wrapper.Draw(screen)
	vector.StrokeLine(screen, 0, 300, GAMEWIDTH, 300, 1, g.wrapper.Color, false)
	g.scroll.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GAMEWIDTH, GAMEHEIGHT
}

func key(k ebiten.Key) bool {
	return ebiten.IsKeyPressed(k)
}