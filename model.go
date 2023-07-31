package rndc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Serializable interface {
	Serialize() ([]byte, error)
}

type CtrlRequest struct {
	Ser   string `json:"_ser"`
	Tim   string `json:"_tim"`
	Exp   string `json:"_exp"`
	Nonce string `json:"_nonce,omitempty"`
}

func (c *CtrlRequest) Serialize() ([]byte, error) {
	var rv bytes.Buffer

	// 序列化 _ser
	if err := serializeStringField(&rv, "_ser", c.Ser); err != nil {
		return nil, err
	}

	// 序列化 _tim
	if err := serializeStringField(&rv, "_tim", c.Tim); err != nil {
		return nil, err
	}

	// 序列化 _exp
	if err := serializeStringField(&rv, "_exp", c.Exp); err != nil {
		return nil, err
	}

	if c.Nonce != "" {
		// 序列化 _nonce
		if err := serializeStringField(&rv, "_nonce", c.Nonce); err != nil {
			return nil, err
		}
	}
	return rv.Bytes(), nil
}

type AuthRequest struct {
	Hash []byte `json:"hsha,omitempty"`
	Md5  []byte `json:"md5,omitempty"`
}

func (s *AuthRequest) Serialize() ([]byte, error) {
	var rv bytes.Buffer

	err := serializeBytesField(&rv, "hsha", s.Hash)
	if err != nil {
		return nil, err
	}
	return rv.Bytes(), nil
}

type CmdRequest struct {
	Auth *AuthRequest `json:"_auth,omitempty"`
	Ctrl *CtrlRequest `json:"_ctrl"`
	Data *DataRequest `json:"_data"`
}

func (s *CmdRequest) Serialize() ([]byte, error) {
	var rv bytes.Buffer

	if s.Auth != nil {
		// 序列化 _auth
		if err := serializeStructField(&rv, "_auth", s.Auth); err != nil {
			return nil, err
		}
	}

	// 序列化 _ctrl
	if err := serializeStructField(&rv, "_ctrl", s.Ctrl); err != nil {
		return nil, err
	}

	if err := serializeStructField(&rv, "_data", s.Data); err != nil {
		return nil, err
	}

	return rv.Bytes(), nil
}

type DataRequest struct {
	Type string `json:"type"`
}

func (s *DataRequest) Serialize() ([]byte, error) {
	var rv bytes.Buffer

	// 序列化 _ctrl
	if err := serializeStringField(&rv, "type", s.Type); err != nil {
		return nil, err
	}

	return rv.Bytes(), nil
}

type CmdResponse struct {
	Auth *AuthResponse `json:"_auth"`
	Ctrl *CtrlResponse `json:"_ctrl"`
	Data *DataResponse `json:"_data"`
}

// 实现 String() 方法以自定义打印格式
func (r *CmdResponse) String() string {
	return fmt.Sprintf("CmdResponse{Auth: %s, Ctrl: %s, Data: %s}", r.Auth, r.Ctrl, r.Data)
}

func (r *CmdResponse) DeSerialize(input []byte) error {
	pos := 0
	for pos < len(input) {
		labelLen := int(input[pos])
		pos++
		label := string(input[pos : pos+labelLen])
		pos += labelLen

		elementType := int(input[pos])
		pos++

		dataLen := int(binary.BigEndian.Uint32(input[pos : pos+4]))
		pos += 4

		data := input[pos : pos+dataLen]
		pos += dataLen

		switch label {
		case "_auth":
			if r.Auth == nil {
				r.Auth = &AuthResponse{}
			}
			err := r.Auth.DeSerialize(data)
			if err != nil {
				return err
			}
		case "_ctrl":
			if r.Ctrl == nil {
				r.Ctrl = &CtrlResponse{}
			}
			err := r.Ctrl.DeSerialize(data)
			if err != nil {
				return err
			}
		case "_data":
			if r.Data == nil {
				r.Data = &DataResponse{}
			}
			err := r.Data.DeSerialize(data)
			if err != nil {
				return err
			}
		default:
			logger.Warn(nil, "未知的字段名: %s, 字段类型: %d", label, elementType)
		}
	}

	return nil
}

type AuthResponse struct {
	Hash []byte `json:"hsha,omitempty"`
	Md5  []byte `json:"md5,omitempty"`
}

// 实现 String() 方法以自定义打印格式
func (r *AuthResponse) String() string {
	return fmt.Sprintf("AuthResponse{Hash: %s, Md5: %s}", r.Hash, r.Md5)
}

func (r *AuthResponse) DeSerialize(input []byte) error {
	pos := 0
	for pos < len(input) {
		labelLen := int(input[pos])
		pos++
		label := string(input[pos : pos+labelLen])
		pos += labelLen

		elementType := int(input[pos])
		pos++

		dataLen := int(binary.BigEndian.Uint32(input[pos : pos+4]))
		pos += 4

		data := input[pos : pos+dataLen]
		pos += dataLen

		switch label {
		case "hsha":
			r.Hash = data
		case "md5":
			r.Md5 = data
		default:
			logger.Warn(nil, "未知的字段名: %s, 字段类型: %d", label, elementType)
		}
	}

	return nil
}

type CtrlResponse struct {
	Ser   string `json:"_ser"`
	Tim   string `json:"_tim"`
	Exp   string `json:"_exp"`
	Rpl   string `json:"_rpl"`
	Nonce string `json:"_nonce,omitempty"`
}

// 实现 String() 方法以自定义打印格式
func (r *CtrlResponse) String() string {
	return fmt.Sprintf("CtrlResponse{Ser: %s, Tim: %s, Exp: %s, Rpl: %s, Nonce: %s}", r.Ser, r.Tim, r.Exp, r.Rpl, r.Nonce)
}

func (r *CtrlResponse) DeSerialize(input []byte) error {
	pos := 0
	for pos < len(input) {
		labelLen := int(input[pos])
		pos++
		label := string(input[pos : pos+labelLen])
		pos += labelLen

		elementType := int(input[pos])
		pos++

		dataLen := int(binary.BigEndian.Uint32(input[pos : pos+4]))
		pos += 4

		data := input[pos : pos+dataLen]
		pos += dataLen

		switch label {
		case "_ser":
			r.Ser = string(data)
		case "_tim":
			r.Tim = string(data)
		case "_exp":
			r.Exp = string(data)
		case "_rpl":
			r.Rpl = string(data)
		case "_nonce":
			r.Nonce = string(data)
		default:
			logger.Warn(nil, "未知的字段名: %s, 字段类型: %d", label, elementType)
		}
	}

	return nil
}

type DataResponse struct {
	Type   string `json:"type"`
	Result string `json:"result"`
	Err    string `json:"err"`
	Text   string `json:"text"`
}

// 实现 String() 方法以自定义打印格式
func (r *DataResponse) String() string {
	return fmt.Sprintf("DataResponse{Type: %s, Result: %s, Err: %s, Text: %s}", r.Type, r.Result, r.Err, r.Text)
}

func (r *DataResponse) DeSerialize(input []byte) error {
	pos := 0
	for pos < len(input) {
		labelLen := int(input[pos])
		pos++
		label := string(input[pos : pos+labelLen])
		pos += labelLen

		elementType := int(input[pos])
		pos++

		dataLen := int(binary.BigEndian.Uint32(input[pos : pos+4]))
		pos += 4

		data := input[pos : pos+dataLen]
		pos += dataLen

		switch label {
		case "type":
			r.Type = string(data)
		case "result":
			r.Result = string(data)
		case "err":
			r.Err = string(data)
		case "text":
			r.Text = string(data)
		default:
			logger.Warn(nil, "未知的字段名: %s, 字段类型: %d", label, elementType)
		}
	}

	return nil
}
