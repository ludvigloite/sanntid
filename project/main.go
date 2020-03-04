package main

import(
    "./fsm"
    "../elevcontroller"
    "fmt"
)


func main(){
    elevcontroller.Initialize()
    fsm.RunElevator()
}
