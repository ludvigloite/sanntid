package main

import(
    "./fsm"
    "./elevcontroller"
    "./config"
    "./elevio"
    "./timer"
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

    fsmChannels := config.FSMChannels{
    Drv_buttons: make(chan elevio.ButtonEvent), 
    Drv_floors: make(chan int),  
    Open_door: make(chan bool), 
    Close_door: make(chan bool),
	}

	go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go timer.DoorTimer(fsmChannels.Close_door,fsmChannels.Open_door,config.DOOR_OPEN_TIME)
    go elevcontroller.CheckAndAddOrder(fsmChannels.Drv_buttons)
    //Legg true på open_door når dør skal åpnes
    //skrives true til close_door når tiden er ute

    fsm.RunElevator(fsmChannels) //kjøre som go?

    for{}
}
