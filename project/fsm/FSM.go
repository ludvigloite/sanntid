package fsm

import(
	"fmt"
	"time"
	//"../elevcontroller"
	"../orderhandler"
	"../elevcontroller"
	"../elevio"
	"../timer"
)

const(
	IDLE = "IDLE"
	ACTIVE = "ACTIVE"
	DOOR_OPEN = "DOOR_OPEN"
	UNDEFINED = "UNDEFINED"
)


func RunElevator(){

	state := IDLE

	//orderhandler.currentFloor
	//orderhandler.currentDir

	DOOR_OPEN_TIME := 3 * time.Second


	drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)
    open_door	:= make(chan bool)
    close_door	:= make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)
    go timer.DoorTimer(close_door,open_door,DOOR_OPEN_TIME) //bytt ut dette med noe i config. 3 secunder
    //Legg true på open_door når dør skal åpnes
    //skrives true til close_door når tiden er ute

    //init
	elevio.SetMotorDirection(elevio.MD_Down)
	

	a := <- drv_floors

	for a == -1{
		a = <- drv_floors
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
			case order := <- drv_buttons:
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
			case order := <- drv_buttons: //Fått inn knappetrykk
				fmt.Println("Knapp er trykket fra ACTIVE ",int(order.Button), order.Floor)
				//elevio.SetButtonLamp(order.Button, order.Floor, true)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				orderhandler.UpdateLights()

			case reachedFloor := <- drv_floors:
				orderhandler.SetCurrentFloor(reachedFloor)
				elevio.SetFloorIndicator(reachedFloor)
				if orderhandler.ShouldStopAtFloor(reachedFloor, orderhandler.GetCurrentOrder(), orderhandler.GetElevatorID()){ //Kan jeg ikke bare ta variablene rett fra orderhandler??
					fmt.Println("stopping at floor")

					elevio.SetDoorOpenLamp(true)
					orderhandler.ClearFloor(reachedFloor)
					open_door <- true
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
			case order := <- drv_buttons: //Fått inn knappetrykk //hvis det er en knapp på denne etasjen skal timeren starte på nytt
				fmt.Println("Knapp er trykket fra ACTIVE ",int(order.Button), order.Floor)
				//elevio.SetButtonLamp(order.Button, order.Floor, true)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				orderhandler.UpdateLights()


			case <- close_door:
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