package fsm

import(
	"fmt"
	"../orderhandler"
	"../elevio"
	"../config"
	"../elevcontroller"
)



func RunElevator(ch config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator, activeElevators map[int]bool, elevator *config.Elevator){
	
    /* 		INIT 	*/
    NuActiveElevators := 0

	elevcontroller.Initialize(elevator)

	floor := <- ch.Drv_floors
	for floor == -1{
		floor = <- ch.Drv_floors
	}

    for _, result := range activeElevators{

    	if result == true{
    		NuActiveElevators++
    	}
    }

    elevator.Active = true
    elevator.ElevRank = NuActiveElevators

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(floor)
	elevator.CurrentFloor = floor
	elevator.CurrentDir = elevio.MD_Stop

	elevcontroller.PrintElevator(*elevator)

	//elevatorMap[elevID] = &elevator //Denne lokale lista vil oppdateres automatisk!
	
	/*		INIT FERDIG		*/
	
	fmt.Println("Heisen er intialisert og venter i etasje nr ", floor)

	for{
		//fmt.Println(activeElevators)

		switch elevatorMap[elevID].CurrentState{
		case config.IDLE:
			
			destination := elevatorMap[elevID].CurrentOrder
			//destination := config.Order{Floor: -1}
			if destination.Floor != -1{
				fmt.Println("Jeg har fått en oppgave i etasje ",destination.Floor,"! Denne skal jeg utføre")

				elevatorMap[elevID].CurrentDir = elevcontroller.GetDirection(*elevatorMap[elevID]) //kanskje jeg må bruke destination istedet for elevator.CurrentOrder. Ting kan fucke segf om currentorder endres! 

				elevio.SetMotorDirection(elevatorMap[elevID].CurrentDir)
				elevatorMap[elevID].CurrentState = config.ACTIVE

				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!

			}

		case config.ACTIVE:
			select{
			case reachedFloor := <- ch.Drv_floors: //treffet et floor
				elevio.SetFloorIndicator(reachedFloor)
				elevatorMap[elevID].CurrentFloor = reachedFloor

				if elevcontroller.ShouldStopAtFloor(*elevatorMap[elevID]){
					fmt.Println("stopping at floor")

					elevio.SetDoorOpenLamp(true)
					ch.Open_door <- true
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevatorMap[elevID].CurrentState = config.DOOR_OPEN

					ch.Stopping_at_floor <- reachedFloor //sender til de andre heisene slik at de kan slette alt i den etasjen.
				}
				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!

			default:
				if elevatorMap[elevID].CurrentDir == elevio.MD_Stop{
					fmt.Println("stopping at floor in ACTIVE")

					elevio.SetDoorOpenLamp(true)
					ch.Open_door <- true
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevatorMap[elevID].CurrentDir = elevio.MD_Stop
					elevatorMap[elevID].CurrentState = config.DOOR_OPEN
					
					go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!
				}

			}


		case config.DOOR_OPEN:

			select{
			case <- ch.Close_door:
	
				fmt.Println("closing door__")
				elevio.SetDoorOpenLamp(false) //slår av lys
				
				orderhandler.ClearCurrentFloor(elevatorMap[elevID])
				ch.LightUpdateCh <- true

				if elevatorMap[elevID].CurrentOrder.Floor == elevatorMap[elevID].CurrentFloor{
					elevatorMap[elevID].CurrentOrder.Floor = -1 //Fjerner currentOrder, siden den har utført den.
					elevatorMap[elevID].CurrentState = config.IDLE
					
				}else{
					//hvis den ikke er ferdig med currentOrder, fortsett i samme retning
					elevatorMap[elevID].CurrentDir = elevcontroller.GetDirection(*elevatorMap[elevID])
					elevio.SetMotorDirection(elevatorMap[elevID].CurrentDir) 
					elevatorMap[elevID].CurrentState = config.ACTIVE
				}
				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!
			}



		case config.UNDEFINED: //??
			//


		default:

		}
	}
}