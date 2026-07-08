package server

import (
	"go-training-backend/app/handlers/category_handlers"
	"go-training-backend/app/handlers/content_schedule_handlers"
	"go-training-backend/app/handlers/curricula_handlers"
	"go-training-backend/app/handlers/image_handlers"
	"go-training-backend/app/handlers/metric_handlers"
	"go-training-backend/app/handlers/training_content_handlers"
	"go-training-backend/app/handlers/training_handlers"
	"go-training-backend/app/middleware/auth_middleware"
	"go-training-backend/app/middleware/ratelimit_middleware"
	"go-training-backend/config"
	"go-training-backend/database"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	log "github.com/thedataflows/go-lib-log"
)

func initialize() {

	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithLogFormat(log.LOG_FORMAT_JSON).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Config Initialization
	config.InitDatabase()

	err = database.Training.Ping()

	if err != nil {
		logger.Fatal().Err(err).Msg("(DATABASE)")
	}
}

func Start() {

	initialize()

	//limiterLow := ratelimit_middleware.NewRateLimiter(10, 1*time.Minute)
	limiterMedium := ratelimit_middleware.NewRateLimiter(30, 1*time.Minute)
	limiterHigh := ratelimit_middleware.NewRateLimiter(60, 1*time.Minute)

	//containerBackoffice := source_middleware.Container("go-backoffice-backend")

	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	http.HandleFunc("GET /health/{$}", metric_handlers.Health)

	http.HandleFunc("GET /trainings/{$}", limiterHigh.RateLimit(training_handlers.GetTrainingsHandler))
	http.HandleFunc("GET /trainings/{id}/{$}", limiterHigh.RateLimit(training_handlers.GetTrainingHandler))
	http.HandleFunc("POST /trainings/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_handlers.CreateTrainingHandler)))
	http.HandleFunc("PUT /trainings/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_handlers.UpdateTrainingHandler)))
	http.HandleFunc("DELETE /trainings/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_handlers.DeleteTrainingHandler)))
	http.HandleFunc("GET /trainings/{id}/curricula/{$}", limiterHigh.RateLimit(training_handlers.GetTrainingCurriculaHandler))
	http.HandleFunc("POST /trainings/{id}/curricula/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_handlers.LinkTrainingCurriculumHandler)))
	http.HandleFunc("DELETE /trainings/{id}/curricula/{curriculum_id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_handlers.UnlinkTrainingCurriculumHandler)))
	http.HandleFunc("GET /trainings/{id}/content/{$}", limiterHigh.RateLimit(training_handlers.GetTrainingContentHandler))
	http.HandleFunc("POST /trainings/{id}/content/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_handlers.LinkTrainingContentHandler)))
	http.HandleFunc("DELETE /trainings/{id}/content/{content_id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_handlers.UnlinkTrainingContentHandler)))
	http.HandleFunc("GET /trainings/{id}/schedules/{$}", limiterHigh.RateLimit(training_handlers.GetTrainingSchedulesHandler))

	http.HandleFunc("GET /curricula/{$}", limiterHigh.RateLimit(curricula_handlers.GetCurriculaHandler))
	http.HandleFunc("GET /curricula/{id}/{$}", limiterHigh.RateLimit(curricula_handlers.GetCurriculumHandler))
	http.HandleFunc("POST /curricula/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(curricula_handlers.CreateCurriculumHandler)))
	http.HandleFunc("PUT /curricula/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(curricula_handlers.UpdateCurriculumHandler)))
	http.HandleFunc("DELETE /curricula/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(curricula_handlers.DeleteCurriculumHandler)))

	http.HandleFunc("GET /schedules/{$}", limiterHigh.RateLimit(content_schedule_handlers.GetContentSchedulesHandler))
	http.HandleFunc("GET /schedules/{id}/{$}", limiterHigh.RateLimit(content_schedule_handlers.GetContentScheduleHandler))
	http.HandleFunc("POST /schedules/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(content_schedule_handlers.CreateContentScheduleHandler)))
	http.HandleFunc("PUT /schedules/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(content_schedule_handlers.UpdateContentScheduleHandler)))
	http.HandleFunc("DELETE /schedules/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(content_schedule_handlers.DeleteContentScheduleHandler)))

	http.HandleFunc("GET /training-content/{$}", limiterHigh.RateLimit(training_content_handlers.GetTrainingContentsHandler))
	http.HandleFunc("GET /training-content/{id}/{$}", limiterHigh.RateLimit(training_content_handlers.GetTrainingContentHandler))
	http.HandleFunc("POST /training-content/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_content_handlers.CreateTrainingContentHandler)))
	http.HandleFunc("PUT /training-content/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_content_handlers.UpdateTrainingContentHandler)))
	http.HandleFunc("DELETE /training-content/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(training_content_handlers.DeleteTrainingContentHandler)))

	http.HandleFunc("GET /images/{$}", limiterHigh.RateLimit(image_handlers.GetImagesHandler))
	http.HandleFunc("GET /images/{id}/{$}", limiterHigh.RateLimit(image_handlers.GetImageHandler))
	http.HandleFunc("POST /images/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(image_handlers.CreateImageHandler)))
	http.HandleFunc("PUT /images/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(image_handlers.UpdateImageHandler)))
	http.HandleFunc("DELETE /images/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(image_handlers.DeleteImageHandler)))

	http.HandleFunc("GET /categories/{$}", limiterHigh.RateLimit(category_handlers.GetCategoriesHandler))
	http.HandleFunc("GET /categories/{id}/{$}", limiterHigh.RateLimit(category_handlers.GetCategoryHandler))
	http.HandleFunc("POST /categories/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(category_handlers.CreateCategoryHandler)))
	http.HandleFunc("PUT /categories/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(category_handlers.UpdateCategoryHandler)))
	http.HandleFunc("DELETE /categories/{id}/{$}", limiterMedium.RateLimit(auth_middleware.IsAuth(category_handlers.DeleteCategoryHandler)))

	logger.Info().Msg("Listening at http://localhost:" + os.Getenv("APP_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), corsMiddleware(http.DefaultServeMux))
	if err != nil {
		return
	}
}

// corsMiddleware enables cross-origin requests from the local dev frontends
// (Vite on a different port). Reflects the request Origin and answers the
// preflight OPTIONS so browser calls are not blocked.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
