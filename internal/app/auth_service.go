package app

import (
	"github.com/ksopin/notes/pkg/auth"
	"github.com/ksopin/notes/pkg/db"
	"sync"
)

var (
	service *auth.Service
	serviceOnce sync.Once
)

func GetAuthService() *auth.Service {
	serviceOnce.Do(func(){
		service = auth.New(nil, db.GetPersistentDB())
	})
	return service
}
