package utils

// MQBackend This is the response structure of the message queue to the backend.
type MQBackend struct {
	ServerStatus  string `json:"server_status"`
	QueuePosition uint16 `json:"queue_position"`
	Message       string `json:"message"`
}

func MQBackendWrapper(serverStatus string, QueuePosition uint16, Message string) *MQBackend {
	return &MQBackend{
		ServerStatus:  serverStatus,
		QueuePosition: QueuePosition,
		Message:       Message,
	}
}
