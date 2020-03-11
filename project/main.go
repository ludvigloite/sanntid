package main

import(
    "./fsm"
    "./elevcontroller"
    "./config"
    "./elevio"
    "./timer"
    "./network/peers"
    "./network/bcast"
    //"./orderhandler"

    "fmt"
    "flag"
)

////////////////////////////////////////////
//              TODO
//  -Nye ordre kommer ikke inn når døra er åpen
//  -Bruke goroutines og channels
//  -Fikse nettverk
//
//
//
//
//
///////////////////////////////////////////

func main(){

    elevIDPtr := flag.Int("elevID",1,"an int")

    flag.Parse()
    fmt.Println("Elevator ID: ", *elevIDPtr)
    elevID := *elevIDPtr

    elevcontroller.Initialize(elevID)

    fsmChannels := config.FSMChannels{
    Drv_buttons: make(chan elevio.ButtonEvent), 
    Drv_floors: make(chan int),  
    Open_door: make(chan bool), 
    Close_door: make(chan bool),
	}

	networkChannels := config.NetworkChannels{
		PeerUpdateCh : make(chan peers.PeerUpdate),
		PeerTxEnable : make(chan bool),
		TransmitterCh : make(chan config.Packet),
		ReceiverCh : make(chan config.Packet),
	}

	go peers.Transmitter(config.SERVER_PORT, "hei22", networkChannels.PeerTxEnable)
	go peers.Receiver(config.SERVER_PORT, networkChannels.PeerUpdateCh)

	go bcast.Transmitter(config.BROADCAST_PORT, networkChannels.TransmitterCh)
	go bcast.Receiver(config.BROADCAST_PORT, networkChannels.ReceiverCh)


	go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go timer.DoorTimer(fsmChannels.Close_door,fsmChannels.Open_door,config.DOOR_OPEN_TIME)
    go elevcontroller.CheckAndAddOrder(fsmChannels.Drv_buttons) //Legg true på open_door når dør skal åpnes //skrives true til close_door når tiden er ute

    go elevcontroller.SendMsg(networkChannels.TransmitterCh)
    go elevcontroller.TestReceiver(networkChannels)

    fsm.RunElevator(fsmChannels) //kjøre som go?
}
