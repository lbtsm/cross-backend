package constants

type StatusOfCross int64

const (
	StatusOfInit StatusOfCross = iota
	StatusOfPending
	StatusOfSend
	StatusOfCompleted
	StatusOfFailed
)

var StatusOfCrossMap = map[StatusOfCross]string{
	StatusOfInit:      "初始化",
	StatusOfPending:   "待处理",
	StatusOfSend:      "已发送",
	StatusOfCompleted: "已完成",
	StatusOfFailed:    "失败",
}
