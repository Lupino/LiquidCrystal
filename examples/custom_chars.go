package main

import (
	"github.com/Lupino/LiquidCrystal"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"time"
)

var bell = []byte{0x4, 0xe, 0xe, 0xe, 0x1f, 0x0, 0x4}
var note = []byte{0x2, 0x3, 0x2, 0xe, 0x1e, 0xc, 0x0}
var clock = []byte{0x0, 0xe, 0x15, 0x17, 0x11, 0xe, 0x0}
var heart = []byte{0x0, 0xa, 0x1f, 0x1f, 0xe, 0x4, 0x0}
var duck = []byte{0x0, 0xc, 0x1d, 0xf, 0xf, 0x6, 0x0}
var check = []byte{0x0, 0x1, 0x3, 0x16, 0x1c, 0x8, 0x0}
var cross = []byte{0x0, 0x1b, 0xe, 0x4, 0xe, 0x1b, 0x0}
var retarrow = []byte{0x1, 0x1, 0x5, 0x9, 0x1f, 0x8, 0x4}

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	lcd := LiquidCrystal.NewLiquidCrystalDriver(firmataAdaptor,
		"LiquidCrystal",
		0x27,
		16,
		2)

	work := func() {
		lcd.Backlight()
		lcd.CreateChar(0, bell)
		lcd.CreateChar(1, note)
		lcd.CreateChar(2, clock)
		lcd.CreateChar(3, heart)
		lcd.CreateChar(4, duck)
		lcd.CreateChar(5, check)
		lcd.CreateChar(6, cross)
		lcd.CreateChar(7, retarrow)
		lcd.Home()

		lcd.Print("Hello world...")
		lcd.SetCursor(0, 1)
		lcd.Print(" i ")
		lcd.Write(3)
		lcd.Print(" arduinos!")
		<-time.After(5 * time.Second)
		go displayKeyCodes(lcd)
	}

	robot := gobot.NewRobot("LiquidCrystal",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{lcd},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}

// display all keycodes
func displayKeyCodes(lcd *LiquidCrystal.LiquidCrystalDriver) {
	var i byte = 0

	for {
		lcd.Clear()
		lcd.Printf("Codes 0x%02x-0x%02x", i, i+16)
		lcd.SetCursor(0, 1)

		for j := 0; j < 16; j++ {
			lcd.Write(i + byte(j))
		}
		i += 16

		<-time.After(4 * time.Second)
	}
}
