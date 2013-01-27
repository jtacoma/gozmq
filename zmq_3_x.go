// +build zmq_3_x

/*
  Copyright 2010-2012 Alec Thomas

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package gozmq

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lzmq
#include <zmq.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

var (
	SNDHWM = IntSocketOption(C.ZMQ_SNDHWM)
	RCVHWM = IntSocketOption(C.ZMQ_SNDHWM)

	// TODO Not documented in the man page...
	//LAST_ENDPOINT       = UInt64SocketOption(C.ZMQ_LAST_ENDPOINT)
	FAIL_UNROUTABLE     = BoolSocketOption(C.ZMQ_FAIL_UNROUTABLE)
	TCP_KEEPALIVE       = IntSocketOption(C.ZMQ_TCP_KEEPALIVE)
	TCP_KEEPALIVE_CNT   = IntSocketOption(C.ZMQ_TCP_KEEPALIVE_CNT)
	TCP_KEEPALIVE_IDLE  = IntSocketOption(C.ZMQ_TCP_KEEPALIVE_IDLE)
	TCP_KEEPALIVE_INTVL = IntSocketOption(C.ZMQ_TCP_KEEPALIVE_INTVL)
	// TODO Make this work.
	//TCP_ACCEPT_FILTER   = IntSocketOption(C.ZMQ_TCP_ACCEPT_FILTER)

	// Message options
	MORE = MessageOption(C.ZMQ_MORE)

	// Send/recv options
	DONTWAIT = SendRecvOption(C.ZMQ_DONTWAIT)
)

type SocketOptions interface {
	SetSockOptInt(option IntSocketOption, value int) error
	SetSockOptInt64(option Int64SocketOption, value int64) error
	SetSockOptUInt64(option UInt64SocketOption, value uint64) error
	SetSockOptString(option StringSocketOption, value string) error
	GetSockOptInt(option IntSocketOption) (value int, err error)
	GetSockOptInt64(option Int64SocketOption) (value int64, err error)
	GetSockOptUInt64(option UInt64SocketOption) (value uint64, err error)
	GetSockOptString(option StringSocketOption) (value string, err error)
	GetSockOptBool(option BoolSocketOption) (value bool, err error)

	// Set options
	SetAffinity(value uint64) error
	SetBacklog(value int) error
	SetIdentity(value string) error
	SetLinger(value int) error
	SetRate(value int64) error
	SetRcvBuf(value uint64) error
	SetRcvHWM(value int) error
	SetReconnectIvl(value int) error
	SetReconnectIvlMax(value int) error
	SetRecoveryIvl(value int64) error
	SetSndBuf(value uint64) error
	SetSndHWM(value int) error
	SetSubscribe(value string) error
	SetTcpKeepalive(value int) error
	SetTcpKeepaliveCnt(value int) error
	SetTcpKeepaliveIdle(value int) error
	SetTcpKeepaliveIntvl(value int) error
	SetUnsubscribe(value string) error

	// Get options
	Affinity() (uint64, error)
	Backlog() (int, error)
	Fd() (int64, error)
	Identity() (string, error)
	Linger() (int, error)
	Rate() (int64, error)
	RcvBuf() (uint64, error)
	RcvHWM() (int, error)
	RcvMore() (rval uint64, e error)
	ReconnectIvl() (int, error)
	ReconnectIvlMax() (int, error)
	RecoveryIvl() (int64, error)
	SndBuf() (uint64, error)
	SndHWM() (int, error)
	SocketType() (uint64, error)
	TcpKeepalive() (int, error)
	TcpKeepaliveCnt() (int, error)
	TcpKeepaliveIdle() (int, error)
	TcpKeepaliveIntvl() (int, error)
}

/* sockopt setters */

func (s *zmqOptions) SetRcvHWM(value int) error {
	return s.SetSockOptInt(RCVHWM, value)
}

func (s *zmqOptions) SetSndHWM(value int) error {
	return s.SetSockOptInt(SNDHWM, value)
}

func (s *zmqOptions) SetTcpKeepalive(value int) error {
	return s.SetSockOptInt(TCP_KEEPALIVE, value)
}

func (s *zmqOptions) SetTcpKeepaliveCnt(value int) error {
	return s.SetSockOptInt(TCP_KEEPALIVE_CNT, value)
}

func (s *zmqOptions) SetTcpKeepaliveIdle(value int) error {
	return s.SetSockOptInt(TCP_KEEPALIVE_IDLE, value)
}

func (s *zmqOptions) SetTcpKeepaliveIntvl(value int) error {
	return s.SetSockOptInt(TCP_KEEPALIVE_INTVL, value)
}

/* sockopt getters */

func (s *zmqOptions) RcvHWM() (int, error) {
	return s.GetSockOptInt(RCVHWM)
}

func (s *zmqOptions) SndHWM() (int, error) {
	return s.GetSockOptInt(SNDHWM)
}

func (s *zmqOptions) TcpKeepalive() (int, error) {
	return s.GetSockOptInt(TCP_KEEPALIVE)
}

func (s *zmqOptions) TcpKeepaliveCnt() (int, error) {
	return s.GetSockOptInt(TCP_KEEPALIVE_CNT)
}

func (s *zmqOptions) TcpKeepaliveIdle() (int, error) {
	return s.GetSockOptInt(TCP_KEEPALIVE_IDLE)
}

func (s *zmqOptions) TcpKeepaliveIntvl() (int, error) {
	return s.GetSockOptInt(TCP_KEEPALIVE_INTVL)
}

// Send a message to the socket.
// int zmq_send (void *s, zmq_msg_t *msg, int flags);
func (s *zmqSocket) Send(data []byte, flags SendRecvOption) error {
	var m C.zmq_msg_t
	// Copy data array into C-allocated buffer.
	size := C.size_t(len(data))

	if C.zmq_msg_init_size(&m, size) != 0 {
		return errno()
	}

	if size > 0 {
		// FIXME Ideally this wouldn't require a copy.
		C.memcpy(unsafe.Pointer(C.zmq_msg_data(&m)), unsafe.Pointer(&data[0]), size) // XXX I hope this works...(seems to)
	}

	if C.zmq_sendmsg(s.s, &m, C.int(flags)) == -1 {
		// zmq_send did not take ownership, free message
		C.zmq_msg_close(&m)
		return errno()
	}
	return nil
}

// Receive a message from the socket.
// int zmq_recv (void *s, zmq_msg_t *msg, int flags);
func (s *zmqSocket) Recv(flags SendRecvOption) (data []byte, err error) {
	// Allocate and initialise a new zmq_msg_t
	var m C.zmq_msg_t
	if C.zmq_msg_init(&m) != 0 {
		err = errno()
		return
	}
	defer C.zmq_msg_close(&m)
	// Receive into message
	if C.zmq_recvmsg(s.s, &m, C.int(flags)) == -1 {
		err = errno()
		return
	}
	// Copy message data into a byte array
	// FIXME Ideally this wouldn't require a copy.
	size := C.zmq_msg_size(&m)
	if size > 0 {
		data = make([]byte, int(size))
		C.memcpy(unsafe.Pointer(&data[0]), C.zmq_msg_data(&m), size)
	} else {
		data = nil
	}
	return
}

// run a zmq_proxy with in, out and capture sockets
func Proxy(in, out, capture Socket) error {
	if C.zmq_proxy(in.apiSocket(), out.apiSocket(), capture.apiSocket()) != 0 {
		return errno()
	}
	return errors.New("zmq_proxy() returned unexpectedly.")
}
