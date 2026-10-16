package server

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// WSConn wraps a hijacked net.Conn for RFC 6455 WebSocket communication.
type WSConn struct {
	conn net.Conn
	bufr *bufio.Reader
	mu   sync.Mutex
}

// Upgrade upgrades an HTTP connection to a WebSocket connection.
func Upgrade(w http.ResponseWriter, r *http.Request) (*WSConn, error) {
	if strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		return nil, errors.New("expected websocket upgrade")
	}

	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return nil, errors.New("missing Sec-WebSocket-Key")
	}

	h := sha1.New()
	h.Write([]byte(key + wsGUID))
	acceptKey := base64.StdEncoding.EncodeToString(h.Sum(nil))

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("http.ResponseWriter does not support hijacking")
	}

	conn, bufrw, err := hijacker.Hijack()
	if err != nil {
		return nil, fmt.Errorf("hijack error: %w", err)
	}

	response := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n\r\n",
		acceptKey,
	)

	if _, err := conn.Write([]byte(response)); err != nil {
		conn.Close()
		return nil, fmt.Errorf("handshake write error: %w", err)
	}

	return &WSConn{
		conn: conn,
		bufr: bufrw.Reader,
	}, nil
}

// Close closes the connection.
func (c *WSConn) Close() error {
	return c.conn.Close()
}

// ReadMessage reads an unmasked or masked frame from the client.
func (c *WSConn) ReadMessage() (byte, []byte, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(c.bufr, header); err != nil {
		return 0, nil, err
	}

	opcode := header[0] & 0x0F
	isMasked := (header[1] & 0x80) != 0
	payloadLen := int64(header[1] & 0x7F)

	if payloadLen == 126 {
		lenBytes := make([]byte, 2)
		if _, err := io.ReadFull(c.bufr, lenBytes); err != nil {
			return 0, nil, err
		}
		payloadLen = int64(binary.BigEndian.Uint16(lenBytes))
	} else if payloadLen == 127 {
		lenBytes := make([]byte, 8)
		if _, err := io.ReadFull(c.bufr, lenBytes); err != nil {
			return 0, nil, err
		}
		payloadLen = int64(binary.BigEndian.Uint64(lenBytes))
	}

	var maskKey []byte
	if isMasked {
		maskKey = make([]byte, 4)
		if _, err := io.ReadFull(c.bufr, maskKey); err != nil {
			return 0, nil, err
		}
	}

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(c.bufr, payload); err != nil {
		return 0, nil, err
	}

	if isMasked {
		for i := 0; i < int(payloadLen); i++ {
			payload[i] ^= maskKey[i%4]
		}
	}

	// Ping frame (0x09) -> reply Pong (0x0A)
	if opcode == 0x09 {
		_ = c.WriteMessage(0x0A, payload)
		return c.ReadMessage()
	}

	// Close frame (0x08)
	if opcode == 0x08 {
		return opcode, nil, io.EOF
	}

	return opcode, payload, nil
}

// WriteMessage sends an unmasked frame to the client.
func (c *WSConn) WriteMessage(opcode byte, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var header []byte
	length := len(payload)

	firstByte := 0x80 | (opcode & 0x0F) // FIN + opcode

	if length <= 125 {
		header = []byte{firstByte, byte(length)}
	} else if length <= 65535 {
		header = make([]byte, 4)
		header[0] = firstByte
		header[1] = 126
		binary.BigEndian.PutUint16(header[2:], uint16(length))
	} else {
		header = make([]byte, 10)
		header[0] = firstByte
		header[1] = 127
		binary.BigEndian.PutUint64(header[2:], uint64(length))
	}

	if _, err := c.conn.Write(header); err != nil {
		return err
	}
	if length > 0 {
		if _, err := c.conn.Write(payload); err != nil {
			return err
		}
	}

	return nil
}

// WriteJSON serializes and sends a JSON text frame.
func (c *WSConn) WriteJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.WriteMessage(0x01, data)
}
