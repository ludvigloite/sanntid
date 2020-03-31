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
	for{
		select{
		case <- LightUpdateCh:
			//fmt.Println("Updating Lights...")
			for i := 0; i < config.NUM_FLOORS; i++{
				elevio.SetButtonLamp(elevio.BT_Cab, i, elevatorMap[elevID].CabOrders[i])

				for j := elevio.BT_HallUp; j != elevio.BT_Cab; j++{
					if i != 0 && j == elevio.BT_HallDown{
						elevio.SetButtonLamp(elevio.BT_HallDown, i, elevatorMap[elevID].HallOrders[i][elevio.BT_HallDown])
					}
					if i != config.NUM_FLOORS && j == elevio.BT_HallUp{
						elevio.SetButtonLamp(elevio.BT_HallUp, i, elevatorMap[elevID].HallOrders[i][elevio.BT_HallUp])
					}
				}
			}
		}
	}
}