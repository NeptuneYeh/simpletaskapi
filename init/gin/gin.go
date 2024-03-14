package gin

import (
	"context"
	"github.com/NeptuneYeh/simpletask/internal/app/controllers"
	"github.com/NeptuneYeh/simpletask/internal/infra/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Module struct {
	Router *gin.Engine
	Server *http.Server
}

func NewModule(store db.Store) *Module {
	r := gin.Default()
	ginModule := &Module{
		Router: r,
	}
	gin.ForceConsoleColor()
	ginModule.setupRoute(store)

	return ginModule
}

// setup route
func (module *Module) setupRoute(store db.Store) {
	// init controller
	taskController := controllers.NewTaskController(store)
	// add routes to router
	module.Router.POST("/tasks", taskController.CreateTask)
	module.Router.GET("/tasks/:id", taskController.GetTask)
	module.Router.GET("/tasks", taskController.ListTask)
	module.Router.PUT("/tasks/:id", taskController.UpdateTask)
	module.Router.DELETE("/tasks/:id", taskController.DeleteTask)
}

// Run gin
func (module *Module) Run(address string) {
	module.Server = &http.Server{
		Addr:    address,
		Handler: module.Router,
	}

	go func() {
		if err := module.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
}

func (module *Module) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := module.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to run Gin shutdown: %v", err)
	}
	return nil
}
