package vhostuser

const (
	VHOST_VRING_SIZE = 32 * 1024 // number of vring buffers, may change later

	VIRTIO_DESC_F_NEXT     = 1 // descriptor continues via 'next' field
	VIRTIO_DESC_F_WRITE    = 2 // write-only descriptor(otherwise read-only)
	VIRTIO_DESC_F_INDIRECT = 4 // buffer contains a list of descriptors

	VHOST_CLIENT_VRING_IDX_RX = 0 // rx vring index
	VHOST_CLIENT_VRING_IDX_TX = 1 // tx vring index
	VHOST_CLIENT_VRING_NUM    = 2 // vring number

	VHOST_USER_VERSION = 0x1
)

// vringDesc is buffer descriptor
type vringDesc struct {
	addr   uint64 // packet data buffer address
	length uint32 // packet data buffer size
	flags  uint16 // vring flags
	next   uint16 // next descriptor in chain
}

type vringAvail struct {
	flags uint16
	idx   uint16
	ring  [VHOST_VRING_SIZE]uint16
}

type vringUsedElem struct {
	id     uint32
	length uint32
}

type vringUsed struct {
	flags uint16
	idx   uint16
	ring  [VHOST_VRING_SIZE]vringUsedElem
}

type vhostVring struct {
	kickfd int
	callfd int
	// FIXME: @denggx how to align struct field
	desc  [VHOST_VRING_SIZE]vringDesc
	avail vringAvail
	used  vringUsed
}

type vring struct {
	kickfd, callfd int
	desc           *vringDesc
	avail          *vringAvail
	used           *vringUsed
	num            uint
	last_avail_idx uint16
	last_used_idx  uint16
}

type vhostUserMemoryRegion struct {
	guestPhysAddr uint64
	memorySize    uint64
	userspaceAddr uint64
	mmapOffset    uint64
}

type vhostUserMemory struct {
	nregions uint32
	padding  uint32
	regions  [VHOST_CLIENT_VRING_NUM]vhostUserMemoryRegion
	fds      [VHOST_CLIENT_VRING_NUM]int
}
