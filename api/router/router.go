package router

import (
	"fmt"

	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/jclebreton/opensirene/api/views"
	"github.com/jclebreton/opensirene/conf"
)

// SetupAndRun creates the router and runs it
func SetupAndRun(gormClient *gorm.DB) error {
	// Create the router
	r := gin.Default()

	// Setup debug mode or not in Gin
	if !conf.C.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup CORS configuration
	if conf.C.Server.Cors.Enabled {
		cc := cors.Config{
			AllowMethods:  []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
			AllowHeaders:  []string{"DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Content-Range", "Range", "Authorization"},
			ExposeHeaders: []string{"DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Content-Range", "Range", "Authorization"},
		}
		if conf.C.Server.Cors.PermissiveMode {
			cc.AllowAllOrigins = true
		} else {
			cc.AllowOrigins = conf.C.Server.Cors.AllowOrigins
		}
		r.Use(cors.New(cc))
	}
	p := ginprom.New(ginprom.Subsystem(conf.C.Prometheus.Prefix), ginprom.Engine(r))
	r.Use(p.Instrument())

	// Route setup
	views := &views.ViewsContext{GormClient: gormClient}
	r.GET("/siret/:id", views.GetSiret)
	r.GET("/siren/:id", views.GetSiren)
	r.GET("/history", views.GetHistory)

	// Run the server
	logrus.WithFields(logrus.Fields{"port": conf.C.Server.Port, "host": conf.C.Server.Host}).Info("Starting server")
	return r.Run(fmt.Sprintf("%s:%d", conf.C.Server.Host, conf.C.Server.Port))
}
