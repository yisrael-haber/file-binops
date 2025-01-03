package main

import (
	"fmt"
	"os"
)

const UNDEFINED_FILE_NAME = "__UNDEFINED_FILE_NAME__"
const UNDEFINED_BINOP = "__UNDEFINED_BINOP__"

const BIG_BUFF_SIZE = 1000 * 1000 * 70

func main() {
	parsed_args, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		display_help()
		return
	}

	e := check_and_apply_force(parsed_args)
	check(e)

	parsed_args.apply_vectorized_function()
}
