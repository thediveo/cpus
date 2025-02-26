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
	"bytes"
	"os"
	"runtime"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/ginkgo/v2/dsl/table"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("cpu sets", func() {

	DescribeTable("parsing",
		func(set Set, expected List) {
			Expect(set.List()).To(Equal(expected))
		},
		Entry("nil set", nil, List{}),
		Entry("all-zeros set", Set{0}, List{}),
		Entry("all-zeros set", Set{0, 0}, List{}),

		// all in first word
		Entry("single cpu #0", Set{1 << 0, 0}, List{{0, 0}}),
		Entry("single cpu #1", Set{1 << 1}, List{{1, 1}}),
		Entry("single cpu #63", Set{1 << 63}, List{{63, 63}}),
		Entry("single cpu #63, none else", Set{1 << 63, 0, 0}, List{{63, 63}}),
		Entry("cpus #1-3", Set{0xe, 0}, List{{1, 3}}),

		// skip first zero words
		Entry("single cpu #64", Set{0, 1 << 0}, List{{64, 64}}),

		// multiple cpu ranges in same word
		Entry("cpu #1-2, #62", Set{1<<62 | 1<<2 | 1<<1}, List{{1, 2}, {62, 62}}),

		// range across boundaries
		Entry("cpus #63-64", Set{1 << 63, 1 << 0}, List{{63, 64}}),
		Entry("cpus #63-127", Set{1 << 63, ^uint64(0)}, List{{63, 127}}),

		// multiple all-1s words
		Entry("cpu #0-127", Set{^uint64(0), ^uint64(0)}, List{{0, 127}}),

		// mixed
		Entry("cpu #0-64", Set{^uint64(0), 1 << 0}, List{{0, 64}}),
		Entry("cpu #0-64, 67", Set{^uint64(0), 1<<3 | 1<<0}, List{{0, 64}, {67, 67}}),
		Entry("cpu #65-127, 129", Set{0, ^uint64(0) - 1, 1 << 1}, List{{65, 127}, {129, 129}}),

		Entry("b/w", Set{0xaa0}, List{{5, 5}, {7, 7}, {9, 9}, {11, 11}}),
		Entry("art", Set{0x5a0}, List{{5, 5}, {7, 8}, {10, 10}}),
	)

	It("gets this process's CPU affinity list, consistent with /proc/self/status data", func() {
		Expect(wordbytesize).To(Equal(uint64(64 /* bits in uint64 */ / 8 /* bits/byte*/)))
		cpulist := Successful(Affinity(os.Getpid())).List()
		Expect(cpulist).NotTo(BeEmpty())
		Expect(setsize.Load()).NotTo(BeZero())

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
		Expect(SetAffinity(0, Set{}.SetRange(oneonly, oneonly))).To(Succeed())

		reducedaffs := Successful(Affinity(0)).List()
		Expect(reducedaffs).To(Equal(List{[2]uint{oneonly, oneonly}}))

		Expect(SetAffinity(0, affs)).To(Succeed())
	})

	It("cannot set empty affinities", func() {
		Expect(SetAffinity(0, Set{})).NotTo(Succeed())
		Expect(SetAffinity(0, Set{0})).NotTo(Succeed())
	})

	Context("textual representation", func() {

		It("handles the empty set correctly", func() {
			Expect(Set{}.String()).To(BeEmpty())
		})

		It("returns a textual list representation", func() {
			s := Set{6, 1}
			Expect(s.String()).To(Equal("1-2,64"))
		})

	})

	When("testing CPUs in sets", func() {

		It("returns correct indices", func() {
			Expect(setBitIndex(32)).To(Equal(0))
			Expect(setBitIndex(32 + 2*64)).To(Equal(2))
		})

		It("returns correct bit masks", func() {
			Expect(setBitMask(32)).To(Equal(uint64(1) << 32))
			Expect(setBitMask(32 + 2*64)).To(Equal(uint64(1) << 32))
		})

		It("correctly tests", func() {
			Expect(Set{2}.IsSet(0)).To(BeFalse())
			Expect(Set{2}.IsSet(1)).To(BeTrue())
			Expect(Set{2}.IsSet(666)).To(BeFalse())
		})

	})

	When("setting ranges", func() {

		It("sets CPU ranges", func() {
			Expect(Set{}.SetRange(1, 1).SetRange(63, 65).String()).To(Equal("1,63-65"))
			Expect(Set{0, 0, 0}.SetRange(63, 65).String()).To(Equal("63-65"))
		})

		It("panics on invalid range", func() {
			Expect(func() {
				Set{}.SetRange(3, 1)
			}).To(Panic())
		})

	})

})
