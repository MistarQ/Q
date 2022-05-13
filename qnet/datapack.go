package qnet

import (
	"Q/qiface"
	"Q/qutils"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type DataPack struct {
}

func NewDataPack() qiface.IDataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	return 4 + 4
}

func (d *DataPack) Pack(msg qiface.IMessage) ([]byte, error) {

	dataBuff := bytes.NewBuffer([]byte{})

	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen())
	if err != nil {
		return nil, err
	}

	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	fmt.Println("msgLen: ", msg.GetDataLen(), ", msg id: ", msg.GetMsgId(), ", proto data:", string(msg.GetData()))

	return dataBuff.Bytes(), nil
}

func (d *DataPack) Unpack(binaryData []byte) (qiface.IMessage, error) {

	dataBuff := bytes.NewBuffer(binaryData)

	msg := &Message{}

	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data received")
	}

	return msg, nil
}
