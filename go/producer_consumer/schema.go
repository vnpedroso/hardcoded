package main

import "math/rand"

type PastryOrder struct {
	recipe  string
	id      int
	success bool
}

func genRandomOrder() PastryOrder {
	return PastryOrder{
		recipe:  recipes[rand.Intn(len(recipes))],
		id:      0,
		success: true,
	}
}
