package ds

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var WindowRatio = 0.5

func init() {
	runtime.LockOSThread()
}

func New(name string, ds [2]Window) (err error) {
	//mutex
	var mutex sync.Mutex

	if err = glfw.Init(); err != nil {
		err = fmt.Errorf("failed to initialize glfw: %v", err)
		return
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	var window *glfw.Window
	window, err = glfw.CreateWindow(800, 600, name, nil, nil)
	if err != nil {
		return
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		return
	}

	glfw.SwapInterval(1) // Enable vsync

	defer func() {
		// 3D window is close
		glfw.Terminate()
	}()

	var w, h, xSplit int
	var focusIndex uint = 0

	window.SetCharCallback(func(w *glfw.Window, r rune) {
		//mutex
		mutex.Lock()
		defer mutex.Unlock()
		//action
		if f := ds[focusIndex].SetCharCallback; f != nil {
			f(r)
		}
	})

	window.SetScrollCallback(func(w *glfw.Window, xoffset, yoffset float64) {
		//mutex
		mutex.Lock()
		defer mutex.Unlock()
		//action
		x, _ := window.GetCursorPos()
		// split by windows
		if int(x) < xSplit {
			if f := ds[0].SetScrollCallback; f != nil {
				focusIndex = 0
				f(xoffset, yoffset)
			}
			return
		}
		if f := ds[1].SetScrollCallback; f != nil {
			focusIndex = 1
			f(xoffset, yoffset)
		}
	})

	window.SetMouseButtonCallback(func(
		w *glfw.Window,
		button glfw.MouseButton,
		action glfw.Action,
		mods glfw.ModifierKey,
	) {
		//mutex
		mutex.Lock()
		defer mutex.Unlock()
		//action
		x, _ := window.GetCursorPos()
		// split by windows
		if int(x) < xSplit {
			if f := ds[0].SetMouseButtonCallback; f != nil {
				focusIndex = 0
				f(button, action, mods, x)
			}
			return
		}
		if f := ds[1].SetMouseButtonCallback; f != nil {
			focusIndex = 1
			f(button, action, mods, x-float64(xSplit))
		}
	})

	for !window.ShouldClose() {
		// windows
		w, h = window.GetSize()
		xSplit = int(float64(w) * WindowRatio)
		// clean
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(1, 1, 1, 1)

		// prepare screen 0
		gl.Viewport(0, 0, int32(xSplit), int32(h))
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		if f := ds[0].Draw; f != nil {
			f()
		}

		// prepare screen 1
		gl.Viewport(int32(xSplit), 0, int32(w-xSplit), int32(h))
 		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		if f := ds[1].Draw; f != nil {
			f()
		}

		// end
		window.MakeContextCurrent()
		window.SwapBuffers()
	}

	return
}

type Window struct {
	SetMouseButtonCallback func(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, xLocal float64)
	SetCharCallback        func(r rune)
	SetScrollCallback      func(xoffset, yoffset float64)
	Draw                   func()
}
