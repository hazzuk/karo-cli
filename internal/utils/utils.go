// SPDX-FileCopyrightText: © 2026 hazzuk
//
// SPDX-License-Identifier: AGPL-3.0-only

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
