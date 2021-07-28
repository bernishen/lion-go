package entities

// Service :  the information of service .
type Service struct {
	BaseModel
	Name        string
	Version     int
	Description string
}

// Instance : the information of service instance .
type Instance struct {
	BaseModel
	ServiceID string
	Address   string
	Port      string
}

// Resource : the entities of service resource.
type Resource struct {
	BaseModel
	ServiceID    string
	Name         string
	GetMethod    string
	PostMethod   string
	PutMethod    string
	DeleteMethod string
}
