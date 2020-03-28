package elevcontroller

import(
	"../elevio"
	"../orderhandler"
	"../config"
	"fmt"
	"time"
	"math/rand"
)


func Initialize(elevator *config.Elevator){
	elevio.SetMotorDirection(elevio.MD_Down)
	elevcontroller.ResetLights()	
	InitQueues(&elevator)
}

func InitQueues(elevator *config.Elevator){
	for i := 0;i < config.NUM_FLOORS; i++ {
		*elevator.CabOrders[i] = false //INIT CABORDERS

		for j := 0;j< config.NUM_HALLBUTTONS; j++{
			*elevator.HallOrders[i][j] = false //INIT HALLORDERS
		}
	}
}

func CheckAndAddOrder(fsmCh config.FSMChannels, netCh config.NetworkChannels){
	//KJØRES SOM GOROUNTINE
	order := config.Order{}
	Msg := config.Packet{}
	for{
		select{
			case buttonpress := <- fsmCh.Drv_buttons: //Fått inn knappetrykk
				fmt.Println("Knapp er trykket! ", int(buttonpress.Button), buttonpress.Floor)
				//orderhandler.AddOrder(buttonpress.Floor, int(buttonpress.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.

				//Msg.Order_list = orderhandler.GetHallOrderQueue()
				Msg.Elev_ID = orderhandler.GetElevID()
				Msg.Elev_rank = orderhandler.GetElevRank()
				Msg.CurrentFloor = orderhandler.GetCurrentFloor()
				//Msg.CurrentOrder = orderhandler.GetCurrentOrder()
				//Msg.State = orderhandler.GetCurrentState()
				Msg.New_current_order_to_who = -1

				order.Floor = buttonpress.Floor
				order.ButtonType = int(buttonpress.Button)
				order.Packet_id = rand.Intn(10000)
				order.Type_action = 1 //Det er en ordre som skal legges til
				order.Approved = false

				Msg.New_order = order

				if orderhandler.IsMaster(){
					netCh.TransmitterCh <- Msg
				}else{ //en slave har fått inn en ordre
					netCh.TransmitterCh <- Msg
				}
				fmt.Println("Har nå sendt avgårde pakke om at knapp er trykket!")


			case received_order := <- netCh.ReceiverCh:

				if received_order.Elev_ID == orderhandler.GetElevID() || (received_order.Elev_rank!=1 && orderhandler.GetElevRank()!=1){		//hvis det er fra deg selv eller du er slave og pakken er fra slave: IGNORE!
					break //Drit i ordre fra deg selv
				}
				MsgOrder := received_order.New_order



				if received_order.New_current_order_to_who == orderhandler.GetElevID(){ //Heisen har fått en ny current order fra Master.
					orderhandler.SetCurrentOrder(MsgOrder.Floor)
					fmt.Println("Jeg har fått beskjed av Master om å utføre ordre i etasje ", MsgOrder.Floor)
					break
				}

				orderhandler.UpdateCurrentFloor(received_order.Elev_ID, received_order.CurrentFloor)

				if MsgOrder.Type_action == -1 && orderhandler.IsMaster() && !MsgOrder.Approved{ //Master må eventuelt fjerne currentOrder for heisen som sender inn slettet ordre.
					if MsgOrder.Floor == orderhandler.GetElevList()[received_order.Elev_ID-1].CurrentOrder.Floor{
						//heisen har nettopp utført sin currentOrder
						deletedOrder := config.Order{}
						deletedOrder.Floor = -1
						deletedOrder.ButtonType = -1

						orderhandler.AddNewCurrentOrder(received_order.Elev_ID,deletedOrder) //sletter heisens currentOrder
					}
				}


				if received_order.New_order.Approved{
					//legg til i ordrekø! evt fjern fra ordrekø!
					if orderhandler.IsMaster(){
						received_order.Elev_ID = orderhandler.GetElevID()
						netCh.TransmitterCh <- received_order
					}

					if MsgOrder.Type_action == 1{
						orderhandler.AddOrder(MsgOrder.Floor, MsgOrder.ButtonType, 0) //elevatorID er 0 om det bare skal legges inn ordre uten at noen tar den.)
					}else{ 
						orderhandler.AddOrder(MsgOrder.Floor, MsgOrder.ButtonType, -1) //fjerner ordre
					}
					fsmCh.LightUpdateCh <- true //hvis ordrekø er endret oppdateres lysene.
					break
				}

				if orderhandler.IsMaster(){
					received_order.Elev_ID = orderhandler.GetElevID()
					netCh.TransmitterCh <- received_order

				}else{ //er slave
					received_order.Elev_ID = orderhandler.GetElevID()
					received_order.New_order.Approved = true
					netCh.TransmitterCh <- received_order

				}

		}
	}
}

func Arbitrator(ch config.NetworkChannels){ //kjøres bare av Master. Master kan bytte underveis. Derfor må det sjekkes hver gang og den må være inni while-loopen // KJØRES SOM GOROUNTINE
	elevList := orderhandler.GetElevList()
	order := config.Order{}
	Msg := config.Packet{}
	for{
		if orderhandler.IsMaster(){
			for i := 0 ; i < config.NUM_ELEVATORS; i++{ //går gjennom heisene
				elevator := elevList[i]	//HUSK AT LISTA ER 0-INDEKSERT. ELEVID 1 ER PLASSERT PÅ INDEX 0
				if elevator.CurrentOrder.Floor == -1{
					//Heis nr i+1 har ingen current orders!
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


//Kan vel kanskje i stedet bare fjerne alle ordre og så kjøre update lights??
func ResetLights(){	//Slår av lyset på alle lys
	numFloors := config.NUM_FLOORS
	elevio.SetDoorOpenLamp(false)
	for i := 0; i < numFloors; i++{
		elevio.SetButtonLamp(elevio.BT_Cab, i, false)
		if i != 0{ //er ikke i første etasje -> kan endre på alle ned_lys 
			elevio.SetButtonLamp(elevio.BT_HallDown,i,false)
		}
		if i != numFloors{ //er ikke i 4 etasje -> kan endre på alle opp_lys
			elevio.SetButtonLamp(elevio.BT_HallUp,i,false)
		}
	}

}

func PrintElevator(elevator config.Elevator){
	fmt.Println()
	fmt.Println("elevID: ",elevator.ElevID,"\t Rank: ",elevator.ElevRank)
	fmt.Println("CurrentOrder = Floor: ",elevator.CurrentOrder.Floor, "\t ButtonType: ",elevator.CurrentOrder.ButtonType)
	fmt.Println("CurrentFloor = ", elevator.CurrentFloor)
	fmt.Println("CurrentState = ", elevator.CurrentState)
}

func GetDirection(currentFloor int, destinationFloor int) elevio.MotorDirection{
	if destinationFloor == -1 || destinationFloor == currentFloor { //enten har den ikke noen retning, eller så er den på riktig floor
		return elevio.MD_Stop

	} else if currentFloor < destinationFloor { //heisen er lavere enn sin destinasjon -> kjører oppover
		return elevio.MD_Up

	} else{
		return elevio.MD_Down
	}
}

func ShouldStopAtFloor(elevator config.Elevator) bool{
	id = elevator.Elev_ID
	currentFloor = elevator.CurrentFloor
	destinationFloor = elevator.CurrentOrder.Floor
	dir = elevator.CurrentDir
	if dir == MD_Stop{ //har ingen ordre eller er på etasjen currentOrder tilsier. KAN FØRE TIL ERROR!!
		return true
	}
	if IsThereOrder(currentFloor,2,elevID){ //Det er en cab order i denne etasjen
		return true
	}
	if IsThereOrder(currentFloor,0,elevID) && dir == 1{ //retning til heis er opp og det er en ordre opp
		return true
	}
	if IsThereOrder(currentFloor,1,elevID) && dir == -1 { //retning til heis er ned og det er en ordre ned
		return true
	}
	return false
}

func IsThereOrder(floor int, buttonType config.ButtonType, elevID int) bool{ //kan kanskje bare implementeres i ShouldStopAtFloor
	if buttonType == BT_Cab{
		if cabOrderQueue.Active[floor] == 0 && cabOrderQueue.ElevID == elevID{
			return true
		}
	} else{
		if hallOrderQueue[floor][buttonType] == 0 || hallOrderQueue[floor][buttonType] == elevID{
			return true
		}
	}
	return false
}

func TestReceiver(ch config.NetworkChannels){
	fmt.Println("Har kommet inn i TestReceiver")
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			/*
		case packet := <-ch.ReceiverCh:
			//fmt.Printf("Received: %#v\n", a.Order_list)
			//fmt.Println("Mottar fra: ",a.ID," OPP \t NED") //opp ned er bare for at man skal forstå ordrekøen.
			//orderhandler.PrintHallOrderQueue(a.Order_list)
			if packet.ID == orderhandler.GetElevID(){
				break //litt usikker på hvordan break funker
			}

			if orderhandler.IsMaster(){ //Du selv er Master
				orderhandler.MergeHallQueues(packet)

			} else if packet.ID ==1 { //du mottar fra Master
				orderhandler.SetHallOrderQueue(packet.Order_list)
			}
			
			//orderhandler.UpdateLights()
			LightUpdateCh <- true

			//UPDATE BARE LIGHTS OM DET ER EN ENDRING. GJØR DETTE INNI DE FORSKJELLIGE FUNKSJONENE
*/
		}
	}
}


func SendMsg(TransmitterCh chan <- config.Packet){
	Msg := config.Packet{}
	for{
		Msg.Order_list = orderhandler.GetHallOrderQueue()
		Msg.Elev_ID = orderhandler.GetElevID()
		Msg.CurrentFloor = orderhandler.GetCurrentFloor()
		Msg.State = orderhandler.GetCurrentState()
		//fmt.Println("Sender kø:   ",orderhandler.GetHallOrderQueue())
		TransmitterCh <- Msg
		//fmt.Println(Msg)
		time.Sleep(1*time.Second)
	}
}

/*func SendMsg(TransmitterCh chan <- config.Packet, NewOrderCh ){ //send bare hvis du har fått inn en ny ordre. Ellers sender man hvert x sekund
	Msg := config.Packet{}
	for{

		Msg.Order_list = orderhandler.GetHallOrderQueue()
		Msg.ID = orderhandler.GetElevID()
		Msg.CurrentFloor = orderhandler.GetCurrentFloor()
		Msg.State = orderhandler.GetCurrentState()
		//fmt.Println("Sender kø:   ",orderhandler.GetHallOrderQueue())
		TransmitterCh <- Msg
		//fmt.Println(Msg)
		time.Sleep(1*time.Second)
	}
}*/


/*	BRUKES IKKE, MEN KANSKJE TIL TESTING SENERE?

func StopElevator(){
	elevio.SetMotorDirection(elevio.MD_Stop)
	OpenDoor(3)
}

func OpenDoor(seconds time.Duration) {
	elevio.SetDoorOpenLamp(true)
	time.Sleep(seconds * time.Second)
	elevio.SetDoorOpenLamp(false)
}
*/