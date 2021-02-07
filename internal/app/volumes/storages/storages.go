package storages

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewFileManager)
