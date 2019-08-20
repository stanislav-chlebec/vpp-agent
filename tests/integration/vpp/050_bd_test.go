//  Copyright (c) 2019 Cisco and/or its affiliates.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at:
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package vpp

import (
	"github.com/ligato/cn-infra/logging/logrus"
	//	"strings"
	"testing"

	vpp_l2 "github.com/ligato/vpp-agent/api/models/vpp/l2"
	//_ "github.com/ligato/vpp-agent/api/models/vpp/l2"
	"github.com/ligato/vpp-agent/plugins/vpp/ifplugin/ifaceidx"
	ifplugin_vppcalls "github.com/ligato/vpp-agent/plugins/vpp/ifplugin/vppcalls"
	_ "github.com/ligato/vpp-agent/plugins/vpp/l2plugin"
	l2plugin_vppcalls "github.com/ligato/vpp-agent/plugins/vpp/l2plugin/vppcalls"
	_ "github.com/ligato/vpp-agent/plugins/vpp/l3plugin"
	//"github.com/ligato/vpp-agent/plugins/vpp/l3plugin/vrfidx"
	"github.com/ligato/vpp-agent/pkg/idxvpp"
	. "github.com/onsi/gomega"
)

const (
	dummyBridgeDomain     = 4
	dummyBridgeDomainName = "bridge_domain"
)

// Input test data for creating bridge domain
var createTestDataInBD *vpp_l2.BridgeDomain = &vpp_l2.BridgeDomain{
	Name:                dummyBridgeDomainName,
	Flood:               true,
	UnknownUnicastFlood: true,
	Forward:             true,
	Learn:               true,
	ArpTermination:      true,
	MacAge:              45,
}

func TestBd(t *testing.T) {
	ctx := setupVPP(t)
	defer ctx.teardownVPP()

	ih := ifplugin_vppcalls.CompatibleInterfaceVppHandler(ctx.vppBinapi, logrus.NewLogger("test"))
	Expect(ih).To(Not(BeNil()), "Handler should be created.")
	const ifName = "loop1"
	ifIdx, err := ih.AddLoopbackInterface(ifName)
	if err != nil {
		t.Fatalf("creating interface failed: %v", err)
	}
	t.Logf("interface created %v", ifIdx)

	ifIndexes := ifaceidx.NewIfaceIndex(logrus.NewLogger("test-if"), "test-if")
	ifIndexes.Put(ifName, &ifaceidx.IfaceMetadata{SwIfIndex: ifIdx})

	bdIdx := idxvpp.NewNameToIndex(logrus.NewLogger("test-if2"), "vpp-bd-index", nil)
	t.Logf("%v", bdIdx)
	h := l2plugin_vppcalls.CompatibleL2VppHandler(ctx.vppBinapi, ifIndexes, bdIdx, logrus.NewLogger("test"))
	t.Logf("%v", h)
	Expect(h).To(Not(BeNil()), "Handler should be created.")
	bd, err := h.DumpBridgeDomains()
	t.Logf("%v", bd)

	err = h.AddBridgeDomain(dummyBridgeDomain, createTestDataInBD)

	Expect(err).ShouldNot(HaveOccurred())
	bd, err = h.DumpBridgeDomains()
	t.Logf("%v", bd[0])

	//the second time te same bd
	err = h.AddBridgeDomain(dummyBridgeDomain, createTestDataInBD)

	Expect(err).Should(HaveOccurred())
	bd, err = h.DumpBridgeDomains()
	t.Logf("%v", bd[0])

	err = h.DeleteBridgeDomain(dummyBridgeDomain)
	Expect(err).ShouldNot(HaveOccurred())
	bd, err = h.DumpBridgeDomains()
	t.Logf("%v", bd)
	err = h.DeleteBridgeDomain(dummyBridgeDomain)
	Expect(err).Should(HaveOccurred())

	err = h.AddBridgeDomain(dummyBridgeDomain, createTestDataInBD)
	Expect(err).ShouldNot(HaveOccurred())
	bd, err = h.DumpBridgeDomains()
	t.Logf("%v", bd)

}
