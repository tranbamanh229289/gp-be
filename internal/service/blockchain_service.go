package service

import (
	"be/internal/domain/blockchain"
	"be/internal/infrastructure/blockchain/ether"
	"be/internal/shared/constant"
	"be/pkg/logger"
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockchainService struct {
	logger *logger.ZapLogger
	ethclient *ether.Ether
	blockchainRepo blockchain.IBlockchainRepository
}	

func NewBlockchainService(logger *logger.ZapLogger, ethclient *ether.Ether, blockchainRepo blockchain.IBlockchainRepository) *BlockchainService {
	return &BlockchainService{
		logger: logger,
		ethclient: ethclient,
		blockchainRepo: blockchainRepo,
	}
}

func (s *BlockchainService) CrawlBlocks(ctx context.Context) error { 
	s.logger.Info("Starting block crawler...")

	lastBlockCrawled := s.blockchainRepo.GetLastBlock(ctx)

	lastBlock, err := s.ethclient.GetClient().BlockNumber(ctx)
	if err != nil {
		return &constant.BadRequest
	}

	var blockDatas []*blockchain.Block
	for blockNum := lastBlockCrawled + 1; blockNum < lastBlock; blockNum ++ {
		block, err := s.ethclient.GetClient().BlockByNumber(ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			return &constant.BadRequest
		}

		blockData := blockchain.Block{
			Number: block.NumberU64(),
			Hash: block.Hash().Hex(),
			Timestamp: block.Time(),
			TxCount: len(block.Transactions()),
			GasUsed: block.GasUsed(),
			CrawledAt: time.Now(),
		}
		blockDatas = append(blockDatas, &blockData)
	}
	err = s.blockchainRepo.SaveBlocks(ctx, blockDatas)

	if err != nil {
			return &constant.InternalServer
	}

	err = s.blockchainRepo.SaveLastBlockNumber(ctx, lastBlock)

	if err != nil {
			return &constant.InternalServer
	}
	
	return nil
}

func (s *BlockchainService) CrawlTransactions(ctx context.Context, address string) error {
	s.logger.Info("Starting transaction crawler...")

	lastBlockCrawled := s.blockchainRepo.GetLastBlock(ctx)

	lastBlock, err := s.ethclient.GetClient().BlockNumber(ctx)
	if err != nil {
		return &constant.BadRequest
	}
	
	chainID, err := s.ethclient.GetClient().NetworkID(ctx)
	if err != nil {
		return &constant.BadRequest
	}
	var txDatas []*blockchain.Transaction
	for blockNum := lastBlockCrawled + 1; blockNum < lastBlock; blockNum ++ {
		block, err := s.ethclient.GetClient().BlockByNumber(ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			return &constant.BadRequest
		}
		
		for _, tx := range block.Transactions() {
			signer := types.NewEIP155Signer(chainID)
			from, err := types.Sender(signer, tx)
			if err != nil {
				return &constant.BadRequest
			}
			to := tx.To()
			if to == nil {
				return &constant.BadRequest
			}
			addr := common.HexToAddress(address)
			if from == addr || *to == addr {
				receipt, err := s.ethclient.GetClient().TransactionReceipt(ctx, tx.Hash())
				if err != nil {
					return &constant.BadRequest
				}

				txData := blockchain.Transaction{
					Hash: tx.Hash().Hex(),
					From: from.Hex(),
					To: to.Hex(),
					Value: tx.Value(),
					BlockNumber: block.NumberU64(),
					Timestamp: block.Time(),
					GasPrice: tx.GasPrice().String(),
					GasUsed: receipt.GasUsed,
					Status: receipt.Status,
					CrawledAt: time.Now(),
				}
				txDatas = append(txDatas, &txData)
			}
		}
	}

	err = s.blockchainRepo.SaveTransactions(ctx, txDatas)

	if err != nil {
			return &constant.InternalServer
	}

	err = s.blockchainRepo.SaveLastBlockNumber(ctx, lastBlock)

	if err != nil {
			return &constant.InternalServer
	}

	return nil
}

func (s *BlockchainService) CrawlEvents(ctx context.Context, contractAddress string) error {
	s.logger.Info("Starting event crawler...")

	lastBlockCrawled := s.blockchainRepo.GetLastBlock(ctx)

	lastBlock, err := s.ethclient.GetClient().BlockNumber(ctx)
	if err != nil {
		return &constant.BadRequest
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(lastBlockCrawled +1)),
		ToBlock: big.NewInt(int64(lastBlock)),
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
	}

	logs, err := s.ethclient.GetClient().FilterLogs(ctx, query)

	if err != nil {
		return &constant.BadRequest
	}
	var events [] *blockchain.Event;
	for _, log :=range logs {
		topics := make([]string, len(log.Topics))
		for i, topic := range log.Topics {
			topics[i] = topic.Hex()
		}

		event := blockchain.Event{
			ContractAddress: log.Address.Hex(),
			EventName: "",
			Topics: topics,
			Data: common.Bytes2Hex(log.Data),
			BlockNumber: log.BlockNumber,
			TxHash: log.TxHash.Hex(),
			CrawledAt: time.Now(),
		}
		events = append(events, &event)
	}

		err = s.blockchainRepo.SaveEvents(ctx, events)

	if err != nil {
			return &constant.InternalServer
	}

	err = s.blockchainRepo.SaveLastBlockNumber(ctx, lastBlock)

	if err != nil {
			return &constant.InternalServer
	}

	return nil
}