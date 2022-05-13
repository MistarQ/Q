package core

import "Q/granBlueFantasy/pb"

func (wm *WorldManager) NewBoss(bossId uint32) *pb.Boss {
	boss := &pb.Boss{
		Id:   bossId,
		Name: "tmpBoss",
		Hp:   100,
		Atk:  10,
		Def:  10,
	}
	return boss
}
