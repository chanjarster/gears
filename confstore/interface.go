package configstore

import (
	"fmt"
	"log"
	"os"
)

// TODO write examples
// A interface to manage config key-values
type Interface interface {
	// register config key
	//
	// duplicated key not allowed
	//
	// default string value must be valid to ValueProcessor
	RegisterKey(key string, defaultValue string, v ValueProcessor) error

	// deregister config key, noop if key doesn't exist
	DeregisterKey(key string)

	// batch update config key-values
	//
	// key must be registered before
	BatchUpdate(kvStrs []*KVStr) []*KVError

	// update config key-value and persist.
	// ValueProcessor will handle string value
	Update(key string, value string) *KVError

	// update config key-value but do not persist.
	// ValueProcessor will handle string value
	UpdateNoPersist(key string, value string) *KVError

	// batch get config key-values
	//
	// if key doesn't exist, no error will return
	BatchGetValues(keys []string) []*KV

	// batch get config key-value strings
	//
	// if key doesn't exist, no error will return
	BatchGetValueString(keys []string) []*KVStr

	// get config key-value
	//
	// panic if key doesn't exist
	MustGetValue(key string) interface{}

	// get config key-value string
	//
	// panic if key doesn't exist
	MustGetValueString(key string) string

	// get config key-value
	//
	// if key doesn't exist, return nil, false
	GetValue(key string) (v interface{}, hit bool)

	// get config key-value string
	//
	// if key doesn't exist, return nil, false
	GetValueString(key string) (v string, hit bool)

	// reset key's current value to default value
	ResetKey(key string)
}

// used to validate string value and convert string value to specific type
type ValueProcessor interface {
	Validate(value string) (ok bool, err string)
	Convert(value string) interface{}
}

// config key persistence layer
type Persister interface {
	Load(s Interface)
	Save(key, value string) error
	BatchSave([]*KVStr) error
	Delete(key string) error
}

type LoadPolicy interface {
	// load config key-values from Persister to Interface
	DoLoad(s Interface, p Persister)
}

type KVStr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (K *KVStr) String() string {
	return fmt.Sprintf("key: %s, value: %s", K.Key, K.Value)
}

type KVError struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Error string `json:"error"`
}

func (K *KVError) String() string {
	return fmt.Sprintf("key: %s, value: %s, error: %s", K.Key, K.Value, K.Error)
}

type KV struct {
	Key   string
	Value interface{}
}

func (K *KV) String() string {
	return fmt.Sprintf("key: %s, value: %s", K.Key, K.Value)
}

var stdLogger = log.New(os.Stdout, "", log.Ldate|log.LstdFlags|log.Llongfile)
var errLogger = log.New(os.Stderr, "", log.Ldate|log.LstdFlags|log.Llongfile)
