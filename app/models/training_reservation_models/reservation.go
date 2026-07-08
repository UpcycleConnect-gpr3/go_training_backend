package training_reservation_models

import (
	"go-training-backend/database"
	"go-training-backend/utils/log"
)

const TABLE = "TRAINING_RESERVATIONS"

type TrainingReservation struct {
	Id              int    `db:"id" json:"id"`
	TrainingId      int    `db:"training_id" json:"training_id"`
	UserId          string `db:"user_id" json:"user_id"`
	StripeSessionId string `db:"stripe_session_id" json:"stripe_session_id"`
	AmountCents     int    `db:"amount_cents" json:"amount_cents"`
	Status          string `db:"status" json:"status"`
	CreatedAt       string `db:"created_at" json:"created_at"`
	UpdatedAt       string `db:"updated_at" json:"updated_at"`
}

// Upsert records a reservation keyed by its Stripe session id: the first call
// (checkout creation) inserts it as "pending", the webhook / status poll updates
// it to "paid".
func Upsert(res TrainingReservation) {
	action := "UPSERT " + TABLE
	_, err := database.Training.Exec(
		"INSERT INTO "+TABLE+" (training_id, user_id, stripe_session_id, amount_cents, status) "+
			"VALUES (?, ?, ?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE status = VALUES(status), amount_cents = VALUES(amount_cents)",
		res.TrainingId, res.UserId, res.StripeSessionId, res.AmountCents, res.Status,
	)
	if err != nil {
		log.Database(action, err)
	}
}

// GetUserReservations returns the reservations of a user (most recent first).
func GetUserReservations(userId string) []TrainingReservation {
	action := "SELECT " + TABLE + " (user)"
	reservations := []TrainingReservation{}
	err := database.Training.Select(
		&reservations,
		"SELECT id, training_id, user_id, stripe_session_id, amount_cents, status, created_at, updated_at "+
			"FROM "+TABLE+" WHERE user_id = ? ORDER BY created_at DESC",
		userId,
	)
	if err != nil {
		log.Database(action, err)
	}
	return reservations
}
