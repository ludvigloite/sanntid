package arbitrator

import(
	"time"
	//"fmt"


	"../config"
	"../elevio"
)

func RankSolver(fsmCh config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator){
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
				go func(){fsmCh.New_state <- *elevatorMap[elevID]}() 
			}
		}
		if !masterExist{
			myRank = 1
			elevatorMap[elevID].ElevRank = myRank
			go func(){fsmCh.New_state <- *elevatorMap[elevID]}() 
		}
		time.Sleep(time.Second)
	}
}

func Arbitrator(ch config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator){ //kjøres bare av Master. Master kan bytte underveis. Derfor må det sjekkes hver gang og den må være inni while-loopen // KJØRES SOM GOROUNTINE
	order := config.Order{}
	for{
		if elevatorMap[elevID].ElevRank == 1 && !elevatorMap[elevID].HasRecentlyBeenDown{
			for i, elevator := range elevatorMap{ //går gjennom heisene.
				if elevator.Active && !elevator.Stuck{
					if elevator.CurrentOrder.Floor == -1{

						//Heis har ingen current orders! Finnes det noen nye ordre?
						time.Sleep(10 * time.Millisecond) //pass på at alle RemoveOrders har blitt synca
						order = GetNewOrder(*elevator, elevatorMap, elevID, i)
						
						if order.Floor != -1{
							//fmt.Println("Gitt ny CurrentOrder til ", i, " i etasje ", order.Floor)

							elevatorMap[i].CurrentOrder = order
							ch.New_current_order <- order
						}
					}
				}
			}
		}
	}
}


func GetNewOrder(elevator config.Elevator, elevatorMap map[int]*config.Elevator, masterElevID int, currentElevID int) config.Order{
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
		if elevator.CabOrders[3]{ //sjekker om det er en caborder i 3 etasje!
			newOrder.Floor = 3
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}

		for i := config.NUM_FLOORS-2; i > currentFloor; i--{ //Er det noen som går oppover over meg?
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !AnotherGoingToFloor(i,elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}

		for i := config.NUM_FLOORS-1; i > currentFloor; i--{ //Er det noen som går nedover over meg?
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !AnotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}

		if elevator.CabOrders[0]{ //sjekker om det er en caborder i 0 etasje!
			newOrder.Floor = 0
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}

		for i := 1; i < currentFloor; i++{ //er det noen som går nedover under meg?
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !AnotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}

		for i := 0; i < currentFloor; i++{ //Er det noen som går oppover under meg?
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !AnotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}
		
	}else{ //DEN GÅR NEDOVER!
		if elevator.CabOrders[0]{ //sjekker om det er en caborder i 0 etasje!
			newOrder.Floor = 0
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}
		for i := 1; i < currentFloor; i++{ //Er det noen som går nedover under meg?
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !AnotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}

		for i := 0; i < currentFloor; i++{ //Er det noen som går oppover under meg?
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !AnotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}

		if elevator.CabOrders[3]{ //sjekker om det er en caborder i 3 etasje!
			newOrder.Floor = 3
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}

		for i := config.NUM_FLOORS-2; i > currentFloor; i--{ //Er det noen som går oppover over meg?
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !AnotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}

		for i := config.NUM_FLOORS-1; i > currentFloor; i--{ //Er det noen som går nedover over meg?
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !AnotherGoingToFloor(i,elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}
	}
	
	return newOrder
}

func AnotherGoingToFloor(floor int, elevatorMap map[int]*config.Elevator, elevID int) bool{
	for _, elevator := range elevatorMap{
		if elevID != elevator.ElevID{
			if elevator.CurrentOrder.Floor == floor || (elevator.CurrentFloor == floor && elevator.CurrentFsmState != config.ACTIVE && elevator.Active){
				return true
			}
		}
	}
	return false
}