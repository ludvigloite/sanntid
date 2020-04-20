//This module contains all help-functions for the elevator.
// - Initialize() initializes the elevator
// - GetDirection() and ShouldStopAtFloor() are used by the FSM
// - LightUpdater() updates all the lights.
package elevcontroller

import(
	"../elevio"
	"../config"
)

func Initialize(elevator *config.Elevator){
	elevio.SetMotorDirection(elevio.MD_Down)
	resetLights()
	initQueues(elevator)
}

func initQueues(elevator *config.Elevator){
	for flr := 0; flr < config.NUM_FLOORS; flr++ {
		elevator.CabOrders[flr] = false

		for btn := elevio.BT_HallUp; btn != elevio.BT_Cab; btn++{
			elevator.HallOrders[flr][btn] = false
		}
	}
}

func resetLights(){
	numFloors := config.NUM_FLOORS
	elevio.SetDoorOpenLamp(false)
	for i := 0; i < numFloors; i++{
		elevio.SetButtonLamp(elevio.BT_Cab, i, false)
		if i != 0{
			elevio.SetButtonLamp(elevio.BT_HallDown,i,false)
		}
		if i != numFloors{
			elevio.SetButtonLamp(elevio.BT_HallUp,i,false)
		}
	}
}

func GetDirection(elevator config.Elevator) elevio.MotorDirection{
	currentFloor := elevator.CurrentFloor
	destinationFloor := elevator.CurrentOrder.Floor

	if destinationFloor == -1 || destinationFloor == currentFloor {
		return elevio.MD_Stop

	} else if currentFloor < destinationFloor {
		return elevio.MD_Up

	} else{
		return elevio.MD_Down
	}
}

func ShouldStopAtFloor(elevator config.Elevator) bool{

	currentFloor := elevator.CurrentFloor
	dir := elevator.CurrentDir

	if elevator.CurrentOrder.Floor == currentFloor{
		return true
	}

	if elevator.CabOrders[currentFloor]{
		return true
	}
	if elevator.HallOrders[currentFloor][elevio.BT_HallUp] && dir == elevio.MD_Up{
		return true
	}
	if elevator.HallOrders[currentFloor][elevio.BT_HallDown] && dir == elevio.MD_Down {
		return true
	}
	return false
}

func LightUpdater(LightUpdateCh <-chan bool, elevatorMap map[int]*config.Elevator, elevID int){
	empty_elevator := config.Elevator{}
	for{
		select{
		case <- LightUpdateCh:
			
			elevator := elevatorMap[elevID]

			if !config.SHOW_ORDERS_WHEN_NETWORK_DOWN && elevator.NetworkDown{
				elevator = &empty_elevator
			}
			for i := 0; i < config.NUM_FLOORS; i++{
				elevio.SetButtonLamp(elevio.BT_Cab, i, elevator.CabOrders[i])

				for j := elevio.BT_HallUp; j != elevio.BT_Cab; j++{
					if i != 0 && j == elevio.BT_HallDown{
						elevio.SetButtonLamp(elevio.BT_HallDown, i, elevator.HallOrders[i][elevio.BT_HallDown])
					}
					if i != config.NUM_FLOORS && j == elevio.BT_HallUp{
						elevio.SetButtonLamp(elevio.BT_HallUp, i, elevator.HallOrders[i][elevio.BT_HallUp])
					}
				}
			}
		}
	}
}