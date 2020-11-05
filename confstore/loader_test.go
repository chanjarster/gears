package confstore

import (
	"testing"
	"time"
)

type noopStore struct {
}

func (n *noopStore) RegisterKey(key string, defaultValue string, v ValueProcessor) error {
	return nil
}

func (n *noopStore) DeregisterKey(key string) {
}

func (n *noopStore) BatchUpdate(kvStrs []*KVStr) []*KVError {
	return nil
}

func (n *noopStore) Update(key string, value string) *KVError {
	return nil
}

func (n *noopStore) UpdateNoPersist(key string, value string) *KVError {
	return nil
}

func (n *noopStore) BatchGetValues(keys []string) []*KV {
	return nil
}

func (n *noopStore) BatchGetValueString(keys []string) []*KVStr {
	return nil
}

func (n *noopStore) MustGetValue(key string) interface{} {
	return nil
}

func (n *noopStore) MustGetValueString(key string) string {
	return ""
}

func (n *noopStore) GetValue(key string) (v interface{}, hit bool) {
	return nil, false
}

func (n *noopStore) GetValueString(key string) (v string, hit bool) {
	return "", false
}

func (n *noopStore) ResetKey(key string) {
}

type mockPersister struct {
	loadCount int
}

func (m *mockPersister) Load(s Interface) error {
	m.loadCount++
	return nil
}

func (m *mockPersister) Save(key, value string) error {
	panic("implement me")
}

func (m *mockPersister) BatchSave(strs []*KVStr) error {
	panic("implement me")
}

func (m *mockPersister) Delete(key string) error {
	panic("implement me")
}

func Test_defaultLoadPolicy_DoLoad(t1 *testing.T) {

	s := &noopStore{}
	p := &mockPersister{}
	t := NewLoadPolicy(time.Millisecond * 10)

	type args struct {
		sleep time.Duration
		s     Interface
		p     *mockPersister
	}
	tests := []struct {
		name          string
		args          args
		wantLoadCount int
	}{
		{
			args:          args{0, s, p},
			wantLoadCount: 1,
		},
		{
			args:          args{0, s, p},
			wantLoadCount: 1,
		},
		{
			args:          args{0, s, p},
			wantLoadCount: 1,
		},
		{
			args:          args{time.Millisecond * 15, s, p},
			wantLoadCount: 2,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if tt.args.sleep > 0 {
				time.Sleep(tt.args.sleep)
			}
			if err := t.DoLoad(tt.args.s, tt.args.p); err != nil {
				t1.Errorf("DoLoad() error = %v", err)
			}
			if got := tt.args.p.loadCount; got != tt.wantLoadCount {
				t1.Errorf("After DoLoad() loadCount = %v, want = %v", got, tt.wantLoadCount)
			}
		})
	}
}
