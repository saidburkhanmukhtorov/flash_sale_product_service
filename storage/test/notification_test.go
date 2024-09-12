package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/flash_sale/flash_sale_product_service/storage/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestNotificationRepo(t *testing.T) {
	db := createDBConnection(t)
	defer db.Close(context.Background())

	notificationRepo := postgres.NewNotificationRepo(db)

	t.Run("SendNotification", func(t *testing.T) {
		// Create a test user
		userID := uuid.NewString()
		createUserQuery := `
			INSERT INTO users (id, username, email, password_hash, full_name, date_of_birth, role, created_at, updated_at, deleted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0)
		`
		_, err := db.Exec(context.Background(), createUserQuery, userID, "testuser", "muxtorovsaidburxon764@gmail.com", "passwordhash", "Test User", "2000-01-01", "user")
		assert.NoError(t, err)
		defer deleteUser(t, db, userID) // Cleanup

		// Test message
		message := "Test notification message"

		// Send the notification
		err = notificationRepo.SendNotification(context.Background(), message)
		assert.NoError(t, err)

		fmt.Println("Notification sent successfully (check your email)")
	})
}

func deleteUser(t *testing.T, db *pgx.Conn, userID string) {
	_, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Fatalf("Error deleting user: %v", err)
	}
}
