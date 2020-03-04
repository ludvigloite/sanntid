///////////////////////////////////////////////////
//	Viktig å huske på at ned første etasje og opp øverste etasje ikke finnes!
//	Selve order_handler blir vel ikke kjørt. Vil verdiene bli intialisert?
//
//
//
//////////////////////////////////////////////////

package order_handler

import{
	"fmt"
}


const numFloors = 4
const numHallButtons = 2

var elevatorID int //kan byttes underveis
var elevatorRank int //bytter underveis
var isMaster = false

type ElevState int
const(
	Idle = 0
	Executing = 1
	Lost = 2
)

var currentOrder int //sier hvilken etasje heisen er på vei til. -1 om den ikke har noen ordre.
var currentFloor int //hvilken etasje er heisen i nå. 0 , 1 , 2 , 3
var currentDir int //hvilken retning har heisen. -1 , 0 , 1. Kun 0 i spesielle tilfeller. Er -1 / 1 også når den stopper i et floor. Den skal jo tross alt videre i samme retning.

type CabOrders struct{
	ElevID int //hvilken elevator cab callsa tilhører
	Active [numFloors]int //hvilke av de fire knappene som er aktive // !!!! -1:inaktiv, 0:aktiv 2:executing?
}

var cabOrderQueue = &CabOrders{} struct //variabelen som kan endres på

var hallOrderQueue = &[numFloors][numHallButtons] int //inneholder en liste med alle hall orders. -1 om inaktiv. 0 om den er aktiv, men ikke tatt. ellers ID til heisen om en av dem skal utføre ordren.
//nullte element er opp, første element er ned.


func InitQueues(hallQueue *[numFloors][numHallButtons] int, cabQueue *CabOrders){
	InitHallQueue(hallQueue)
	InitCabQueue(cabQueue)
}


func InitHallQueue(queue *[numFloors][numHallButtons] int){
	for i := 0; i < numFloors; i++{
		for j := 0; j < numHallButtons; j++{
			queue[i][j] = -1
		}
	}
}

func InitCabQueue(queue *CabOrders){
	queue.ElevID = elevatorID
	for i := 0; i < numFloors; i++{
		queue.Active[i] = -1
	}
}

func SetCurrentFloor(int floor){
	currentFloor = floor
}
func SetCurrentDir(int dir){
	currentDir = dir
}

////////////// ARBITRATOR UNDER ? //////////////

func WhatElevatorShouldTakeOrder(){ //Evt whatOrderSHouldthisElevatorTake()??

}





////////////// ARBITRATOR OVER ? //////////////



func IsThereOrder(floor int, buttonType int, elevID int) bool{ //buttontype: 0=opp 1=ned 2=cabOrder //kan kanskje bare implementeres i ShouldStopAtFloor //kan kanskje fjerne elevID og heller bruke global variabel
	if buttontype == 2{
		if cabOrderQueue.Active[floor] == 0 && cabOrderQueue.ElevID == elevID{
			return true
		}
	} else{
		if hallOrderQueue[floor][buttonType] == 0 || hallOrderQueue[floor][buttonType] == elevID{
			return true
		}
	}
	return false
}



func ShouldStopAtFloor(currentFloor int, currentOrder int, elevID int) bool{
	dir := getDirection(currentFloor,currentOrder) //-1, 0 eller 1
	if dir == 0{ //har ingen ordre eller er på etasjen currentOrder tilsier
		return true
	}
	if IsThereOrder(currentOrder,2,elevID){ //Det er en cab order i denne etasjen
		return true
	}
	if IsThereOrder(currentOrder,0,elevID) && dir == 1{ //retning til heis er opp og det er en ordre opp
		return true
	}
	if IsThereOrder(currentOrder,1,elevID) && dir == -1 { //retning til heis er ned og det er en ordre ned
		return true
	}
	return false
}

func ClearFloor(floor int) void{ //fjerner alle ordre i denne etasjen fra køene. Kan bare utføres av heisen selv
	hallOrderQueue[floor][0] = -1
	hallOrderQueue[floor][1] = -1
	cabOrderQueue.Active[floor] = -1
}

func AddOrder(floor int, buttonType int) void{
	if buttonType == 2{ //caborder
		cabOrderQueue.Active[floor] = 0 //active
	} else{
		hallOrderQueue[floor][buttonType] = 0 //active
	}
}

func UpdateLights() void{ //vet ikke om i og j blir riktig???? //Kan sikkert gjøres mer effektiv. NumHallButtons er jo bare 2..Evt lage en funskjon for hall-lights og en for cab-lights
	for i := 0; i < numFloors; i++{
		if cabOrderQueue.Active[i] ==-1 {
			elevio.SetButtonLamp(2, i, false)
		} else{
			elevio.SetButtonLamp(2, i, true)
		}

		for j := 0; j < numHallButtons; j++{
			if i != 0 && i != 3{ //hvis det ikke er 1 etasje eller 4 etasje.
				if hallOrderQueue[i][j] == -1{
				elevio.SetButtonLamp(j, i, false)
				} else{
				elevio.SetButtonLamp(j, i, true)
				}
			}
		}
	}
	if hallOrderQueue[0][0] == -1{ //første etasje opp
	elevio.SetButtonLamp(0, 0, false)
	} else{
		elevio.SetButtonLamp(0, 0, true)
	}

	if hallOrderQueue[numFloors][1] == -1{ //4 etasje ned
	elevio.SetButtonLamp(1, numFloors, false)
	} else{
		elevio.SetButtonLamp(1, numFloors, true)
	}
}

func GetDirection(currentFloor int, currentOrder int) int{
	if currentOrder == -1 || currentOrder == currentFloor { //enten har den ikke noen retning, eller så er den på riktig floor
		return 0

	} else if currentFloor < currentOrder { //heisen er lavere enn sin destinasjon -> kjører oppover
		return 1

	} else{
		return -1
	}

}
