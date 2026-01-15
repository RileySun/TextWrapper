package textwrapper

import(
	"log"
	"bytes"
	"bufio"
	"strings"
	"image/color"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	
	"golang.org/x/text/language"
)

type TextWrapper struct {
	X, Y, W, H float64 //Width must be greater than 0, 
	Color color.NRGBA //Any image.Color should work
	Scroll, ScrollMax, ScrollCurrentMax, ScrollVisible int //Based off height (H)
	size, lineHeight float64 //Font Size & Line Height
	face *text.GoTextFace
	faceSource *text.GoTextFaceSource
	originalText, finalText []string
}

//Create
func NewTextWrapper(width float64, height float64, fontData []byte) *TextWrapper {
	//Make sure width and height are correct
	if int(width) <= 0 {
		log.Fatal("TextWrapper - NewTextWrapper: ", "TextWrapper width must be greater than 0.")
	}
	if int(height) == 0 || int(height) < -1 {
		log.Fatal("TextWrapper - NewTextWrapper: ", "TextWrapper height must be greater than 0 or equal -1 for infinite height.")
	}
	
	//Get Face Source
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal("TextWrapper - NewTextWrapper: ", err.Error())
	}
	
	//Create
	wrapper := &TextWrapper {
		X:0,Y:0,W:width, H:height,
		size:18,
		faceSource:source,
		face: &text.GoTextFace{
			Source:    source,
			Direction: text.DirectionLeftToRight,
			Size:      18,
			Language:  language.English,
		},
		Color:color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	}
	
	wrapper.lineHeight = wrapper.size * 4.5
	
	return wrapper
} //Width must be greater than 0, height can be -1 for infinite, but needs to not be 0

//Render
func (t *TextWrapper) Draw(screen *ebiten.Image) {
	startIndex := 0
	for i := t.Scroll; i < t.ScrollCurrentMax; i++ {
		offset := t.lineHeight * float64(startIndex)
		op := &text.DrawOptions{}
		op.GeoM.Translate(t.X, t.Y + offset)
		op.ColorScale.Scale(NRGBAtoFloat32(t.Color))
		text.Draw(screen, t.finalText[i], t.face, op)
		startIndex++
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

func (t *TextWrapper) calculateScroll() {
	t.Scroll, t.ScrollMax = 0, len(t.finalText)
	
	//If infinite, skip the rest
	if t.H == -1 {
		t.ScrollVisible = t.ScrollMax
		t.ScrollCurrentMax = t.ScrollMax
		return
	}
	
	//Calculate how many lines can be visible using the height of the wrapper
	_, singleLineHeight := text.Measure("Example", t.face, t.lineHeight)
	t.ScrollVisible = int(t.H/singleLineHeight)
	t.ScrollCurrentMax = t.Scroll + t.ScrollVisible
}

//Actions 
//Set text allows you to use mutltiple lines of input in the form of a
//string slice. New lines will be inserted where needed as the text
//wraps across the bounds set by the W (width) property of the TextWrapper
func (t *TextWrapper) SetText(newText []string) {
	//Find newlines first
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
	
	//errors with only one line of text
	if len(t.finalText) == 1 {
		t.finalText = append(t.finalText, "")
	}
	
	t.calculateScroll()
}

func (t *TextWrapper) SetTextFromBytes(byt []byte) {
	scanner := bufio.NewScanner(bufio.NewReader(bytes.NewReader(byt)))
	scanner.Split(bufio.ScanLines)
	
	var textLines []string
	for scanner.Scan() {
		textLine := scanner.Text()
		if textLine != "" {
			textLines = append(textLines, textLine + "\n")	
		}
	}
	
	t.SetText(textLines)
} //Set text from a loaded byte slice (file)

func (t *TextWrapper) SetSize(newSize float64, lineHeight float64) {
	t.face.Size = newSize
	t.lineHeight = lineHeight
	t.SetText(t.originalText)
}

func (t *TextWrapper) GetFace() *text.GoTextFace {
	return t.face
}

func (t *TextWrapper) ScrollUp(scrollAmount int) {
	//Dont scroll if not a scrolling Wrapper
	if t.H == -1 {
		return
	}  

	if t.Scroll - scrollAmount > 0 {
		t.ScrollCurrentMax -= scrollAmount
		t.Scroll -= scrollAmount
	} else {
		t.Scroll = 0
		t.ScrollCurrentMax = t.ScrollVisible
	}
}//Scroll up by scroll amount if possible

func (t *TextWrapper) ScrollDown(scrollAmount int) {
	//Dont scroll if not a scrolling Wrapper
	if t.H == -1 {
		return
	}  

	if t.Scroll + scrollAmount < t.ScrollMax - t.ScrollVisible {
		t.ScrollCurrentMax += scrollAmount
		t.Scroll += scrollAmount
	}
}//Scroll down by scroll amount if possible

//Color Util
func NRGBAtoFloat32(newColor color.NRGBA) (float32, float32, float32, float32) {
	r := (float32(newColor.R) + 0.5)/256
	g := (float32(newColor.G) + 0.5)/256
	b := (float32(newColor.B) + 0.5)/256
	a := float32(1)
	return r, g, b, a
}