package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	config "github.com/vantihovich/go_tasks/tree/master/swagger/configuration"
	handlers "github.com/vantihovich/go_tasks/tree/master/swagger/handlers"
	"github.com/vantihovich/go_tasks/tree/master/swagger/models"
)

type DB struct {
	pool *pgxpool.Pool
	cfg  string
}

func New(cfg config.App) (db DB) {
	db.cfg = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	return db
}

func (db *DB) Open() error {
	pool, err := pgxpool.Connect(context.Background(), db.cfg)
	if err != nil {
		log.WithError(err).Error("Unable to connect to database")
		return err
	}
	log.Info("Successfully connected to DB")
	db.pool = pool
	return nil
}

func (db *DB) FindByLoginAndPwd(ctx context.Context, login, password string) (*models.User, error) {
	var user *models.User = &models.User{}
	err := db.pool.QueryRow(ctx, "SELECT user_id FROM users WHERE login=$1 AND password=$2", login, password).Scan(&user.UserID)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.WithField("User not found", err).Debug("Valid error when login is not found")
			return nil, handlers.ErrNoRows
		}
		log.WithError(err).Error("err executing or parsing the request to DB")
		return nil, err
	}

	log.WithFields(log.Fields{"user_id": user.UserID}).Debug("User found in DB")
	return user, nil
}

func (db *DB) CheckIfLoginExists(ctx context.Context, login string) (bool, error) {
	var user *models.User = &models.User{}
	err := db.pool.QueryRow(ctx, "SELECT user_id FROM users WHERE login=$1 ", login).Scan(&user.UserID)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.WithField("User not found", err).Debug("Valid error when login is not found")
			return false, nil
		}
		log.WithError(err).Error("err executing or parsing the request to DB")
		return false, err
	}
	log.Debug("Provided login found in DB")
	return true, nil
}

func (db *DB) AddNewUser(ctx context.Context, login, password, firstName, lastName, email string, socialMediaLinks []string) error {
	stmnt := "INSERT INTO users (login, password, first_name, last_name, email, social_media_links) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.pool.Exec(ctx, stmnt, login, password, firstName, lastName, email, socialMediaLinks)
	if err != nil {
		log.WithError(err).Error("err executing the DB request to add new user")
		return err
	}

	log.Debug("Successfully added user to DB")
	return nil
}
