package confstore

import (
	"github.com/go-redis/redis/v7"
	"sync"
)

var (
	NoopPersister Persister = &noopPersister{} // Never persist config key-values
)

// never persist config key-values
type noopPersister struct{}

func (n *noopPersister) Load(s Interface) {
}

func (n *noopPersister) Save(key, value string) error {
	return nil
}

func (n *noopPersister) BatchSave(strs []*KVStr) error {
	return nil
}

func (n *noopPersister) Delete(key string) error {
	return nil
}

// Create a redis Persister
func NewRedisPersister(redisClient *redis.Client, configKeyRoot string) Persister {
	return &redisPersister{
		redisClient:   redisClient,
		configRootKey: configKeyRoot,
	}
}

type redisPersister struct {
	redisClient   *redis.Client
	configRootKey string
}

func (r *redisPersister) Load(s Interface) {
	result, err := r.redisClient.HGetAll(r.configRootKey).Result()
	if err != nil {
		errLogger.Println("load configs from redis failed:", err)
		return
	}
	for k, v := range result {
		kvError := s.UpdateNoPersist(k, v)
		if kvError != nil {
			stdLogger.Printf("warning: load key[%s] value[%s] error: %s\n", kvError.Key, kvError.Value, kvError.Error)
		}
	}
	stdLogger.Printf("%d config keys loaded from redis\n", len(result))
}

func (r *redisPersister) Save(key, value string) error {
	return r.BatchSave([]*KVStr{{key, value}})
}

func (r *redisPersister) BatchSave(kvStrs []*KVStr) error {
	args := make([]string, 0, 2*len(kvStrs))
	for _, kvStr := range kvStrs {
		args = append(args, kvStr.Key)
		args = append(args, kvStr.Value)
	}
	if len(args) == 0 {
		return nil
	}
	_, error := r.redisClient.HMSet(r.configRootKey, args).Result()
	if error != nil {
		return error
	}
	return nil
}

func (r *redisPersister) Delete(key string) error {
	_, err := r.redisClient.HDel(r.configRootKey, key).Result()
	return err
}

func NewSimpleLoadPolicy() LoadPolicy {
	return &simpleLoadPolicy{
		loadLock: &sync.Mutex{},
	}
}

type simpleLoadPolicy struct {
	loadLock *sync.Mutex
}

func (d *simpleLoadPolicy) DoLoad(s Interface, p Persister) {
	d.loadLock.Lock()
	defer d.loadLock.Unlock()
	p.Load(s)
}
