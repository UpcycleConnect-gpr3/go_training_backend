package payment_handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"go-training-backend/app/middleware/auth_middleware"
	"go-training-backend/app/models/training_models"
	"go-training-backend/app/models/training_reservation_models"
	"go-training-backend/utils/log"
	"go-training-backend/utils/request"
	"go-training-backend/utils/response"
	"go-training-backend/utils/stripe"
)

type createCheckoutDTO struct {
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

// CreateTrainingCheckoutHandler — POST /trainings/{id}/checkout (auth required)
// Creates a one-time Stripe Checkout Session to reserve and pay a training.
// The amount is derived from the training price server side — the client never
// sends a price.
func CreateTrainingCheckoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	userId := auth_middleware.GetUserId(r.Context())
	if userId == "" {
		response.NewErrorMessage(w, response.ErrUserNotFound, http.StatusUnauthorized)
		return
	}

	id := request.Request(r, "id").ConvertToInt(w)
	if id == -1 {
		return
	}

	var dto createCheckoutDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.NewErrorMessage(w, response.ErrJson, http.StatusBadRequest)
		return
	}

	var training training_models.Training
	if err := training.Get([]string{"id", "name", "price"}, "id = ?", id); err != nil {
		response.NewErrorMessage(w, response.ErrTrainingNotFound, http.StatusNotFound)
		return
	}

	// A free training has nothing to pay.
	if training.Price <= 0 {
		response.NewErrorMessage(w, response.ErrInvalidBody, http.StatusBadRequest)
		return
	}

	amountCents := int64(math.Round(training.Price * 100))

	session, err := stripe.CreatePaymentCheckoutSession(
		training.Name,
		amountCents,
		dto.SuccessURL,
		dto.CancelURL,
		userId,
		map[string]string{
			"type":        "training_reservation",
			"training_id": request.Request(r, "id").Value(),
		},
	)
	if err != nil {
		log.Info("stripe create training session: " + err.Error())
		response.NewErrorMessage(w, response.ErrStripe, http.StatusBadGateway)
		return
	}

	training_reservation_models.Upsert(training_reservation_models.TrainingReservation{
		TrainingId:      training.Id,
		UserId:          userId,
		StripeSessionId: session.ID,
		AmountCents:     int(amountCents),
		Status:          "pending",
	})

	response.NewSuccessData(w, map[string]string{"url": session.URL})
}

// GetTrainingPaymentStatusHandler — GET /trainings/payments/session/{id} (auth)
// Reads the session from Stripe, reflects the result in DB, returns the status.
func GetTrainingPaymentStatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Api(r)

	id := request.Request(r, "id").Value()
	if id == "" {
		response.NewErrorMessage(w, response.ErrInvalidValue, http.StatusNotFound)
		return
	}

	session, err := stripe.GetCheckoutSession(id)
	if err != nil {
		log.Info("stripe get training session: " + err.Error())
		response.NewErrorMessage(w, response.ErrStripe, http.StatusBadGateway)
		return
	}

	status := "pending"
	switch session.PaymentStatus {
	case "paid", "no_payment_required":
		status = "paid"
		trainingId := 0
		if _, ok := session.Metadata["training_id"]; ok {
			// best-effort parse; ignore error, reservation already exists
			for _, c := range session.Metadata["training_id"] {
				if c < '0' || c > '9' {
					trainingId = 0
					break
				}
				trainingId = trainingId*10 + int(c-'0')
			}
		}
		training_reservation_models.Upsert(training_reservation_models.TrainingReservation{
			TrainingId:      trainingId,
			UserId:          session.ClientReferenceID,
			StripeSessionId: session.ID,
			AmountCents:     int(session.AmountTotal),
			Status:          "paid",
		})
	case "unpaid":
		status = "unpaid"
	}

	response.NewSuccessData(w, map[string]string{
		"status":         status,
		"customer_email": session.CustomerDetails.Email,
	})
}
