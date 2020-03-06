//funksjons -& variablenavn må oppdateres etter hva som blir benyttet i de andre modulene

package timer

import(
  "time"
)

//var timerFlag int

//func timer_start(){
//  timerFlag = 1
//}

func DoorTimer(finished chan<- bool, start <-chan bool, doorOpenTime time.Duration) { //må kjøres som goroutine

	doorTimer := time.NewTimer(doorOpenTime)

  //empty the channel -> not concurrent receivers
	if !doorTimer.Stop() {
		<-doorTimer.C
	}

	for {
		select {
		case <-start:
			doorTimer.Reset(doorOpenTime)
		case <-doorTimer.C:
			finished <- true
		}
	}
}



/*

func StartDoorTimer(finished chan<- bool, doorOpenTime time.Duration) {

	doorTimer := time.NewTicker(3 * time.Second)

	for { //evig while. Funksjonen må kjøres som en goroutine
		select {
		case <-doorTimer.C: //ikke gå ut av evig while før doortimer gir ut C, som betyr at tiden er ute.
			finished <- true
			return
		}
	}
}




func hasOrders(globalState GlobalElevator) bool {
	for f := range globalState.HallRequests { //floor
		for b := range globalState.HallRequests[f] { //button
			if globalState.HallRequests[f][b] {
				return true
			}
		}
	}
	return false
}


func MotorTimer(timeout chan<- bool, globalState <-chan GlobalElevator, motorMotionTimer time.Duration) {

    floorMap := make(map[string]int)
    motorTimerEnabled:= false
  	motorTimer := time.NewTimer(timeout)

  	for {
  		select {
  		case newGlobalState := <-globalState:
  			//motortimer is enabled when there exists hall orders
  			motorTimerEnabled = hasOrders(newGlobalState)

  			// Reset timer if an elevator has reached a new floor
  			for newElevID, newElev := range newGlobalState.Elevators {
  				if floor, ok := floorMap[newElevID]; ok {
  					if floor != newElev.Floor {
  						if motorTimer.Stop() {
  							motorTimer.Reset(timeout)
  						}
  					}
  				}
  				floorMap[newElevID] = newElev.Floor
  			}

  	//motortimer timed out
    case <- motorTimer:
  			timeout <- true
  			mototTimer.Reset(timeout)

  		default:
  			if !motorTimerEnabled && motorTimer.Stop() {
  				motorTimer.Reset(timeout)
  			}
  		}
  	}
  }

// Variabler mest sannsynlig nødvendig i FSM:
  //doorTimerFinished <-chan bool,
	//doorTimerStart chan<- bool,
  //motorTimeout <-chan bool,
	//motorUpdateState chan<- GlobalElevator

*/