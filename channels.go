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
	//Selecting the channel and causing a deadlock
	chan_string := make(chan string)
	chan_int := make(chan int)
	go func() {
		select {
		case string_value := <-chan_string:
			fmt.Println(string_value)
		case int_value := <-chan_int:
			fmt.Println(fmt.Sprintf("Example %s\\.2", int_value))
		}
	}()
	chan_int <- 7
	chan_string <- "Example 7.1"
}

func example_8() {
	//range over a channel
}

func channels() {
	fmt.Println("Running examples...")
	example_1()
	example_2()
	example_3()
	example_4()
	example_5()
	example_6()
	example_7()
}
