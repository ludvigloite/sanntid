package main

import(
    "./fsm"
    "./elevcontroller"
    "./config"
    //"fmt"
)

////////////////////////////////////////////
//              TODO
//  -Nye ordre kommer ikke inn når døra er åpen
//  -Bruke goroutines og channels
//  -Fikse nettverk
//
//
//
//
//
///////////////////////////////////////////

func main(){


    elevcontroller.Initialize()

    fsmChannels := config.Channels{
    drv_buttons: make(chan elevio.ButtonEvent), 
    drv_floors: make(chan int),  
    open_door: make(chan bool), 
    close_door: make(chan bool),
	}

	go elevio.PollButtons(New_Channels.drv_buttons)
    go elevio.PollFloorSensor(New_Channels.drv_floors)
    go timer.DoorTimer(New_Channels.close_door,New_Channels.open_door,config.DOOR_OPEN_TIME)
    //Legg true på open_door når dør skal åpnes
    //skrives true til close_door når tiden er ute

    fsm.RunElevator(fsmChannels) //kjøre som go?
}
