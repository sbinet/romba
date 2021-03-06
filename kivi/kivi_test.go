// Copyright (c) 2013 Uwe Hoffmann. All rights reserved.

/*
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package kivi

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBasic(t *testing.T) {
	root, err := ioutil.TempDir("", "kivi_test")
	if err != nil {
		t.Fatalf("cannot open tempdir: %v", err)
	}

	kdb, err := Open(root, keySizeSha1)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	key := randomBytes(t, keySizeSha1)
	value := randomBytes(t, 50)

	err = kdb.Put(key, value)
	if err != nil {
		t.Fatal("failed to insert")
	}

	kdb.Flush()

	sval, err := kdb.Get(key)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	if sval == nil {
		t.Fatal("no value found")
	}

	if !bytes.Equal(value, sval) {
		t.Fatal("values differ")
	}
}

func TestMultipleDataFiles(t *testing.T) {
	root, err := ioutil.TempDir("", "kivi_test")
	if err != nil {
		t.Fatalf("cannot open tempdir: %v", err)
	}

	kdb, err := Open(root, keySizeSha1)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	key1 := randomBytes(t, keySizeSha1)
	value1 := randomBytes(t, 50)

	err = kdb.Put(key1, value1)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	err = kdb.Close()
	if err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	kdb, err = Open(root, keySizeSha1)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	key2 := randomBytes(t, keySizeSha1)
	value2 := randomBytes(t, 50)

	err = kdb.Put(key2, value2)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	err = kdb.Close()
	if err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	kdb, err = Open(root, keySizeSha1)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	sval1, err := kdb.Get(key1)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	if sval1 == nil {
		t.Fatal("no value found")
	}

	if !bytes.Equal(value1, sval1) {
		t.Fatal("values differ")
	}

	sval2, err := kdb.Get(key2)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	if sval2 == nil {
		t.Fatal("no value found")
	}

	if !bytes.Equal(value2, sval2) {
		t.Fatal("values differ")
	}
}

func TestNoKeydirFile(t *testing.T) {
	root, err := ioutil.TempDir("", "kivi_test")
	if err != nil {
		t.Fatalf("cannot open tempdir: %v", err)
	}

	kdb, err := Open(root, keySizeSha1)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	key1 := randomBytes(t, keySizeSha1)
	value1 := randomBytes(t, 50)

	err = kdb.Put(key1, value1)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	err = kdb.Close()
	if err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	err = os.Remove(filepath.Join(root, keydirFilename+"_0"))
	if err != nil {
		t.Fatalf("failed to delete keydir file: %v", err)
	}

	kdb, err = Open(root, keySizeSha1)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	sval1, err := kdb.Get(key1)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	if sval1 == nil {
		t.Fatal("no value found")
	}

	if !bytes.Equal(value1, sval1) {
		t.Fatal("values differ")
	}
}
