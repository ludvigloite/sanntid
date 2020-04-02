package fsm

import(
	"fmt"
	//"../orderhandler"
	"../elevio"
	"../config"
	"../elevcontroller"
)



func RunElevator(ch config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator, elevator *config.Elevator){

    /* 		INIT 	*/

	elevcontroller.Initialize(elevator)

	floor := <- ch.Drv_floors
	for floor == -1{
		floor = <- ch.Drv_floors
	}

	NuActiveElevators := 0

	for _, elevator := range elevatorMap{
		if elevator.Active{
			NuActiveElevators++
		}
	}


    elevator.Active = true
    elevator.ElevRank = NuActiveElevators + 1
    fmt.Println("JEG HAR RANK ",elevator.ElevRank)

	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(floor)
	elevator.CurrentFloor = floor

	/*		INIT FERDIG		*/

	ch.New_state <- *elevator

	fmt.Println("Heisen er intialisert og venter i etasje nr ", floor)

	for{

		switch elevatorMap[elevID].CurrentState{
		case config.IDLE:

			destination := elevatorMap[elevID].CurrentOrder
			if destination.Floor != -1{
				//fmt.Println("Jeg har fått en oppgave i etasje ",destination.Floor,"! Denne skal jeg utføre")

				elevatorMap[elevID].CurrentDir = elevcontroller.GetDirection(*elevatorMap[elevID]) //kanskje jeg må bruke destination istedet for elevator.CurrentOrder. Ting kan fucke segf om currentorder endres!

				elevio.SetMotorDirection(elevatorMap[elevID].CurrentDir)
				elevatorMap[elevID].CurrentState = config.ACTIVE

				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!

			}

		case config.ACTIVE:
			select{
			case reachedFloor := <- ch.Drv_floors: //treffet et floor
				fmt.Println("Passerte etasje ", reachedFloor)
				elevio.SetFloorIndicator(reachedFloor)
				elevatorMap[elevID].CurrentFloor = reachedFloor

				if elevcontroller.ShouldStopAtFloor(*elevatorMap[elevID]){
					//fmt.Println("stopping at floor")

					elevio.SetDoorOpenLamp(true)
					ch.Open_door <- true
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevatorMap[elevID].CurrentState = config.DOOR_OPEN

					ch.Stopping_at_floor <- reachedFloor //sender til de andre heisene slik at de kan slette alt i den etasjen.
				}
				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!

			default:

				if elevatorMap[elevID].CurrentDir == elevio.MD_Stop{
					//fmt.Println("stopping at floor I am already in")

					elevio.SetDoorOpenLamp(true)
					ch.Open_door <- true
					ch.Stopping_at_floor <- elevatorMap[elevID].CurrentFloor

					elevio.SetMotorDirection(elevio.MD_Stop)
					elevatorMap[elevID].CurrentState = config.DOOR_OPEN

					go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!
				}

			}


		case config.DOOR_OPEN:

			select{
			case <- ch.Close_door:

				fmt.Println("closing door__")
				elevio.SetDoorOpenLamp(false) //slår av lys

				elevatorMap[elevID].CurrentState = config.IDLE

				if elevatorMap[elevID].CurrentOrder.Floor == elevatorMap[elevID].CurrentFloor{
					elevatorMap[elevID].CurrentOrder.Floor = -1 //Fjerner currentOrder, siden den har utført den.
					//fmt.Println("CurrentOrder er fjernet!")
				}

				go func(){ch.New_state <- *elevatorMap[elevID]}() //sender kun sin egen Elevator!
			}



		case config.UNDEFINED: //??
			//


		default:

		}
	}
}
