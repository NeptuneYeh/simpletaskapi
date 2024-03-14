package init

import (
	"github.com/NeptuneYeh/simpletask/init/gin"
	"github.com/NeptuneYeh/simpletask/internal/infra/db"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type MainInitProcess struct {
	storeModule db.Store
	ginModule   *gin.Module
	OsChannel   chan os.Signal
}

func NewMainInitProcess() *MainInitProcess {
	// init Store
	store := db.NewStore()
	channel := make(chan os.Signal, 1)
	return &MainInitProcess{
		storeModule: store,
		ginModule:   gin.NewModule(store),
		OsChannel:   channel,
	}
}

// Run run gin module
func (m *MainInitProcess) Run() {
	m.ginModule.Run("0.0.0.0:8080")

	// register os signal for graceful shutdown
	signal.Notify(m.OsChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-m.OsChannel
	log.Fatalf("Received signal: " + s.String())
	m.ginModule.Shutdown()
}
