/**
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	terraTesting "github.com/gruntwork-io/terratest/modules/testing"
)

// GetLoggerFromT returns a logger based on test verbosity
func GetLoggerFromT() *logger.Logger {
	if testing.Verbose() {
		return logger.Default
	}
	return logger.Discard
}

// TestFileLogger is a logger that writes to disk instead of stdout.
// This is useful when you want to redirect verbose logs of long running tests to disk.
type TestFileLogger struct {
	pth string
	w   io.WriteCloser
}

// NewTestFileLogger returns a TestFileLogger logger that can be used with the WithLogger option.
func NewTestFileLogger(t *testing.T, pth string) (*logger.Logger, func(t *testing.T)) {
	f, err := os.OpenFile(pth, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("error opening file %s: %v", pth, err)
	}
	fl := TestFileLogger{
		pth: pth,
		w:   f,
	}
	return logger.New(fl), fl.Close
}

func (fl TestFileLogger) Logf(t terraTesting.TestingT, format string, args ...interface{}) {
	logger.DoLog(t, 3, fl.w, fmt.Sprintf(format, args...))
}

func (fl TestFileLogger) Close(t *testing.T) {
	if err := fl.w.Close(); err != nil {
		t.Fatalf("error closing file logger %s: %v", fl.pth, err)
	}
}
