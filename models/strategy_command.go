package models

const (
	ST_ALIGN_PICCOLO   = 0x03
	ST_ALIGN_GRANDE    = 0x01
	ST_ENABLE_STARTER  = 0x04
	ST_DISABLE_STARTER = 0x05
)

//Position rappresent a 2D point with angle
type StrategyCommand struct {
	CMD          uint8
	FLAGS        uint8
	ELAPSED_TIME int16
}
