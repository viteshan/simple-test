package main

import (
	"fmt"
	"testing"

	"github.com/asaskevich/EventBus"
)

func calculator1(a int, b int) {
	fmt.Printf("calculator1:\t%d\n", a+b)
}
func calculator2(a int, b int) {
	fmt.Printf("calculator2:\t%d\n", a+b)
}
func TestEventBus(t *testing.T) {
	bus := EventBus.New()
	bus.Subscribe("main:calculator", calculator1)
	bus.Subscribe("main:calculator", calculator2)
	bus.Publish("main:calculator", 20, 40)
	bus.Unsubscribe("main:calculator", calculator1)
	bus.Unsubscribe("main:calculator", calculator2)
}
