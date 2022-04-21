package main

import "fmt"

//common interface for invoker and receiver
type command interface {
	execute() int
}

//bet aggregator contains the number of total bets placed , this is receiver
type betAggregator struct {
	totalBets int
}

//creating new instance of the receiver
func NewBetAggregator() *betAggregator {
	const totalBets = 0
	return &betAggregator{
		totalBets: totalBets,
	}
}

//placeBet command is represented as a struct, that contains betAggregator (receiver) as an attribute and
//is executed with the receiver
type placeBetCommand struct {
	aggregator *betAggregator
}

//the placeBet command itself, it implements the "command" interface with it`s "execute" method
func (pb *placeBetCommand) execute() int {
	//increase the number of placed bets
	pb.aggregator.totalBets += 1
	fmt.Printf("The bet is placed, total number of bets is: %d\n", pb.aggregator.totalBets)

	return pb.aggregator.totalBets
}

//placeBet - method for creating the instance of the command
func (ba *betAggregator) placeBet() command {
	return &placeBetCommand{
		aggregator: ba,
	}
}

// creating the struct that will execute the command
type BookMaker struct {
	Commands []command
}

func (bm *BookMaker) executeCommand() {
	for _, c := range bm.Commands {
		c.execute()
	}
}

func main() {
	ba := NewBetAggregator()

	tasks := []command{
		ba.placeBet(),
		ba.placeBet(),
		ba.placeBet(),
	}

	placers := &BookMaker{} //creating the instance of executor

	for _, task := range tasks {
		placers.Commands = append(placers.Commands, task)
	}

	placers.executeCommand()

	fmt.Println("total bets are:", ba.totalBets)
}
