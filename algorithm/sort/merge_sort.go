package sort

import (
	"math"
)

// 归并排序
// O(nlogn)
func MergeSort(data []int) {
	if data == nil {
		return
	}
	mergeSort(data, 0, len(data)-1)
}

func mergeSort(data []int, left, right int) {
	if left < right {
		mid := (left + right) / 2
		mergeSort(data, left, mid)
		mergeSort(data, mid+1, right)
		merge(data, left, mid, right)
	}
}

func merge(data []int, left, mid, right int) {
	if !(left <= mid && mid <= right) {
		return
	}
	l := make([]int, 0, mid-left+2)
	r := make([]int, 0, right-mid+1)

	l = append(l, data[left:mid+1]...)
	l = append(l, math.MaxInt32)

	r = append(r, data[mid+1:right+1]...)
	r = append(r, math.MaxInt32)

	x, y := 0, 0
	for i := left; i < right+1; i++ {
		if l[x] > r[y] {
			data[i] = r[y]
			y++
		} else {
			data[i] = l[x]
			x++
		}
	}
}
