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
		Active: true,
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
	
	/*		INIT FERDIG		*/
	fmt.Println("Heisen er intialisert og venter i etasje nr ", floor)

	for{

		switch *elevatorMap[elevID].CurrentState{
		case config.IDLE:
			
			destination := *elevatorMap[elevID].CurrentOrder
			if destination.Floor != -1{
				fmt.Println("Jeg har fått en oppgave i etasje ",destination.Floor,"! Denne skal jeg utføre")

				*elevatorMap[elevID].CurrentDir = elevcontroller.GetDirection(*elevatorMap[elevID]) //kanskje jeg må bruke destination istedet for elevator.CurrentOrder. Ting kan fucke segf om currentorder endres! 

				elevio.SetMotorDirection(*elevatorMap[elevID].CurrentDir)
				*elevatorMap[elevID].CurrentState = config.ACTIVE

				go func(){ch.New_state <- *elevatorMap[elevID]} //sender kun sin egen Elevator!

			}

		case config.ACTIVE:
			select{
			case reachedFloor := <- ch.Drv_floors: //treffet et floor
				elevio.SetFloorIndicator(reachedFloor)
				*elevatorMap[elevID].CurrentFloor = reachedFloor

				if elevcontroller.ShouldStopAtFloor(*elevatorMap[elevID]){
					fmt.Println("stopping at floor")

					elevio.SetDoorOpenLamp(true)
					ch.Open_door <- true
					elevio.SetMotorDirection(elevio.MD_Stop)
					*elevatorMap[elevID].CurrentState = config.DOOR_OPEN

					ch.Stopping_at_floor <- reachedFloor //sender til de andre heisene slik at de kan slette alt i den etasjen.
				}
				go func(){ch.New_state <- *elevatorMap[elevID]} //sender kun sin egen Elevator!

			default:
				if *elevatorMap[elevID].CurrentDir == config.MD_Stop{
					fmt.Println("stopping at floor in ACTIVE")

					elevio.SetDoorOpenLamp(true)
					ch.Open_door <- true
					elevio.SetMotorDirection(elevio.MD_Stop)
					*elevatorMap[elevID].CurrentDir = elevio.MD_Stop
					*elevatorMap[elevID].CurrentState = config.DOOR_OPEN
					
					go func(){ch.New_state <- *elevatorMap[elevID]} //sender kun sin egen Elevator!
				}

			}


		case config.DOOR_OPEN:

			select{
			case <- ch.Close_door:
	
				fmt.Println("closing door__")
				elevio.SetDoorOpenLamp(false) //slår av lys
				
				orderhandler.ClearCurrentFloor(&elevatorMap[elevID])
				ch.LightUpdateCh <- true

				if *elevatorMap[elevID].CurrentOrder.Floor == *elevatorMap[elevID].CurrentFloor{
					*elevatorMap[elevID].CurrentOrder.Floor = -1 //Fjerner currentOrder, siden den har utført den.
					*elevatorMap[elevID].CurrentState = config.IDLE
					
				}else{
					//hvis den ikke er ferdig med currentOrder, fortsett i samme retning
					*elevatorMap[elevID].CurrentDir = elevcontroller.GetDirection(*elevatorMap[elevID])
					elevio.SetMotorDirection(*elevatorMap[elevID].CurrentDir) 
					*elevatorMap[elevID].CurrentState = config.ACTIVE
				}
				go func(){ch.New_state <- *elevatorMap[elevID]} //sender kun sin egen Elevator!
			}



		case config.UNDEFINED: //??
			//


		default:

		}
	}
}