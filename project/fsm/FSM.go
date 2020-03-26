package fsm

import(
	"fmt"
	"../orderhandler"
	"../elevio"
	"../config"
)



func RunElevator(ch config.FSMChannels, elevID int, elevatorList *[config.NUM_ELEVATORS] config.Elevator){
	state := config.IDLE
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
	fmt.Println("Heisen er intialisert og venter i etasje nr ", a)
	/*		INIT FERDIG		*/

	//INITIALISER ELEV I ELEVATOR_LIST!!
	elevator = config.Elevator{
		ElevID: elevID,
		ElevRank: 3,
		CurrentOrder: config.Order{Floor:-1, ButtonType:-1},
		CurrentFloor: a,
		CurrentState: config.IDLE,
	}
	elevatorList[id] = elevator
	
	stateIsChanged := true

	for{
		

		switch state{
		case config.IDLE: //heis er IDLE. Skal ikke gjøre noe med mindre den får knappetrykk eller får inn en ordre som skal utføres
			
			elevator := orderhandler.GetElevList()[orderhandler.GetElevID()-1]
			destination := elevator.CurrentOrder
			if destination.Floor != -1{
				fmt.Println("Jeg har fått en oppgave! Denne skal jeg utføre")
				//orderhandler.AddOrder(newOrder.Floor, newOrder.ButtonType, orderhandler.GetElevID())
				elevator.CurrentOrder = destination
				//orderhandler.SetCurrentOrder(destination.Floor)

				HER JOBBER JEG NÅ!!! prøver å kutte ut alle lokale variable for currentFloor osv!!

				orderhandler.SetCurrentDir(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()))

				elevio.SetMotorDirection(elevio.MotorDirection(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder())))
				state = config.ACTIVE
				stateIsChanged = true

				orderhandler.SetCurrentState(1)
			}
			if stateIsChanged{
				go func(){ch.New_state <- }
			}

		case config.ACTIVE:
			//orderhandler.UpdateLights()
			select{
			case reachedFloor := <- ch.Drv_floors:
				orderhandler.SetCurrentFloor(reachedFloor)
				elevio.SetFloorIndicator(reachedFloor)
				if orderhandler.ShouldStopAtFloor(reachedFloor, orderhandler.GetCurrentOrder(), orderhandler.GetElevID()){
					fmt.Println("stopping at floor")

					elevio.SetDoorOpenLamp(true)
					//orderhandler.ClearFloor(reachedFloor)
					//orderhandler.UpdateLights()
					ch.Open_door <- true

					elevio.SetMotorDirection(elevio.MD_Stop)//
					state = config.DOOR_OPEN
					orderhandler.SetCurrentState(2)
					//orderhandler.UpdateLights()
				}

			default:
				if orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()) == 0{
					fmt.Println("stopping at floor in ACTIVE")

					elevio.SetDoorOpenLamp(true)
					//orderhandler.ClearFloor(orderhandler.GetCurrentFloor())
					//orderhandler.UpdateLights()
					ch.Open_door <- true

					elevio.SetMotorDirection(elevio.MD_Stop)//
					state = config.DOOR_OPEN
					orderhandler.SetCurrentState(2)
					//orderhandler.UpdateLights()
				}

			}




		case config.DOOR_OPEN:
			//elevio.SetMotorDirection(elevio.MD_Stop)

			select{
			case <- ch.Close_door:
	
				fmt.Println("closing door__")
				elevio.SetDoorOpenLamp(false) //slår av lys
				
				orderhandler.ClearFloor(orderhandler.GetCurrentFloor()) //
				//orderhandler.UpdateLights() //
				ch.LightUpdateCh <- true


				//orderhandler.UpdateLights()

				if orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder()) == 0 {
					//kommet frem til enden.
					orderhandler.SetCurrentOrder(-1)

					state = config.IDLE
					orderhandler.SetCurrentState(0)
				}else{
					elevio.SetMotorDirection(elevio.MotorDirection(orderhandler.GetDirection(orderhandler.GetCurrentFloor(), orderhandler.GetCurrentOrder())))

					state = config.ACTIVE
					orderhandler.SetCurrentState(1)
				}
			//orderhandler.UpdateLights()



			default:

			}



		case config.UNDEFINED: //??
			//


		default:

		}
	}
}