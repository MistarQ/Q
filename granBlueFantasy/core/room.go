package core

import (
	"Q/granBlueFantasy/pb"
	"fmt"
)

func (wm *WorldManager) CreateRoom(bossId uint32, playerId uint32) *pb.Room {

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	room := &pb.Room{
		Id:         TheRoomCount,
		Boss:       wm.NewBoss(bossId),
		PlayerList: make([]*pb.Player, 1),
	}
	player := TheWorldManager.GetPlayerByPID(playerId)
	room.PlayerList = append(room.PlayerList, player)

	TheWorldManager.Rooms[TheRoomCount] = room

	fmt.Println(TheWorldManager.Rooms)
	return room
}

func (wm *WorldManager) JoinRoom(roomId uint32, playerId uint32) *pb.Room {

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	room := TheWorldManager.Rooms[roomId]
	player := TheWorldManager.GetPlayerByPID(playerId)
	room.PlayerList = append(room.PlayerList, player)
	return room
}

func (wm *WorldManager) GetRoom(roomId uint32) *pb.Room {

	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	room := TheWorldManager.Rooms[roomId]
	return room
}
