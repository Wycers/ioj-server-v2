package repositories

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewJudgementRepository)

//var MockProviderSet = wire.NewSet(wire.InterfaceValue(new(Repository),new(MockJudgementsRepository)))
