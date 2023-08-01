/**
 * Copyright 2023 Google LLC
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

package utils_test

import (
	"errors"
	"fmt"
	"testing"
)

// inspectableT wraps testing.T, overriding testing behavior to make error cases retrievable.
type inspectableT struct {
	*testing.T
	err error
}

func (it *inspectableT) Error(args ...interface{}) {
	it.addError(args...)
}

func (it *inspectableT) Errorf(format string, args ...interface{}) {
	a := append([]interface{}{format}, args)
	it.addError(a)
}

func (it *inspectableT) Fatal(args ...interface{}) {
	it.addError(args...)
}

func (it *inspectableT) Fatalf(format string, args ...interface{}) {
	a := append([]interface{}{format}, args)
	it.addError(a)
}

func (it *inspectableT) addError(args ...interface{}) {
	s := fmt.Sprint(args...)
	it.err = errors.Join(it.err, errors.New(s))
}
