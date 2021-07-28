package request

// Service is the information of the service for register.
type Service struct {
	Name         string
	Version      int
	Description  string
	Address      string
	Port         string
	ResouceName  string
	GetMethod    string
	PostMethod   string
	PutMethod    string
	DeleteMethod string
}
