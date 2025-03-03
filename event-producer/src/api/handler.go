package api

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"rest-server/schema"
	"time"

	"rest-server/logger"
	"rest-server/redis"
	"rest-server/utils"
)

const (
	streamCreateEntryKey = "CreateStream"
	streamEntryKey       = "EventStream"
)

type handler struct {
	logger      *logger.Log
	ctx         context.Context
	redisClient *redis.RedisClient
}

func (h *handler) defaultEndpoint(w http.ResponseWriter, _ *http.Request) {
	status := http.StatusNotFound
	utils.WriteResponse(w, status, utils.ErrorResponse{Status: status, Title: "Endpoint Not Found", Description: "Endpoint Not Found"})
}

func (h *handler) eventEndpoint(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK

	if r.Method != http.MethodPost {
		status = http.StatusMethodNotAllowed
		utils.WriteResponse(w, status, utils.ErrorResponse{Status: status, Title: r.Method, Description: "Method Not Allowed"})
		return
	}

	var eventMessage schema.EventMessage
	err := json.NewDecoder(r.Body).Decode(&eventMessage)
	if err != nil {
		status = http.StatusBadRequest
		utils.WriteResponse(w, status, utils.ErrorResponse{Status: status, Title: "Error unmarshalling body", Description: err.Error()})
		return
	}

	notificationMessage := schema.NotificationMessage{
		ID:          rand.Int(),
		Description: eventMessage.Description,
		Timestamp:   time.Now(),
	}

	notificationMessageJson, err := utils.Marshal(notificationMessage)
	if err != nil {
		status = http.StatusInternalServerError
		h.logger.Error("Marshal encode json error: %v", err)
		utils.WriteResponse(w, status, utils.ErrorResponse{Status: status, Title: "Marshal Error", Description: "Marshal Error"})
		return
	}

	err = h.redisClient.WriteEntry(h.ctx, streamEntryKey, notificationMessageJson)
	if err != nil {
		status = http.StatusInternalServerError
		h.logger.Error("Redis write entry error: %v", err)
		utils.WriteResponse(w, status, utils.ErrorResponse{Status: status, Title: "Redis Client error", Description: "Redis Client error"})
		return
	}

	utils.WriteStatusCode(w, status)
}
