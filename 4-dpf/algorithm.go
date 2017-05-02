package main

import (
	"math"

	"gopkg.in/distil.v4"
)

//This is our distillate algorithm
type DisplacementPFDistiller struct {
	// This line is required. It says this struct inherits some useful
	// default methods.
	distil.DistillateTools

	// For the frequency distillate, we make use of a rebase stage, to do
	// that, we need to know the intended frequency of the stream
	basefreq int64
}

// PadSnap is a rebase stage that will adjust the incoming data to strictly
// appear on a timebase of the given frequency (hence the 'rebase'). Any
// values that do not appear with exactly the right time are snapped to
// the nearest time, and any duplicates are dropped. In addition, any missing
// values are replaced by NaN (hence pad). The advantage of this is that it
// simplifies calculations that refer to values across time, you can rest
// assured that a value 1s ago is exactly basefreq samples away, even if
// there were holes in the data or if there were duplicates. Note that in
// the presence of duplicate data there is ZERO GUARANTEE as to WHICH of the
// multiple duplicate values you receive. In general this makes algorithms
// that compare across time quite useless, as the real time difference
// between the points a fixed interval apart will experience extreme jitter.
// The default implementation (in DistillateTools) returns RebasePassthrough
func (d *DisplacementPFDistiller) Rebase() distil.Rebaser {
	return distil.RebasePadSnap(d.basefreq)
}

// This is our main algorithm. It will automatically be called with chunks
// of data that require processing by the engine.
func (d *DisplacementPFDistiller) Process(in *distil.InputSet, out *distil.OutputSet) {
	/* Output 0 is real_power.
	 * Output 1 is reactive_power.
	 */
	var ns int = in.NumSamples(0)
	var i int
	for i = 0; i < ns; i++ {
		var time int64 = in.Get(0, i).T

		var angI = in.Get(0, i).V
		var angV = in.Get(1, i).V

		var dpf = 100 * math.Cos((angV*math.Pi/180)-(angI*math.Pi/180))

		if !math.IsNaN(dpf) {
			out.Add(0, time, dpf)
		}
	}
}
