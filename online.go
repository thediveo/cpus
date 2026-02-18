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
	"bytes"
	"path/filepath"

	"github.com/thediveo/faf"
)

// Online returns the List of logical CPUs that are currently online. In case of
// any errors when the online CPUs cannot be determined, it returns an empty CPU
// list; but it never returns a nil List.
//
// For background information, see also [How CPU topology info is exported via
// sysfs] in the Linux kernel documentation.
//
// [How CPU topology info is exported via sysfs]:
// https://www.kernel.org/doc/html/latest/admin-guide/cputopology.html#
func Online() List {
	return online("/")
}

func online(sysfsroot string) List {
	contents, ok := faf.ReadFile(filepath.Join(sysfsroot, "sys/devices/system/cpu/online"), nil)
	if !ok {
		return List{}
	}
	cpus, err := NewList(bytes.TrimSuffix(contents, []byte{'\n'}))
	if err != nil {
		return List{}
	}
	return cpus
}
