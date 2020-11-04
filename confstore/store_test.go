package confstore

import (
	"reflect"
	"testing"
)

func Test_memStore_RegisterKey(t *testing.T) {

	m := NewStore(NoopPersister, NoopLoadPolicy)

	type args struct {
		key          string
		defaultValue string
		vp           ValueProcessor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{key: "", defaultValue: "true", vp: Bool},
			wantErr: true,
		},
		{
			args:    args{key: "foo", defaultValue: "1", vp: Bool},
			wantErr: true,
		},
		{
			args:    args{key: "foo", defaultValue: "1", vp: Int},
			wantErr: false,
		},
		{
			args:    args{key: "foo", defaultValue: "true", vp: Bool},
			wantErr: true,
		},
		{
			args:    args{key: "foo.bar", defaultValue: "1", vp: Int},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.RegisterKey(tt.args.key, tt.args.defaultValue, tt.args.vp); (err != nil) != tt.wantErr {
				t.Errorf("RegisterKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memStore_BatchUpdate(t *testing.T) {
	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)

	type args struct {
		kvStrs []*KVStr
	}
	var tests = []struct {
		name string
		args args
		want []*KVError
	}{
		{
			args: args{kvStrs: []*KVStr{
				{Key: "", Value: ""},
			}},
			want: []*KVError{},
		},
		{
			args: args{kvStrs: []*KVStr{
				{Key: "bar", Value: ""},
			}},
			want: []*KVError{{
				Key:   "bar",
				Value: "",
				Error: "key[bar] not registered",
			}},
		},
		{
			args: args{kvStrs: []*KVStr{
				{Key: "foo", Value: "2"},
			}},
			want: []*KVError{},
		},
		{
			args: args{kvStrs: []*KVStr{
				{Key: "foo", Value: "a"},
			}},
			want: []*KVError{{
				Key:   "foo",
				Value: "a",
				Error: "not int value",
			}},
		},
		{
			args: args{kvStrs: []*KVStr{
				{Key: "foo", Value: "2"},
				{Key: "foo", Value: "a"},
			}},
			want: []*KVError{{
				Key:   "foo",
				Value: "a",
				Error: "not int value",
			}},
		},
		{
			args: args{kvStrs: []*KVStr{
				{Key: "foo", Value: "a"},
				{Key: "foo", Value: "2"},
			}},
			want: []*KVError{{
				Key:   "foo",
				Value: "a",
				Error: "not int value",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.BatchUpdate(tt.args.kvStrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memStore_Update(t *testing.T) {
	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want *KVError
	}{
		{
			args: args{"foo", "a"},
			want: &KVError{Key: "foo", Value: "a", Error: "not int value"},
		},
		{
			args: args{"foo", "1"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.Update(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memStore_BatchGetValues(t *testing.T) {
	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)
	m.RegisterKey("bar", "true", Bool)

	m.Update("bar", "false")

	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		args args
		want []*KV
	}{
		{
			args: args{[]string{"bar", "foo", ""}},
			want: []*KV{
				{Key: "bar", Value: false},
				{Key: "foo", Value: 1},
			},
		},
		{
			args: args{[]string{"loo", "zip", ""}},
			want: []*KV{},
		},
		{
			args: args{[]string{"bar", "loo", "foo", ""}},
			want: []*KV{
				{Key: "bar", Value: false},
				{Key: "foo", Value: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.BatchGetValues(tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchGetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memStore_BatchGetValueString(t *testing.T) {
	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)
	m.RegisterKey("bar", "true", Bool)

	m.Update("bar", "false")

	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		args args
		want []*KVStr
	}{
		{
			args: args{[]string{"bar", "foo", ""}},
			want: []*KVStr{
				{Key: "bar", Value: "false"},
				{Key: "foo", Value: "1"},
			},
		},
		{
			args: args{[]string{"loo", "zip", ""}},
			want: []*KVStr{},
		},
		{
			args: args{[]string{"bar", "loo", "foo", ""}},
			want: []*KVStr{
				{Key: "bar", Value: "false"},
				{Key: "foo", Value: "1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.BatchGetValueString(tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchGetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memStore_GetValue(t *testing.T) {
	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)
	m.RegisterKey("bar", "true", Bool)

	m.Update("bar", "false")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantV   interface{}
		wantHit bool
	}{
		{
			args:    args{"foo"},
			wantV:   1,
			wantHit: true,
		},
		{
			args:    args{"bar"},
			wantV:   false,
			wantHit: true,
		},
		{
			args:    args{"zoo"},
			wantV:   nil,
			wantHit: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotV, gotHit := m.GetValue(tt.args.key)
			if !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("GetValue() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotHit != tt.wantHit {
				t.Errorf("GetValue() gotHit = %v, want %v", gotHit, tt.wantHit)
			}
		})
	}
}

func Test_memStore_MustGetValue(t *testing.T) {

	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)
	m.RegisterKey("bar", "true", Bool)

	m.Update("bar", "false")

	type args struct {
		key string
	}
	tests := []struct {
		name      string
		args      args
		wantV     interface{}
		wantPanic bool
	}{
		{
			args:      args{"foo"},
			wantV:     1,
			wantPanic: false,
		},
		{
			args:      args{"bar"},
			wantV:     false,
			wantPanic: false,
		},
		{
			args:      args{"zoo"},
			wantV:     nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := func() {
				gotV := m.MustGetValue(tt.args.key)
				if !reflect.DeepEqual(gotV, tt.wantV) {
					t.Errorf("MustGetValue() gotV = %v, want %v", gotV, tt.wantV)
				}
			}
			if tt.wantPanic {
				shouldPanic(t, f)
			} else {
				shouldNotPanic(t, f)
			}

		})
	}
}

func Test_memStore_GetValueString(t *testing.T) {
	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)
	m.RegisterKey("bar", "true", Bool)

	m.Update("bar", "false")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantV   string
		wantHit bool
	}{
		{
			args:    args{"foo"},
			wantV:   "1",
			wantHit: true,
		},
		{
			args:    args{"bar"},
			wantV:   "false",
			wantHit: true,
		},
		{
			args:    args{"zoo"},
			wantV:   "",
			wantHit: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotV, gotHit := m.GetValueString(tt.args.key)
			if gotV != tt.wantV {
				t.Errorf("GetValueString() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotHit != tt.wantHit {
				t.Errorf("GetValueString() gotHit = %v, want %v", gotHit, tt.wantHit)
			}
		})
	}
}

func Test_memStore_MustGetValueString(t *testing.T) {

	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)
	m.RegisterKey("bar", "true", Bool)

	m.Update("bar", "false")

	type args struct {
		key string
	}
	tests := []struct {
		name      string
		args      args
		wantV     string
		wantPanic bool
	}{
		{
			args:      args{"foo"},
			wantV:     "1",
			wantPanic: false,
		},
		{
			args:      args{"bar"},
			wantV:     "false",
			wantPanic: false,
		},
		{
			args:      args{"zoo"},
			wantV:     "",
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			f := func() {
				gotV := m.MustGetValueString(tt.args.key)
				if gotV != tt.wantV {
					t.Errorf("GetValueString() gotV = %v, want %v", gotV, tt.wantV)
				}
			}
			if tt.wantPanic {
				shouldPanic(t, f)
			} else {
				shouldNotPanic(t, f)
			}

		})
	}

}

func Test_memStore_ResetKey(t *testing.T) {

	m := NewStore(NoopPersister, NoopLoadPolicy)

	m.RegisterKey("foo", "1", Int)
	m.RegisterKey("bar", "true", Bool)

	m.Update("bar", "false")

	m.ResetKey("bar")

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantV   string
		wantHit bool
	}{
		{
			args:    args{"foo"},
			wantV:   "1",
			wantHit: true,
		},
		{
			args:    args{"bar"},
			wantV:   "true",
			wantHit: true,
		},
		{
			args:    args{"zoo"},
			wantV:   "",
			wantHit: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotV, gotHit := m.GetValueString(tt.args.key)
			if gotV != tt.wantV {
				t.Errorf("GetValueString() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotHit != tt.wantHit {
				t.Errorf("GetValueString() gotHit = %v, want %v", gotHit, tt.wantHit)
			}
		})
	}

}

func shouldPanic(t *testing.T, f func()) {
	defer func() { recover() }()
	f()
	t.Errorf("did not panicked")
}

func shouldNotPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panicked")
		}
	}()
	f()
}
