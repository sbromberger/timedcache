package timedmap

import (
	"math/rand"
	"testing"
	"time"
)

func benchmarkTimedMapSet(size int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := New[int, int](time.Duration(30 * time.Second))
		for i := 0; i < size; i++ {
			m.Set(i, i)
		}
	}
}

func benchmarkMapSet(size int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		// m := map[int]int{}
		m := make(map[int]int, size*2)
		for i := 0; i < size; i++ {
			m[i] = i
		}
	}

}

func benchmarkTimedMapGet(size int, b *testing.B) {
	m := New[int, int](time.Duration(30 * time.Second))
	accesses := make([]int, size*2)
	for i := 0; i < size; i++ {
		m.Set(i, i)
		accesses[i] = i
		accesses[size+i] = size + i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(accesses), func(i, j int) { accesses[i], accesses[j] = accesses[j], accesses[i] })
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := range accesses {
			_, _ = m.Get(i)
		}
	}
}

func benchmarkMapGet(size int, b *testing.B) {
	m := map[int]int{}
	accesses := make([]int, size*2)
	for i := 0; i < size; i++ {
		m[i] = i
		accesses[i] = i
		accesses[size+i] = size + i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(accesses), func(i, j int) { accesses[i], accesses[j] = accesses[j], accesses[i] })
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := range accesses {
			_ = m[i]
		}
	}
}

func BenchmarkTimedMapSet1(b *testing.B) { benchmarkTimedMapSet(10, b) }

func BenchmarkTimedMapSet2(b *testing.B) { benchmarkTimedMapSet(100, b) }
func BenchmarkTimedMapSet3(b *testing.B) { benchmarkTimedMapSet(1000, b) }
func BenchmarkTimedMapSet4(b *testing.B) { benchmarkTimedMapSet(10000, b) }
func BenchmarkTimedMapSet5(b *testing.B) { benchmarkTimedMapSet(100000, b) }
func BenchmarkTimedMapSet6(b *testing.B) { benchmarkTimedMapSet(1000000, b) }

func BenchmarkMapSet1(b *testing.B) { benchmarkMapSet(10, b) }
func BenchmarkMapSet2(b *testing.B) { benchmarkMapSet(100, b) }
func BenchmarkMapSet3(b *testing.B) { benchmarkMapSet(1000, b) }
func BenchmarkMapSet4(b *testing.B) { benchmarkMapSet(10000, b) }
func BenchmarkMapSet5(b *testing.B) { benchmarkMapSet(100000, b) }
func BenchmarkMapSet6(b *testing.B) { benchmarkMapSet(1000000, b) }

func BenchmarkTimedMapGet1(b *testing.B) { benchmarkTimedMapGet(10, b) }
func BenchmarkTimedMapGet2(b *testing.B) { benchmarkTimedMapGet(100, b) }
func BenchmarkTimedMapGet3(b *testing.B) { benchmarkTimedMapGet(1000, b) }
func BenchmarkTimedMapGet4(b *testing.B) { benchmarkTimedMapGet(10000, b) }
func BenchmarkTimedMapGet5(b *testing.B) { benchmarkTimedMapGet(100000, b) }
func BenchmarkTimedMapGet6(b *testing.B) { benchmarkTimedMapGet(1000000, b) }

func BenchmarkMapGet1(b *testing.B) { benchmarkMapGet(10, b) }
func BenchmarkMapGet2(b *testing.B) { benchmarkMapGet(100, b) }

func BenchmarkMapGet3(b *testing.B) { benchmarkMapGet(1000, b) }
func BenchmarkMapGet4(b *testing.B) { benchmarkMapGet(10000, b) }
func BenchmarkMapGet5(b *testing.B) { benchmarkMapGet(100000, b) }
func BenchmarkMapGet6(b *testing.B) { benchmarkMapGet(1000000, b) }
