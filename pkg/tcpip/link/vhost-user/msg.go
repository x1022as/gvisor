package vhostuser

import (
	"fmt"
	"net"
)

type uint32 vhostUserRequest

const (
	// vhost-user request type
	VHOST_USER_NONE                  = 0
	VHOST_USER_GET_FEATURES          = 1
	VHOST_USER_SET_FEATURES          = 2
	VHOST_USER_SET_OWNER             = 3
	VHOST_USER_RESET_OWNER           = 4
	VHOST_USER_SET_MEM_TABLE         = 5
	VHOST_USER_SET_LOG_BASE          = 6
	VHOST_USER_SET_LOG_FD            = 7
	VHOST_USER_SET_VRING_NUM         = 8
	VHOST_USER_SET_VRING_ADDR        = 9
	VHOST_USER_SET_VRING_BASE        = 10
	VHOST_USER_GET_VRING_BASE        = 11
	VHOST_USER_SET_VRING_KICK        = 12
	VHOST_USER_SET_VRING_CALL        = 13
	VHOST_USER_SET_VRING_ERR         = 14
	VHOST_USER_GET_PROTOCOL_FEATURES = 15
	VHOST_USER_SET_PROTOCOL_FEATURES = 16
	VHOST_USER_GET_QUEUE_NUM         = 17
	VHOST_USER_SET_VRING_ENABLE      = 18
	VHOST_USER_SEND_RARP             = 19
	VHOST_USER_NET_SET_MTU           = 20
	VHOST_USER_SET_SLAVE_REQ_FD      = 21
	VHOST_USER_IOTLB_MSG             = 22
	VHOST_USER_CRYPTO_CREATE_SESS    = 26
	VHOST_USER_CRYPTO_CLOSE_SESS     = 27
	VHOST_USER_POSTCOPY_ADVISE       = 28
	VHOST_USER_POSTCOPY_LISTEN       = 29
	VHOST_USER_POSTCOPY_END          = 30
	VHOST_USER_MAX                   = 31

	// other constants
	MAX_PAYLOAD_SIZE = 1024
)

type vhostUserMsg struct {
	request vhostUserRequest
	flags   uint32
	size    uint32
	// TODO: @denggx how to express multi message type into one like union in c
	payload [MAX_PAYLOAD_SIZE]byte
}

type vhostUserClient struct {
	sockPath string
	conn     *net.UnixConn
	cfd      int
}

func initVhostClient(path string) (vhostUserClient, error) {
	client := vhostUserClient{}
	addr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return client, fmt.Errorf("Resolve unix socket path %s failed: %s", path, err)
	}
	client.sockPath = path

	c, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return client, fmt.Errorf("Dial unix socket %s failed: %s", path, err)
	}
	client.conn = c
	cf, err := c.File()
	if err != nil {
		return client, fmt.Errorf("Get file from connection failed: %s", err)
	}
	client.cfd = int(cf.Fd())

	return client, nil
}

func vhostUserOneTimeRequest(req vhostUserRequest) bool {
	switch req {
	case VHOST_USER_SET_OWNER:
		return true
	case VHOST_USER_RESET_OWNER:
		return true
	case VHOST_USER_SET_MEM_TABLE:
		return true
	case VHOST_USER_GET_QUEUE_NUM:
		return true
	case VHOST_USER_NET_SET_MTU:
		return true
	default:
		return false
	}
}

func vhostUserGetU64(ep *endpoint, idx int, req vhostUserRequest) (uint64, error) {
	msg := vhostUserMsg{
		request: req,
		flags:   VHOST_USER_VERSION,
	}

	if vhostUserOneTimeRequest(req) && idx != 0 {
		return 0, nil
	}

	if err := vhostUserWrite(ep, msg, nil, 0); err != nil {
		return 0, fmt.Errorf("Vhost user write request(%s) failed: %s", requestStr(req), err)
	}

	if err := vhostUserRead(ep, &msg); err != nil {
		return 0, fmt.Errorf("Vhost user read response(%s) failed: %s", requestString(req), err)
	}

	if msg.request != req {
		return 0, fmt.Errorf("vhost user response type not match, expected %s, get %s", requestStr(req), requestStr(msg.request))
	}

	if msg.size != unsafe.Sizeof(uint64(0)) {
		return 0, fmt.Errorf("vhost user response payload size not match, expected %d, get %d", unsafe.Sizeof(uint64(0)), msg.size)
	}

	return *((*uint64)(unsafe.Pointer(&msg.payload))), nil
}

func vhostUserWrite(ep *endpoint, msg vhostUserMsg, fds []int, fdnum int) error {

}

func vhostUserRead(ep *endpoint, msg *vhostUserMsg) error {

}

/*
   if (vhost_user_read(dev, &msg) < 0) {
       return -1;
   }

   if (msg.hdr.request != request) {
       error_report("Received unexpected msg type. Expected %d received %d",
                    request, msg.hdr.request);
       return -1;
   }

   if (msg.hdr.size != sizeof(msg.payload.u64)) {
       error_report("Received bad msg size.");
       return -1;
   }

   *u64 = msg.payload.u64;

   return 0;

*/

func (ep *endpoint) vhostUserGetFeatures(idx int) (uint64, error) {
	return vhostUserGetU64(ep, idx, VHOST_USER_GET_FEATURES)
}

func requestStr(req vhostUserRequest) string {
	switch req {
	case VHOST_USER_NONE:
		return "none"
	case VHOST_USER_GET_FEATURES:
		return "vhost_user_get_features"
	case VHOST_USER_SET_FEATURES:
		return "vhost_user_set_features"
	case VHOST_USER_SET_OWNER:
		return "vhost_user_set_owner"
	case VHOST_USER_RESET_OWNER:
		return "vhost_user_reset_owner"
	case VHOST_USER_SET_MEM_TABLE:
		return "vhost_user_set_mem_table"
	case VHOST_USER_SET_LOG_BASE:
		return "vhost_user_set_log_base"
	case VHOST_USER_SET_LOG_FD:
		return "vhost_user_set_log_fd"
	case VHOST_USER_SET_VRING_NUM:
		return "vhost_user_set_vring_num"
	case VHOST_USER_SET_VRING_ADDR:
		return "vhost_user_set_vring_addr"
	case VHOST_USER_SET_VRING_BASE:
		return "vhost_user_set_vring_base"
	case VHOST_USER_GET_VRING_BASE:
		return "vhost_user_get_vring_base"
	case VHOST_USER_SET_VRING_KICK:
		return "vhost_user_set_vring_kick"
	case VHOST_USER_SET_VRING_CALL:
		return "vhost_user_set_vring_call"
	case VHOST_USER_SET_VRING_ERR:
		return "vhost_user_set_vring_err"
	case VHOST_USER_GET_PROTOCOL_FEATURES:
		return "vhost_user_get_protocol_features"
	case VHOST_USER_SET_PROTOCOL_FEATURES:
		return "vhost_user_set_protocol_features"
	case VHOST_USER_GET_QUEUE_NUM:
		return "vhost_user_get_queue_num"
	case VHOST_USER_SET_VRING_ENABLE:
		return "vhost_user_set_vring_enable"
	case VHOST_USER_SEND_RARP:
		return "vhost_user_send_rarp"
	case VHOST_USER_NET_SET_MTU:
		return "vhost_user_net_set_mtu"
	case VHOST_USER_SET_SLAVE_REQ_FD:
		return "vhost_user_set_slave_req_fd"
	case VHOST_USER_IOTLB_MSG:
		return "vhost_user_iotlb_msg"
	case VHOST_USER_CRYPTO_CREATE_SESS:
		return "vhost_user_crypto_create_sess"
	case VHOST_USER_CRYPTO_CLOSE_SESS:
		return "vhost_user_crypto_close_sess"
	case VHOST_USER_POSTCOPY_ADVISE:
		return "vhost_user_postcopy_advise"
	case VHOST_USER_POSTCOPY_LISTEN:
		return "vhost_user_postcopy_listen"
	case VHOST_USER_POSTCOPY_END:
		return "vhost_user_postcopy_end"
	default:
		return "invalid vhost-user request"
	}
}
