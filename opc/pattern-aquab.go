package opc

// Spatial Stripes
//   Creates spatial sine wave stripes: x in the red channel, y--green, z--blue
//   Also makes a white dot which moves down the strip non-spatially in the order
//   that the LEDs are indexed.

import (
	"github.com/longears/pixelslinger/colorutils"
	"github.com/longears/pixelslinger/midi"
	//"math"
	"time"
	"fmt"
)


func MakePatternAquaB(locations []float64) ByteThread {

	var (
		FLUSH_LENGTH = 2.0 //seconds
		FLUSH_REFILL = 5.0 //seconds
		FLUSH_REST = 1.0 //seconds
		FLUSH_CYCLE = FLUSH_LENGTH+FLUSH_REST+FLUSH_REFILL
		//calculate the positions in the [0,1] domain when flush and refill will happen
		// order is: rest, flush, refill
		FLUSH_POINT = 1.0 - ((FLUSH_LENGTH+FLUSH_REFILL)/FLUSH_CYCLE)
		REFILL_POINT = 1.0 - (FLUSH_REFILL/FLUSH_CYCLE)
		UNIT_REST = FLUSH_REST / FLUSH_CYCLE
		UNIT_REST_FLUSH = (FLUSH_REST + FLUSH_LENGTH) / FLUSH_CYCLE
		UNIT_FLUSH = FLUSH_LENGTH / FLUSH_CYCLE
		UNIT_REFILL = FLUSH_REFILL / FLUSH_CYCLE
		)
	// get bounding box
	n_pixels := len(locations) / 3
	var max_coord_x, max_coord_y, max_coord_z float64
	var min_coord_x, min_coord_y, min_coord_z float64
	for ii := 0; ii < n_pixels; ii++ {
		x := locations[ii*3+0]
		y := locations[ii*3+1]
		z := locations[ii*3+2]
				if ii == 0 || x > max_coord_x { max_coord_x = x }
				if ii == 0 || y > max_coord_y { max_coord_y = y }
				if ii == 0 || z > max_coord_z { max_coord_z = z }
				if ii == 0 || x < min_coord_x { min_coord_x = x }
				if ii == 0 || y < min_coord_y { min_coord_y = y }
				if ii == 0 || z < min_coord_z { min_coord_z = z }
	}

	return func(bytesIn chan []byte, bytesOut chan []byte, midiState *midi.MidiState) {
		for bytes := range bytesIn {
			n_pixels := len(bytes) / 3
			t := float64(time.Now().UnixNano())/1.0e9 - 9.4e8
			// fill in bytes slice
			for ii := 0; ii < n_pixels; ii++ {
				//--------------------------------------------------------------------------------

				// make moving stripes for x, y, and z
				x := locations[ii*3+0]
				y := locations[ii*3+1]
				z := locations[ii*3+2]
				//r := colorutils.Cos(x, t/4, 1, 0, 0.7) // offset, period, min, max

				s1 := colorutils.Cos((z*0.3)*(y*0.3)*(x*0.3)    , t/8, 1, 0.05, 0.7)
				s3 := colorutils.Cos(x*0.2 + y*0.8, t/16, 1, 0.05, 0.5)
				s2 := colorutils.Cos(z*0.2, t/20, 1, 0.05, 0.8)
				s4 := 2.1*s1

				r:= 0.0
				g:= 0.0
				b:= 0.0

				// number of colors
				nc := 2.0


				//lightcyan
				r += (0.875*s1) / nc
				g += (1.000*s1) / nc
				b += (1.000*s1) / nc

				//cyan
				r += (0.000*s2) / nc
				g += (1.000*s2) / nc
				b += (1.000*s2) / nc

				// aquamarine
				r += (0.495 *s3) /nc
				g += (1.000 *s3) /nc
				b += (0.831 *s3) /nc

				// teal
				r += (0.000 *s4) /nc
				g += (0.600 *s4) /nc
				b += (0.700 *s4) /nc


				// flushing

				//position in cycle on the interval [0,1]
				flush_cyc_pos := colorutils.PosMod2(t/FLUSH_CYCLE, 1)
				//flush part
				if flush_cyc_pos > FLUSH_POINT && flush_cyc_pos < REFILL_POINT {
					flush_amount := 1- (flush_cyc_pos - UNIT_REST) / UNIT_FLUSH
					z_amount := (z - min_coord_z) / max_coord_z
					if flush_amount < z_amount {
						r = 0
						g = 0
						b = 0
					}

				}
				//refill part
				if flush_cyc_pos > REFILL_POINT {
					flush_amount :=  (flush_cyc_pos - UNIT_REST_FLUSH) / UNIT_REFILL
					z_amount := (z - min_coord_z) / max_coord_z
					if flush_amount < z_amount {
						r = 0
						g = 0
						b = 0
					}
				}
				fmt.Sprintf("%v",FLUSH_LENGTH)
				fmt.Sprintf("%v",FLUSH_REST)
				fmt.Sprintf("%v",FLUSH_CYCLE)

				bytes[ii*3+0] = colorutils.FloatToByte(r)
				bytes[ii*3+1] = colorutils.FloatToByte(g)
				bytes[ii*3+2] = colorutils.FloatToByte(b)

				//--------------------------------------------------------------------------------
			}
			bytesOut <- bytes
		}
	}
}
