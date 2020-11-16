package sort

import "sort"

func ShellSort(data []int) {
	for gap := len(data) / 2; gap > 0; gap /= 2 {
		for i := 0; i < gap; i++ {
			// 以步长gap的数组对所有元素进行插入排序
			for j := i; j < len(data); j += gap {
				k := j - gap
				for k >= 0 && data[k] > data[k+gap] {
					data[k], data[k+gap] = data[k+gap], data[k]
					k -= gap
				}
			}
		}
	}
}

func ShellSortByItf(data sort.Interface) {
	for gap := data.Len() / 2; gap > 0; gap /= 2 {
		for i := 0; i < gap; i++ {
			// 以步长gap的数组对所有元素进行插入排序
			for j := i; j < data.Len(); j += gap {
				k := j - gap
				for k >= 0 && data.Less(k+gap, k) {
					data.Swap(k+gap, k)
					k -= gap
				}
			}
		}
	}
}
