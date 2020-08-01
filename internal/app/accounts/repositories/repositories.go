package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(New)

//var MockProviderSet = wire.NewSet(wire.InterfaceValue(new(Repository),new(MockUsersRepository)))
