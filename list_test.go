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
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/ginkgo/v2/dsl/table"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("cpu lists", func() {

	DescribeTable("generating textual representations",
		func(list List, expected string) {
			Expect(list.String()).To(Equal(expected))
		},
		Entry(nil, List{{1, 1}, {2, 42}, {666, 666}}, "1,2-42,666"),
		Entry(nil, List{{2, 42}}, "2-42"),
		Entry(nil, List{{2, 42}, {777, 778}}, "2-42,777-778"),
	)

	When("parsing lists from text", func() {

		It("returns nothing from nothing", func() {
			Expect(NewList([]byte(""))).To(Equal(List{}))
		})

		It("returns a single cpu", func() {
			Expect(NewList([]byte("42"))).To(Equal(List{[2]uint{42, 42}}))
		})

		It("returns a single range", func() {
			Expect(NewList([]byte("42-666"))).To(Equal(List{[2]uint{42, 666}}))
		})

		It("returns multiple individual CPUs", func() {
			Expect(NewList([]byte("42,666"))).To(Equal(List{[2]uint{42, 42}, [2]uint{666, 666}}))
		})

		It("altogether", func() {
			Expect(NewList([]byte("1-42,666,1000-1001"))).To(
				Equal(List{[2]uint{1, 42}, [2]uint{666, 666}, [2]uint{1000, 1001}}))
		})

		DescribeTable("parsing errors",
			func(s string, msg string) {
				Expect(NewList([]byte(s))).Error().To(MatchError(msg))
			},
			Entry(nil, "abc", "expected unsigned integer number"),
			Entry(nil, "0abc", "expected '-' or ','"),
			Entry(nil, "1-z", "expected unsigned integer number"),
			Entry(nil, "0-0abc", "expected ','"),
		)

	})

	It("converts a list into a set", func() {
		Expect(List{}.Set().String()).To(BeEmpty())
		Expect(Successful(NewList([]byte("3,5,666"))).Set().String()).To(Equal("3,5,666"))
	})

	DescribeTable("overlapping lists",
		func(l1, l2 string, overlapping bool) {
			Expect(Successful(NewList([]byte(l1))).
				Overlap(Successful(NewList([]byte(l2))))).To(Equal(overlapping))
		},
		Entry(nil, "", "", false),
		Entry(nil, "1", "5", false),
		Entry(nil, "1-2", "5-7", false),
		Entry(nil, "5-7", "1-2", false),
		Entry(nil, "1,7,19", "3-5,6-8", true),
		Entry(nil, "3-5,6-8", "1,7,19", true),
		Entry(nil, "7", "1-3,5-999", true),
	)

	DescribeTable("removing CPUs",
		func(l string, cpu int, remainers string) {
			c, rem := Successful(NewList([]byte(l))).Remove()
			Expect(c).To(Equal(uint(cpu)))
			Expect(rem.String()).To(Equal(remainers))
		},
		Entry(nil, "1,3", 1, "3"),
		Entry(nil, "1-2", 1, "2"),
		Entry(nil, "1-3", 1, "2-3"),
		Entry(nil, "5", 5, ""),
	)

	It("panics when there are no cpus to remove", func() {
		Expect(func() {
			_, _ = List{}.Remove()
		}).To(Panic())
	})

})
