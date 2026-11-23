package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yudai0820/go-samurai-store/backend/internal/database"
)

// --- 1. 型定義（struct） ---

// お問い合わせ登録リクエストの構造体
type InquiryCreateRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required"`
}

// お問い合わせ一覧の各項目の構造体
type Inquiry struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// --- 2. ハンドラ定義 ---

// お問い合わせ情報を新規登録する関数
func CreateInquiryHandler(c *gin.Context) {
	var req InquiryCreateRequest
	// HTTPリクエストボディをInquiryCreateRequest構造体にバインドする
	// 必須項目や形式が異なる場合はエラーを返す
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("お問い合わせリクエストのバインドエラー: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力内容が正しくありません"})
		return
	}

	// お問い合わせをデータベースに登録する
	db := database.GetDB()
	insertQuery := "INSERT INTO inquiries (name, email, message) VALUES (?, ?, ?)"
	_, err := db.Exec(insertQuery, req.Name, req.Email, req.Message)
	if err != nil {
		log.Printf("お問い合わせ登録エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	// 登録成功のレスポンスを返す
	c.JSON(http.StatusOK, gin.H{"message": "お問い合わせを送信しました"})
}

// お問い合わせ一覧を取得する関数
func ListInquiriesHandler(c *gin.Context) {
	// 作成日時が新しい順にすべてのお問い合わせを取得する
	db := database.GetDB()
	query := "SELECT id, name, email, message, created_at FROM inquiries ORDER BY created_at DESC"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("お問い合わせ一覧取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}
	// 関数を抜けるときに必ずrowsを閉じる
	defer rows.Close()

	// 0件の場合にnullではなく空配列[]を返すため、空のスライスで初期化する
	inquiries := []Inquiry{}
	for rows.Next() {
		var inquiry Inquiry
		if err := rows.Scan(
			&inquiry.ID,
			&inquiry.Name,
			&inquiry.Email,
			&inquiry.Message,
			&inquiry.CreatedAt,
		); err != nil {
			log.Printf("お問い合わせデータのスキャンエラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
			return
		}
		inquiries = append(inquiries, inquiry)
	}

	// ループ中に発生したエラーを確認する
	if err := rows.Err(); err != nil {
		log.Printf("お問い合わせ一覧の読み込みエラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	// JSONとしてレスポンスを返す
	c.JSON(http.StatusOK, inquiries)
}
