package ds_test

import (
	"image"

	"fmt"
	"testing"
	"time"

	"github.com/Konstantin8105/ds"
)

func Test(t *testing.T) {
	var (
		ws [2]ds.Window
		ch = make(chan func() (fus bool), 1000)
		tr = ds.DemoSpiral{}
		d3 = ds.DemoCube{}
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
	go func() {
		screen.Run(&quit)
	}()
	for it, tc := range []func(){
		func() {
			screen.ChangeRatio(0.5)
			d3.Alpha = 10.0
			d3.Betta = 10.0
		},
		func() {
			screen.ChangeRatio(0.5)
			d3.Alpha = 10.0
			d3.Betta = 30.0
		},
		func() {
			screen.ChangeRatio(0.3)
			d3.Alpha = 10.0
			d3.Betta = 30.0
		},
	} {
		t.Run(fmt.Sprintf("%02d", it), func(t *testing.T) {
			time.Sleep(500 * time.Millisecond)
			tc()
			screen.Screenshot(func(img image.Image) {
				//compare.TestPng(t, fmt.Sprintf("test.%02d.png", it), img)
			})
		})
	}
	t.Run("Example", func(t *testing.T) {
		d3.Alpha = 10.0
		d3.Betta = 10.0
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
	})
	close(quit)
}
