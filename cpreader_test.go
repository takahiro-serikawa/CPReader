package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestInts(t *testing.T) {
	f, err := os.Open("test-ints.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cr := NewReader(f)

	ints := []int{0,
		9223372036854775807,  // math.MaxInt64
		-9223372036854775808, // math.MinInt64
		1_000_000_000_000_000_000,
		999_999_999_999_999_999,
		+999_999_999_999_999_999,
		-999_999_999_999_999_999}

	for _, v0 := range ints {
		if v := cr.Int(); v != v0 {
			t.Errorf("%v != %v\n", v, v0)
		}
	}

	ints32 := []int32{
		2147483647,  // math.MaxInt32
		-2147483648, // math.MinInt32
		1_000_000_000,
		999_999_999,
		+999_999_999,
		-999_999_999}

	for _, v0 := range ints32 {
		if v := cr.Int32(); v != v0 {
			t.Errorf("%v != %v\n", v, v0)
		}
	}

	uints64 := []uint64{18446744073709551615} // math.MaxUint64
	for _, v0 := range uints64 {
		if v := cr.Uint64(); v != v0 {
			t.Errorf("%v != %v\n", v, v0)
		}
	}

	uints32 := []uint32{4294967295} // math.MaxUint32
	for _, v0 := range uints32 {
		if v := cr.Uint32(); v != v0 {
			t.Errorf("%v != %v\n", v, v0)
		}
	}

	//cr.Int() FAIL
}

func TestFloats(t *testing.T) {
	f, err := os.Open("test-floats.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cr := NewReader(f)

	floats := []float64{
		0,
		0.0,
		0.,
		.0,
		+0,
		-0,
		00,
		01,
		1.79769e+308,
		+1.79769e+308,
		-1.79769e+308,
		4.94066e-324,
		+4.94066e-324,
		-4.94066e-324,
		999999999,
		999.999999,
		123.456789,
		.456789,
	}

	for i, v0 := range floats {
		if v := cr.Float64(); v != v0 {
			t.Errorf("#%v: %v != %v\n", i, v, v0)
		}
	}
}

func MakeTestDataPerm(dummy_txt string, n, offset int) {
	rand.Seed(time.Now().UnixNano())

	o, err := os.Create(dummy_txt)
	if err != nil {
		panic(err)
	}
	defer o.Close()
	wr := bufio.NewWriter(o)
	defer wr.Flush()

	fmt.Fprintln(wr, n)
	aa := rand.Perm(n)
	for i, a := range aa {
		if i > 0 {
			fmt.Fprint(wr, " ")
		}
		fmt.Fprint(wr, a+offset)
	}
	fmt.Fprintln(wr)
}

func MakeTestHW(dummy_txt string, h, w int) {
	rand.Seed(time.Now().UnixNano())

	o, err := os.Create(dummy_txt)
	if err != nil {
		panic(err)
	}
	defer o.Close()
	wr := bufio.NewWriter(o)
	defer wr.Flush()

	fmt.Fprintln(wr, h, w)
	for i := 0; i < h; i++ {
		bb := make([]byte, w)
		for j := 0; j < w; j++ {
			bb[j] = "abcdefghijklmnopqrstuvwxyz"[rand.Int()%26]
		}
		fmt.Fprintln(wr, string(bb))
	}
}

func TestMain(m *testing.M) {
	// before all...
	err := os.MkdirAll(TESTDATA, 0777)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(TESTDATA + "perm_1e7I.txt"); err != nil {
		MakeTestDataPerm(TESTDATA+"perm_1e7I.txt", 10_000_000, 0)
	}
	if _, err := os.Stat(TESTDATA + "perm_2e7I.txt"); err != nil {
		MakeTestDataPerm(TESTDATA+"perm_2e7I.txt", 20_000_001, -10_000_000)
	}

	if _, err := os.Stat(TESTDATA + "HW10Kx10K.txt"); err != nil {
		MakeTestHW(TESTDATA+"HW10Kx10K.txt", 10_000, 10_000)
	}
	if _, err := os.Stat(TESTDATA + "HW1x100M.txt"); err != nil {
		MakeTestHW(TESTDATA+"HW1x100M.txt", 1, 100_000_000)
	}
	if _, err := os.Stat(TESTDATA + "HW100Mx1.txt"); err != nil {
		MakeTestHW(TESTDATA+"HW100Mx1.txt", 100_000_000, 1)
	}

	code := m.Run()

	// after all...

	os.Exit(code)
}

// cr.Int(); 順列 0 .. +10,000,000 まで
func TestPerm_1e7I(t *testing.T) {
	n, sum := LoadPerm(t, TESTDATA+"perm_1e7I.txt")

	//N := 10_000_000
	if (n*(n-1))/2 != sum {
		t.Error("sum miss matched")
	}
}

// cr.Int(); 順列 -10,000,000 .. +10,000,000 まで
func TestPerm_2e7I(t *testing.T) {
	n, sum := LoadPerm(t, TESTDATA+"perm_2e7I.txt")
	_ = n
	if sum != 0 {
		t.Error("sum miss matched")
	}
}

func LoadPerm(t *testing.T, dummy_txt string) (n, sum int) {
	f, err := os.Open(dummy_txt)
	if err != nil {
		panic(err)
	}

	time0 := time.Now()
	defer func() { fmt.Fprintf(os.Stderr, "%v %v\n", filepath.Base(dummy_txt), time.Since(time0)) }()

	cr := NewReader(f /*os.Stdin*/)
	N := cr.Int()
	sum = 0
	for i := 0; i < N; i++ {
		a := cr.Int()
		sum += a
	}

	return N, sum
}

// 性能比較用 bufio.Scanner sc.Buffer(...)
func BenchmarkBufio_Perm_2e7I(b *testing.B) {
	b.ResetTimer()

	f, err := os.Open(TESTDATA + "perm_2e7I.txt")
	if err != nil {
		panic(err)
	}

	time0 := time.Now()
	defer func() { fmt.Fprintf(os.Stderr, "bufio.Scanner; perm_2e7I %v\n", time.Since(time0)) }()

	sc := bufio.NewScanner(f /*os.Stdin*/)
	nbytes := 11*2_000_000_001 + 2
	sc.Buffer(make([]byte, nbytes), nbytes)
	sc.Scan()
	N, _ := strconv.Atoi(sc.Text())

	sc.Scan()
	ss := strings.Split(sc.Text(), " ")
	sum := 0
	for i := 0; i < N; i++ {
		a, _ := strconv.Atoi(ss[i])
		sum += a
	}

	if sum != 0 {
		b.Error("sum miss matched")
	}
}

// 性能比較用 bufio.Scanner, ScanWords版
func BenchmarkBufio_Perm_2e7I2(b *testing.B) {
	b.ResetTimer()

	f, err := os.Open(TESTDATA + "perm_2e7I.txt")
	if err != nil {
		panic(err)
	}

	time0 := time.Now()
	defer func() { fmt.Fprintf(os.Stderr, "bufio.Scanner; perm_2e7I %v\n", time.Since(time0)) }()

	sc := bufio.NewScanner(f /*os.Stdin*/)
	sc.Split(bufio.ScanWords)
	sc.Scan()
	N, _ := strconv.Atoi(sc.Text())

	sum := 0
	for i := 0; i < N; i++ {
		sc.Scan()
		a, _ := strconv.Atoi(sc.Text())
		sum += a
	}

	if sum != 0 {
		b.Error("sum miss matched")
	}
}

func Benchmark_HW10Kx10K(b *testing.B) {
	loadHW(b, TESTDATA+"HW10Kx10K.txt")
}

func Benchmark_HW1x1000000(b *testing.B) {
	loadHW(b, TESTDATA+"HW1x100M.txt")
}

func Benchmark_HW1000000x1(b *testing.B) {
	loadHW(b, TESTDATA+"HW100Mx1.txt")
}

func loadHW(b *testing.B, dummy_txt string) {
	b.ResetTimer()

	f, err := os.Open(dummy_txt)
	if err != nil {
		panic(err)
	}

	time0 := time.Now()
	defer func() { fmt.Fprintf(os.Stderr, "%v %v\n", filepath.Base(dummy_txt), time.Since(time0)) }()

	cr := NewReader(f)
	h, w := cr.Int(), cr.Int()
	for i := 0; i < h; i++ {
		s := cr.Line()
		if len(s) != w {
			b.Errorf("length miss matched; %v len(s) != %v (w)\n", len(s), w)
		}
	}
}

func Benchmark_ScHW10Kx10K(b *testing.B) {
	ScloadHW(b, TESTDATA+"HW10Kx10K.txt")
}

func Benchmark_ScHW1x100M(b *testing.B) {
	ScloadHW(b, TESTDATA+"HW1x100M.txt")
}

func Benchmark_ScHW100Mx1(b *testing.B) {
	ScloadHW(b, TESTDATA+"HW100Mx1.txt")
}

func ScloadHW(b *testing.B, dummy_txt string) {
	b.ResetTimer()

	f, err := os.Open(dummy_txt)
	if err != nil {
		panic(err)
	}

	time0 := time.Now()
	defer func() { fmt.Fprintf(os.Stderr, "bufio.Scanner; %v %v\n", filepath.Base(dummy_txt), time.Since(time0)) }()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 100_000_000+2), 100_000_000+2)
	sc.Scan()
	ss := strings.Split(sc.Text(), " ")
	h, _ := strconv.Atoi(ss[0])
	w, _ := strconv.Atoi(ss[1])
	for i := 0; i < h; i++ {
		sc.Scan()
		s := sc.Text()
		if len(s) != w {
			b.Errorf("length miss matched; %v len(s) != %v (w)\n", len(s), w)
		}
	}
}

// 性能比較; Word()切り出し + ParseInt()
func TestPerm_2e7I_P(t *testing.T) {
	n, sum := LoadPermP(t, TESTDATA+"perm_2e7I.txt")
	_ = n
	if sum != 0 {
		t.Error("sum miss matched")
	}
}

func LoadPermP(t *testing.T, dummy_txt string) (n, sum int) {
	f, err := os.Open(dummy_txt)
	if err != nil {
		panic(err)
	}

	time0 := time.Now()
	defer func() { fmt.Fprintf(os.Stderr, "%v %v\n", filepath.Base(dummy_txt), time.Since(time0)) }()

	cr := NewReader(f /*os.Stdin*/)
	N := cr.IntP()
	sum = 0
	for i := 0; i < N; i++ {
		a := cr.IntP()
		sum += a
	}

	return N, sum
}

func (cr *Reader) IntP() int {
	w := cr.Word()
	i, err := strconv.ParseInt(w, 10, 64)
	if err != nil {
		panic(err)
	}

	//cr.skipCRLF(b)

	return int(i)
}
