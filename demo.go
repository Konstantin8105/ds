//go:build ignore

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Konstantin8105/ds"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
	ws[0].SetMouseButtonCallback = func(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, x, y float64) {
		fmt.Fprintf(os.Stdout, "Click on window 0:[%v,%v]\n", x, y)
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
	ws[1].SetMouseButtonCallback = func(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, x, y float64) {
		fmt.Fprintf(os.Stdout, "Click on window 1:[%v,%v]\n", x, y)
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
