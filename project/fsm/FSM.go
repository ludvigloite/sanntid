package fsm

import(
	"fmt"
	"../orderhandler"
	"../elevio"
	"../config"
)



func RunElevator(ch config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator){
	
    /* 		INIT 	*/

    p := <-ch.PeerUpdateCh
	NuActiveElevators := length(p.Peers)

	elevator = config.Elevator{
		ElevID: elevID,
		ElevRank: NuActiveElevators, //Dette fikses ved at man sjekker hvor mange heiser som er online.
		CurrentOrder: config.Order{Floor:-1, ButtonType:-1}, //usikker på om denne initialiseringen funker.
		CurrentFloor: -1,
		CurrentDir: elevio.MD_Down,
		CurrentState: config.IDLE,
		CabOrders: [config.NUM_FLOORS]bool{},
		HallOrders: [config.NUM_FLOORS][config.NUM_HALLBUTTONS]bool{},
	}

	elevcontroller.Initialize(&elevator)

	floor := <- ch.Drv_floors
	for floor == -1{
		floor = <- ch.Drv_floors
	}

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(floor)
	elevator.CurrentFloor = floor
	elevator.CurrentDir = elevio.MD_Stop

	elevcontroller.PrintElevator(elevator)

	elevatorMap[elevID] = &elevator //Denne lokale lista vil oppdateres automatisk!
	
	stateIsChanged := true

	/*		INIT FERDIG		*/
	fmt.Println("Heisen er intialisert og venter i etasje nr ", floor)

	for{

		switch *elevatorList[elevID].CurrentState{
		case config.IDLE:
			
			destination := *elevatorMap[elevID].CurrentOrder
			if destination.Floor != -1{
				fmt.Println("Jeg har fått en oppgave i etasje ",destination.Floor,"! Denne skal jeg utføre")

				*elevatorMap[elevID].CurrentDir = elevcontroller.GetDirection(*elevatorMap[elevID].CurrentFloor, *elevatorMap[elevID].CurrentOrder.Floor) //kanskje jeg må bruke destination istedet for elevator.CurrentOrder. Ting kan fucke segf om currentorder endres! 

				elevio.SetMotorDirection(*elevatorMap[elevID].CurrentDir)
				*elevatorMap[elevID].CurrentState = config.ACTIVE
				stateIsChanged = true

			}
			if stateIsChanged{
				go func(){ch.New_state <- *elevatorMap[elevID]} //sender kun sin egen Elevator!
				stateIsChanged = false
			}

		case config.ACTIVE:
			select{
			case reachedFloor := <- ch.Drv_floors:
				stateIsChanged = true
				elevio.SetFloorIndicator(reachedFloor)
				*elevatorMap[elevID].CurrentFloor = reachedFloor

				if orderhandler.ShouldStopAtFloor(*elevatorMap[elevID]){
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