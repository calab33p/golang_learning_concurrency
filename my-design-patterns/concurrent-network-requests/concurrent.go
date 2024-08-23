package main

import (
	"fmt"
	"sync"
)

const NumThreads = 5

type RequestResponse struct {
	input  string
	output string
}

func main() {

	inputs := []string{"Hello", "World", "These", "Are", "My", "inputs"}
	results := ProcessData(inputs)

	for _, res := range results {
		fmt.Printf("input: %s, output: %s\n", res.input, res.output)
	}
}

func ProcessData(inputs []string) []RequestResponse {
	var result []RequestResponse

	lookupChan := make(chan RequestResponse, NumThreads)
	resultChan := make(chan []RequestResponse, NumThreads)

	var wg sync.WaitGroup
	wg.Add(NumThreads)

	fmt.Printf("Performing network processing for %d inputs...\n", len(inputs))

	// create goroutines
	for j := 0; j < NumThreads; j++ {

		//log.Debug().Msgf("Creating goroutine %d", j)
		fmt.Printf("Creating goroutine %d\n", j)

		go NetworkThread(lookupChan, resultChan, &wg)
	}

	// send inputs to channel for goroutines to pick up
	for _, in := range inputs {
		fmt.Printf("	Adding input %s to lookup channel\n", in)
		lookupChan <- RequestResponse{input: in, output: ""}
	}
	//log.Debug().Msgf("Closing lookup chan of size %d", len(lookupChan))
	fmt.Printf("Closing lookup chan of size %d\n", len(lookupChan))
	close(lookupChan)

	//log.Debug().Msg("Starting to read results")
	//log.Debug().Msgf("Size of result chan: %d", len(resultChan))
	fmt.Println("Starting to read results")
	fmt.Printf("Size of result chan: %d\n", len(resultChan))

	//TODO: close resultChan first and use range instead?
	for k := 0; k < len(resultChan); k++ {
		rsp := <-resultChan
		//log.Debug().Msgf("Got result: &v", rsp)
		fmt.Printf("Got result: %v\n", rsp)
		for _, r := range rsp {
			//log.Debug().Msgf("Got result: &s", r)
			fmt.Printf("Got result: %s\n", r)
			result = append(result, r)
		}
	}

	//log.Debug().Msg("Waiting on worker threads")
	fmt.Println("Waiting on worker threads")
	wg.Wait()

	//log.Debug().Msg("Closing result chan")
	fmt.Println("Closing result chan")
	close(resultChan)

	return result

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
