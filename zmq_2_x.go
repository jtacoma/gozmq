// +build !zmq_3_x

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
import "unsafe"

var (
	RECOVERY_IVL_MSEC = Int64SocketOption(C.ZMQ_RECOVERY_IVL_MSEC)
	SWAP              = Int64SocketOption(C.ZMQ_SWAP)
	MCAST_LOOP        = Int64SocketOption(C.ZMQ_MCAST_LOOP)
	HWM               = UInt64SocketOption(C.ZMQ_HWM)
	NOBLOCK           = SendRecvOption(C.ZMQ_NOBLOCK)
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
	SetHWM(value uint64) error
	SetIdentity(value string) error
	SetLinger(value int) error
	SetMcastLoop(value int64) error
	SetRate(value int64) error
	SetRcvBuf(value uint64) error
	SetReconnectIvl(value int) error
	SetReconnectIvlMax(value int) error
	SetRecoveryIvl(value int64) error
	SetRecoveryIvlMsec(value int64) error
	SetSndBuf(value uint64) error
	SetSubscribe(value string) error
	SetSwap(value int64) error
	SetUnsubscribe(value string) error

	// Get options
	Affinity() (uint64, error)
	Backlog() (int, error)
	Fd() (int64, error)
	HWM() (uint64, error)
	Identity() (string, error)
	Linger() (int, error)
	McastLoop() (int64, error)
	Rate() (int64, error)
	RcvBuf() (uint64, error)
	RcvMore() (rval uint64, e error)
	ReconnectIvl() (int, error)
	ReconnectIvlMax() (int, error)
	RecoveryIvl() (int64, error)
	RecoveryIvlMsec() (int64, error)
	SndBuf() (uint64, error)
	SocketType() (uint64, error)
	Swap() (int64, error)
}

/* sockopt setters */

func (s *zmqOptions) SetHWM(value uint64) error {
	return s.SetSockOptUInt64(HWM, value)
}

func (s *zmqOptions) SetSwap(value int64) error {
	return s.SetSockOptInt64(SWAP, value)
}

func (s *zmqOptions) SetRecoveryIvlMsec(value int64) error {
	return s.SetSockOptInt64(RECOVERY_IVL_MSEC, value)
}

func (s *zmqOptions) SetMcastLoop(value int64) error {
	return s.SetSockOptInt64(MCAST_LOOP, value)
}

/* sockopt getters */

func (s *zmqOptions) HWM() (uint64, error) {
	return s.GetSockOptUInt64(HWM)
}

func (s *zmqOptions) Swap() (int64, error) {
	return s.GetSockOptInt64(SWAP)
}

func (s *zmqOptions) RecoveryIvlMsec() (int64, error) {
	return s.GetSockOptInt64(RECOVERY_IVL_MSEC)
}

func (s *zmqOptions) McastLoop() (int64, error) {
	return s.GetSockOptInt64(MCAST_LOOP)
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

	if C.zmq_send(s.s, &m, C.int(flags)) != 0 {
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
	if C.zmq_recv(s.s, &m, C.int(flags)) != 0 {
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
