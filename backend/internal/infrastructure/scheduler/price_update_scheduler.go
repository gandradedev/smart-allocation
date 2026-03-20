package scheduler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	assetusecase "smart-allocation/internal/application/usecase/asset"
	domainrepo "smart-allocation/internal/domain/repository"
)

// PriceUpdateScheduler runs a daily job at midnight to refresh all asset prices.
type PriceUpdateScheduler struct {
	cron         *cron.Cron
	repo         domainrepo.AssetRepository
	priceUpdater assetusecase.UpdateAssetPriceUseCase
}

func NewPriceUpdateScheduler(repo domainrepo.AssetRepository, priceUpdater assetusecase.UpdateAssetPriceUseCase) *PriceUpdateScheduler {
	return &PriceUpdateScheduler{
		cron:         cron.New(),
		repo:         repo,
		priceUpdater: priceUpdater,
	}
}

// Start registers the midnight job and starts the scheduler.
func (s *PriceUpdateScheduler) Start() {
	s.cron.AddFunc("0 0 * * *", s.updateAllPrices)
	s.cron.Start()
	log.Println("Price update scheduler started (runs daily at midnight)")
}

// Stop gracefully shuts down the scheduler.
func (s *PriceUpdateScheduler) Stop() {
	s.cron.Stop()
}

func (s *PriceUpdateScheduler) updateAllPrices() {
	ctx := context.Background()

	assets, err := s.repo.FindAll(ctx)
	if err != nil {
		log.Printf("scheduler: failed to fetch assets: %v", err)
		return
	}

	log.Printf("scheduler: updating prices for %d assets", len(assets))

	for _, a := range assets {
		if err := s.priceUpdater.Execute(ctx, a.Ticker); err != nil {
			log.Printf("scheduler: failed to update price for %s: %v", a.Ticker, err)
		} else {
			log.Printf("scheduler: updated price for %s", a.Ticker)
		}
	}
}
