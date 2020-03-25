package elevcontroller

import(
	"../elevio"
	"../orderhandler"
	"../config"
	"fmt"
	"time"
)


func Initialize(elevID int, localhost string){
    elevio.Init(localhost, config.NUM_FLOORS) //"localhost:15657"
	InitializeLights(config.NUM_FLOORS)
	orderhandler.SetElevatorID(elevID)
	orderhandler.InitQueues()

    //Wipe alle ordre til nå??
}

func CheckAndAddOrder(ch config.FSMChannels){
	for{
		select{
			case order := <- ch.Drv_buttons: //Fått inn knappetrykk
				fmt.Println("Knapp er trykket ", int(order.Button), order.Floor)
				orderhandler.AddOrder(order.Floor, int(order.Button),0) //0 fordi det bare skal legges til ordre. Ingen har tatt den enda.
				//orderhandler.UpdateLights()
				ch.LightUpdateCh <- true
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

func TestReceiver(ch config.NetworkChannels){
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
			
			orderhandler.UpdateLights()

			//UPDATE BARE LIGHTS OM DET ER EN ENDRING. 

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