// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package catalog

import (
	"fmt"
	"io"
)

type TableMVCCNode struct {
	// history schema
	Schema *Schema
}

func NewEmptyTableMVCCNode() *TableMVCCNode {
	return &TableMVCCNode{}
}

func (e *TableMVCCNode) CloneAll() *TableMVCCNode {
	return &TableMVCCNode{
		Schema: e.Schema.Clone(),
	}
}

func (e *TableMVCCNode) CloneData() *TableMVCCNode {
	return e.CloneAll()
}

func (e *TableMVCCNode) String() string {
	return fmt.Sprintf("schema.v.%d", e.Schema.Version)
}

// for create drop in one txn
func (e *TableMVCCNode) Update(un *TableMVCCNode) {
	e.Schema = un.Schema
}

func (e *TableMVCCNode) WriteTo(w io.Writer) (n int64, err error) {
	var schemaBuf []byte
	if schemaBuf, err = e.Schema.Marshal(); err != nil {
		return
	}
	if _, err = w.Write(schemaBuf); err != nil {
		return
	}
	n += int64(len(schemaBuf))
	return
}

func (e *TableMVCCNode) ReadFrom(r io.Reader) (n int64, err error) {
	e.Schema = NewEmptySchema("")
	if n, err = e.Schema.ReadFrom(r); err != nil {
		return
	}
	return
}

type TableNode struct {
	// The latest schema. A shortcut to the schema in the last mvvcnode.
	// TODO(aptend): use atomic.Pointer?
	schema *Schema
}

func (node *TableNode) WriteTo(w io.Writer) (n int64, err error) {
	// do not writeTo inherit from mvvcnode in replay phrase
	// reference: function onReplayCreateTable and onReplayUpdateTable
	return
}

func (node *TableNode) ReadFrom(r io.Reader) (n int64, err error) {
	// do not readFrom, inherit from mvvcnode in replay phrase
	// reference: function onReplayCreateTable and onReplayUpdateTable
	return
}
