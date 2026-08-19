package utils

import (
	"fmt"
	"os"
)

func Check(err error) {

	if err != nil {
		fmt.Println("unexpected error: ", err)
		os.Exit(1)
	}

}
