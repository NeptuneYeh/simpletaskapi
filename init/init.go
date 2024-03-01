package init

import (
	"github.com/NeptuneYeh/simpletask/init/gin"
	"github.com/NeptuneYeh/simpletask/internal/infra/db"
)

type MainInitProcess struct {
	storeModule db.Store
	ginModule   *gin.Module
}

func NewMainInitProcess() *MainInitProcess {
	// init Store
	store := db.NewStore()
	return &MainInitProcess{
		storeModule: store,
		ginModule:   gin.NewModule(store),
	}
}

// Run run gin module
func (m *MainInitProcess) Run() {
	m.ginModule.Run("0.0.0.0:8080")
}
