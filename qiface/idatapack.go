package qiface

type IDataPack interface {
	GetHeadLen() uint32 // get head len

	Pack(msg IMessage) ([]byte, error) // pack

	Unpack([]byte) (IMessage, error) // unpack
}
