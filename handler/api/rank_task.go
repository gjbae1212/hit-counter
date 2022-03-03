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
	dailyHitTTL := time.Hour * 24 * 7 // 7 days.

	// If a domain is 'github.com', it is calculating ranks.
	if task.Domain == githubGroup && task.Path != "" {
		// Calculate visiting count based on gitHub.com.
		if _, err := task.Counter.IncreaseRankOfDaily(ctx, githubGroup, task.Path, task.CreatedAt, dailyHitTTL); err != nil {
			return err
		}
		if _, err := task.Counter.IncreaseRankOfTotal(ctx, githubGroup, task.Path); err != nil {
			return err
		}

		// Calculate sum of visiting count for github projects based on github profile.
		seps := strings.Split(task.Path, "/")
		if len(seps) >= 2 && seps[1] != "" {
			if _, err := task.Counter.IncreaseRankOfDaily(ctx, githubProfileSumGroup, seps[1], task.CreatedAt, dailyHitTTL); err != nil {
				return err
			}
			if _, err := task.Counter.IncreaseRankOfTotal(ctx, githubProfileSumGroup, seps[1]); err != nil {
				return err
			}
		}
	}

	// Calculate domain visiting count for daily and total.
	if _, err := task.Counter.IncreaseRankOfDaily(ctx, domainGroup, task.Domain, task.CreatedAt, dailyHitTTL); err != nil {
		return err
	}
	if _, err := task.Counter.IncreaseRankOfTotal(ctx, domainGroup, task.Domain); err != nil {
		return err
	}
	return nil
}
