// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package datavalues

import (
	"bytes"
	"unsafe"

	"base/docs"
)

type ValueTuple struct {
	fields []IDataValue
}

func MakeTuple(v ...IDataValue) IDataValue {
	return &ValueTuple{fields: v}
}

func ZeroTuple() IDataValue {
	return &ValueTuple{fields: nil}
}

func (v *ValueTuple) Size() uintptr {
	size := unsafe.Sizeof(*v)
	for _, field := range v.fields {
		size += field.Size()
	}
	return size
}

func (v *ValueTuple) Show() []byte {
	result := make([][]byte, len(v.fields))
	for i := range v.fields {
		result[i] = v.fields[i].Show()
	}
	return bytes.Join(result, []byte{0x01})
}

func (v *ValueTuple) GetType() Type {
	return TypeTuple
}

func (v *ValueTuple) AsSlice() []IDataValue {
	return v.fields
}

func (v *ValueTuple) Compare(other IDataValue) (Comparison, error) {
	otherv := other.(*ValueTuple)
	for i := range v.fields {
		cmp, err := v.fields[i].Compare(otherv.fields[i])
		if err != nil {
			return 0, err
		}
		switch cmp {
		case Equal:
			continue
		default:
			return cmp, nil
		}
	}
	return 0, nil
}

func (v *ValueTuple) Document() docs.Documentation {
	return docs.Text("Tuple")
}