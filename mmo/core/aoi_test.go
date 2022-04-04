package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	fmt.Println(aoiMgr)

}

func TestAOIManager_GetSurroundingGridsByGId(t *testing.T) {

	aoiMgr := NewAOIManager(0, 250, 5, 0, 300, 6)

	for gid, _ := range aoiMgr.grids {
		//grids := aoiMgr.GetSurroundGridsByGid(gid)
		grids := aoiMgr.GetSurroundingGridsByGId(gid)
		fmt.Println("gid: ", gid, "grids len = ", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Printf("surrounding grid IDs are %v\n", gIDs)
	}

}
