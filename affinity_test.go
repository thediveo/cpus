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
	"os"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("task/thread CPU affinities", func() {

	It("gets this process's CPU affinity list, consistent with /proc/self/status data", func() {
		Expect(elementBytesSize).To(Equal(uint64(64 /* bits in uint64 */ / 8 /* bits/byte*/)))
		cpulist := Successful(Affinity(os.Getpid())).List()
		Expect(cpulist).NotTo(BeEmpty())
		Expect(systemSetSize.Load()).NotTo(BeZero())

		var prefix = []byte("Cpus_allowed_list:\t")
		var allowedList List
		for line := range bytes.Lines(Successful(os.ReadFile("/proc/self/status"))) {
			if !bytes.HasPrefix(line, prefix) {
				continue
			}
			allowedList = Successful(NewList(line[len(prefix) : len(line)-1]))
		}
		Expect(cpulist).To(Equal(allowedList))
	})

	It("changes this process's CPU affinity", func() {
		runtime.LockOSThread() // don't unlock, throw away the tainted task

		affs := Successful(Affinity(0))
		oneonly, _ := affs.List().Remove()
		Expect(Set{}.AddRange(oneonly, oneonly).PinTask(0)).To(Succeed())

		reducedaffs := Successful(Affinity(0)).List()
		Expect(reducedaffs).To(Equal(List{[2]uint{oneonly, oneonly}}))

		Expect(affs.PinTask(0)).To(Succeed())
	})

	It("cannot set empty affinities", func() {
		Expect(SetAffinity(0, Set{})).NotTo(Succeed())
		Expect(SetAffinity(0, Set{0})).NotTo(Succeed())
	})

})
