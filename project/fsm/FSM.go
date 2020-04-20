//This module contains the FSM for the individual elevators. It writes to FSM channels when it needs to communicate with other modules.
package fsm

import(
	"fmt"
	

	"../elevio"
	"../config"
	"../elevcontroller"
)

func RunElevator(ch config.FSMChannels, elevator *config.Elevator){
	
    /* 		INIT 			*/

	elevcontroller.Initialize(elevator)


	floor := <- ch.Drv_floors
	for floor == -1{
		floor = <- ch.Drv_floors
	}

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(floor)
	elevator.CurrentFloor = floor
	elevator.Active = true

	ch.New_state <- *elevator

	fmt.Println("The elevator is initialized and waiting on floor ", floor)
	fmt.Println()

	/*		INIT FINISHED		*/

	for{

		switch elevator.CurrentFsmState{
		case config.IDLE:
			
			destination := elevator.CurrentOrder
			if destination.Floor != -1{
				elevator.CurrentDir = elevcontroller.GetDirection(*elevator) 
				elevio.SetMotorDirection(elevator.CurrentDir)
				elevator.CurrentFsmState = config.ACTIVE

				if elevator.CurrentDir == elevio.MD_Stop{ //Received new order at currentFloor
					elevio.SetDoorOpenLamp(true)
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevator.CurrentFsmState = config.DOOR_OPEN

					ch.Open_door <- true
					ch.Stopping_at_floor <- elevator.CurrentFloor
					ch.Watchdog_updater <- true
				}

				ch.New_state <- *elevator
			}


		case config.ACTIVE:

			select{
			case reachedFloor := <- ch.Drv_floors:
				elevio.SetFloorIndicator(reachedFloor)
				elevator.CurrentFloor = reachedFloor
				elevator.Stuck = false
				ch.Watchdog_updater <- true

				elevio.SetMotorDirection(elevio.MD_Stop)
				elevator.CurrentFsmState = config.IDLE

				if elevcontroller.ShouldStopAtFloor(*elevator){

					elevio.SetDoorOpenLamp(true)
					ch.Open_door <- true
					elevator.CurrentFsmState = config.DOOR_OPEN

					ch.Stopping_at_floor <- reachedFloor //Sending order to all the other elevators to delete HallOrders at this floor
				}
				ch.New_state <- *elevator
			}


		case config.DOOR_OPEN:

			select{
			case <- ch.Close_door:	
				elevio.SetDoorOpenLamp(false)

				elevator.CurrentFsmState = config.IDLE

				if elevator.CurrentOrder.Floor == elevator.CurrentFloor{
					elevator.CurrentOrder.Floor = -1 //currentOrder is finished
				}
				
				ch.New_state <- *elevator
			}


		}
	}
}