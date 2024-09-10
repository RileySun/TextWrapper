package textwrapper

import(
	"log"
	"bytes"
	"strings"
	"image/color"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	
	"golang.org/x/text/language"
)

type TextWrapper struct {
	X, Y, W float64
	Color color.Color
	size, lineHeight float64 
	face *text.GoTextFace
	faceSource *text.GoTextFaceSource
	originalText, finalText []string
}

//Create
func NewTextWrapper(fontData []byte) *TextWrapper {
	//Get Face Source
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}
	
	//Create
	wrapper := &TextWrapper {
		X:0,Y:0,W:20,
		size:18,
		faceSource:source,
		face: &text.GoTextFace{
			Source:    source,
			Direction: text.DirectionLeftToRight,
			Size:      18,
			Language:  language.English,
		},
	}
	
	wrapper.lineHeight = wrapper.size * 4.5
	
	return wrapper
}

//Render
func (t *TextWrapper) Draw(screen *ebiten.Image) {
	for i, tx := range t.finalText {
		offset := t.lineHeight * float64(i)
		op := &text.DrawOptions{}
		//Translate
		op.GeoM.Translate(t.X, t.Y + offset)
		//Color
		if t.Color != nil {
			r, g, b, a := t.Color.RGBA()
			op.ColorM.Scale(float64(r), float64(g), float64(b), float64(a))
		}
		
		text.Draw(screen, tx, t.face, op)
	}
}

//Utils
func (t *TextWrapper) findNewLines(rawText []string) []string {
	var out []string
	for _, s := range rawText {
		split := strings.Split(s ,"\n")
		out = append(out, split...)
	}
	return out
}

func (t *TextWrapper) split(newText string, textWidth float64) (string, string) {
	maxWidth := t.W - 5
	diff := maxWidth/float64(textWidth)
	length := len([]rune(newText))
	newLength := int(float64(length)*diff)
	lastIndex := strings.LastIndex(newText[:newLength], " ")
	
	//What if there are no spaces?
	if lastIndex == -1 {
		return newText[:newLength], newText[newLength:]
	}
	
	return newText[:lastIndex], newText[lastIndex+1:]
}

//Actions 
//Set text allows you to use mutltiple lines of input in the form of a
//string slice. New lines will be inserted where needed as the text
//wraps across the bounds set by the W (width) property of the TextWrapper
func (t *TextWrapper) SetText(newText []string) {
	//Find newliens first
	newLineText := t.findNewLines(newText)
	
	
	t.originalText = newLineText
	var output []string
	
	for _, textLine := range newLineText {
		maxWidth := t.W - 5
		w, _ := text.Measure(textLine, t.face, t.lineHeight)
		var currentText string = textLine
		
		if maxWidth < w {
			loopDone := false
			for !loopDone {
				textWidth, _ := text.Measure(currentText, t.face, t.lineHeight)
				if maxWidth < textWidth {
					var out string
					out, currentText = t.split(currentText, textWidth)
					output = append(output, out)
				} else {
					output = append(output, currentText)
					loopDone = true
				}
			}
		} else {
			output = append(output, textLine)
		}
	}
	
	t.finalText = output
}

func (t *TextWrapper) SetSize(newSize float64, lineHeight float64) {
	t.face.Size = newSize
	t.lineHeight = lineHeight
	t.SetText(t.originalText)
}

func (t *TextWrapper) GetFace() *text.GoTextFace {
	return t.face
}