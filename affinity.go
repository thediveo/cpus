// Copyright 2026 Harald Albrecht.
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
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Affinity returns the affinity CPUList (list of CPU ranges) of the
// process with the passed PID. Otherwise, it returns an error. If pid is zero,
// then the affinity CPU list of the calling thread is returned (make sure to
// have the OS-level thread locked to the calling go routine in this case).
//
// Notes:
//   - we don't use [unix.SchedGetaffinity] as this is tied to the fixed size
//     [unix.CPUSet] type; instead, we dynamically figure out the size needed
//     and cache the size internally.
//   - retrieving the affinity CPU mask and then speed-running it to
//     generate the range list is roughly two orders of magnitude faster than
//     fetching “/proc/$PID/status” and looking for the “Cpus_allowed_list”,
//     because generating the broad status procfs file is expensive.
func Affinity(tid int) (Set, error) {
	var set Set

	setlenStart := systemSetSize.Load()
	setlen := setlenStart
	for {
		set = make(Set, setlen)
		// see also:
		// https://man7.org/linux/man-pages/man2/sched_setaffinity.2.html; we
		// use RawSyscall here instead of Syscall as we know that
		// SYS_SCHED_GETAFFINITY does not block, following Go's stdlib
		// implementation.
		_, _, e := unix.RawSyscall(unix.SYS_SCHED_GETAFFINITY,
			uintptr(tid),
			uintptr(setlen*elementBytesSize),
			uintptr(unsafe.Pointer(&set[0])))
		if e != 0 {
			if e == unix.EINVAL {
				setlen *= 2
				continue
			}
			return nil, e
		}
		// Set the new size; if this fails because another go routine already
		// upped the set size, retry until we either notice that we're smaller
		// than what was set as the new set size, or we succeed in setting the
		// size.
		for !systemSetSize.CompareAndSwap(setlenStart, setlen) {
			setlenStart = systemSetSize.Load()
			if setlenStart > setlen {
				break
			}
		}
		break
	}
	return set, nil
}

// SetAffinity sets the CPU affinities for the specified task/process, returning
// nil on success. Otherwise, it returns an error. It is an error trying to
// remove any CPU affinity by specifying an effectively empty Set.
//
// See also the functional equivalent [Set.PinTask].
func SetAffinity(tid int, cpus Set) error {
	if len(cpus) == 0 {
		return syscall.EINVAL
	}
	_, _, e := unix.RawSyscall(unix.SYS_SCHED_SETAFFINITY,
		uintptr(tid), uintptr(uint64(len(cpus))*elementBytesSize), uintptr(unsafe.Pointer(&cpus[0])))
	if e != 0 {
		return e
	}
	return nil
}
