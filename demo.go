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

func main() {
	get := func() float64 {
		return (rand.Float64() * 2) - 1
	}
	size := 50
	var ws [2]ds.Window
	ws[0].Draw = func() {
		for i := 0; i < size; i++ {
			gl.Begin(gl.QUADS)
			gl.Color4d(0.8, float64(i)/float64(size), 0.1, 0.5)
			{
				gl.Vertex2d(-get(), -get())
				gl.Vertex2d(-get(), +get())
				gl.Vertex2d(+get(), +get())
				gl.Vertex2d(+get(), -get())
			}
			gl.End()
		}
	}
	ws[0].SetMouseButtonCallback = func(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, x, y float64) {
		fmt.Fprintf(os.Stdout, "Click on window 0:[%v,%v]\n", x, y)
	}
	ws[1].Draw = func() {
		for i := 0; i < size; i++ {
			gl.Begin(gl.QUADS)
			gl.Color4d(0.1, 0.3, float64(i)/float64(size), 0.5)
			{
				gl.Vertex2d(-get(), -get())
				gl.Vertex2d(-get(), +get())
				gl.Vertex2d(+get(), +get())
				gl.Vertex2d(+get(), -get())
			}
			gl.End()
		}
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

	err := ds.New("Demo", ws, make(chan func(), 1000))
	if err != nil {
		panic(err)
	}
}
