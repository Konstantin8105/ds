//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Konstantin8105/ds"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Triangles struct {
	points [][4][2]float64
	color  float64
}

func (o *Triangles) SetMouseButtonCallback(
	button glfw.MouseButton,
	action glfw.Action,
	mods glfw.ModifierKey,
	x, y float64,
) {
	fmt.Fprintf(os.Stdout, "Click on window 0:[%v,%v]\n", x, y)
}
func (o *Triangles) SetCharCallback(r rune) {
}
func (o *Triangles) SetScrollCallback(xcursor, ycursor float64, xoffset, yoffset float64) {
}
func (o *Triangles) SetCursorPosCallback(xpos float64, ypos float64) {
}
func (o *Triangles) SetKeyCallback(
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	fmt.Fprintf(os.Stdout, "Key on window 0: %v\n", key)
}
func (o *Triangles) Draw(x, y, w, h int32) {
	size := 100
	if len(o.points) == 0 {
		if true {
			// spiral triangles
			base := [][2]float64{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}
			o.points = make([][4][2]float64, size)
			for range size {
				o.points = append(o.points,
					[4][2]float64{base[0], base[1], base[2], base[3]})
				for j := range 4 {
					for k := range 2 {
						base[j][k] = base[j][k] + (base[j+1][k]-base[j][k])/
							(0.2*float64(size))
					}
				}
				base[4] = base[0]
			}
		} else {
			// random
			// return from -1 to +1
			get := func() float64 {
				return (rand.Float64() * 2) - 1
			}
			o.points = make([][4][2]float64, size)
			for i := 0; i < size; i++ {
				for j := 0; j < 4; j++ {
					o.points[i][j] = [2]float64{get(), get()}
				}
			}
		}
	}
	for i, ps := range o.points {
		gl.Begin(gl.QUADS)
		p := float64(i) / float64(size)
		if i%2 == 0 {
			p = float64(size-i) / float64(size)
		}
		gl.Color4d(0.8, p, o.color, 0.5)
		for _, p := range ps {
			gl.Vertex2d(p[0], p[1])
		}
		gl.End()
	}
}

type D3 struct {
	alpha, betta float64
	actions      *chan ds.Action
}

func (o *D3) SetMouseButtonCallback(
	button glfw.MouseButton,
	action glfw.Action,
	mods glfw.ModifierKey,
	x, y float64,
) {
	*o.actions <- func() (fus bool) {
		fmt.Fprintf(os.Stdout, "Click on window 1:[%v,%v]\n", x, y)
		return false
	}
}
func (o *D3) SetCharCallback(r rune) {
}
func (o *D3) SetCursorPosCallback(xpos float64, ypos float64) {
}
func (o *D3) SetScrollCallback(xcursor, ycursor float64, xoffset, yoffset float64) {
}
func (o *D3) SetKeyCallback(
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey,
) {
	fmt.Fprintf(os.Stdout, "Key on window 1: %v\n", key)
}
func (o *D3) Draw(x, y, w, h int32) {
	gl.Viewport(int32(x), int32(y), int32(w), int32(h))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	var ratio float64
	ratio = float64(w) / float64(h)
	ymax := 0.2 * 8000
	scale := 1.0
	gl.Ortho(-scale*ratio, scale*ratio, -scale, scale, float64(-ymax), float64(ymax))

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Translated(0, 0, 0) // TODO ?
	gl.Rotated(o.betta, 1.0, 0.0, 0.0)
	gl.Rotated(o.alpha, 0.0, 1.0, 0.0)

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

func main() {
	var ws [2]ds.Window
	ch := make(chan func() (fus bool), 1000)

	tr := Triangles{color: float64(1)}
	ws[0] = &tr

	d3 := D3{actions: &ch}
	ws[1] = &d3

	screen, err := ds.New("Demo", ws, &ch)
	if err != nil {
		panic(err)
	}

	quit := make(chan struct{})

	d3.betta = 10
	go func() {
		var t float64
		for {
			ch <- func() (fus bool) {
				// t := time.Now().Second()
				// d3.alpha = 360 * float64(t) / 60
				// d3.betta = 360 * float64(t) / 60
				t += 0.05
				d3.alpha = 360 * t
				//d3.betta = 360 * t
				return true
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			ch <- func() (fus bool) {
				t := time.Now().Second()
				screen.ChangeRatio(float64(t)/60*0.8 + 0.1)
				return true
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		time.Sleep(50 * time.Second)
		close(quit)
	}()

	screen.Run(&quit)
}
