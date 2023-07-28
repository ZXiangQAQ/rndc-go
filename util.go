package rndc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// paddingBHash 填充 bHash 到 88 字节
func paddingBHash(bHash []byte) ([]byte, error) {
	padLength := 88 - len(bHash)
	if padLength < 0 {
		return nil, fmt.Errorf("bHash 长度必须小于等于88字节")
	}
	bHashBytes := bHash
	bHashBytes = append(bHashBytes, bytes.Repeat([]byte{0x00}, padLength)...)
	return bHashBytes, nil
}

// unpackHeader 使用binary包的Unpack方法解包数据
func unpackHeader(header []byte) (uint32, uint32, error) {
	var length, version uint32
	err := binary.Read(bytes.NewReader(header), binary.BigEndian, &length)
	if err != nil {
		return 0, 0, fmt.Errorf("解包失败, 无法解析length: %s", err.Error())
	}
	err = binary.Read(bytes.NewReader(header[4:]), binary.BigEndian, &version)
	if err != nil {
		return 0, 0, fmt.Errorf("解包失败, 无法解析version: %s", err.Error())
	}
	return length, version, nil
}

func packHeader(msg []byte) ([]byte, error) {
	var rv bytes.Buffer
	if err := binary.Write(&rv, binary.BigEndian, uint32(len(msg)+4)); err != nil {
		return nil, fmt.Errorf("包长:%d, 封包失败, %s", len(msg), err.Error())
	}
	if err := binary.Write(&rv, binary.BigEndian, uint32(ProtocolVersion)); err != nil {
		return nil, fmt.Errorf("RnDc协议版本:%d, 封包失败, %s", ProtocolVersion, err.Error())
	}
	rv.Write(msg)
	return rv.Bytes(), nil
}
