package blockchain

import (
	"context"
)

type IBlockchainRepository interface {
	GetLastBlock(ctx context.Context) uint64
	SaveBlocks(ctx context.Context, blocks []*Block) error
	SaveLastBlockNumber(ctx context.Context, blockNumber uint64) error
	SaveTransactions(ctx context.Context, txs []*Transaction) error
	SaveEvents(ctx context.Context, events []*Event) error
}