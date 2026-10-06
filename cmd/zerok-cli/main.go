// Copyright (c) 2026 Nrupal Akolkar
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"

	"github.com/zerok-vault/zerok-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}