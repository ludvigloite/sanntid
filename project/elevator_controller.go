package elevController

import(
	"fmt"

)

//ha timer med her??

func Initialize(){

	InitializeElevator()
	InitializeLights()

	//Gjør det som main starter med.

}

func InitializeElevator(){
	//kjør ned til etasjen under etasje
	//når treffer floor. Sett floor.
}

func InitializeLights(numFloors int, numHallButtons int){ //NB: Endra her navn til numHallButtons
	//Slår av lyset på alle lys
	elevio.SetDoorOpenLamp(false)
	for i := 0; i < numFloors; i++{
		elevio.SetButtonLamp(2, i, false)
		if i != 0{ //er ikke i første etasje -> kan endre på alle ned_lys 
			elevio.SetButtonLamp(1,i,false)
		}
		if i != numFloors{ //er ikke i 4 etasje -> kan endre på alle opp_lys
			elevio.SetButtonLamp(0,i,false)
		}
	}

}

func OpenDoor(){

}
