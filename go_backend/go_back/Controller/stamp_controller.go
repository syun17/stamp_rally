package controller

import (
	"go_back/go_backend/go_back/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSpots(c *gin.Context) {
	c.JSON(http.StatusOK, model.Spots)
}

func PostStamp(c *gin.Context) {
	var req model.StampRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stamp, status, errMsg := model.AcquireStamp(req.UserID, req.SpotID)

	if errMsg != "" {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusCreated, stamp)
}

func PostStampByNfc(c *gin.Context) {
	var req model.NfcStampRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	spotID := 0
	for _, s := range model.Spots {
		if s.NfcUID == req.NfcUID {
			spotID = s.ID
			break
		}
	}
	if spotID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Spot not found"})
		return
	}
	stamp, status, errMsg := model.AcquireStamp(req.UserID, spotID)
	if errMsg != "" {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusCreated, stamp)
}

func GetUserStamps(c *gin.Context) {
	userID := c.Param("user_id")

	model.Mu.Lock()
	defer model.Mu.Unlock()

	userStamps := model.Stamps[userID]
	if userStamps == nil {
		userStamps = []model.Stamp{}
	}
	c.JSON(http.StatusOK, userStamps)
}

func GetBingo(c *gin.Context) {
	userID := c.Param("user_id")

	model.Mu.Lock()
	defer model.Mu.Unlock()

	stampedSet := make(map[int]bool)
	for _, s := range model.Stamps[userID] {
		stampedSet[s.SpotID] = true
	}

	stampedIDs := make([]int, 0, len(stampedSet))
	for id := range stampedSet {
		stampedIDs = append(stampedIDs, id)
	}

	completedLines := [][]int{}
	for _, line := range model.BingoLines {
		bingo := true
		for _, id := range line {
			if !stampedSet[id] {
				bingo = false
				break
			}
		}
		if bingo {
			completedLines = append(completedLines, line)
		}
	}

	c.JSON(http.StatusOK, model.BingoResult{
		StampedIDs: stampedIDs,
		BingoCount: len(completedLines),
		BingoLines: completedLines,
		IsComplete: len(stampedSet) == len(model.Spots),
	})
}
