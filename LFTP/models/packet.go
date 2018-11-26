package models

import (
	"bytes"
	"fmt"
)

type Head struct {
	PkgNum   byte
	ACK      byte
	Status   byte
	Finished byte
	Length   rune
}

type Body struct {
	Data []byte
}

type Packet struct {
	Head
	Body
}

// 初始化一个Packet包
func NewPacket(pkgNum, ack, status, finished byte, data []byte) *Packet {
	packet := &Packet{}
	packet.PkgNum = pkgNum
	packet.ACK = ack
	packet.Status = status
	packet.Finished = finished
	packet.Length = rune(len(data))
	packet.Data = data
	return packet
}

// 将Packet转化为[]byte（封包）
func (packet *Packet) ToBytes() []byte {
	var bytesBuf bytes.Buffer
	bytesBuf.WriteByte(packet.PkgNum)
	bytesBuf.WriteByte(packet.ACK)
	bytesBuf.WriteByte(packet.Status)
	bytesBuf.WriteByte(packet.Finished)
	bytesBuf.WriteRune(packet.Length)
	bytesBuf.Write(packet.Data)
	return bytesBuf.Bytes()
}

// 将[]bytes解析到packet包中（拆包）
func (packet *Packet) FromBytes(buf []byte) {
	// buf = bytes.TrimRight(buf, "\x00")
	bytesBuf := bytes.NewBuffer(buf)
	var err error
	packet.PkgNum, err = bytesBuf.ReadByte()
	checkErr(err)
	packet.ACK, err = bytesBuf.ReadByte()
	checkErr(err)
	packet.Status, err = bytesBuf.ReadByte()
	checkErr(err)
	packet.Finished, err = bytesBuf.ReadByte()
	checkErr(err)
	length, _, err := bytesBuf.ReadRune()
	checkErr(err)
	packet.Length = length
	packet.Data = bytesBuf.Next(int(packet.Length))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
