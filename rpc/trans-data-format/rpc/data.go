package rpc

import (
	"encoding/binary"
	"io"
	"net"
)

type Session struct {
	Conn net.Conn
}

func NewSession(conn net.Conn) *Session {
	return &Session{conn}
}

func (s *Session) Write(data []byte) error {
	// 定义写数据的格式
	// 4字节头部 + 可变体的长度
	buf := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	copy(buf[4:], data)
	_, err := s.Conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) Read() ([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(s.Conn, header)
	if err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(s.Conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
