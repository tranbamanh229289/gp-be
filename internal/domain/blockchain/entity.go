package blockchain

import (
	"math/big"
	"time"
)


type Block struct {
	Number uint64					`json:"number"`
	Hash string 					`json:"hash"`
	Timestamp uint64      `json:"timestamp"`
	TxCount int				`json:"tx_count"`
	GasUsed uint64				`json:"gas_used"`
	CrawledAt time.Time		`json:"crawled_at"`
}

type Transaction struct {
	Hash string						`json:"hash"`
	BlockNumber uint64		`json:"block_number"`
	From string						`json:"from"`
	To string							`json:"to"`
	Value *big.Int				`json:"value"`
	Timestamp uint64			`json:"timestamp"`
	GasPrice string				`json:"gas_price"`
	GasUsed uint64				`json:"gas_used"`
	Status uint64					`json:"status"`
	CrawledAt time.Time		`json:"crawled_at"`
}

type Event struct {
	ContractAddress string	`json:"contract_address"`
	EventName string				`json:"event_name"`
	Topics []string					`json:"topics"`
	Data string							`json:"data"`
	BlockNumber uint64			`json:"block_number"`
	TxHash string						`json:"tx_hash"`
	CrawledAt time.Time			`json:"crawled_at"`
}

type TokenTransfer struct {
	TokenAddress string			`json:"token_address"`
	From string							`json:"from"`
	To string								`json:"to"`
}

type Balance struct {
	Address string					`json:"address"`
	Balance string					`json:"balance"`
	Block uint64						`json:"block"`
}