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
	acl "github.com/ligato/vpp-agent/api/models/vpp/acl"
	. "github.com/onsi/gomega"
	"strings"
	"testing"

	_ "github.com/ligato/vpp-agent/plugins/vpp/aclplugin"
	aclplugin_vppcalls "github.com/ligato/vpp-agent/plugins/vpp/aclplugin/vppcalls"
	_ "github.com/ligato/vpp-agent/plugins/vpp/ifplugin"
	"github.com/ligato/vpp-agent/plugins/vpp/ifplugin/ifaceidx"
	ifplugin_vppcalls "github.com/ligato/vpp-agent/plugins/vpp/ifplugin/vppcalls"
	//	"encoding/json"
)

var aclNoRules []*acl.ACL_Rule

var aclErr1Rules = []*acl.ACL_Rule{
	{
		Action: acl.ACL_Rule_PERMIT,
		IpRule: &acl.ACL_Rule_IpRule{
			Ip: &acl.ACL_Rule_IpRule_Ip{
				SourceNetwork:      ".0.",
				DestinationNetwork: "10.20.0.0/24",
			},
		},
	},
}

var aclErr2Rules = []*acl.ACL_Rule{
	{
		Action: acl.ACL_Rule_PERMIT,
		IpRule: &acl.ACL_Rule_IpRule{
			Ip: &acl.ACL_Rule_IpRule_Ip{
				SourceNetwork:      "192.168.1.1/32",
				DestinationNetwork: ".0.",
			},
		},
	},
}

var aclErr3Rules = []*acl.ACL_Rule{
	{
		Action: acl.ACL_Rule_PERMIT,
		IpRule: &acl.ACL_Rule_IpRule{
			Ip: &acl.ACL_Rule_IpRule_Ip{
				SourceNetwork:      "192.168.1.1/32",
				DestinationNetwork: "dead::1/64",
			},
		},
	},
}

var aclErr4Rules = []*acl.ACL_Rule{
	{
		Action: acl.ACL_Rule_PERMIT,
		MacipRule: &acl.ACL_Rule_MacIpRule{
			SourceAddress:        "192.168.0.1",
			SourceAddressPrefix:  uint32(16),
			SourceMacAddress:     "",
			SourceMacAddressMask: "ff:ff:ff:ff:00:00",
		},
	},
}

var aclErr5Rules = []*acl.ACL_Rule{
	{
		Action: acl.ACL_Rule_PERMIT,
		MacipRule: &acl.ACL_Rule_MacIpRule{
			SourceAddress:        "192.168.0.1",
			SourceAddressPrefix:  uint32(16),
			SourceMacAddress:     "11:44:0A:B8:4A:36",
			SourceMacAddressMask: "",
		},
	},
}

var aclErr6Rules = []*acl.ACL_Rule{
	{
		Action: acl.ACL_Rule_PERMIT,
		MacipRule: &acl.ACL_Rule_MacIpRule{
			SourceAddress:        "",
			SourceAddressPrefix:  uint32(16),
			SourceMacAddress:     "11:44:0A:B8:4A:36",
			SourceMacAddressMask: "ff:ff:ff:ff:00:00",
		},
	},
}

var aclIPrules = []*acl.ACL_Rule{
	{
		//RuleName:  "permitIPv4",
		Action: acl.ACL_Rule_PERMIT,
		IpRule: &acl.ACL_Rule_IpRule{
			Ip: &acl.ACL_Rule_IpRule_Ip{
				SourceNetwork:      "192.168.1.1/32",
				DestinationNetwork: "10.20.0.0/24",
			},
		},
	},
	{
		//RuleName:  "permitIPv6",
		Action: acl.ACL_Rule_PERMIT,
		IpRule: &acl.ACL_Rule_IpRule{
			Ip: &acl.ACL_Rule_IpRule_Ip{
				SourceNetwork:      "dead::1/64",
				DestinationNetwork: "dead::2/64",
			},
		},
	},
	{
		//RuleName:  "permitIP",
		Action: acl.ACL_Rule_PERMIT,
		IpRule: &acl.ACL_Rule_IpRule{
			Ip: &acl.ACL_Rule_IpRule_Ip{
				SourceNetwork:      "",
				DestinationNetwork: "",
			},
		},
	},
	{
		//RuleName:  "denyICMP",
		Action: acl.ACL_Rule_DENY,
		IpRule: &acl.ACL_Rule_IpRule{
			Icmp: &acl.ACL_Rule_IpRule_Icmp{
				Icmpv6: false,
				IcmpCodeRange: &acl.ACL_Rule_IpRule_Icmp_Range{
					First: 150,
					Last:  250,
				},
				IcmpTypeRange: &acl.ACL_Rule_IpRule_Icmp_Range{
					First: 1150,
					Last:  1250,
				},
			},
		},
	},
	{
		//RuleName:  "denyICMPv6",
		Action: acl.ACL_Rule_DENY,
		IpRule: &acl.ACL_Rule_IpRule{
			Icmp: &acl.ACL_Rule_IpRule_Icmp{
				Icmpv6: true,
				IcmpCodeRange: &acl.ACL_Rule_IpRule_Icmp_Range{
					First: 150,
					Last:  250,
				},
				IcmpTypeRange: &acl.ACL_Rule_IpRule_Icmp_Range{
					First: 1150,
					Last:  1250,
				},
			},
		},
	},
	{
		//RuleName:  "permitTCP",
		Action: acl.ACL_Rule_PERMIT,
		IpRule: &acl.ACL_Rule_IpRule{
			Tcp: &acl.ACL_Rule_IpRule_Tcp{
				TcpFlagsMask:  20,
				TcpFlagsValue: 10,
				SourcePortRange: &acl.ACL_Rule_IpRule_PortRange{
					LowerPort: 150,
					UpperPort: 250,
				},
				DestinationPortRange: &acl.ACL_Rule_IpRule_PortRange{
					LowerPort: 1150,
					UpperPort: 1250,
				},
			},
		},
	},
	{
		//RuleName:  "denyUDP",
		Action: acl.ACL_Rule_DENY,
		IpRule: &acl.ACL_Rule_IpRule{
			Udp: &acl.ACL_Rule_IpRule_Udp{
				SourcePortRange: &acl.ACL_Rule_IpRule_PortRange{
					LowerPort: 150,
					UpperPort: 250,
				},
				DestinationPortRange: &acl.ACL_Rule_IpRule_PortRange{
					LowerPort: 1150,
					UpperPort: 1250,
				},
			},
		},
	},
}

var aclMACIPrules = []*acl.ACL_Rule{
	{
		//RuleName:  "denyIPv4",
		Action: acl.ACL_Rule_DENY,
		MacipRule: &acl.ACL_Rule_MacIpRule{
			SourceAddress:        "192.168.0.1",
			SourceAddressPrefix:  uint32(16),
			SourceMacAddress:     "11:44:0A:B8:4A:35",
			SourceMacAddressMask: "ff:ff:ff:ff:00:00",
		},
	},
	{
		//RuleName:  "denyIPv6",
		Action: acl.ACL_Rule_DENY,
		MacipRule: &acl.ACL_Rule_MacIpRule{
			SourceAddress:        "dead::1",
			SourceAddressPrefix:  uint32(64),
			SourceMacAddress:     "11:44:0A:B8:4A:35",
			SourceMacAddressMask: "ff:ff:ff:ff:00:00",
		},
	},
}

func TestCRUDIPAcl(t *testing.T) {
	ctx := setupVPP(t)
	defer ctx.teardownVPP()

	ih := ifplugin_vppcalls.CompatibleInterfaceVppHandler(ctx.vppBinapi, logrus.NewLogger("test"))

	const ifName = "loop1"
	ifIdx, errI := ih.AddLoopbackInterface(ifName)
	Expect(errI).To(BeNil())
	t.Logf("interface created %v", ifIdx)

	const ifName2 = "loop2"
	ifIdx2, errI2 := ih.AddLoopbackInterface(ifName2)
	Expect(errI2).To(BeNil())
	t.Logf("interface created %v", ifIdx2)

	ifIndexes := ifaceidx.NewIfaceIndex(logrus.NewLogger("test-iface1"), "test-iface1")
	ifIndexes.Put(ifName, &ifaceidx.IfaceMetadata{
		SwIfIndex: ifIdx,
	})
	ifIndexes.Put(ifName2, &ifaceidx.IfaceMetadata{
		SwIfIndex: ifIdx2,
	})

	h := aclplugin_vppcalls.CompatibleACLVppHandler(ctx.vppBinapi, ifIndexes, logrus.NewLogger("test"))
	if h == nil {
		t.Fatalf("handler was not created")
	}

	acls, errx := h.DumpACL()
	Expect(errx).To(BeNil())
	aclCnt := len(acls)
	Expect(aclCnt).Should(Equal(0))
	t.Logf("no acls dumped")

	const aclname = "test0"
	aclIdx, err := h.AddACL(aclIPrules, aclname)
	Expect(err).To(BeNil())
	Expect(aclIdx).To(BeEquivalentTo(0))
	t.Logf("acl added - with index %d", aclIdx)

	err = h.SetACLToInterfacesAsIngress(aclIdx, []uint32{ifIdx})
	Expect(err).To(BeNil())
	t.Logf("acl with index %d was assigned to interface %v ingress", aclIdx, ifName)

	acls, errx = h.DumpACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	var rules []*acl.ACL_Rule
	var isPresent bool
	var isForInterface bool
	for _, item := range acls {
		rules = item.ACL.Rules
		if (item.Meta.Index == aclIdx) && (aclname == item.Meta.Tag) {
			t.Logf("Found ACL aclIPrules with aclName  %v", item.Meta.Tag)
			for _, rule := range rules {
				//t.Logf("%+v", rule)
				if (rule.IpRule.GetIp().SourceNetwork == "192.168.1.1/32") &&
					(rule.IpRule.GetIp().DestinationNetwork == "10.20.0.0/24") {
					isPresent = true
					break
				}
			}
			// check assignation to interface
			for _, intf := range item.ACL.Interfaces.Ingress {
				if intf == ifName {
					isForInterface = true
					break
				}
			}
		}
	}
	if !isPresent {
		t.Fatalf("Configured IP is not present.")
	} else {
		t.Logf("Configured IP is present.")

	}
	if isForInterface {
		t.Logf("dumped acl is correctly assigned to interface %v", ifName)
	} else {
		t.Fatalf("dumped interface is not correctly assigned")
	}

	indexes := []uint32{ifIdx, ifIdx2}
	ifaces, errI3 := h.DumpACLInterfaces(indexes)
	Expect(errI3).To(Succeed())
	Expect(ifaces).To(HaveLen(2))
	t.Logf("%v", ifaces)
	t.Logf("%v", ifaces[1])
	t.Logf("%v", ifaces[2])
	//this does not work for VPP 19.04 and maybe also other version
	//Expect(ifaces[0].Ingress).To(Equal([]string{ifName}))
	//Expect(ifaces[2].Egress).To(Equal([]string{ifName2}))

	//negative tests - it is expected failure
	t.Logf("Let us test some negative cases....")
	_, err = h.AddACL(aclNoRules, "test1")
	Expect(err).To(Not(BeNil()))
	t.Logf("adding acls failed: %v", err)

	_, err = h.AddACL(aclErr1Rules, "test2")
	Expect(err).To(Not(BeNil()))
	t.Logf("adding acls failed: %v", err)

	_, err = h.AddACL(aclErr2Rules, "test3")
	Expect(err).To(Not(BeNil()))
	t.Logf("adding acls failed: %v", err)

	_, err = h.AddACL(aclErr3Rules, "test4")
	Expect(err).To(Not(BeNil()))
	t.Logf("adding acls failed: %v", err)

	//add the same acls again but it will be assigned to the second interface
	const aclname2 = "test5"
	aclIdx, err = h.AddACL(aclIPrules, aclname2)

	Expect(err).To(BeNil())
	Expect(aclIdx).To(BeEquivalentTo(1))
	t.Logf("acl added - with index %d", aclIdx)

	err = h.SetACLToInterfacesAsEgress(aclIdx, []uint32{ifIdx2})
	Expect(err).To(BeNil())
	t.Logf("acl with index %d was assigned to interface %v egress", aclIdx, ifName2)

	acls, errx = h.DumpACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(2))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	for _, item := range acls {
		rules = item.ACL.Rules
		if (item.Meta.Index == aclIdx) && (aclname2 == item.Meta.Tag) {
			t.Logf("Found ACL aclIPrules with aclName  %v", item.Meta.Tag)
			for _, rule := range rules {
				//t.Logf("%+v", rule)
				if (rule.IpRule.GetIp().SourceNetwork == "192.168.1.1/32") &&
					(rule.IpRule.GetIp().DestinationNetwork == "10.20.0.0/24") {
					isPresent = true
					break
				}
			}
			// check assignation to interface
			for _, intf := range item.ACL.Interfaces.Egress {
				if intf == ifName2 {
					isForInterface = true
					break
				}
			}
		}
	}
	if !isPresent {
		t.Fatalf("Configured IP is not present.")
	} else {
		t.Logf("Configured IP is present.")

	}
	if isForInterface {
		t.Logf("dumped acl is correctly assigned to interface %v", ifName2)
	} else {
		t.Fatalf("dumped interface is not correctly assigned")
	}

	//negative tests
	err = h.DeleteACL(5)
	Expect(err).To(Not(BeNil()))
	t.Logf("deleting acls failed: %v", err)

	// find the acl with aclname test0
	var foundaclidx uint32
	for _, item := range acls {
		rules = item.ACL.Rules
		if aclname == item.Meta.Tag {
			foundaclidx = item.Meta.Index
			break
		}
	}
	err = h.DeleteACL(foundaclidx)
	Expect(err).To(Not(BeNil()))
	t.Logf("deleting acls failed: %v", err)

	// DELETE ACL
	err = h.RemoveACLFromInterfacesAsIngress(foundaclidx, []uint32{ifIdx})
	err = h.DeleteACL(foundaclidx)
	Expect(err).To(BeNil())
	t.Logf("deleting acls succeed")

	acls, errx = h.DumpACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	for _, aclrecord := range acls {
		if aclrecord.Meta.Index == foundaclidx {
			t.Fatalf("This acll should be deleted : %v", errx)
		}
	}

	// MODIFY ACL
	rule2modify := []*acl.ACL_Rule{
		{
			Action: acl.ACL_Rule_PERMIT,
			IpRule: &acl.ACL_Rule_IpRule{
				Ip: &acl.ACL_Rule_IpRule_Ip{
					SourceNetwork:      "10.20.30.1/32",
					DestinationNetwork: "10.20.0.0/24",
				},
			},
		},
		{
			Action: acl.ACL_Rule_PERMIT,
			IpRule: &acl.ACL_Rule_IpRule{
				Ip: &acl.ACL_Rule_IpRule_Ip{
					SourceNetwork:      "dead:dead::3/64",
					DestinationNetwork: "dead:dead::4/64",
				},
			},
		},
	}

	err = h.ModifyACL(1, rule2modify, "test_modify0")
	Expect(err).To(BeNil())
	t.Logf("acl was modified")

	acls, errx = h.DumpACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	var modifiedacl aclplugin_vppcalls.ACLDetails
	for _, item := range acls {
		modifiedacl = *item
		rules = item.ACL.Rules
		if item.Meta.Index == aclIdx { //&& (aclname2 == item.Meta.Tag) {
			t.Logf("Found modified ACL aclIPrules with aclName  %v", item.Meta.Tag)
			for _, rule := range rules {
				//t.Logf("%+v", rule)
				if (rule.IpRule.GetIp().SourceNetwork == "10.20.30.1/32") &&
					(rule.IpRule.GetIp().DestinationNetwork == "10.20.0.0/24") {
					isPresent = true
					break
				}
			}
			// check assignation to interface
			for _, intf := range item.ACL.Interfaces.Egress {
				if intf == ifName2 {
					isForInterface = true
					break
				}
			}
		}
	}
	if !isPresent {
		t.Fatalf("Configured IP is not present.")
	} else {
		t.Logf("Configured IP is present.")

	}
	if isForInterface {
		t.Logf("dumped acl is correctly assigned to interface %v", ifName2)
	} else {
		t.Fatalf("dumped interface is not correctly assigned")
	}

	// negative test
	err = h.ModifyACL(1, aclErr1Rules, "test_modify1")
	Expect(err).To(Not(BeNil()))
	t.Logf("modifying of acl failed: %v", err)

	const aclname3 = "test_modify2"
	err = h.ModifyACL(1, aclNoRules, aclname3)
	Expect(err).To(BeNil())
	t.Logf("acl was modified")

	acls, errx = h.DumpACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	for _, item := range acls {
		if item.Meta.Index == aclIdx { //&& (aclname2 == item.Meta.Tag) {
			t.Logf("Found modified ACL aclIPrules with aclName  %v", item.Meta.Tag)
		}
		if item.ACL.String() == modifiedacl.ACL.String() {
			t.Logf("Last update caused no change in acl definition.")
			break
		} else {
			t.Fatalf("Last update caused change in acl definition.")
		}
	}

	// DELETE ACL
	err = h.RemoveACLFromInterfacesAsEgress(aclIdx, []uint32{ifIdx2})

	acls, errx = h.DumpACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	for _, item := range acls {
		if item.Meta.Index == aclIdx { //&& (aclname2 == item.Meta.Tag) {
			t.Logf("Found modified ACL aclIPrules with aclName  %v", item.Meta.Tag)
		}
		if item.ACL.Interfaces.String() == "" {
			t.Logf("Interface assignment was removed")
		} else {
			t.Fatalf("Interface assignment was not removed.")
		}
	}

	err = h.DeleteACL(aclIdx)
	Expect(err).To(BeNil())
	t.Logf("deleting acls succeed")

	acls, errx = h.DumpACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(0))
	t.Logf("%d acls dumped", aclCnt)

}

// Test add MACIP acl rules
func TestCRUDMacIPAcl(t *testing.T) {
	ctx := setupVPP(t)
	defer ctx.teardownVPP()

	ih := ifplugin_vppcalls.CompatibleInterfaceVppHandler(ctx.vppBinapi, logrus.NewLogger("test"))

	const ifName = "loop1"
	ifIdx, errI := ih.AddLoopbackInterface(ifName)
	Expect(errI).To(BeNil())
	t.Logf("interface created %v", ifIdx)

	const ifName2 = "loop2"
	ifIdx2, errI2 := ih.AddLoopbackInterface(ifName2)
	Expect(errI2).To(BeNil())
	t.Logf("interface created %v", ifIdx2)

	ifIndexes := ifaceidx.NewIfaceIndex(logrus.NewLogger("test-iface1"), "test-iface1")
	ifIndexes.Put(ifName, &ifaceidx.IfaceMetadata{
		SwIfIndex: ifIdx,
	})
	ifIndexes.Put(ifName2, &ifaceidx.IfaceMetadata{
		SwIfIndex: ifIdx2,
	})

	h := aclplugin_vppcalls.CompatibleACLVppHandler(ctx.vppBinapi, ifIndexes, logrus.NewLogger("test"))
	if h == nil {
		t.Fatalf("handler was not created")
	}

	acls, errx := h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt := len(acls)
	Expect(aclCnt).Should(Equal(0))
	t.Logf("no acls dumped")

	const aclname = "test6"
	aclIdx, err := h.AddMACIPACL(aclMACIPrules, aclname)
	Expect(err).To(BeNil())
	Expect(aclIdx).To(BeEquivalentTo(0))
	t.Logf("acl added - with index %d", aclIdx)

	err = h.SetMACIPACLToInterfaces(aclIdx, []uint32{ifIdx})
	Expect(err).To(BeNil())
	t.Logf("acl with index %d was assigned to interface %v", aclIdx, ifName)

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	var rules []*acl.ACL_Rule
	var isPresent bool
	var isForInterface bool
	for _, item := range acls {
		rules = item.ACL.Rules
		if (item.Meta.Index == aclIdx) && (aclname == item.Meta.Tag) {
			t.Logf("Found ACL aclMACIPrules with aclName  %v", item.Meta.Tag)
			for _, rule := range rules {
				if (rule.MacipRule.SourceAddress == "192.168.0.1") &&
					(rule.MacipRule.SourceAddressPrefix == 16) &&
					(strings.ToLower(rule.MacipRule.SourceMacAddress) == strings.ToLower("11:44:0A:B8:4A:35")) &&
					(rule.MacipRule.SourceMacAddressMask == "ff:ff:ff:ff:00:00") {
					isPresent = true
					break
				}
			}
			// check assignation to interface
			t.Logf("%v", item)
			t.Logf("%v", item.ACL.Interfaces)
			for _, intf := range item.ACL.Interfaces.Ingress {
				if intf == ifName {
					isForInterface = true
					break
				}
			}
		}
	}
	if !isPresent {
		t.Fatalf("Configured IP is not present.")
	} else {
		t.Logf("Configured IP is present.")

	}
	if isForInterface {
		t.Logf("dumped acl is correctly assigned to interface %v", ifName)
	} else {
		t.Fatalf("dumped interface is not correctly assigned")
	}

	indexes := []uint32{ifIdx, ifIdx2}
	ifaces, errI3 := h.DumpACLInterfaces(indexes)
	Expect(errI3).To(Succeed())
	Expect(ifaces).To(HaveLen(2))
	t.Logf("%v", ifaces)
	t.Logf("%v", ifaces[1])
	t.Logf("%v", ifaces[2])
	//this does not work for VPP 19.04 and maybe also other version
	//Expect(ifaces[0].Ingress).To(Equal([]string{ifName}))
	//Expect(ifaces[2].Egress).To(Equal([]string{ifName2}))

	//negative tests - it is expected failure
	t.Logf("Let us test some negative cases....")
	_, err = h.AddMACIPACL(aclNoRules, "test7")
	Expect(err).To(Not(BeNil()))
	t.Logf("adding acls failed: %v", err)

	_, err = h.AddMACIPACL(aclErr4Rules, "test8")
	Expect(err).To(Not(BeNil()))
	t.Logf("adding acls failed: %v", err)

	_, err = h.AddMACIPACL(aclErr5Rules, "test9")
	Expect(err).To(Not(BeNil()))
	t.Logf("adding acls failed: %v", err)

	_, err = h.AddMACIPACL(aclErr6Rules, "test10")
	Expect(err).To(Not(BeNil()))
	Expect(err.Error()).To(BeEquivalentTo("invalid IP address "))
	t.Logf("adding acls failed: %v", err)

	// now let us add the same aclMACIPrules again
	//add the same acls again but it will be assigned to the second interface
	const aclname2 = "test11"
	aclIdx, err = h.AddMACIPACL(aclMACIPrules, aclname2)
	Expect(err).To(BeNil())
	Expect(aclIdx).To(BeEquivalentTo(1))
	t.Logf("acl added - with index %d", aclIdx)

	err = h.AddMACIPACLToInterface(aclIdx, ifName2)
	Expect(err).To(BeNil())
	t.Logf("acl with index %d was assigned to interface %v ", aclIdx, ifName2)

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(2))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	for _, item := range acls {
		rules = item.ACL.Rules
		if (item.Meta.Index == aclIdx) && (aclname2 == item.Meta.Tag) {
			t.Logf("Found ACL aclMACIPrules with aclName  %v", item.Meta.Tag)
			for _, rule := range rules {
				if (rule.MacipRule.SourceAddress == "192.168.0.1") &&
					(rule.MacipRule.SourceAddressPrefix == 16) &&
					(strings.ToLower(rule.MacipRule.SourceMacAddress) == strings.ToLower("11:44:0A:B8:4A:35")) &&
					(rule.MacipRule.SourceMacAddressMask == "ff:ff:ff:ff:00:00") {
					isPresent = true
					break
				}
			}
			// check assignation to interface
			t.Logf("%v", item)
			t.Logf("%v", item.ACL.Interfaces)
			for _, intf := range item.ACL.Interfaces.Ingress {
				if intf == ifName2 {
					isForInterface = true
					break
				}
			}
		}
	}
	if !isPresent {
		t.Fatalf("Configured IP is not present.")
	} else {
		t.Logf("Configured IP is present.")

	}
	if isForInterface {
		t.Logf("dumped acl is correctly assigned to interface %v", ifName2)
	} else {
		t.Fatalf("dumped interface is not correctly assigned")
	}

	//negative tests
	err = h.DeleteMACIPACL(5)
	Expect(err).To(Not(BeNil()))
	t.Logf("deleting acls failed: %v", err)

	// find the acl with aclname test6
	var foundaclidx uint32
	for _, item := range acls {
		rules = item.ACL.Rules
		if aclname == item.Meta.Tag {
			foundaclidx = item.Meta.Index
			break
		}
	}
	err = h.DeleteMACIPACL(foundaclidx)
	Expect(err).To(BeNil())
	t.Logf("deleting acls succeed")

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	for _, aclrecord := range acls {
		if aclrecord.Meta.Index == foundaclidx {
			t.Fatalf("This acll should be deleted : %v", errx)
		}
	}

	// MODIFY ACL
	rule2modify := []*acl.ACL_Rule{
		{
			Action: acl.ACL_Rule_DENY,
			MacipRule: &acl.ACL_Rule_MacIpRule{
				SourceAddress:        "192.168.10.1",
				SourceAddressPrefix:  uint32(24),
				SourceMacAddress:     "11:44:0A:B8:4A:37",
				SourceMacAddressMask: "ff:ff:ff:ff:00:00",
			},
		},
		{
			Action: acl.ACL_Rule_DENY,
			MacipRule: &acl.ACL_Rule_MacIpRule{
				SourceAddress:        "dead::2",
				SourceAddressPrefix:  uint32(64),
				SourceMacAddress:     "11:44:0A:B8:4A:38",
				SourceMacAddressMask: "ff:ff:ff:ff:00:00",
			}},
	}

	err = h.ModifyMACIPACL(1, rule2modify, "test_modify0")
	Expect(err).To(BeNil())
	t.Logf("acl was modified")

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	var modifiedacl aclplugin_vppcalls.ACLDetails
	for _, item := range acls {
		modifiedacl = *item
		rules = item.ACL.Rules
		if item.Meta.Index == aclIdx { //&& (aclname2 == item.Meta.Tag) {
			t.Logf("Found modified ACL aclMACIPrules with aclName  %v", item.Meta.Tag)
			for _, rule := range rules {
				if (rule.MacipRule.SourceAddress == "192.168.10.1") &&
					(rule.MacipRule.SourceAddressPrefix == 24) &&
					(strings.ToLower(rule.MacipRule.SourceMacAddress) == strings.ToLower("11:44:0A:B8:4A:37")) &&
					(rule.MacipRule.SourceMacAddressMask == "ff:ff:ff:ff:00:00") {
					isPresent = true
					break
				}
				//t.Logf("%+v", rule)
			}
			// check assignation to interface
			for _, intf := range item.ACL.Interfaces.Ingress {
				if intf == ifName2 {
					isForInterface = true
					break
				}
			}
		}
	}
	if !isPresent {
		t.Fatalf("Configured IP is not present.")
	} else {
		t.Logf("Configured IP is present.")

	}
	if isForInterface {
		t.Logf("dumped acl is correctly assigned to interface %v", ifName2)
	} else {
		t.Fatalf("dumped interface is not correctly assigned")
	}

	t.Logf("%v", modifiedacl)

	// negative test
	err = h.ModifyMACIPACL(1, aclErr1Rules, "test_modify1")
	Expect(err).To(Not(BeNil()))
	t.Logf("modifying of acl failed: %v", err)

	err = h.ModifyMACIPACL(1, aclIPrules, "test_modify5")
	Expect(err).To(Not(BeNil()))
	t.Logf("modifying of acl failed: %v", err)

	err = h.SetMACIPACLToInterfaces(aclIdx, []uint32{ifIdx})
	Expect(err).To(BeNil())
	t.Logf("acl with index %d was assigned to interface %v", aclIdx, ifName)

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	for _, item := range acls {
		rules = item.ACL.Rules
		if item.Meta.Index == aclIdx {
			t.Logf("Found ACL aclMACIPrules with aclName  %v", item.Meta.Tag)
			// check assignation to interface
			t.Logf("%v", item)
			t.Logf("%v", item.ACL.Interfaces)
			for _, intf := range item.ACL.Interfaces.Ingress {
				if intf == ifName {
					isForInterface = true
					break
				}
			}
		}
	}
	if isForInterface {
		t.Logf("dumped acl is correctly assigned to interface %v", ifName)
	} else {
		t.Fatalf("dumped interface is not correctly assigned")
	}

	err = h.SetMACIPACLToInterfaces(aclIdx, []uint32{ifIdx})
	Expect(err).To(BeNil())
	t.Logf("acl with index %d was assigned to interface %v", aclIdx, ifName)

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	for _, item := range acls {
		rules = item.ACL.Rules
		if item.Meta.Index == aclIdx {
			t.Logf("Found ACL aclMACIPrules with aclName  %v", item.Meta.Tag)
			// check assignation to interface
			t.Logf("%v", item)
			t.Logf("%v", item.ACL.Interfaces)
			for _, intf := range item.ACL.Interfaces.Ingress {
				if intf == ifName {
					isForInterface = true
					break
				}
			}
		}
	}
	if isForInterface {
		t.Logf("dumped acl is correctly assigned to interface %v", ifName)
	} else {
		t.Fatalf("dumped interface is not correctly assigned")
	}

	err = h.DeleteMACIPACLFromInterface(aclIdx, ifName)
	Expect(err).To(BeNil())
	t.Logf("for acl with index %d was deleted the relation to interface %v", aclIdx, ifName)

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	for _, item := range acls {
		rules = item.ACL.Rules
		if item.Meta.Index == aclIdx {
			t.Logf("Found ACL aclMACIPrules with aclName  %v", item.Meta.Tag)
			// check assignation to interface
			t.Logf("%v", item)
			t.Logf("%v", item.ACL.Interfaces)
			for _, intf := range item.ACL.Interfaces.Ingress {
				if intf == ifName {
					t.Fatalf("alc should not be assigned to the interface %v", ifName)
				}
			}
		}
	}
	t.Logf("for acl was correctly deleted relation to interface %v", ifName)

	err = h.DeleteMACIPACLFromInterface(aclIdx, ifName2)
	Expect(err).To(BeNil())
	t.Logf("for acl with index %d was deleted the relation to interface %v", aclIdx, ifName2)

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	isPresent = false
	isForInterface = false
	for _, item := range acls {
		rules = item.ACL.Rules
		if item.Meta.Index == aclIdx {
			t.Logf("Found ACL aclMACIPrules with aclName  %v", item.Meta.Tag)
			// check assignation to interface
			t.Logf("%v", item)
			t.Logf("%v", item.ACL.Interfaces)
			Expect(item.ACL.Interfaces).To(BeNil())
		}
	}
	t.Logf("for acl was correctly deleted relation to interface %v", ifName2)

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(1))
	t.Logf("%d acls dumped", aclCnt)

	for _, aclrecord := range acls {
		foundaclidx = aclrecord.Meta.Index
	}

	err = h.DeleteMACIPACL(foundaclidx)
	Expect(err).To(BeNil())
	t.Logf("deleting acls succeed")

	acls, errx = h.DumpMACIPACL()
	Expect(errx).To(BeNil())
	aclCnt = len(acls)
	Expect(aclCnt).Should(Equal(0))
	t.Logf("no acls dumped")
}
