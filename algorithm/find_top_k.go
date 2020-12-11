// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.

package find

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"github.com/CharlesPu/flamingo/container/queue"
)

// FindTopK
func FindTop50(path string) []uint32 {
	return findTopK(50, path)
}

const (
	fileMaxN = 4096
)

// inspired by Map-Reduce from kubernetes scheduler
// Map:
// step1. currently process every path(num N)
// step2. for each path, read each line, use Insert Sort to get top K nums
// step3. get N files' top K nums result

// Reduce:
// step1. traverse each path result, build Max Heap(size N)
// step2. get top value from Heap, which is the top i num of all numbers, and append it to final result
// step3. modify the top value to next maximum number, besides, maxHeapify heap(sink top value)
// step4. continuously execute step2-step3 until the length of final result is K
func findTopK(k int, path string) []uint32 {
	resN := make([][]uint32, 0, fileMaxN)
	mu := sync.Mutex{}
	wg := &sync.WaitGroup{}

	// map
	var traverseDir func(fn string)
	traverseDir = func(dir string) {
		rd, err := ioutil.ReadDir(dir)
		if err != nil {
			return
		}
		for i := range rd {
			if rd[i].IsDir() {
				traverseDir(dir + "/" + rd[i].Name())
				continue
			}

			wg.Add(1)
			go func(fn string) {
				tmp := mapper(k, fn)
				mu.Lock()
				resN = append(resN, tmp)
				mu.Unlock()
				wg.Done()
			}(dir + "/" + rd[i].Name())
		}
	}
	traverseDir(path)

	wg.Wait()

	// reduce
	res := reducer(k, resN)

	return res
}

func mapper(k int, fileName string) []uint32 {
	f, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	defer f.Close()

	// fmt.Println("start to process path", fileName)
	r := bufio.NewReader(f)
	res := make([]uint32, 0, k+1) // last one is cache
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		// process line
		n, convErr := strconv.Atoi(string(line))
		if convErr != nil {
			continue
		}
		// append to end
		num := uint32(n)
		if len(res) >= k+1 {
			res = res[:k+1]
			res[k] = num
		} else {
			res = append(res, num)
		}
		// find correct position by insert sort
		i := len(res) - 1
		j := i - 1
		for j >= 0 && res[j] < res[j+1] {
			res[j], res[j+1] = res[j+1], res[j]
			j--
		}
	}

	l := len(res)
	if l > k {
		l = k
	}
	// fmt.Println("end process path", fileName)

	return res[:l]
}

type (
	entry struct {
		val     uint32
		idx     int
		fileIdx int
	}
)

func reducer(k int, resN [][]uint32) []uint32 {
	res := make([]uint32, 0, k)
	// build max heap
	maxHeap := queue.NewMaxPriorityQueue(len(resN), func(a interface{}, b interface{}) bool {
		return a.(*entry).val < b.(*entry).val
	})

	for i := 0; i < len(resN); i++ {
		maxHeap.Insert(&entry{val: resN[i][0], idx: 0, fileIdx: i})
	}

	// maxHeapify until len(res) is k
	for i := 0; i < k; i++ {
		e, ok := maxHeap.Peek().(*entry)
		if !ok || e == nil {
			break
		}
		res = append(res, e.val)
		nextIdx := e.idx + 1
		if nextIdx >= len(resN[e.fileIdx]) { // when one file's top K result is fully processed, remove it
			maxHeap.Pop()
		} else {
			maxHeap.SwapTop(&entry{val: resN[e.fileIdx][e.idx+1], idx: e.idx + 1, fileIdx: e.fileIdx})
		}
	}

	return res
}
