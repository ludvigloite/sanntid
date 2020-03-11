package elevcontroller

import(
	//"fmt"
	"time"
	"../elevio"
	"../orderhandler"
	"../config"
	//"../timer"

)

//ha timer med her??

func Initialize(){
    elevio.Init("localhost:15657", config.NUM_FLOORS)
    orderhandler.SetElevatorID(1) //BØR IKKE HARDKODES!!
    //Wipe alle ordre til nå??

    orderhandler.InitQueues()
	InitializeLights(orderhandler.GetNumFloors())
	
	//InitializeElevator()
	//Gjør det som main starter med.

}

func checkAndAddOrder(Drv_buttons <- chan elevio.ButtonEvent){
	for{
		select{
			case order := <- Drv_buttons: //Fått inn knappetrykk
				fmt.Println("Knapp er trykket ", int(order.Button), order.Floor)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				orderhandler.UpdateLights()
	}
}

/*
func InitializeElevator(){
	//hasHitFloor := false
	//kjør ned til etasjen under etasje
	//når treffer floor. Sett floor.

	drv_floors  := make(chan int)
	fmt.Println("Prøver å starte goroutine")
    go elevio.PollFloorSensor(drv_floors)
    fmt.Println("ferdig med goroutine")

	elevio.SetMotorDirection(elevio.MD_Down)

	a := <- drv_floors

	elevio.SetMotorDirection(elevio.MD_Stop)
	orderhandler.SetCurrentFloor(a)
	orderhandler.SetCurrentDir(0)
	elevio.SetFloorIndicator(a)
	fmt.Println("Heisen er intialisert og venter i etasje nr ", a)

	close(drv_floors) //vet ikke om funker?


	for !hasHitFloor {
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
	
}
*/


//Kan vel kanskje i stedet bare fjerne alle ordre og så kjøre update lights??
func InitializeLights(numFloors int){ //NB: Endra her navn til numHallButtons
	//Slår av lyset på alle lys
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

func StopElevator(){
	elevio.SetMotorDirection(elevio.MD_Stop)
	OpenDoor(3)
}

func OpenDoor(seconds time.Duration) {
	elevio.SetDoorOpenLamp(true)
	time.Sleep(seconds * time.Second)
	elevio.SetDoorOpenLamp(false)
}
