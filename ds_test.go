package ds_test

import (
	"image"
	"testing"
	"time"

	"github.com/Konstantin8105/compare"
	"github.com/Konstantin8105/ds"
)

func Test(t *testing.T) {
	var (
		ws [2]ds.Window
		ch = make(chan func() (fus bool), 1000)
		tr = ds.NewDemoSpiral(1)
		d3 = ds.DemoCube{
			Alpha: 10,
			Betta: 10,
		}
	)
	// create screen
	ws[0] = &tr
	ws[1] = &d3
	screen, err := ds.New("Demo", ws, &ch)
	if err != nil {
		t.Fatal(err)
	}
	// quit channel
	quit := make(chan struct{})
	pause := func() {
		time.Sleep(500 * time.Millisecond)
	}
	go func() {
		screen.Run(&quit)
	}()
	pause()
	screen.Screenshot(func(img image.Image) {
		compare.TestPng(t, "test.00.png", img)
	})
	pause()
	ch <- func() (fus bool) {
		d3.Alpha = 20
		return false
	}
	pause()
	screen.Screenshot(func(img image.Image) {
		compare.TestPng(t, "test.01.png", img)
	})
	pause()
	ch <- func() (fus bool) {
		screen.ChangeRatio(0.3)
		return false
	}
	pause()
	screen.Screenshot(func(img image.Image) {
		compare.TestPng(t, "test.02.png", img)
	})
	pause()
	go func() {
		var t float64
		for {
			ch <- func() (fus bool) {
				t += 0.05
				d3.Alpha = 360 * t
				return false
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()
	go func() {
		for {
			ch <- func() (fus bool) {
				t := time.Now().Second()
				screen.ChangeRatio(float64(t)/60.0*0.8 + 0.1)
				return false
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	time.Sleep(50 * time.Second)
	close(quit)
}
