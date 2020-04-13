package orderhandler

import(
	"../elevio"
	"../config"
	//"fmt"
)
/*
func ClearCurrentFloor(elevator *config.Elevator){
	currentFloor := elevator.CurrentFloor

	elevator.HallOrders[currentFloor][elevio.BT_HallDown] = false
	elevator.HallOrders[currentFloor][elevio.BT_HallUp] = false
	elevator.CabOrders[currentFloor] = false

}*/


func LightUpdater(LightUpdateCh <-chan bool, elevatorMap map[int]*config.Elevator, elevID int){ //DENNE MÅ ENDRES SLIK AT DEN BARE ENDRER LYS OM DET FAKTISK ER EN ENDRING!!
	empty_elevator := config.Elevator{}
	for{
		select{
		case <- LightUpdateCh:
			//fmt.Println("Updating Lights...")
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