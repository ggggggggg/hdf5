// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hdf5

// #include "hdf5.h"
// #include <stdlib.h>
// #include <string.h>
// inline static
// hid_t _go_hdf5_H5P_DEFAULT() { return H5P_DEFAULT; }
// hid_t _go_hdf5_H5P_DATASET_CREATE() { return H5P_DATASET_CREATE; }
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type PropType C.hid_t

type PropList struct {
	Identifier
}

var (
	H5P_DATASET_CREATE PropType = PropType(C._go_hdf5_H5P_DATASET_CREATE())
)

var (
	P_DEFAULT *PropList = newPropList(C._go_hdf5_H5P_DEFAULT())
)

func newPropList(id C.hid_t) *PropList {
	p := &PropList{Identifier{id}}
	runtime.SetFinalizer(p, (*PropList).finalizer)
	return p
}

func (p *PropList) finalizer() {
	if err := p.Close(); err != nil {
		panic(fmt.Errorf("error closing PropList: %s", err))
	}
}

// NewPropList creates a new PropList as an instance of a property list class.
func NewPropList(cls_id PropType) (*PropList, error) {
	hid := C.H5Pcreate(C.hid_t(cls_id))
	if err := checkID(hid); err != nil {
		return nil, err
	}
	return newPropList(hid), nil
}

func (p *PropList) SetChunk(chunkDims []uint) error {
	rank := len(chunkDims)
	var c_chunkDims *_Ctype_ulonglong
	if chunkDims != nil {
		c_chunkDims = (*C.hsize_t)(unsafe.Pointer(&chunkDims))
	}
	c_rank := (_Ctype_int)(rank)
	err := h5err(C.H5Pset_chunk(p.id, c_rank, c_chunkDims))
	return err
}

// Close terminates access to a PropList.
func (p *PropList) Close() error {
	if p.id == 0 {
		return nil
	}
	err := h5err(C.H5Pclose(p.id))
	p.id = 0
	return err
}

// Copy copies an existing PropList to create a new PropList.
func (p *PropList) Copy() (*PropList, error) {
	hid := C.H5Pcopy(p.id)
	if err := checkID(hid); err != nil {
		return nil, err
	}
	return newPropList(hid), nil
}
