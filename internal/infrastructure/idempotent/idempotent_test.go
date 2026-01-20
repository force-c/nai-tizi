package idempotent

import (
	"context"
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试用的内存数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 创建表
	err = db.Exec(`
		CREATE TABLE idempotent_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			msg_id VARCHAR(255) NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

// TestIdempotent_Execute_FirstTime 测试首次执行消息
func TestIdempotent_Execute_FirstTime(t *testing.T) {
	db := setupTestDB(t)
	idem := New(db)

	msgID := "msg_001"
	executed := false

	err := idem.Execute(context.Background(), msgID, func() error {
		executed = true
		return nil
	})

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if !executed {
		t.Error("expected function to be executed")
	}
}

// TestIdempotent_Execute_Duplicate 测试重复消息
func TestIdempotent_Execute_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	idem := New(db)

	msgID := "msg_002"

	// 首次执行
	err := idem.Execute(context.Background(), msgID, func() error {
		return nil
	})
	if err != nil {
		t.Fatalf("first execution failed: %v", err)
	}

	// 重复执行
	executed := false
	err = idem.Execute(context.Background(), msgID, func() error {
		executed = true
		return nil
	})

	if err != ErrDuplicateMessage {
		t.Errorf("expected ErrDuplicateMessage, got: %v", err)
	}

	if executed {
		t.Error("function should not be executed for duplicate message")
	}
}

// TestIdempotent_Execute_BusinessError 测试业务逻辑返回错误
func TestIdempotent_Execute_BusinessError(t *testing.T) {
	db := setupTestDB(t)
	idem := New(db)

	msgID := "msg_003"
	businessErr := errors.New("business error")

	err := idem.Execute(context.Background(), msgID, func() error {
		return businessErr
	})

	if err != businessErr {
		t.Errorf("expected business error, got: %v", err)
	}

	// 验证业务逻辑失败时事务回滚,消息可以被重试
	// 因为幂等记录没有被插入
	executed := false
	err = idem.Execute(context.Background(), msgID, func() error {
		executed = true
		return nil
	})

	if err != nil {
		t.Errorf("expected successful retry, got: %v", err)
	}

	if !executed {
		t.Error("function should be executed for retry after business error")
	}
}

// TestIdempotent_Execute_DifferentMessages 测试不同的消息
func TestIdempotent_Execute_DifferentMessages(t *testing.T) {
	db := setupTestDB(t)
	idem := New(db)

	messages := []string{"msg_101", "msg_102", "msg_103"}
	executionCount := 0

	for _, msgID := range messages {
		err := idem.Execute(context.Background(), msgID, func() error {
			executionCount++
			return nil
		})
		if err != nil {
			t.Errorf("execution for %s failed: %v", msgID, err)
		}
	}

	if executionCount != len(messages) {
		t.Errorf("expected %d executions, got: %d", len(messages), executionCount)
	}
}

// TestIdempotent_Execute_EmptyMsgID 测试空的 msgID
func TestIdempotent_Execute_EmptyMsgID(t *testing.T) {
	db := setupTestDB(t)
	idem := New(db)

	err := idem.Execute(context.Background(), "", func() error {
		return nil
	})

	// 空 msgID 应该会导致数据库错误或被正常处理
	// 根据实际需求可能需要在代码中添加验证
	if err == nil {
		t.Log("empty msgID was processed without error")
	}
}

// BenchmarkIdempotent_Execute 性能测试
func BenchmarkIdempotent_Execute(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.Exec(`
		CREATE TABLE idempotent_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			msg_id VARCHAR(255) NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)

	idem := New(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msgID := "msg_" + string(rune(i))
		_ = idem.Execute(context.Background(), msgID, func() error {
			return nil
		})
	}
}
