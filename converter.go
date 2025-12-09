package tries

import (
	"encoding/binary"
	"fmt"
)

type (
	converter[T TrieKey] interface {
		Load(value T) error
		Next() (value uint8, ok bool)
	}

	converterUInt8[T TrieKey] struct {
		position uint8
		value    uint8
	}

	converterInt8[T TrieKey] struct {
		position uint8
		value    uint8
	}

	converterUInt16[T TrieKey] struct {
		position uint8
		value    [2]uint8
	}

	converterInt16[T TrieKey] struct {
		position uint8
		value    [2]uint8
	}

	converterUInt32[T TrieKey] struct {
		position uint8
		value    [4]uint8
	}

	converterInt32[T TrieKey] struct {
		position uint8
		value    [4]uint8
	}

	converterUInt64[T TrieKey] struct {
		position uint8
		value    [8]uint8
	}

	converterInt64[T TrieKey] struct {
		position uint8
		value    [8]uint8
	}

	converterString[T TrieKey] struct {
		position int
		value    string
	}

	converterInt8Slice[TItem ~uint8 | ~int8, TKey TrieKey] struct {
		position int
		value    []TItem
	}

	converterIntSlice[TItem ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint64 | ~int64, TKey TrieKey] struct {
		position     int
		subConverter converter[TItem]
		value        []TItem
	}

	converterTransforms[T TrieKey] struct {
		subConverter converter[T]
		transforms   []TransformFunc
	}
)

// ----- uint8 -----
func (this *converterUInt8[T]) Load(value T) error {
	this.position = 0
	this.value = any(value).(uint8)
	return nil
}

func (this *converterUInt8[T]) Next() (value uint8, ok bool) {
	const oneByte = 0
	if this.position > oneByte {
		return 0, false
	}

	this.position = 1
	return this.value, true
}

// ----- int8 -----
func (this *converterInt8[T]) Load(value T) error {
	this.position = 0
	this.value = uint8(any(value).(int8)) //nolint:gosec // this casting is fine
	return nil
}

func (this *converterInt8[T]) Next() (value uint8, ok bool) {
	const oneByte = 0
	if this.position > oneByte {
		return 0, false
	}

	this.position = 1
	return this.value, true
}

// ----- uint16 -----
func (this *converterUInt16[T]) Load(value T) error {
	this.position = 0
	binary.BigEndian.PutUint16(this.value[:2], any(value).(uint16))
	return nil
}

func (this *converterUInt16[T]) Next() (value uint8, ok bool) {
	const twoBytes = 1
	if this.position > twoBytes {
		return 0, false
	}

	value = this.value[this.position]
	this.position++
	return value, true
}

// ----- int16 -----
func (this *converterInt16[T]) Load(value T) error {
	this.position = 0
	binary.BigEndian.PutUint16(this.value[:2], uint16(any(value).(int16))) //nolint:gosec // this casting is fine
	return nil
}

func (this *converterInt16[T]) Next() (value uint8, ok bool) {
	const twoBytes = 1
	if this.position > twoBytes {
		return 0, false
	}

	value = this.value[this.position]
	this.position++
	return value, true
}

// ----- uint32 -----
func (this *converterUInt32[T]) Load(value T) error {
	this.position = 0
	binary.BigEndian.PutUint32(this.value[:4], any(value).(uint32))
	return nil
}

func (this *converterUInt32[T]) Next() (value uint8, ok bool) {
	const fourBytes = 3
	if this.position > fourBytes {
		return 0, false
	}

	value = this.value[this.position]
	this.position++
	return value, true
}

// ----- int32 -----
func (this *converterInt32[T]) Load(value T) error {
	this.position = 0
	binary.BigEndian.PutUint32(this.value[:4], uint32(any(value).(int32))) //nolint:gosec // this casting is fine
	return nil
}

func (this *converterInt32[T]) Next() (value uint8, ok bool) {
	const fourBytes = 3
	if this.position > fourBytes {
		return 0, false
	}

	value = this.value[this.position]
	this.position++
	return value, true
}

// ----- uint64 -----
func (this *converterUInt64[T]) Load(value T) error {
	this.position = 0
	binary.BigEndian.PutUint64(this.value[:8], any(value).(uint64))
	return nil
}

func (this *converterUInt64[T]) Next() (value uint8, ok bool) {
	const eightBytes = 7
	if this.position > eightBytes {
		return 0, false
	}

	value = this.value[this.position]
	this.position++
	return value, true
}

// ----- int64 -----
func (this *converterInt64[T]) Load(value T) error {
	this.position = 0
	binary.BigEndian.PutUint64(this.value[:8], uint64(any(value).(int64))) //nolint:gosec // this casting is fine
	return nil
}

func (this *converterInt64[T]) Next() (value uint8, ok bool) {
	const eightBytes = 7
	if this.position > eightBytes {
		return 0, false
	}

	value = this.value[this.position]
	this.position++
	return value, true
}

// ----- string -----
func (this *converterString[T]) Load(value T) error {
	this.position = 0
	switch v := any(value).(type) {
	case string:
		this.value = v
		return nil
	default:
		return fmt.Errorf("%w: unable to convert %T to a string", ErrorBadTrieKey, value)
	}
}

func (this *converterString[T]) Next() (value uint8, ok bool) {
	if this.position >= len(this.value) {
		return 0, false
	}

	value = this.value[this.position]
	this.position++
	return value, true
}

// ----- []int8 -----
func (this *converterInt8Slice[TItem, TKey]) Load(value TKey) error {
	this.position = 0
	switch v := any(value).(type) {
	case []TItem:
		this.value = v
		return nil
	default:
		return fmt.Errorf("%w: unable to convert %T to an integer slice", ErrorBadTrieKey, value)
	}
}

func (this *converterInt8Slice[TItem, TKey]) Next() (value uint8, ok bool) {
	if this.position >= len(this.value) {
		return 0, false
	}

	value = uint8(this.value[this.position])
	this.position++
	return value, true
}

// ----- []int(x) -----
func (this *converterIntSlice[TItem, TKey]) Load(value TKey) error {
	this.position = 0
	switch v := any(value).(type) {
	case []TItem:
		this.value = v
		if len(this.value) > 0 {
			this.subConverter.Load(this.value[0])
		}

		return nil
	default:
		return fmt.Errorf("%w: unable to convert %T to an integer slice", ErrorBadTrieKey, value)
	}
}

func (this *converterIntSlice[TItem, TKey]) Next() (value uint8, ok bool) {
	if this.position >= len(this.value) {
		return 0, false
	}

	value, ok = this.subConverter.Next()
	if !ok {
		this.position++
		if this.position >= len(this.value) {
			return 0, false
		}

		this.subConverter.Load(this.value[this.position])
		value, _ = this.subConverter.Next()
	}

	return value, true
}

// ----- transforms -----
func (this *converterTransforms[T]) Load(value T) error {
	return this.subConverter.Load(value)
}

func (this *converterTransforms[T]) Next() (value uint8, ok bool) {
	for {
		value, ok = this.subConverter.Next()
		if !ok {
			return value, ok
		}

		passes := true
		for _, transform := range this.transforms {
			value, passes = transform(value)
			if !passes {
				break
			}
		}

		if passes {
			return value, ok
		}
	}
}

func selectConverter[T TrieKey]() (converter converter[T], err error) {
	var dummy T
	switch any(dummy).(type) {
	case uint8:
		return new(converterUInt8[T]), nil
	case int8:
		return new(converterInt8[T]), nil
	case uint16:
		return new(converterUInt16[T]), nil
	case int16:
		return new(converterInt16[T]), nil
	case uint32:
		return new(converterUInt32[T]), nil
	case int32:
		return new(converterInt32[T]), nil
	case uint64:
		return new(converterUInt64[T]), nil
	case int64:
		return new(converterInt64[T]), nil
	case string:
		return new(converterString[T]), nil
	case []uint8:
		return new(converterInt8Slice[uint8, T]), nil
	case []int8:
		return new(converterInt8Slice[int8, T]), nil
	case []uint16:
		return &converterIntSlice[uint16, T]{subConverter: new(converterUInt16[uint16])}, nil
	case []int16:
		return &converterIntSlice[int16, T]{subConverter: new(converterInt16[int16])}, nil
	case []uint32:
		return &converterIntSlice[uint32, T]{subConverter: new(converterUInt32[uint32])}, nil
	case []int32:
		return &converterIntSlice[int32, T]{subConverter: new(converterInt32[int32])}, nil
	case []uint64:
		return &converterIntSlice[uint64, T]{subConverter: new(converterUInt64[uint64])}, nil
	case []int64:
		return &converterIntSlice[int64, T]{subConverter: new(converterInt64[int64])}, nil
	default:
		return nil, fmt.Errorf("%w: no converter is defined for type %T", ErrorBadTrieKey, dummy)
	}
}

func wrapConverter[T TrieKey](converter converter[T], transforms []TransformFunc) converter[T] {
	return &converterTransforms[T]{
		subConverter: converter,
		transforms:   transforms,
	}
}
