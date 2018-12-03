package common

import (
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/goloop/module"
	"github.com/pkg/errors"
	"math/big"
)

const (
	TypeAddress = iota + codec.TypeCustom
	TypeInt
)

type typeCodec struct{}

var TypeCodec = &typeCodec{}

func (*typeCodec) Decode(tag uint8, data []byte) (interface{}, error) {
	switch tag {
	case TypeAddress:
		return NewAddress(data), nil
	case TypeInt:
		i := new(HexInt)
		i.SetBytes(data)
		return i, nil
	default:
		return 0, errors.Errorf("InvalidTypeTag:%d", tag)
	}
}

func (*typeCodec) Encode(o interface{}) (uint8, []byte, error) {
	switch obj := o.(type) {
	case module.Address:
		return TypeAddress, obj.Bytes(), nil
	case *big.Int:
		return TypeInt, BigIntToBytes(obj), nil
	case *HexInt:
		return TypeInt, obj.Bytes(), nil
	case int:
		return TypeInt, Int64ToBytes(int64(obj)), nil
	case int16:
		return TypeInt, Int64ToBytes(int64(obj)), nil
	case int32:
		return TypeInt, Int64ToBytes(int64(obj)), nil
	case int64:
		return TypeInt, Int64ToBytes(obj), nil
	case uint:
		return TypeInt, Int64ToBytes(int64(obj)), nil
	case uint16:
		return TypeInt, Int64ToBytes(int64(obj)), nil
	case uint32:
		return TypeInt, Int64ToBytes(int64(obj)), nil
	case uint64:
		var bi big.Int
		bi.SetUint64(o.(uint64))
		return TypeInt, BigIntToBytes(&bi), nil
	default:
		return 0, nil, errors.Errorf("UnknownType:%T", obj)
	}
}

func MarshalAny(obj interface{}) ([]byte, error) {
	return codec.MarshalAny(TypeCodec, obj)
}

func EncodeAny(obj interface{}) (interface{}, error) {
	return codec.EncodeAny(TypeCodec, obj)
}

func UnmarshalAny(bs []byte) (interface{}, error) {
	return codec.UnmarshalAny(TypeCodec, bs)
}