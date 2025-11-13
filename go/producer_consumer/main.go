package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const maxOrders = 1_000_000
const ovens = 100

var wg sync.WaitGroup
var mtx sync.Mutex

var successOrders, failedOrders, totalOrders int
var recipes = [3]string{
	"croissant",
	"pâte à choux",
	"pain au chocolat",
}

func main() {

	//logging definitions
	logger := log.New(os.Stdout, "Go Brasserie: ", log.Ldate|log.Ltime|log.Lshortfile)
	cl := coloredLogger{
		logger: logger,
		colors: map[string]string{
			"failure": Red,
			"warning": Yellow,
			"success": Green,
			"info":    Cyan,
		},
	}

	//start program
	cl.ColoredPrintf("info", "La Brasserie est ouverte!") //cyan

	time.Sleep(time.Second * 2)

	//create cashierToKitchen and kitchenToTable channels
	cashierToKitchen := make(chan PastryOrder)
	kitchenToTable := make(chan PastryOrder)

	//add maxOrders to wait group
	wg.Add(maxOrders)

	//generate and send orders to the kitchen
	dispatchOrders(cashierToKitchen, cl)
	bakeOrders(cashierToKitchen, ovens, kitchenToTable, cl)

	var msg string
	for result := range kitchenToTable {
		if result.success {
			msg = fmt.Sprintf("order %d successful", result.id)
			cl.ColoredPrintf("success", msg) //green
		} else {
			msg = fmt.Sprintf("order %d failed", result.id)
			cl.ColoredPrintf("failure", msg) //red
		}
	}

	time.Sleep(time.Second * 2)

	// summarize results
	cl.ColoredPrintf("info", "total orders = %v", totalOrders)     //cyan
	cl.ColoredPrintf("info", "failed orders = %v", failedOrders)   //cyan
	cl.ColoredPrintf("info", "success orders = %v", successOrders) //cyan

	//end program
	cl.ColoredPrintf("info", "La Brasserie est fermée!") //cyan
}
