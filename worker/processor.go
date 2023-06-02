package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/techschool/simplebank/db/sqlc"
)

const (
	QUEUECRITICAL = "critical"
	QUEUEDEFAULT  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessorTaskVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt *asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QUEUECRITICAL: 10,
				QUEUEDEFAULT:  5,
			},
		},
	)
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TASKSENDVERIFYEMAIL, processor.ProcessorTaskVerifyEmail)

	return processor.server.Start(mux)
}
