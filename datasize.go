package datasize

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Size uint64

const (
	Kilobyte Size = 1000 * Byte
	Megabyte      = 1000 * Kilobyte
	Gigabyte      = 1000 * Megabyte
	Terabyte      = 1000 * Gigabyte
	Petabyte      = 1000 * Terabyte
)

const (
	Byte     Size = 1
	Kibibyte      = 1024 * Byte
	Mebibyte      = 1024 * Kibibyte
	Gibibyte      = 1024 * Mebibyte
	Tebibyte      = 1024 * Gibibyte
	Pebibyte      = 1024 * Tebibyte
)

var sizeRegex = regexp.MustCompile(`([0-9]*)(\.[0-9]*)?([a-z]+)`)

var units = []Size{
	Pebibyte, Petabyte,
	Tebibyte, Terabyte,
	Gibibyte, Gigabyte,
	Mebibyte, Megabyte,
	Kibibyte, Kilobyte,
}

func Parse(s string) (Size, error) {
	if s == "" {
		return 0, errors.New("datasize: invalid Size: empty")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return Size(f), nil
	}
	ss := sizeRegex.FindStringSubmatch(strings.ToLower(s))
	if len(ss) == 0 {
		return 0, fmt.Errorf("datasize: invalid Size format: %q", s)
	}
	f, err = strconv.ParseFloat(ss[1]+ss[2], 64)
	if err != nil {
		return 0, err
	}
	sz, err := suffixSize(ss[3])
	if err != nil {
		return 0, err
	}
	return Size(f * float64(sz)), nil
}

func suffixSize(suffix string) (Size, error) {
	switch suffix {
	case "b":
		return Byte, nil
	case "kb":
		return Kilobyte, nil
	case "mb":
		return Megabyte, nil
	case "gb":
		return Gigabyte, nil
	case "tb":
		return Terabyte, nil
	case "pb":
		return Petabyte, nil
	case "kib":
		return Kibibyte, nil
	case "mib":
		return Mebibyte, nil
	case "gib":
		return Gibibyte, nil
	case "tib":
		return Tebibyte, nil
	case "pib":
		return Pebibyte, nil
	default:
		return 0, fmt.Errorf("datasize: invalid Size unit suffix: %q", suffix)
	}
}

func sizeSuffix(unit Size) string {
	switch unit {
	default:
		return "B"
	case Petabyte:
		return "PB"
	case Pebibyte:
		return "PiB"
	case Terabyte:
		return "TB"
	case Tebibyte:
		return "TiB"
	case Gigabyte:
		return "GB"
	case Gibibyte:
		return "GiB"
	case Megabyte:
		return "MB"
	case Mebibyte:
		return "MiB"
	case Kilobyte:
		return "kB"
	case Kibibyte:
		return "KiB"
	}
}

func (s Size) Floor() Size {
	for _, unit := range units {
		if s >= unit {
			return (s / unit) * unit
		}
	}
	return s
}

func (s Size) Round() Size {
	for _, unit := range units {
		if s >= unit {
			return Size(math.Round(float64(s)/float64(unit))) * unit
		}
	}
	return s
}

func (s Size) String() string {
	switch {
	case s == 0:
		return "0B"
	case s%Petabyte == 0:
		return format(s.Petabytes(), "PB")
	case s >= Pebibyte:
		return format(s.Pebibytes(), "PiB")
	case s%Terabyte == 0:
		return format(s.Terabytes(), "TB")
	case s >= Tebibyte:
		return format(s.Tebibytes(), "TiB")
	case s%Gigabyte == 0:
		return format(s.Gigabytes(), "GB")
	case s >= Gibibyte:
		return format(s.Gibibytes(), "GiB")
	case s%Megabyte == 0:
		return format(s.Megabytes(), "MB")
	case s >= Mebibyte:
		return format(s.Mebibytes(), "MiB")
	case s%Kilobyte == 0:
		return format(s.Kilobytes(), "kB")
	case s >= Kibibyte:
		return format(s.Kibibytes(), "KiB")
	default:
		return fmt.Sprintf("%dB", s)
	}
}

func format(size float64, suffix string) string {
	if math.Floor(size) == size {
		return fmt.Sprintf("%d%s", int64(size), suffix)
	}
	return fmt.Sprintf("%.2f%s", size, suffix)
}

func (s Size) Bytes() uint64 {
	return uint64(s)
}

func (s Size) Kilobytes() float64 {
	return float64(s) / float64(Kilobyte)
}

func (s Size) Megabytes() float64 {
	return float64(s) / float64(Megabyte)
}

func (s Size) Gigabytes() float64 {
	return float64(s) / float64(Gigabyte)
}

func (s Size) Terabytes() float64 {
	return float64(s) / float64(Terabyte)
}

func (s Size) Petabytes() float64 {
	return float64(s) / float64(Petabyte)
}

func (s Size) Kibibytes() float64 {
	return float64(s) / float64(Kibibyte)
}

func (s Size) Mebibytes() float64 {
	return float64(s) / float64(Mebibyte)
}

func (s Size) Gibibytes() float64 {
	return float64(s) / float64(Gibibyte)
}

func (s Size) Tebibytes() float64 {
	return float64(s) / float64(Tebibyte)
}

func (s Size) Pebibytes() float64 {
	return float64(s) / float64(Pebibyte)
}

type sizeFlag struct {
	*Size
}

func Flag(name, value, description string) *Size {
	sz, err := Parse(value)
	if err != nil {
		panic(fmt.Sprintf("Invalid Size value for flag --%q: %q", name, value))
	}
	return FlagVar(flag.CommandLine, &sz, name, sz, description)
}

func FlagVar(fs *flag.FlagSet, s *Size, name string, value Size, description string) *Size {
	*s = value
	f := &sizeFlag{s}
	fs.Var(f, name, description)
	return f.Size
}

func (f *sizeFlag) Get() any {
	return *f.Size
}

func (f *sizeFlag) Set(s string) error {
	sz, err := Parse(s)
	if err != nil {
		return err
	}
	*f.Size = sz
	return nil
}
