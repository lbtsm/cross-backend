package dto

import "go-admin/common/dto"

type CrossListReq struct {
	dto.Pagination `search:"-"`
	SrcChain       string `form:"srcChain"  search:"type:exact;column:src_chain;table:cross_info" comment:"来源链"`
	DstChain       string `form:"dstChain"  search:"type:exact;column:dst_chain;table:cross_info" comment:"目标链"`
	SrcTxHash      string `form:"srcTxHash"  search:"type:exact;column:src_tx_hash;table:cross_info" comment:"交易哈希，可以是源链、中继链、目标链的交易哈希"`
	OrderId        string `form:"orderId"  search:"type:exact;column:order_id;table:cross_info" comment:"跨链订单ID"`
}

func (m *CrossListReq) GetNeedSearch() interface{} {
	return *m
}

type SingleCrossReq struct {
	Id int `json:"id"`
}

func (s *SingleCrossReq) GetId() interface{} {
	return s.Id
}
