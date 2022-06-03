package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
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

var err error

func cmdTest() {
	cmd := exec.Command("ifconfig", "-a")
	grepCmd := exec.Command("grep", "utun5")

	// Get ps's stdout and attach it to grep's stdin.
	grepCmd.Stdin, err = cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	grepCmd.Stdout = os.Stdout

	if err := grepCmd.Start(); err != nil {
		fmt.Println(err)
		return
	}

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	if err := grepCmd.Wait(); err != nil {
		fmt.Println(err)
		return
	}
}

func shellTest() {
	cmd := exec.Command("sh", "-c", "ifconfig -a | grep utun5")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
