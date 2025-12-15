package models

import (
	"time"
)

type CrossInfo struct {
	Id           int       `json:"id" gorm:"primaryKey;autoIncrement"`                //
	Project      string    `json:"project" gorm:"size:64"`                            //
	OrderId      string    `json:"orderId" gorm:"size:128" uniqueIndex:unique_index"` //
	SrcChain     string    `json:"srcChain" gorm:"size:64"`
	SrcTxHash    string    `json:"srcTxHash" gorm:"size:128"` //
	SrcInfo      string    `json:"srcInfo" gorm:"type:text"`
	RelayTxHash  string    `json:"relayTxHash" gorm:"size:128"`
	RelayInfo    string    `json:"relayInfo" gorm:"type:text"`
	DstChain     string    `json:"dstChain" gorm:"size:64"`
	DstTxHash    string    `json:"dstTxHash" gorm:"size:128"` //
	DstInfo      string    `json:"dstInfo" gorm:"type:text"`
	MapDstTxHash string    `json:"mapDstTxHash" gorm:"size:128"`
	MapDstInfo   string    `json:"mapDstInfo" gorm:"type:text"`
	Status       int64     `json:"status" gorm:"size:20"` //
	CostTime     int64     `json:"costTime" gorm:"size:20"`
	CreatedAt    time.Time `json:"createdAt"` //
	UpdatedAt    time.Time `json:"updatedAt"` //
	CreateBy     string    `json:"createBy"`  //
	UpdateBy     string    `json:"updateBy"`  //
}

func (CrossInfo) TableName() string {
	return "cross_info"
}
