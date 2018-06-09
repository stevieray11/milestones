package main

import (
	"fmt"
	"strconv"
	"bytes"
	//"fmt"
	"image"
	"log"
	"image/color"
	//"math/rand"
	"github.com/hajimehoshi/ebiten"
	//"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/enjoykarma/milestones/src/images"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/enjoykarma/milestones/unit"
	c "github.com/enjoykarma/milestones/config"
	g "github.com/enjoykarma/milestones/game"
	m "github.com/enjoykarma/milestones/gamemap"
	
)


var (
	tilesImage *ebiten.Image
	player1 = unit.Unit{"Player1", 100, 120, 120, 0, 1, 0, 0, []unit.Checkpoint{}}
	unit1 = unit.Unit{}

)


func init() {

	img, _, err := image.Decode(bytes.NewReader(images.Snow_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)


	types1 := [6]int {0, 0, 0, 0, 1, 1}
	types2 := [6]int {4, 4, 4, 4, 4, 2}
	types3 := [6]int {5, 5, 5, 5, 5, 5}

	m.GameMapLayer1 = m.GenerateGameMap(types1)
	m.GameMapLayer2 = m.GenerateGameMap(types2) 
	
	m.GameMapLayerGrid = m.GenerateGameMap(types3) 
	
	
	unit.InitUnits()
	g.InitGame()

	

	unit1 = unit.Unit{"Unit 1", 100, 120, 120, 0, 1, 0, 0, []unit.Checkpoint{}}

}


var layers []int

func update(screen *ebiten.Image) error {
	
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})

	// Draw map layers
	 
	m.GameMapLayer1.DrawLevelMapLayer(screen)
	m.GameMapLayer2.DrawLevelMapLayer(screen)
	m.GameMapLayerGrid.DrawLevelMapLayer(screen)

	
	unit1.ExecuteActionsUnit(screen)
	
	player1.ExecuteActionsPlayer1(screen)
	
	g.ExecuteActions(screen)

	
	ebitenutil.DebugPrint(screen, fmt.Sprintf(strconv.Itoa(player1.CoordX) + " : " + strconv.Itoa(player1.CoordY) + "\n" + strconv.Itoa(player1.DestinationX) + " : " + strconv.Itoa(player1.DestinationY)))

	return nil
}

func main() {

	//layers = getRandomLayers()
	if err := ebiten.Run(update, c.ScreenWidth, c.ScreenHeight, 2, "Tiles (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}

	fmt.Print(m.GameMapLayer1)
	// test := g.GetImage("Cursor_png")
	// fmt.Print(test)

}



func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}




