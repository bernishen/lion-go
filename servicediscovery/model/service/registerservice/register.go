package registerservice

import (
	"github.com/Berni-Shen/lion-go/servicediscovery/model/dao/dbpool"
	"github.com/Berni-Shen/lion-go/servicediscovery/model/dao/entities"
	"github.com/Berni-Shen/lion-go/servicediscovery/model/service/request"
	"github.com/Berni-Shen/lion-go/utils/exception"
)

// func Register() *utils.Exception {

// }

func saveService(service *request.Service) *exception.Exception {
	db, err := dbpool.Take()
	if err != nil {
		return exception.NewException(exception.Error, 1001, err.Message)
	}

	var s *entities.Service
	db.First(s, "Name=? and Version=?", service.Name, service.Version)
	if s == nil {
		s = new(entities.Service)
		s.Name = service.Name
		s.Version = service.Version
		s.Description = service.Description
		db.Create(s)
	}

	var r *entities.Resource
	db.First(r, "ServiceID=?", s.ID)
	if r == nil {
		r = new(entities.Resource)
		r.ID = s.ID
		r.Name = service.ResouceName
		r.GetMethod = service.GetMethod
		r.PostMethod = service.PostMethod
		r.PutMethod = service.PutMethod
		r.DeleteMethod = service.DeleteMethod
		db.Create(r)
	}

	var i *entities.Instance
	db.First(i, "ServiceID=?", s.ID)
	if i == nil {
		i = new(entities.Instance)
		i.ID = s.ID
		i.Address = service.Address
		i.Port = service.Port
		db.Create(i)
	}

	return nil
}

func find() {

}
