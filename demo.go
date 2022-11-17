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

type Opengl struct {
	color float64
}

func (o *Opengl) SetMouseButtonCallback(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, x, y float64) {
	fmt.Fprintf(os.Stdout, "Click on window 0:[%v,%v]\n", x, y)
}
func (o *Opengl) SetCharCallback(r rune) {
}
func (o *Opengl) SetScrollCallback(xoffset, yoffset float64) {
}
func (o *Opengl) Draw() {
	get := func() float64 {
		return (rand.Float64() * 2) - 1
	}
	size := 50
	for i := 0; i < size; i++ {
		gl.Begin(gl.QUADS)
		gl.Color4d(0.8, float64(i)/float64(size), o.color, 0.5)
		{
			gl.Vertex2d(-get(), -get())
			gl.Vertex2d(-get(), +get())
			gl.Vertex2d(+get(), +get())
			gl.Vertex2d(+get(), -get())
		}
		gl.End()
	}
}

func main() {
	var ws [2]ds.Window
	for i := range ws {
		o := new(Opengl)
		o.color = float64(i)
		ws[i] = o
	}
	screen, err := ds.New("Demo", ws, make(chan func(), 1000))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			t := time.Now().Second()
			screen.ChangeRatio(float64(t) / 60)
			time.Sleep(time.Second)
		}
	}()

	screen.Run()
}
