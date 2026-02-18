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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("online cpus", func() {

	It("returns an empty list when the online cpu list cannot be accessed", func() {
		Expect(online("_test/online")).To(Equal(List{}))
	})

	It("returns an empty list when the online cpu list is empty or broken", func() {
		Expect(online("_test/online/empty")).To(Equal(List{}))
		Expect(online("_test/online/broken")).To(Equal(List{}))
	})

	It("correctly reads the online cpu list", func() {
		Expect(online("_test/online/1,3-42")).To(
			Equal(Successful(NewList([]byte("1,3-42")))))
	})

	It("fetches the system's online cpu list", func() {
		Expect(Online()).NotTo(BeEmpty())
	})

})
