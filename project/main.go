package main

import(
    "./elevio"
    "fmt"
)


func main(){

    numFloors := 4

    elevio.Init("localhost:15657", numFloors)

    var d elevio.MotorDirection = elevio.MD_Up
    //elevio.SetMotorDirection(d)

    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)


    for {
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            //Endringer
            b := <- drv_floors
            if a.Floor < b {
              elevio.SetMotorDirection(elevio.MD_Down)
            }else if a.Floor >  b {
              elevio.SetMotorDirection(elevio.MD_Up)
            }else{
              elevio.SetMotorDirection(elevio.MD_Stop)
              elevio.SetButtonLamp(a.Button, a.Floor, false)
              //Fjern andre lamper i samme etasje
            }

        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)
            if a == numFloors-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)
            //Endringer
            elevio.SetFloorIndicator(a)
            /*if ((elevio.getButton(elevio.BT_HallUP, a) && d == elevio.MD_Up) || (elevio.getButton(elevio.BT_HallDown, a) && d == elevio.MD_Down) || elevio.getButton(elevio.BT_Cab, a)) {
              d = elevio.MD_Stop
            }*/

        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }
    }
}
