package arbitrator

import(
	"fmt"
	"../config"
	"../elevio"
	//"time"
	"../elevcontroller"
)

func GetNewOrder(elevatorMap map[int]*config.Elevator) config.Order{
	newOrder := config.Order{}
	for i := 0; i < config.NUM_FLOORS; i++{

	}
	newOrder.Floor = -1
	newOrder.ButtonType = elevio.BT_HallDown
	return newOrder
}

func Arbitrator(ch config.FSMChannels, elevID int, elevatorMap map[int]*config.Elevator){ //kjøres bare av Master. Master kan bytte underveis. Derfor må det sjekkes hver gang og den må være inni while-loopen // KJØRES SOM GOROUNTINE
	order := config.Order{}
	for{

		if elevatorMap[elevID].ElevRank == 1{
			for i, elevator := range elevatorMap{ //går gjennom heisene. USIKKER PÅ OM DETTE FUNKER..
				//fmt.Println("Sjekker heis nr", i)
				elevcontroller.PrintElevator(*elevator)

				
				if elevator.CurrentOrder.Floor == -1{
					//Heis har ingen current orders! Finnes det noen nye ordre?
					order = GetNewOrder(elevatorMap)

					//LAG NY GetNewOrder()!! DENNE MÅ TA HENSYN TIL AT ORDRE KAN VÆRE TATT AV NOEN ANDRE. MÅ ALTSÅ SJEKKE AT INGEN ANDRE HAR DEN SOM CURRENT ORDER.
					//newOrder := orderhandler.GetNewOrder(elevator.CurrentFloor, i) //antar at hall_order_list er oppdatert!
					
					if order.Floor != -1{
						//orderhandler.AddNewCurrentOrder(elevator.ElevID, newOrder)
						//orderhandler.AddOrder(newOrder.Floor, newOrder.ButtonType, elevator.ElevID)

						//Det finnes en ordre!!
						order.Sender_elev_ID = elevID
						order.Sender_elev_rank = elevatorMap[elevID].ElevRank
						order.Receiver_elev = 1

						elevatorMap[i].CurrentOrder = order
						
						fmt.Println("Ny currentOrder til heis nr ",i)

						ch.New_current_order <- order

						
					}

				}
			}
		}

	}

}