package config


import( 
	"time"
	"../elevio"
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
	SERVER_PORT 		= 15657 //ENDRES
	BROADCAST_PORT		= 15898 //ENDRES
	BROADCAST_INTERVAL 	= 200 * time.Millisecond


)

type FSMChannels struct {
	Drv_buttons 		chan elevio.ButtonEvent
    Drv_floors  		chan int
    Open_door			chan bool
    Close_door			chan bool
}

type NETWORKChannels struct{

}