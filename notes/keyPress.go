package main

import (
    "bufio"
    "fmt"
    "os"
    "os/signal"
)

func main() {
	// ^C is SIGINT
    c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
	    for sig := range c {
	        // sig is a ^C, handle it
	        fmt.Printf("%s \n", sig)
	        fmt.Println("Waiting for jobs to finish...")
	        
	    }
	}()

	// KEY PRESS
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')

    fmt.Printf("Input Char Is : %v", string([]byte(input)[0]))

    // fmt.Printf("You entered: %v", []byte(input))
}