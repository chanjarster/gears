package confstore

import (
	"errors"
	"fmt"
	"github.com/modern-go/reflect2"
	"sync"
)

type defaultImpl struct {
	kvProcessorLock *sync.RWMutex
	kvProcessor     map[string]ValueProcessor

	kvLock       *sync.RWMutex
	kvDefault    map[string]interface{}
	kvDefaultStr map[string]string
	kv           map[string]interface{}
	kvStr        map[string]string
	persister    Persister
	loadPolicy   LoadPolicy
}


func NewStore(persister Persister, policy LoadPolicy) Interface {
	return &defaultImpl{
		kvProcessorLock: &sync.RWMutex{},
		kvProcessor:     make(map[string]ValueProcessor),
		kvLock:          &sync.RWMutex{},
		kvDefault:       make(map[string]interface{}),
		kvDefaultStr:    make(map[string]string),
		kv:              make(map[string]interface{}),
		kvStr:           make(map[string]string),
		persister:       persister,
		loadPolicy:      policy,
	}
}

func (m *defaultImpl) RegisterKey(key string, defaultValue string, vp ValueProcessor) error {

	if key == "" {
		return errors.New("key is required")
	}
	if reflect2.IsNil(vp) {
		return errors.New("ValueProcessor is required")
	}

	m.kvProcessorLock.Lock()
	defer m.kvProcessorLock.Unlock()

	m.kvLock.Lock()
	defer m.kvLock.Unlock()

	if _, hit := m.kvProcessor[key]; hit {
		return errors.New(fmt.Sprintf("key[%s] ValueProcessor duplicated", key))
	}

	if _, hit := m.kvDefault[key]; hit {
		return errors.New(fmt.Sprintf("key[%s] default value duplicated", key))
	}

	if ok, err := vp.Validate(defaultValue); !ok {
		return errors.New(err)
	}

	m.kvDefault[key] = vp.Convert(defaultValue)
	m.kvDefaultStr[key] = defaultValue

	m.kvProcessor[key] = vp

	return nil
}

func (m *defaultImpl) DeregisterKey(key string) {
	if key == "" {
		return
	}

	m.kvProcessorLock.Lock()
	defer m.kvProcessorLock.Unlock()

	m.kvLock.Lock()
	defer m.kvLock.Unlock()

	delete(m.kvProcessor, key)
	delete(m.kv, key)
	delete(m.kvStr, key)
	delete(m.kvDefault, key)
	delete(m.kvDefaultStr, key)

	err := m.persister.Delete(key)
	if err != nil {
		errLogger.Println("perister delete key error:", err)
	}

}

func (m *defaultImpl) BatchUpdate(kvStrs []*KVStr) []*KVError {
	errors := make([]*KVError, 0, len(kvStrs))
	if len(kvStrs) == 0 {
		return nil
	}

	for _, kvStr := range kvStrs {
		error := m.UpdateNoPersist(kvStr.Key, kvStr.Value)
		if error != nil {
			errors = append(errors, error)
		}
	}

	// batch persist kvStrs
	errorKeys := make(map[string]bool)
	for _, kvError := range errors {
		errorKeys[kvError.Key] = true
	}
	okKvStrs := make([]*KVStr, 0, len(kvStrs))
	for _, kvStr := range kvStrs {
		if errorKeys[kvStr.Key] {
			continue
		}
		okKvStrs = append(okKvStrs, kvStr)
	}
	m.persister.BatchSave(okKvStrs)
	return errors
}

func (m *defaultImpl) Update(key string, value string) *KVError {
	err := m.UpdateNoPersist(key, value)
	if err != nil {
		err2 := m.persister.Save(key, value)
		if err2 != nil {
			errLogger.Printf("error happened when persist key[%s] value[%s]\n", key, value)
			return nil
		}
	}
	return err
}

func (m *defaultImpl) UpdateNoPersist(key string, value string) *KVError {
	if key == "" {
		return nil
	}

	m.kvProcessorLock.RLock()
	p, hit := m.kvProcessor[key]
	m.kvProcessorLock.RUnlock()

	if !hit {
		return &KVError{
			Key:   key,
			Value: value,
			Error: fmt.Sprintf("key[%s] not registered", key),
		}
	}

	if ok, err := p.Validate(value); !ok {
		return &KVError{
			Key:   key,
			Value: value,
			Error: err,
		}
	}

	m.kvLock.Lock()
	m.kv[key] = p.Convert(value)
	m.kvStr[key] = value
	m.kvLock.Unlock()
	return nil
}

func (m *defaultImpl) BatchGetValues(keys []string) []*KV {
	m.loadPolicy.DoLoad(m, m.persister)

	kvs := make([]*KV, 0, len(keys))
	for _, key := range keys {
		if key == "" {
			continue
		}
		v, hit := m.GetValue(key)
		if !hit {
			continue
		}
		kvs = append(kvs, &KV{key, v})
	}
	return kvs
}

func (m *defaultImpl) BatchGetValueString(keys []string) []*KVStr {
	m.loadPolicy.DoLoad(m, m.persister)

	kvs := make([]*KVStr, 0, len(keys))
	for _, key := range keys {
		if key == "" {
			continue
		}
		v, hit := m.GetValueString(key)
		if !hit {
			continue
		}
		kvs = append(kvs, &KVStr{key, v})
	}
	return kvs
}

func (m *defaultImpl) MustGetValue(key string) interface{} {
	m.loadPolicy.DoLoad(m, m.persister)

	v, hit := m.GetValue(key)
	if !hit {
		panic(fmt.Sprintf("config key[%s] missing", key))
	}
	return v
}

func (m *defaultImpl) MustGetValueString(key string) string {
	m.loadPolicy.DoLoad(m, m.persister)

	v, hit := m.GetValueString(key)
	if !hit {
		panic(fmt.Sprintf("config key[%s] missing", key))
	}
	return v
}

func (m *defaultImpl) GetValue(key string) (v interface{}, hit bool) {
	m.loadPolicy.DoLoad(m, m.persister)

	m.kvLock.RLock()
	defer m.kvLock.RUnlock()

	v, hit = m.kv[key]
	if !hit {
		v, hit = m.kvDefault[key]
		return v, hit
	}
	return v, hit

}

func (m *defaultImpl) GetValueString(key string) (v string, hit bool) {
	m.loadPolicy.DoLoad(m, m.persister)

	m.kvLock.RLock()
	defer m.kvLock.RUnlock()

	v, hit = m.kvStr[key]
	if !hit {
		v, hit = m.kvDefaultStr[key]
		return v, hit

	}
	return v, hit

}

func (m *defaultImpl) ResetKey(key string) {
	m.kvLock.Lock()
	defer m.kvLock.Unlock()

	delete(m.kv, key)
	delete(m.kvStr, key)
	err := m.persister.Delete(key)
	if err != nil {
		errLogger.Println("perister delete key error:", err)
	}
}