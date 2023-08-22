package ds

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var windowRatio = 0.5

func init() {
	runtime.LockOSThread()
}

type Screen struct {
	ds           [2]Window
	w, h, xSplit int
	actions      *chan func()
	window       *glfw.Window
}

func (sc *Screen) ChangeRatio(newRatio float64) {
	if newRatio < 0 {
		return
	}
	if 1 < newRatio {
		return
	}
	if newRatio < 0.1 {
		newRatio = 0.1
	}
	if 0.9 < newRatio {
		newRatio = 0.9
	}
	(*sc.actions) <- func() {
		windowRatio = newRatio
		sc.initRatio()
	}
	return
}

func (sc *Screen) initRatio() {
	sc.w, sc.h = sc.window.GetSize()
	sc.xSplit = int(float64(sc.w) * windowRatio)
}

// New return windows.
// Minimal `actions = make(chan func(), 1000)`.
func New(name string, ds [2]Window, actions *chan func()) (sc *Screen, err error) {
	if actions == nil {
		err = fmt.Errorf("nil action channel")
		return
	}
	// initialization screen
	sc = new(Screen)
	sc.actions = actions
	sc.ds = ds

	//initialization gl
	if err = glfw.Init(); err != nil {
		err = fmt.Errorf("failed to initialize glfw: %v", err)
		return
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	sc.window, err = glfw.CreateWindow(800, 600, name, nil, nil)
	if err != nil {
		return
	}
	sc.window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		return
	}

	glfw.SwapInterval(1) // Enable vsync

	sc.initRatio()

	// var w, h, xSplit int
	var focusIndex uint = 0

	sc.window.SetCharCallback(func(w *glfw.Window, r rune) {
		//action
		if f := ds[focusIndex].SetCharCallback; f != nil {
			*actions <- func() { f(r) }
		}
	})

	sc.window.SetScrollCallback(func(w *glfw.Window, xoffset, yoffset float64) {
		//action
		x, _ := sc.window.GetCursorPos()
		// split by windows
		if int(x) < sc.xSplit {
			if f := ds[0].SetScrollCallback; f != nil {
				*actions <- func() {
					f(xoffset, yoffset)
					focusIndex = 0
				}
			}
			return
		}
		if f := ds[1].SetScrollCallback; f != nil {
			*actions <- func() {
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

	sc.window.SetMouseButtonCallback(func(
		w *glfw.Window,
		button glfw.MouseButton,
		action glfw.Action,
		mods glfw.ModifierKey,
	) {
		//action
		x, y := sc.window.GetCursorPos()
		// split by windows
		if int(x) < sc.xSplit {
			if f := ds[0].SetMouseButtonCallback; f != nil {
				*actions <- func() {
					f(button, action, mods, x, y)
					focusIndex = 0
				}
			}
			return
		}
		if f := ds[1].SetMouseButtonCallback; f != nil {
			*actions <- func() {
				f(button, action, mods, x-float64(sc.xSplit), y)
				focusIndex = 1
			}
		}
	})

	return
}

func (sc *Screen) Run() {
	defer func() {
		// 3D window is close
		glfw.Terminate()
	}()
	for !sc.window.ShouldClose() {
		// clean
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(1, 1, 1, 1)

		// prepare screen 0
		gl.Viewport(0, 0, int32(sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		if f := sc.ds[0].Draw; f != nil {
			f()
		}

		// prepare screen 1
		gl.Viewport(int32(sc.xSplit), 0, int32(sc.w-sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		if f := sc.ds[1].Draw; f != nil {
			f()
		}

		// actions func run
		select {
		case f := <-(*sc.actions):
			f()
		default:
		}

		// end
		sc.window.MakeContextCurrent()
		sc.window.SwapBuffers()
	}
	return
}

type Window interface {
	SetMouseButtonCallback(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, x, y float64)
	SetCharCallback(r rune)
	SetScrollCallback(xoffset, yoffset float64)
	Draw()
}
