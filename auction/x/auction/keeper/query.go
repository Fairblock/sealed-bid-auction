package keeper

import (
	"auction/x/auction/types"
)

var _ types.QueryServer = Keeper{}
