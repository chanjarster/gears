package event

import (
	"fmt"
	"testing"
)

func BenchmarkFanOutBus_GoDispatch1(b *testing.B) {
	benchmarkFanOutBus_GoDispatch(b, 1)
}

func BenchmarkFanOutBus_GoDispatch10(b *testing.B) {
	benchmarkFanOutBus_GoDispatch(b, 10)
}

func BenchmarkFanOutBus_GoDispatch100(b *testing.B) {
	benchmarkFanOutBus_GoDispatch(b, 100)
}

func benchmarkFanOutBus_GoDispatch(b *testing.B, recvCount int) {
	bus := NewFanOutBus(0)
	bus.GoDispatch()

	for i := 0; i < recvCount; i++ {
		recv := bus.NewRecv(fmt.Sprintf("#%d", i), 1024)
		go func() {
			for range recv.C {
				// just consuming
			}
		}()
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bus.C <- 1
		}
	})

}
