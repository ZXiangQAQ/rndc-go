package rndc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func serializeStringField(rv *bytes.Buffer, key string, val string) error {
	// 序列化 key
	// 一个字节存储 key 的长度
	rv.WriteByte(byte(len(key)))
	// 后面字节写入 key 本身的值
	rv.WriteString(key)

	// 序列化 val
	// 一个字节可能是标志位
	if err := binary.Write(rv, binary.BigEndian, uint8(1)); err != nil {
		return fmt.Errorf("序列化失败, %s字段值的标识位写入失败: %s\n", key, err.Error())
	}
	// 4个字节的长度存储 val 的长度, 大端序编码(右对齐)
	if err := binary.Write(rv, binary.BigEndian, uint32(len(val))); err != nil {
		return fmt.Errorf("序列化失败, %s字段的值%s写入失败: %s\n", key, val, err.Error())
	}
	rv.WriteString(val)
	return nil
}

func serializeBytesField(rv *bytes.Buffer, key string, val []byte) error {
	// 序列化 key
	// 一个字节存储 key 的长度
	rv.WriteByte(byte(len(key)))
	// 后面字节写入 key 本身的值
	rv.WriteString(key)

	// 序列化 val
	// 一个字节可能是标志位
	if err := binary.Write(rv, binary.BigEndian, uint8(1)); err != nil {
		return fmt.Errorf("序列化失败, %s字段值的标识位写入失败: %s\n", key, err.Error())
	}
	// 4个字节的长度存储 val 的长度, 大端序编码(右对齐)
	if err := binary.Write(rv, binary.BigEndian, uint32(len(val))); err != nil {
		return fmt.Errorf("序列化失败, %s字段的值%s写入失败: %s\n", key, val, err.Error())
	}
	rv.Write(val)
	return nil
}

func serializeStructField(rv *bytes.Buffer, key string, val Serializable) error {
	// 1. 序列化 key
	// 1.1 一个字节存储 key 的长度
	rv.WriteByte(byte(len(key)))
	// 1.2 后面字节写入 key 本身的值
	rv.WriteString(key)

	// 2. 序列化 val
	// 2.1 一个字节可能是标志位
	if err := binary.Write(rv, binary.BigEndian, uint8(2)); err != nil {
		return fmt.Errorf("序列化失败, %s字段值的标识位写入失败: %s\n", key, err.Error())
	}
	// 2.2 序列化 val 的值
	// 2.2.1 先拿到 val 的 bytes
	valBytes, err := val.Serialize()
	if err != nil {
		return fmt.Errorf("序列化失败, %s字段的值%s写入失败: %s\n", key, val, err.Error())
	}
	// 2.2.2 存储 val 的长度, 4个字节的长度存储 val 的长度, 大端序编码(右对齐)
	if err := binary.Write(rv, binary.BigEndian, uint32(len(valBytes))); err != nil {
		return fmt.Errorf("序列化失败, %s字段的值%s写入失败: %s\n", key, valBytes, err.Error())
	}
	// 2.2.1.3 存储 val 的值
	rv.Write(valBytes)
	return nil
}
