package config


import( 
	"time"
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

type Channels struct{
}