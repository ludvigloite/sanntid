package main

import(
    "./fsm"
    //"./elevcontroller"
    "./config"
    "./elevio"
    "./timer"
    "./orderhandler"
    "./network/peers"
    "./network/bcast"
    "./arbitrator"
    "./network"
    //"./orderhandler"
    "strconv"
    "fmt"
    "flag"
    "math/rand"
    "time"
)

//KANSKJE LAGE EN FUNKSJON SOM HELE TIDEN SJEKKER AT EN OG BARE EN HEIS ER MASTER
//KANSKJE LAGE EN FUNKSJON SOM SJEKKER AT CURRENTORDER IKKE ER TATT AV NOEN SOM ER INAKTIVE. EVT WATCHDOG.

func main(){

    elevIDPtr := flag.Int("elevID",42,"elevator ID")
    portPtr := flag.String("port","","port to connect to Simulator")

    flag.Parse()
    fmt.Println("Elevator ID: ", *elevIDPtr, " ,  Port Number: ", *portPtr)
    elevID := *elevIDPtr
    port := *portPtr

    rand.Seed(time.Now().UnixNano()) //genererer seed til randomizer.

    elevio.Init("localhost:"+port, config.NUM_FLOORS)

    elevatorMap := make(map[int]*config.Elevator)
    //activeElevators := make(map[int]bool) //activeElevators[elevID] = false/true)
    //kanskje må disse initialiseres til at alle er unactive. De vil få riktig konfig med en gang de får første beskjeden om hvilke som er active.

    elevator := config.Elevator{
        Active: false,
        ElevID: elevID,
        ElevRank: -1, //Dette fikses ved at man sjekker hvor mange heiser som er online.
        CurrentOrder: config.Order{Floor:-1, ButtonType:-1}, //usikker på om denne initialiseringen funker.
        CurrentFloor: -1,
        CurrentDir: elevio.MD_Down,
        CurrentState: config.IDLE,
        CabOrders: [config.NUM_FLOORS]bool{},
        HallOrders: [config.NUM_FLOORS][config.NUM_HALLBUTTONS]bool{},
    }

    firstElevator := elevator
    firstElevator.ElevID = 1
    elevatorMap[1] = &firstElevator
    secondElevator := elevator
    secondElevator.ElevID = 2
    elevatorMap[2] = &secondElevator
    thirdElevator := elevator
    thirdElevator.ElevID = 3
    elevatorMap[3] = &thirdElevator

    elevatorMap[elevID] = &elevator
    


    fsmChannels := config.FSMChannels{
        Drv_buttons: make(chan elevio.ButtonEvent), 
        Drv_floors: make(chan int),  
        Open_door: make(chan bool), 
        Close_door: make(chan bool),
        LightUpdateCh: make(chan bool),
        New_state: make(chan config.Elevator),
        New_current_order: make(chan config.Order),
        Stopping_at_floor: make(chan int),
    }

    networkChannels := config.NetworkChannels{
        PeerUpdateCh : make(chan peers.PeerUpdate),
        PeerTxEnable : make(chan bool),
        TransmittOrderCh: make(chan config.Order),
        ReceiveOrderCh: make(chan config.Order),
        TransmittElevStateCh: make(chan config.Elevator),
        ReceiveElevStateCh: make(chan config.Elevator),
        TransmittCurrentOrderCh: make(chan config.Order),
        ReceiveCurrentOrderCh: make(chan config.Order),
        TransmittCabOrderBackupCh: make(chan map[string][config.NUM_FLOORS]bool),
        ReceiveCabOrderBackupCh:make(chan map[string][config.NUM_FLOORS]bool),
    }

    go peers.Transmitter(config.SERVER_PORT, strconv.Itoa(elevID), networkChannels.PeerTxEnable)
    go peers.Receiver(config.SERVER_PORT, networkChannels.PeerUpdateCh)

    go bcast.Transmitter(config.BROADCAST_ORDER_PORT, networkChannels.TransmittOrderCh)
    go bcast.Receiver(config.BROADCAST_ORDER_PORT, networkChannels.ReceiveOrderCh)

    go bcast.Transmitter(config.BROADCAST_ELEV_STATE_PORT, networkChannels.TransmittElevStateCh)
    go bcast.Receiver(config.BROADCAST_ELEV_STATE_PORT, networkChannels.ReceiveElevStateCh)

    go bcast.Transmitter(config.BROADCAST_CURRENT_ORDER_PORT, networkChannels.TransmittCurrentOrderCh)
    go bcast.Receiver(config.BROADCAST_CURRENT_ORDER_PORT, networkChannels.ReceiveCurrentOrderCh)

    go bcast.Transmitter(config.BROADCAST_CAB_BACKUP_PORT, networkChannels.TransmittCabOrderBackupCh)
    go bcast.Receiver(config.BROADCAST_CAB_BACKUP_PORT, networkChannels.ReceiveCabOrderBackupCh)


    go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go timer.DoorTimer(fsmChannels.Close_door,fsmChannels.Open_door,config.DOOR_OPEN_TIME) //Legg true på open_door når dør skal åpnes //skrives true til close_door når tiden er ute
    go orderhandler.LightUpdater(fsmChannels.LightUpdateCh, elevatorMap, elevID)

    go network.Sender(fsmChannels, networkChannels, elevID, elevatorMap)
    go network.Receiver(networkChannels,fsmChannels, elevID, elevatorMap)
    go arbitrator.Arbitrator(fsmChannels, elevID, elevatorMap)

    //go elevcontroller.PrintElevators_withTime(elevatorMap, config.SEND_ELEV_CYCLE)


    fsm.RunElevator(fsmChannels, elevID, elevatorMap, &elevator) //kjøre som go?
}
