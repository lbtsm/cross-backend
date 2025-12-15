package models

import "go-admin/internal/constants"

type CrossData struct {
	TxHash           string `json:"tx_hash"`
	Topic            string `json:"topic"`
	Height           int64  `json:"height"`
	OrderId          string `json:"order_id"`
	LogIndex         uint   `json:"log_index"`
	Chain            string `json:"chain"`
	ChainAndGasLimit string `json:"chain_and_gas_limit"`
	Timestamp        int64  `json:"timestamp"`
}

type CrossSet struct {
	Src    *CrossData              `json:"src"`
	Relay  *CrossData              `json:"relay"`
	Dest   *CrossData              `json:"dest"`
	MapDst *CrossData              `json:"map_dest"`
	Now    int64                   `json:"now"`
	Status constants.StatusOfCross `json:"status"`
}

type CrossMapping struct {
	Key      string    `json:"key"`
	CrossSet *CrossSet `json:"cross_set"`
}

type CrossListResponse struct {
	Data []*CrossMapping `json:"data"`
}
