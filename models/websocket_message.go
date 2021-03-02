package models

//WebSocketMessage rappresents a base I2C payload
type WebSocketMessage struct {
	Command string      `json:"command"`
	Payload interface{} `json:"data"`
}
