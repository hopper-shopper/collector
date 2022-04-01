package contracts

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/patrickmn/go-cache"
	"github.com/steschwa/hopper-analytics-collector/constants"
)

type (
	OnChainClient struct {
		Connection *ethclient.Client
		Cache      *cache.Cache
	}

	ZoneContract interface {
		TotalBaseShare(opts *bind.CallOpts) (*big.Int, error)
		TotalVeShare(opts *bind.CallOpts) (*big.Int, error)
	}
)

func NewOnChainClient() (*OnChainClient, error) {
	client, err := ethclient.Dial(constants.AVAX_RPC)
	if err != nil {
		return &OnChainClient{}, err
	}

	return &OnChainClient{
		Connection: client,
		Cache:      cache.New(time.Minute*1, time.Second*30),
	}, nil
}

// ----------------------------------------
// Adventure callers getter
// ----------------------------------------

func (client *OnChainClient) getPondCaller() (*AdventurePondCaller, error) {
	return NewAdventurePondCaller(common.HexToAddress(constants.ADVENTURE_POND_CONTRACT), client.Connection)
}
func (client *OnChainClient) getStreamCaller() (*AdventureStreamCaller, error) {
	return NewAdventureStreamCaller(common.HexToAddress(constants.ADVENTURE_STREAM_CONTRACT), client.Connection)
}
func (client *OnChainClient) getSwampCaller() (*AdventureSwampCaller, error) {
	return NewAdventureSwampCaller(common.HexToAddress(constants.ADVENTURE_SWAMP_CONTRACT), client.Connection)
}
func (client *OnChainClient) getRiverCaller() (*AdventureRiverCaller, error) {
	return NewAdventureRiverCaller(common.HexToAddress(constants.ADVENTURE_RIVER_CONTRACT), client.Connection)
}
func (client *OnChainClient) getForestCaller() (*AdventureForestCaller, error) {
	return NewAdventureForestCaller(common.HexToAddress(constants.ADVENTURE_FOREST_CONTRACT), client.Connection)
}
func (client *OnChainClient) getGreatLakeCaller() (*AdventureGreatLakeCaller, error) {
	return NewAdventureGreatLakeCaller(common.HexToAddress(constants.ADVENTURE_GREAT_LAKE_CONTRACT), client.Connection)
}
func (client *OnChainClient) getBallotCaller() (*BallotCaller, error) {
	return NewBallotCaller(common.HexToAddress(constants.BALLOT_CONTRACT), client.Connection)
}

func (client *OnChainClient) getAdventureCaller(adventure constants.Adventure) (ZoneContract, error) {
	switch adventure {
	case constants.AdventurePond:
		return client.getPondCaller()
	case constants.AdventureStream:
		return client.getStreamCaller()
	case constants.AdventureSwamp:
		return client.getSwampCaller()
	case constants.AdventureRiver:
		return client.getRiverCaller()
	case constants.AdventureForest:
		return client.getForestCaller()
	case constants.AdventureGreatLake:
		return client.getGreatLakeCaller()
	}

	// should never happen
	return nil, errors.New("unknown adventure")
}

func getContractByAdventure(adventure constants.Adventure) string {
	switch adventure {
	case constants.AdventurePond:
		return constants.ADVENTURE_POND_CONTRACT
	case constants.AdventureStream:
		return constants.ADVENTURE_STREAM_CONTRACT
	case constants.AdventureSwamp:
		return constants.ADVENTURE_SWAMP_CONTRACT
	case constants.AdventureRiver:
		return constants.ADVENTURE_RIVER_CONTRACT
	case constants.AdventureForest:
		return constants.ADVENTURE_FOREST_CONTRACT
	case constants.AdventureGreatLake:
		return constants.ADVENTURE_GREAT_LAKE_CONTRACT
	}

	return ""
}

// ----------------------------------------
// Contract read wrappers
// ----------------------------------------

func (client *OnChainClient) GetTotalBaseShares(adventure constants.Adventure) (*big.Int, error) {
	cacheKey := fmt.Sprintf("%s.total-base-shares", adventure)

	if total, found := client.Cache.Get(cacheKey); found {
		return total.(*big.Int), nil
	}

	caller, err := client.getAdventureCaller(adventure)
	if err != nil {
		return big.NewInt(0), err
	}

	totalBaseShares, err := caller.TotalBaseShare(nil)
	if err != nil {
		return big.NewInt(0), err
	}

	client.Cache.Set(cacheKey, totalBaseShares, cache.DefaultExpiration)
	return totalBaseShares, nil
}

func (client *OnChainClient) GetVotesByAdventure(adventure constants.Adventure) (*big.Int, error) {
	cacheKey := fmt.Sprintf("%s.total-votes", adventure)

	if total, found := client.Cache.Get(cacheKey); found {
		return total.(*big.Int), nil
	}

	caller, err := client.getBallotCaller()
	if err != nil {
		return big.NewInt(0), err
	}

	adventureContract := getContractByAdventure(adventure)

	votes, err := caller.ZonesVotes(nil, common.HexToAddress(adventureContract))
	if err != nil {
		return big.NewInt(0), err
	}

	client.Cache.Set(cacheKey, votes, cache.DefaultExpiration)
	return votes, nil
}
