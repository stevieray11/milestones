package gamemap

import (
	"math"

	//"fmt"
	"image"
	"math/rand"
	"bytes"
	c "github.com/enjoykarma/milestones/config"
	g "github.com/enjoykarma/milestones/game"
	"github.com/hajimehoshi/ebiten"
	"github.com/enjoykarma/milestones/src/images"
)
var (
	GameMapLayer1 Layer
	GameMapLayer2 Layer
	GameMapLayerGrid Layer
)

var TileFieldStartX int 
var TileFieldStartY int

var TileRedFieldStartX int 
var TileRedFieldStartY int

var TileImage  *ebiten.Image

var cursorX float64
var cursorY float64

var tileX int
var tileY int

var tilesCountX int
var tilesCountY int

type Area struct {
	Type int
	X float64
	Y float64
}

type Layer []Area



var (
	TypesAndTilesCoords = map[int][]int{
		
		0 : {56, 8}, // Просто поле
		1 : {71, 42}, // Поле с текстурой
		2 : {49, 83, 16, 21}, // Бочка (16х21)
		3 : {376, 169, 32, 47}, // Башня 
		4 : {194, 27}, // Прозрачный
		5 : {176, 26}, // Сетка
	}
)


func GetAreaByCoords(x int, y int) (a Area) {
	for _, a := range GameMapLayer2{
		if (a.IsBlockedArea()) {
			continue
		}
		if (a.IsCoordsInArea(x, y) ) {
			return a
		}
	}
	return a
}


func (a *Area) LeftArea() (al Area){
	x := (a.X - c.TileSize) - (c.TileSize / 2)
	y := a.Y + (c.TileSize / 2)
	return GetAreaByCoords(int(x), int(y))
}

func (a *Area) RightArea() (al Area){
	x := (a.X + c.TileSize) - (c.TileSize / 2)
	y := a.Y + (c.TileSize / 2)
	return GetAreaByCoords(int(x), int(y))
}

func (a *Area) UpArea() (al Area){
	x := a.X - (c.TileSize / 2)
	y := a.Y - (c.TileSize / 2)
	return GetAreaByCoords(int(x), int(y))
}

func (a *Area) DownArea() (al Area){
	x := a.X - (c.TileSize / 2)
	y := a.Y + c.TileSize + (c.TileSize / 2)
	return GetAreaByCoords(int(x), int(y))
}



func (a *Area) IsBlockedArea() bool {
	if (a.Type == 2){
		return true
	}
	return false
}

func (a *Area) GetCenterPoint() (X int, Y int) { 
	return int(a.X) - (int(c.TileSize) /2), int(a.Y) + (int(c.TileSize) /2)
}

func (a *Area) GetCenterPointFloat64() (X float64, Y float64) { 
	return a.X - (c.TileSize /2), a.Y + (c.TileSize /2)
}

func (a *Area) IsCoordsInArea(x int, y int) bool {
	if ( ( (x > (int(a.X) - int(c.TileSize))) && (x < int(a.X)) ) && ( (y > int(a.Y)) && (y < (int(a.Y) + int(c.TileSize) )) ) ) {
		return true
	}
	return false
}

func (a *Area) GetDistanceTo(bX float64, bY float64) float64 {
	
	aX, aY := a.GetCenterPointFloat64()
	return math.Sqrt(math.Pow(bX - aX, 2) + math.Pow(bY - aY, 2))
}

func (a *Area) IsCloserToDestination(bX float64, bY float64, DestinationX float64, DestinationY float64) bool {
	
	aX, aY := a.GetCenterPointFloat64()

	Da := math.Sqrt(math.Pow(DestinationX - aX, 2) + math.Pow(DestinationY - aY, 2))
	Db := math.Sqrt(math.Pow(DestinationX - bX, 2) + math.Pow(DestinationY - bY, 2))
	//fmt.Println(Da)
	//fmt.Println(Db)
	if (Da < Db) {
		//fmt.Println(a.X, a.Y, bX, bY, "a is closer than b", Da, Db)
	
		return true
	}
	
	return false
}

func (a * Area) Mark(screen *ebiten.Image) {
		//=============
		img, _, _ := image.Decode(bytes.NewReader(images.Cursor_png))	
		ebitenImage, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		//=============
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(a.X - 16, a.Y)
		screen.DrawImage(ebitenImage, op)
}

// {Area{0, 0, 0}, Area{1, 0, 0}, Area{0, 0, 0}},
// 	{Area{1, 0, 0}, Area{0, 0, 0}, Area{1, 0, 0}},

func GenerateGameMap(types [6]int) (layer Layer){
	

	tilesCount := int((c.ScreenWidth * c.ScreenHeight) / ( int(c.TileSize) * int(c.TileSize) ))
	 
	layer = make([]Area, tilesCount) 

	for i := 0; i < tilesCount; i++ {
		layer[i] = Area{types[rand.Intn(len(types))],0,0}
	}
	
	TileImage = g.GetImage("Colony_png")

	tilesCountX = c.ScreenWidth / int(c.TileSize)
	tilesCountY = c.ScreenHeight / int(c.TileSize)

	//fmt.Print(tilesCount)
	
	
	return layer

}

func (GameMapLayer *Layer) DrawLevelMapLayer(screen *ebiten.Image){
	
	
	cursorX = 0
	cursorY = 0

	tileX = 0
	tileY = 0

	var counter int = 0

	for i, a := range *GameMapLayer {
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(cursorX, cursorY)
		tileX = TypesAndTilesCoords[a.Type][0]
		tileY = TypesAndTilesCoords[a.Type][1]
		r := image.Rect(tileX, tileY, tileX+int(c.TileSize), tileY+int(c.TileSize))
		op.SourceRect = &r
		screen.DrawImage(TileImage, op)
		
		if (counter >= tilesCountX ){
			counter = 0
			cursorY += c.TileSize
			cursorX = 0
		} else {
			counter++
			cursorX += c.TileSize
		}
		(*GameMapLayer)[i].X = cursorX
		(*GameMapLayer)[i].Y = cursorY		
	}


	
	
}

func GetRandomLayers() (layers []int ){
	
	// size := 210
	// tiles := [10]int {1,126,126,126,1,1,1,1,1, 1}
	// layers = make([]int, size)
	// for i := 0; i < size; i++ {
	// 	layers[i] = tiles[rand.Intn(len(tiles))]
	// }

	// return
	return
}