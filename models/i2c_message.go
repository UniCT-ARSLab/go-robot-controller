package models

type i2cMessage struct {
	Command byte
	Val1    int16
	Val2    int16
	Val3    int16
	End     byte
}
