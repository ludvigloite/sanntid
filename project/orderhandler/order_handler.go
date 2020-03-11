///////////////////////////////////////////////////
//	Viktig å huske på at ned første etasje og opp øverste etasje ikke finnes!
//	Selve order_handler blir vel ikke kjørt. Vil verdiene bli intialisert?
//
//	Kanskje bare la alt skje i orderhandler. Eneste som blir sendt mellom FSM og orderhandler er ordre som skal legges til. Dette gjøres med goroutine
//
//////////////////////////////////////////////////


package orderhandler

import(
	"fmt"
	"../elevio"
)

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
type Order struct{
	Floor int
	ButtonType int
}


var cabOrderQueue = &CabOrders{}//variabelen som kan endres på

var hallOrderQueue = &[numFloors][numHallButtons] int{} //inneholder en liste med alle hall orders. -1 om inaktiv. 0 om den er aktiv, men ikke tatt. ellers ID til heisen om en av dem skal utføre ordren.
//nullte element er opp, første element er ned.


func InitQueues(){
	InitCabQueue(cabOrderQueue)
	InitHallQueue(hallOrderQueue)
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

func SetElevatorID(ID int){
	elevatorID = ID
}

func SetCurrentFloor(floor int){
	currentFloor = floor
}

func SetCurrentDir(dir int){
	currentDir = dir
}

func SetCurrentOrder(floor int){
	currentOrder = floor
}
func GetCurrentOrder()int {return currentOrder}
func GetCurrentDir()int {return currentDir}
func GetCurrentFloor()int {return currentFloor}
func GetElevatorID()int {return elevatorID}
func GetNumFloors()int {return numFloors}
func GetNumHallButtons()int {return numHallButtons}

////////////// ARBITRATOR UNDER ? //////////////

func WhatElevatorShouldTakeOrder(){ //Evt whatOrderSHouldthisElevatorTake()??

}





////////////// ARBITRATOR OVER ? //////////////



func IsThereOrder(floor int, buttonType int, elevID int) bool{ //buttontype: 0=opp 1=ned 2=cabOrder //kan kanskje bare implementeres i ShouldStopAtFloor //kan kanskje fjerne elevID og heller bruke global variabel
	if buttonType == 2{
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
	dir := GetDirection(currentFloor,currentOrder) //-1, 0 eller 1
	if dir == 0{ //har ingen ordre eller er på etasjen currentOrder tilsier
		return true
	}
	if IsThereOrder(currentFloor,2,elevID){ //Det er en cab order i denne etasjen
		return true
	}
	if IsThereOrder(currentFloor,0,elevID) && dir == 1{ //retning til heis er opp og det er en ordre opp
		return true
	}
	if IsThereOrder(currentFloor,1,elevID) && dir == -1 { //retning til heis er ned og det er en ordre ned
		return true
	}
	return false
}

func ClearFloor(floor int){ //fjerner alle ordre i denne etasjen fra køene. Kan bare utføres av heisen selv
	//gjør det noe at den setter -1 til 1 etasje ned og 4 etasje opp??
	hallOrderQueue[floor][0] = -1
	hallOrderQueue[floor][1] = -1
	cabOrderQueue.Active[floor] = -1
}

func AddOrder(floor int, buttonType int, elevatorID int){ //elevatorID er 0 om det bare skal legges inn ordre uten at noen tar den.
	if buttonType == 2{ //caborder
		cabOrderQueue.Active[floor] = elevatorID //active
	} else{
		hallOrderQueue[floor][buttonType] = elevatorID //active
	}
	fmt.Println(cabOrderQueue.Active)
	fmt.Println(hallOrderQueue)
}

func UpdateLights(){ //vet ikke om i og j blir riktig???? //Kan sikkert gjøres mer effektiv. NumHallButtons er jo bare 2..Evt lage en funskjon for hall-lights og en for cab-lights
	for i := 0; i < numFloors; i++{
		if cabOrderQueue.Active[i] ==-1 {
			elevio.SetButtonLamp(elevio.BT_Cab, i, false)
		} else{
			elevio.SetButtonLamp(elevio.BT_Cab, i, true)
		}

		for j := 0; j < numHallButtons; j++{
			if i != 0 && j == 1{ //hvis det ikke er 1 etasje eller 4 etasje.
				if hallOrderQueue[i][j] == -1{
				elevio.SetButtonLamp(elevio.BT_HallDown, i, false)
				} else{
				elevio.SetButtonLamp(elevio.BT_HallDown, i, true)
				}
			}
			if i != numFloors && j == 0{
				if hallOrderQueue[i][j] == -1{
				elevio.SetButtonLamp(elevio.BT_HallUp, i, false)
				} else{
				elevio.SetButtonLamp(elevio.BT_HallUp, i, true)
				}
			}
		}
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


func GetNewOrder() Order{ //returnerer en ordre med floor: -1 om det ikke er noen ordre.
	newOrder := Order{}
	for i := 0; i < numFloors; i++{
		if IsThereOrder(i, 2, elevatorID){
			//det finnes en cab order
			newOrder.Floor = i 
			newOrder.ButtonType = 2
			return newOrder
		}
		for j := 0; j < numHallButtons; j++{
			if IsThereOrder(i, j, elevatorID){
				newOrder.Floor = i
				newOrder.ButtonType = j
				return newOrder
			}
			
		}
	}
	newOrder.Floor = -1
	newOrder.ButtonType = -1
	return newOrder
}

