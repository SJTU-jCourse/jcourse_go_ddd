package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"jcourse_go/internal/domain/event"
)

type asynqEventBus struct {
	client *asynq.Client
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewAsynqEventBus(redisAddr string) (event.EventBusPublisher, error) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Queues: map[string]int{
				"default": 10,
			},
		},
	)

	mux := asynq.NewServeMux()

	return &asynqEventBus{
		client: client,
		server: server,
		mux:    mux,
	}, nil
}

func (b *asynqEventBus) Register(eventType event.Type, handler event.Handler) error {
	taskType := fmt.Sprintf("event:%d", eventType)

	b.mux.HandleFunc(taskType, func(ctx context.Context, task *asynq.Task) error {
		var eventPayload map[string]interface{}
		if err := json.Unmarshal(task.Payload(), &eventPayload); err != nil {
			return err
		}

		baseEvent := event.NewBaseEvent(event.Type(eventPayload["type"].(float64)), nil)
		baseEvent.SetPayload(nil)

		payload, err := b.unmarshalPayload(eventType, task.Payload())
		if err != nil {
			return err
		}
		baseEvent.SetPayload(payload)

		return handler.Handle(ctx, baseEvent)
	})

	return nil
}

func (b *asynqEventBus) Dispatch(ctx context.Context, events ...event.Event) error {
	for _, event := range events {
		taskType := fmt.Sprintf("event:%d", event.Type())

		data, err := event.ToJSON()
		if err != nil {
			return err
		}

		task := asynq.NewTask(taskType, data)
		info, err := b.client.Enqueue(task)
		if err != nil {
			return err
		}

		if info == nil {
			return fmt.Errorf("failed to enqueue event: %d", event.Type())
		}
	}
	return nil
}

func (b *asynqEventBus) Publish(ctx context.Context, events ...event.Event) error {
	return b.Dispatch(ctx, events...)
}

func (b *asynqEventBus) Start() error {
	return b.server.Run(b.mux)
}

func (b *asynqEventBus) Shutdown() error {
	b.client.Close()
	b.server.Shutdown()
	return nil
}

func (b *asynqEventBus) unmarshalPayload(eventType event.Type, data []byte) (event.Payload, error) {
	switch eventType {
	case event.TypeReviewCreated, event.TypeReviewModified:
		var payload event.ReviewPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		return &payload, nil
	default:
		return nil, fmt.Errorf("unknown event type: %d", eventType)
	}
}
