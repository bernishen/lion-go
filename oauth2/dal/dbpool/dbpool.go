package dbpool

import (
	"github.com/bernishen/exception"
	"time"

	// gorm use and connection to mysql.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type dbPool struct {
	bufpool chan *gorm.DB
}

const poolsize = 10

var pool dbPool
var currentCount int

func init() {
	pool.bufpool = make(chan *gorm.DB, poolsize)
	currentCount = 0
}

// Take is  connection
func Take() (*gorm.DB, *exception.Exception) {
	if currentCount >= poolsize {
		select {
		case ret := <-pool.bufpool:
			return ret, nil
		case <-time.After(time.Second * 2):
			return nil, exception.NewException(exception.Error, 1001, "time out")
		}
	} else {
		select {
		case ret := <-pool.bufpool:
			return ret, nil
		default:
			db, err := gorm.Open("mysql", "root:123@tcp(192.168.56.50:3306)/lion-sys?charset=utf8&parseTime=True&loc=Local")

			if err != nil {
				return nil, exception.NewException(exception.Error, 1002, "Create db connection error:"+err.Error())
			}
			currentCount++
			return db, nil
		}
	}
}

// Put a dbconnection into pool.
func Put(db *gorm.DB) {
	select {
	case pool.bufpool <- db:
	default:
		_ = db.Close()
	}
}
