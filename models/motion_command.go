package models

const (
	MC_STOP            = 0x82
	MC_BRAKE           = 0x83
	MC_SET_POSITION    = 0x84
	MC_FW_TO_DISTANCE  = 0x85
	MC_ROTATE_RELATIVE = 0x88
	MC_SET_SPEED       = 0x8C
)

//Position rappresent a 2D point with angle
type MotionCommand struct {
	CMD     uint8
	PARAM_1 int16
	PARAM_2 int16
	PARAM_3 int16
	flags   int8
}
