package textwrapper

import(
	"log"
	"bytes"
	"strings"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	
	"golang.org/x/text/language"
)

type TextWrapper struct {
	X, Y, W float64
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
		op.GeoM.Translate(t.X, t.Y + offset)
		text.Draw(screen, tx, t.face, op)
	}
}

//Actions 
//Set text allows you to use mutltiple lines of input in the form of a
//string slice. New lines will be inserted where needed as the text
//wraps across the bounds set by the W (width) property of the TextWrapper
func (t *TextWrapper) SetText(newText []string) {
	t.originalText = newText
	var output []string
	
	for _, textLine := range newText {
		maxWidth := t.W - 5
		w, _ := text.Measure(textLine, t.face, t.lineHeight)
		var currentText string = textLine
		
		if maxWidth < w {
			loopDone := false
			for !loopDone {
				textWidth, _ := text.Measure(currentText, t.face, t.lineHeight)
				if maxWidth < textWidth {
					var out string
					out, currentText = t.Split(currentText, textWidth)
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

func (t *TextWrapper) Split(newText string, textWidth float64) (string, string) {
	maxWidth := t.W - 5
	diff := maxWidth/float64(textWidth)
	length := len([]rune(newText))
	newLength := int(float64(length)*diff)
	lastIndex := strings.LastIndex(newText[:newLength], " ")
	return newText[:lastIndex], newText[lastIndex+1:]
}

func (t *TextWrapper) SetSize(newSize float64, lineHeight float64) {
	t.face.Size = newSize
	t.lineHeight = lineHeight
	t.SetText(t.originalText)
}