package ds

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var windowRatio = 0.5

func init() {
	runtime.LockOSThread()
}

type Action = func() (forceUpdateScreen bool)

type Screen struct {
	ds           [2]Window
	focusIndex   int
	w, h, xSplit int
	actions      *chan Action
	window       *glfw.Window
}

func (sc *Screen) UpdateWindow(pos int, w Window) {
	if sc.actions == nil {
		return
	}
	if w == nil {
		return
	}
	if pos < 0 {
		return
	}
	if 1 < pos {
		return
	}
	*sc.actions <- func() (forceUpdateScreen bool) {
		sc.ds[pos] = w
		return true
	}
}

func (sc *Screen) ChangeRatio(newRatio float64) {
	if newRatio < 0 {
		return
	}
	if math.Abs(newRatio-windowRatio) < 1e-6 {
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
	(*sc.actions) <- func() (forceUpdateScreen bool) {
		windowRatio = newRatio
		sc.initRatio()
		return true
	}
	return
}

func (sc *Screen) initRatio() {
	sc.w, sc.h = sc.window.GetSize()
	sc.xSplit = int(float64(sc.w) * windowRatio)
}

// New return windows.
// Minimal `actions = make(chan func(), 1000)`.
func New(
	name string,
	ds [2]Window,
	actions *chan Action,
) (
	sc *Screen,
	err error,
) {
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
		glfw.Terminate()
		return
	}

	// glfw.SwapInterval(1) // Enable vsync

	sc.initRatio()

	sc.focusIndex = 0 // default value

	// func (w *Window) SetCharCallback(cbfun CharCallback) (previous CharCallback)
	//     SetCharCallback sets the character callback which is called when a Unicode
	//     character is input.
	//
	//     The character callback is intended for Unicode text input. As it deals with
	//     characters, it is keyboard layout dependent, whereas the key callback is
	//     not. Characters do not map 1:1 to physical keys, as a key may produce zero,
	//     one or more characters. If you want to know whether a specific physical key
	//     was pressed or released, see the key callback instead.
	//
	//     The character callback behaves as system text input normally does and will
	//     not be called if modifier keys are held down that would prevent normal text
	//     input on that platform, for example a Super (Command) key on OS X or Alt key
	//     on Windows. There is a character with modifiers callback that receives these
	//     events.
	sc.window.SetCharCallback(func(w *glfw.Window, r rune) {
		//action
		if f := sc.ds[sc.focusIndex].SetCharCallback; f != nil {
			*actions <- func() (fus bool) {
				f(r)
				return false
			}
		}
	})

	// func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback)
	//     SetScrollCallback sets the scroll callback which is called when a scrolling
	//     device is used, such as a mouse wheel or scrolling area of a touchpad.
	sc.window.SetScrollCallback(func(w *glfw.Window, xoffset, yoffset float64) {
		//action
		x, y := sc.window.GetCursorPos()
		// split by windows
		if int(x) < sc.xSplit {
			if f := sc.ds[0].SetScrollCallback; f != nil {
				*actions <- func() (fus bool) {
					f(x, y, xoffset, yoffset)
					sc.focusIndex = 0
					return false
				}
			}
			return
		}
		if f := sc.ds[1].SetScrollCallback; f != nil {
			*actions <- func() (fus bool) {
				f(x-float64(sc.xSplit), y, xoffset, yoffset)
				sc.focusIndex = 1
				return false
			}
		}
	})

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
	sc.window.SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			//action
			if f := sc.ds[sc.focusIndex].SetKeyCallback; f != nil {
				*actions <- func() (fus bool) {
					f(key, scancode, action, mods)
					return false
				}
			}
		})

	// TODO:
	// func (w *Window) SetCursorEnterCallback(cbfun CursorEnterCallback) (previous CursorEnterCallback)
	//     SetCursorEnterCallback the cursor boundary crossing callback which is called
	//     when the cursor enters or leaves the client area of the window.

	// func SetClipboardString(str string)
	//     SetClipboardString sets the system clipboard to the specified UTF-8 encoded
	//     string.
	//
	//     This function may only be called from the main thread.

	// func (w *Window) SetCursorPosCallback(cbfun CursorPosCallback) (previous CursorPosCallback)
	//
	//	SetCursorPosCallback sets the cursor position callback which is called when
	//	the cursor is moved. The callback is provided with the position relative to
	//	the upper-left corner of the client area of the window.
	sc.window.SetCursorPosCallback(
		func(w *glfw.Window, xpos, ypos float64) {
			// action
			if sc.focusIndex == 1 {
				xpos = xpos - float64(sc.xSplit)
			}
			if f := sc.ds[sc.focusIndex].SetCursorPosCallback; f != nil {
				*actions <- func() (fus bool) {
					f(xpos, ypos)
					return false
				}
			}
		})

	// func (w *Window) SetMouseButtonCallback(cbfun MouseButtonCallback) (previous MouseButtonCallback)
	//     SetMouseButtonCallback sets the mouse button callback which is called when a
	//     mouse button is pressed or released.
	//
	//     When a window loses focus, it will generate synthetic mouse button release
	//     events for all pressed mouse buttons. You can tell these events from
	//     user-generated events by the fact that the synthetic ones are generated
	//     after the window has lost focus, i.e. Focused will be false and the focus
	//     callback will have already been called.
	sc.window.SetMouseButtonCallback(func(
		w *glfw.Window,
		button glfw.MouseButton,
		action glfw.Action,
		mods glfw.ModifierKey,
	) {
		//action
		x, y := sc.window.GetCursorPos()
		switch action {
		case glfw.Press: // The key or button was pressed.
			if x < float64(sc.xSplit) {
				sc.focusIndex = 0
			} else {
				sc.focusIndex = 1
				x = x - float64(sc.xSplit)
			}
			if f := sc.ds[sc.focusIndex].SetMouseButtonCallback; f != nil {
				*actions <- func() (fus bool) {
					f(button, action, mods, x, y)
					return false
				}
			}
		default:
			// The key or button was released.
			// case glfw.Release:
			if sc.focusIndex == 1 {
				x = x - float64(sc.xSplit)
			}
			if f := sc.ds[sc.focusIndex].SetMouseButtonCallback; f != nil {
				*actions <- func() (fus bool) {
					f(button, action, mods, x, y)
					return false
				}
			}
		}
	})

	return
}

func (sc *Screen) Screenshot(afterSave func(img image.Image)) {
	*sc.actions <- func() bool { return true }
	*sc.actions <- func() bool {
		// flush opengl
		gl.Flush()
		gl.Finish()
		// get pixels from opengl
		sizeX := sc.w
		sizeY := sc.h
		size := sizeX * sizeY
		data := make([]uint8, 4*size)
		gl.ReadPixels(0, 0, int32(sizeX), int32(sizeY),
			gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&data[0]))
		// create picture
		img := image.NewNRGBA(image.Rect(0, 0, sizeX, sizeY))
		for y := 0; y < sizeY; y++ {
			for x := 0; x < sizeX; x++ {
				c := data[4*(x+(sizeY-1-y)*sizeX):]
				img.Set(x, y, color.NRGBA{R: c[0], G: c[1], B: c[2], A: c[3]})
			}
		}
		// run after save
		if f := afterSave; f != nil {
			f(img)
		}
		// update screen
		return true
	}
}

func (sc *Screen) Run(quit *chan struct{}) {
	defer func() {
		sc.window.Destroy()
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
		gl.Viewport(0, 0, int32(sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		if f := sc.ds[0].Draw; f != nil {
			f(0, 0, int32(sc.xSplit), int32(sc.h))
		}

		// prepare screen 1
		gl.Viewport(int32(sc.xSplit), 0, int32(sc.w-sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Viewport(int32(sc.xSplit), 0, int32(sc.w-sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		if f := sc.ds[1].Draw; f != nil {
			f(int32(sc.xSplit), 0, int32(sc.w-sc.xSplit), int32(sc.h))
		}

		// separator
		gl.Viewport(int32(sc.xSplit), 0, int32(1), int32(sc.h))
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Viewport(int32(sc.xSplit), 0, int32(1), int32(sc.h))
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		{
			gl.LineWidth(1)
			gl.Color3f(0.7, 0.7, 0.7)
			gl.Begin(gl.LINES)
			gl.Vertex2i(0, -1)
			gl.Vertex2i(0, +1)
			gl.End()
		}

		// end
		sc.window.MakeContextCurrent()
		sc.window.SwapBuffers()

		// actions
		// run first funcs
		for i, size := 0, 50; i < size; i++ {
			forceUpdateScreen := false
			select {
			case f, ok := <-(*sc.actions):
				if !ok {
					// probably closed channel
					break
				}
				// TODO: if action time long for example 10 minutes,
				// then screen is freeze.
				forceUpdateScreen = f()
				if forceUpdateScreen {
					break
				}
			default:
				break
			}
			if forceUpdateScreen {
				break
			}
		}

		// quit
		isQuit := false
		select {
		case _, ok := <-*quit:
			if !ok {
				isQuit = true
			}
		default:
		}
		if isQuit {
			break
		}

		// update ratio
		sc.initRatio()
	}
	return
}

type Window interface {
	SetMouseButtonCallback(
		button glfw.MouseButton,
		action glfw.Action,
		mods glfw.ModifierKey,
		xcursor, ycursor float64,
	)
	SetCharCallback(r rune)
	SetScrollCallback(
		xcursor, ycursor float64,
		xoffset, yoffset float64,
	)
	SetKeyCallback(
		key glfw.Key,
		scancode int,
		action glfw.Action,
		mods glfw.ModifierKey,
	)
	SetCursorPosCallback(
		xpos float64,
		ypos float64,
	)
	Draw(x, y, w, h int32)
}
