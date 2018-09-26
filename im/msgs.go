package im

import (
	"bytes"
	"encoding/gob"
	"errors"
	"goim/helpers"
	"io"
	"log"
)

const (
	MessageType_Login1 = iota
	MessageType_Login2
	MessageType_Login3
	MessageType_Info
)

type MessageBase struct {
	CheckSum uint32
}

type InfoMsg struct {
	Text string
}

type LoginMsg1 struct {
	MessageBase
	PublicKey string
}

func UnMarshal(b []byte, msg interface{}) error {
	var buf = bytes.Buffer{}
	buf.Write(b)
	// Create a decoder and receive a value.
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(msg)
	if err != nil {
		log.Fatal("decode:", err)
		return err
	}
	return nil
}

type LoginMsg2 struct {
	MessageBase
	//encrypted message with
	//RSA public key
	EncryptedText []byte
}

type LoginMsg3 struct {
	MessageBase
	//DecryptedText with private
	//key
	DecryptedText []byte
}

func Marshal(o interface{}) ([]byte, error) {
	var buf bytes.Buffer
	// Create an encoder and send a value.
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(o)
	if err != nil {
		log.Fatal("encode:", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func WriteInfoMessage(w io.Writer, msg string) error {
	m := &InfoMsg{
		Text: msg,
	}
	bs, err := Marshal(m)
	if err != nil {
		return err
	}
	l, err2 := helpers.WriteMessage(w, MessageType_Info, bs)
	if err2 != nil {
		return err2
	}
	if l != len(bs) {
		return errors.New("write failed!")
	}

	return nil
}

func GetInfoMessage(r io.Reader) (*InfoMsg, error) {
	t, buf, e := helpers.ReadMessage(r)
	if e != nil {
		return nil, e
	}
	if t != MessageType_Info {
		return nil, errors.New("Wrong message sequence!")
	}
	var m1 InfoMsg
	err := UnMarshal(buf, &m1)
	if err != nil {
		return nil, err
	}
	return &m1, nil
}
