package task_api

import (
	"context"
	"fmt"
	"time"

	"github.com/ozoncp/ocp-facade-api/internal/models"

	taskApi "github.com/ozoncp/ocp-task-api/pkg/ocp-task-api"

	"google.golang.org/grpc"
)

// Api интерфейс взаимодейстия с ocp-task-api
type Api interface {
	DescribeTask(
		ctx context.Context,
		taskId uint64,
	) (*models.Task, error)
}

// Settings настройки для Api
type Settings struct {
	Timeout time.Duration
}

type api struct {
	settings Settings
	client   taskApi.OcpTaskApiClient
}

// NewApi возвращает имплементацию с поддержкой таймаута
func NewApi(
	settings Settings,
	conn *grpc.ClientConn,
) Api {
	client := taskApi.NewOcpTaskApiClient(conn)
	return &api{
		settings: settings,
		client:   client,
	}
}

func (a *api) DescribeTask(
	ctx context.Context,
	taskId uint64,
) (*models.Task, error) {

	ctx, cancel := context.WithTimeout(ctx, a.settings.Timeout)
	defer cancel()

	req := &taskApi.DescribeTaskV1Request{
		TaskId: taskId,
	}
	resp, err := a.client.DescribeTaskV1(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("describe task with id <%d>: %w", taskId, err)
	}

	task := &models.Task{
		Id:          resp.Task.Id,
		Description: resp.Task.Description,
	}

	return task, nil
}
