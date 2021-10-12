package spinlock

import (
	"spinlock/internal"
	"sync"
	"sync/atomic"
	"testing"
)

func fibonacci(num int) int {
	if num < 2 {
		return 1
	}
	return fibonacci(num-1) + fibonacci(num-2)
}

// 控制并发量
const concurency = 1000

// 模拟其他任务
const othergrt = 0

func BenchmarkLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < othergrt; j++ {
			go func() {
				fibonacci(20)
			}()
		}
		var mu sync.Mutex
		var wg sync.WaitGroup
		wg.Add(concurency)
		cnt := 0
		for j := 0; j < concurency; j++ {
			go func() {
				for k := 0; k < 100; k++ {
					mu.Lock()
					cnt++
					mu.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkSpinLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < othergrt; j++ {
			go func() {
				fibonacci(20)
			}()
		}
		sp := internal.NewSpinLock()
		var wg sync.WaitGroup
		wg.Add(concurency)
		cnt := 0
		for j := 0; j < concurency; j++ {
			go func() {
				for k := 0; k < 100; k++ {
					sp.Lock()
					cnt++
					sp.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkAtomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < othergrt; j++ {
			go func() {
				fibonacci(20)
			}()
		}
		sp := internal.NewSpinLock()
		var wg sync.WaitGroup
		wg.Add(concurency)
		var cnt int32 = 0
		for j := 0; j < concurency; j++ {
			go func() {
				for k := 0; k < 100; k++ {
					sp.Lock()
					atomic.AddInt32(&cnt, 1)
					sp.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkLongTaskMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < othergrt; j++ {
			go func() {
				fibonacci(20)
			}()
		}
		var mu sync.Mutex
		var wg sync.WaitGroup
		wg.Add(concurency)
		for j := 0; j < concurency; j++ {
			go func() {
				for k := 0; k < 100; k++ {
					mu.Lock()
					fibonacci(20)
					mu.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkLongTaskSpinLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < othergrt; j++ {
			go func() {
				fibonacci(20)
			}()
		}
		sp := internal.NewSpinLock()
		var wg sync.WaitGroup
		wg.Add(concurency)
		for j := 0; j < concurency; j++ {
			go func() {
				for k := 0; k < 100; k++ {
					sp.Lock()
					fibonacci(20)
					sp.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

//
//func TestSpinLock(t *testing.T) {
//	sp:=internal.NewSpinLock()
//	cnt:=0
//	var wg sync.WaitGroup
//	wg.Add(3)
//	for j:=0;j<3;j++{
//		go func() {
//			for k:=0;k<100;k++{
//				sp.Lock()
//				cnt++
//				fmt.Println(cnt)
//				sp.Unlock()
//			}
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//}
