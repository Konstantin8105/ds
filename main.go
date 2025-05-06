//go:build ignore

package main

import (
	"fmt"
	"image"
	"os"
	"sync"
	"time"

	"github.com/Konstantin8105/compare"
	"github.com/Konstantin8105/ds"
)

type checker struct {
	iserror bool
	err     error
}

func (c *checker) Errorf(format string, args ...any) {
	c.iserror = true
	c.err = fmt.Errorf(format, args...)
}

func main() {
	var (
		ws [2]ds.Window
		ch = make(chan func() (fus bool), 1000)
		tr = ds.DemoSpiral{}
		d3 = ds.DemoCube{}
	)
	// create screen
	ws[0] = &tr
	ws[1] = &d3
	var screen *ds.Screen
	// quit channel
	quit := make(chan struct{})
	screen, err := ds.New("Demo", ws, &ch)
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(time.Second)
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
			time.Sleep(500 * time.Millisecond)
			fmt.Fprintf(os.Stdout, "Test: %01d\n", it)
			tc()
			var wg sync.WaitGroup
			wg.Add(1)
			screen.Screenshot(func(img image.Image) {
				defer wg.Done()
				var t checker
				compare.TestPng(&t, fmt.Sprintf("test.%02d.png", it), img)
				if t.iserror {
					fmt.Fprintf(os.Stdout, "TODO: Not clear error: %v\n", t.err)
				}
				fmt.Println(t)
			})
			wg.Wait()
		}
		fmt.Fprintf(os.Stdout, "Test: free move\n")
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
		time.Sleep(5 * time.Second)
		close(quit)
	}()
	screen.Run(&quit)
}
