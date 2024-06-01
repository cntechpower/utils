package async_task

import (
	"context"

	"github.com/hibiken/asynq"

	"github.com/cntechpower/utils/log"
)

// mux maps a type to a handler
var mux = asynq.NewServeMux()

var srv *asynq.Server

var cli *asynq.Client

var logger *log.Header

func RegisterTaskHandlerFunc(taskType string, handler func(context.Context, *asynq.Task) error) {
	mux.HandleFunc(taskType, handler)
}

func RegisterTaskHandler(taskType string, handler asynq.Handler) {
	mux.Handle(taskType, handler)
}

func AddTask(taskType string, payload []byte, opts ...asynq.Option) error {
	task := asynq.NewTask(taskType, payload, opts...)
	taskInfo, err := cli.Enqueue(task, opts...)
	if err != nil {
		logger.Errorf("add task %v  error %v", task, err)
		return err
	}
	logger.Infof("add task %v success, taskInfo %v", task, taskInfo)
	return nil
}

func Init(redisUri string) error {
	srv = asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisUri},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)
	cli = asynq.NewClient(asynq.RedisClientOpt{Addr: redisUri})

	logger = log.NewHeader("async_task")
	return srv.Start(mux)
}

func Stop() {
	if srv != nil {
		srv.Stop()
	}
	if cli != nil {
		err := cli.Close()
		if err != nil {
			logger.Errorf("close cli error", err)
		}
	}
}
