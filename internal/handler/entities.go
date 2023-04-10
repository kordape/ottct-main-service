package handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Entity struct {
	Id          string
	TwitterId   string
	DisplayName string
}

type EntityManager struct {
	storage EntityStorage
}

type EntityStorage interface {
	GetEntity(id string) (*Entity, error)
	GetEntities() ([]Entity, error)
}

func NewEntityManager(entityStorage EntityStorage) *EntityManager {
	return &EntityManager{
		storage: entityStorage,
	}
}

func (m EntityManager) GetEntity(id string, log *logrus.Entry) (entity *Entity, err error) {
	entity, err = m.storage.GetEntity(id)
	if err != nil {
		log.WithError(err).Error("[EntityManager] Failed to get entity by id")
		return nil, fmt.Errorf("[EntityManager] storage error: %w", err)
	}

	return
}

func (m EntityManager) GetEntities(log *logrus.Entry) (entities []Entity, err error) {
	entities, err = m.storage.GetEntities()
	if err != nil {
		log.WithError(err).Error("[EntityManager] Failed to get entities")
		return nil, fmt.Errorf("[EntityManager] storage error: %w", err)
	}

	return
}
