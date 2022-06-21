package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/mojbro/gocoa"
)

var indicator *gocoa.ProgressIndicator

const maxValue = 100.00

func main() {
	gocoa.InitApplication()
	gocoa.OnApplicationDidFinishLaunching(func() {
		fmt.Println("App running!")
	})
	wnd := gocoa.NewWindow("Hello World!", 150, 150, 300, 200)
	quit := make(chan bool, 1)

	// Change me button
	// currentTitle, nextTitle := "Change me!", "Change me again!"
	// changeButton := gocoa.NewButton(75, 125, 150, 25)
	// changeButton.SetTitle(currentTitle)
	// changeButton.OnClick(func() {
	// 	changeButton.SetTitle(nextTitle)
	// 	currentTitle, nextTitle = nextTitle, currentTitle
	// })
	// wnd.AddButton(changeButton)

	// Start button
	startBtn := gocoa.NewButton(75, 125, 150, 25)
	startBtn.SetTitle("Start")
	startBtn.OnClick(func() {

		if runtime.NumGoroutine() < 2 {
			go vpnCheck(quit, quit)
		}

		// go func(ch chan bool) {
		// 	if safeCheck(ch) {
		// 		vpnCheck(quit, quit)
		// 	}
		// }(quit)
	})
	wnd.AddButton(startBtn)

	// Stop button
	stopBtn := gocoa.NewButton(75, 80, 150, 25)
	stopBtn.SetTitle("Stop")
	stopBtn.OnClick(func() {
		go func(ch chan bool) {
			if safeCheck(ch) {
				quit <- true
				//close(ch)
			}
		}(quit)
	})
	wnd.AddButton(stopBtn)

	// Quit button
	quitButton := gocoa.NewButton(75, 50, 150, 25)
	quitButton.SetTitle("Quit")
	quitButton.OnClick(func() { gocoa.TerminateApplication() })
	wnd.AddButton(quitButton)

	wnd.MakeKeyAndOrderFront()
	gocoa.RunApplication()
}

func vpnCheck(
	recvCh <-chan bool,
	sendCh chan<- bool) {

	for {
		select {
		case <-recvCh:
			fmt.Println("goroutine stop")
			return
		default:
			fmt.Println("goroutine running...")
			fmt.Println(runtime.NumGoroutine())
			time.Sleep(3000 * time.Millisecond)
			sendCh <- true
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func safeCheck(ch <-chan bool) bool {
	select {
	case <-ch:
		return false
	default:
	}
	return true
}

func channelTest() {
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				stop <- struct{}{}
				return
			default:
				testCheck()
			}

			time.Sleep(500 * time.Millisecond)
		}
	}(ctx)

	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()

	<-stop
	fmt.Println("finish")
}

func testCheck() {
	fmt.Println("test")
}

func wgTest() {
	var wait sync.WaitGroup
	wait.Add(3)

	go func() {
		defer wait.Done()
		fmt.Println("goroutine 1")
	}()

	go func() {
		defer wait.Done()
		fmt.Println("goroutine 2")
	}()

	go func() {
		defer wait.Done()
		fmt.Println("goroutine 3")
	}()

	wait.Wait()
}
