package main

import (
	"fmt"
	"io"
	"os"
)

func display_help() {
	fmt.Println("USAGE: file_binops.exe --if INPUT_FILE_NAME --op BINOP --comp VAL [OPTION]")
	fmt.Println("Available Options:")
	fmt.Println("\t--help            Present help for command.")
	fmt.Println("\t--of              Provide name for the output file to write result to. by default, \"modified_\" is prepended to name of file.")
	fmt.Printf("\t--buf             Provide buffer size, in bytes, to use to read to and write from to apply the operator. By default it is %d.\n", BIG_BUFF_SIZE)
	fmt.Println("\t--quiet           Turns off printing info during processing. By default it is printed.")
	fmt.Print("\t--force           Forces overwriting of file if it exists before running. By default it is off.\n\n")
}

func byte_count_iec(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func file_exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func check_and_apply_force(pa ParsedArgs) error {
	if file_exists(pa.ofn) && !pa.force {
		return fmt.Errorf("file %s already exists, to override it use the --force flag", pa.ofn)
	}

	if file_exists(pa.ofn) && pa.force {
		err := os.Remove(pa.ofn)
		return err
	}

	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check_eof_error(e error) bool {
	if e != nil {
		if e == io.EOF {
			return true
		}
		panic(e)
	}
	return false
}
