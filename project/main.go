package main

import(
	"strconv"
    "fmt"
    "flag"


    "./fsm"
    "./elevcontroller"
    "./config"
    "./elevio"
    "./timer"
    "./network/peers"
    "./network/bcast"
    "./arbitrator"
    "./network"
)


func main(){

	fmt.Print("In Case of Network Shutdown: \nShow Orders: ",config.SHOW_ORDERS_WHEN_NETWORK_DOWN,"\tAdd HallOrders: ", config.ADD_HALL_ORDERS_WHEN_NETWORK_DOWN, "\n\n")
	
	/*	INIT 	*/

	// Taking arguments from command line
    elevIDPtr := 	flag.Int("elevID",42,"elevator ID")
    portPtr := 		flag.String("port","","port to connect to Simulator")
    flag.Parse()

    fmt.Println("Elevator ID: ", *elevIDPtr, " ,  Port Number: ", *portPtr)
    elevID := 	*elevIDPtr
    port := 	*portPtr

    elevio.Init("localhost:"+port, config.NUM_FLOORS)

    elevatorMap := make(map[int]*config.Elevator)

    elevator := config.Elevator{
        NetworkDown: 			true, //Must be initialized to true to handle the case that network is down from initialization
        HasRecentlyBeenDown: 	true,
        ElevID: 				elevID,
        CurrentOrder: 			config.Order{Floor:-1, ButtonType:-1}, //starting by giving no currentOrder. floor == -1 means no currentOrder.
        CurrentFsmState: 		config.IDLE,
    }

    //initializing all elevators in the local elevatorMap. These will be updated by peer updates once it is online.
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
        Drv_buttons: 		make(chan elevio.ButtonEvent), 
        Drv_floors: 		make(chan int),  
        Open_door: 			make(chan bool), 
        Close_door: 		make(chan bool),
        LightUpdateCh: 		make(chan bool),
        New_state: 			make(chan config.Elevator),
        New_current_order: 	make(chan config.Order),
        Stopping_at_floor: 	make(chan int),
        Watchdog_updater: 	make(chan bool),
    }

    networkChannels := config.NetworkChannels{
        PeerUpdateCh : 				make(chan peers.PeerUpdate),
        PeerTxEnable : 				make(chan bool),
        TransmittOrderCh: 			make(chan config.Order),
        ReceiveOrderCh: 			make(chan config.Order),
        TransmittElevStateCh: 		make(chan config.Elevator),
        ReceiveElevStateCh: 		make(chan config.Elevator),
        TransmittCurrentOrderCh: 	make(chan config.Order),
        ReceiveCurrentOrderCh: 		make(chan config.Order),
        TransmittCabOrderBackupCh: 	make(chan map[string][config.NUM_FLOORS]bool),
        ReceiveCabOrderBackupCh:	make(chan map[string][config.NUM_FLOORS]bool),
    }

    /*	INIT FINISHED	*/

    /* STARTING GOROUTINES, ALWAYS RUNNING	*/

    go peers.Transmitter(	config.SERVER_PORT, strconv.Itoa(elevID), networkChannels.PeerTxEnable)
    go peers.Receiver(		config.SERVER_PORT, networkChannels.PeerUpdateCh)

    go bcast.Transmitter(	config.BROADCAST_ORDER_PORT, networkChannels.TransmittOrderCh)
    go bcast.Receiver(		config.BROADCAST_ORDER_PORT, networkChannels.ReceiveOrderCh)

    go bcast.Transmitter(	config.BROADCAST_ELEV_STATE_PORT, networkChannels.TransmittElevStateCh)
    go bcast.Receiver(		config.BROADCAST_ELEV_STATE_PORT, networkChannels.ReceiveElevStateCh)

    go bcast.Transmitter(	config.BROADCAST_CURRENT_ORDER_PORT, networkChannels.TransmittCurrentOrderCh)
    go bcast.Receiver(		config.BROADCAST_CURRENT_ORDER_PORT, networkChannels.ReceiveCurrentOrderCh)

    go bcast.Transmitter(	config.BROADCAST_CAB_BACKUP_PORT, networkChannels.TransmittCabOrderBackupCh)
    go bcast.Receiver(		config.BROADCAST_CAB_BACKUP_PORT, networkChannels.ReceiveCabOrderBackupCh)


    go elevio.PollButtons(		fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(	fsmChannels.Drv_floors)

    go timer.DoorTimer(			fsmChannels.Close_door, fsmChannels.Open_door, 							config.DOOR_OPEN_TIME)
    go timer.WatchDogTimer(		fsmChannels, 									elevatorMap, elevID, 	config.WATCHDOG_TIME)
    go timer.HasBeenDownTimer(													elevatorMap, elevID, 	config.HAS_BEEN_DOWN_BUFFER)

	go network.Sender(	fsmChannels, networkChannels, elevID, elevatorMap)
    go network.Receiver(fsmChannels, networkChannels, elevID, elevatorMap)

    go elevcontroller.LightUpdater(	fsmChannels.LightUpdateCh, 		elevatorMap, elevID)
    go arbitrator.Arbitrator(		fsmChannels.New_current_order, 	elevatorMap, elevID)
    go arbitrator.RankSolver(		fsmChannels.New_state, 			elevatorMap, elevID)

    /* ALL GOROUTINES UP AND RUNNING */


    /* STARTING FSM */
    fsm.RunElevator(fsmChannels, &elevator)
}
