package queries

import (
	"database/sql"
	"fmt"

	"coding-winner/internal/models"
)

type UserDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
}

// CreateUser creates a new user
func CreateUser(db UserDB, user *models.User) error {
	query := `
		INSERT INTO users (discord_id, atcoder_username)
		VALUES ($1, $2)
		ON CONFLICT (discord_id) DO UPDATE
		SET atcoder_username = EXCLUDED.atcoder_username,
		    updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.Exec(query, user.DiscordID, user.AtCoderUsername)
	return err
}

// GetUser retrieves a user by Discord ID
func GetUser(db UserDB, discordID string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE discord_id = $1`
	err := db.Get(&user, query, discordID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByAtCoderUsername retrieves a user by AtCoder username
func GetUserByAtCoderUsername(db UserDB, username string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE atcoder_username = $1`
	err := db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all registered users
func GetAllUsers(db UserDB) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT * FROM users ORDER BY created_at DESC`
	err := db.Select(&users, query)
	return users, err
}

// GetServerUsers retrieves all users in a specific server
func GetServerUsers(db UserDB, serverID string) ([]*models.User, error) {
	// Note: This requires tracking which users are in which server
	// For now, we'll return all users and filter in Discord Bot logic
	return GetAllUsers(db)
}

// DeleteUser deletes a user
func DeleteUser(db UserDB, discordID string) error {
	query := `DELETE FROM users WHERE discord_id = $1`
	result, err := db.Exec(query, discordID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// UserExists checks if a user exists
func UserExists(db UserDB, discordID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE discord_id = $1)`
	err := db.Get(&exists, query, discordID)
	return exists, err
}
