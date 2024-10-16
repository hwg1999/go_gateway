package codec

import (
	"encoding/binary"
	"errors"
	"io"
)

const Msg_Header = "12345678"

func Encode(bytesBuf io.Writer, content string) error {
	// msg_header + content_len + content
	// 8 + 4 + content_len
	err := binary.Write(bytesBuf, binary.BigEndian, []byte(Msg_Header))
	if err != nil {
		return err
	}

	cLen := int32(len([]byte(content)))
	err = binary.Write(bytesBuf, binary.BigEndian, cLen)
	if err != nil {
		return err
	}

	err = binary.Write(bytesBuf, binary.BigEndian, []byte(content))
	if err != nil {
		return err
	}

	return nil
}

func Decode(bytesBuf io.Reader) (bodyBuf []byte, err error) {
	MagicBuf := make([]byte, len(Msg_Header))
	_, err = io.ReadFull(bytesBuf, MagicBuf)
	if err != nil {
		return nil, err
	}

	if string(MagicBuf) != Msg_Header {
		return nil, errors.New("msg_header error")
	}

	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(bytesBuf, lenBuf)
	if err != nil {
		return nil, err
	}

	len := binary.BigEndian.Uint32(lenBuf)
	bodyBuf = make([]byte, len)
	_, err = io.ReadFull(bytesBuf, bodyBuf)
	if err != nil {
		return nil, err
	}

	return bodyBuf, err
}
