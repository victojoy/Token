package main

import "fmt"
import "strconv"
import "os"
import "time"
import "math/rand"

type Token struct {
 	data string
 	recipient int
 	ttl int
 } 

func send(prev <-chan Token, next chan<- Token, id int) {
				for  {	
					token := <- prev				
					if token.data == "Check" {
						next <- token
						close(next)
						return
					}
					fmt.Printf("Get by %d: %d\n", id, token.ttl)
					token.ttl--
					

					if token.recipient == id {
						fmt.Printf("Complited!\n")												
						token.data = "Check"
						next <- token	
						close(next)
						
						return
					} else if  token.ttl < -1 {
						fmt.Println("Ah Shit, Here We Go Again!\n")						
						token.data = "Check"
						next <- token	
						close(next)
						return
					} else {
						next <- token			  				
					}
				}
}
			
func main() {
	if (len(os.Args) < 2) {
		fmt.Println("USAGE: go run 'Channels.go' <Number>")
		os.Exit(0);
	}

	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	fmt.Printf("N=%d\n", N)
	if N <= 1 {
		fmt.Printf("Invalid N=%d, change N\n", N)
		os.Exit(2);
	}


	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	msg := Token{"Random message", r.Intn(N), r.Intn(N*2)}
	fmt.Println(msg);


	ch := make(chan Token)
	last := ch
	for i := 1; i < N; i = i + 1 {
		current := make(chan Token)
		go send(last, current, i);
		last = current;
	}	 
	
	go send(last, ch, 0)	
	ch <- msg
	<- last
}