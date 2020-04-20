//This module handles new orders being given to the elevators by the Master-elevator. It also makes sure that there are always one and only one Master on the network.
package arbitrator

import(
	"time"


	"../config"
	"../elevio"
)

//This function makes sure that there always is one and only one master.
func RankSolver(New_state chan<- config.Elevator, elevatorMap map[int]*config.Elevator, elevID int){
	masterExist := false
	for{
		masterExist = false
		myRank := elevatorMap[elevID].ElevRank

		for id,elev := range elevatorMap{
			if elev.ElevRank == 1 && elev.Active{
				masterExist = true
			}
			if id != elevID && elev.Active && elev.ElevRank == myRank{				
				if myRank != 1{
					myRank--
				}else if myRank != 3{
					myRank++
				}
				if myRank==1{
					masterExist = true
				}
				elevatorMap[elevID].ElevRank = myRank
				New_state <- *elevatorMap[elevID]
			}
		}
		if !masterExist{
			myRank = 1
			elevatorMap[elevID].ElevRank = myRank
			New_state <- *elevatorMap[elevID]
		}
		time.Sleep(time.Second)
	}
}

//This function can only be ran by Master, and gives out orders to all the elevators.
func Arbitrator(New_current_order chan<- config.Order, elevatorMap map[int]*config.Elevator, elevID int){
	order := config.Order{}
	for{
		if elevatorMap[elevID].ElevRank == 1 && !elevatorMap[elevID].HasRecentlyBeenDown{
			for i, elevator := range elevatorMap{
				if elevator.Active && !elevator.Stuck{
					if elevator.CurrentOrder.Floor == -1{ //Elevator has no currentOrder

						time.Sleep(10 * time.Millisecond) //Make sure all RemoveOrders are synced.
						order = getNewOrder(*elevator, elevatorMap, elevID, i)
						
						if order.Floor != -1{
							elevatorMap[i].CurrentOrder = order
							New_current_order <- order
						}
					}
				}
			}
		}
	}
}

//Finding new order for the spesific elevator
func getNewOrder(elevator config.Elevator, elevatorMap map[int]*config.Elevator, masterElevID int, currentElevID int) config.Order{
	newOrder := config.Order{
		Sender_elev_ID: masterElevID,
		Sender_elev_rank: 1,
		Receiver_elev: currentElevID,
		Floor: -1,
		ButtonType: elevio.BT_HallDown,
	}

	masterElev := elevatorMap[masterElevID]
	currentFloor := elevator.CurrentFloor

	if currentFloor == -1{
		return newOrder
	}

	//Is there any order at my CurrentFloor?
	if elevator.CabOrders[currentFloor]{
		newOrder.Floor = currentFloor
		newOrder.ButtonType = elevio.BT_Cab
		return newOrder
	}
	if masterElev.HallOrders[currentFloor][elevio.BT_HallUp]{
		newOrder.Floor = currentFloor
		newOrder.ButtonType = elevio.BT_HallUp
		return newOrder
	}
	if masterElev.HallOrders[currentFloor][elevio.BT_HallDown]{
		newOrder.Floor = currentFloor
		newOrder.ButtonType = elevio.BT_HallDown
		return newOrder
	}

	
	if elevator.CurrentDir == elevio.MD_Up{
		if elevator.CabOrders[3]{ //Is there any cab orders in 3rd floor?
			newOrder.Floor = 3
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}

		for i := config.NUM_FLOORS-2; i > currentFloor; i--{ //Is there any order going upwards over me?
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !anotherGoingToFloor(i,elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}

		for i := config.NUM_FLOORS-1; i > currentFloor; i--{ //Is there any order going downwards over me?
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !anotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}

		if elevator.CabOrders[0]{ //Is there any cab order in 0 floor?
			newOrder.Floor = 0
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}

		for i := 1; i < currentFloor; i++{ //Is there any order going downwards under me?
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !anotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}

		for i := 0; i < currentFloor; i++{ //Is there any order going upwards under me?
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !anotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}
		
	}else{ //Going downwards
		if elevator.CabOrders[0]{ //Is there any cab order in 0 floor?
			newOrder.Floor = 0
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}
		for i := 1; i < currentFloor; i++{ //Is there any order going downwards under me?
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !anotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}

		for i := 0; i < currentFloor; i++{ //Is there any order going upwards under me?
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !anotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}

		if elevator.CabOrders[3]{ //Is there any cab order in 3rd floor?
			newOrder.Floor = 3
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}

		for i := config.NUM_FLOORS-2; i > currentFloor; i--{ //Is there any order going upwards over me?
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !anotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}

		for i := config.NUM_FLOORS-1; i > currentFloor; i--{ //Is there any order going downwards over me?
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !anotherGoingToFloor(i,elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}
	}
	
	return newOrder
}

//checking if another elevator is going to the floor
func anotherGoingToFloor(floor int, elevatorMap map[int]*config.Elevator, elevID int) bool{
	for _, elevator := range elevatorMap{
		if elevID != elevator.ElevID{
			if elevator.CurrentOrder.Floor == floor || (elevator.CurrentFloor == floor && elevator.CurrentFsmState != config.ACTIVE && elevator.Active){
				return true
			}
		}
	}
	return false
}