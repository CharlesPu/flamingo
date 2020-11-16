package sort

import "sort"

// O(nlogn)
func QuickSort(data []int) {
	quickSort(data, 0, len(data)-1)
}

func quickSort(data []int, left, right int) {
	if left < right {
		mid := partition(data, left, right)
		quickSort(data, left, mid-1)
		quickSort(data, mid+1, right)
	}
}

func partition(data []int, left, right int) int {
	cursor := left - 1
	// 以data[right]的大小作为分割依据
	for i := left; i <= right; i++ { // 保证cursor走过的都是比data[right]小的
		if data[i] < data[right] {
			cursor++
			data[cursor], data[i] = data[i], data[cursor]
		}
	}
	// ...( < [right])... [cursor] ...... ...... [right]
	cursor++
	data[cursor], data[right] = data[right], data[cursor]
	// ...( < [right])...  ......  [right]  ...  [cursor]
	// 								cursor
	return cursor
}

func QuickSortByItf(data sort.Interface) {
	quickSortByItf(data, 0, data.Len()-1)
}

func quickSortByItf(data sort.Interface, left, right int) {
	if left < right {
		mid := partitionByItf(data, left, right)
		quickSortByItf(data, left, mid-1)
		quickSortByItf(data, mid+1, right)
	}
}

func partitionByItf(data sort.Interface, left, right int) int {
	cursor := left - 1
	// 以data[right]的大小作为分割依据
	for i := left; i <= right; i++ { // 保证cursor走过的都是比data[right]小的
		if data.Less(i, right) {
			cursor++
			data.Swap(cursor, i)
		}
	}
	cursor++
	data.Swap(cursor, right)

	return cursor
}
