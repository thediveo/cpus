// Copyright 2024 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cpus

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thediveo/faf"
)

// List is a list of CPU [from...to] ranges. CPU numbers are starting from zero.
type List [][2]uint

// String returns the CPU list in textual format, with the individual ranges
// “x-y” separated by “,” and single CPU ranges collapsed into “x”.
func (l List) String() string {
	var b strings.Builder
	for idx, cpurange := range l {
		if idx > 0 {
			b.WriteString(",")
		}
		if cpurange[0] == cpurange[1] {
			b.WriteString(fmt.Sprintf("%d", cpurange[0]))
			continue
		}
		b.WriteString(fmt.Sprintf("%d-%d", cpurange[0], cpurange[1]))
	}
	return b.String()
}

// NewList returns a new CPU List for the given textual list format.
func NewList(b []byte) (List, error) {
	bs := faf.NewBytestring(b)
	l := List{}
	for {
		// nothing more, we're at the end of text/line, so we're successfully
		// done.
		if bs.EOL() {
			return l, nil
		}
		// we now expect a CPU number and if there is nothing else following,
		// we're also done, adding the CPU number as a single CPU range to our
		// list.
		from, ok := bs.Uint64()
		if !ok {
			return nil, errors.New("expected unsigned integer number")
		}
		if bs.EOL() {
			return append(l, [2]uint{uint(from), uint(from)}), nil
		}
		// Either this is a from-to range or another range should follow...
		switch ch, _ := bs.Next(); ch {
		case '-':
			// a range, so get the end of the range and then add the range to
			// our list. If nothing else follows, then we're done.
			to, ok := bs.Uint64()
			if !ok {
				return nil, errors.New("expected unsigned integer number")
			}
			l = append(l, [2]uint{uint(from), uint(to)})
			if bs.EOL() {
				return l, nil
			}
			// another CPU number (or range) is expected to follow, separated by
			// ",", so we check for a necessary comma. Then we start over with
			// the next CPU number or range.
			ch, _ = bs.Next()
			if ch != ',' {
				return nil, errors.New("expected ','")
			}
		case ',':
			// a single CPU number, and more to follow; so add this single CPU
			// range and then rinse and repeat.
			l = append(l, [2]uint{uint(from), uint(from)})
		default:
			return nil, errors.New("expected '-' or ','")
		}
	}
}
