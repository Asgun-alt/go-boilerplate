package repository

import (
	"database/sql"
	"fmt"
	"go-boilerplate/domain/user/constant"
	"go-boilerplate/domain/user/model"
	pkgconstant "go-boilerplate/pkg/constant"
	"go-boilerplate/pkg/session"
)

const createUserSQL = `
    INSERT INTO users (
		user_id, 
		username, 
		password, 
		email)
    VALUES ($1, $2, $3, $4)
    `

func (r *userRepository) Create(
	sess *session.Session,
	user *model.User,
) error {
	_, err := r.db.ExecContext(sess.Ctx, createUserSQL,
		user.UserID,
		user.Username,
		user.Password,
		user.Email)
	if err != nil {
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}

const querySelectUserByIDSQL = `
	SELECT 
		user_id, 
		username, 
		password, 
		email, 
		last_login_at, 
		created_at, 
		updated_at, 
		deleted_at,
		is_email_verified
    FROM users
    WHERE user_id = $1
	AND deleted_at IS NULL `

func (r *userRepository) FindByID(
	sess *session.Session,
	userID string,
) (*model.User, error) {
	var user model.User

	rows, err := r.db.QueryxContext(
		sess.Ctx,
		querySelectUserByIDSQL,
		userID)
	if err != nil {
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.IsEmailVerified)
		if err != nil {
			sess.SetError(pkgconstant.ErrorDatabase, err)
			return nil, err
		}
	} else {
		err = fmt.Errorf("user not found")
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return nil, err
	}

	return &user, nil
}

const querySelectUserByUsernameSQL = `
	SELECT 
		user_id, 
		username, 
		password, 
		email, 
		last_login_at, 
		created_at, 
		updated_at, 
		deleted_at,
		is_email_verified
    FROM users
    WHERE username = $1
	AND deleted_at IS NULL `

func (r *userRepository) FindByUsername(
	sess *session.Session,
	username string,
) (*model.User, error) {
	var user model.User

	rows, err := r.db.QueryxContext(
		sess.Ctx,
		querySelectUserByUsernameSQL,
		username)
	if err != nil {
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.IsEmailVerified)
		if err != nil {
			sess.SetError(pkgconstant.ErrorDatabase, err)
			return nil, err
		}
	} else {
		sess.SetError(pkgconstant.ErrorDatabase, constant.ErrUserNotFound)
		return nil, err
	}

	return &user, nil
}

const querySelectUserByEmailSQL = `
	SELECT 
		user_id, 
		username, 
		password, 
		email, 
		last_login_at, 
		created_at, 
		updated_at, 
		deleted_at,
		is_email_verified
    FROM users
    WHERE email = $1
	AND deleted_at IS NULL `

func (r *userRepository) FindByEmail(
	sess *session.Session,
	email string,
) (*model.User, error) {
	var user model.User

	rows, err := r.db.QueryxContext(
		sess.Ctx,
		querySelectUserByEmailSQL,
		email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.IsEmailVerified)
		if err != nil {
			sess.SetError(pkgconstant.ErrorDatabase, err)
			return nil, err
		}
	} else {
		sess.SetError(pkgconstant.ErrorDatabase, constant.ErrUserNotFound)
		return nil, err
	}

	return &user, nil
}

const queryUpdateUserSQL = `
	UPDATE users
        SET username = $2, 
		password = $3, 
		email = $4, 
		updated_at = now()
    WHERE user_id = $1 
	AND deleted_at IS NULL`

func (r *userRepository) Update(
	sess *session.Session,
	user *model.User,
) error {
	_, err := r.db.ExecContext(sess.Ctx,
		queryUpdateUserSQL,
		user.UserID,
		user.Username,
		user.Password,
		user.Email)
	if err != nil {
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}

const querySoftDeleteUserSQL = `
	UPDATE users
        SET deleted_at = now()
    WHERE user_id = $1 
	AND deleted_at IS NULL`

func (r *userRepository) Delete(
	sess *session.Session,
	userID string,
) error {
	_, err := r.db.ExecContext(sess.Ctx,
		querySoftDeleteUserSQL,
		userID)
	if err != nil {
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}

//nolint:gosec // keep it like this for now
const querySelectUserByEmailAndPasswordSQL = `
	SELECT 
		user_id, 
		username, 
		password, 
		email, 
		last_login_at, 
		created_at, 
		updated_at, 
		deleted_at,
		is_email_verified
    FROM users
    WHERE email = $1
	AND password = $2
	AND deleted_at IS NULL `

func (r *userRepository) FindByEmailAndPassword(
	sess *session.Session,
	email,
	password string,
) (*model.User, error) {
	var user model.User

	rows, err := r.db.QueryxContext(
		sess.Ctx,
		querySelectUserByEmailAndPasswordSQL,
		email,
		password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.IsEmailVerified)
		if err != nil {
			sess.SetError(pkgconstant.ErrorDatabase, err)
			return nil, err
		}
	} else {
		sess.SetError(pkgconstant.ErrorDatabase, constant.ErrUserNotFound)
		return nil, err
	}

	return &user, nil
}

const queryUpdateIsEmailVerified = `
	UPDATE users
        SET is_email_verified = $2, 
		updated_at = now()
    WHERE email = $1 
	AND deleted_at IS NULL`

func (r *userRepository) UpdateIsEmailVerified(
	sess *session.Session,
	email string,
	isVerified bool,
) error {
	_, err := r.db.ExecContext(sess.Ctx,
		queryUpdateIsEmailVerified,
		email,
		isVerified)
	if err != nil {
		sess.SetError(pkgconstant.ErrorDatabase, err)
		return err
	}

	return nil
}
