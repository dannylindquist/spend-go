package sqlite

import (
	"context"
	"strconv"
	"time"

	sg "github.com/dannylindquist/spend-go"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email string, password string) (*sg.User, error) {
	tx, err := s.db.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	result, err := tx.Exec(`insert into user(email,password) values (?, ?)`, email, password)

	if err != nil {
		return nil, err
	}
	
	rowId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := &sg.User{}
	var (
		createdAt string
		updatedAt string
	)
	err = tx.QueryRowContext(ctx, `
	select
	 id,
	 email,
	 password,
	 createdAt,
	 updatedAt 
	from user where rowId = ?`, rowId).Scan(&user.ID, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	cint, _ := strconv.ParseInt(createdAt, 10, 64)
	tm := time.Unix(cint, 0)
	user.CreatedAt = &tm

	uint, _ := strconv.ParseInt(updatedAt, 10, 64)
	utm := time.Unix(uint, 0)
	user.UpdatedAt = &utm

	tx.Commit()
	return user, nil
}
