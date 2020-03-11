package fsm

import(
	"fmt"
	//"../elevcontroller"
	"../orderhandler"
	"../elevcontroller"
	"../elevio"
)

const(
	IDLE = "IDLE"
	ACTIVE = "ACTIVE"
	UNDEFINED = "UNDEFINED"
)


func RunElevator(){

	state := IDLE

	//orderhandler.currentFloor
	//orderhandler.currentDir


	drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)

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
					elevcontroller.StopElevator()
					orderhandler.ClearFloor(reachedFloor)
					if orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()) == 0 {
						//kommet frem til enden.
						orderhandler.SetCurrentOrder(-1)
						state = IDLE
					}else{
						elevio.SetMotorDirection(elevio.MotorDirection(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder())))
					}
				}

			default:
				if orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()) == 0{
					elevcontroller.OpenDoor(3)
					orderhandler.ClearFloor(orderhandler.GetCurrentFloor())
					state = IDLE
				}

			}









		case UNDEFINED: //??
			//


		}
	}
}