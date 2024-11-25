// Copyright (c) HashiCorp, Inc.

package cronitor

import "errors"

var (
	ErrFailedGetMonitor    = errors.New("failed to get monitor details")
	ErrFailedCreateMonitor = errors.New("failed to create monitor")
	ErrFailedDeleteMonitor = errors.New("failed to delete monitor")
)
