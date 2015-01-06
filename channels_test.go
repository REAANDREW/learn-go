package learngo

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const (
	STOP uint16 = 0x0001
)

type Instance struct {
	State uint16
}

func (instance *Instance) Stop() {
	instance.State = STOP
}

func Test_WriteAValueToAChannelAndReadItBackAgain(t *testing.T) {
	channel := make(chan int, 1)
	channel <- 1
	val := <-channel
	fmt.Printf("So we got %d\n", val)
}

func Test_WithoutBufferingZeeChannel(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	channel := make(chan int)
	go func() {
		val := <-channel
		fmt.Printf("So we got %d\n", val)
		time.Sleep(time.Second * 1)
		wg.Done()
	}()
	channel <- 2
	wg.Wait()
	fmt.Println("ok we are now done")
}

func Test_structs(t *testing.T) {
	i := Instance{}
	fmt.Println(i.State)
	i.Stop()
	fmt.Println(i.State)
}
