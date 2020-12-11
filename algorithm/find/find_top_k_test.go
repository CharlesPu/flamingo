package find

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

const (
	tmpDirAlpha = "alpha"
	tmpDirBeta  = "beta"
)

// fileNum=100000   lineNum=100000
// GOMAXPROCS=12      b.N=1      136986870300ns(2.28min)/op
func BenchmarkFindTop50_GOMAXPROCS_12(b *testing.B) {
	genTestFiles(tmpDirBeta, 100000, 100000)
	defer cleanTestFiles(tmpDirBeta)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindTop50(tmpDirBeta)
	}
}

// fileNum=100000   lineNum=100000
// GOMAXPROCS=1        b.N=1      631946887400ns(10.53min)/op
func BenchmarkFindTop50_GOMAXPROCS_1(b *testing.B) {
	genTestFiles(tmpDirBeta, 100000, 100000)
	defer cleanTestFiles(tmpDirBeta)

	runtime.GOMAXPROCS(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindTop50(tmpDirBeta)
	}
}

func TestFindTop50(t *testing.T) {
	genTestFiles(tmpDirBeta, 100000, 100000)
	defer cleanTestFiles(tmpDirBeta)

	res := FindTop50(tmpDirBeta)
	t.Log(res)
}

func TestGenTestFiles(t *testing.T) {
	genTestFiles(tmpDirAlpha, 100, 100, []uint32{1231231, 823123, 213418})
}

// choose random file and line to set expectResults,
// hence a set of files with expected results
// Note: expectResults[0] must be sorted
func genTestFiles(dir string, fileNum, lineNum int, expectResults ...[]uint32) {
	isNotRandom := false
	var min uint32
	var hashRes map[string]uint32

	if len(expectResults) != 0 && len(expectResults[0]) != 0 {
		isNotRandom = true
		hashRes = make(map[string]uint32, len(expectResults[0]))
		min = expectResults[0][len(expectResults[0])-1]

		setIfNew := func(idx string, value uint32) bool {
			_, ok := hashRes[idx]
			if ok {
				return false
			}
			hashRes[idx] = value
			return true
		}

		for _, v := range expectResults[0] {
			for !setIfNew(fmt.Sprintf("%d-%d", rand.Intn(fileNum), rand.Intn(lineNum)), v) {
			}
		}
	}

	os.Mkdir(dir, os.ModePerm)
	rand.Seed(time.Now().Unix())
	for i := 0; i < fileNum; i++ {
		f, err := os.Create(dir + fmt.Sprintf("/%d.txt", i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		for j := 0; j < lineNum; j++ {
			if !isNotRandom {
				f.WriteString(strconv.Itoa(rand.Intn(int(min))) + "\n")
				continue
			}
			v, exist := hashRes[fmt.Sprintf("%d-%d", i, j)]
			if !exist {
				f.WriteString(strconv.Itoa(rand.Intn(int(min))) + "\n")
				continue
			}
			f.WriteString(strconv.Itoa(int(v)) + "\n")

		}
		f.Close()
	}
}

func cleanTestFiles(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		fmt.Println(err)
	}
}

func TestReducer(t *testing.T) {
	type args struct {
		k    int
		resN [][]uint32
	}
	tests := []struct {
		name string
		args args
		want []uint32
	}{
		{
			args: args{
				k: 1,
				resN: [][]uint32{
					{3, 3, 1},
					{8, 7, 6, 4},
				},
			},
			want: []uint32{8},
		},
		{
			args: args{
				k: 2,
				resN: [][]uint32{
					{3, 3, 1},
					{8, 7, 6, 4},
				},
			},
			want: []uint32{8, 7},
		},
		{
			args: args{
				k: 5,
				resN: [][]uint32{
					{5, 4, 1},
					{8, 7, 6, 3},
				},
			},
			want: []uint32{8, 7, 6, 5, 4},
		},
		{
			args: args{
				k: 100,
				resN: [][]uint32{
					{5, 3, 1},
					{8, 7, 6, 4, 0},
					{9, 8, 8, 2},
				},
			},
			want: []uint32{9, 8, 8, 8, 7, 6, 5, 4, 3, 2, 1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reducer(tt.args.k, tt.args.resN); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapper(t *testing.T) {
	type args struct {
		k        int
		fileName string
	}
	tests := []struct {
		name string
		args args
		want []uint32
	}{
		{
			want: []uint32{10000, 9999, 8888},
		},
		{
			want: []uint32{1234567890, 234567890, 34567890, 4567890, 567890},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// use expect results to generate files
			genTestFiles(tmpDirAlpha, 1, 10000, tt.want)
			defer cleanTestFiles(tmpDirAlpha)

			if got := mapper(len(tt.want), tmpDirAlpha+"/0.txt"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindTopK(t *testing.T) {
	type args struct {
		k       int
		fileNum int
		lineNum int
		path    string
	}
	tests := []struct {
		name string
		args args
		want []uint32
	}{
		{
			args: args{
				fileNum: 100,
				lineNum: 100,
			},
			want: []uint32{10000, 9999, 8888},
		},
		{
			args: args{
				fileNum: 100,
				lineNum: 100,
			},
			want: []uint32{1234567890, 234567890, 34567890, 4567890, 567890},
		},
		{
			args: args{
				fileNum: 2,
				lineNum: 3,
			},
			want: []uint32{1234567890, 234567890, 34567890, 4567890, 567890, 67890},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// use expect results to generate files
			genTestFiles(tmpDirAlpha, tt.args.fileNum, tt.args.lineNum, tt.want)
			defer cleanTestFiles(tmpDirAlpha)

			if got := findTopK(len(tt.want), tmpDirAlpha); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findTopK() = %v, want %v", got, tt.want)
			}
		})
	}
}
