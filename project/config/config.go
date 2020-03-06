package config


import( 
	"time"
	"../elevio"
)

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

type Channels struct {
	drv_buttons 		chan elevio.ButtonEvent
    drv_floors  		chan int
    open_door			chan bool
    close_door			chan bool
}