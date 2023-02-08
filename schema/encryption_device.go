package schema

type EncryptionDevice struct {
	DeviceId   int64  `json:"device_id"`
	SessionId  int    `json:"session_id"`
	LaunchCode int    `json:"launch_code"`
	Location   string `json:"location"`
}
