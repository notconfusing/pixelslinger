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

//func makeSpiral() ByteThread{
//	return func(x,y, t, rings, SPIRAL_tightness, SPIRAL_speed, SPIRAL_thickness float64){
//
//		spiral_amount :=
//		return spiral_amount
//	}
//}

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

				//tao := 2*math.Pi

				SPIRAL_tightness := -0.1
				SPIRAL_thickness := 0.05
				SPIRAL_speed := 3.0
				SPIRAL_rings := 4

				tt := t * SPIRAL_speed

				y = z
				lhs := math.Sqrt((x*x + y*y))
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

				r := 0.0
				g := on_amount
				b := 0.0

				bytes[ii*3+0] = colorutils.FloatToByte(r)
				bytes[ii*3+1] = colorutils.FloatToByte(g)
				bytes[ii*3+2] = colorutils.FloatToByte(b)

				//--------------------------------------------------------------------------------
			}
			bytesOut <- bytes
		}
	}
}
