package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID int

	MinX int

	MaxX int

	MinY int

	MaxY int

	playerIDs map[int]bool

	PIDLock sync.RWMutex
}

func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

func (g *Grid) Add(playerID int) {
	g.PIDLock.Lock()
	defer g.PIDLock.Unlock()

	g.playerIDs[playerID] = true
}

func (g *Grid) Remove(playerID int) {
	g.PIDLock.Lock()
	defer g.PIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.PIDLock.RLock()
	defer g.PIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
