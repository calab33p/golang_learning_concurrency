package main

import (
	"fmt"
	"sync"
)

const NumThreads = 1

type RequestResponse struct {
	input  string
	output string
}

func main() {

	inputs := []string{"Hello", "World", "These", "Are", "My", "Inputs"}
	results := ProcessData(inputs)

	for _, res := range results {
		fmt.Printf("input: %s, output: %s\n", res.input, res.output)
	}
}

func ProcessData(inputs []string) []RequestResponse {
	var result []RequestResponse

	var wg sync.WaitGroup
	wg.Add(NumThreads)

	fmt.Printf("Performing network processing for %d inputs...\n", len(inputs))

	lookupChan := Generator(inputs)
	//resultChan := FanIn(responseChans)
	resultChan := make(chan []RequestResponse)
	//resultChan := make(chan []RequestResponse, NumThreads)

	// fan-out to worker  goroutines
	for j := 0; j < NumThreads; j++ {

		//log.Debug().Msgf("Creating goroutine %d", j)
		fmt.Printf("Creating goroutine %d\n", j)

		go NetworkThread(lookupChan, resultChan, &wg)
	}

	//log.Debug().Msg("Starting to read results")
	//log.Debug().Msgf("Size of result chan: %d", len(resultChan))
	fmt.Println("Starting to read results")
	//fmt.Printf("Size of result chan: %d\n", len(resultChan))

	//log.Debug().Msg("Waiting on worker threads")
	fmt.Println("Waiting on worker threads")

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for rsp := range resultChan {
		fmt.Printf("Got result: %v\n", rsp)
		result = append(result, rsp...)
	}

	return result

}

func Generator(inputs []string) <-chan RequestResponse {
	//lookupChan := make(chan RequestResponse, NumThreads)
	lookupChan := make(chan RequestResponse)

	go func() {
		// send inputs to channel for goroutines to pick up
		for _, in := range inputs {
			fmt.Printf("	Adding input %s to lookup channel\n", in)
			lookupChan <- RequestResponse{input: in, output: ""}
		}
		//log.Debug().Msgf("Closing lookup chan of size %d", len(lookupChan))
		fmt.Printf("Closing lookup chan of size %d\n", len(lookupChan))
		close(lookupChan)
	}()
	return lookupChan

}

func NetworkThread(lc <-chan RequestResponse, rc chan<- []RequestResponse, wg *sync.WaitGroup) {

	defer wg.Done()

	for rr := range lc {

		var result []RequestResponse

		in := rr.input

		result = append(result, RequestResponse{input: in, output: "Processed"})

		//log.Debug().Msgf("Sending %v to response channel", res)
		fmt.Printf("Sending %v to response channel\n", result)
		rc <- result
	}
	//log.Debug().Msg("Ending thread")
	fmt.Println("Ending thread")
}
