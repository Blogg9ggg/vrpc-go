package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}

// ReadHeader 将 header 从 conn 读取到 h 变量
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// ReadBody 将 body 从 conn 读取到 body 变量
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// Write 将 header 和 body 写入到 buf
func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		err = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()

	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header: ", err)
		return
	}

	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body: ", err)
		return
	}

	return
}
