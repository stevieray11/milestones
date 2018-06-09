package game

import (

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
	"log"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"image"
	"bytes"
	"github.com/enjoykarma/milestones/src/images"
	//c "github.com/enjoykarma/milestones/config"
)

type Game struct {

	PlayerDestinationX int
	PlayerDestinationY int
	Font font.Face
	EventRegistry []string
	CameraX float64
	CameraY float64

}

var (
	Game1 = Game{0, 0, nil, nil, 0, 0}
)

func InitGame() {
		// Fonts
		tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
		if err != nil {
			log.Fatal(err)
		}
		Game1.Font = truetype.NewFace(tt, &truetype.Options{
			Size:    10,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		Game1.CameraX = 0
		Game1.CameraY = 0

		
}

func ExecuteActions(screen *ebiten.Image) {
	
	if (ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
		
		img, _, err := image.Decode(bytes.NewReader(images.Cursor_png))
		if err != nil {
			log.Fatal(err)
		}
		x, y := ebiten.CursorPosition()
		ebitenImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x - 7), float64(y - 7))
		screen.DrawImage(ebitenImage, op)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		// x, y := ebiten.CursorPosition()

	}

}


func GetImage(imageName string) (EImage *ebiten.Image) {

	var err error
 	var img image.Image
	
	switch imageName {

	case "Cursor_png":
		img, _, err = image.Decode(bytes.NewReader(images.Cursor_png))
	case "Colony_png":
		img, _, err = image.Decode(bytes.NewReader(images.Colony_png))
	default:
		img = nil
		err = nil
		break
	}
	
	if err != nil {
		log.Fatal(err)
	}

	EImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	return EImage

}


func (g *Game) CameraMoveLeft() {
	g.CameraX--

}
func (g *Game) CameraMoveRight() {
	g.CameraX++
}
func (g *Game) CameraMoveUp() {
	g.CameraY--
}
func (g *Game) CameraMoveDown() {
	g.CameraY++
}
