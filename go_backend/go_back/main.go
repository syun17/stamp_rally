package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ----------------------------------------------------------------
// 構造体（struct）定義
// Goでは struct を使ってデータの型を定義する。
// json:"..." タグは JSON のキー名を指定する（APIレスポンスに使われる）。
// ----------------------------------------------------------------

// Spot はスタンプスポット（訪問場所）を表す
type Spot struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NfcUID      string `json:"nfc_uid"` // NFCタグのUID（リーダーから読み取れる物理ID）
}

// Stamp はスタンプを押した記録を表す
type Stamp struct {
	UserID    string    `json:"user_id"`
	SpotID    int       `json:"spot_id"`
	StampedAt time.Time `json:"stamped_at"` // time.Time はGoの日時型
}

// StampRequest は手動スタンプ取得のリクエストボディ
// binding:"required" は Gin のバリデーション機能で、
// フィールドが空の場合に自動的にエラーを返す
type StampRequest struct {
	UserID string `json:"user_id" binding:"required"`
	SpotID int    `json:"spot_id" binding:"required"`
}

// NfcStampRequest はNFCタグ読み取りによるスタンプ取得のリクエストボディ
type NfcStampRequest struct {
	UserID string `json:"user_id" binding:"required"`
	NfcUID string `json:"nfc_uid" binding:"required"`
}

// BingoResult はビンゴ状況のレスポンス
type BingoResult struct {
	StampedIDs []int   `json:"stamped_ids"` // 取得済みスポットIDの一覧
	BingoCount int     `json:"bingo_count"` // 達成済みビンゴライン数
	BingoLines [][]int `json:"bingo_lines"` // 達成済みラインのスポットID一覧
	IsComplete bool    `json:"is_complete"` // 全スポット制覇フラグ
}

// ----------------------------------------------------------------
// ビンゴライン定義
// 3x3 グリッドの spot_id 配置イメージ:
//
//	[1][2][3]
//	[4][5][6]
//	[7][8][9]
//
// 横3・縦3・斜め2 = 計8ライン
// ----------------------------------------------------------------
var bingoLines = [][]int{
	{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, // 横
	{1, 4, 7}, {2, 5, 8}, {3, 6, 9}, // 縦
	{1, 5, 9}, {3, 5, 7}, // 斜め
}

// ----------------------------------------------------------------
// インメモリデータストア
// DBを使わず、サーバーのメモリ上にデータを保持する。
// サーバーを再起動するとデータはリセットされる。
//
// sync.Mutex（ミューテックス）は複数のリクエストが同時にデータを
// 読み書きしてもデータが壊れないようにするための排他制御。
// Goはgoroutineで並行処理するため、共有データには必須。
// ----------------------------------------------------------------
var (
	mu    sync.Mutex
	spots = []Spot{
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
	// map[string][]Stamp は「ユーザーIDをキー、スタンプ一覧を値」とするマップ
	stamps = make(map[string][]Stamp)
)

// ----------------------------------------------------------------
// main 関数
// Goプログラムのエントリーポイント。サーバーの初期化と起動を行う。
// ----------------------------------------------------------------
func main() {
	// gin.Default() はロギングとパニック時の自動リカバリを持つルーターを生成する
	r := gin.Default()

	// ----------------------------------------------------------------
	// CORS（Cross-Origin Resource Sharing）設定
	// ブラウザはセキュリティ上、異なるオリジン（ドメイン+ポート）への
	// リクエストをデフォルトでブロックする。
	// フロント(localhost:5173) → バックエンド(localhost:8080) は
	// オリジンが異なるため、明示的に許可する必要がある。
	// ----------------------------------------------------------------
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Vite開発サーバーのURL
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
	}))

	// ----------------------------------------------------------------
	// ルーティング定義
	// r.Group("/api") で共通プレフィックス /api をまとめて設定できる。
	// ----------------------------------------------------------------
	api := r.Group("/api")
	{
		api.GET("/spots", getSpots)                // スポット一覧取得
		api.POST("/stamps", postStamp)             // 手動スタンプ（デバッグ用）
		api.POST("/stamps/nfc", postStampByNfc)    // NFCタグ読み取りによるスタンプ
		api.GET("/stamps/:user_id", getUserStamps) // ユーザーのスタンプ一覧
		api.GET("/bingo/:user_id", getBingo)       // ビンゴ状況確認
	}

	// :8080 でサーバーを起動。エラーが起きたらログを出してプロセスを終了する
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

// ----------------------------------------------------------------
// ハンドラー関数
// *gin.Context はリクエスト情報の読み取りとレスポンス送信を担うオブジェクト。
// ----------------------------------------------------------------

// GET /api/spots — スポット一覧をそのまま返す
func getSpots(c *gin.Context) {
	c.JSON(http.StatusOK, spots)
}

// POST /api/stamps — 手動でスタンプを押す（デバッグ・管理用途）
func postStamp(c *gin.Context) {
	var req StampRequest
	// ShouldBindJSON はリクエストボディをJSONとして解析し req に格納する。
	// binding:"required" を付けたフィールドが空ならここでエラーになる。
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 共通ロジックに委譲。複数値の返却はGoの多値返却で行う。
	stamp, status, errMsg := acquireStamp(req.UserID, req.SpotID)
	if errMsg != "" {
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusCreated, stamp)
}

// POST /api/stamps/nfc — NFCタグのUIDを受け取り、対応スポットのスタンプを押す
func postStampByNfc(c *gin.Context) {
	var req NfcStampRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 受け取ったNFC UIDをスポット一覧から線形探索してspot_idを解決する。
	// スポット数が少ないので線形探索で十分。
	spotID := 0
	for _, s := range spots {
		if s.NfcUID == req.NfcUID {
			spotID = s.ID
			break
		}
	}
	// UID が登録されていない場合は 404 を返す
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

// GET /api/stamps/:user_id — ユーザーの取得済みスタンプ一覧を返す
func getUserStamps(c *gin.Context) {
	// :user_id はURLパスパラメータ。c.Param で取得する。
	userID := c.Param("user_id")

	mu.Lock()
	defer mu.Unlock() // defer は関数終了時に必ず実行される。ロック解放を確実にするために使う。

	userStamps := stamps[userID]
	if userStamps == nil {
		// nil をそのままJSONにすると null になるため、空スライスを返す
		userStamps = []Stamp{}
	}
	c.JSON(http.StatusOK, userStamps)
}

// GET /api/bingo/:user_id — ユーザーのビンゴ状況を計算して返す
func getBingo(c *gin.Context) {
	userID := c.Param("user_id")

	mu.Lock()
	defer mu.Unlock()

	// map[int]bool を使ってO(1)で「スタンプ済みか」を判定できるようにする。
	// スライスのまま使うとビンゴ判定のたびに全件走査が必要になり非効率。
	stampedSet := make(map[int]bool)
	for _, s := range stamps[userID] {
		stampedSet[s.SpotID] = true
	}

	// map のキーをスライスに変換（JSONレスポンス用）
	stampedIDs := make([]int, 0, len(stampedSet))
	for id := range stampedSet {
		stampedIDs = append(stampedIDs, id)
	}

	// 各ビンゴラインについて、全スポットがスタンプ済みか確認する
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
		// 取得済み数がスポット総数と一致すれば全制覇
		IsComplete: len(stampedSet) == len(spots),
	})
}

// ----------------------------------------------------------------
// acquireStamp — スタンプ取得の共通ロジック
//
// postStamp と postStampByNfc の両方から呼ばれる。
// 重複処理を避けるために関数として切り出している。
//
// 戻り値: (スタンプデータ, HTTPステータスコード, エラーメッセージ)
// Goは複数の値を返せる。エラー時は stamp が nil になる。
// ----------------------------------------------------------------
func acquireStamp(userID string, spotID int) (*Stamp, int, string) {
	// スポットの存在チェック
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

	// ここからはデータ書き込みを伴うのでロックを取得する
	mu.Lock()
	defer mu.Unlock()

	// 同じユーザーが同じスポットを2回スタンプしようとした場合は 409 Conflict を返す
	for _, s := range stamps[userID] {
		if s.SpotID == spotID {
			return nil, http.StatusConflict, "already stamped"
		}
	}

	// スタンプレコードを作成してメモリに追加する
	// & をつけることでポインタ（参照）を取得する
	stamp := &Stamp{
		UserID:    userID,
		SpotID:    spotID,
		StampedAt: time.Now(), // サーバー側の現在時刻を記録
	}
	// *stamp でポインタが指す実体をスライスに追加する
	stamps[userID] = append(stamps[userID], *stamp)
	return stamp, http.StatusCreated, ""
}
