package unit

import (
	//"reflect"
	"fmt"
	"strconv"
	"image"
	"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"bytes"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
	"log"
	//"math/rand"
	c "github.com/enjoykarma/milestones/config"
	g "github.com/enjoykarma/milestones/game"
	m "github.com/enjoykarma/milestones/gamemap"
)

var (
	playerImage *ebiten.Image
)

type Unit struct {

	Name string
	Health int
	CoordX int
	CoordY int
	TileX int
	Speed int
	DestinationX int
	DestinationY int
	Route []Checkpoint
}

type Checkpoint struct {
	X int
	Y int
}

func (p *Unit) DrawStatic(screen *ebiten.Image, tileX int, tileY int, positionX float64, positionY float64) {
	
	op := &ebiten.DrawImageOptions{}
	positionX = positionX - (c.PlayertileWidth/2)
	positionY = positionY - (c.PlayertileHeight/2)

	op.GeoM.Translate(positionX, positionY)

	sx := tileX * int(c.PlayertileWidth)
	sy := tileY * int(c.PlayertileHeight)

	r := image.Rect(sx, sy, sx+int(c.PlayertileWidth), sy+int(c.PlayertileHeight))
	op.SourceRect = &r
	screen.DrawImage(playerImage, op)
	// fmt.Print(g.Game1)
	
	text.Draw(screen, p.Name, g.Game1.Font, int(positionX) - 10, int(positionY), color.White)
}

func InitUnits() {

	//Player1
	img, _, err := image.Decode(bytes.NewReader(images.Player_png))
	if err != nil {
		log.Fatal(err)
	}
	playerImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)


}

func (p *Unit) MovementPossibility() bool{
	

	return true
}

func (p *Unit) TerminateAction(){
	
	p.DestinationX = 0
	p.DestinationY = 0
}

func (p *Unit) IsDestinationOnBlockedArea(customCoords bool, x int, y int) bool{
	
	if (!customCoords) {
		x = p.DestinationX
		y = p.DestinationY
	}
	
	var (
		Xleft int
		Xright int
		Ybottom int
		Ytop int
	)
	for _, a := range m.GameMapLayer2{
		Ytop = int(a.Y) 
		Ybottom = int(a.Y) + int(c.TileSize)
		Xleft = int(a.X) - int(c.TileSize) 
		Xright = int(a.X)
		if (a.IsBlockedArea() && (y >= Ytop && y <= Ybottom ) && (x >= Xleft && x <= Xright)){
			
			if (customCoords) {
				//fmt.Println("CUSTOM DESTINATION ON BLOCKED!!!")
			} else {
				//fmt.Println("PLAYER DESTINATION ON BLOCKED!!!")
			}
			

			return true
		}
		
	}
	return false
	
}


func (p *Unit) ChangeCoord(coordType string, target int) (r bool){
	
	var (
		Xleft int
		Xright int
		Ybottom int
		Ytop int
		//changedCoord int
	)

	switch coordType {
		
	case "X" :
		for _, a := range m.GameMapLayer2{
			Xleft = int(a.X) - int(c.TileSize) - c.PlayerSpaceFromBlockSize
			Xright = int(a.X) + c.PlayerSpaceFromBlockSize
			Ytop = int(a.Y) - c.PlayerSpaceFromBlockSize
			Ybottom = int(a.Y) + int(c.TileSize) + c.PlayerSpaceFromBlockSize
			if (a.IsBlockedArea() && (target >= Xleft && target <= Xright) && (p.CoordY >= Ytop && p.CoordY <= Ybottom) ) {
				if (p.IsDestinationOnBlockedArea(false, 0, 0)){
					p.TerminateAction()
				}
				return false
			}
		}	
		
		p.CoordX = target
		return true

		break

	case "Y":
		for _, a := range m.GameMapLayer2{
			Ytop = int(a.Y) - c.PlayerSpaceFromBlockSize
			Ybottom = int(a.Y) + int(c.TileSize) + c.PlayerSpaceFromBlockSize
			Xleft = int(a.X) - int(c.TileSize) - c.PlayerSpaceFromBlockSize
			Xright = int(a.X) + c.PlayerSpaceFromBlockSize
			if (a.IsBlockedArea() && (target >= Ytop && target <= Ybottom ) && (p.CoordX >= Xleft && p.CoordX <= Xright)){
				if (p.IsDestinationOnBlockedArea(false, 0, 0)){
					p.TerminateAction()
				}
				return false	
			}
		}
		p.CoordY = target
		return true

		break

	default:
		break
	
	}

	return false
	
}

func (p *Unit) GetDistanceToTarget(coordType string) (distance int) {
	
	switch coordType {
	
	case "X":
		if (p.CoordX > p.DestinationX){
			distance = p.CoordX - p.DestinationX
		} else {
			distance = p.DestinationX - p.CoordX
		}
		break
	case "Y":
		
		if (p.CoordY > p.DestinationY){
			distance = p.CoordY - p.DestinationY
		} else {
			distance = p.DestinationY - p.CoordY
		}

		break
	}

	return distance
}


func (p *Unit) MoveLeft() (r bool, distanceToTarget int) {
	
	if (p.IsDestinationOnBlockedArea(true, p.CoordX - p.Speed - c.PlayerSpaceFromBlockSize, p.CoordY)){
		//p.TerminateAction()
		return false, 0
	}
	
	r = p.ChangeCoord("X", p.CoordX - p.Speed)

	// fmt.Println("MoveLeft")
	// fmt.Println(rand.Intn(10000000))
	// fmt.Println(p.IsDestinationOnBlockedArea(false, 0, 0))
	

	

	return r, p.GetDistanceToTarget("X")
} 

func (p *Unit) MoveRight() (r bool, distanceToTarget int) {
	
		
	if (p.IsDestinationOnBlockedArea(true, p.CoordX + p.Speed + c.PlayerSpaceFromBlockSize, p.CoordY)){
		// p.TerminateAction()
		return false, 0
	}
	
	r = p.ChangeCoord("X", p.CoordX + p.Speed)

	// fmt.Println("MoveRight")
	// fmt.Println(rand.Intn(10000000))
	// fmt.Println(p.IsDestinationOnBlockedArea(false, 0, 0))



	return r, p.GetDistanceToTarget("X")
} 

func (p *Unit) MoveDown() (r bool, distanceToTarget int) {


		
	if (p.IsDestinationOnBlockedArea(true, p.CoordX, p.CoordY + p.Speed + c.PlayerSpaceFromBlockSize)){
		p.TerminateAction()
		return false, 0
	}

	r = p.ChangeCoord("Y", p.CoordY + p.Speed)

	// fmt.Println("MoveDown")
	// fmt.Println(r)
	// fmt.Println(rand.Intn(10000000))
	// fmt.Println(p.IsDestinationOnBlockedArea(false, 0, 0))



	return r, p.GetDistanceToTarget("Y")
} 

func (p *Unit) MoveUp() (r bool, distanceToTarget int) {

	if (p.IsDestinationOnBlockedArea(true, p.CoordX, p.CoordY - p.Speed - c.PlayerSpaceFromBlockSize)){
		p.TerminateAction()
		return false, 0
	}

	
	r = p.ChangeCoord("Y", p.CoordY - p.Speed)

	// fmt.Println("MoveUp")
	// fmt.Println(rand.Intn(10000000))
	// fmt.Println(p.IsDestinationOnBlockedArea(false, 0, 0))


	return r, p.GetDistanceToTarget("Y")

} 

func (p *Unit) MoveTo(screen *ebiten.Image) bool {

	destinationX := "none"
	destinationY := "none"
	XDistance := 0
	YDistance := 0
	terminate := false
	var r bool


	if (p.DestinationX != p.CoordX || p.DestinationY != p.CoordY) {
		
		if (p.DestinationX != p.CoordX ){
			if (p.CoordX > p.DestinationX ){
				// Move Left
				r, XDistance = p.MoveLeft()
				if (!r){
					terminate = true
				}
				destinationX = "left"
			} else {
				// Move Right
				r, XDistance = p.MoveRight()
				if (!r){
					terminate = true
				}
				destinationX = "right"
			}
		}
	 	if (p.DestinationY != p.CoordY) {
			if (p.CoordY > p.DestinationY ) {
				// Move Down
				r, YDistance = p.MoveUp()
				destinationY = "up"
			} else {
				// Move Up
				r, YDistance = p.MoveDown()
				destinationY = "down"
			}
		}
		// fmt.Println(p.DestinationX)
		// fmt.Println(p.DestinationY)

	} else {
		p.TerminateAction()
	}
	
	if (terminate){
		p.TerminateAction()
	}
	

	p.DrawAnimation(destinationX, destinationY, XDistance, YDistance, screen)

	return true
}

func (p *Unit) DrawAnimation(destinationX string, destinationY string, XDistance int, YDistance int, screen *ebiten.Image) {
	
	var tileY int


	if (destinationY == "up"){
		tileY = 2
	} else if (destinationY == "down"){
		tileY = 0
	}
	if (destinationX == "left") {
		tileY = 3
	} else if (destinationX == "right") {
		tileY = 1
	}

	if ((XDistance + YDistance) > 0){
		if ((XDistance + YDistance ) % 5 == 0) {
			if (p.TileX < 3) {
				p.TileX++
			} else {
				p.TileX = 0
			}
		} 
	}
	//ebitenutil.DebugPrint(screen, fmt.Sprintf(strconv.Itoa(XDistance) + " : " + strconv.Itoa(YDistance)))
	//fmt.Println(moveCount)
	p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
}


func (u *Unit) GenerateRoute() {
	if (u.DestinationX > 0 || u.DestinationY > 0){
		u.Route = []Checkpoint{
			{1, 0},
			{2, 0},
		}
		
	}
}


func (u *Unit) ExecuteActionsUnit(screen *ebiten.Image) {
	
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
	
		u.DestinationX, u.DestinationY = ebiten.CursorPosition()
		u.MoveTo(screen)
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n x: " + strconv.Itoa(unit1.DestinationX) + "y: " +  strconv.Itoa(unit1.DestinationY) ))
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else if (u.DestinationX > 0 || u.DestinationY > 0){ 
		u.MoveTo(screen)
		
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else {
		u.DrawStatic(screen, 0, 0, float64(u.CoordX), float64(u.CoordY))
	}
	

}

func (u *Unit) ExecuteActionsPlayer1(screen *ebiten.Image) {
	
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	
		u.DestinationX, u.DestinationY = ebiten.CursorPosition()
		u.GenerateRoute()
		u.MoveTo(screen)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n x: " + strconv.Itoa(u.DestinationX) + "y: " +  strconv.Itoa(u.DestinationY) ))
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else if (u.DestinationX > 0 || u.DestinationY > 0){ 
		fmt.Println(u)
		u.MoveTo(screen)
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else {
		u.DrawStatic(screen, 0, 0, float64(u.CoordX), float64(u.CoordY))
	}
	

}