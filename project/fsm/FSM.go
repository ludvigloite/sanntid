package fsm

import(
	"fmt"
	"../orderhandler"
	"../elevio"
	"../config"
)

const(
	IDLE = "IDLE"
	ACTIVE = "ACTIVE"
	DOOR_OPEN = "DOOR_OPEN"
	UNDEFINED = "UNDEFINED"
)


func RunElevator(ch config.FSMChannels){
	state := IDLE
	orderhandler.SetCurrentState(0)


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
	/*		INIT FERDIG		*/
	


	for{
		//orderhandler.UpdateLights()

		switch state{
		case IDLE: //heis er IDLE. Skal ikke gjøre noe med mindre den får knappetrykk eller får inn en ordre som skal utføres
			newOrder := orderhandler.GetNewOrder()
			if newOrder.Floor != -1{
				fmt.Println("Det er funnet en ordre! Denne skal jeg utføre")
				orderhandler.AddOrder(newOrder.Floor, newOrder.ButtonType, orderhandler.GetElevID())
				orderhandler.SetCurrentOrder(newOrder.Floor)
				orderhandler.SetCurrentDir(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()))

				elevio.SetMotorDirection(elevio.MotorDirection(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder())))
				state = ACTIVE
				orderhandler.SetCurrentState(1)
			}

		case ACTIVE:
			//orderhandler.UpdateLights()
			select{
			case reachedFloor := <- ch.Drv_floors:
				orderhandler.SetCurrentFloor(reachedFloor)
				elevio.SetFloorIndicator(reachedFloor)
				if orderhandler.ShouldStopAtFloor(reachedFloor, orderhandler.GetCurrentOrder(), orderhandler.GetElevID()){
					fmt.Println("stopping at floor")

					elevio.SetDoorOpenLamp(true)
					orderhandler.ClearFloor(reachedFloor)
					orderhandler.UpdateLights()
					ch.Open_door <- true

					elevio.SetMotorDirection(elevio.MD_Stop)//
					state = DOOR_OPEN
					orderhandler.SetCurrentState(2)
					//orderhandler.UpdateLights()
				}

			default:
				if orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()) == 0{
					fmt.Println("stopping at floor")

					elevio.SetDoorOpenLamp(true)
					orderhandler.ClearFloor(orderhandler.GetCurrentFloor())
					orderhandler.UpdateLights()
					ch.Open_door <- true

					elevio.SetMotorDirection(elevio.MD_Stop)//
					state = DOOR_OPEN
					orderhandler.SetCurrentState(2)
					//orderhandler.UpdateLights()
				}

			}




		case DOOR_OPEN:
			//elevio.SetMotorDirection(elevio.MD_Stop)

			select{
			case <- ch.Close_door:
	
				fmt.Println("closing door__")
				elevio.SetDoorOpenLamp(false) //slår av lys
				orderhandler.UpdateLights()

				if orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()) == 0 {
					//kommet frem til enden.
					orderhandler.SetCurrentOrder(-1)

					state = IDLE
					orderhandler.SetCurrentState(0)
				}else{
					elevio.SetMotorDirection(elevio.MotorDirection(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder())))

					state = ACTIVE
					orderhandler.SetCurrentState(1)
				}
			//orderhandler.UpdateLights()



			default:

			}



		case UNDEFINED: //??
			//


		default:

		}
	}
}