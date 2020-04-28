package setup

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/ajay340/SearchBreaches.me/database"
	"github.com/ajay340/SearchBreaches.me/echofw"
	"github.com/ajay340/SearchBreaches.me/populate"
	"gopkg.in/robfig/cron.v2"
)

func GetCurrentDirectoryPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func SetDBPath() {
	_, filename, _, _ := runtime.Caller(1)
	d := path.Join(path.Dir(filename), "../../postgres_config.json")
	database.SetDBPath(d)
}

func Setup() {
	SetDBPath()
	cJob := cron.New()
	cJob.AddFunc("@daily", func() {
		populate.PopulateWKBreaches()
		fmt.Println("Finished updating Wikipedia breaches")
	})
	cJob.AddFunc("@daily", func() {
		populate.PopulateHIBPBreaches()
		fmt.Println("Finished updating HaveIbeenPwned breaches")
	})
	cJob.AddFunc("@daily", func() {
		database.DeleteExpireSessions()
		fmt.Println("Finished deleting expired session cookies")
	})
	echofw.StartServer(GetCurrentDirectoryPath())
}
