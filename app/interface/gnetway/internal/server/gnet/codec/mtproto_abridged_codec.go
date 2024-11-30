// Copyright (c) 2021-present,  Teamgram Studio (https://teamgram.io).
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codec

import (
	"encoding/binary"
	"fmt"
)

// https://core.telegram.org/mtproto#tcp-transport
//
// There is an abridged version of the same protocol:
// if the client sends 0xef as the first byte (**important:** only prior to the very first data packet),
// then packet length is encoded by a single byte (0x01..0x7e = data length divided by 4;
// or 0x7f followed by 3 length bytes (little endian) divided by 4) followed
// by the data themselves (sequence number and CRC32 not added).
// In this case, server responses look the same (the server does not send 0xefas the first byte).
//

type AbridgedCodec struct {
	*AesCTR128Crypto
	state     int
	packetLen [4]byte
}

func newMTProtoAbridgedCodec(crypto *AesCTR128Crypto) *AbridgedCodec {
	return &AbridgedCodec{
		AesCTR128Crypto: crypto,
		state:           WAIT_PACKET_LENGTH_1,
	}
}

// Encode encodes frames upon server responses into TCP stream.
func (c *AbridgedCodec) Encode(conn CodecWriter, msg []byte) ([]byte, error) {
	// b := message.Encode() d
	sb := make([]byte, 4)
	// minus padding
	size := len(msg) / 4

	if size < 127 {
		sb = []byte{byte(size)}
	} else {
		binary.LittleEndian.PutUint32(sb, uint32(size<<8|127))
	}

	buf := append(sb, msg...)
	return c.Encrypt(buf), nil
}

// Decode decodes frames from TCP stream via specific implementation.
func (c *AbridgedCodec) Decode(conn CodecReader) (bool, []byte, error) {
	var (
		in      innerBuffer
		buf     []byte
		n       int
		err     error
		needAck bool
	)

	in, _ = conn.Peek(-1)
	// log.Debugf("connId: %d, n = %d", conn.ConnID(), len(in))
	if len(in) == 0 {
		return false, nil, nil
	}

	switch c.state {
	case WAIT_PACKET_LENGTH_1:
		if buf, err = in.readN(1); err != nil {
			return false, nil, ErrUnexpectedEOF
		}
		buf = c.Decrypt(buf)
		c.packetLen[0] = buf[0]
		_, _ = conn.Discard(1)

		/*
			### Quick ack
			Some of the TCP transports listed above support quick ACKs: quick ACKs are a way for clients to get quick
			receipt acknowledgements for packets.

			To request a quick ack for a specific outgoing payload, clients must set the MSB of an appropriate field in the
			transport envelope (as described in the documentation for each transport protocol above).

			Also, clients must generate and store a quick ACK token, associating it with the outgoing MTProto payload, by:

			- Taking the first 32 bits of the SHA256 of the encrypted portion of the payload prepended by 32 bytes from the
			  authorization key (the same hash generated when computing the message key, except that instead of taking the
			  middle 128 bits, the first 32 bits are taken instead).
			- Setting the MSB of the last byte to 1: in other words, treat the 32 bits generated above as a little-endian integer,
			  then add 0x80000000 to it (i.e. ack_token = msg_key_long[0:4] | (1 << 31) on a little-endian system).

			Once the payload is successfully received, decrypted and accepted for processing by the server, the server will send
			back the same quick ACK token we generated above, using the encoding described in the documentation for each
			transport protocol.

			Note that reception of a quick ACK does not indicate that any of the RPC queries contained in the message have
			succeeded, failed or finished execution at all, it simply indicates that they have been received, decrypted and
			accepted for processing by the server.

			The server will still send msgs_ack constructors for content-related constructors and methods contained in
			payloads which were quick ACKed, as well as replies/errors for methods and constructors, as usual.
		*/
		needAck = c.packetLen[0]>>7 == 1
		_ = needAck

		n = int(c.packetLen[0] & 0x7f)
		if n < 0x7f {
			c.state = WAIT_PACKET_LENGTH_1_PACKET
			n = n << 2
			// log.Debugf("n = %d", n)
		} else {
			c.state = WAIT_PACKET_LENGTH_3
			if buf, err = in.readN(3); err != nil {
				return false, nil, ErrUnexpectedEOF
			}
			buf = c.Decrypt(buf)
			c.packetLen[1] = buf[0]
			c.packetLen[2] = buf[1]
			c.packetLen[3] = buf[2]
			_, _ = conn.Discard(3)

			c.state = WAIT_PACKET_LENGTH_3_PACKET
			n = (int(c.packetLen[1]) | int(c.packetLen[2])<<8 | int(c.packetLen[3])<<16) << 2
			// log.Debugf("n = %d", n)
			if n > MAX_MTPRORO_FRAME_SIZE {
				// TODO(@benqi): close conn
				return false, nil, fmt.Errorf("too large data(%d)", n)
			}
		}
		if buf, err = in.readN(n); err != nil {
			return false, nil, ErrUnexpectedEOF
		} else if len(buf) <= 8 {
			// TODO: fix
			return false, nil, ErrUnexpectedEOF
		}

		buf = c.Decrypt(buf)
		_, _ = conn.Discard(n)
		c.state = WAIT_PACKET_LENGTH_1

		// message := mtproto.NewMTPRawMessage(int64(binary.LittleEndian.Uint64(buf)), 0, TRANSPORT_TCP)
		// _ = message.Decode(buf)

		return needAck, buf, nil
	case WAIT_PACKET_LENGTH_1_PACKET:
		n = int(c.packetLen[0]&0x7f) << 2
		if buf, err = in.readN(n); err != nil {
			return false, nil, ErrUnexpectedEOF
		} else if len(buf) <= 8 {
			// TODO: fix
			return false, nil, ErrUnexpectedEOF
		}
		// log.Debugf("n = %d", n)

		buf = c.Decrypt(buf)
		_, _ = conn.Discard(n)
		c.state = WAIT_PACKET_LENGTH_1

		// message := mtproto.NewMTPRawMessage(int64(binary.LittleEndian.Uint64(buf)), 0, TRANSPORT_TCP)
		// _ = message.Decode(buf)

		return needAck, buf, nil
	case WAIT_PACKET_LENGTH_3:
		if buf, err = in.readN(3); err != nil {
			return false, nil, ErrUnexpectedEOF
		}
		buf = c.Decrypt(buf)
		c.packetLen[1] = buf[0]
		c.packetLen[2] = buf[1]
		c.packetLen[3] = buf[2]
		_, _ = conn.Discard(3)

		c.state = WAIT_PACKET_LENGTH_3_PACKET
		n = (int(c.packetLen[1]) | int(c.packetLen[2])<<8 | int(c.packetLen[3])<<16) << 2
		// log.Debugf("n = %d", n)
		if n > MAX_MTPRORO_FRAME_SIZE {
			// TODO(@benqi): close conn
			return false, nil, fmt.Errorf("too large data(%d)", n)
		}
		if buf, err = in.readN(n); err != nil {
			return false, nil, ErrUnexpectedEOF
		} else if len(buf) <= 8 {
			// TODO: fix
			return false, nil, ErrUnexpectedEOF
		}

		buf = c.Decrypt(buf)
		_, _ = conn.Discard(n)
		c.state = WAIT_PACKET_LENGTH_1

		// message := mtproto.NewMTPRawMessage(int64(binary.LittleEndian.Uint64(buf)), 0, TRANSPORT_TCP)
		// _ = message.Decode(buf)

		return needAck, buf, nil
	case WAIT_PACKET_LENGTH_3_PACKET:
		n = (int(c.packetLen[1]) | int(c.packetLen[2])<<8 | int(c.packetLen[3])<<16) << 2
		// log.Debugf("n = %d", n)
		if n > MAX_MTPRORO_FRAME_SIZE {
			// TODO(@benqi): close conn
			return false, nil, fmt.Errorf("too large data(%d)", n)
		}
		if buf, err = in.readN(n); err != nil {
			return false, nil, ErrUnexpectedEOF
		} else if len(buf) <= 8 {
			// TODO: fix
			return false, nil, ErrUnexpectedEOF
		}

		buf = c.Decrypt(buf)
		_, _ = conn.Discard(n)
		c.state = WAIT_PACKET_LENGTH_1

		// message := mtproto.NewMTPRawMessage(int64(binary.LittleEndian.Uint64(buf)), 0, TRANSPORT_TCP)
		// _ = message.Decode(buf)

		return needAck, buf, nil
	}

	// TODO(@benqi): close conn
	return false, nil, fmt.Errorf("unknown error")
}
