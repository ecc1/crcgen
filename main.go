package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

var (
	prog       = path.Base(os.Args[0])
	stdUsage   = flag.Usage
	goPackage  = os.Getenv("GOPACKAGE")
	crcSize    = flag.Int("size", 0, "CRC size in `bits` (8 or 16)")
	polynomial = flag.Int("poly", 0, "CRC `polynomial`")
	outputFile = flag.String("output", "", "output `filename` (default: crc<size>_table.go)")
)

// Usage prints the program's usage information.
func Usage() {
	stdUsage()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "See https://en.wikipedia.org/wiki/Cyclic_redundancy_check#Standards_and_common_use")
	fmt.Fprintln(os.Stderr, "for more information about CRC polynomials.")
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if goPackage == "" {
		fmt.Fprintf(os.Stderr, "%s: GOPACKAGE environment variable must be set\n\n", prog)
		Usage()
	}
	if *crcSize != 8 && *crcSize != 16 {
		fmt.Fprintf(os.Stderr, "%s: CRC size must be 8 or 16 bits\n\n", prog)
		Usage()
	}
	if *polynomial == 0 {
		fmt.Fprintf(os.Stderr, "%s: a nonzero polynomial must be specified\n\n", prog)
		Usage()
	}
	if *outputFile == "" {
		*outputFile = fmt.Sprintf("crc%d_table.go", *crcSize)
	}
	f := setup()
	switch *crcSize {
	case 8:
		genCRC8(f)
	case 16:
		genCRC16(f)
	}
}

func setup() io.WriteCloser {
	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("%s: %v\n", prog, err)
	}
	fmt.Fprintf(f, "// Generated by \"%s %s\": do not edit!\n\n", prog, strings.Join(os.Args[1:], " "))
	fmt.Fprintf(f, "package %s\n\n", goPackage)
	return f
}

// genCRC8 generates the lookup table for a CRC-8 calculation.
func genCRC8(f io.WriteCloser) {
	poly := uint8(*polynomial)
	fmt.Fprintf(f, "// Lookup table for CRC-8 calculation with polyomial 0x%02X.\n", poly)
	fmt.Fprintf(f, "var crc8Table = []uint8{\n")
	for i := 0; i < 256; i++ {
		res := uint8(i)
		for n := 0; n < 8; n++ {
			c := res & (1 << 7)
			res <<= 1
			if c != 0 {
				res ^= poly
			}
		}
		if i%8 == 0 {
			fmt.Fprintf(f, "\t")
		} else {
			fmt.Fprintf(f, " ")
		}
		fmt.Fprintf(f, "0x%02X,", res)
		if (i+1)%8 == 0 {
			fmt.Fprintf(f, "\n")
		}
	}
	fmt.Fprintf(f, "}\n")
	_ = f.Close()
}

// genCRC16 generates the lookup table for a CRC-16 calculation.
func genCRC16(f io.WriteCloser) {
	poly := uint16(*polynomial)
	fmt.Fprintf(f, "// Lookup table for CRC-16 calculation with polynomial 0x%04X.\n", poly)
	fmt.Fprintf(f, "var crc16Table = []uint16{\n")
	for i := 0; i < 256; i++ {
		res := uint16(0)
		b := uint16(i << 8)
		for n := 0; n < 8; n++ {
			c := (res ^ b) & (1 << 15)
			res <<= 1
			b <<= 1
			if c != 0 {
				res ^= poly
			}
		}
		if i%8 == 0 {
			fmt.Fprintf(f, "\t")
		} else {
			fmt.Fprintf(f, " ")
		}
		fmt.Fprintf(f, "0x%04X,", res)
		if (i+1)%8 == 0 {
			fmt.Fprintf(f, "\n")
		}
	}
	fmt.Fprintf(f, "}\n")
	_ = f.Close()
}
