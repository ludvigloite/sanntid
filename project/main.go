package main

import(
    "./fsm"
    "./elevcontroller"
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
    go fsm.RunElevator()
}
