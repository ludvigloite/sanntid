package config


import( 
	"time"
	"../elevio"
	"../network/peers"
)

const(
	NUM_FLOORS 			= 4
	NUM_HALLBUTTONS 	= 2
	NUM_ELEVATORS		= 3
	DOOR_OPEN_TIME 		= 3 * time.Second
	SEND_ELEV_CYCLE		= 5 * time.Second
	NUM_PACKETS			= 3
	//WATCHDOG_TIMEOUT	= 5 * time.Second
)

const(
	SERVER_PORT 					= 12346
	BROADCAST_ORDER_PORT			= 12347
	BROADCAST_CURRENT_ORDER_PORT	= 12348
	BROADCAST_ELEV_STATE_PORT		= 12349
	BROADCAST_CAB_BACKUP_PORT		= 12350
	//BROADCAST_INTERVAL 				= 200 * time.Millisecond

)

type State int
const(
	IDLE State 	= 0
	ACTIVE 		= 1
	DOOR_OPEN 	= 2
	UNDEFINED 	= 3
)
/*
type Type_Action int
const{
	ADD Type_Action = 1
	REMOVE = -1
}*/

type Order struct{
	Sender_elev_ID 		int
	Sender_elev_rank 	int
	Floor 				int
	ButtonType 			elevio.ButtonType
	Should_add			bool //-1 hvis ordre skal slettes, 1 hvis ordre blir lagt til.
	Packet_id 			int
	Approved 			bool
	Receiver_elev 		int
}

/*type ButtonType int
const (
	BT_HallUp   ButtonType = 0  //////VIKTIG!!
	BT_HallDown            = 1
	BT_Cab                 = 2
)*/

type Elevator struct{
	Active bool
	ElevID int
	ElevRank int
	CurrentOrder Order
	CurrentFloor int
	CurrentDir elevio.MotorDirection
	CurrentState State
	CabOrders [NUM_FLOORS]bool
	HallOrders [NUM_FLOORS][NUM_HALLBUTTONS]bool
}

type FSMChannels struct {
	Drv_buttons 		chan elevio.ButtonEvent
    Drv_floors  		chan int
    Open_door			chan bool
    Close_door			chan bool
    LightUpdateCh		chan bool
    New_state			chan Elevator
    New_current_order 	chan Order
    Stopping_at_floor	chan int
}

type NetworkChannels struct{
	PeerUpdateCh 				chan peers.PeerUpdate
	PeerTxEnable 				chan bool
	TransmittOrderCh 			chan Order
	ReceiveOrderCh 				chan Order
	TransmittElevStateCh 		chan Elevator
	ReceiveElevStateCh 			chan Elevator
	TransmittCurrentOrderCh		chan Order
	ReceiveCurrentOrderCh		chan Order
	TransmittCabOrderBackupCh 	chan map[string][NUM_FLOORS]bool
	ReceiveCabOrderBackupCh		chan map[string][NUM_FLOORS]bool
}
