package metric_handlers

import (
	"go-upcycle_connect-backend/utils/log"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	log.Api(r)
}
