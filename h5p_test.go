// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hdf5

import (
	"testing"
)

func TestPropList(t *testing.T) {
	propList, err := NewPropList(H5P_DATASET_CREATE)
	if err != nil {
		t.Error(err)
	}
	err = propList.SetChunk([]uint{100, 100})
	if err != nil {
		t.Error(err)
	}

}
