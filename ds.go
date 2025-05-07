package ds

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

type Action = func() (forceUpdateScreen bool)

type Screen struct {
	ds           [2]Window
	focusIndex   int
	w, h, xSplit int
	windowRatio  float64
	actions      *chan Action
	window       *glfw.Window
}

// func (sc *Screen) UpdateWindow(pos int, w Window) {
// 	if sc.actions == nil {
// 		return
// 	}
// 	if w == nil {
// 		return
// 	}
// 	if pos < 0 {
// 		return
// 	}
// 	if 1 < pos {
// 		return
// 	}
// 	*sc.actions <- func() (forceUpdateScreen bool) {
// 		sc.ds[pos] = w
// 		return true
// 	}
// }

func (sc *Screen) ChangeRatio(ratio float64) {
	ratio = math.Min(math.Max(ratio, 0.1), 0.9)
	// Acceptable same ratio, for example:
	// * after update window sizes
	// * update screen
	w, h := sc.window.GetSize()
	if sc.w == w && sc.h == h && math.Abs(ratio-sc.windowRatio) < 1e-6 {
		return
	}
	(*sc.actions) <- func() (forceUpdateScreen bool) {
		sc.windowRatio = ratio
		sc.w, sc.h = sc.window.GetSize()
		sc.xSplit = int(float64(sc.w) * sc.windowRatio)
		return true
	}
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

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Visible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	// glfw.WindowHint(glfw.ContextCreationAPI, glfw.NativeContextAPI)
	// glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLAPI)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	// glfw.WindowHint(glfw.Samples, 4) // smooth

	sc.window, err = glfw.CreateWindow(800, 600, name, nil, nil)
	if err != nil {
		return
	}
	// move function MakeContextCurrent from loop for avoid problem:
	//
	// X Error of failed request:  BadAccess (attempt to access private resource denied)
	// Major opcode of failed request:  153 (GLX)
	// Minor opcode of failed request:  5 (X_GLXMakeCurrent)
	// Serial number of failed request:  197
	// Current serial number in output stream:  197
	sc.window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err = gl.Init(); err != nil {
		glfw.Terminate()
		return
	}

	defer func() {
		sc.ChangeRatio(0.5)
	}()

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
		f := sc.ds[sc.focusIndex].SetCharCallback
		{
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
			f := sc.ds[0].SetScrollCallback
			*actions <- func() (fus bool) {
				f(x, y, xoffset, yoffset)
				sc.focusIndex = 0
				return false
			}
			return
		}
		{
			f := sc.ds[1].SetScrollCallback
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
			f := sc.ds[sc.focusIndex].SetKeyCallback
			*actions <- func() (fus bool) {
				f(key, scancode, action, mods)
				return false
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
			{
				f := sc.ds[sc.focusIndex].SetCursorPosCallback
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
			{
				f := sc.ds[sc.focusIndex].SetMouseButtonCallback
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
			{
				f := sc.ds[sc.focusIndex].SetMouseButtonCallback
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
		// get pixels from opengl
		sizeX := sc.w
		sizeY := sc.h
		size := sizeX * sizeY
		img := image.NewNRGBA(image.Rect(0, 0, sizeX, sizeY))
		if 0 < size {
			gl.Finish()
			gl.Flush()
			data := make([]uint8, 4*size)
			gl.ReadPixels(0, 0, int32(sizeX), int32(sizeY),
				gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&data[0]))
			// debug information
			{
				empty := true
				for _, d := range data {
					if d != 0 {
						empty = false
					}
				}
				if empty {
					fmt.Fprintf(os.Stdout, "Screenshoot is empty\n")
				}
			}
			// create picture
			for y := 0; y < sizeY; y++ {
				for x := 0; x < sizeX; x++ {
					c := data[4*(x+(sizeY-1-y)*sizeX):]
					img.Set(x, y, color.NRGBA{R: c[0], G: c[1], B: c[2], A: c[3]})
				}
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
		// update ratio after window size change
		sc.ChangeRatio(sc.windowRatio)
		// events
		glfw.PollEvents()
		// clean
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(1, 1, 1, 1)
		// gl.ClearDepth(1)
		// gl.DepthFunc(gl.LEQUAL)
		// prepare screen 0
		gl.Viewport(0, 0, int32(sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Viewport(0, 0, int32(sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		{
			f := sc.ds[0].Draw
			f(0, 0, int32(sc.xSplit), int32(sc.h))
		}
		// prepare screen 1
		gl.Viewport(int32(sc.xSplit), 0, int32(sc.w-sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		gl.Viewport(int32(sc.xSplit), 0, int32(sc.w-sc.xSplit), int32(sc.h))
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		{
			f := sc.ds[1].Draw
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
		sc.window.SwapBuffers()

		// actions
		// run first funcs
		for i, size := 0, 50; i < size; i++ {
			var reset bool
			select {
			case f, ok := <-(*sc.actions):
				if !ok {
					// probably closed channel
					break
				}
				// if action time long for example 10 minutes,
				// then screen is freeze.
				reset = f() // forceUpdateScreen
			default:
				reset = true
				break
			}
			if reset {
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
	}
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
