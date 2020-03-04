package fsm

import(
	"fmt"
	//"../elevcontroller"
	"../orderhandler"
)

const(
	IDLE = "IDLE"
	ACTIVE = "ACTIVE"
	RESET = "RESET"
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




	for{
		switch state{
		case IDLE: //heis er IDLE. Skal ikke gjøre noe med mindre den får knappetrykk eller får inn en ordre som skal utføres
			select{
			case order := <- drv_buttons:
				fmt.Println("Knapp er trykket fra IDLE")
				//elevio.SetButtonLamp(order.Button, order.Floor, true)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				orderhandler.UpdateLights()

			default:
				//sjekk om en ordre skal utføres
				newOrder := orderhandler.GetNewOrder()
				if newOrder.Floor != -1{
					//Det er funnet en ordre
					fmt.Println("Det er funnet en ordre! Denne skal jeg utføre")
					orderhandler.AddOrder(newOrder.Floor, newOrder.ButtonType, orderhandler.elevatorID)
					orderhandler.SetCurrentOrder(newOrder.Floor)
					orderhandler.SetCurrentDir(orderhandler.GetDirection(orderhandler.currentFloor, orderhandler.currentOrder))

					state = ACTIVE
				}
			}




		case ACTIVE:
			elevio.SetMotorDirection(elevio.SetMotorDirection(orderhandler.GetDirection(orderhandler.currentFloor, orderhandler.currentOrder)))
			
			select{
			case order := <- drv_buttons: //Fått inn knappetrykk
				fmt.Println("Knapp er trykket fra ACTIVE")
				//elevio.SetButtonLamp(order.Button, order.Floor, true)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				orderhandler.UpdateLights()

			case reachedFloor := <- drv_floors:
				orderhandler.currentFloor = reachedFloor
				elevio.SetFloorIndicator(floor)
				if orderhandler.ShouldStopAtFloor(floor, orderhandler.currentOrder, orderhandler.elevID){ //Kan jeg ikke bare ta variablene rett fra orderhandler??
					fmt.Println("stopping at floor")
					elevcontroller.StopElevator()
					orderhandler.ClearFloor()
					if orderhandler.GetDirection(orderhandler.currentFloor, orderhandler.currentOrder) == 0 {
						//kommet frem til enden.
						orderhandler.SetCurrentOrder(-1)
						state = IDLE
					}else{
						elevio.SetMotorDirection(elevio.SetMotorDirection(orderhandler.GetDirection(orderhandler.currentFloor, orderhandler.currentOrder)))
					}
				}

			default:


			}







		case RESET:
			//


		}
	}
}