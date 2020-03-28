package arbitrator

import(
	"fmt"
	"../config"
	"../elevio"
)

func GetNewOrder(elevatorMap map[int]*config.Elevator){

}

func Arbitrator(ch config.NetworkChannels, elevID int, elevatorMap map[int]*config.Elevator){ //kjøres bare av Master. Master kan bytte underveis. Derfor må det sjekkes hver gang og den må være inni while-loopen // KJØRES SOM GOROUNTINE
	order := config.Order{}
	for{
		if *elevatorMap[elevID].ElevRank == 1{
			for i, elevator in range *elevatorMap{ //går gjennom heisene. USIKKER PÅ OM DETTE FUNKER..
				
				if elevator.CurrentOrder.Floor == -1{
					//Heis har ingen current orders! Finnes det noen nye ordre?
					newOrder := 

					LAG NY GetNewOrder()!! DENNE MÅ TA HENSYN TIL AT ORDRE KAN VÆRE TATT AV NOEN ANDRE. MÅ ALTSÅ SJEKKE AT INGEN ANDRE HAR DEN SOM CURRENT ORDER.


					newOrder := orderhandler.GetNewOrder(elevator.CurrentFloor, i) //antar at hall_order_list er oppdatert!
					
					if newOrder.Floor != -1{
						orderhandler.AddNewCurrentOrder(elevator.ElevID, newOrder)
						orderhandler.AddOrder(newOrder.Floor, newOrder.ButtonType, elevator.ElevID)

						if 	orderhandler.GetElevID() == i+1{		//Master skal gi ordre til seg selv.
							fmt.Println("Jeg som master tilegner en ordre til meg selv")

							orderhandler.SetCurrentOrder(newOrder.Floor)

						} else{
							fmt.Println("Jeg som master har en ordre til heis nr ",i+1)

							Msg.New_current_order_to_who = i+1
							order.Floor = newOrder.Floor
							order.ButtonType = newOrder.ButtonType
							Msg.New_order = order
							ch.TransmitterCh <- Msg
						}
					}

				}
			}
		}

	}

}