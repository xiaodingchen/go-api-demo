package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"test.local/internal"
	"test.local/pkg/utils"
	"test.local/pkg/xgin/middleware"
	"test.local/pkg/xtrace"
	"time"
)

var api = &cobra.Command{
	Use:   "api",
	Short: "api",
	Long:  "api service",
	Run: func(cmd *cobra.Command, args []string) {
		ex, _ := os.Executable()
		addr := viper.GetString("server.addr")
		log.Println("api service start, addr:", addr, "time:", time.Now().Format(utils.DateFormatTimestamp))
		overseer.Run(overseer.Config{
			Address:          addr,
			Program:          prog,
			Fetcher:          &fetcher.File{Path: ex},
			TerminateTimeout: 5 * time.Second,
			Debug:            true,
		})
	},
}

func prog(state overseer.State) {
	log.Println("api service run, addr:", state.Address, "time:", state.StartedAt.Format(utils.DateFormatTimestamp))
	debug := viper.GetBool("server.debug")
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 引入gin框架
	engine := gin.New()
	engine.Use(gin.Recovery(), middleware.Request(), middleware.Logger(debug), xtrace.TraceLogger())
	internal.Init(engine)

	srv := http.Server{
		Handler:      engine,
		ReadTimeout:  viper.GetDuration("server.readTimeout") * time.Millisecond,
		WriteTimeout: viper.GetDuration("server.writeTimeout") * time.Millisecond,
	}
	srv.Serve(state.Listener)
}