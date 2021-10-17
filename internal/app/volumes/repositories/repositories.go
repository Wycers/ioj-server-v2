package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewRepository)

//var MockProviderSet = wire.NewSet(wire.InterfaceValue(new(repository),new(MockFilesRepository)))
