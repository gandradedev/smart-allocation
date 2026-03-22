package repository

import (
	"context"
	"database/sql"
	"strings"

	"smart-allocation/internal/domain/entity"
	domainerrors "smart-allocation/internal/domain/errors"
)

type AssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(ctx context.Context, a *entity.Asset) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO assets (ticker, asset_type, quantity, price, ceiling_price, target_percent, icon, currency)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, a.Ticker, a.AssetType, a.Quantity, a.Price, a.CeilingPrice, a.TargetPercent, a.Icon, a.Currency)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return domainerrors.NewAlreadyExistsError("Asset already registered in portfolio")
		}
		return err
	}
	return nil
}

func (r *AssetRepository) FindAll(ctx context.Context) ([]*entity.Asset, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ticker, asset_type, quantity, price, ceiling_price, target_percent, icon, currency
		FROM assets
		ORDER BY asset_type, ticker
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*entity.Asset
	for rows.Next() {
		a := &entity.Asset{}
		if err := rows.Scan(&a.Ticker, &a.AssetType, &a.Quantity, &a.Price, &a.CeilingPrice, &a.TargetPercent, &a.Icon, &a.Currency); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

func (r *AssetRepository) FindByTicker(ctx context.Context, ticker string) (*entity.Asset, error) {
	a := &entity.Asset{}
	err := r.db.QueryRowContext(ctx, `
		SELECT ticker, asset_type, quantity, price, ceiling_price, target_percent, icon, currency
		FROM assets
		WHERE ticker = ?
	`, ticker).Scan(&a.Ticker, &a.AssetType, &a.Quantity, &a.Price, &a.CeilingPrice, &a.TargetPercent, &a.Icon, &a.Currency)

	if err == sql.ErrNoRows {
		return nil, domainerrors.NewNotFoundError("Asset not found")
	}
	return a, err
}

func (r *AssetRepository) TotalValue(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(quantity * price), 0) FROM assets`).Scan(&total)
	return total, err
}

func (r *AssetRepository) Update(ctx context.Context, ticker string, a *entity.Asset) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE assets
		SET asset_type = ?, quantity = ?, ceiling_price = ?, target_percent = ?,
		    updated_at = CURRENT_TIMESTAMP
		WHERE ticker = ?
	`, a.AssetType, a.Quantity, a.CeilingPrice, a.TargetPercent, ticker)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domainerrors.NewNotFoundError("Asset not found")
	}
	return nil
}

func (r *AssetRepository) UpdatePrice(ctx context.Context, ticker string, price float64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE assets
		SET price = ?, updated_at = CURRENT_TIMESTAMP
		WHERE ticker = ?
	`, price, ticker)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domainerrors.NewNotFoundError("Asset not found")
	}
	return nil
}

func (r *AssetRepository) UpdateMetadata(ctx context.Context, ticker string, price float64, icon, currency string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE assets
		SET price = ?, icon = ?, currency = ?, updated_at = CURRENT_TIMESTAMP
		WHERE ticker = ?
	`, price, icon, currency, ticker)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domainerrors.NewNotFoundError("Asset not found")
	}
	return nil
}

func (r *AssetRepository) Delete(ctx context.Context, ticker string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM assets WHERE ticker = ?`, ticker)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domainerrors.NewNotFoundError("Asset not found")
	}
	return nil
}
