package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	if err := db.Create(entity).Error; err != nil {
		r.Log.WithError(err).Error("Failed to create entity")
		return err
	}
	return nil
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	if err := db.Save(entity).Error; err != nil {
		r.Log.WithError(err).Error("Failed to update entity")
		return err
	}
	return nil
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	if err := db.Delete(entity).Error; err != nil {
		r.Log.WithError(err).Error("Failed to delete entity")
		return err
	}
	return nil
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error
	if err != nil {
		r.Log.WithError(err).WithField("id", id).Error("Failed to count entity by ID")
	}
	return total, err
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) error {
	err := db.Where("id = ?", id).Take(entity).Error
	if err != nil {
		r.Log.WithError(err).WithField("id", id).Error("Failed to find entity by ID")
	}
	return err
}
