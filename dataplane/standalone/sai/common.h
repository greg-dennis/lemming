// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef DATAPLANE_STANDALONE_SAI_COMMON_H_
#define DATAPLANE_STANDALONE_SAI_COMMON_H_

#include <glog/logging.h>

#include <algorithm>
#include <memory>
#include <string>

#include "dataplane/proto/acl.grpc.pb.h"
#include "dataplane/proto/bfd.grpc.pb.h"
#include "dataplane/proto/bmtor.grpc.pb.h"
#include "dataplane/proto/bridge.grpc.pb.h"
#include "dataplane/proto/buffer.grpc.pb.h"
#include "dataplane/proto/common.pb.h"
#include "dataplane/proto/counter.grpc.pb.h"
#include "dataplane/proto/debug_counter.grpc.pb.h"
#include "dataplane/proto/dtel.grpc.pb.h"
#include "dataplane/proto/fdb.grpc.pb.h"
#include "dataplane/proto/generic_programmable.grpc.pb.h"
#include "dataplane/proto/hash.grpc.pb.h"
#include "dataplane/proto/hostif.grpc.pb.h"
#include "dataplane/proto/ipmc.grpc.pb.h"
#include "dataplane/proto/ipmc_group.grpc.pb.h"
#include "dataplane/proto/ipsec.grpc.pb.h"
#include "dataplane/proto/isolation_group.grpc.pb.h"
#include "dataplane/proto/l2mc.grpc.pb.h"
#include "dataplane/proto/l2mc_group.grpc.pb.h"
#include "dataplane/proto/lag.grpc.pb.h"
#include "dataplane/proto/macsec.grpc.pb.h"
#include "dataplane/proto/mcast_fdb.grpc.pb.h"
#include "dataplane/proto/mirror.grpc.pb.h"
#include "dataplane/proto/mpls.grpc.pb.h"
#include "dataplane/proto/my_mac.grpc.pb.h"
#include "dataplane/proto/nat.grpc.pb.h"
#include "dataplane/proto/neighbor.grpc.pb.h"
#include "dataplane/proto/next_hop.grpc.pb.h"
#include "dataplane/proto/next_hop_group.grpc.pb.h"
#include "dataplane/proto/policer.grpc.pb.h"
#include "dataplane/proto/port.grpc.pb.h"
#include "dataplane/proto/qos_map.grpc.pb.h"
#include "dataplane/proto/queue.grpc.pb.h"
#include "dataplane/proto/route.grpc.pb.h"
#include "dataplane/proto/router_interface.grpc.pb.h"
#include "dataplane/proto/rpf_group.grpc.pb.h"
#include "dataplane/proto/samplepacket.grpc.pb.h"
#include "dataplane/proto/scheduler.grpc.pb.h"
#include "dataplane/proto/scheduler_group.grpc.pb.h"
#include "dataplane/proto/srv6.grpc.pb.h"
#include "dataplane/proto/stp.grpc.pb.h"
#include "dataplane/proto/switch.grpc.pb.h"
#include "dataplane/proto/system_port.grpc.pb.h"
#include "dataplane/proto/tam.grpc.pb.h"
#include "dataplane/proto/tunnel.grpc.pb.h"
#include "dataplane/proto/udf.grpc.pb.h"
#include "dataplane/proto/virtual_router.grpc.pb.h"
#include "dataplane/proto/vlan.grpc.pb.h"
#include "dataplane/proto/wred.grpc.pb.h"

extern "C" {
#include "inc/sai.h"
}

extern std::unique_ptr<lemming::dataplane::sai::Acl::Stub> acl;
extern std::unique_ptr<lemming::dataplane::sai::Bfd::Stub> bfd;
extern std::unique_ptr<lemming::dataplane::sai::Buffer::Stub> buffer;
extern std::unique_ptr<lemming::dataplane::sai::Bmtor::Stub> bmtor;
extern std::unique_ptr<lemming::dataplane::sai::Bridge::Stub> bridge;
extern std::unique_ptr<lemming::dataplane::sai::Counter::Stub> counter;
extern std::unique_ptr<lemming::dataplane::sai::DebugCounter::Stub>
    debug_counter;
extern std::unique_ptr<lemming::dataplane::sai::Dtel::Stub> dtel;
extern std::unique_ptr<lemming::dataplane::sai::Fdb::Stub> fdb;
extern std::unique_ptr<lemming::dataplane::sai::GenericProgrammable::Stub>
    generic_programmable;
extern std::unique_ptr<lemming::dataplane::sai::Hash::Stub> hash;
extern std::unique_ptr<lemming::dataplane::sai::Hostif::Stub> hostif;
extern std::unique_ptr<lemming::dataplane::sai::IpmcGroup::Stub> ipmc_group;
extern std::unique_ptr<lemming::dataplane::sai::Ipmc::Stub> ipmc;
extern std::unique_ptr<lemming::dataplane::sai::Ipsec::Stub> ipsec;
extern std::unique_ptr<lemming::dataplane::sai::IsolationGroup::Stub>
    isolation_group;
extern std::unique_ptr<lemming::dataplane::sai::L2mcGroup::Stub> l2mc_group;
extern std::unique_ptr<lemming::dataplane::sai::L2mc::Stub> l2mc;
extern std::unique_ptr<lemming::dataplane::sai::Lag::Stub> lag;
extern std::unique_ptr<lemming::dataplane::sai::Macsec::Stub> macsec;
extern std::unique_ptr<lemming::dataplane::sai::Mirror::Stub> mirror;
extern std::unique_ptr<lemming::dataplane::sai::McastFdb::Stub> mcast_fdb;
extern std::unique_ptr<lemming::dataplane::sai::Mpls::Stub> mpls;
extern std::unique_ptr<lemming::dataplane::sai::MyMac::Stub> my_mac;
extern std::unique_ptr<lemming::dataplane::sai::Nat::Stub> nat;
extern std::unique_ptr<lemming::dataplane::sai::Neighbor::Stub> neighbor;
extern std::unique_ptr<lemming::dataplane::sai::NextHopGroup::Stub>
    next_hop_group;
extern std::unique_ptr<lemming::dataplane::sai::NextHop::Stub> next_hop;
extern std::unique_ptr<lemming::dataplane::sai::Policer::Stub> policer;
extern std::unique_ptr<lemming::dataplane::sai::Port::Stub> port;
extern std::unique_ptr<lemming::dataplane::sai::QosMap::Stub> qos_map;
extern std::unique_ptr<lemming::dataplane::sai::Queue::Stub> queue;
extern std::unique_ptr<lemming::dataplane::sai::Route::Stub> route;
extern std::unique_ptr<lemming::dataplane::sai::RouterInterface::Stub>
    router_interface;
extern std::unique_ptr<lemming::dataplane::sai::RpfGroup::Stub> rpf_group;
extern std::unique_ptr<lemming::dataplane::sai::Samplepacket::Stub>
    samplepacket;
extern std::unique_ptr<lemming::dataplane::sai::SchedulerGroup::Stub>
    scheduler_group;
extern std::unique_ptr<lemming::dataplane::sai::Scheduler::Stub> scheduler;
extern std::unique_ptr<lemming::dataplane::sai::Srv6::Stub> srv6;
extern std::unique_ptr<lemming::dataplane::sai::Stp::Stub> stp;
extern std::shared_ptr<lemming::dataplane::sai::Switch::Stub> switch_;
extern std::unique_ptr<lemming::dataplane::sai::SystemPort::Stub> system_port;
extern std::unique_ptr<lemming::dataplane::sai::Tam::Stub> tam;
extern std::unique_ptr<lemming::dataplane::sai::Tunnel::Stub> tunnel;
extern std::unique_ptr<lemming::dataplane::sai::Udf::Stub> udf;
extern std::unique_ptr<lemming::dataplane::sai::VirtualRouter::Stub>
    virtual_router;
extern std::unique_ptr<lemming::dataplane::sai::Vlan::Stub> vlan;
extern std::unique_ptr<lemming::dataplane::sai::Wred::Stub> wred;

std::string convert_from_ip_addr(sai_ip_addr_family_t addr_family,
                                 const sai_ip_addr_t &addr);
std::string convert_from_ip_address(const sai_ip_address_t &val);
lemming::dataplane::sai::RouteEntry convert_from_route_entry(
    const sai_route_entry_t &entry);
lemming::dataplane::sai::IpPrefix convert_from_ip_prefix(
    const sai_ip_prefix_t &ip_prefix);

sai_ip_addr_t convert_to_ip_addr(std::string val);
sai_ip_address_t convert_to_ip_address(std::string str);
sai_route_entry_t convert_to_route_entry(
    const lemming::dataplane::sai::RouteEntry &entry);
sai_ip_prefix_t convert_to_ip_prefix(
    const lemming::dataplane::sai::IpPrefix &ip_prefix);
std::vector<sai_port_oper_status_notification_t> convert_to_oper_status(
    const lemming::dataplane::sai::PortStateChangeNotificationResponse &resp);

lemming::dataplane::sai::NeighborEntry convert_from_neighbor_entry(
    const sai_neighbor_entry_t &entry);

sai_neighbor_entry_t convert_to_neighbor_entry(
    const lemming::dataplane::sai::NeighborEntry &entry);

void convert_to_acl_capability(
    sai_acl_capability_t &out,
    const lemming::dataplane::sai::ACLCapability &in);

lemming::dataplane::sai::AclActionData convert_from_acl_action_data(
    const sai_acl_action_data_t &in, sai_object_id_t id);

lemming::dataplane::sai::AclActionData convert_from_acl_action_data_action(
    const sai_acl_action_data_t &in, sai_int32_t id);

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data_ip_type(
    const sai_acl_field_data_t &in, sai_int32_t type, sai_int32_t mask);

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data(
    const sai_acl_field_data_t &in, sai_ip4_t data, sai_ip4_t mask);

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data(
    const sai_acl_field_data_t &in, sai_uint8_t data, sai_uint8_t mask);

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data(
    const sai_acl_field_data_t &in, sai_uint16_t data, sai_uint16_t mask);

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data(
    const sai_acl_field_data_t &in, sai_object_id_t data);

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data_ip6(
    const sai_acl_field_data_t &in, const sai_ip6_t data, const sai_ip6_t mask);

lemming::dataplane::sai::AclFieldData convert_from_acl_field_data_mac(
    const sai_acl_field_data_t &in, const sai_mac_t data, const sai_mac_t mask);

// copy_list copies a scalar proto list to an attribute.
// Note: It is expected that the attribute list contains preallocated memory.
template <typename T, typename S>
void copy_list(S *dst, const google::protobuf::RepeatedField<T> &src,
               uint32_t *attr_len) {
  // It's not safe to just memcpy this because in some cases to proto types are
  // larger than the corresponding sai types.
  *attr_len =
      static_cast<uint32_t>(std::min(static_cast<int>(*attr_len), src.size()));
  for (uint32_t i = 0; i < *attr_len; i++) {
    dst[i] = src[i];
  }
}

class PortStateReactor
    : public grpc::ClientReadReactor<
          lemming::dataplane::sai::PortStateChangeNotificationResponse> {
 public:
  PortStateReactor(std::shared_ptr<lemming::dataplane::sai::Switch::Stub> stub,
                   sai_port_state_change_notification_fn callback) {
    this->callback = callback;
    lemming::dataplane::sai::PortStateChangeNotificationRequest req;
    stub->async()->PortStateChangeNotification(&context, &req, this);
    StartRead(&resp);
    StartCall();
  }

  void OnReadDone(bool ok) override {
    if (!ok) return;
    std::vector<sai_port_oper_status_notification_t> v =
        convert_to_oper_status(resp);
    callback(v.size(), v.data());
    StartRead(&resp);
  }

  void OnDone(const grpc::Status &status) override {
    if (status.ok()) {
      LOG(INFO) << "PortStateChangeNotification RPC succeeded.";
    } else {
      LOG(ERROR) << "PortStateChangeNotification RPC failed.";
    }
  }

 private:
  grpc::ClientContext context;
  lemming::dataplane::sai::PortStateChangeNotificationResponse resp;
  sai_port_state_change_notification_fn callback;
};

#endif  // DATAPLANE_STANDALONE_SAI_COMMON_H_
