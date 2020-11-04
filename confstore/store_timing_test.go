package configstore

import (
	"math/rand"
	"testing"
)

var (
	_dict = []string{
		"foo", "bar", "zoo", "alice",
		"daniel", "fox", "dog", "fish",
		"cat", "mouse", "kangaroo",
		"dwarf", "insect", "zebra", "sheep",
		"pig", "bear", "duck", "lizard", "shepard",
	}
)

func randWord() string {
	return _dict[rand.Intn(len(_dict))]
}

func Benchmark_memStore(b *testing.B) {

	m := NewStore(NoopPersister, NoopLoadPolicy)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			for i := 0; i < 10; i++ {
				m.RegisterKey(randWord(), randWord(), String)
			}

			m.DeregisterKey(randWord())
			m.Update(randWord(), randWord())
			m.BatchUpdate([]*KVStr{
				{randWord(), randWord()},
				{randWord(), randWord()},
			})
			m.GetValueString(randWord())
			m.GetValue(randWord())
			m.BatchGetValues([]string{randWord(), randWord()})
			m.BatchGetValueString([]string{randWord(), randWord()})
			m.ResetKey(randWord())
		}
	})

}
