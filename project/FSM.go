package fsm

import(
	"fmt"


)

const(
	IDLE = "IDLE"
	ACTIVE = "ACTIVE"
	RESET = "RESET"
)


func RunElevator(){

	state := IDLE
	//orderHandler.currentFloor
	//orderHandler.currentDir


	drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)




	for{
		switch state{
		case UDLE: //heis er IDLE. Skal ikke gjøre noe med mindre den får knappetrykk
			select{
			case a := <- elevController.drv_buttons
			}




		case ACTIVE:
			//







		case RESET:
			//


		}
	}
}