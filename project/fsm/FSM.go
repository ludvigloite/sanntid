package fsm

import(
	"fmt"
	//"time"
	//"../elevcontroller"
	"../orderhandler"
	"../elevcontroller"
	"../elevio"
	"../config"
	//"../timer"
)

const(
	IDLE = "IDLE"
	ACTIVE = "ACTIVE"
	DOOR_OPEN = "DOOR_OPEN"
	UNDEFINED = "UNDEFINED"
)

func Init(floor_scanner <- chan int){
	elevio.SetMotorDirection(elevio.MD_Down)

	a := <- floor_scanner

	for a == -1{
		a = <- floor_scanner
	}
	

	elevio.SetMotorDirection(elevio.MD_Stop)
	orderhandler.SetCurrentFloor(a)
	orderhandler.SetCurrentDir(0)
	elevio.SetFloorIndicator(a)
	fmt.Println("Heisen er intialisert og venter i etasje nr ", a+1)
	//state := IDLE
}


func RunElevator(ch config.FSMChannels){
	state := IDLE

	///////////////////////////////////
	//	ch.Drv_buttons
    //	ch.Drv_floors
    //	ch.Open_door
    //	ch.Close_door
    //////////////////////////////////

    /* 		INIT 	*/
    
	elevio.SetMotorDirection(elevio.MD_Down)

	a := <- ch.Drv_floors

	for a == -1{
		a = <- ch.Drv_floors
	}

	

	elevio.SetMotorDirection(elevio.MD_Stop)
	orderhandler.SetCurrentFloor(a)
	orderhandler.SetCurrentDir(0)
	elevio.SetFloorIndicator(a)
	fmt.Println("Heisen er intialisert og venter i etasje nr ", a+1)
	


	for{
		switch state{
		case IDLE: //heis er IDLE. Skal ikke gjøre noe med mindre den får knappetrykk eller får inn en ordre som skal utføres
			//fmt.Println("JEG ER I IDLE")
			orderhandler.UpdateLights()
			select{
			case order := <- ch.Drv_buttons:
				fmt.Println("Knapp er trykket fra IDLE ",int(order.Button), order.Floor)
				//elevio.SetButtonLamp(order.Button, order.Floor, true)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				orderhandler.UpdateLights()

			default:
				//sjekk om en ordre skal utføres
				newOrder := orderhandler.GetNewOrder()
				if newOrder.Floor != -1{
					//Det er funnet en ordre
					fmt.Println("Det er funnet en ordre! Denne skal jeg utføre")
					orderhandler.AddOrder(newOrder.Floor, newOrder.ButtonType, orderhandler.GetElevatorID())
					orderhandler.SetCurrentOrder(newOrder.Floor)
					orderhandler.SetCurrentDir(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()))

					state = ACTIVE
				}
			}




		case ACTIVE:
			//fmt.Println("JEG ER I ACTIVE")
			elevio.SetMotorDirection(elevio.MotorDirection(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder())))
			orderhandler.UpdateLights()
			select{
			case order := <- ch.Drv_buttons: //Fått inn knappetrykk
				fmt.Println("Knapp er trykket fra ACTIVE ",int(order.Button), order.Floor)
				//elevio.SetButtonLamp(order.Button, order.Floor, true)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				orderhandler.UpdateLights()

			case reachedFloor := <- ch.Drv_floors:
				orderhandler.SetCurrentFloor(reachedFloor)
				elevio.SetFloorIndicator(reachedFloor)
				if orderhandler.ShouldStopAtFloor(reachedFloor, orderhandler.GetCurrentOrder(), orderhandler.GetElevatorID()){ //Kan jeg ikke bare ta variablene rett fra orderhandler??
					fmt.Println("stopping at floor")

					elevio.SetDoorOpenLamp(true)
					orderhandler.ClearFloor(reachedFloor)
					ch.Open_door <- true
					state = DOOR_OPEN

				}

			default:
				if orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()) == 0{
					elevcontroller.OpenDoor(3)
					orderhandler.ClearFloor(orderhandler.GetCurrentFloor())
					state = IDLE
				}

			}




		case DOOR_OPEN:
			elevio.SetMotorDirection(elevio.MD_Stop)
			fmt.Println("JEG ER I DOOR_OPEN")
/*
			if orderhandler.ShouldStopAtFloor(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder(), orderhandler.GetElevatorID()){
				open_door <- true
			}
*/

			select{
			case order := <- ch.Drv_buttons: //Fått inn knappetrykk //hvis det er en knapp på denne etasjen skal timeren starte på nytt
				fmt.Println("Knapp er trykket fra ACTIVE ",int(order.Button), order.Floor)
				//elevio.SetButtonLamp(order.Button, order.Floor, true)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				orderhandler.UpdateLights()


			case <- ch.Close_door:
				//døren skal lukkes. Det har gått 3 sek.
				fmt.Println("close door___")
				elevio.SetDoorOpenLamp(false) //slår av lys


				if orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()) == 0 {
					//kommet frem til enden.
					orderhandler.SetCurrentOrder(-1)
					state = IDLE
				}else{
					elevio.SetMotorDirection(elevio.MotorDirection(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder())))
					state = ACTIVE
				}



			default:

			}



		case UNDEFINED: //??
			//


		default:

		}
	}
}