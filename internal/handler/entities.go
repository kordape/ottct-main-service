package handler

import (
	"fmt"

	"github.com/kordape/ottct-main-service/pkg/logger"
)

type Entity struct {
	Id          string
	TwitterId   string
	DisplayName string
}

type EntityManager struct {
	storage EntityStorage
	log     logger.Interface
}

type EntityStorage interface {
	GetEntity(id string) (*Entity, error)
	GetEntities() ([]Entity, error)
}

func NewEntityManager(entityStorage EntityStorage, log logger.Interface) *EntityManager {
	return &EntityManager{
		storage: entityStorage,
		log:     log,
	}
}

func (m EntityManager) GetEntity(id string) (entity *Entity, err error) {
	entity, err = m.storage.GetEntity(id)
	if err != nil {
		m.log.Error(fmt.Errorf("[EntityManager] Failed to get entity by id: %w", err))
		return nil, fmt.Errorf("[EntityManager] storage error: %w", err)
	}

	return
}

func (m EntityManager) GetEntities() (entities []Entity, err error) {
	entities, err = m.storage.GetEntities()
	if err != nil {
		m.log.Error(fmt.Errorf("[EntityManager] Failed to get entities: %w", err))
		return nil, fmt.Errorf("[EntityManager] storage error: %w", err)
	}

	return
}