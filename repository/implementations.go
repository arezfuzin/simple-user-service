package repository

import (
	"context"
	"strconv"
)

func (r *Repository) Registration(ctx context.Context, input RegistrationInput) (output RegistrationOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO users (full_name, phone_number, password_hash) VALUES ($1, $2, $3) RETURNING id", input.FullName, input.PhoneNumber, input.Password).Scan(&output.UserId)
	if err != nil {
		return
	}
	return
}

func (r *Repository) Login(ctx context.Context, input LoginInput) (output LoginOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT * FROM users WHERE phone_number = $1", input.PhoneNumber).Scan(&output.Id, &output.FullName, &output.PhoneNumber, &output.Password)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetMyProfile(ctx context.Context, input GetMyProfileInput) (output GetMyProfileOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT full_name, phone_number FROM users WHERE id = $1", input.Id).Scan(&output.FullName, &output.PhoneNumber)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateMyProfile(ctx context.Context, input UpdateMyProfileInput) (output UpdateMyProfileOutput, err error) {
	query := "UPDATE users SET"
	params := []interface{}{}

	if input.FullName != nil {
		query += " full_name = $" + strconv.Itoa(len(params)+1) + ","
		params = append(params, input.FullName)
	}

	if input.PhoneNumber != nil {
		query += " phone_number = $" + strconv.Itoa(len(params)+1) + ","
		params = append(params, input.PhoneNumber)
	}

	query = query[:len(query)-1] // Remove the trailing comma

	query += " WHERE id = $" + strconv.Itoa(len(params)+1) + " RETURNING full_name, phone_number"
	params = append(params, input.Id)

	err = r.Db.QueryRowContext(ctx, query, params...).Scan(&output.FullName, &output.PhoneNumber)
	if err != nil {
		return
	}
	return
}
