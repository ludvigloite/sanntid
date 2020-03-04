package elevController

import(
	"fmt"
	"../elevio"

)

//ha timer med her??

func Initialize(){
    elevio.Init("localhost:15657", orderHandler.numFloors)
    //Wipe alle ordre til nå??
 

    orderHandler.InitQueues(orderHandler.cabOrderQueue, orderHandler.hallOrderQueue)
	InitializeLights(orderHandler.numFloors)
	InitializeElevator()
	//Gjør det som main starter med.

}

func InitializeElevator(){
	hasHitFloor := false
	//kjør ned til etasjen under etasje
	//når treffer floor. Sett floor.

	drv_floors  := make(chan int)
    go elevio.PollFloorSensor(drv_floors)

	elevio.SetMotorDirection(elevio.MD_Down)

	a <- drv_floors

	elevio.SetMotorDirection(elevio.MD_Stop)
	orderHandler.setFloor(a)
	orderHandler.setDir(0)
	elevio.SetFloorIndicator(a)
	fmt.Println("Heisen er intialisert og venter i etasje nr ", a)

	close(drv_floors) //vet ikke om funker?


	/*for !hasHitFloor {
		select{ //skal runne helt til den treffer et floor
		case a := <- drv_floors:
			//fmt.Printf("%+v\n", a)
			hasHitFloor := true
			elevio.SetMotorDirection(elevio.MD_Stop)
			orderHandler.setFloor(a)
			orderHandler.setDir(0)
			elevio.SetFloorIndicator(a)
			fmt.Println("Heisen er intialisert og venter i etasje nr ", a)
		}
	}
	hasHitFloor := false //Resetter din til false for å unngå problemer. Kan kanskje fjerne?
	*/
}

func InitializeLights(numFloors int){ //NB: Endra her navn til numHallButtons
	//Slår av lyset på alle lys
	elevio.SetDoorOpenLamp(false)
	for i := 0; i < numFloors; i++{
		elevio.SetButtonLamp(2, i, false)
		if i != 0{ //er ikke i første etasje -> kan endre på alle ned_lys 
			elevio.SetButtonLamp(1,i,false)
		}
		if i != numFloors{ //er ikke i 4 etasje -> kan endre på alle opp_lys
			elevio.SetButtonLamp(0,i,false)
		}
	}

}

func OpenDoor(){

}
