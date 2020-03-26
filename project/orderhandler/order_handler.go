///////////////////////////////////////////////////
//	Viktig å huske på at ned første etasje og opp øverste etasje ikke finnes!
//	Selve order_handler blir vel ikke kjørt. Vil verdiene bli intialisert?
//
//	Kanskje bare la alt skje i orderhandler. Eneste som blir sendt mellom FSM og orderhandler er ordre som skal legges til. Dette gjøres med goroutine
//
//////////////////////////////////////////////////


package orderhandler

import(
	"../elevio"
	"../config"
	"fmt"
)


var elevatorID int //kan byttes underveis
var elevatorRank int //bytter underveis
var isMaster = false

var currentOrder int //sier hvilken etasje heisen er på vei til. -1 om den ikke har noen ordre.
var currentFloor int //hvilken etasje er heisen i nå. 0 , 1 , 2 , 3
var currentDir int //hvilken retning har heisen. -1 , 0 , 1. Kun 0 i spesielle tilfeller. Er -1 / 1 også når den stopper i et floor. Den skal jo tross alt videre i samme retning.
var currentState int //0: IDLE, 1: ACTIVE, 2: DOOR_OPEN, 3:UNDEFINED

var hallOrderQueue = &[config.NUM_FLOORS][config.NUM_HALLBUTTONS] int{} //inneholder en liste med alle hall orders. -1 om inaktiv. 0 om den er aktiv, men ikke tatt. ellers ID til heisen om en av dem skal utføre ordren.
//nullte element er opp, første element er ned.



type CabOrders struct{
	ElevID int //hvilken elevator cab callsa tilhører
	Active [config.NUM_FLOORS]int //hvilke av de fire knappene som er aktive // !!!! -1:inaktiv, 0:aktiv 2:executing?
}
type Order struct{
	Floor int
	ButtonType int
}

var cabOrderQueue = &CabOrders{}//variabelen som kan endres på

var currentOrderList = &[config.NUM_ELEVATORS] Order{}
var currentFloorList = &[config.NUM_ELEVATORS] int{}

func IsMaster()bool{
	if elevatorRank == 1{
		return true
	}
	return false
}

func UpdateCurrentOrderList(elevID int, current_order Order){ //evt bytte ut int med en Ordre? Tror egt vi bare trenger FloorNr, men er kanskje mer leselig om man tar med hele order.
	*currentOrderList[elevID] = current_order
}
func UpdateCurrentFloorList(elevID int, current_floor int){ //evt bytte ut int med en Ordre? Tror egt vi bare trenger FloorNr, men er kanskje mer leselig om man tar med hele order.
	*currentOrderList[elevID] = current_floor
}


func SetElevatorRank(rank int){elevatorRank = rank}
func SetElevatorID(ID int){elevatorID = ID}
func SetCurrentFloor(floor int){currentFloor = floor}
func SetCurrentDir(dir int){currentDir = dir}
func SetCurrentOrder(floor int){currentOrder = floor}
func SetHallOrderQueue(queue [config.NUM_FLOORS][config.NUM_HALLBUTTONS] int){ *hallOrderQueue = queue}
func SetCurrentState(state int){currentState = state}

func GetCurrentOrder()int {return currentOrder}
func GetCurrentDir()int {return currentDir}
func GetCurrentFloor()int {return currentFloor}
func GetElevID()int {return elevatorID}
func GetElevRank()int {return elevatorRank}
func GetHallOrderQueue()[config.NUM_FLOORS][config.NUM_HALLBUTTONS] int{return *hallOrderQueue}
func GetCurrentState()int{return currentState}
func GetCurrentOrderList()[config.NUM_ELEVATORS] Order{return *currentOrderList}
func GetCurrentFloorList()[config.NUM_ELEVATORS] int{return *currentFloorList}




func InitQueues(){
	InitCabQueue(cabOrderQueue)
	InitHallQueue(hallOrderQueue)
	InitCurrentOrderList(currentOrderList)
}

func InitHallQueue(queue *[config.NUM_FLOORS][config.NUM_HALLBUTTONS] int){
	for i := 0; i < config.NUM_FLOORS; i++{
		for j := 0; j < config.NUM_HALLBUTTONS; j++{
			queue[i][j] = -1
		}
	}
}

func InitCabQueue(queue *CabOrders){
	queue.ElevID = elevatorID
	fmt.Println("->>",queue.ElevID)
	for i := 0; i < config.NUM_FLOORS; i++{
		queue.Active[i] = -1
	}
}

func InitCurrentOrderList(list *[config.NUM_ELEVATORS]int){
	for i:=0;i<config.NUM_ELEVATORS;i++{
		list[i].Floor = -1
		list[i].ButtonType = -1
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



func GetNewOrder(elevCurrentFloor int, elevID int) Order{ //returnerer en ordre med floor: -1 om det ikke er noen ordre.
	newOrder := Order{}
	if IsThereOrder(elevCurrentFloor,0,elevID){
		newOrder.Floor = elevCurrentFloor
		newOrder.ButtonType = 0
		return newOrder
	}else if IsThereOrder(elevCurrentFloor,1,elevID){
		newOrder.Floor = elevCurrentFloor
		newOrder.ButtonType = 1
		return newOrder
	}

	for i := 0; i < config.NUM_FLOORS; i++{
		if IsThereOrder(i, 2, elevID){
			//det finnes en cab order
			newOrder.Floor = i 
			newOrder.ButtonType = 2
			return newOrder
		}
		for j := 0; j < config.NUM_HALLBUTTONS; j++{
			if IsThereOrder(i, j, elevID){
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

func AddOrder(floor int, buttonType int, elevatorID int){ //elevatorID er 0 om det bare skal legges inn ordre uten at noen tar den.
	if buttonType == 2{ //caborder
		cabOrderQueue.Active[floor] = elevatorID //active
	} else{
		hallOrderQueue[floor][buttonType] = elevatorID //active
	}
	//UpdateLights()
}



func IsThereOrder(floor int, buttonType int, elevID int) bool{ //buttontype: 0=opp 1=ned 2=cabOrder //kan kanskje bare implementeres i ShouldStopAtFloor //kan kanskje fjerne elevID og heller bruke global variabel
	if buttonType == 2{
		//fmt.Println(elevID," ... ", cabOrderQueue.ElevID)
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

func LightUpdater(LightUpdateCh <-chan bool){
	for{
		select{
		case <-LightUpdateCh:
			fmt.Println("Updating Lights...")
			for i := 0; i < config.NUM_FLOORS; i++{
				if cabOrderQueue.Active[i] ==-1 {
					elevio.SetButtonLamp(elevio.BT_Cab, i, false)
				} else{
					elevio.SetButtonLamp(elevio.BT_Cab, i, true)
				}

				for j := 0; j < config.NUM_HALLBUTTONS; j++{
					if i != 0 && j == 1{ //hvis det ikke er 1 etasje eller 4 etasje.
						if hallOrderQueue[i][j] == -1{
							elevio.SetButtonLamp(elevio.BT_HallDown, i, false)
						} else{
							elevio.SetButtonLamp(elevio.BT_HallDown, i, true)
						}
					}
					if i != config.NUM_FLOORS && j == 0{
						if hallOrderQueue[i][j] == -1{
							elevio.SetButtonLamp(elevio.BT_HallUp, i, false)
						} else{
							elevio.SetButtonLamp(elevio.BT_HallUp, i, true)
						}
					}
				}
			}
		}
	}
}

/*
func UpdateLights(){ //vet ikke om i og j blir riktig???? //Kan sikkert gjøres mer effektiv. NumHallButtons er jo bare 2..Evt lage en funskjon for hall-lights og en for cab-lights
	for i := 0; i < config.NUM_FLOORS; i++{
		if cabOrderQueue.Active[i] ==-1 {
			elevio.SetButtonLamp(elevio.BT_Cab, i, false)
		} else{
			elevio.SetButtonLamp(elevio.BT_Cab, i, true)
		}

		for j := 0; j < config.NUM_HALLBUTTONS; j++{
			if i != 0 && j == 1{ //hvis det ikke er 1 etasje eller 4 etasje.
				if hallOrderQueue[i][j] == -1{
				elevio.SetButtonLamp(elevio.BT_HallDown, i, false)
				} else{
				elevio.SetButtonLamp(elevio.BT_HallDown, i, true)
				}
			}
			if i != config.NUM_FLOORS && j == 0{
				if hallOrderQueue[i][j] == -1{
				elevio.SetButtonLamp(elevio.BT_HallUp, i, false)
				} else{
				elevio.SetButtonLamp(elevio.BT_HallUp, i, true)
				}
			}
		}
	}
}
*/


func MergeHallQueues(elev2 config.Packet){

	for i := 0; i < config.NUM_FLOORS; i++{
	
		for j := 0; j < config.NUM_HALLBUTTONS; j++{
			if i != 0 && j == 1{ //hvis det ikke er 1 etasje eller 4 etasje.
				hallOrderQueue[i][j] = PrioritizeNumbers(elev2, i, j)
			}
			if i != config.NUM_FLOORS && j == 0{
				hallOrderQueue[i][j] = PrioritizeNumbers(elev2, i, j)
			}
		}
	}
}

func PrioritizeNumbers(elev2 config.Packet, i int, j int) int{
	order1 := hallOrderQueue[i][j]
	order2 := elev2.Order_list[i][j]

	if order1 == order2{
		return order1
	}
	if (order1 == -1 && currentFloor==i && currentState == 2) || (order2 == -1 && elev2.CurrentFloor==i && elev2.State == 2){
		return -1
	}

	if order1 == elevatorID {
		return order1
	}else if order2 == elev2.ID{
		return order2
	}
	return 0
}

func PrintHallOrderQueue(hallOrderQueue [config.NUM_FLOORS][config.NUM_HALLBUTTONS]int){
	for i := 0; i < config.NUM_FLOORS; i++{
		fmt.Print(i+1," etasje: \t")
		for j := 0; j < config.NUM_HALLBUTTONS; j++{
			if i != 0 && j == 1{ //hvis det ikke er 1 etasje eller 4 etasje.
				fmt.Print(hallOrderQueue[i][j], "\t")
			}
			if i != config.NUM_FLOORS && j == 0{
				fmt.Print(hallOrderQueue[i][j], "\t")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}




