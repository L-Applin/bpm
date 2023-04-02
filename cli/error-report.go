// Copyright (c) 2023 Olivier Lepage-Applin. All rights reserved.

package main

import (
	"bpm/log"
)

func ReportError(err error) error {
	log.ErrorE(err)
	return err
}
