// This file contains linux specific optimized functions.

package gfx

import (
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/skelterjohn/go.wde/xgb"
)

func init() {
	clearFunc = clear
}

// clear clears the destination image of ctx.
// Returns true if it was successful.
func clear(ctx *Context) bool {
	xi, ok := ctx.dst.(*xgb.Image)
	if !ok {
		return false
	}

	col := xgraphics.BGRAModel.Convert(ctx.col).(xgraphics.BGRA)
	r := ctx.dst.Bounds()

	pix := xi.Pix

	// Do first line "manually":
	offs := xi.PixOffset(r.Min.X, r.Min.Y)
	firstLine := pix[offs : offs+r.Dx()*4]
	for i := 0; i < len(firstLine); i += 4 {
		firstLine[i] = col.B
		firstLine[i+1] = col.G
		firstLine[i+2] = col.R
		firstLine[i+3] = col.A
	}

	// Then copy the first line:
	pix = pix[offs:]
	for y := r.Min.Y + 1; y < r.Max.Y; y++ {
		pix = pix[xi.Stride:]
		copy(pix, firstLine)
	}

	return true
}
