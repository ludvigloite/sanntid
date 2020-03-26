package main

import(
    "./fsm"
    "./elevcontroller"
    "./config"
    "./elevio"
    "./timer"
    "./orderhandler"
    "./network/peers"
    "./network/bcast"
    //"./orderhandler"
    "strconv"
    "fmt"
    "flag"
    "math/rand"
    "time"
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

    //Eventuelt bør det funke å ha portnummer som elevID?

    elevIDPtr := flag.Int("elevID",42,"elevator ID")
    portPtr := flag.String("port","","port to connect to Simulator")

    flag.Parse()
    fmt.Println("Elevator ID: ", *elevIDPtr, " ,  Port Number: ", *portPtr)
    elevID := *elevIDPtr
    port := *portPtr

    rand.Seed(time.Now().UnixNano()) //genererer seed til randomizer.

    elevcontroller.Initialize(elevID, "localhost:"+port)

    fsmChannels := config.FSMChannels{
        Drv_buttons: make(chan elevio.ButtonEvent), 
        Drv_floors: make(chan int),  
        Open_door: make(chan bool), 
        Close_door: make(chan bool),
        LightUpdateCh: make(chan bool),
    }

    networkChannels := config.NetworkChannels{
        PeerUpdateCh : make(chan peers.PeerUpdate),
        PeerTxEnable : make(chan bool),
        TransmitterCh : make(chan config.Packet),
        ReceiverCh : make(chan config.Packet),
    }

    go peers.Transmitter(config.SERVER_PORT, strconv.Itoa(elevID), networkChannels.PeerTxEnable)
    go peers.Receiver(config.SERVER_PORT, networkChannels.PeerUpdateCh)

    go bcast.Transmitter(config.BROADCAST_PORT, networkChannels.TransmitterCh)
    go bcast.Receiver(config.BROADCAST_PORT, networkChannels.ReceiverCh)


    go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go timer.DoorTimer(fsmChannels.Close_door,fsmChannels.Open_door,config.DOOR_OPEN_TIME) //Legg true på open_door når dør skal åpnes //skrives true til close_door når tiden er ute
    go elevcontroller.CheckAndAddOrder(fsmChannels,networkChannels)
    go orderhandler.LightUpdater(fsmChannels.LightUpdateCh)

    //go elevcontroller.SendMsg(networkChannels.TransmitterCh)
    go elevcontroller.TestReceiver(networkChannels)
    go elevcontroller.Arbitrator(networkChannels)

    fsm.RunElevator(fsmChannels) //kjøre som go?
}
