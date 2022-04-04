package core

import "fmt"

type AOIManager struct {
	MinX int

	MaxX int

	CountX int

	MinY int

	MaxY int

	CountY int

	grids map[int]*Grid
}

func NewAOIManager(minX, maxX, countX, minY, maxY, countY int) *AOIManager {
	aoiMGr := &AOIManager{
		MinX:   minX,
		MaxX:   maxX,
		CountX: countX,
		MinY:   minY,
		MaxY:   maxY,
		CountY: countY,
		grids:  make(map[int]*Grid),
	}

	for y := 0; y < countY; y++ {
		for x := 0; x < countX; x++ {
			gid := y*countX + x
			aoiMGr.grids[gid] = NewGrid(gid,
				aoiMGr.MinX+x*aoiMGr.gridWidth(),
				aoiMGr.MinX+(x+1)*aoiMGr.gridWidth(),
				aoiMGr.MinY+y*aoiMGr.gridLength(),
				aoiMGr.MinY+(y+1)*aoiMGr.gridLength())
		}
	}
	return aoiMGr
}

func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CountX
}

func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CountY
}

func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d, MaxX:%d, countX:%d, minY:%d, maxY:%d, countY:%d\n Grids in AOIManager: ",
		m.MinX, m.MaxX, m.CountX, m.MinY, m.MaxY, m.CountY,
	)

	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

func (m *AOIManager) GetSurroundingGridsByGId(gID int) (grids []*Grid) {

	if _, ok := m.grids[gID]; !ok {
		return
	}

	grids = append(grids, m.grids[gID])

	idx := gID % m.CountX

	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}

	if idx < m.CountX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	gidsX := make([]int, 0, len(grids))

	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	for _, v := range gidsX {
		idy := v / m.CountX
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CountX])
		}
		if idy < m.CountY-1 {
			grids = append(grids, m.grids[v+m.CountX])
		}
	}
	return
}

func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CountX + idx
}

func (m *AOIManager) GetPidsByPos(x, y float32) (plyerIDs []int) {

	gID := m.GetGidByPos(x, y)
	grids := m.GetSurroundingGridsByGId(gID)

	for _, grid := range grids {
		plyerIDs = append(plyerIDs, grid.GetPlayerIDs()...)
		fmt.Printf("-> grid ID : %d, pids :%v", grid.GID, grid.GetPlayerIDs())
	}
	return
}

func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

func (m *AOIManager) GetPidsByGid(gID int) (players []int) {
	return m.grids[gID].GetPlayerIDs()
}

func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(pID)
}
