package tt

import (
	"fmt"
)

func Start(project, task string, tags []string, timestamp string) (*GRpcTimer, error) {
	t := &GRpcTimer{
		Start:   timestamp,
		Project: project,
		Task:    task,
		Tags:    tags,
	}
	err := t.Validate()
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}
	// err = db.SaveTimer(t)
	// if err != nil {
	// 	return nil, fmt.Errorf("start: %w", err)
	// }
	return t, nil
}
