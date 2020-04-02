package elevcontroller

import(
	"../elevio"
	//"../orderhandler"
	"../config"
	"fmt"
	"time"
	//"math/rand"
)
/*
func RankChecker2(elevID int, elevatorMap map[int]*config.Elevator){

	NuActiveElevators := 0
	for !elevatorMap[elevID].Active{}


	for _, elevator := range elevatorMap{
		if elevator.Active{
			NuActiveElevators++
		}
	}
	elevatorMap[elevID].ElevRank = NuActiveElevators

}

func RankChecker(elevatorMap map[int]*config.Elevator){ //kjøres som goroutine. Sjekker at det alltid er en og kun en master.
	ID_to_Master := -1
	nuMasters := 0
	for{
		nuMasters = 0
		for i, elevator := range elevatorMap{
			if elevator.ElevRank == 1{
				nuMasters++
				ID_to_Master = i
			}
		}
		if nuMasters == 0{
			for i, elevator := range elevatorMap{
				if elevator.Active{
					nuMasters++
					ID_to_Master = i
				}
			}
		}
		if nuMasters > 1{

		}
	}
}
*/

func Initialize(elevator *config.Elevator){
	elevio.SetMotorDirection(elevio.MD_Down)
	ResetLights()	
	InitQueues(elevator)
}

func InitQueues(elevator *config.Elevator){
	for i := 0;i < config.NUM_FLOORS; i++ {
		elevator.CabOrders[i] = false //INIT CABORDERS

		for j := elevio.BT_HallUp; j< config.NUM_HALLBUTTONS; j++{
			elevator.HallOrders[i][j] = false //INIT HALLORDERS
		}
	}
}

//Kan vel kanskje i stedet bare fjerne alle ordre og så kjøre update lights??
func ResetLights(){	//Slår av lyset på alle lys
	numFloors := config.NUM_FLOORS
	elevio.SetDoorOpenLamp(false)
	for i := 0; i < numFloors; i++{
		elevio.SetButtonLamp(elevio.BT_Cab, i, false)
		if i != 0{ //er ikke i første etasje -> kan endre på alle ned_lys 
			elevio.SetButtonLamp(elevio.BT_HallDown,i,false)
		}
		if i != numFloors{ //er ikke i 4 etasje -> kan endre på alle opp_lys
			elevio.SetButtonLamp(elevio.BT_HallUp,i,false)
		}
	}

}

func PrintElevator(elevator config.Elevator){
	fmt.Println()
	fmt.Println("elevID: ",elevator.ElevID,"\t Rank: ",elevator.ElevRank)
	fmt.Println("CurrentOrder = Floor: ",elevator.CurrentOrder.Floor, "\t ButtonType: ",elevator.CurrentOrder.ButtonType)
	fmt.Println("CurrentFloor = ", elevator.CurrentFloor)
	fmt.Println("CurrentState = ", elevator.CurrentState)
	fmt.Println("Hallorders   = ")
	for i := 0; i< config.NUM_FLOORS;i++{
		for j := elevio.BT_HallUp; j != elevio.BT_Cab; j++{
			fmt.Print(elevator.HallOrders[i][j],"\t")
		}
		fmt.Println()
	}
	fmt.Println()
}

func PrintElevators_withTime(elevatorMap map[int]*config.Elevator, openTime time.Duration){
	for{
		for _, elevator := range elevatorMap{
			fmt.Println()
			fmt.Println("elevID: ",elevator.ElevID,"\t Rank: ",elevator.ElevRank)
			fmt.Println("CurrentOrder = Floor: ",elevator.CurrentOrder.Floor, "\t ButtonType: ",elevator.CurrentOrder.ButtonType)
			fmt.Println("CurrentFloor = ", elevator.CurrentFloor)
			fmt.Println("CurrentState = ", elevator.CurrentState)
			fmt.Println("Hallorders   = ")
			for i := 0; i< config.NUM_FLOORS;i++{
				for j := elevio.BT_HallUp; j != elevio.BT_Cab; j++{
					fmt.Print(elevator.HallOrders[i][j],"\t")
				}
				fmt.Println()
			}
			fmt.Println()
			time.Sleep(300*time.Millisecond)
		}
		time.Sleep(openTime)
	}
}

func GetDirection(elevator config.Elevator) elevio.MotorDirection{
	currentFloor := elevator.CurrentFloor
	destinationFloor := elevator.CurrentOrder.Floor
	if destinationFloor == -1 || destinationFloor == currentFloor { //enten har den ikke noen retning, eller så er den på riktig floor
		//fmt.Print("\n\n\n\n SETTER DIRECTION TIL STOPP!! \n\n\n\n")
		return elevio.MD_Stop

	} else if currentFloor < destinationFloor { //heisen er lavere enn sin destinasjon -> kjører oppover
		return elevio.MD_Up

	} else{
		return elevio.MD_Down
	}
}

func ShouldStopAtFloor(elevator config.Elevator) bool{
	currentFloor := elevator.CurrentFloor
	dir := elevator.CurrentDir
	destinationFloor := elevator.CurrentOrder.Floor
	if currentFloor == destinationFloor{
		return true
	}
	if dir == elevio.MD_Stop{ //har ingen ordre eller er på etasjen currentOrder tilsier. KAN FØRE TIL ERROR!!
		return true
	}
	if elevator.CabOrders[currentFloor]{ //Det er en cab order i denne etasjen
		return true
	}
	if elevator.HallOrders[currentFloor][elevio.BT_HallUp] && dir == elevio.MD_Up{ //retning til heis er opp og det er en ordre opp
		return true
	}
	if elevator.HallOrders[currentFloor][elevio.BT_HallDown] && dir == elevio.MD_Down { //retning til heis er ned og det er en ordre ned
		return true
	}
	return false
}