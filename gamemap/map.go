package gamemap

import (

	"fmt"
	"image"
	"math/rand"
	// u "github.com/enjoykarma/milestones/unit"
	c "github.com/enjoykarma/milestones/config"
	g "github.com/enjoykarma/milestones/game"
	"github.com/hajimehoshi/ebiten"
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




func (a *Area) IsBlockedArea() bool {
	if (a.Type == 2){
		return true
	}
	return false
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

	fmt.Print(tilesCount)
	
	
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