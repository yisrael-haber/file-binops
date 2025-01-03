package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type ParsedArgs struct {
	ifn     string
	ofn     string
	op      string
	comp    int
	force   bool
	bufsize int
	verbose bool
}

const XOR = "XOR"
const AND = "AND"
const OR = "OR"

func parseArgs(args []string) (ParsedArgs, error) {
	if (len(args) < 2) || (args[0] == "--help") {
		return ParsedArgs{}, errors.New("")
	}

	available_binops := []string{XOR, AND, OR}

	if_name := UNDEFINED_FILE_NAME
	of_name := UNDEFINED_FILE_NAME
	binop := UNDEFINED_BINOP
	comp := -1
	force := false
	bufsize := BIG_BUFF_SIZE
	verbose := true

	for idx := range args {
		if args[idx] == "--if" {
			if_name = args[idx+1]
		}

		if args[idx] == "--of" {
			of_name = args[idx+1]
		}

		if args[idx] == "--op" {
			binop = args[idx+1]
		}

		if args[idx] == "--comp" {
			comp_cand, err := strconv.Atoi(args[idx+1])
			if err != nil {
				return ParsedArgs{}, err
			}
			comp = comp_cand
		}

		if args[idx] == "--force" {
			force = true
		}

		if args[idx] == "--buf" {
			bufsize_cand, err := strconv.Atoi(args[idx+1])
			if err != nil {
				return ParsedArgs{}, err
			}
			bufsize = bufsize_cand
		}

		if args[idx] == "--quiet" {
			verbose = false
		}
	}

	if if_name == UNDEFINED_FILE_NAME {
		return ParsedArgs{}, errors.New("ERROR: input file name not provided. Make sure that the format you use is --input-file IF_NAME")
	}

	if binop == UNDEFINED_BINOP {
		return ParsedArgs{}, errors.New("ERROR: Binop not provided. Make sure that the format you use is --binop BINOP")
	}

	if !slices.Contains(available_binops, binop) {
		return ParsedArgs{}, fmt.Errorf("ERROR: Unsupported binop %s. Currently supported binops: %v", binop, available_binops)
	}

	if comp < 0 || comp > 255 {
		return ParsedArgs{}, fmt.Errorf("ERROR: comparator %d cannot be used in the program. Value must represent a byte value", comp)
	}

	if bufsize < 0 {
		return ParsedArgs{}, fmt.Errorf("ERROR: buffer size option must be positive, got %d", bufsize)
	}

	if of_name == UNDEFINED_FILE_NAME {
		dir, file := filepath.Split(if_name)
		new_file_name := strings.Join([]string{"modified", file}, "_")
		of_name = filepath.Join(dir, new_file_name)
	}

	return ParsedArgs{
		ifn:     if_name,
		ofn:     of_name,
		op:      binop,
		comp:    comp,
		force:   force,
		bufsize: bufsize,
		verbose: verbose,
	}, nil
}

func (pa ParsedArgs) apply_vectorized_function() error {
	big_buf := make([]byte, pa.bufsize)

	f, err := os.Open(pa.ifn)
	check(err)
	defer f.Close()

	ofalgs := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	of, err := os.OpenFile(pa.ofn, ofalgs, 0644)
	check(err)
	defer f.Close()

	ifile_info, err := os.Stat(pa.ifn)
	if err != nil {
		return err
	}

	size_in_bytes := ifile_info.Size()

	pa.print_if_verbose("Processing input file with size: %s\n", byte_count_iec(size_in_bytes))

	total_read_bytes := int64(0)

	for {

		ba, err := f.Read(big_buf)
		if check_eof_error(err) {
			break
		}
		total_read_bytes += int64(ba)
		total_read_bytes_readable := byte_count_iec(total_read_bytes)
		percent := 100 * (float64(total_read_bytes) / float64(size_in_bytes))

		pa.print_if_verbose("Read %d bytes (%s), out of %d (%f percent)...", total_read_bytes, total_read_bytes_readable, size_in_bytes, percent)

		switch pa.op {
		case XOR:
			for i := range ba {
				big_buf[i] ^= byte(pa.comp)
			}
		case AND:
			for i := range ba {
				big_buf[i] &= byte(pa.comp)
			}
		case OR:
			for i := range ba {
				big_buf[i] |= byte(pa.comp)
			}
		}

		n, err := of.Write(big_buf[:ba])
		check(err)

		if n < ba {
			panic("bytes written to output file were of different length than bytes read from input file")
		}

		pa.print_if_verbose("Done !\n")
	}

	return nil
}

func (pa ParsedArgs) print_if_verbose(str string, args ...any) {
	if pa.verbose {
		if len(args) > 0 {
			fmt.Printf(str, args...)
			return
		}
		fmt.Print(str)
	}
}
