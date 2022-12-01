// ComProReader - competitive programming input routines for golang
package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const TESTDATA = "./testdata/"

// sample main
func main() {
	//MakeTestDataPerm(TESTDATA+"perm_1e7I.txt", 10_000_000, 0)
	// f, err := os.Open(TESTDATA + "perm_1e7I.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// os.Stdin = f

	// time0 := time.Now()

	cr := NewReader(os.Stdin)
	N := cr.Int()
	A := make([]int, N)
	for i := 0; i < N; i++ {
		A[i] = cr.Int()
	}

	// fmt.Fprintf(os.Stderr, "%v\n", time.Since(time0))

	sum := 0
	for _, a := range A {
		sum += a
	}
	fmt.Println((N * (N - 1)) / 2)
	fmt.Println(sum)
}

// copy below from here
type Reader struct {
	io.Reader
	buf   []byte
	index int
	//EOF   bool
}

func NewReader(r *os.File) *Reader {
	cr := &Reader{Reader: r}
	cr.buf = make([]byte, 0, 64*1024)
	cr.nextRead()
	return cr
}

func (cr *Reader) nextRead() {
	n, err := cr.Read(cr.buf[:cap(cr.buf)])
	/*if err == io.EOF {
		cr.EOF = true
	} else */if err != nil {
		panic(err)
	}
	cr.buf = cr.buf[:n]
	cr.index = 0
}

func (cr *Reader) Byte() (b byte) {
	if cr.index >= len(cr.buf) {
		cr.nextRead()
	}
	b = cr.buf[cr.index]
	cr.index++
	return
}

func (cr *Reader) skipCRLF(b byte) { // for MS-DOS
	if b == '\r' {
		cr.Byte() // == '\n'
	}
}

func (cr *Reader) Uint64() (u uint64) {
	u = 0
	b := cr.Byte()
	if '0' > b || b > '9' {
		panic("not decimal")
	}
	for '0' <= b && b <= '9' {
		u = 10*u + uint64(b-'0')
		b = cr.Byte()
	}
	cr.skipCRLF(b)
	return
}

func (cr *Reader) Int64() int64 {
	sign := int64(+1)
	b := cr.Byte()
	if b == '-' {
		sign = -1
	} else if b != '+' {
		cr.index--
	}
	return sign * int64(cr.Uint64())
}

func (cr *Reader) Int() int       { return int(cr.Int64()) }
func (cr *Reader) Int32() int32   { return int32(cr.Int64()) }
func (cr *Reader) Uint() uint     { return uint(cr.Int64()) }
func (cr *Reader) Uint32() uint32 { return uint32(cr.Int64()) }

// func (cr *Reader) II() (int, int) {
// 	return cr.Int(), cr.Int()
// }

// func (cr *Reader) III() (int, int, int) {
// 	return cr.Int(), cr.Int(), cr.Int()
// }

// func (cr *Reader) A(n int) []int {
// 	a := make([]int, n)
// 	for i := 0; i < n; i++ {
// 		a[i] = cr.Int()
// 	}
// 	return a
// }

// func (cr *Reader) AA(n int) ([]int, []int) {
// 	a := make([]int, n)
// 	b := make([]int, n)
// 	for i := 0; i < n; i++ {
// 		a[i], b[i] = cr.Int(), cr.Int()
// 	}
// 	return a, b
// }

func (cr *Reader) Float64x() float64 {
	sign := +1.0
	b := cr.Byte()
	if b == '+' {
		b = cr.Byte()
	} else if b == '-' {
		sign = -1.0
		b = cr.Byte()
	}

	f := 0.0
	for '0' <= b && b <= '9' {
		f = 10*f + float64(b-'0')
		b = cr.Byte()
	}

	e := 0
	if b == '.' {
		b = cr.Byte()
		for '0' <= b && b <= '9' {
			f = 10*f + float64(b-'0')
			b = cr.Byte()
			e--
		}
	}

	if b == 'e' || b == 'E' {
		e += cr.Int()
	}

	cr.skipCRLF(b)

	return math.Copysign(f*math.Pow10(e), sign)
}

func (cr *Reader) Float64() float64 {
	w := cr.Word()
	f, err := strconv.ParseFloat(w, 64)
	if err != nil {
		panic(err)
	}

	//cr.skipCRLF(b)

	return f
}

func (cr *Reader) Slice() (result []byte) {
	if cr.index >= len(cr.buf) {
		cr.nextRead()
	}

	result = cr.buf[cr.index:]
	cr.index = len(cr.buf)
	return
}

func (cr *Reader) Line() string {
	if cr.index >= len(cr.buf) {
		cr.nextRead()
	}

	buf := cr.buf[cr.index:]
	i := bytes.IndexByte(buf, '\n')
	if i >= 0 {
		cr.index += i + 1
		s := string(buf[:i])
		return cutCR(s)
	}

	ss := []string{string(buf)}
	for {
		cr.nextRead()

		i := bytes.IndexByte(cr.buf, '\n')
		if i >= 0 {
			cr.index = i + 1
			ss = append(ss, string(cr.buf[:i]))
			break
		} else {
			ss = append(ss, string(cr.buf))
		}
	}
	s := strings.Join(ss, "")
	return cutCR(s)
}

func cutCR(s string) string {
	l := len(s)
	if l > 0 && s[l-1] == '\r' {
		s = s[:l-1]
	}
	return s
}

func (cr *Reader) Word() string {
	if cr.index >= len(cr.buf) {
		cr.nextRead()
	}

	buf := cr.buf[cr.index:]
	i := bytes.IndexAny(buf, "\n ")
	if i >= 0 {
		cr.index += i + 1
		s := string(buf[:i])
		return cutCR(s)
	}

	ss := []string{string(buf)}
	for {
		cr.nextRead()

		i := bytes.IndexAny(cr.buf, "\n ")
		if i >= 0 {
			ss = append(ss, string(cr.buf[:i]))
			cr.index += i + 1
			break
		} else {
			ss = append(ss, string(cr.buf))
		}
	}
	s := strings.Join(ss, "")
	return cutCR(s)
}
