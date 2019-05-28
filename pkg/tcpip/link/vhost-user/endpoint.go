// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package vhostuser provides the implemention of vhostuser data-link layer
// endpoints. Such endpoints would transfer outbound packets to vhost-user
// through vring.
//
// vhostuser endpoints can be used in the networking stack by calling New() to
// create a new endpoint, and then passing it as an argument to
// Stack.CreateNIC().
package vhostuser

import (
	"gvisor.googlesource.com/gvisor/pkg/tcpip"
	"gvisor.googlesource.com/gvisor/pkg/tcpip/buffer"
	"gvisor.googlesource.com/gvisor/pkg/tcpip/stack"
)

// Options specify the details about vhost-user endpoint to be created.
type Options struct {
	MTU        uint32
	Addresses  tcpio.LinkAddress
	SocketPath string
}

type endpoint struct {
	// TODO: @denggx, neccessary fields needed by vhost-user endpoint
	dispatcher stack.NetworkDispatcher
	client     vhostUserClient
	memory     vhostUserMemory
	vringTable [VIRTQUEUE_NUM]vring
}

// New creates a new vhost-user endpoint. This link-layer endpoint would transfer
// outbound packets into the other end of vhost-user through vring buffer.
func New(opt *Options) (tcpip.LinkEndpointID, error) {
	// TODO: @denggx, specific operations here to intialize vhost-user and network
	// related parameters
	ep := endpoint{}

	client, err := initVhostClient(opt.SockPath)
	if err != nil {
		// err handler
		return nil, fmt.Errorf("init vhost client failed: %s", err)
	}
	ep.client = client

	// buffer size should be sum of net header and packet size
	buf_size := bufSizeFromMTU(opt.MTU)

	// memory region size should be sum of vring size and all buffers size
	mem_size := getMemSize(buf_size)

	vum := vhostUserMemory{
		nregions: VIRTQUEUE_NUM,
	}
	for i := 0; i < vum.nregions; i++ {
		fd, size, offset, addr, err := initShmFd(mem_size)
		if err != nil {
			// err handler
			return nil, fmt.Errorf("init share memory failed: %s", err)
		}
		vum.fds[i] = fd

		mr := memoryRegion{
			guestPhysAddr: addr,
			memorySize:    size,
			userspaceAddr: addr,
			mmapOffset:    offset,
		}
		vum.regions[i] = mr
	}
	ep.memory = vum

	if err := vhostUserSetMemTable(vum); err != nil {
		return nil, fmt.Errorf("set memory table failed:%s", err)
	}

	for i := 0; i < vum.nregions; i++ {
		vr, err := initVring(vum.regions[i], vum.fds[i])
		if err != nil {
			return nil, fmt.Errorf("init vring failed:%s", err)
		}
		ep.vringtable[i] = vr
		if err := vhostUserSetVring(vr, i); err != nil {
			return nil, fmt.Errorf("set vring failed: %s", err)
		}

	}
	return stack.RegisterLinkEndpoint(&endpoint{})
}

// Attach implements stack.LinkEndpoint.Attach. It just saves the stack network-
// layer dispatcher for later use when packets need to be dispatched.
func (e *endpoint) Attach(dispatcher stack.NetworkDispatcher) {
	// TODO: @denggx, Attach will launch a dispatch_loop goroutine
	// dispatch_loop is responsible for handling packets delivering
	e.dispatcher = dispatcher
}

// IsAttached implements stack.LinkEndpoint.IsAttached.
func (e *endpoint) IsAttached() bool {
	return e.dispatcher != nil
}

// MTU implements stack.LinkEndpoint.MTU.
func (*endpoint) MTU() uint32 {
	// FIXME: @denggx, should not be a constant here
	return 65536
}

// Capabilities implements stack.LinkEndpoint.Capabilities.
func (*endpoint) Capabilities() stack.LinkEndpointCapabilities {
	// FIXME: @denggx
	return stack.CapabilityRXChecksumOffload | stack.CapabilityTXChecksumOffload | stack.CapabilitySaveRestore | stack.CapabilityLoopback
}

// MaxHeaderLength implements stack.LinkEndpoint.MaxHeaderLength.
func (*endpoint) MaxHeaderLength() uint16 {
	// FIXME: @denggx
	return 0
}

// LinkAddress returns the link address of this endpoint.
func (*endpoint) LinkAddress() tcpip.LinkAddress {
	// FIXME: @denggx
	return ""
}

// WritePacket implements stack.LinkEndpoint.WritePacket. It delivers outbound
// packets to the network-layer dispatcher.
func (e *endpoint) WritePacket(_ *stack.Route, _ *stack.GSO, hdr buffer.Prependable, payload buffer.VectorisedView, protocol tcpip.NetworkProtocolNumber) *tcpip.Error {
	views := make([]buffer.View, 1, 1+len(payload.Views()))
	views[0] = hdr.View()
	views = append(views, payload.Views()...)
	// vv := buffer.NewVectorisedView(len(views[0])+payload.Size(), views)
	_ := buffer.NewVectorisedView(len(views[0])+payload.Size(), views)

	// TODO: @denggx, write packet to the other end of vhost-user

	return nil
}
