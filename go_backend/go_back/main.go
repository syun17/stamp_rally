package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// スタンプスポット
type Spot struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NfcUID      string `json:"nfc_uid"` // NFCタグのUID
}

// スタンプ（押した記録）
type Stamp struct {
	UserID    string    `json:"user_id"`
	SpotID    int       `json:"spot_id"`
	StampedAt time.Time `json:"stamped_at"`
}

// StampRequest はスタンプ取得リクエスト（手動）
type StampRequest struct {
	UserID string `json:"user_id" binding:"required"`
	SpotID int    `json:"spot_id" binding:"required"`
}

// NfcStampRequest はNFCタグ読み取りによるスタンプ取得リクエスト
type NfcStampRequest struct {
	UserID string `json:"user_id" binding:"required"`
	NfcUID string `json:"nfc_uid" binding:"required"`
}

// BingoResult はビンゴ状況レスポンス
type BingoResult struct {
	StampedIDs []int   `json:"stamped_ids"`
	BingoCount int     `json:"bingo_count"`
	BingoLines [][]int `json:"bingo_lines"`
	IsComplete bool    `json:"is_complete"`
}

// 3x3ビンゴのライン定義（spot_id は 1〜9）
var bingoLines = [][]int{
	{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, // 横
	{1, 4, 7}, {2, 5, 8}, {3, 6, 9}, // 縦
	{1, 5, 9}, {3, 5, 7}, // 斜め
}

// インメモリデータストア
var (
	mu    sync.Mutex
	spots = []Spot{
		// IDは1から9まで、NFC UIDはダミー値
		// いったんは東京の有名スポットを例にしてみる
		{ID: 1, Name: "東京タワー", Description: "東京の象徴的な電波塔", NfcUID: "04:AB:CD:EF:01"},
		{ID: 2, Name: "浅草寺", Description: "東京最古の寺院", NfcUID: "04:AB:CD:EF:02"},
		{ID: 3, Name: "渋谷スクランブル交差点", Description: "世界有数の混雑交差点", NfcUID: "04:AB:CD:EF:03"},
		{ID: 4, Name: "新宿御苑", Description: "広大な都市公園", NfcUID: "04:AB:CD:EF:04"},
		{ID: 5, Name: "上野動物園", Description: "日本最古の動物園", NfcUID: "04:AB:CD:EF:05"},
		{ID: 6, Name: "秋葉原", Description: "電気街・サブカルチャーの聖地", NfcUID: "04:AB:CD:EF:06"},
		{ID: 7, Name: "銀座", Description: "日本屈指の高級商業地", NfcUID: "04:AB:CD:EF:07"},
		{ID: 8, Name: "六本木ヒルズ", Description: "複合エンタメ施設", NfcUID: "04:AB:CD:EF:08"},
		{ID: 9, Name: "東京スカイツリー", Description: "世界最高クラスの電波塔", NfcUID: "04:AB:CD:EF:09"},
	}
	// stamps[userID] = []Stamp
	stamps = make(map[string][]Stamp)
)

func main() {
	r := gin.Default()

	// CORSミドルウェア設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
	}))

	api := r.Group("/api")
	{
		api.GET("/spots", getSpots)
		api.POST("/stamps", postStamp)
		api.POST("/stamps/nfc", postStampByNfc)
		api.GET("/stamps/:user_id", getUserStamps)
		api.GET("/bingo/:user_id", getBingo)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

// GET /api/spots
func getSpots(c *gin.Context) {
	c.JSON(http.StatusOK, spots)
}

// POST /api/stamps  （手動スタンプ / デバッグ用）
func postStamp(c *gin.Context) {
	var req StampRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stamp, status, errMsg := acquireStamp(req.UserID, req.SpotID)
	if errMsg != "" {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusCreated, stamp)
}

// POST /api/stamps/nfc  （NFCタグ読み取りによるスタンプ取得）
func postStampByNfc(c *gin.Context) {
	var req NfcStampRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// NFC UIDからスポットを解決
	spotID := 0
	for _, s := range spots {
		if s.NfcUID == req.NfcUID {
			spotID = s.ID
			break
		}
	}
	if spotID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "NFC tag not registered"})
		return
	}

	stamp, status, errMsg := acquireStamp(req.UserID, spotID)
	if errMsg != "" {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusCreated, stamp)
}

// GET /api/stamps/:user_id
func getUserStamps(c *gin.Context) {
	userID := c.Param("user_id")

	mu.Lock()
	defer mu.Unlock()

	userStamps := stamps[userID]
	if userStamps == nil {
		userStamps = []Stamp{}
	}
	c.JSON(http.StatusOK, userStamps)
}

// GET /api/bingo/:user_id  ビンゴ状況を返す
func getBingo(c *gin.Context) {
	userID := c.Param("user_id")

	mu.Lock()
	defer mu.Unlock()

	stampedSet := make(map[int]bool)
	for _, s := range stamps[userID] {
		stampedSet[s.SpotID] = true
	}

	stampedIDs := make([]int, 0, len(stampedSet))
	for id := range stampedSet {
		stampedIDs = append(stampedIDs, id)
	}

	completedLines := [][]int{}
	for _, line := range bingoLines {
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

	c.JSON(http.StatusOK, BingoResult{
		StampedIDs: stampedIDs,
		BingoCount: len(completedLines),
		BingoLines: completedLines,
		IsComplete: len(stampedSet) == len(spots),
	})
}

// acquireStamp はスタンプ取得の共通ロジック
func acquireStamp(userID string, spotID int) (*Stamp, int, string) {
	validSpot := false
	for _, s := range spots {
		if s.ID == spotID {
			validSpot = true
			break
		}
	}
	if !validSpot {
		return nil, http.StatusNotFound, "spot not found"
	}

	mu.Lock()
	defer mu.Unlock()

	for _, s := range stamps[userID] {
		if s.SpotID == spotID {
			return nil, http.StatusConflict, "already stamped"
		}
	}

	stamp := &Stamp{
		UserID:    userID,
		SpotID:    spotID,
		StampedAt: time.Now(),
	}
	stamps[userID] = append(stamps[userID], *stamp)
	return stamp, http.StatusCreated, ""
}
