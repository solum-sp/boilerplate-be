package repositories

import (
	"context"
	"fmt"
	"time"

	"proposal-template/models"

	"gorm.io/gorm"
)




type GenericDAO[T any] struct {
	db *gorm.DB
	tableName string
}

func NewGenericDAO[T any](db *gorm.DB, tableName string) *GenericDAO[T] {
	return &GenericDAO[T]{
		db: db,
		tableName: tableName,
	}
}

func (dao *GenericDAO[T]) GetByColumn(ctx context.Context, column string, value interface{}) (*T, error) {
	var obj T
	fmt.Println("colum:", column)
	err := dao.db.WithContext(ctx).
		Table(dao.tableName).
		Where(fmt.Sprintf("%s = ?", column), value).
		First(&obj).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error retrieving data: %w", err)
	}
	return &obj, nil
}


func (dao *GenericDAO[T]) List(ctx context.Context, paging model.Paging, query *gorm.DB) ([]T, error) {
	paging.Validate()

	var results []T
	offset := (paging.Page - 1) * paging.Limit

	err := query.WithContext(ctx).
		Table(dao.tableName).
		Limit(paging.Limit).
		Offset(offset).
		Order("id DESC").
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %w", err)
	}
	return results, nil
}

func (dao *GenericDAO[T]) Create(ctx context.Context, model T) (uint, error) {
	now := time.Now().UTC()

	// Set timestamps if the struct supports it
	if v, ok := any(model).(interface{ SetCreatedAt(time.Time) }); ok {
		v.SetCreatedAt(now)
	}
	if v, ok := any(model).(interface{ SetUpdatedAt(time.Time) }); ok {
		v.SetUpdatedAt(now)
	}

	// Insert only non-zero fields (ignore empty fields)
	err := dao.db.WithContext(ctx).
		Table(dao.tableName).
		Omit("ID"). // Exclude ID if it's auto-generated
		Create(&model).Error

	if err != nil {
		return 0, fmt.Errorf("error inserting data: %w", err)
	}

	// Extract ID (assuming `ID` is the primary key)
	var idField uint
	if v, ok := any(model).(interface{ GetID() uint }); ok {
		idField = v.GetID()
	}

	return idField, nil
}


