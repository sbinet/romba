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

package archive

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"github.com/uwedeportivo/torrentzip/cgzip"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

const (
	zipSuffix  = ".zip"
	gzipSuffix = ".gz"
)

type hashes struct {
	crc  []byte
	md5  []byte
	sha1 []byte
}

func newHashes() *hashes {
	rs := new(hashes)
	rs.crc = make([]byte, 0, crc32.Size)
	rs.md5 = make([]byte, 0, md5.Size)
	rs.sha1 = make([]byte, 0, sha1.Size)
	return rs
}

func (hh *hashes) forFile(inpath string) error {
	file, err := os.Open(inpath)
	if err != nil {
		return err
	}
	defer file.Close()

	return hh.forReader(file)
}

func (hh *hashes) forReader(in io.Reader) error {
	br := bufio.NewReader(in)

	hSha1 := sha1.New()
	hMd5 := md5.New()
	hCrc := cgzip.NewCrc32()

	w := io.MultiWriter(hSha1, hMd5, hCrc)

	_, err := io.Copy(w, br)
	if err != nil {
		return err
	}

	hh.crc = hCrc.Sum(hh.crc[0:0])
	hh.md5 = hMd5.Sum(hh.md5[0:0])
	hh.sha1 = hSha1.Sum(hh.sha1[0:0])

	return nil
}

func hashesForFile(inpath string) (*hashes, error) {
	file, err := os.Open(inpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return hashesForReader(file)
}

func hashesForReader(in io.Reader) (*hashes, error) {
	hSha1 := sha1.New()
	hMd5 := md5.New()
	hCrc := crc32.NewIEEE()

	w := io.MultiWriter(hSha1, hMd5, hCrc)

	_, err := io.Copy(w, in)
	if err != nil {
		return nil, err
	}

	res := new(hashes)
	res.crc = hCrc.Sum(nil)
	res.md5 = hMd5.Sum(nil)
	res.sha1 = hSha1.Sum(nil)

	return res, nil
}

func sha1ForFile(inpath string) ([]byte, error) {
	file, err := os.Open(inpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return sha1ForReader(file)
}

func sha1ForReader(in io.Reader) ([]byte, error) {
	h := sha1.New()

	_, err := io.Copy(h, in)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func pathFromSha1HexEncoding(root, hexStr, suffix string) string {
	prefix := hexStr[0:8]
	pieces := make([]string, 6)

	pieces[0] = root
	for i := 0; i < 4; i++ {
		pieces[i+1] = prefix[2*i : 2*i+2]
	}
	pieces[5] = hexStr + suffix

	return filepath.Join(pieces...)
}

func pathExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func archive(outpath string, r io.Reader) error {
	br := bufio.NewReader(r)

	err := os.MkdirAll(filepath.Dir(outpath), 0777)
	if err != nil {
		return err
	}

	outfile, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	bufout := bufio.NewWriter(outfile)
	defer bufout.Flush()

	zipWriter := cgzip.NewWriter(bufout)
	defer zipWriter.Close()

	_, err = io.Copy(zipWriter, br)
	if err != nil {
		return err
	}

	return nil
}