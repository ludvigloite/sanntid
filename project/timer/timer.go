// This module handles all timer-related tasks.

package timer

import(
  "time"
  "fmt"


  "../config"
)

//Put true on open_door when door should be opened. It will be written true to close_door when time is up.
func DoorTimer(finished chan<- bool, start <-chan bool, doorOpenTime time.Duration) {

	doorTimer := time.NewTimer(doorOpenTime)

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

//This function is ran only once. It exist so that an elevator that only had network trouble should not get cab orders from other elevators. 
//Backed up cab orders are only given when the elevator must be restarted.
func HasBeenDownTimer(elevatorMap map[int]*config.Elevator, elevID int, hasBeenDownBufferTime time.Duration) {

	time.Sleep(hasBeenDownBufferTime)
	elevatorMap[elevID].HasRecentlyBeenDown = false
}

//Handles motor failure error.
func WatchDogTimer(fsmCh config.FSMChannels, elevatorMap map[int]*config.Elevator, elevID int, watchDogTime time.Duration) {
	WatchDogTimer := time.NewTimer(watchDogTime)

	if !WatchDogTimer.Stop() && elevatorMap[elevID].CurrentOrder.Floor != -1 && elevatorMap[elevID].CurrentFsmState == config.ACTIVE {
		<-WatchDogTimer.C
	}

	go func(){ //update WatchDog every second as long as it is IDLE or DOOR_OPEN
		for{
			if elevatorMap[elevID].CurrentFsmState != config.ACTIVE{
				WatchDogTimer.Stop()
				WatchDogTimer.Reset(watchDogTime)
				time.Sleep(time.Second)
			}
		}
	 }()

	for {
		select {
		case <-fsmCh.Watchdog_updater:
			WatchDogTimer.Stop()
			WatchDogTimer.Reset(watchDogTime)

		case <-WatchDogTimer.C:
			fmt.Println("WatchDog Released")
			elevatorMap[elevID].CurrentOrder.Floor = -1
			elevatorMap[elevID].Stuck = true
			fsmCh.New_state <- *elevatorMap[elevID]			
		}
	}
}