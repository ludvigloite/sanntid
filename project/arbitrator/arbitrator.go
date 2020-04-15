package arbitrator

import(
	"time"
	"fmt"


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
		//elevatorMap[elevID].Active = true
		if elevatorMap[elevID].ElevRank == 1{
			for i, elevator := range elevatorMap{ //går gjennom heisene.
				if elevator.Active && !elevator.Stuck{
					if elevator.CurrentOrder.Floor == -1{
						//fmt.Println("LEDIG HEIS")

						//Heis har ingen current orders! Finnes det noen nye ordre?
						time.Sleep(10 * time.Millisecond) //pass på at alle RemoveOrders har blitt synca
						order = GetNewOrder(*elevator, elevatorMap, elevID, i)
						
						if order.Floor != -1{
							fmt.Println()
							fmt.Println("JEG HAR GITT NY CurrentOrder! elev: ",i, " floor: ", order.Floor)
							fmt.Println()
							
							elevatorMap[i].CurrentOrder = order
							ch.New_current_order <- order
							//time.Sleep(10 * time.Millisecond)
						}
					}
				}
			}
		}
	}
}
/*
func Arbitrator2(ch config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator){ //kjøres bare av Master. Master kan bytte underveis. Derfor må det sjekkes hver gang og den må være inni while-loopen // KJØRES SOM GOROUNTINE
	order := config.Order{}
	for{
		//elevatorMap[elevID].Active = true
		readyElevators := [config.NUM_ELEVATORS+1] bool {} //0 element blir ikke brukt
		if elevatorMap[elevID].ElevRank == 1{
			for i, elevator := range elevatorMap{ //går gjennom heisene.
				if elevator.Active && !elevator.Stuck{
					if elevator.CurrentOrder.Floor == -1{
						readyElevators[i] = true
					}
				}
			}
			order = FindNewOrder(elevatorMap, elevID)
			if order.Floor != -1{
				bestElevID := GetBestElevator(order.Floor, elevatorMap)
			} 
		}
	}
}

func GetBestElevator(floor int, elevatorMap map[int]*config.Elevator, readyElevators [config.NUM_ELEVATORS] bool) int{
	minDist := 10
	minElev := -1

	for elevID := 1; elevID < config.NUM_ELEVATORS; elevID ++{
		if readyElevators[elevID]{
			dist := Abs(elevatorMap.CurrentFloor - floor)
			if dist < minDist{
				minDist = dist
				minElev = elevID
			}
		}
	}
	return minElev
}

func FindNewOrder(elevatorMap map[int]*config.Elevator, elevID int) config.Order{
	newOrder := config.Order{
		Sender_elev_ID: elevID,
		Sender_elev_rank: 1,
		Floor: -1,
		ButtonType: elevio.BT_HallDown,
	}

	for flr := 0; flr < config.NUM_FLOORS; flr++{
		if elevatorMap[elevID].CabOrders[flr]{
			newOrder.Floor = flr
			newOrder.ButtonType = elevio.BT_Cab
			return newOrder
		}
		for btn := elevio.BT_HallUp; btn != elevio.BT_Cab; btn++{
			if elevatorMap[elevID].HallOrders[flr][btn]{
				newOrder.Floor = flr
				newOrder.ButtonType = btn
				return newOrder
			}
		}
	}
	return newOrder
}
*/

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
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !AnotherGoingToFloor(i,elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}
		for i := 1; i < config.NUM_FLOORS; i++{ //går fra 2 etasje til 4 etasje
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !AnotherGoingToFloor(i, elevatorMap,currentElevID){
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
			if masterElev.HallOrders[i][elevio.BT_HallDown]{
				if !AnotherGoingToFloor(i, elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallDown
					return newOrder
				}
			}
		}
		for i := config.NUM_FLOORS-2; i > -1; i--{ //går fra 3 etasje til 1 etasje. Øverste etasje kan ikke ha opp-knapp!
			if masterElev.HallOrders[i][elevio.BT_HallUp]{
				if !AnotherGoingToFloor(i,elevatorMap,currentElevID){
					newOrder.Floor = i
					newOrder.ButtonType = elevio.BT_HallUp
					return newOrder
				}
			}
		}
	}
	
	return newOrder
}

func AnotherGoingToFloor(floor int, elevatorMap map[int]*config.Elevator, elevID int) bool{
	for _, elevator := range elevatorMap{
		//fmt.Println("ElevID: ",elevator.ElevID, " currentOrder: ", elevator.CurrentOrder.Floor)
		if elevID != elevator.ElevID{
			if elevator.CurrentOrder.Floor == floor || (elevator.CurrentFloor == floor && elevator.CurrentFsmState != config.ACTIVE && elevator.Active){
				//fmt.Println("elevID: ", elevator.ElevID, " currentFloor: ",elevator.CurrentFloor)
				return true
			}
		}
	}
	return false
}