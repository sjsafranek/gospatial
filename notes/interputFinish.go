package main

import (
    // "bufio"
    "fmt"
    "os"
    "os/signal"
    "sync"
    "time"
)


func main() {
	messages := make(chan int)
	var wg sync.WaitGroup

	// https://gobyexample.com/signals
	// ^C is SIGINT
    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
        sig := <-sigs
        fmt.Println()
        fmt.Println(sig)
        fmt.Println("Finishing jobs...")
        done <- true
    }()

	fmt.Println("awaiting signal")
    // <-done
    // fmt.Println("exiting")


    // you can also add these one at 
    // a time if you need to 
	for i:=0; i<=100000;i++{
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Second * 3)
			messages <- i
		}(i)
	}

    go func() {
        for i := range messages {
            fmt.Println(i)
        }
        done <- true
        close(sigs)
    }()



	<-done

    wg.Wait()

    fmt.Println("exiting")

	// // KEY PRESS
 //    reader := bufio.NewReader(os.Stdin)
 //    input, _ := reader.ReadString('\n')
 //    fmt.Printf("Input Char Is : %v", string([]byte(input)[0]))

}


