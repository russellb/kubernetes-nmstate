package helper

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	nmstatev1alpha1 "github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1"
)

var (
	badYaml = nmstatev1alpha1.State("}")
	empty   = nmstatev1alpha1.State("")

	noBridges = nmstatev1alpha1.State(`interfaces:
  - name: eth1
    type: ethernet
    state: up
  - name: bond1
    type: bond
    state: up
    link-aggregation:
      mode: active-backup
      slaves:
        - eth1
      options:
        miimon: '120'
`)
	noBridgesUp = nmstatev1alpha1.State(`interfaces:
  - name: eth1
    type: ethernet
    state: up
  - name: br1
    type: linux-bridge
    state: down
  - name: br2
    type: linux-bridge
    state: absent
`)
	someBridgesUp = nmstatev1alpha1.State(`interfaces:
  - name: eth1
    type: ethernet
    state: up
  - name: eth2
    type: ethernet
    state: up
  - name: br1
    type: linux-bridge
    state: up
    bridge:
      options:
        stp:
          enabled: false
      port:
        - name: eth1
          stp-hairpin-mode: false
          stp-path-cost: 100
          stp-priority: 32
  - name: br2
    type: linux-bridge
    state: up
    bridge:
      options:
        stp:
          enabled: false
      port:
        - name: eth2
          stp-hairpin-mode: false
          stp-path-cost: 100
          stp-priority: 32
  - name: br3
    type: linux-bridge
    state: down
  - name: br4
    type: linux-bridge
    state: absent
`)
)

var _ = Describe("Network desired state bridge parser", func() {
	var (
		obtainedBridges []string
		desiredState    nmstatev1alpha1.State
		err             error
	)
	JustBeforeEach(func() {
		obtainedBridges, err = getBridgesUp(desiredState)
	})
	Context("when desired state is not a yaml", func() {
		BeforeEach(func() {
			desiredState = badYaml
		})
		It("should return error", func() {
			Expect(err).To(HaveOccurred())
		})
	})
	Context("when desired state is empty", func() {
		BeforeEach(func() {
			desiredState = empty
		})
		It("should return empty slice", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(obtainedBridges).To(BeEmpty())
		})
	})
	Context("when there is no bridges", func() {
		BeforeEach(func() {
			desiredState = noBridges
		})
		It("should return empty slice", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(obtainedBridges).To(BeEmpty())
		})
	})
	Context("when there are no bridges up", func() {
		BeforeEach(func() {
			desiredState = noBridgesUp
		})
		It("should return empty slice", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(obtainedBridges).To(BeEmpty())
		})
	})
	Context("when there are bridges up", func() {
		BeforeEach(func() {
			desiredState = someBridgesUp
		})
		It("should return the slice with the bridges", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(obtainedBridges).To(ConsistOf("br1", "br2"))
		})
	})
})
