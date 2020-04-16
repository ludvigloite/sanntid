package fsm

import(
	"fmt"
	

	"../elevio"
	"../config"
	"../elevcontroller"
)

func RunElevator(ch config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator, elevator *config.Elevator){
	
    /* 		INIT 			*/

	elevcontroller.Initialize(elevator)


	floor := <- ch.Drv_floors
	for floor == -1{
		floor = <- ch.Drv_floors
	}

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(floor)
	elevator.CurrentFloor = floor //Siden det er snakk om pekere vil dette være det samme som elevatorMap[elevID].CurrentFloor = floor
	elevator.Active = true

	ch.New_state <- *elevator

	fmt.Println("Heisen er intialisert og venter i etasje nr ", floor)
	fmt.Println()

	/*		INIT FERDIG		*/

	//elevcontroller.PrintElevator(elevator)

	for{

		switch elevatorMap[elevID].CurrentFsmState{
		case config.IDLE:
			
			destination := elevatorMap[elevID].CurrentOrder
			if destination.Floor != -1{
				elevatorMap[elevID].CurrentDir = elevcontroller.GetDirection(*elevatorMap[elevID]) 
				elevio.SetMotorDirection(elevatorMap[elevID].CurrentDir)
				elevatorMap[elevID].CurrentFsmState = config.ACTIVE

				if elevatorMap[elevID].CurrentDir == elevio.MD_Stop{ //kommet ny ordre der du allerede er
					elevio.SetDoorOpenLamp(true)
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevatorMap[elevID].CurrentFsmState = config.DOOR_OPEN

					ch.Open_door <- true
					ch.Stopping_at_floor <- elevatorMap[elevID].CurrentFloor
					ch.Watchdog_updater <- true
				}

				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!
			}


		case config.ACTIVE:

			select{
			case reachedFloor := <- ch.Drv_floors:
				elevio.SetFloorIndicator(reachedFloor)
				elevatorMap[elevID].CurrentFloor = reachedFloor
				elevatorMap[elevID].Stuck = false
				ch.Watchdog_updater <- true

				elevio.SetMotorDirection(elevio.MD_Stop)
				elevatorMap[elevID].CurrentFsmState = config.IDLE

				if elevcontroller.ShouldStopAtFloor(*elevatorMap[elevID]){

					elevio.SetDoorOpenLamp(true)
					ch.Open_door <- true
					elevatorMap[elevID].CurrentFsmState = config.DOOR_OPEN

					ch.Stopping_at_floor <- reachedFloor //sender til de andre heisene slik at de kan slette alt i den etasjen.
				}
				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!
			}


		case config.DOOR_OPEN:

			select{
			case <- ch.Close_door:	
				elevio.SetDoorOpenLamp(false)

				elevatorMap[elevID].CurrentFsmState = config.IDLE

				if elevatorMap[elevID].CurrentOrder.Floor == elevatorMap[elevID].CurrentFloor{
					elevatorMap[elevID].CurrentOrder.Floor = -1 //Fjerner currentOrder, siden den har utført den.
				}
				
				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!
			}


		}
	}
}