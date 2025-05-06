package ds_test

import (
	"testing"
	"time"

	"github.com/Konstantin8105/ds"
)

func Test(t *testing.T) {
	var ws [2]ds.Window
	ch := make(chan func() (fus bool), 1000)

	tr := ds.NewDemoSpiral(1)
	ws[0] = &tr

	d3 := ds.DemoCube{}
	ws[1] = &d3

	screen, err := ds.New("Demo", ws, &ch)
	if err != nil {
		t.Fatal(err)
	}

	quit := make(chan struct{})

	d3.Betta = 10
	go func() {
		var t float64
		for {
			ch <- func() (fus bool) {
				// t := time.Now().Second()
				// d3.alpha = 360 * float64(t) / 60
				// d3.betta = 360 * float64(t) / 60
				t += 0.05
				d3.Alpha = 360 * t
				//d3.betta = 360 * t
				return false // true
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()

	go func() {
		for {
			ch <- func() (fus bool) {
				t := time.Now().Second()
				screen.ChangeRatio(float64(t)/60.0*0.8 + 0.1)
				return false // true
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		time.Sleep(50 * time.Second)
		close(quit)
	}()

	screen.Run(&quit)
}
