package main

import (
	"time"
	"github.com/embeddedgo/kendryte/devboard/maixbit/board/leds"
)

func main() {
	led := leds.Blue
	for {
		led.SetOn()
		time.Sleep(time.Second)
		led.SetOff()
		time.Sleep(time.Second)
	}
}

