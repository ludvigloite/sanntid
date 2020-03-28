package config


import( 
	"time"
	"../elevio"
	"../network/peers"
)

/* DETTE KAN IKKE HARDKODES
const(
	ELEV_ID				= const(os.Args[1])
	ELEV_RANK			= 1
)
*/

const(
	NUM_FLOORS 			= 4
	NUM_HALLBUTTONS 	= 2
	NUM_ELEVATORS		= 3
	DOOR_OPEN_TIME 		= 3 * time.Second 
	WATCHDOG_TIMEOUT	= 5 * time.Second
)

const(
	SERVER_PORT 				= 12346//15647//15371 //ENDRES
	BROADCAST_ORDER_PORT		= 12347//16732 //ENDRES
	BROADCAST_CURRENT_ORDER_PORT= 12348
	BROADCAST_ELEV_STATE_PORT	= 12349
	BROADCAST_INTERVAL 			= 200 * time.Millisecond

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
	ButtonType 			config.ButtonType
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

/*
type Packet struct {
	Elev_ID                		int
	Elev_rank 					int
	New_order 					Order
	New_current_order_to_who 	int
	Timestamp         			int
	Error_id          			int
	State             			int //0:Idle, 1: Active, 2: Door_open, 3: UNDEFINED
	Current_order     			int
	Message_nr        			int
	Order_list        			[NUM_FLOORS][NUM_HALLBUTTONS] int
	Confirmed_orders  			[3][4]int
	Rank 						int //bytter underveis
	CurrentFloor 				int //hvilken etasje er heisen i nå. 0 , 1 , 2 , 3
	CurrentDir 					int //hvilken retning har heisen. -1 , 0 , 1. Kun 0 i spesielle tilfeller. Er -1 / 1 også når den stopper i et floor. Den skal jo tross alt videre i samme retning.
}*/

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
	PeerUpdatAeCh 			chan peers.PeerUpdate
	PeerTxEnable 			chan bool
	TransmittOrderCh 		chan Order
	ReceiveOrderCh 			chan Order
	TransmittElevStateCh 	chan Elevator
	ReceiveElevStateCh 		chan Elevator
	TransmittCurrentOrderCh	chan Order
	ReceiveCurrentOrderCh	chan Order
}
