LiquidCrystal library for [Gobot](http://gobot.io)
=========================================

This is a LiquidCrystal library for the [Gobot](http://gobot.io) to interface with liquid crystal (LCD) displays.

Install
-------

```bash
go get -v github.com/Lupino/LiquidCrystal
```

Example
-------

```go
package main

import (
    "github.com/Lupino/LiquidCrystal"
    "github.com/hybridgroup/gobot"
    "github.com/hybridgroup/gobot/platforms/firmata"
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
        lcd.Print("Hello World!")
    }

    robot := gobot.NewRobot("LiquidCrystal",
        []gobot.Connection{firmataAdaptor},
        []gobot.Device{lcd},
        work,
    )

    gbot.AddRobot(robot)

    gbot.Start()
}
```
