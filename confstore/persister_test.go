package confstore

import (
	"github.com/chanjarster/gears/confs"
	"os"
	"reflect"
	"testing"
)

type mockStore struct {
	data map[string]string
}

func (m *mockStore) LoadFromPersister() error {
	panic("implement me")
}

func (m *mockStore) RegisterKey(key string, defaultValue string, v ValueProcessor) error {
	panic("implement me")
}

func (m *mockStore) DeregisterKey(key string) {
	panic("implement me")
}

func (m *mockStore) BatchUpdate(kvStrs []*KVStr) []*KVError {
	panic("implement me")
}

func (m *mockStore) Update(key string, value string) *KVError {
	panic("implement me")
}

func (m *mockStore) UpdateNoPersist(key string, value string) *KVError {
	m.data[key] = value
	return nil
}

func (m *mockStore) BatchGetValues(keys []string) []*KV {
	panic("implement me")
}

func (m *mockStore) BatchGetValueString(keys []string) []*KVStr {
	panic("implement me")
}

func (m *mockStore) MustGetValue(key string) interface{} {
	panic("implement me")
}

func (m *mockStore) MustGetValueString(key string) string {
	panic("implement me")
}

func (m *mockStore) GetValue(key string) (v interface{}, hit bool) {
	panic("implement me")
}

func (m *mockStore) GetValueString(key string) (v string, hit bool) {
	panic("implement me")
}

func (m *mockStore) ResetKey(key string) {
	panic("implement me")
}

func Test_redisPersister_Load(t *testing.T) {

	val, hit := os.LookupEnv("INTEGRATION_TEST")
	if !hit || val != "true" {
		t.Skip("skip integration test")
	}

	redisClient := confs.NewRedisClient(&confs.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Pool:     10,
		MinIdle:  1,
	}, nil)
	defer redisClient.Close()

	store := &mockStore{
		data: make(map[string]string),
	}

	r := &redisPersister{
		redisClient:   redisClient,
		configRootKey: "_foo_",
	}

	type args struct {
		key, value string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			args: args{"foo", "bar"},
			want: map[string]string{"foo": "bar"},
		},
		{
			args: args{"foo", "zoo"},
			want: map[string]string{"foo": "zoo"},
		},
		{
			args: args{"bar", "foo"},
			want: map[string]string{"foo": "zoo", "bar": "foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisClient.FlushAll()
			r.Save(tt.args.key, tt.args.value)
			r.Load(store)
			if got := store.data; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("After Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisPersister_BatchSave(t *testing.T) {

	val, hit := os.LookupEnv("INTEGRATION_TEST")
	if !hit || val != "true" {
		t.Skip("skip integration test")
	}

	redisClient := confs.NewRedisClient(&confs.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Pool:     10,
		MinIdle:  1,
	}, nil)
	defer redisClient.Close()

	store := &mockStore{
		data: make(map[string]string),
	}

	r := &redisPersister{
		redisClient:   redisClient,
		configRootKey: "_foo_",
	}

	type args struct {
		kvStrs []*KVStr
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			args: args{[]*KVStr{{"foo", "bar"}, {"zoo", "foo"}}},
			want: map[string]string{"foo": "bar", "zoo": "foo"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisClient.FlushAll()
			r.BatchSave(tt.args.kvStrs)
			r.Load(store)
			if got := store.data; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("After Load() got = %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_redisPersister_Delete(t *testing.T) {
	val, hit := os.LookupEnv("INTEGRATION_TEST")
	if !hit || val != "true" {
		t.Skip("skip integration test")
	}

	redisClient := confs.NewRedisClient(&confs.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Pool:     10,
		MinIdle:  1,
	}, nil)
	defer redisClient.Close()

	store := &mockStore{
		data: make(map[string]string),
	}

	r := &redisPersister{
		redisClient:   redisClient,
		configRootKey: "_foo_",
	}

	r.BatchSave([]*KVStr{{"foo", "bar"}, {"zoo", "foo"},{"bar","foo"}})
	r.Delete("foo")
	r.Load(store)

	want := map[string]string{"bar": "foo", "zoo": "foo"}
	if got := store.data; !reflect.DeepEqual(got, want) {
		t.Errorf("After Load() got = %v, want %v", got, want)
	}

}
