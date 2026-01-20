package idempotent

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	// ErrDuplicateMessage 重复消息错误
	ErrDuplicateMessage = errors.New("duplicate message")
)

// Idempotent 幂等性处理器
type Idempotent struct {
	db *gorm.DB
}

// idempotentRecord 幂等性记录
type idempotentRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	MsgID     string    `gorm:"column:msg_id;uniqueIndex;size:255;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName 指定表名
func (idempotentRecord) TableName() string {
	return "idempotent_records"
}

// New 创建幂等性处理器
func New(db *gorm.DB) *Idempotent {
	return &Idempotent{db: db}
}

// Execute 执行幂等性检查
// 如果 msgID 已存在（重复消息），返回 ErrDuplicateMessage
// 如果 msgID 不存在，插入记录并执行函数 fn
func (i *Idempotent) Execute(ctx context.Context, msgID string, fn func() error) error {
	// 开启事务
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 尝试插入记录
		record := &idempotentRecord{
			MsgID: msgID,
		}

		// 尝试创建记录，如果已存在会因为唯一索引失败
		err := tx.Create(record).Error
		if err != nil {
			// 检查是否是唯一约束冲突
			if errors.Is(err, gorm.ErrDuplicatedKey) || isUniqueViolation(err) {
				return ErrDuplicateMessage
			}
			return err
		}

		// 执行业务函数
		return fn()
	})
}

// isUniqueViolation 检查是否是唯一约束冲突错误
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// PostgreSQL: "duplicate key value violates unique constraint"
	// MySQL: "Duplicate entry"
	// SQLite: "UNIQUE constraint failed"
	return contains(errStr, "duplicate key") ||
		contains(errStr, "Duplicate entry") ||
		contains(errStr, "UNIQUE constraint failed")
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexIgnoreCase(s, substr) >= 0
}

// indexIgnoreCase 不区分大小写查找子串
func indexIgnoreCase(s, substr string) int {
	sLower := toLower(s)
	substrLower := toLower(substr)
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return i
		}
	}
	return -1
}

// toLower 转换为小写
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
