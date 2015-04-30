package LiquidCrystal

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"time"
)

const (
    // commands
    LCD_CLEARDISPLAY byte = 0x01
    LCD_RETURNHOME byte = 0x02
    LCD_ENTRYMODESET byte = 0x04
    LCD_DISPLAYCONTROL byte = 0x08
    LCD_CURSORSHIFT byte = 0x10
    LCD_FUNCTIONSET byte = 0x20
    LCD_SETCGRAMADDR byte = 0x40
    LCD_SETDDRAMADDR byte = 0x80

    // flags for display entry mode
    LCD_ENTRYRIGHT byte = 0x00
    LCD_ENTRYLEFT byte = 0x02
    LCD_ENTRYSHIFTINCREMENT byte = 0x01
    LCD_ENTRYSHIFTDECREMENT byte = 0x00

    // flags for display on/off control
    LCD_DISPLAYON byte = 0x04
    LCD_DISPLAYOFF byte = 0x00
    LCD_CURSORON byte = 0x02
    LCD_CURSOROFF byte = 0x00
    LCD_BLINKON byte = 0x01
    LCD_BLINKOFF byte = 0x00

    // flags for display/cursor shift
    LCD_DISPLAYMOVE byte = 0x08
    LCD_CURSORMOVE byte = 0x00
    LCD_MOVERIGHT byte = 0x04
    LCD_MOVELEFT byte = 0x00

    // flags for function set
    LCD_8BITMODE byte = 0x10
    LCD_4BITMODE byte = 0x00
    LCD_2LINE byte = 0x08
    LCD_1LINE byte = 0x00
    LCD_5x10DOTS byte = 0x04
    LCD_5x8DOTS byte = 0x00

    // flags for backlight control
    LCD_BACKLIGHT byte = 0x08
    LCD_NOBACKLIGHT byte = 0x00

    En byte = 1 << 2  // Enable bit
    Rw byte = 1 << 1  // Read/Write bit
    Rs byte = 1 << 0  // Register select bit
)

var _ gobot.Driver = (*LiquidCrystalDriver)(nil)

type LiquidCrystalDriver struct {
	name        string
	connection  i2c.I2c
    addr        byte
    backlight   byte
    cols        int
    rows        int
    charsize    int
    displayfunc byte
    displayctrl byte
    displaymode byte
}

// NewLiquidCrystalDriver creates a new driver with specified name and i2c interface
func NewLiquidCrystalDriver(a i2c.I2c, name string, addr byte, cols int, rows int) *LiquidCrystalDriver {
	return &LiquidCrystalDriver{
		name:       name,
		connection: a,
        addr:       addr,
        backlight:  LCD_BACKLIGHT,
        cols:       cols,
        rows:       rows,
        charsize:   int(LCD_5x8DOTS),
	}
}

func (h *LiquidCrystalDriver) SetCharSize(size int) {
    h.charsize = size
}

func (h *LiquidCrystalDriver) Name() string                 { return h.name }
func (h *LiquidCrystalDriver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start initialized the LIDAR
func (h *LiquidCrystalDriver) Start() (errs []error) {
	if err := h.connection.I2cStart(h.addr); err != nil {
		return []error{err}
	}
    h.displayfunc = LCD_4BITMODE | LCD_1LINE | LCD_5x8DOTS
    if h.rows > 1 {
        h.displayfunc |= LCD_2LINE
    }
	// for some 1 line displays you can select a 10 pixel high font
    if h.charsize != 0 && h.rows == 1 {
        h.displayfunc |= LCD_5x10DOTS
    }
	// SEE PAGE 45/46 FOR INITIALIZATION SPECIFICATION!
	// according to datasheet, we need at least 40ms after power rises above 2.7V
	// before sending commands. Arduino can turn on way befer 4.5V so we'll wait 50
    <-time.After(50 * time.Millisecond)

	// Now we pull both RS and R/W low to begin commands
    h.expanderWrite(h.backlight) // reset expanderand turn backlight off (Bit 8 =1)
    <-time.After(1 * time.Second)

	//put the LCD into 4 bit mode
	// this is according to the hitachi HD44780 datasheet
	// figure 24, pg 46

	// we start in 8bit mode, try to set 4 bit mode
    h.write4bits(0x03 << 4)
    <-time.After(4500 * time.Microsecond) // wait min 4.1ms

	// second try
    h.write4bits(0x03 << 4)
    <-time.After(4500 * time.Microsecond) // wait min 4.1ms

	// third go!
    h.write4bits(0x03 << 4)
    <-time.After(150 * time.Microsecond)

	// finally, set to 4-bit interface
	h.write4bits(0x02 << 4)

	// set # lines, font size, etc.
	h.command(LCD_FUNCTIONSET | h.displayfunc)

	// turn the display on with no cursor or blinking default
	h.displayctrl = LCD_DISPLAYON | LCD_CURSOROFF | LCD_BLINKOFF
	h.Display()

	// clear it off
	h.Clear()

	// Initialize to default text direction (for roman languages)
	h.displaymode = LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT

	// set the entry mode
	h.command(LCD_ENTRYMODESET | h.displaymode)

	h.Home()
	return
}

// Halt returns true if devices is halted successfully
func (h *LiquidCrystalDriver) Halt() (errs []error) {
    h.Clear()
    h.NoBacklight()
    h.NoCursor()
    h.NoDisplay()
    return
}

/********** high level commands, for the user! */
func (h *LiquidCrystalDriver) Clear() (error) {
    var err = h.command(LCD_CLEARDISPLAY)// clear display, set cursor position to zero
    <-time.After(2 * time.Millisecond)// this command takes a long time!
    return err
}

func (h *LiquidCrystalDriver) Home() (error) {
    var err = h.command(LCD_RETURNHOME)// set cursor position to zero

    <-time.After(2 * time.Millisecond)// this command takes a long time!
    return err
}

func (h *LiquidCrystalDriver) SetCursor(col int, row int) (error){
	var row_offsets = []byte{ 0x00, 0x40, 0x14, 0x54 }
	if row > h.rows {
		row = h.rows - 1;    // we count rows starting w/0
	}
	return h.command(LCD_SETDDRAMADDR | (byte(col) + row_offsets[row]))
}

// Turn the display on/off (quickly)
func (h *LiquidCrystalDriver) NoDisplay() (error) {
	h.displayctrl &= ^LCD_DISPLAYON
	return h.command(LCD_DISPLAYCONTROL | h.displayctrl)
}
func (h *LiquidCrystalDriver) Display()  (error) {
	h.displayctrl |= LCD_DISPLAYON
	return h.command(LCD_DISPLAYCONTROL | h.displayctrl)
}

// Turns the underline cursor on/off
func (h *LiquidCrystalDriver) NoCursor()  (error) {
	h.displayctrl &= ^LCD_CURSORON
	return h.command(LCD_DISPLAYCONTROL | h.displayctrl)
}
func (h *LiquidCrystalDriver) Cursor()  (error) {
	h.displayctrl |= LCD_CURSORON
	return h.command(LCD_DISPLAYCONTROL | h.displayctrl)
}

// Turn on and off the blinking cursor
func (h *LiquidCrystalDriver) NoBlink()  (error) {
	h.displayctrl &= ^LCD_BLINKON
	return h.command(LCD_DISPLAYCONTROL | h.displayctrl)
}
func (h *LiquidCrystalDriver) Blink()  (error) {
	h.displayctrl |= LCD_BLINKON
	return h.command(LCD_DISPLAYCONTROL | h.displayctrl)
}

// These commands scroll the display without changing the RAM
func (h *LiquidCrystalDriver) ScrollDisplayLeft()  (error) {
	return h.command(LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVELEFT)
}
func (h *LiquidCrystalDriver) ScrollDisplayRight()  (error) {
	return h.command(LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVERIGHT)
}

// This is for text that flows Left to Right
func (h *LiquidCrystalDriver) LeftToRight()  (error) {
	h.displaymode |= LCD_ENTRYLEFT
	return h.command(LCD_ENTRYMODESET | h.displaymode)
}

// This is for text that flows Right to Left
func (h *LiquidCrystalDriver) RightToLeft()  (error) {
	h.displaymode &= ^LCD_ENTRYLEFT
	return h.command(LCD_ENTRYMODESET | h.displaymode)
}

// This will 'right justify' text from the cursor
func (h *LiquidCrystalDriver) Autoscroll()  (error) {
	h.displaymode |= LCD_ENTRYSHIFTINCREMENT
	return h.command(LCD_ENTRYMODESET | h.displaymode)
}

// This will 'left justify' text from the cursor
func (h *LiquidCrystalDriver) NoAutoscroll()  (error) {
	h.displaymode &= ^LCD_ENTRYSHIFTINCREMENT
	return h.command(LCD_ENTRYMODESET | h.displaymode)
}

// Allows us to fill the first 8 CGRAM locations
// with custom characters
func (h *LiquidCrystalDriver) CreateChar(location byte, charmap []byte) (error) {
	location &= 0x7; // we only have 8 locations 0-7
    err := h.command(LCD_SETCGRAMADDR | (location << 3))
    for _, char := range charmap {
		h.Write(char)
    }
    return err
}

// Turn the (optional) backlight off/on
func (h *LiquidCrystalDriver) NoBacklight() {
	h.backlight=LCD_NOBACKLIGHT
	h.expanderWrite(0)
}

func (h *LiquidCrystalDriver) Backlight() {
	h.backlight=LCD_BACKLIGHT
	h.expanderWrite(0)
}

/*********** mid level commands, for sending data/cmds */

func (h *LiquidCrystalDriver) command(value byte) (error) {
    return h.send(value, 0)
}

func (h *LiquidCrystalDriver) Write(value byte) (int, error) {
    if err := h.send(value, Rs); err != nil {
        return 0, err
    }
    return 1, nil
}

/************ low level data pushing commands **********/
func (h *LiquidCrystalDriver) expanderWrite(data byte) error {
    return h.connection.I2cWrite([]byte{data | h.backlight})
}

func (h *LiquidCrystalDriver) pulseEnable(data byte) (err error) {
	if err = h.expanderWrite(data | En); err != nil {    // En high
        return
    }
    <-time.After(1 * time.Microsecond)  // enable pulse must be >450ns

	if err = h.expanderWrite(data & ^En); err != nil {// En low
        return
    }
    <-time.After(50 * time.Microsecond) // commands need > 37us to settle
    return
}

func (h *LiquidCrystalDriver) write4bits(value byte) (err error) {
    if err = h.expanderWrite(value); err != nil {
        return
    }
    if err = h.pulseEnable(value); err != nil {
        return
    }
    return
}

func (h *LiquidCrystalDriver) send(value, mode byte) (err error) {
    var highnib = value & 0xf0
    var lownib = (value << 4) & 0xf0
    if err = h.write4bits(highnib | mode); err != nil {
        return
    }
    if err = h.write4bits(lownib | mode); err != nil {
        return
    }
    return
}

func (h *LiquidCrystalDriver) LoadCustomCharacter(char_num byte, rows []byte){
	h.CreateChar(char_num, rows)
}

func (h *LiquidCrystalDriver) SetBacklight(new_val bool){
	if new_val {
		h.Backlight()		// turn backlight on
	} else {
		h.NoBacklight()		// turn backlight off
	}
}

func (h *LiquidCrystalDriver) Print(str string) {
    var charmap = []byte(str)
    for _, char := range charmap {
      h.Write(char)
    }
}
