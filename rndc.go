package rndc

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	algosMap = map[string]int{
		"md5":    157,
		"sha1":   161,
		"sha224": 162,
		"sha256": 163,
		"sha384": 164,
		"sha512": 165,
	}
)

type RnDC struct {
	host   string
	algo   string
	secret []byte
	ser    int
	nonce  string
	conn   net.Conn
}

// NewRNDCClient create rndc client
// host - (ip, port) tuple
// algo - HMAC algorithm: one of md5, sha1, sha224, sha256, sha384, sha512
//
//	(with optional prefix 'hmac-')
//
// secret - HMAC secret, base64 encoded
func NewRNDCClient(host, algo, secret string) (*RnDC, error) {
	rand.Seed(time.Now().UnixNano())
	c := &RnDC{
		host:   host,
		secret: []byte(secret),
		ser:    rand.Intn(1 << 24),
	}

	algo = strings.ToLower(algo)
	if strings.HasPrefix(algo, "hmac-") {
		algo = algo[5:]
	}

	switch algo {
	case "md5":
		c.algo = algo
	case "sha1":
		c.algo = algo
	case "sha224":
		c.algo = algo
	case "sha256":
		c.algo = algo
	case "sha384":
		c.algo = algo
	case "sha512":
		c.algo = algo
	default:
		return nil, errors.New("unsupported HMAC algorithm")
	}

	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	c.secret = decodedSecret
	if err := c.connectLogin(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *RnDC) connectLogin() error {
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return err
	}
	c.conn = conn
	resp, err := c.command("null")
	if err != nil {
		return err
	}
	if resp.Data.Result != "0" {
		return errors.New(resp.Data.Err)
	}
	c.nonce = resp.Ctrl.Nonce
	return nil
}

func (c *RnDC) command(cmd string) (*CmdResponse, error) {
	msg, err := c.prepMessage(cmd)
	if err != nil {
		return nil, err
	}
	sent, err := c.conn.Write(msg)
	if err != nil {
		return nil, err
	}
	if sent != len(msg) {
		return nil, fmt.Errorf("消息发送失败")
	}

	header := make([]byte, 8)
	_, err = io.ReadFull(c.conn, header)
	if err != nil {
		return nil, errors.New("无法读取返回头信息")
	}

	length, version, err := unpackHeader(header)
	if err != nil {
		return nil, err
	}

	if version != ProtocolVersion {
		return nil, fmt.Errorf("RnDC协议版本错误, 服务器支持版本: %d, 客户端版本: %d", version, ProtocolVersion)
	}

	message := make([]byte, length-4)
	_, err = io.ReadFull(c.conn, message)
	if err != nil {
		return nil, fmt.Errorf("无法读取返回内容: %s", err.Error())
	}

	resp := &CmdResponse{}
	if err = resp.DeSerialize(message); err != nil {
		return nil, fmt.Errorf("返回内容反序列化失败: %s", err.Error())
	}
	return resp, nil
}

func (c *RnDC) Call(cmd string) (*CmdResponse, error) {
	return c.command(cmd)
}

func (c *RnDC) prepMessage(cmd string) ([]byte, error) {
	// -------------- 1. 构造发送的消息结构 -----------------
	c.ser++
	now := time.Now().Unix()

	preMsg := CmdRequest{
		Ctrl: &CtrlRequest{
			Ser: strconv.FormatInt(int64(c.ser), 10),
			Tim: strconv.FormatInt(now, 10),
			Exp: strconv.FormatInt(now+60, 10),
		},
		Data: &DataRequest{Type: cmd},
	}
	if c.nonce != "" {
		preMsg.Ctrl.Nonce = c.nonce
	}

	// -------------- 2. 计算发送消息的摘要并加密 -----------------
	msg, err := preMsg.Serialize()
	if err != nil {
		return nil, err
	}

	// 对 msg 进行校验加密
	Hash := c.calculateHMAC(msg)
	bHash := make([]byte, base64.StdEncoding.EncodedLen(len(Hash)))
	base64.StdEncoding.Encode(bHash, Hash)
	authMsg, err := c.calculateAuthData(bHash)
	if err != nil {
		return nil, err
	}
	preMsg.Auth = &AuthRequest{Hash: authMsg}

	msg, err = preMsg.Serialize()
	if err != nil {
		return nil, err
	}

	return packHeader(msg)
}

func (c *RnDC) calculateHMAC(data []byte) []byte {
	var h func() hash.Hash

	switch c.algo {
	case "md5":
		h = md5.New
	case "sha1":
		h = sha1.New
	//case "sha224":
	//	h = sha224.New
	case "sha256":
		h = sha256.New
	case "sha384":
		h = sha512.New384
	case "sha512":
		h = sha512.New
	default:
		panic(fmt.Errorf("不支持的算法类型: %s", c.algo))
	}

	ha := hmac.New(h, c.secret)
	ha.Write(data)
	return ha.Sum(nil)
}

func (c *RnDC) calculateAuthData(bHash []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(byte(algosMap[c.algo]))
	bv, _ := paddingBHash(bHash)
	err := binary.Write(&buf, binary.LittleEndian, bv)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}
