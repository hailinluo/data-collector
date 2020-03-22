package task

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

type Task interface {
	Spec() string
	Run()
}

type TaskHub struct {
	*cron.Cron
}

func InitTaskHub() *TaskHub {
	// Second | Minute | Hour | Dom | Month | Dow | Descriptor,
	controller := &TaskHub{Cron: cron.New(cron.WithSeconds())}
	controller.Start()
	return controller
}

func (c *TaskHub) AddTask(t Task) (cron.EntryID, error) {
	if t == nil {
		return 0, errors.New("with empty task.")
	}

	id, err := c.AddFunc(t.Spec(), func() {
		t.Run()
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *TaskHub) Close() error {
	ctx := c.Stop()
	select {
	case <-ctx.Done():
	}
	return nil
}
