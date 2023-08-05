package datasize

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		str string
		sz  Size
	}{
		{"0", 0},
		{"0tb", 0},
		{"100", 100 * Byte},
		{"1024", 1024 * Byte},
		{"1Gb", Gigabyte},
		{"1KiB", Kibibyte},
		{"1Pb", Petabyte},
		{"1PiB", Pebibyte},
		{"1b", Byte},
		{"1gib", Gibibyte},
		{"1kb", Kilobyte},
		{"1mb", Megabyte},
		{"1mib", Mebibyte},
		{"1tb", Terabyte},
		{"1tib", Tebibyte},
		{"510.85MB", Size(510.85 * float64(Megabyte))},
		{"5TB", 5 * Terabyte},
	}
	for _, test := range tests {
		found, err := Parse(test.str)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found != test.sz {
			t.Errorf("Parse(%q): expected: %q, found: %q", test.str, test.sz, found)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		str string
		sz  Size
	}{
		{"0B", 0},
		{"1B", Byte},
		{"1GB", Gigabyte},
		{"1GiB", Gibibyte},
		{"1KiB", Kibibyte},
		{"1MB", Megabyte},
		{"1MiB", Mebibyte},
		{"1PB", Petabyte},
		{"1PiB", Pebibyte},
		{"1TB", Terabyte},
		{"1TiB", Tebibyte},
		{"1kB", Kilobyte},
		{"2PB", 2 * Petabyte},
		{"407.50KiB", Size(407.5 * float64(Kibibyte))},
		{"407B", 407 * Byte},
	}
	for _, test := range tests {
		if found := test.sz.String(); found != test.str {
			t.Errorf("%d.String(): expected: %q, found: %q", test.sz, test.str, found)
		}
	}
}

func mustParse(s string) Size {
	sz, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return sz
}

func TestUnits(t *testing.T) {
	for i := 0; i < len(units); i++ {
		if expected := sizeSuffix(units[i]); !strings.HasSuffix(units[i].String(), expected) {
			t.Errorf("expected suffix: %q, found: %q", expected, units[i])
		}
		if i != len(units)-1 && units[i] <= units[i+1] {
			t.Errorf("%s >= %s", units[i], units[i+1])
		}
	}
}

func TestFloor(t *testing.T) {
	type test struct {
		sz       Size
		expected Size
	}
	tests := []test{
		{0, 0},
		{10, 10},
		{mustParse("1.013MB"), Megabyte},
		{mustParse("1.014GB"), Gigabyte},
		{mustParse("1.014KB"), Kilobyte},
		{mustParse("1.114PB"), Petabyte},
		{mustParse("1.12KiB"), Kibibyte},
		{mustParse("1.55MiB"), Mebibyte},
		{mustParse("1.608PiB"), Pebibyte},
		{mustParse("1.99GiB"), Gibibyte},
	}
	for _, unit := range units {
		tests = append(tests, test{unit, unit})
	}
	for _, test := range tests {
		found := test.sz.Floor()
		if found != test.expected {
			t.Errorf("%d.Floor(): expected: %q, found: %q", test.sz, test.expected, found)
		}
		if strings.Contains(found.String(), ".") {
			t.Errorf("%d.Floor() == %s: contains decimal", test.sz, found)
		}
	}
}

func TestRound(t *testing.T) {
	type test struct {
		sz       Size
		expected Size
	}
	tests := []test{
		{0, 0},
		{mustParse("1.23KiB"), Kibibyte},
		{mustParse("1.489PiB"), Pebibyte},
		{mustParse("1.4999MiB"), Mebibyte},
		{mustParse("1.8PiB"), 2 * Pebibyte},
		{mustParse("2.51GiB"), 3 * Gibibyte},
		{mustParse("2.9KiB"), 3 * Kibibyte},
		{mustParse("4.27GiB"), 4 * Gibibyte},
		{mustParse("4.87GiB"), 5 * Gibibyte},
	}
	for _, unit := range units {
		tests = append(tests, test{unit, unit})
	}
	for _, test := range tests {
		found := test.sz.Round()
		if found != test.expected {
			t.Errorf("%d.Round(): expected: %q, found: %q", test.sz, test.expected, found)
		}
		if strings.Contains(found.String(), ".") {
			t.Errorf("%d.Round() == %s: contains decimal", test.sz, found)
		}
	}
}

func TestRoundString(t *testing.T) {
	tests := []Size{0, 1023, mustParse("2.24KiB"), mustParse("1.52GiB")}
	tests = append(tests, units...)
	for _, size := range tests {
		str := size.String()
		found, err := Parse(str)
		if err != nil {
			t.Errorf("Parse(%q.String()): error: %v", size, err)
		}
		if found != size {
			t.Errorf("Parse(Size(%d).String()): expected: %q, found: %q", size, size, found)
		}
	}
}
