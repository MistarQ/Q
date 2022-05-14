package core

import (
	"Q/granBlueFantasy/pb"
	"errors"
	"fmt"
	"strconv"
)

func (wm *WorldManager) CreateRoom(bossId uint32, playerId uint32) (*pb.Room, error) {

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	if _, ok := TheWorldManager.Rooms[TheRoomCount]; ok {
		return nil, errors.New("该房间已经存在")
	}

	room := &pb.Room{
		Id:        TheRoomCount,
		Boss:      wm.NewBoss(bossId),
		PlayerMap: make(map[uint32]*pb.Player, 1),
	}
	room.PlayerMap[playerId] = wm.GetPlayerByPID(playerId)
	TheWorldManager.Rooms[TheRoomCount] = room
	TheRoomCount += 1
	for _, rm := range TheWorldManager.Rooms {
		fmt.Println("roomid" + strconv.Itoa(int(rm.Id)))
	}
	return room, nil
}

func (wm *WorldManager) JoinRoom(roomId uint32, playerId uint32) *pb.Room {

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	room := TheWorldManager.Rooms[roomId]
	if _, ok := room.PlayerMap[playerId]; ok {
		return room
	}

	room.PlayerMap[playerId] = wm.GetPlayerByPID(playerId)
	return room
}

func (wm *WorldManager) GetRoom(roomId uint32) *pb.Room {

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	room := TheWorldManager.Rooms[roomId]
	return room
}

func (wm *WorldManager) LeaveRoom(roomId uint32, playerId uint32) {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

}
