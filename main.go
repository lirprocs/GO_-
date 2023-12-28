package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	var wg sync.WaitGroup
	list1 := make([]int, n)
	list2 := make([]int, n)
	wg.Add(2 * n)
	go func() {
		for i := 0; i < n; i++ {
			x1 := <-in1
			go func(p int) {
				defer wg.Done()
				list1[p] = fn(x1)
			}(i)
		}
	}()
	go func() {
		for i := 0; i < n; i++ {
			x2 := <-in2
			go func(p int) {
				defer wg.Done()
				list2[p] = fn(x2)
			}(i)
		}

	}()
	go func() {
		wg.Wait()
		for i := 0; i < n; i++ {
			out <- list1[i] + list2[i]
		}
	}()
}

const N = 5

func main() {

	fn := func(x int) int {
		time.Sleep(time.Duration(rand.Int31n(N)) * time.Second)
		return x * 2
	}
	in1 := make(chan int, N)
	in2 := make(chan int, N)
	out := make(chan int, N)

	start := time.Now()
	merge2Channels(fn, in1, in2, out, N+1)
	for i := 0; i < N+1; i++ {
		in1 <- i
		in2 <- i
	}

	orderFail := false
	EvenFail := false
	for i, prev := 0, 0; i < N; i++ {
		c := <-out
		if c%2 != 0 {
			EvenFail = true
		}
		if prev >= c && i != 0 {
			orderFail = true
		}
		prev = c
		fmt.Println(c)
	}
	if orderFail {
		fmt.Println("порядок нарушен")
	}
	if EvenFail {
		fmt.Println("Есть не четные")
	}
	duration := time.Since(start)
	if duration.Seconds() > N {
		fmt.Println("Время превышено")
	}
	fmt.Println("Время выполнения: ", duration)
}
