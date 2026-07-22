package migrations

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestUpFromEmptyDatabase(t *testing.T) {
	dsn := os.Getenv("QUICK_ADMIN_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("QUICK_ADMIN_TEST_POSTGRES_DSN is not set")
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
}
