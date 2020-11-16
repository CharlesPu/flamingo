package sort

import (
	"testing"
)

var a = []int{
	2, 1, 5, 4, 10, 6, 12, 3,
}

func TestInsertSort(t *testing.T) {
	InsertSort(a)
	t.Log(a)
}

func BenchmarkInsertSort(t *testing.B) {
	for i := 0; i < t.N; i++ {
		InsertSort(a)
	}
}

func TestMergeSort(t *testing.T) {
	MergeSort(a)
	t.Log(a)
}

func BenchmarkMergeSort(t *testing.B) {
	for i := 0; i < t.N; i++ {
		MergeSort(a)
	}
}

func TestHeapSort(t *testing.T) {
	HeapSort(a)
	t.Log(a)
}

func BenchmarkHeapSort(t *testing.B) {
	for i := 0; i < t.N; i++ {
		HeapSort(a)
	}
}

func TestQuickSort(t *testing.T) {
	QuickSort(a)
	t.Log(a)
}

func BenchmarkQuickSort(t *testing.B) {
	for i := 0; i < t.N; i++ {
		QuickSort(a)
	}
}

func TestBubbleSort(t *testing.T) {
	BubbleSort(a)
	t.Log(a)
}

func BenchmarkBubbleSort(t *testing.B) {
	for i := 0; i < t.N; i++ {
		BubbleSort(a)
	}
}

func TestShellSort(t *testing.T) {
	ShellSort(a)
	t.Log(a)
}

func BenchmarkShellSort(t *testing.B) {
	for i := 0; i < t.N; i++ {
		ShellSort(a)
	}
}
