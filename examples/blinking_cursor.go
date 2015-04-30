package main

import (
	"github.com/Lupino/LiquidCrystal"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"time"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	lcd := LiquidCrystal.NewLiquidCrystalDriver(firmataAdaptor,
		"LiquidCrystal",
		0x27,
		16,
		2)

	work := func() {
		var blinking = false
		lcd.Cursor()

		for {
			if blinking {
				lcd.Clear()
				lcd.Print("No cursor blink")
				lcd.NoBlink()
				blinking = false
			} else {
				lcd.Clear()
				lcd.Print("Cursor blink")
				lcd.Blink()
				blinking = true
			}
			<-time.After(4 * time.Second)
		}
	}

	robot := gobot.NewRobot("LiquidCrystal",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{lcd},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
