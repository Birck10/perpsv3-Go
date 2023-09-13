package services

import (
	"context"
	"github.com/gateway-fm/perpsv3-Go/errors"
	"github.com/gateway-fm/perpsv3-Go/pkg/logger"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/gateway-fm/perpsv3-Go/contracts/coreGoerli"
	"github.com/gateway-fm/perpsv3-Go/contracts/perpsMarketGoerli"
	"github.com/gateway-fm/perpsv3-Go/contracts/spotMarketGoerli"
	"github.com/gateway-fm/perpsv3-Go/models"
)

// IService is a service layer interface
type IService interface {
	// RetrieveTrades is used to get logs from the "OrderSettled" event preps market contract within given block range
	RetrieveTrades(fromBlock uint64, toBLock *uint64) ([]*models.Trade, error)

	// RetrieveTradesLimit is used to get all trades and their additional data from the contract with given block search
	// limit. For most public RPC providers the value for limit is 20 000 blocks
	RetrieveTradesLimit(limit uint64) ([]*models.Trade, error)

	// RetrieveOrders is used to get logs from the "OrderCommitted" event preps market contract within given block range
	RetrieveOrders(fromBlock uint64, toBLock *uint64) ([]*models.Order, error)

	// RetrieveOrdersLimit is used to get all orders and their additional data from the contract with given block search
	// limit. For most public RPC providers the value for limit is 20 000 blocks
	RetrieveOrdersLimit(limit uint64) ([]*models.Order, error)

	// RetrieveMarketUpdates is used to get logs from the "MarketUpdated" event preps market contract within given block
	// range
	RetrieveMarketUpdates(fromBlock uint64, toBLock *uint64) ([]*models.MarketUpdate, error)

	// RetrieveMarketUpdatesLimit is used to get all market updates and their additional data from the contract with given block search
	// limit. For most public RPC providers the value for limit is 20 000 blocks
	RetrieveMarketUpdatesLimit(limit uint64) ([]*models.MarketUpdate, error)

	// RetrieveLiquidations is used to get logs from the "PositionLiquidated" event preps market contract within given block
	// range
	RetrieveLiquidations(fromBlock uint64, toBLock *uint64) ([]*models.Liquidation, error)

	// RetrieveLiquidationsLimit is used to get all liquidations and their additional data from the contract with given block search
	// limit. For most public RPC providers the value for limit is 20 000 blocks
	RetrieveLiquidationsLimit(limit uint64) ([]*models.Liquidation, error)

	// GetPosition is used to get "Position" data struct from the latest block from the perps market with given data
	GetPosition(accountID *big.Int, marketID *big.Int) (*models.Position, error)

	// GetMarketMetadata is used to get market metadata by given market ID
	GetMarketMetadata(marketID *big.Int) (*models.MarketMetadata, error)

	// FormatAccount is used to get account, and it's additional data from the contract by given account id
	FormatAccount(id *big.Int) (*models.Account, error)

	// FormatAccounts is used to get all accounts and their additional data from the contract
	FormatAccounts() ([]*models.Account, error)

	// FormatAccountsLimit is used to get all accounts and their additional data from the contract with given block search
	// limit. For most public RPC providers the value for limit is 20 000 blocks
	FormatAccountsLimit(limit uint64) ([]*models.Account, error)
}

// Service is an implementation of IService interface
type Service struct {
	rpcClient             *ethclient.Client
	core                  *coreGoerli.CoreGoerli
	coreFirstBlock        uint64
	spotMarket            *spotMarketGoerli.SpotMarketGoerli
	spotMarketFirstBlock  uint64
	perpsMarket           *perpsMarketGoerli.PerpsMarketGoerli
	perpsMarketFirstBlock uint64
}

// NewService is used to get instance of Service
func NewService(
	rpc *ethclient.Client,
	core *coreGoerli.CoreGoerli,
	coreFirstBlock uint64,
	spotMarket *spotMarketGoerli.SpotMarketGoerli,
	spotMarketFirstBlock uint64,
	perpsMarket *perpsMarketGoerli.PerpsMarketGoerli,
	perpsMarketFirstBlock uint64,
) IService {
	return &Service{
		rpcClient:             rpc,
		core:                  core,
		coreFirstBlock:        coreFirstBlock,
		spotMarket:            spotMarket,
		spotMarketFirstBlock:  spotMarketFirstBlock,
		perpsMarket:           perpsMarket,
		perpsMarketFirstBlock: perpsMarketFirstBlock,
	}
}

// getIterationsForLimitQuery is used to get iterations of querying data from the contract with given rpc limit for blocks
// and latest block number. Limit by default (if given limit is 0) is set to 20 000 blocks
func (s *Service) getIterationsForLimitQuery(limit uint64) (iterations uint64, lastBlock uint64, err error) {
	lastBlock, err = s.rpcClient.BlockNumber(context.Background())
	if err != nil {
		logger.Log().WithField("layer", "Service-getIterationsForLimitQuery").Errorf("get latest block rpc error: %v", err.Error())
		return 0, 0, errors.GetRPCProviderErr(err, "BlockNumber")
	}

	if limit == 0 {
		limit = 20000
	}

	iterations = (lastBlock-s.perpsMarketFirstBlock)/limit + 1

	return iterations, lastBlock, nil
}

// getFilterOptsPerpsMarket is used to get options for event filtering on perps market contract
func (s *Service) getFilterOptsPerpsMarket(fromBlock uint64, toBLock *uint64) *bind.FilterOpts {
	if fromBlock == 0 {
		fromBlock = s.perpsMarketFirstBlock
	}

	return &bind.FilterOpts{
		Start:   fromBlock,
		End:     toBLock,
		Context: context.Background(),
	}
}
