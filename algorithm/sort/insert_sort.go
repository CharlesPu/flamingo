// Copyright (c) Huawei Technologies Co., Ltd. 2012-2019. All rights reserved.
package sort

import (
	"sort"
)

// 插入排序
func InsertSort(data []int) {
	for i := 0; i < len(data); i++ {
		j := i - 1
		// 保证data[i]找到它理想的位置
		// 遍历data[i]之前的所有，比data[i]大的都跑到后面去，最终data[i]找到了自己的位置
		for j >= 0 && data[j] > data[j+1] {
			data[j], data[j+1] = data[j+1], data[j]
			j--
		}
	}
}

func InsertSortByItf(data sort.Interface) {
	for i := 0; i < data.Len(); i++ {
		j := i - 1
		for j >= 0 && data.Less(j+1, j) {
			data.Swap(j, j+1)
			j--
		}
	}
}
