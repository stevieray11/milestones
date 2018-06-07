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
				fmt.Println("CUSTOM DESTINATION ON BLOCKED!!!")
			} else {
				fmt.Println("PLAYER DESTINATION ON BLOCKED!!!")
			}
			

			return true
		}
		
	}
	return false
	
}

func (p *Unit) ChangeCoord(coordType string, target int) (r bool, moveType string, strafeChangedCoord int){
	
	var (
		Xleft int
		Xright int
		Ybottom int
		Ytop int
		changedCoord int
	)


	strafeChangedCoord = 0



	switch coordType {
	case "X" :

		for _, a := range m.GameMapLayer2{
			Xleft = int(a.X) - int(c.TileSize) - c.PlayerSpaceFromBlockSize
			Xright = int(a.X) + c.PlayerSpaceFromBlockSize
		
			Ytop = int(a.Y) - c.PlayerSpaceFromBlockSize
			Ybottom = int(a.Y) + int(c.TileSize) + c.PlayerSpaceFromBlockSize
			if (a.IsBlockedArea() && (target >= Xleft && target <= Xright) && (p.CoordY >= Ytop && p.CoordY <= Ybottom) ) {
				
				if (p.IsDestinationOnBlockedArea(false, 0, 0)){

					// fmt.Println("Blocked X")
					// fmt.Println(target)
					// fmt.Println(a)
					// fmt.Println(Xleft)
					// fmt.Println(Xright)
					// fmt.Println(rand.Intn(100000))
					p.TerminateAction()
					return false, "", strafeChangedCoord

				} else {
					if (p.DestinationY > p.CoordY) {
						changedCoord = p.CoordY + p.Speed + c.PlayerSpaceFromBlockSize
					} else if (p.DestinationY == p.CoordY) {
						
						return false, "", strafeChangedCoord
						//changedCoord = p.CoordY + p.Speed + c.PlayerSpaceFromBlockSize
					} else {
						changedCoord = p.CoordY - p.Speed - c.PlayerSpaceFromBlockSize
					}

					if (!p.IsDestinationOnBlockedArea(true, p.CoordX, changedCoord)){
					
						if (p.DestinationY > p.CoordY) {
							changedCoord = p.CoordY + p.Speed
						} else if (p.DestinationY == p.CoordY) {
							return false, "", strafeChangedCoord
							//changedCoord = p.CoordY - p.Speed
						} else {
							changedCoord = p.CoordY - p.Speed
						}
						strafeChangedCoord = changedCoord
						//p.CoordY = changedCoord
						//p.ChangeCoord("Y", changedCoord)
						return true, "strafeY", strafeChangedCoord
					}
				
					return false, "", strafeChangedCoord
				}
			}
		}
		
		p.CoordX = target

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
					return false, "", strafeChangedCoord

				} else {

					if (p.DestinationX > p.CoordX) {
						changedCoord = p.CoordX + p.Speed + c.PlayerSpaceFromBlockSize
					} else if (p.DestinationX == p.CoordX){
						return false, "", strafeChangedCoord
						//changedCoord = p.CoordX - p.Speed - c.PlayerSpaceFromBlockSize
					} else {
						changedCoord = p.CoordX - p.Speed - c.PlayerSpaceFromBlockSize
					}

					if (!p.IsDestinationOnBlockedArea(true, changedCoord, p.CoordY)){
						//p.ChangeCoord("X", changedCoord)
					
						if (p.DestinationX > p.CoordX) {
							changedCoord = p.CoordX + p.Speed
						} else if (p.DestinationX == p.CoordX){
							//changedCoord = p.CoordX - p.Speed
							return false, "", strafeChangedCoord
						} else {
							changedCoord = p.CoordX - p.Speed 
						}
						
						strafeChangedCoord = changedCoord
						//p.CoordX = changedCoord
						//p.ChangeCoord("X", changedCoord)
						return true, "strafeX", strafeChangedCoord
					}
					
					return false, "", strafeChangedCoord

				}

			}
		}
		p.CoordY = target

		break

	default:
		break
		
	}

	return true, "", strafeChangedCoord
}

func (p *Unit) MoveTo(screen *ebiten.Image) bool {

	destinationX := "none"
	destinationY := "none"

	//var moved bool = false
	var (
		r bool
		strafeChoord int
		moveType string
		moveCount int = 0
	)
	tileY := 0

	XDistance := 0
	YDistance := 0

	if (p.DestinationX != p.CoordX || p.DestinationY != p.CoordY) {
		
		if (p.DestinationX != p.CoordX ){
			
			if (p.CoordX > p.DestinationX ){
				if (p.Speed > p.CoordX) {
					p.CoordX = p.DestinationX
					XDistance = 0
				} else {
					r, moveType, strafeChoord = p.ChangeCoord("X", p.CoordX - p.Speed)
					if (!r) {
						p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
						return false
					} else {
						
						if (moveType == "strafeY") {
							fmt.Println("strafeY")
							r, _, _ = p.ChangeCoord("Y", strafeChoord)
							if (!r) {
								p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
								return false
							} else {
								moveCount++
							}
						} else {
							destinationX = "left"
							XDistance = p.CoordX - p.DestinationX
							moveCount++
						}//moved = true
					}
					
				}
				
			} else {
				r, moveType, strafeChoord = p.ChangeCoord("X", p.CoordX + p.Speed)
				if (!r) {
					p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
					return false
				} else {
					if (moveType == "strafeY") {
						fmt.Println("strafeY")
						r, _, _ = p.ChangeCoord("Y", strafeChoord)
						if (!r) {
							p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
							return false
						} else {
							moveCount++
						}
					} else {
						moveCount++
						destinationX = "right"
						XDistance = p.DestinationX - p.CoordX 
					}//moved = true
				}
				
				
			}
		}
		if (p.DestinationY != p.CoordY) {
			if (p.CoordY > p.DestinationY ){
				if (p.Speed > p.CoordY) {
					p.CoordY = p.DestinationY
					YDistance = 0
				} else {
					r, moveType, strafeChoord = p.ChangeCoord("Y", p.CoordY - p.Speed)
					if (!r) {
						p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
						return false
					} else {
						if (moveType == "strafeX") {
							fmt.Println("strafeX")
							r, _, _ = p.ChangeCoord("X", strafeChoord)
							if (!r) {
								p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
								return false
							} else {
								moveCount++
							}
						} else {
							moveCount++
							destinationY = "down"
							YDistance = p.CoordY - p.DestinationY
						}

					}
			}

			} else {
				
				r, moveType, strafeChoord = p.ChangeCoord("Y", p.CoordY + p.Speed)
				if (!r) {
					p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
					return false
				} else {
					if (moveType == "strafeX") {
						fmt.Println("strafeX")
						r, _, _ = p.ChangeCoord("X", strafeChoord)
						if (!r) {
							p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))
							return false
						} else {
							moveCount++
						}
					} else {
						moveCount++
						destinationY = "up"
						YDistance = p.DestinationY - p.CoordY
					}
					
				}

				
			}
		}
		
	} else {
		p.TerminateAction()
		//p.DrawStatic(screen, 0, 0, float64(p.CoordX), float64(p.CoordY))
	}


	if (destinationY == "up"){
		tileY = 0
	} else if (destinationY == "down"){
		tileY = 2
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
	fmt.Println(moveCount)
	p.DrawStatic(screen, p.TileX, tileY, float64(p.CoordX), float64(p.CoordY))

	return true
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
		u.MoveTo(screen)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n x: " + strconv.Itoa(u.DestinationX) + "y: " +  strconv.Itoa(u.DestinationY) ))
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else if (u.DestinationX > 0 || u.DestinationY > 0){ 
		u.MoveTo(screen)
		//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	} else {
		u.DrawStatic(screen, 0, 0, float64(u.CoordX), float64(u.CoordY))
	}
	

}