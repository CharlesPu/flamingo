package sort

import "sort"

// O(nlogn)
func HeapSort(data []int) {
	l := len(data)
	buildMaxHeap(data, l)
	heapSize := l
	for i := l - 1; i >= 0; i-- {
		data[i], data[0] = data[0], data[i]
		heapSize--
		maxHeapIfy(data, 0, heapSize)
	}
}

func buildMaxHeap(data []int, l int) {
	for i := l / 2; i >= 0; i-- {
		maxHeapIfy(data, i, l)
	}
}

// 从顶往下调整大顶堆
func maxHeapIfy(data []int, idx, heapSize int) {
	l := leftChild(idx)
	r := rightChild(idx)

	maxIdx := idx
	if l < heapSize && data[l] > data[maxIdx] {
		maxIdx = l
	}
	if r < heapSize && data[r] > data[maxIdx] {
		maxIdx = r
	}
	if maxIdx != idx {
		data[maxIdx], data[idx] = data[idx], data[maxIdx]
		maxHeapIfy(data, maxIdx, heapSize)
	}
}

func leftChild(i int) int {
	return 2*(i+1) - 1
}

func rightChild(i int) int {
	return 2 * (i + 1)
}

func HeapSortByItf(data sort.Interface) {
	l := data.Len()
	buildMaxHeapByItf(data, l)
	heapSize := l
	for i := l - 1; i >= 0; i-- {
		data.Swap(i, 0)
		heapSize--
		maxHeapIfyByItf(data, 0, heapSize)
	}
}

func buildMaxHeapByItf(data sort.Interface, l int) {
	for i := l / 2; i >= 0; i-- {
		maxHeapIfyByItf(data, i, l)
	}
}

// 从顶往下调整大顶堆
func maxHeapIfyByItf(data sort.Interface, idx, heapSize int) {
	l := leftChild(idx)
	r := rightChild(idx)

	maxIdx := idx
	if l < heapSize && data.Less(maxIdx, l) {
		maxIdx = l
	}
	if r < heapSize && data.Less(maxIdx, r) {
		maxIdx = r
	}
	if maxIdx != idx {
		data.Swap(maxIdx, idx)
		maxHeapIfyByItf(data, maxIdx, heapSize)
	}
}
