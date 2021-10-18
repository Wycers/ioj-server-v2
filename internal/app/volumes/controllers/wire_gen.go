// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package controllers

import (
	"github.com/google/wire"
	"github.com/infinity-oj/server-v2/internal/app/volumes/repositories"
	"github.com/infinity-oj/server-v2/internal/app/volumes/services"
	"github.com/infinity-oj/server-v2/internal/app/volumes/storages"
	"github.com/infinity-oj/server-v2/internal/pkg/configs"
	"github.com/infinity-oj/server-v2/internal/pkg/database"
	"github.com/infinity-oj/server-v2/internal/pkg/log"
)

// Injectors from wire.go:

func CreateVolumesController(cf string, sto storages.Storage, rep repositories.Repository) (Controller, error) {
	viper, err := configs.New(cf)
	if err != nil {
		return nil, err
	}
	options, err := log.NewOptions(viper)
	if err != nil {
		return nil, err
	}
	logger, err := log.New(options)
	if err != nil {
		return nil, err
	}
	service := services.NewVolumeService(logger, sto, rep)
	controller := New(logger, service)
	return controller, nil
}

// wire.go:

var testProviderSet = wire.NewSet(log.ProviderSet, configs.ProviderSet, database.ProviderSet, services.ProviderSet, ProviderSet)
