// 1. Задача  Написать функцию, которая принимает канал и число N типа int, после получения N значений из канала, функция должна вернуть срез с этими значениями

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func RandNumbers(length, max int) []int {
	var s []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		s = append(s, rand.Intn(max))
	}

	return s
}

func writeToChan(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)
	for _, v := range RandNumbers(100, 100) {
		ch <- v
	}
}

// функция, для объеденения каналов
func mergeChan(ch ...chan int) chan int {
	newCh := make(chan int, len(ch))

	var wg sync.WaitGroup
	for _, cc := range ch {
		wg.Add(1)
		go func(c chan int) {
			for val := range c {
				newCh <- val
			}
		}(cc)
	}
	go func() {
		wg.Wait()
		close(newCh)
	}()

	return newCh
}

// функция, которая принимает канал и число N типа int, после получения N значений из канала, функция должна вернуть срез с этими значениями
func readFromChan(ch <-chan int, n int) []int {
	var wg sync.WaitGroup
	results := make([]int, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			results[j] = <-ch
		}(i)
	}

	wg.Wait()
	return results
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)
	ch4 := make(chan int)
	var wg sync.WaitGroup

	wg.Add(4)

	go writeToChan(ch1, &wg)
	go writeToChan(ch2, &wg)
	go writeToChan(ch3, &wg)
	go writeToChan(ch4, &wg)

	mergedChan := mergeChan(ch1, ch2, ch3, ch4)

	go func() {
		for val := range mergedChan {
			fmt.Println(val)
		}
	}()

	wg.Wait()
}
