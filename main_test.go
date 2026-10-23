package main

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

func TestGenerateRandomElements(t *testing.T) {
	for _, size := range []int{-1, 0, 1, 7, 8, 9, 1000} {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			data := generateRandomElements(size)
			wantLength := max(size, 0)
			if len(data) != wantLength {
				t.Fatalf("length = %d, want %d", len(data), wantLength)
			}
			for i, value := range data {
				if value <= 0 {
					t.Fatalf("data[%d] = %d, want a positive integer", i, value)
				}
			}
		})
	}
}

func TestMaximum(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	tests := []struct {
		name string
		data []int
		want int
	}{
		{"nil", nil, 0},
		{"empty", []int{}, 0},
		{"single", []int{42}, 42},
		{"zero", []int{0}, 0},
		{"first", []int{9, 2, 1}, 9},
		{"middle", []int{1, 9, 2}, 9},
		{"last", []int{1, 2, 9}, 9},
		{"duplicates", []int{7, 7, 7}, 7},
		{"negative", []int{-10, -2, -7}, -2},
		{"mixed", []int{-3, 0, 4, -2}, 4},
		{"eight_chunks", []int{1, 2, 3, 4, 5, 6, 7, 8}, 8},
		{"remainder", []int{1, 2, 3, 4, 5, 6, 7, 8, 99}, 99},
		{"minimum_int", []int{minInt, minInt}, minInt},
		{"maximum_int", []int{minInt, maxInt, 0}, maxInt},
	}
	for _, search := range []struct {
		name string
		fn   func([]int) int
	}{{"sequential", maximum}, {"parallel", maxChunks}} {
		for _, tt := range tests {
			t.Run(search.name+"/"+tt.name, func(t *testing.T) {
				before := slices.Clone(tt.data)
				if got := search.fn(tt.data); got != tt.want {
					t.Errorf("maximum = %d, want %d", got, tt.want)
				}
				if !slices.Equal(tt.data, before) {
					t.Error("search modified its input")
				}
			})
		}
	}
}

func TestMaxChunksPartitionBoundaries(t *testing.T) {
	// Every possible position must be inspected for every size near chunk boundaries.
	for size := 1; size <= CHUNKS*3+1; size++ {
		for position := 0; position < size; position++ {
			data := make([]int, size)
			for i := range data {
				data[i] = -10
			}
			data[position] = -1
			if got := maxChunks(data); got != -1 {
				t.Fatalf("size=%d position=%d: got %d, want -1", size, position, got)
			}
		}
	}
}

func TestMaxChunksMatchesReference(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for _, size := range []int{31, 32, 33, 1000, 10001} {
		data := make([]int, size)
		for i := range data {
			data[i] = rng.Intn(2000001) - 1000000
		}
		want := slices.Max(data)
		if got := maxChunks(data); got != want {
			t.Fatalf("size=%d: got %d, want %d", size, got, want)
		}
	}
}
