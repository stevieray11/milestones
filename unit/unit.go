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
	"github.com/enjoykarma/milestones/src/images"
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


	fmt.Println("action Terminated")
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
				fmt.Println("CUSTOM DESTINATION ON BLOCKED!!!")
			} else {
				fmt.Println("PLAYER DESTINATION ON BLOCKED!!!")
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


func getNextCheckpointRecursive(routesIn []Checkpoint, start m.Area, u *Unit) ([]Checkpoint) {


	if (start.IsBlockedArea() ){
		return routesIn
	}


	left := start.LeftArea()
	right := start.RightArea()
	up := start.UpArea()
	down := start.DownArea()

	var dist float64
	var minDist float64
	var startDist float64
	var nextArea m.Area

	startX, startY := start.GetCenterPointFloat64()

	fmt.Println("start-", start, " R-", right, " L-", left, " D-", down, " U-", up )
	startDist = start.GetDistanceTo(float64(u.DestinationX), float64(u.DestinationY))	
	minDist = right.GetDistanceTo(float64(u.DestinationX), float64(u.DestinationY))	
	fmt.Println("dist START ", startDist)
	fmt.Println("dist R ", minDist)
	if (right.IsCloserToDestination(startX, startY, float64(u.DestinationX), float64(u.DestinationY)) && !right.IsBlockedArea() ) {
		//x, y := right.GetCenterPoint()
		//routesIn = append(routesIn, Checkpoint{int(x), int(y)})
		nextArea = right
		fmt.Println("possible right")
	}


	dist = up.GetDistanceTo(float64(u.DestinationX), float64(u.DestinationY))
	fmt.Println("dist U", dist)
	if (up.IsCloserToDestination(startX, startY, float64(u.DestinationX), float64(u.DestinationY)) && !up.IsBlockedArea() ) {
		if (dist < minDist && dist > 0) {
			minDist = dist
			nextArea = up
			fmt.Println("possible up")
		}
	}

	dist = down.GetDistanceTo(float64(u.DestinationX), float64(u.DestinationY))
	fmt.Println("dist D", dist)
	if (down.IsCloserToDestination(startX, startY, float64(u.DestinationX), float64(u.DestinationY)) && !down.IsBlockedArea() ) {			
		if (dist < minDist && dist > 0) {
			minDist = dist
			nextArea = down
			fmt.Println("possible down")
		}
	}

	dist = left.GetDistanceTo(float64(u.DestinationX), float64(u.DestinationY))
	fmt.Println("dist L", dist)
	if (left.IsCloserToDestination(startX, startY, float64(u.DestinationX), float64(u.DestinationY)) && !left.IsBlockedArea() ) {
		if (dist < minDist && dist > 0) {
			nextArea = left
			fmt.Println("possible left")
		}
	}
	



	if (nextArea.IsCoordsInArea(u.DestinationX, u.DestinationY)) {
		fmt.Println("finish")
		return routesIn
	}

	fmt.Println("min Dist", minDist)
	// fmt.Println(nextArea)
	if (minDist > 0 && nextArea != m.Area{}){
		x, y := nextArea.GetCenterPoint()
		routesIn = append(routesIn, Checkpoint{int(x), int(y)})
		routesIn = getNextCheckpointRecursive(routesIn, nextArea, u)
	}
	

	return routesIn
}

func (u *Unit) GenerateRoute(screen *ebiten.Image) bool {
	if (u.DestinationX > 0 || u.DestinationY > 0){
		var PointA = Checkpoint{u.CoordX, u.CoordY}
		//var PointB = Checkpoint{}

		if (u.IsDestinationOnBlockedArea(false, 0, 0)){
			u.TerminateAction()
			return false
		}

		u.Route = []Checkpoint{PointA}
		u.Route = getNextCheckpointRecursive(u.Route, m.GetAreaByCoords(u.CoordX, u.CoordY), u)


	}

	return false

}



func (u *Unit) ExecuteActionsUnit(screen *ebiten.Image) {
	
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
	

		
		
		X, Y := ebiten.CursorPosition()
		a := m.GetAreaByCoords(X, Y)

		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n type: " + strconv.Itoa(a.Type) + "x: " +  strconv.Itoa(int(a.X)) + "y:" + strconv.Itoa(int(a.Y)) ))
		//u.MoveTo(screen)
		
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else if (u.DestinationX > 0 || u.DestinationY > 0){ 
		//u.MoveTo(screen)
		
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else {
		//u.DrawStatic(screen, 0, 0, float64(u.CoordX), float64(u.CoordY))
	}
	

}

func (u *Unit) ExecuteActionsPlayer1(screen *ebiten.Image) {
	
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
	
		u.DestinationX, u.DestinationY = ebiten.CursorPosition()
		u.GenerateRoute(screen)
		fmt.Println(u.Route)
		u.MoveTo(screen)
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n x: " + strconv.Itoa(u.DestinationX) + "y: " +  strconv.Itoa(u.DestinationY) ))
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else if (u.DestinationX > 0 || u.DestinationY > 0){ 
		fmt.Println(u.Route)
		
		u.MoveTo(screen)
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else {
		u.DrawStatic(screen, 0, 0, float64(u.CoordX), float64(u.CoordY))
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		u.TerminateAction()
	}

	u.DrawRoute(screen)

}


func (u *Unit) DrawRoute(screen *ebiten.Image) {
	
	for _, cp := range u.Route {
		for _, a := range m.GameMapLayer2{
			if (a.IsBlockedArea()) {
				continue
			}
			if (a.IsCoordsInArea(cp.X, cp.Y) ) {
				a.Mark(screen)
			}
		}
	}

}