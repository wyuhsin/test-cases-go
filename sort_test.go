package tests

import (
	"math/rand"
	"reflect"
	"slices"
	"strconv"
	"testing"
	"time"
)

func generateRandomSlice(size int) []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	slice := make([]int, size)

	for i := range size {
		slice[i] = rng.Intn(1000)
	}

	return slice
}

func BubbleSort(arr []int) []int {
	n := len(arr)
	if n <= 1 {
		return arr
	}

	for i := range n - 1 {
		swapped := false
		for j := range n - i - 1 {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}

	return arr
}

func SelectionSort(arr []int) []int {
	n := len(arr)
	if n <= 1 {
		return arr
	}

	for i := range n - 1 {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		arr[i], arr[minIdx] = arr[minIdx], arr[i]
	}

	return arr
}

func InsertionSort(arr []int) []int {
	n := len(arr)
	if n <= 1 {
		return arr
	}

	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1

		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}

	return arr
}

func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	l, r := 0, 0

	for l < len(left) && r < len(right) {
		if left[l] <= right[r] {
			result = append(result, left[l])
			l++
		} else {
			result = append(result, right[r])
			r++
		}
	}

	result = append(result, left[l:]...)
	result = append(result, right[r:]...)

	return result
}

func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	return quickSortRecursive(arr, 0, len(arr)-1)
}

func quickSortRecursive(arr []int, low, high int) []int {
	if low < high {
		pivotIndex := partition(arr, low, high)
		quickSortRecursive(arr, low, pivotIndex-1)
		quickSortRecursive(arr, pivotIndex+1, high)
	}
	return arr
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func HeapSort(arr []int) []int {
	n := len(arr)
	if n <= 1 {
		return arr
	}

	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}

	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		heapify(arr, i, 0)
	}

	return arr
}

func heapify(arr []int, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr[left] > arr[largest] {
		largest = left
	}

	if right < n && arr[right] > arr[largest] {
		largest = right
	}

	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

func TestBubbleSort(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Unsorted array", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{"Already sorted array", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted array", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Array with negative numbers", []int{-5, 2, -8, 0, 3}, []int{-8, -5, 0, 2, 3}},
		{"Array with duplicate elements", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"Single element array", []int{42}, []int{42}},
		{"Empty array", []int{}, []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := BubbleSort(slices.Clone(tc.input))
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("BubbleSort failed: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestSelectionSort(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Unsorted array", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{"Already sorted array", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted array", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Array with negative numbers", []int{-5, 2, -8, 0, 3}, []int{-8, -5, 0, 2, 3}},
		{"Array with duplicate elements", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"Single element array", []int{42}, []int{42}},
		{"Empty array", []int{}, []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SelectionSort(slices.Clone(tc.input))
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("SelectionSort failed: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestInsertionSort(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Unsorted array", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{"Already sorted array", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted array", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Array with negative numbers", []int{-5, 2, -8, 0, 3}, []int{-8, -5, 0, 2, 3}},
		{"Array with duplicate elements", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"Single element array", []int{42}, []int{42}},
		{"Empty array", []int{}, []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := InsertionSort(slices.Clone(tc.input))
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("InsertionSort failed: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestMergeSort(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Unsorted array", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{"Already sorted array", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted array", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Array with negative numbers", []int{-5, 2, -8, 0, 3}, []int{-8, -5, 0, 2, 3}},
		{"Array with duplicate elements", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"Single element array", []int{42}, []int{42}},
		{"Empty array", []int{}, []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MergeSort(slices.Clone(tc.input))
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("MergeSort failed: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestQuickSort(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Unsorted array", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{"Already sorted array", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted array", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Array with negative numbers", []int{-5, 2, -8, 0, 3}, []int{-8, -5, 0, 2, 3}},
		{"Array with duplicate elements", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"Single element array", []int{42}, []int{42}},
		{"Empty array", []int{}, []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := QuickSort(slices.Clone(tc.input))
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("QuickSort failed: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestHeapSort(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"Unsorted array", []int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
		{"Already sorted array", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"Reverse sorted array", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"Array with negative numbers", []int{-5, 2, -8, 0, 3}, []int{-8, -5, 0, 2, 3}},
		{"Array with duplicate elements", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"Single element array", []int{42}, []int{42}},
		{"Empty array", []int{}, []int{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := HeapSort(slices.Clone(tc.input))
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("HeapSort failed: got %v, want %v", result, tc.expected)
			}
		})
	}
}

func BenchmarkSortingAlgorithms(b *testing.B) {
	sizes := []int{10}
	algorithms := []struct {
		name string
		sort func([]int) []int
	}{
		{"BubbleSort", BubbleSort},
		{"SelectionSort", SelectionSort},
		{"InsertionSort", InsertionSort},
		{"MergeSort", MergeSort},
		{"QuickSort", QuickSort},
		{"HeapSort", HeapSort},
	}

	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			for _, algo := range algorithms {
				b.Run(algo.name, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						b.StopTimer()
						input := generateRandomSlice(size)
						b.StartTimer()
						algo.sort(input)
					}
				})
			}
		})
	}
}
