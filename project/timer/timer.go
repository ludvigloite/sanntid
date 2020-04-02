package timer

import(
  "time"
)



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
