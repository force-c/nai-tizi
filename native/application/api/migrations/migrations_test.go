package migrations

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestUpFromEmptyDatabase(t *testing.T) {
	dsn := os.Getenv("LIGHTNING_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("LIGHTNING_TEST_POSTGRES_DSN is not set")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := Up(db); err != nil {
		t.Fatalf("Up() error = %v", err)
	}

	for _, table := range []string{"s_user", "s_role", "s_api_permission", "biz_attachment"} {
		var exists bool
		if err := db.QueryRow(`SELECT to_regclass('public.' || $1) IS NOT NULL`, table).Scan(&exists); err != nil {
			t.Fatal(err)
		}
		if !exists {
			t.Errorf("table %s was not created", table)
		}
	}

	var grantTypes string
	if err := db.QueryRow(`SELECT grant_type FROM s_auth_client WHERE client_id = 'web-admin'`).Scan(&grantTypes); err != nil {
		t.Fatal(err)
	}
	if grantTypes != "password,email,sms,wechat" {
		t.Fatalf("web-admin grant types = %q", grantTypes)
	}

	for _, index := range []string{"idx_s_user_email", "idx_s_user_phonenumber", "idx_s_user_open_id", "idx_s_user_union_id", "idx_user_role"} {
		var unique bool
		if err := db.QueryRow(`SELECT indisunique FROM pg_index WHERE indexrelid = to_regclass($1)`, index).Scan(&unique); err != nil {
			t.Fatal(err)
		}
		if !unique {
			t.Errorf("index %s is not unique", index)
		}
	}

	var configIndexUnique bool
	if err := db.QueryRow(`SELECT indisunique FROM pg_index WHERE indexrelid = to_regclass('idx_s_config_code')`).Scan(&configIndexUnique); err != nil {
		t.Fatal(err)
	}
	if !configIndexUnique {
		t.Fatal("idx_s_config_code is not unique")
	}
	var runtimeConfigCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM s_config WHERE code IN ('integration.wechat', 'integration.sms', 'integration.email', 'auth.captcha', 'scheduler')`).Scan(&runtimeConfigCount); err != nil {
		t.Fatal(err)
	}
	if runtimeConfigCount != 5 {
		t.Fatalf("runtime configuration count = %d, want 5", runtimeConfigCount)
	}
	var hasVersion bool
	if err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = 'public' AND table_name = 's_config' AND column_name = 'version')`).Scan(&hasVersion); err != nil {
		t.Fatal(err)
	}
	if hasVersion {
		t.Fatal("s_config unexpectedly has a version column")
	}
}
