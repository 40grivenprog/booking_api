package db

import (
	"context"
	"database/sql"
)

// WithTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, it's committed.
//
// Example usage:
//
//	err := db.WithTransaction(ctx, database.DB, func(q *db.Queries) error {
//		// Create client
//		client, err := q.CreateClient(ctx, &db.CreateClientParams{
//			FirstName: "John",
//			LastName:  "Doe",
//			ChatID:    sql.NullInt64{Int64: 123456, Valid: true},
//			Locale:    "en",
//		})
//		if err != nil {
//			return err
//		}
//
//		// Create subscription in the same transaction
//		return q.CreateSubscription(ctx, &db.CreateSubscriptionParams{
//			ClientID:       client.ID,
//			ProfessionalID: professionalID,
//		})
//	})
func WithTransaction(ctx context.Context, db *sql.DB, fn func(*Queries) error) error {
	// Begin transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Create queries instance with transaction
	q := New(db).WithTx(tx)

	// Execute the function
	if err := fn(q); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// WithTransactionTx executes a function within an existing transaction.
// This is useful when you need to chain multiple transaction operations.
//
// Example usage:
//
//	tx, err := db.BeginTx(ctx, nil)
//	if err != nil {
//		return err
//	}
//	defer tx.Rollback()
//
//	q := db.New(db).WithTx(tx)
//	if err := db.WithTransactionTx(ctx, tx, func(q *db.Queries) error {
//		// Your operations here
//		return nil
//	}); err != nil {
//		return err
//	}
//
//	return tx.Commit()
func WithTransactionTx(ctx context.Context, tx *sql.Tx, fn func(*Queries) error) error {
	// Create queries instance with transaction
	q := New(nil).WithTx(tx)

	// Execute the function
	return fn(q)
}
