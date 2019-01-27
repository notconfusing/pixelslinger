package opc

// Spatial Stripes
//   Creates spatial sine wave stripes: x in the red channel, y--green, z--blue
//   Also makes a white dot which moves down the strip non-spatially in the order
//   that the LEDs are indexed.

import (
	"github.com/longears/pixelslinger/colorutils"
	"github.com/longears/pixelslinger/midi"
	"math"
	"time"
)

func Spiral(x, y, t, SPIRAL_tightness, SPIRAL_speed, SPIRAL_thickness float64, SPIRAL_rings int) float64 {
	lhs := math.Sqrt((x*x + y*y))

	tt := t * SPIRAL_speed
	rhs_num := y*math.Cos(tt) + x*math.Sin(tt)
	rhs_den := x*math.Cos(tt) - y*math.Sin(tt)
	rhs_canonical := math.Atan2(rhs_num, rhs_den) * SPIRAL_tightness

	on_amount := 0.0
	for ring := 0; ring <= SPIRAL_rings; ring++ {
		SPIRAL_offset := SPIRAL_tightness * (2 * math.Pi)
		multSign := +1.0
		if SPIRAL_tightness < 0 {
			multSign = -1.0
		}
		rhs := rhs_canonical + (multSign * (float64(ring) * SPIRAL_offset))
		diff := lhs - rhs
		abs_diff := math.Abs(diff)
		if abs_diff < SPIRAL_thickness {
			on_amount += 1
		}
	}
	return on_amount
}

func MakePatternArchimedes(locations []float64) ByteThread {
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
				y = z //actually need to target x-z plane

				spiral1 := Spiral(x, y, t, 0.1, 2, 0.05, 4)
				spiral2 := Spiral(x, y, t, -0.1, 4, 0.05, 4)
				spiral3 := Spiral(x, y, t, -0.05, 6, 0.02, 8)

				r := spiral3
				g := spiral1
				b := spiral2

				bytes[ii*3+0] = colorutils.FloatToByte(r)
				bytes[ii*3+1] = colorutils.FloatToByte(g)
				bytes[ii*3+2] = colorutils.FloatToByte(b)

				//--------------------------------------------------------------------------------
			}
			bytesOut <- bytes
		}
	}
}
