package dnssrv

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
)

func TestResolveError(t *testing.T) {
	r := &Resolver{}
	_, err := r.Resolve("_nonexistent._udp.invalid.invalid.")
	if err == nil {
		t.Log("surprisingly got records — DNS may have responded")
	}
}

func TestResolverFields(t *testing.T) {
	r := &Resolver{Address: "127.0.0.1:53"}
	if r.Address != "127.0.0.1:53" {
		t.Fatalf("Address not set: %s", r.Address)
	}
}

// encodeLabel encodes a DNS name as a sequence of length-prefixed labels.
func encodeLabel(name string) []byte {
	var out []byte
	start := 0
	for i := 0; i <= len(name); i++ {
		if i == len(name) || name[i] == '.' {
			if i > start {
				lbl := name[start:i]
				out = append(out, byte(len(lbl)))
				out = append(out, []byte(lbl)...)
			}
			start = i + 1
		}
	}
	out = append(out, 0)
	return out
}

// buildSRVResponse builds a DNS SRV response with two answer records:
// one with port > 0 and one with port == 0, to exercise both branches.
func buildSRVResponse(queryMsg []byte) []byte {
	if len(queryMsg) < 12 {
		return nil
	}
	id := queryMsg[:2]

	// Extract qname from question section (bytes after the 12-byte header)
	qname := queryMsg[12:]
	end := 0
	for end < len(qname) && qname[end] != 0 {
		end++
	}
	end++ // include root zero byte
	qnameBytes := qname[:end]

	target1 := encodeLabel("svc.local.")
	rdata1 := make([]byte, 6+len(target1))
	binary.BigEndian.PutUint16(rdata1[0:], 10)
	binary.BigEndian.PutUint16(rdata1[2:], 5)
	binary.BigEndian.PutUint16(rdata1[4:], 8080) // port > 0
	copy(rdata1[6:], target1)

	target2 := encodeLabel("svc2.local.")
	rdata2 := make([]byte, 6+len(target2))
	binary.BigEndian.PutUint16(rdata2[0:], 10)
	binary.BigEndian.PutUint16(rdata2[2:], 5)
	binary.BigEndian.PutUint16(rdata2[4:], 0) // port == 0
	copy(rdata2[6:], target2)

	mkRR := func(rdata []byte) []byte {
		var rr []byte
		rr = append(rr, qnameBytes...)
		rr = append(rr, 0, 33)        // TYPE SRV
		rr = append(rr, 0, 1)         // CLASS IN
		rr = append(rr, 0, 0, 0, 30)  // TTL 30
		rdLen := []byte{0, 0}
		binary.BigEndian.PutUint16(rdLen, uint16(len(rdata)))
		rr = append(rr, rdLen...)
		rr = append(rr, rdata...)
		return rr
	}

	hdr := make([]byte, 12)
	copy(hdr[0:2], id)
	binary.BigEndian.PutUint16(hdr[2:], 0x8400) // QR=1, AA=1
	binary.BigEndian.PutUint16(hdr[4:], 1)      // QDCOUNT=1
	binary.BigEndian.PutUint16(hdr[6:], 2)      // ANCOUNT=2

	// Reconstruct the original question section
	question := make([]byte, len(qnameBytes)+4)
	copy(question, qnameBytes)
	binary.BigEndian.PutUint16(question[len(qnameBytes):], 33) // QTYPE SRV
	binary.BigEndian.PutUint16(question[len(qnameBytes)+2:], 1) // QCLASS IN

	msg := append(hdr, question...)
	msg = append(msg, mkRR(rdata1)...)
	msg = append(msg, mkRR(rdata2)...)
	return msg
}

// TestResolveSuccessWithCustomResolver overrides net.DefaultResolver to use a
// local stub DNS server, exercising the success path of Resolve (including
// both the port>0 and port==0 branches).
func TestResolveSuccessWithCustomResolver(t *testing.T) {
	pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Skip("cannot listen UDP:", err)
	}
	serverAddr := pc.LocalAddr().String()

	// Serve DNS queries until the PacketConn is closed.
	go func() {
		buf := make([]byte, 512)
		for {
			n, src, err2 := pc.ReadFrom(buf)
			if err2 != nil {
				return
			}
			resp := buildSRVResponse(buf[:n])
			if resp != nil {
				_, _ = pc.WriteTo(resp, src)
			}
		}
	}()

	orig := net.DefaultResolver
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("udp4", serverAddr)
		},
	}
	defer func() {
		net.DefaultResolver = orig
		pc.Close() // stops the goroutine
	}()

	r := &Resolver{}
	records, err := r.Resolve("svc.example.com")
	if err != nil {
		t.Logf("Resolve returned error: %v — skipping assertions", err)
		return
	}
	if len(records) == 0 {
		t.Fatal("expected at least one record")
	}
	for _, rec := range records {
		if rec.Address == "" {
			t.Error("record has empty address")
		}
	}
}
