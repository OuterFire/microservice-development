package api

import (
	"context"
	"errors"
	"net/http"
	"rest-server/logger"
	"rest-server/utils"
	"sync"
	"time"

	"rest-server/redis"
)

type ApiService struct {
	logger      *logger.Log
	ctx         context.Context
	cancel      context.CancelFunc
	redisClient *redis.RedisClient
	muxHandler  *utils.MuxHandler
	server      *http.Server
	cfg         *config
}

func NewApiService() *ApiService {
	log := logger.NewLogger("ApiService")
	ctx, cancel := context.WithCancel(context.Background())

	var cfg config
	newConfig(&cfg)

	redisClient := redis.NewRedisClient()
	muxHandler := utils.NewMuxHandler()
	setHandler(muxHandler, redisClient, ctx, log)

	return &ApiService{
		logger:      log,
		ctx:         ctx,
		cancel:      cancel,
		redisClient: redisClient,
		muxHandler:  muxHandler,
		cfg:         &cfg,
	}
}

func setHandler(mux *utils.MuxHandler, redisClient *redis.RedisClient, ctx context.Context, logger *logger.Log) {
	h := &handler{
		logger:      logger,
		ctx:         ctx,
		redisClient: redisClient,
	}
	mux.AddHandlers("/", h.defaultEndpoint)
	mux.AddHandlers("/event", h.eventEndpoint)

	for endpoint, handlerFunc := range mux.Handlers {
		mux.Mux.HandleFunc(endpoint, handlerFunc)
	}
}

func (s *ApiService) Start() {
	s.logger.Debug("Starting ApiService")

	var wg sync.WaitGroup
	wg.Add(2)

	go s.redisClientLiveness(&wg)

	go func() {
		err := s.startServer(&wg)
		if err != nil {
			s.logger.Error("Error starting server: %v", err)
		}
	}()
	wg.Wait()

	if s.ctx.Err() != nil {
		s.logger.Warn("ApiService is shutting down: %v", s.ctx.Err().Error())
		return
	}
}

func (s *ApiService) Stop() error {
	s.logger.Warn("Stopping ApiService")

	s.cancel()

	return nil
}

func (s *ApiService) startServer(wg *sync.WaitGroup) error {
	s.logger.Debug("Starting server connection on port 8080")
	defer wg.Done()
	server := utils.NewHttpServer(s.muxHandler.Mux, s.cfg.port, s.cfg.writeTimeout)
	s.server = server
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *ApiService) redisClientLiveness(wg *sync.WaitGroup) {
	s.logger.Debug("Starting redis connection")
	defer wg.Done()

	for {
		err := s.redisConnect()
		if errors.Is(err, context.Canceled) {
			s.logger.Warn("Exist connection check: %v", s.ctx.Err().Error())
			return
		}

		if err == nil {
			s.logger.Debug("Redis client connected and stream created")
			err = s.streamConnectionCheck()
			if errors.Is(err, context.Canceled) {
				s.logger.Warn("Exist connection check: %v", s.ctx.Err().Error())
				return
			}
		}

		s.logger.Error("Error with Redis Client: %v", err)

		select {
		case <-time.After(2 * time.Second):
			s.logger.Debug("Restarting redis connection")
			continue
		case <-s.ctx.Done():
			s.logger.Warn("Exist redis connection: %v", s.ctx.Err().Error())
			return
		}
	}
}

func (s *ApiService) redisConnect() error {
	s.redisClient.Connect()

	err := s.redisClient.WriteEntry(s.ctx, streamCreateEntryKey, "stream created")
	if err != nil {
		return err
	}
	return nil
}

func (s *ApiService) streamConnectionCheck() error {
	s.logger.Debug("Starting connection check")

	ticker := time.NewTicker(pingCheckInterval)
	var err error
	for err = s.redisClient.Ping(s.ctx); err == nil; err = s.redisClient.Ping(s.ctx) {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case <-ticker.C:
			ticker.Reset(pingCheckInterval)
		}
	}

	return err
}
