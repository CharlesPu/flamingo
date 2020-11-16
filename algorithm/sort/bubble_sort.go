package sort

import "sort"

func BubbleSort(data []int) {
	low := 0
	high := len(data) - 1

	for low < high {
		for i := low; i < high; i++ {
			if data[i] > data[i+1] {
				data[i], data[i+1] = data[i+1], data[i]
			}
		}
		high--

		for i := high; i > low; i-- {
			if data[i] < data[i-1] {
				data[i], data[i-1] = data[i-1], data[i]
			}
		}
		low++
	}
}

func BubbleSortByItf(data sort.Interface) {
	low := 0
	high := data.Len() - 1

	for low < high {
		for i := low; i < high; i++ {
			if data.Less(i+1, i) {
				data.Swap(i, i+1)
			}
		}
		high--

		for i := high; i > low; i-- {
			if data.Less(i, i-1) {
				data.Swap(i, i-1)
			}
		}
		low++
	}
}
