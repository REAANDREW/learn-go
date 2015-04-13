package main

import (
	"fmt"
	"time"
)

func example_1() {
	//Write to and read from a channel
	in := make(chan string)
	go func() {
		result := <-in
		fmt.Println(result)
	}()
	in <- "Example 1"
}

func example_2_write(writeOnly chan<- string) {
	writeOnly <- "Example 2"
}

func example_2() {
	//Passing a channel as writeonly
	in := make(chan string)
	go example_2_write(in)
	result := <-in
	fmt.Println(result)
}

func example_3_read(readOnly <-chan string) {
	result := <-readOnly
	fmt.Println(result)
}

func example_3() {
	//Passing a channel as readonly
	in := make(chan string)
	go example_3_read(in)
	in <- "Example 3"
}

func example_4() {
	//Block waiting for a writer
	in := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		in <- "Example 4"
	}()
	result := <-in
	fmt.Println(result)
}

func example_5() {
	//Block waiting for a reader
	in := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		result := <-in
		fmt.Println(result)
	}()
	in <- "Example 5"
}

func example_6() {
	//Buffered channels
	in := make(chan string, 1)
	in <- "Example 6.1"
	go func() {
		time.Sleep(2 * time.Second)
		in <- "Example 6.2"
	}()
	fmt.Println(<-in)
	fmt.Println(<-in)
}

func example_7() {
	//Selecting the channel to read in a non-blocking way
	chan_string := make(chan string)
	chan_int := make(chan int)
	go func() {
		for {
			select {
			case string_value := <-chan_string:
				fmt.Println(string_value)
			case int_value := <-chan_int:
				fmt.Println(fmt.Sprintf("Example %d.2", int_value))
				break
			default:
			}
		}
	}()
	time.Sleep(1 * time.Second)
	chan_string <- "Example 7.1"
	chan_int <- 7
}

func example_8() {
	//Select a channel to write in a non-blocking way
	chan_string := make(chan string)
	chan_int := make(chan int)
	go func() {
		for {
			select {
			case chan_string <- "Example 8.1":
			case chan_int <- 8:
			default:
			}
		}
	}()
	time.Sleep(1 * time.Second)
	fmt.Println(<-chan_string)
	fmt.Println(fmt.Sprintf("Example %d.2", <-chan_int))
}

func example_9() {
	//range over a channel
	chan_int := make(chan int, 8)
	for i := 0; i < 8; i++ {
		chan_int <- i
	}
	close(chan_int)
	for value := range chan_int {
		fmt.Println(fmt.Sprintf("Example 9.%d", value+1))
	}
}

func channels() {
	fmt.Println("Running examples...")
	//example_1()
	//example_2()
	//example_3()
	//example_4()
	//example_5()
	//example_6()
	//example_7()
	//example_8()
	example_9()
}
