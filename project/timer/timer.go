package timer

/*
import(
  "time"
)

var timerFlag int

func timer_start(){
  timerFlag = 1
}

func MotorTimer(finished chan<- bool, start <-chan bool, motorMotionTimer time.Duration)

  motorTimer := time.NewTimer(4 * time.Second)



  for {
    select {
    case <-start:
      motorTimer.Reset(motorMotionTimer)
    case <-motorTimer.:
      finished <- true
    }
  }
}


func DoorTimer(finished chan<- bool, start <-chan bool, doorOpenTime time.Duration) {

	doorTimer := time.NewTimer(3 * time.Second)

	if !doorTimer.Stop() {
		<-doorTimer.C
	}

	for {
		select {
		case <-start:
			doorTimer.Reset(doorOpenTime)
		case <-doorTimer.:
			finished <- true
		}
	}
}

motorTimedOut.Stop()


*/