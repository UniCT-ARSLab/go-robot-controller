package models

//I2CMessage rappresents a base I2C payload
type I2CMessage struct {
	Command byte
	Val1    int16
	Val2    int16
	Val3    int16
	End     byte
}
