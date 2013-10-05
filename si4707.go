//  Copyright (c) Marty Schoch
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package si4707

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"bitbucket.org/gmcbay/i2c"
	"github.com/mschoch/go-rds"
	"github.com/stianeikeland/go-rpio"
)

const I2C_ADDR = 0x63

const POWER_UP_TIME_MS = 1000

// Command Bytes
const COMMAND_POWER_UP = 0x01        // Powerup
const COMMAND_GET_REV = 0x10         // Revision info
const COMMAND_POWER_DOWN = 0x11      // Powerdown
const COMMAND_SET_PROPERTY = 0x12    // Sets property value
const COMMAND_GET_PROPERTY = 0x13    // Gets property value
const COMMAND_GET_INT_STATUS = 0x14  // Read interrupt status bits
const COMMAND_WB_TUNE_FREQ = 0x50    // Selects WB tuning frequency
const COMMAND_WB_TUNE_STATUS = 0x52  // Gets status of previous WB_TUNE_FREQ
const COMMAND_WB_RSQ_STATUS = 0x53   // RSQ of current channel
const COMMAND_WB_SAME_STATUS = 0x54  // SAME info for current channel
const COMMAND_WB_ASQ_STATUS = 0x55   // Gets status of 1050 Hz alert tone
const COMMAND_WB_AGC_STATUS = 0x57   // Gets status of AGC
const COMMAND_WB_AGC_OVERRIDE = 0x58 // Enable or disable WB AGC
const COMMAND_GPO_CTL = 0x80         // Configures GPO3 as output or hi-z
const COMMAND_GPO_SET = 0x81         // Sets GPO3 output level

type Device struct {
	bus    *i2c.I2CBus
	busNum byte
	addr   byte
}

func (d *Device) Init(busNum byte) (err error) {
	return d.InitCustomAddr(I2C_ADDR, busNum)
}

func (d *Device) InitCustomAddr(addr, busNum byte) (err error) {

	// do some manual GPIO to initialize the device
	err = rpio.Open()
	if err != nil {
		return err
	}

	pin23 := rpio.Pin(23)
	pin23.Output()

	pin23.Low()
	time.Sleep(1 * time.Second)
	pin23.High()
	time.Sleep(1 * time.Second)

	rpio.Close()

	if d.bus, err = i2c.Bus(busNum); err != nil {
		return
	}

	d.busNum = busNum
	d.addr = addr

	// wait max powerup time
	time.Sleep(110 * time.Millisecond)

	return
}

func (d *Device) PowerUp() {
	// 0x53 - GP02 output enabled, crystal osc enabled, WB receive mode
	// 0x05 - Use analog audio outputs
	err := d.bus.WriteByteBlock(d.addr, COMMAND_POWER_UP, []byte{0x53, 0x05})
	if err != nil {
		log.Printf("error writing: %v")
	}
}

func (d *Device) GetRev() int {

}

func (d *Device) writeCommand() {

}
