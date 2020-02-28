package order_handler

import{
	"fmt"
}


const numFloors = 4
const numButtons = 3

var elevatorID int //kan byttes underveis
var elevatorRank int //bytter underveis
var isMaster = false

type ElevState int
const(
	Idle = 0
	Executing = 1
	Lost = 2
)

var currentOrder int

type CabOrderQueue struct{
	ID int //hvilken elevator cab callsa tilhører
	Active [numFloors]int //hvilke av de fire knappene som er aktive
}

var hallOrderQueue = &[numFloors][numButtons] int //inneholder en liste med alle hall orders. -1 om utilegnet. ID til heisen om en av dem skal utføre ordren



func initQueue(queue *){ //må fikses
	for i := 0
	//INIT alle til 0

}

func shouldStopAtFloor(currentFloor int, currentOrder int) bool{
	getDirection(currentFloor,currentOrder)


}

func getDirection(currentFloor int, currentOrder int) int{


}