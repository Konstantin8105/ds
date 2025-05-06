package ds

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type DemoSpiral struct {
	points [40][4][2]float64
	init   bool
}

func (o *DemoSpiral) SetMouseButtonCallback(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, x, y float64) {
	fmt.Fprintf(os.Stdout, "Click on window DemoSpiral:[%v,%v]\n", x, y)
}
func (o *DemoSpiral) SetCharCallback(r rune)                                               {}
func (o *DemoSpiral) SetScrollCallback(xcursor, ycursor float64, xoffset, yoffset float64) {}
func (o *DemoSpiral) SetCursorPosCallback(xpos float64, ypos float64)                      {}
func (o *DemoSpiral) SetKeyCallback(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Fprintf(os.Stdout, "Key on window DemoSpiral: %v\n", key)
}
func (o *DemoSpiral) Draw(x, y, w, h int32) {
	size := len(o.points)
	if !o.init {
		o.init = true
		base := [][2]float64{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}
		for i := 0; i < size; i++ {
			o.points[i] = [4][2]float64{base[0], base[1], base[2], base[3]}
			for j := 0; j < 4; j++ {
				for k := 0; k < 2; k++ {
					base[j][k] = base[j][k] + (base[j+1][k]-base[j][k])/
						(0.2*float64(size))
				}
			}
			base[4] = base[0]
		}
	}
	for i, ps := range o.points {
		green := float64(i) / float64(size)
		if i%2 == 0 {
			green = float64(size-i) / float64(size)
		}
		blue := green
		gl.Color4d(0.8, green, blue, 0.5)
		gl.Begin(gl.QUADS)
		for _, p := range ps {
			gl.Vertex2d(p[0], p[1])
		}
		gl.End()
	}
}

type DemoCube struct {
	Alpha, Betta float64
}

func (o *DemoCube) SetMouseButtonCallback(
	button glfw.MouseButton,
	action glfw.Action,
	mods glfw.ModifierKey,
	x, y float64,
) {
	fmt.Fprintf(os.Stdout, "Click on window DemoCube:[%v,%v]\n", x, y)
}
func (o *DemoCube) SetCharCallback(r rune)                                               {}
func (o *DemoCube) SetCursorPosCallback(xpos float64, ypos float64)                      {}
func (o *DemoCube) SetScrollCallback(xcursor, ycursor float64, xoffset, yoffset float64) {}
func (o *DemoCube) SetKeyCallback(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Fprintf(os.Stdout, "Key on window DemoCube: %v\n", key)
}
func (o *DemoCube) Draw(x, y, w, h int32) {
	gl.Viewport(int32(x), int32(y), int32(w), int32(h))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	ratio := float64(w) / float64(h)
	ymax := 0.2 * 8000
	scale := 1.0
	gl.Ortho(-scale*ratio, scale*ratio, -scale, scale, float64(-ymax), float64(ymax))

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Translated(0, 0, 0) // TODO ?
	gl.Rotated(o.Betta, 1.0, 0.0, 0.0)
	gl.Rotated(o.Alpha, 0.0, 1.0, 0.0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	defer func() {
		gl.DepthFunc(gl.LESS)
		gl.Disable(gl.DEPTH_TEST)
	}()

	// cube
	size := 0.1
	gl.Color3d(0.1, 0.7, 0.1)
	gl.Begin(gl.LINES)
	{
		gl.Vertex3d(-size, -size, -size)
		gl.Vertex3d(+size, -size, -size)

		gl.Vertex3d(-size, -size, -size)
		gl.Vertex3d(-size, +size, -size)

		gl.Vertex3d(-size, -size, -size)
		gl.Vertex3d(-size, -size, +size)

		gl.Vertex3d(+size, +size, +size)
		gl.Vertex3d(-size, +size, +size)

		gl.Vertex3d(+size, +size, +size)
		gl.Vertex3d(+size, -size, +size)

		gl.Vertex3d(+size, +size, +size)
		gl.Vertex3d(+size, +size, -size)
	}
	gl.End()

	gl.PointSize(5)
	gl.Color3d(0.2, 0.8, 0.5)
	gl.Begin(gl.POINTS)
	{
		gl.Vertex3d(-size, -size, -size)
		gl.Vertex3d(+size, -size, -size)
		gl.Vertex3d(-size, +size, -size)
		gl.Vertex3d(-size, -size, +size)
		gl.Vertex3d(+size, +size, -size)
		gl.Vertex3d(+size, -size, +size)
		gl.Vertex3d(-size, +size, +size)
		gl.Vertex3d(+size, +size, +size)
	}
	gl.End()

	gl.Color3d(0.5, 0.8, 0.2)
	gl.Begin(gl.TRIANGLES)
	{
		gl.Vertex3d(-size, -size, -size)
		gl.Vertex3d(+size, -size, -size)
		gl.Vertex3d(-size, +size, -size)
	}
	gl.End()

	gl.Color3d(0.8, 0.2, 0.5)
	gl.Begin(gl.TRIANGLES)
	{
		gl.Vertex3d(-size, -size, +size)
		gl.Vertex3d(+size, +size, -size)
		gl.Vertex3d(+size, -size, +size)
	}
	gl.End()
}
