package ds

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var WindowRatio = 0.5

func init() {
	runtime.LockOSThread()
}

// New return windows.
// Minimal `actions = make(chan func(), 1000)`.
//
func New(name string, ds [2]Window, actions chan func()) (err error) {
	//initialization
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
		//action
		if f := ds[focusIndex].SetCharCallback; f != nil {
			actions <- func() { f(r) }
		}
	})

	window.SetScrollCallback(func(w *glfw.Window, xoffset, yoffset float64) {
		//action
		x, _ := window.GetCursorPos()
		// split by windows
		if int(x) < xSplit {
			if f := ds[0].SetScrollCallback; f != nil {
				actions <- func() {
					f(xoffset, yoffset)
					focusIndex = 0
				}
			}
			return
		}
		if f := ds[1].SetScrollCallback; f != nil {
			actions <- func() {
				f(xoffset, yoffset)
				focusIndex = 1
			}
		}
	})

	// TODO:

	// func (w *Window) SetKeyCallback(cbfun KeyCallback) (previous KeyCallback)
	//     SetKeyCallback sets the key callback which is called when a key is pressed,
	//     repeated or released.
	//
	//     The key functions deal with physical keys, with layout independent key
	//     tokens named after their values in the standard US keyboard layout. If you
	//     want to input text, use the SetCharCallback instead.
	//
	//     When a window loses focus, it will generate synthetic key release events
	//     for all pressed keys. You can tell these events from user-generated events
	//     by the fact that the synthetic ones are generated after the window has lost
	//     focus, i.e. Focused will be false and the focus callback will have already
	//     been called.

	// func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback)
	//     SetCursorEnterCallback the cursor boundary crossing callback which is called
	//     when the cursor enters or leaves the client area of the window.

	// func SetClipboardString(str string)
	//     SetClipboardString sets the system clipboard to the specified UTF-8 encoded
	//     string.
	//
	//     This function may only be called from the main thread.

	window.SetMouseButtonCallback(func(
		w *glfw.Window,
		button glfw.MouseButton,
		action glfw.Action,
		mods glfw.ModifierKey,
	) {
		//action
		x, y := window.GetCursorPos()
		// split by windows
		if int(x) < xSplit {
			if f := ds[0].SetMouseButtonCallback; f != nil {
				actions <- func() {
					f(button, action, mods, x, y)
					focusIndex = 0
				}
			}
			return
		}
		if f := ds[1].SetMouseButtonCallback; f != nil {
			actions <- func() {
				f(button, action, mods, x-float64(xSplit), y)
				focusIndex = 1
			}
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

		// actions func run
		select {
		case f := <-actions:
			f()
		default:
		}

		// end
		window.MakeContextCurrent()
		window.SwapBuffers()
	}

	return
}

type Window interface {
	SetMouseButtonCallback(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, x, y float64)
	SetCharCallback(r rune)
	SetScrollCallback(xoffset, yoffset float64)
	Draw()
}
