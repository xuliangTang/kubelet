package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

// Until: 根据 channel 的关闭或者 context Done 的信号来结束对指定函数的轮询操作
// Poll：不只是会根据 channel 或者 context 来决定结束轮询，还会判断轮询函数的返回值来决定是否结束
// Wait: 会根据 WaitFor 函数返回的 channel 来触发函数执行
// Backoff：会根据 Backoff 返回的时间间隔来循环触发函数的执行

func main() {
	//testForever()
	//testUntil()
	//testWaitFor()
	testBackoff()
}

func testForever() {
	wait.Forever(func() {
		fmt.Println(time.Now().Unix())
	}, time.Second*1)
}

func testUntil() {
	// 3秒后结束
	ch := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 3)
		ch <- struct{}{}
	}()

	wait.Until(func() {
		fmt.Println(time.Now().Unix())
	}, time.Second*1, ch)
}

func testWaitFor() {
	ch := make(chan struct{})
	wait.WaitFor(func(done <-chan struct{}) <-chan struct{} {
		c := make(chan struct{})
		// 每隔2秒返回一次
		go func() {
			for {
				time.Sleep(time.Second * 2)
				c <- struct{}{}
			}
		}()
		return c
	}, func() (done bool, err error) {
		// 当读取到c的值时被触发
		fmt.Println(time.Now().Unix())
		// 返回true或err时终止
		return false, nil
	}, ch)
}

func testBackoff() {
	ch := make(chan struct{})
	wait.BackoffUntil(func() {
		fmt.Println(time.Now().Unix())
	}, wait.NewJitteredBackoffManager(time.Second*1, 0.0, &clock.RealClock{}), true, ch)
}
