package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewFileManager)

//var MockProviderSet = wire.NewSet(wire.InterfaceValue(new(Repository),new(MockFilesRepository)))
