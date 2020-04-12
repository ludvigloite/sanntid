package arbitrator

import(
	//"fmt"
	"../config"
	"../elevio"
	//"time"
	//"../elevcontroller"
)

func GetNewOrder(elevator config.Elevator, elevatorMap map[int]*config.Elevator, masterElev int, currentElev int) config.Order{
	newOrder := config.Order{
		Sender_elev_ID: masterElev,
		Sender_elev_rank: 1,
		Receiver_elev: currentElev,
		Floor: -1,
		ButtonType: elevio.BT_HallDown,
	}

	currentFloor := elevator.CurrentFloor
	if currentFloor == -1{
		return newOrder
	}

	if elevator.CabOrders[currentFloor]{
		newOrder.Floor = currentFloor
		newOrder.ButtonType = elevio.BT_Cab
		return newOrder
	}
	if elevator.HallOrders[currentFloor][elevio.BT_HallUp]{
		newOrder.Floor = currentFloor
		newOrder.ButtonType = elevio.BT_HallUp
		return newOrder
	}
	if elevator.HallOrders[currentFloor][elevio.BT_HallDown]{
		newOrder.Floor = currentFloor
		newOrder.ButtonType = elevio.BT_HallDown
		return newOrder
	}

	
	if elevator.CurrentDir == elevio.MD_Up{
		if elevator.CabOrders[3]{ //sjekker om det er en caborder i 4 etasje!
			newOrder.Floor = 3
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}
		for i := config.NUM_FLOORS-2; i > -1; i--{ //går fra 3 etasje til 1 etasje. Øverste etasje kan ikke ha opp-knapp!
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if elevator.HallOrders[i][elevio.BT_HallUp]{
				if !AnotherGoingToFloor(i,elevatorMap){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}
		for i := 1; i < config.NUM_FLOORS; i++{ //går fra 2 etasje til 4 etasje
			if elevator.HallOrders[i][elevio.BT_HallDown]{
				if !AnotherGoingToFloor(i, elevatorMap){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}
	}else{ //DEN GÅR NEDOVER!
		if elevator.CabOrders[0]{ //sjekker om det er en caborder i 1 etasje!
			newOrder.Floor = 0
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}
		for i := 1; i < config.NUM_FLOORS; i++{ //går fra 2 etasje til 4 etasje
			if elevator.CabOrders[i]{
				newOrder.Floor = i
				newOrder.ButtonType = elevio.BT_Cab
				return newOrder
			}
			if elevator.HallOrders[i][elevio.BT_HallDown]{
				if !AnotherGoingToFloor(i, elevatorMap){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}
		for i := config.NUM_FLOORS-2; i > -1; i--{ //går fra 3 etasje til 1 etasje. Øverste etasje kan ikke ha opp-knapp!
			if elevator.HallOrders[i][elevio.BT_HallUp]{
				if !AnotherGoingToFloor(i,elevatorMap){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}
	}
	
	return newOrder
}

func AnotherGoingToFloor(floor int, elevatorMap map[int]*config.Elevator) bool{
	for _, elevator := range elevatorMap{
		if elevator.CurrentOrder.Floor == floor{
			return true
		}
	}
	return false
}

func Arbitrator(ch config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator){ //kjøres bare av Master. Master kan bytte underveis. Derfor må det sjekkes hver gang og den må være inni while-loopen // KJØRES SOM GOROUNTINE
	order := config.Order{}
	for{
		if elevatorMap[elevID].ElevRank == 1{
			for i, elevator := range elevatorMap{ //går gjennom heisene.
				if elevator.Active && !elevator.Stuck{
					if elevator.CurrentOrder.Floor == -1{

						//Heis har ingen current orders! Finnes det noen nye ordre?
						order = GetNewOrder(*elevator, elevatorMap, elevID, i)
						
						if order.Floor != -1{
							elevatorMap[i].CurrentOrder = order
							ch.New_current_order <- order
						}
					}
				}
			}
		}
	}
}