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

	"slices"

	"github.com/thediveo/faf"
)

// List is a list of CPU [from...to] ranges. CPU numbers are starting from zero.
type List [][2]uint

// String returns the CPU list in textual format, with the individual ranges
// “x-y” separated by “,” and single CPU ranges collapsed into “x” (instead of
// “x-x”).
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

// NewList returns a new CPU List for the given textual list format. If the text
// is malformed then an error is returned instead.
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

// Set returns the CPU Set corresponding with this list.
func (l List) Set() Set {
	if len(l) == 0 {
		return Set{}
	}
	// Do last range first to allocate only once.
	var s Set
	for i := range l {
		r := l[len(l)-i-1]
		s = s.AddRange(r[0], r[1])
	}
	return s
}

// IsOverlapping returns true if this List overlaps with another List.
//
// Both lists must be in canonical form where all ranges are ordered from lowest
// to highest and never overlap within the same list.
func (l List) IsOverlapping(another List) bool {
	// We assume canonical list form here, that is, all ranges within a list are
	// ordered from lowest to highest and never overlapping within a list.
	r2idx := 0
	for _, r1 := range l {
		for {
			// If we're exhausted our second range list to compare with, we're
			// done: there can't be any overlap.
			if r2idx >= len(another) {
				return false
			}
			// We're positively done if the current first range and the current
			// second range overlap.
			if r1[1] >= another[r2idx][0] && r1[0] <= another[r2idx][1] {
				return true
			}
			// When the current range from the second list now is beyond the
			// current range from the first list we need to advance to the next
			// range from that first list and then rinse and repeat.
			if another[r2idx][0] > r1[1] {
				break
			}
			// Process ranges from the second list while we've yet to catch up
			// to the current first list range.
			r2idx++
		}
	}
	return false
}

// Overlap returns the overlap of this List with another List as a new List. If
// the range lists are not overlapping, then an empty new List is returned.
func (l List) Overlap(another List) List {
	overlaps := List{}
	r2idx := 0
	for _, r1 := range l {
		for {
			// If we're exhausted our second range list to compare with, we're
			// done: there can't be any more overlap.
			if r2idx >= len(another) {
				return overlaps
			}
			// If we have overlap, then add the range where the lists overlap to
			// the result. In contrast to just detecting an overlap we then
			// carry on, as there might be more overlaps in the store for us.
			if r1[1] >= another[r2idx][0] && r1[0] <= another[r2idx][1] {
				from := max(r1[0], another[r2idx][0])
				to := min(r1[1], another[r2idx][1])
				overlaps = append(overlaps, [2]uint{from, to})
			}
			// Depending on whether the second ranges end lies beyond the first
			// ranges end we either need to move on to the next first range, or
			// next second range, respectively.
			if another[r2idx][1] > r1[1] {
				break
			}
			r2idx++
		}
	}
	return overlaps
}

// Remove the lowest CPU from the specified List, returning the CPU number
// together with a new List of remaining CPUs.
//
// The Remove operation is useful to pick individual and available (“online”)
// CPUs after first getting the List of CPU affinities for a task/process.
func (l List) Remove() (cpu uint, remaining List) {
	if len(l) == 0 {
		panic("cannot remove from empty List")
	}
	lowestRange := l[0]
	if lowestRange[0] < lowestRange[1] {
		// There will still be CPUs in the lowest range after we've removed the
		// CPU at the beginning of the range...
		cpu = lowestRange[0]
		return cpu, append(List{[2]uint{cpu + 1, lowestRange[1]}}, l[1:]...)
	}
	// We've exhausted the lowest range after we've removed the last CPU from
	// that range, so we return the remaining ranges, throwing away the now
	// empty lowest range...
	return lowestRange[0], slices.Clone(l[1:])
}
