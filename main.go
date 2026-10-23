package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

// generateRandomElements returns size positive random integers.
// Non-positive sizes produce an empty slice.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return []int{}
	}
	data := make([]int, size)
	for i := range data {
		data[i] = rand.Intn(1_000_000_000) + 1
	}
	return data
}

// maximum returns the largest element, or zero for an empty slice.
func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}
	maxValue := data[0]
	for _, value := range data[1:] {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

// maxChunks finds the maximum using up to CHUNKS independent goroutines.
// Small slices use one worker per element; empty slices return zero.
func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}
	chunkCount := min(CHUNKS, len(data))
	maxValues := make([]int, chunkCount)
	var wg sync.WaitGroup
	wg.Add(chunkCount)
	for i := 0; i < chunkCount; i++ {
		// These boundaries cover all elements, including the remainder.
		start := i * len(data) / chunkCount
		end := (i + 1) * len(data) / chunkCount
		go func(index int, chunk []int) {
			defer wg.Done()
			// Each worker owns a separate element, so no mutex is needed.
			maxValues[index] = maximum(chunk)
		}(i, data[start:end])
	}
	wg.Wait()
	return maximum(maxValues)
}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	data := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	start := time.Now()
	maxValue := maximum(data)
	elapsed := time.Since(start).Microseconds()

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed)

	fmt.Printf("Ищем максимальное значение в %d горутинах\n", CHUNKS)
	start = time.Now()
	maxValue = maxChunks(data)
	elapsed = time.Since(start).Microseconds()

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed)
}
