package postgres

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"github.com/flash_sale/flash_sale_product_service/config"
	"github.com/jackc/pgx/v5"
)

type NotificationRepo struct {
	db *pgx.Conn
}

func NewNotificationRepo(db *pgx.Conn) *NotificationRepo {
	return &NotificationRepo{
		db: db,
	}
}

// UserEmail represents an email
type UserEmail struct {
	Email string
}

// SendNotification retrieves all users emails from the database.
func (r *NotificationRepo) SendNotification(ctx context.Context, message string) error {
	var emails []UserEmail
	query := `
        SELECT email
        FROM users
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get all emails: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user UserEmail
		if err := rows.Scan(&user.Email); err != nil {
			return fmt.Errorf("failed to scan email: %w", err)
		}
		emails = append(emails, user)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating emails: %w", err)
	}
	for _, email := range emails {
		log.Println(email.Email)
		if err := sendEmail(email.Email, message); err != nil {
			return err
		}
	}
	return nil
}

// SendEmail sends an email to the specified recipient.
func sendEmail(recipient string, message string) error {
	cfg := config.Load()
	// Set up authentication
	auth := smtp.PlainAuth("", cfg.EmailSender, cfg.EmailPassword, cfg.EmailHost)

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", cfg.EmailHost, cfg.EmailPort),
		auth,
		cfg.EmailFromAddress,
		[]string{recipient},
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
