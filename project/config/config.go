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
	DOOR_OPEN_TIME 		= 3 * time.Second 
	WATCHDOG_TIMEOUT	= 5 * time.Second
)

const(
	SERVER_PORT 		= 15647//15371 //ENDRES
	BROADCAST_PORT		= 16569//16732 //ENDRES
	BROADCAST_INTERVAL 	= 200 * time.Millisecond


)

type Packet struct {
	Message 			string
	Iter    			int
	ID                	int
	Timestamp         	int
	Error_id          	int
	State             	int
	Current_order     	int
	Message_nr        	int
	Order_list        	[NUM_FLOORS][NUM_HALLBUTTONS] int
	Confirmed_orders  	[3][4]int
}

type FSMChannels struct {
	Drv_buttons 		chan elevio.ButtonEvent
    Drv_floors  		chan int
    Open_door			chan bool
    Close_door			chan bool
}

type NetworkChannels struct{
	PeerUpdateCh 		chan peers.PeerUpdate
	PeerTxEnable 		chan bool
	TransmitterCh 		chan Packet
	ReceiverCh 			chan Packet
}