//go:build ignore

package main

import (
	"time"

	"github.com/Konstantin8105/ds"
	"github.com/go-gl/gl/v2.1/gl"
)

func main() {
	var ws [2]ds.Window
	ws[0].Draw = func() {
		gl.Begin(gl.QUADS)
		gl.Color3d(0.8, 0.1, 0.1)
		{
			gl.Vertex2d(-0.99, -0.99)
			gl.Vertex2d(-0.99, +0.99)
			gl.Vertex2d(+0.99, +0.99)
			gl.Vertex2d(+0.99, -0.99)
		}
		gl.End()
	}
	ws[1].Draw = func() {
		gl.Begin(gl.QUADS)
		gl.Color3d(0.1, 0.3, 0.9)
		{
			gl.Vertex2d(-0.99, -0.99)
			gl.Vertex2d(-0.99, +0.99)
			gl.Vertex2d(+0.99, +0.99)
			gl.Vertex2d(+0.99, -0.99)
		}
		gl.End()
	}

	go func() {
		for {
			t := time.Now().Second()
			ds.WindowRatio = float64(t) / 60.0
		}
	}()

	err := ds.New("Demo", ws)
	if err != nil {
		panic(err)
	}
}
