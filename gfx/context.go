// Package gfx contains drawing utilities for a draw.Image destination.
package gfx

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Context provides primitive drawing operations on a draw.Image.
// All operations use the drawing color set by SetColor().
type Context struct {
	// dst is the image to draw on
	dst draw.Image

	// col is the drawing color
	col color.Color
}

// NewContext returns a new Context initialized with black drawing color.
func NewContext(dst draw.Image) *Context {
	return &Context{
		dst: dst,
		col: color.NRGBA{0, 0, 0, 255},
	}
}

// SetColor sets the drawing color.
func (ctx *Context) SetColor(col color.Color) {
	ctx.col = col
}

// SetColorRGBA sets the drawing color.
func (ctx *Context) SetColorRGBA(r, g, b, a byte) {
	ctx.SetColor(color.RGBA{r, g, b, a})
}

// Clear clears the destination image.
func (ctx *Context) Clear() {
	draw.Draw(ctx.dst, ctx.dst.Bounds(), &image.Uniform{ctx.col}, image.ZP, draw.Src)
}

// Point draws a point.
func (ctx *Context) Point(x, y int) {
	ctx.dst.Set(x, y, ctx.col)
}

// HLine draws a horizontal line.
func (ctx *Context) HLine(x1, x2, y int) {
	dst, col := ctx.dst, ctx.col
	for x := x1; x <= x2; x++ {
		dst.Set(x, y, col)
	}
}

// VLine draws a vertical line.
func (ctx *Context) VLine(y1, y2, x int) {
	dst, col := ctx.dst, ctx.col
	for y := y1; y <= y2; y++ {
		dst.Set(x, y, col)
	}
}

// Rectangle draws a rectangle.
func (ctx *Context) Rectangle(x1, y1, width, height int) {
	x2, y2 := x1+width-1, y1+height-1
	ctx.HLine(x1, x2, y1)
	ctx.HLine(x1, x2, y2)
	ctx.VLine(y1+1, y2-1, x1)
	ctx.VLine(y1+1, y2-1, x2)
}

// FillCircle draws a filled circle.
// The Midpoint circle algorithm is used which is detailed here:
// https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func (ctx *Context) FillCircle(x0, y0, rad int) {
	for x, y, err := rad, 0, 0; x > 0; {
		ctx.HLine(x0-x, x0+x, y0-y)
		ctx.HLine(x0-x, x0+x, y0+y)

		if err <= 0 {
			y++
			err += 2*y + 1
		}
		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

// DrawString draws a string.
// The y coordinate is the bottom line of the text.
func (ctx *Context) DrawString(s string, x, y int) {
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  ctx.dst,
		Src:  image.NewUniform(ctx.col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(s)
}

// abs reutrns the absolute of an int.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DrawLine draws a line.
// The Bresenham's line algorithm is used which is detailed here:
// https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func (ctx *Context) DrawLine(x0, y0, x1, y1 int) {
	plotLineLow := func(x0, y0, x1, y1 int) {
		dx, dy, yi := x1-x0, y1-y0, 1
		if dy < 0 {
			yi, dy = -1, -dy
		}
		D, y := 2*dy-dx, y0
		for x := x0; x <= x1; x++ {
			ctx.Point(x, y)
			if D > 0 {
				y, D = y+yi, D-2*dx
			}
			D += 2 * dy
		}
	}
	plotLineHigh := func(x0, y0, x1, y1 int) {
		dx, dy, xi := x1-x0, y1-y0, 1
		if dx < 0 {
			xi, dx = -1, -dx
		}
		D, x := 2*dx-dy, x0
		for y := y0; y <= y1; y++ {
			ctx.Point(x, y)
			if D > 0 {
				x, D = x+xi, D-2*dy
			}
			D += 2 * dx
		}
	}

	if abs(y1-y0) < abs(x1-x0) {
		if x0 > x1 {
			plotLineLow(x1, y1, x0, y0)
		} else {
			plotLineLow(x0, y0, x1, y1)
		}
	} else {
		if y0 > y1 {
			plotLineHigh(x1, y1, x0, y0)
		} else {
			plotLineHigh(x0, y0, x1, y1)
		}
	}
}
