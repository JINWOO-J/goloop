package scoredb

import (
	"github.com/icon-project/goloop/common/crypto"
	"github.com/pkg/errors"
)

type DictDB struct {
	key   []byte
	store StateStore
	depth int
}

func NewDictDB(store StateStore, name interface{}, depth int) *DictDB {
	kbytes := ToKey(DictDBPrefix, name)
	return &DictDB{
		key:   kbytes,
		store: store,
		depth: depth,
	}
}

func (d *DictDB) keyBytesForKeys(keys ...interface{}) []byte {
	return AppendKeys(d.key, keys...)
}

func (d *DictDB) GetDB(keys ...interface{}) *DictDB {
	if len(keys) >= d.depth {
		return nil
	}

	kbytes := d.keyBytesForKeys(keys...)

	return &DictDB{
		key:   kbytes,
		store: d.store,
		depth: d.depth - len(keys),
	}
}

func (d *DictDB) Get(keys ...interface{}) Value {
	if len(keys) != d.depth {
		return nil
	}

	kbytes := d.keyBytesForKeys(keys...)

	bs, err := d.store.GetValue(crypto.SHA3Sum256(kbytes))
	if err != nil || bs == nil {
		return nil
	}
	return NewValueFromBytes(bs)
}

func (d *DictDB) Set(params ...interface{}) error {
	if len(params) != d.depth+1 {
		return errors.New("IllegalArgument")
	}

	kbytes := d.keyBytesForKeys(params[:len(params)-1]...)
	v := params[len(params)-1]

	return d.store.SetValue(crypto.SHA3Sum256(kbytes), ToBytes(v))
}

func (d *DictDB) Delete(kv ...interface{}) error {
	if len(kv) != d.depth {
		return errors.New("IllegalArgument")
	}

	kbytes := d.keyBytesForKeys(kv...)

	return d.store.DeleteValue(crypto.SHA3Sum256(kbytes))
}