package api_handler

import (
	"context"
	"strings"
	"time"

	"github.com/gjbae1212/hit-counter/counter"
)

const (
	domainGroup           = "domain"
	githubGroup           = "github.com"
	githubProfileSumGroup = "github.com-profile-sum"
)

type RankTask struct {
	Counter   counter.Counter
	Domain    string
	Path      string
	CreatedAt time.Time
}

// Process is a specific method implemented Task interface in async task.
func (task *RankTask) Process(ctx context.Context) error {
	// If a domain is 'github.com', it is calculating ranks.
	if task.Domain == githubGroup && task.Path != "" {
		// Calculate visiting count based on github.com.
		if _, err := task.Counter.IncreaseRankOfDaily(githubGroup, task.Path, task.CreatedAt); err != nil {
			return err
		}
		if _, err := task.Counter.IncreaseRankOfTotal(githubGroup, task.Path); err != nil {
			return err
		}

		// Calculate sum of visiting count for github projects based on github profile.
		seps := strings.Split(task.Path, "/")
		if len(seps) >= 2 && seps[1] != "" {
			if _, err := task.Counter.IncreaseRankOfDaily(githubProfileSumGroup, seps[1], task.CreatedAt); err != nil {
				return err
			}
			if _, err := task.Counter.IncreaseRankOfTotal(githubProfileSumGroup, seps[1]); err != nil {
				return err
			}
		}
	}

	// Calculate visiting count for daily and total.
	if _, err := task.Counter.IncreaseRankOfDaily(domainGroup, task.Domain, task.CreatedAt); err != nil {
		return err
	}
	if _, err := task.Counter.IncreaseRankOfTotal(domainGroup, task.Domain); err != nil {
		return err
	}

	return nil
}
