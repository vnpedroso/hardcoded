package main

import "fmt"

func dispatchOrders(cashierToKitchen chan<- PastryOrder, cl coloredLogger) {
	go func() {
		for id := 1; id <= maxOrders; id++ {
			newOrder := genRandomOrder()
			totalOrders++
			newOrder.id = id
			cashierToKitchen <- newOrder
			cl.logger.Printf("%s order id %d sent to kitchen", newOrder.recipe, newOrder.id)
		}
		close(cashierToKitchen)
	}()
}

func bakeOrders(cashierToKitchen <-chan PastryOrder, ovens int, kitchenToTable chan<- PastryOrder, cl coloredLogger) {
	for ov := 0; ov < ovens; ov++ {
		go func(ovenId int) {
			for order := range cashierToKitchen {
				order.success = NotSoRandomSuccess()
				if order.success {
					mtx.Lock()
					successOrders++
					mtx.Unlock()
				} else {
					mtx.Lock()
					failedOrders++
					mtx.Unlock()
				}
				kitchenToTable <- order
				wg.Done()
			}
			ovenTurnOffMsg := fmt.Sprintf("Oven %d turning off...", ov)
			cl.ColoredPrintf("warning", ovenTurnOffMsg) //yellow
		}(ov)
	}

	go func() {
		wg.Wait()
		close(kitchenToTable)
	}()
}
