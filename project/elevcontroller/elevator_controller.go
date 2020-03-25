package elevcontroller

import(
	"../elevio"
	"../orderhandler"
	"../config"
	"fmt"
	"time"
	"math/rand"
)


func Initialize(elevID int, localhost string){
    elevio.Init(localhost, config.NUM_FLOORS) //"localhost:15657"
	InitializeLights(config.NUM_FLOORS)
	orderhandler.SetElevatorID(elevID)
	orderhandler.InitQueues()

    //Wipe alle ordre til nå??
}

func CheckAndAddOrder(fsmCh config.FSMChannels, netCh config.NetworkChannels){
	//KJØRES SOM GOROUNTINE
	order := config.Order{}
	Msg := config.Packet{}
	for{
		select{
			case buttonpress := <- fsmCh.Drv_buttons: //Fått inn knappetrykk
				fmt.Println("Knapp er trykket ", int(buttonpress.Button), buttonpress.Floor)
				//orderhandler.AddOrder(buttonpress.Floor, int(buttonpress.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.

				Msg.Order_list = orderhandler.GetHallOrderQueue()
				Msg.ID = orderhandler.GetElevID()
				Msg.CurrentFloor = orderhandler.GetCurrentFloor()
				Msg.State = orderhandler.GetCurrentState()

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


			case received_order := <- netCh.ReceiverCh:

				if received_order.ID == orderhandler.GetElevID() || (received_order.ID!=1 && orderhandler.GetElevID()!=1){		//hvis det er fra deg selv eller du er slave og pakken er fra slave: IGNORE!
				break //litt usikker på hvordan break funker
				}

				if received_order.New_order.Approved{
					//legg til i ordrekø! evt fjern fra ordrekø!
					if received_order.New_order.Type_action == 1{
						orderhandler.AddOrder(received_order.New_order.Floor, received_order.New_order.ButtonType, 0) //elevatorID er 0 om det bare skal legges inn ordre uten at noen tar den.)
					}else{
						orderhandler.AddOrder(received_order.New_order.Floor, received_order.New_order.ButtonType, -1)
					}
					fsmCh.LightUpdateCh <- true //hvis ordrekø er endret oppdateres lysene.
					break
				}

				if orderhandler.IsMaster(){
					netCh.TransmitterCh <- received_order

				}else{ //er slave
					received_order.New_order.Approved = true
					netCh.TransmitterCh <- received_order

				}

		}
	}
}

//Kan vel kanskje i stedet bare fjerne alle ordre og så kjøre update lights??
func InitializeLights(numFloors int){ //NB: Endra her navn til numHallButtons
	//Slår av lyset på alle lys
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

func TestReceiver(ch config.NetworkChannels, LightUpdateCh chan <- bool){
	fmt.Println("Har kommet inn i TestReceiver")
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

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

		}
	}
}


func SendMsg(TransmitterCh chan <- config.Packet){
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