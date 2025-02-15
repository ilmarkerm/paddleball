package main

//
// Copyright (c) 2019 Tony Sarendal <tony@polarcap.org>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//

import (
	//"bytes"
	"log/slog"
	"net"
	"time"
	//"encoding/gob"
	//"github.com/goccy/go-json"
)

func server(ipport string, key int64, lport int, hport int) {
	req := payload{}
	var err error
	var conn net.PacketConn

	conn, err = net.ListenPacket("udp", ipport)
	if err != nil {
		fatal("server:" + err.Error())
	}

	packetBuffer := NewWritePositionBuffer(65536)

	var addr net.Addr
	for {
		packetBuffer.WritePos, addr, err = conn.ReadFrom(packetBuffer.Data)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		err = packetBuffer.UnmarshalPayload(&req)
		packetBuffer.Reset()
		if err != nil {
			slog.Error("error", "message", err, "addr", addr)
			continue
		}
		if req.Key != key {
			continue
		}

		req.Sts = time.Now().UnixNano()
		req.Lport = int64(lport)
		req.Hport = int64(hport)

		err = req.MarshalPayload(packetBuffer)
		if err != nil {
			fatal(err.Error())
		}

		_, err = conn.WriteTo(packetBuffer.Data[:packetBuffer.WritePos], addr)
		packetBuffer.Reset()
		if err != nil {
			slog.Error(err.Error())
			continue
		}
	}
}
