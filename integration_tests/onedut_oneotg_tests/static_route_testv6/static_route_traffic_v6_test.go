// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/open-traffic-generator/snappi/gosnappi"
	"github.com/openconfig/gribigo/fluent"
	"github.com/openconfig/lemming/gnmi/fakedevice"
	"github.com/openconfig/lemming/integration_tests/binding"
	"github.com/openconfig/ondatra"
	"github.com/openconfig/ondatra/gnmi"
	"github.com/openconfig/ondatra/gnmi/oc"
	"github.com/openconfig/ondatra/gnmi/oc/ocpath"
	"github.com/openconfig/ondatra/gnmi/otg/otgpath"
	"github.com/openconfig/ygnmi/ygnmi"
	"github.com/openconfig/ygot/ygot"
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.Get(".."))
}

// Settings for configuring the baseline testbed with the test
// topology.
//
// The testbed consists of ate:port1 -> dut:port1
// and dut:port2 -> ate:port2.
//
//   - ate:port1 -> dut:port1 subnet 2001::aaaa:bbbb:cccc/99
//   - ate:port2 -> dut:port2 subnet 2001::aaab:bbbb:cccc/99
const (
	ipv6PrefixLen     = 99
	ateDstNetCIDR     = "2003::/49"
	ateIndirectNH     = "2002::"
	ateIndirectNHCIDR = ateIndirectNH + "/49"
	nhIndex           = 1
	nhgIndex          = 42
	nhIndex2          = 2
	nhgIndex2         = 52
	nhIndex3          = 3
	nhgIndex3         = 62
)

// Attributes bundles some common attributes for devices and/or interfaces.
// It provides helpers to generate appropriate configuration for OpenConfig
// and for an ATETopology.  All fields are optional; only those that are
// non-empty will be set when configuring an interface.
type Attributes struct {
	IPv4    string
	IPv6    string
	MAC     string
	Name    string // Interface name, only applied to ATE ports.
	Desc    string // Description, only applied to DUT interfaces.
	IPv4Len uint8  // Prefix length for IPv4.
	IPv6Len uint8  // Prefix length for IPv6.
	MTU     uint16
}

var (
	dutPort1 = Attributes{
		Desc:    "dutPort1",
		IPv6:    "2001::aaaa:bbbb:aa",
		IPv6Len: ipv6PrefixLen,
	}

	atePort1 = Attributes{
		Name:    "port1",
		MAC:     "02:00:01:01:01:01",
		IPv6:    "2001::aaaa:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
	}

	dutPort2 = Attributes{
		Desc:    "dutPort2",
		IPv6:    "2001::aaab:bbbb:aa",
		IPv6Len: ipv6PrefixLen,
	}

	atePort2 = Attributes{
		Name:    "port2",
		MAC:     "02:00:02:01:01:01",
		IPv6:    "2001::aaab:bbbb:bb",
		IPv6Len: ipv6PrefixLen,
	}
)

// configureATE configures port1 and port2 on the ATE.
func configureATE(t *testing.T, ate *ondatra.ATEDevice) gosnappi.Config {
	otg := ate.OTG()
	top := otg.NewConfig(t)

	top.Ports().Add().SetName(atePort1.Name)
	i1 := top.Devices().Add().SetName(atePort1.Name)
	eth1 := i1.Ethernets().Add().SetName(atePort1.Name + ".Eth").
		SetPortName(i1.Name()).SetMac(atePort1.MAC)
	eth1.Ipv6Addresses().Add().SetName(i1.Name() + ".IPv6").
		SetAddress(atePort1.IPv6).SetGateway(dutPort1.IPv6).
		SetPrefix(int32(atePort1.IPv6Len))

	top.Ports().Add().SetName(atePort2.Name)
	i2 := top.Devices().Add().SetName(atePort2.Name)
	eth2 := i2.Ethernets().Add().SetName(atePort2.Name + ".Eth").
		SetPortName(i2.Name()).SetMac(atePort2.MAC)
	eth2.Ipv6Addresses().Add().SetName(i2.Name() + ".IPv6").
		SetAddress(atePort2.IPv6).SetGateway(dutPort2.IPv6).
		SetPrefix(int32(atePort2.IPv6Len))
	return top
}

var gatewayMap = map[Attributes]Attributes{
	atePort1: dutPort1,
	atePort2: dutPort2,
}

// configInterfaceDUT configures the interface with the Addrs.
func configInterfaceDUT(i *oc.Interface, a *Attributes) *oc.Interface {
	i.Description = ygot.String(a.Desc)
	i.Type = oc.IETFInterfaces_InterfaceType_ethernetCsmacd

	s := i.GetOrCreateSubinterface(0)
	s.Enabled = ygot.Bool(true)
	s6 := s.GetOrCreateIpv6()
	s6a := s6.GetOrCreateAddress(a.IPv6)
	s6a.PrefixLength = ygot.Uint8(ipv6PrefixLen)

	return i
}

// configureDUT configures port1 and port2 on the DUT.
func configureDUT(t *testing.T, dut *ondatra.DUTDevice) {
	p1 := dut.Port(t, "port1")
	i1 := &oc.Interface{Name: ygot.String(p1.Name())}
	gnmi.Replace(t, dut, ocpath.Root().Interface(p1.Name()).Config(), configInterfaceDUT(i1, &dutPort1))

	p2 := dut.Port(t, "port2")
	i2 := &oc.Interface{Name: ygot.String(p2.Name())}
	gnmi.Replace(t, dut, ocpath.Root().Interface(p2.Name()).Config(), configInterfaceDUT(i2, &dutPort2))

	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Subinterface(0).Ipv6().Address(dutPort1.IPv6).Ip().State(), time.Minute, dutPort1.IPv6)
	gnmi.Await(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Subinterface(0).Ipv6().Address(dutPort2.IPv6).Ip().State(), time.Minute, dutPort2.IPv6)
}

func waitOTGARPEntry(t *testing.T) {
	ate := ondatra.ATE(t, "ate")

	val, ok := gnmi.WatchAll(t, ate.OTG(), otgpath.Root().InterfaceAny().Ipv6NeighborAny().LinkLayerAddress().State(), time.Minute, func(v *ygnmi.Value[string]) bool {
		return v.IsPresent()
	}).Await(t)
	if !ok {
		t.Fatal("failed to get neighbor")
	}
	lla, _ := val.Val()
	t.Logf("Neighbor %v", lla)
}

// testTraffic generates traffic flow from source network to
// destination network via srcEndPoint to dstEndPoint and checks for
// packet loss and returns loss percentage as float.
func testTraffic(t *testing.T, ate *ondatra.ATEDevice, top gosnappi.Config, srcEndPoint, dstEndPoint Attributes, dur time.Duration) float32 {
	otg := ate.OTG()
	waitOTGARPEntry(t)
	top.Flows().Clear().Items()
	flowipv6 := top.Flows().Add().SetName("Flow")
	flowipv6.Metrics().SetEnable(true)
	flowipv6.TxRx().Device().
		SetTxNames([]string{srcEndPoint.Name + ".IPv6"}).
		SetRxNames([]string{dstEndPoint.Name + ".IPv6"})
	flowipv6.Duration().SetChoice("continuous")
	flowipv6.Packet().Add().Ethernet()
	v6 := flowipv6.Packet().Add().Ipv6()
	v6.Src().SetValue(srcEndPoint.IPv6)
	v6.Dst().Increment().SetStart("2003::").SetCount(250)
	otg.PushConfig(t, top)

	otg.StartTraffic(t)
	time.Sleep(dur)
	t.Logf("Stop traffic")
	otg.StopTraffic(t)

	time.Sleep(5 * time.Second)

	txPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().OutPkts().State())
	rxPkts := gnmi.Get(t, otg, gnmi.OTG().Flow("Flow").Counters().InPkts().State())
	lossPct := (txPkts - rxPkts) * 100 / txPkts
	return float32(lossPct)
}

// awaitTimeout calls a fluent client Await, adding a timeout to the context.
func awaitTimeout(ctx context.Context, c *fluent.GRIBIClient, t testing.TB, timeout time.Duration) error {
	subctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Await(subctx, t)
}

// testCounters test packet counters and should be called after testTraffic
func testCounters(t *testing.T, dut *ondatra.DUTDevice, wantTxPkts, wantRxPkts uint64) {
	got := gnmi.Get(t, dut, ocpath.Root().Interface(dut.Port(t, "port1").Name()).Counters().InPkts().State())
	t.Logf("DUT port 1 in-pkts: %d", got)
	if got < wantTxPkts {
		t.Errorf("DUT got less packets (%d) than OTG sent (%d)", got, wantTxPkts)
	}

	got = gnmi.Get(t, dut, ocpath.Root().Interface(dut.Port(t, "port2").Name()).Counters().OutPkts().State())
	t.Logf("DUT port 2 out-pkts: %d", got)
	if got < wantRxPkts {
		t.Errorf("DUT got sent less packets (%d) than OTG received (%d)", got, wantRxPkts)
	}
}

// TestIPv6Entry tests a single IPv6Entry forwarding entry.
func TestIPv6Entry(t *testing.T) {
	dut := ondatra.DUT(t, "dut")
	configureDUT(t, dut)

	ate := ondatra.ATE(t, "ate")
	ateTop := configureATE(t, ate)
	ate.OTG().PushConfig(t, ateTop)

	cases := []struct {
		desc   string
		routes []*oc.NetworkInstance_Protocol_Static
	}{{
		desc: "Single next-hop",
		routes: []*oc.NetworkInstance_Protocol_Static{{
			Prefix: ygot.String(ateDstNetCIDR),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString(atePort2.IPv6),
					Recurse: ygot.Bool(true),
				},
			},
		}},
	}, {
		desc: "Recursive next-hop",
		routes: []*oc.NetworkInstance_Protocol_Static{{
			Prefix: ygot.String(ateIndirectNHCIDR),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString(atePort2.IPv6),
					Recurse: ygot.Bool(true),
				},
			},
		}, {
			Prefix: ygot.String(ateDstNetCIDR),
			NextHop: map[string]*oc.NetworkInstance_Protocol_Static_NextHop{
				"single": {
					Index:   ygot.String("single"),
					NextHop: oc.UnionString(ateIndirectNH),
					Recurse: ygot.Bool(true),
				},
			},
		}},
	}}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// Install Static Routes
			staticp := ocpath.Root().NetworkInstance(fakedevice.DefaultNetworkInstance).
				Protocol(oc.PolicyTypes_INSTALL_PROTOCOL_TYPE_STATIC, fakedevice.StaticRoutingProtocol)
			for _, route := range tc.routes {
				gnmi.Replace(t, dut, staticp.Static(*route.Prefix).Config(), route)
				//gnmi.Await(t, dut, staticp.Static(*route.Prefix).State(), 30*time.Second, route)
				time.Sleep(10 * time.Second)
				gotRoute := gnmi.Get(t, dut, staticp.Static(*route.Prefix).State())
				fmt.Println(cmp.Diff(route, gotRoute))
			}

			// Send some traffic to make sure neighbor cache is warmed up on the dut.
			testTraffic(t, ate, ateTop, atePort1, atePort2, 1*time.Second)

			loss := testTraffic(t, ate, ateTop, atePort1, atePort2, 15*time.Second)
			if loss > 1 {
				t.Errorf("Loss: got %g, want <= 1", loss)
			}

			// Delete routes and test that there is no more traffic.
			// TODO: Test this when sysrib supports route deletions.
			//for _, route := range tc.routes {
			//	gnmi.Delete(t, dut, staticp.Static(*route.Prefix).Config())
			//}
			//loss = testTraffic(t, ate, ateTop, atePort1, atePort2, 5*time.Second)
			//if loss != 100 {
			//	t.Errorf("Loss: got %g, want 100", loss)
			//}
		})
	}
}