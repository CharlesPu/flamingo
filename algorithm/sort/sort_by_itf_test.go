package sort

import (
	"sort"
	"testing"
)

var b = sort.IntSlice{
	2, 1, 5, 4, 10, 6, 12, 3, 4,
}

func TestInsertSortByItf(t *testing.T) {
	InsertSortByItf(b)
	t.Log(b)
}

func BenchmarkInsertSortByItf(t *testing.B) {
	for i := 0; i < t.N; i++ {
		InsertSortByItf(b)
	}
}

func TestHeapSortByItf(t *testing.T) {
	HeapSortByItf(b)
	t.Log(b)
}

func BenchmarkHeapSortByItf(t *testing.B) {
	for i := 0; i < t.N; i++ {
		HeapSortByItf(b)
	}
}

func TestQuickSortByItf(t *testing.T) {
	QuickSortByItf(b)
	t.Log(b)
}

func BenchmarkQuickSortByItf(t *testing.B) {
	for i := 0; i < t.N; i++ {
		QuickSortByItf(b)
	}
}

func TestBubbleSortByItf(t *testing.T) {
	BubbleSortByItf(b)
	t.Log(b)
}

func BenchmarkBubbleSortByItf(t *testing.B) {
	for i := 0; i < t.N; i++ {
		BubbleSortByItf(b)
	}
}

func TestShellSortByItf(t *testing.T) {
	ShellSortByItf(b)
	t.Log(b)
}

func BenchmarkShellSortByItf(t *testing.B) {
	for i := 0; i < t.N; i++ {
		ShellSortByItf(b)
	}
}
