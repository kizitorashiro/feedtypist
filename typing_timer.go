package main

import (
	"time"
	//"fmt"
)

func sendElapsedTime(startTime time.Time, timerCh chan int) {
	currentTime := time.Now()
	elapsedTime := (int)((currentTime.Sub(startTime)).Seconds())
	timerCh <- elapsedTime
}

func timerLoop(timerCh chan int, timerOpCh chan int) {
	timerEnable := false
	startTime := time.Now()

timerloop:
	for {
		t := time.NewTimer(time.Second)
		defer t.Stop()

		select {
		case <-t.C:
			if timerEnable {
				sendElapsedTime(startTime, timerCh)
			}
		case opCode, ok := <-timerOpCh:
			if !ok {
				break timerloop
			}
			switch opCode {
			case 0:
				// start timer
				timerEnable = true
				startTime = time.Now()
			case 1:
				// stop timer
				timerEnable = false
				sendElapsedTime(startTime, timerCh)
			}
		}
	}
	/*
	   timerloop:
	   	for {
	   		//fmt.Println("tiemr loop")
	   		if len(timeOpCh) > 0 {
	   			//fmt.Println("timer receive")
	   			opCode, ok := <-timeOpCh
	   			if !ok {
	   				break timerloop
	   			}
	   			switch opCode {
	   			case 0:
	   				// start timer
	   				timerEnable = true
	   				startTime = time.Now()
	   			case 1:
	   				// stop timer
	   				timerEnable = false
	   				sendElapsedTime(startTime, timerCh)
	   			}
	   		}
	   		if timerEnable {
	   			sendElapsedTime(startTime, timerCh)
	   		}
	   		time.Sleep(500 * time.Millisecond)
	   	}
	*/
}
