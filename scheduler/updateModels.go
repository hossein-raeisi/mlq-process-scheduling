package scheduler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sort"
	"time"
)

type UpdateLog interface {
	ATF() time.Time
}

func (cu *CPUUsage) ATF() time.Time {
	return cu.Start
}

func (proc *Process) ATF() time.Time {
	return proc.AT
}

func Display(updateChannel chan UpdateLog) {
	updates := make([]UpdateLog, len(updateChannel))
	for i := 0; len(updateChannel) != 0; i++ {
		update, ok := <-updateChannel
		if !ok {
			panic("error in gui!")
		}
		updates[i] = update
	}

	sort.Slice(updates, func(i, j int) bool {
		return updates[i].ATF().Before(updates[j].ATF())
	})

	router := gin.Default()
	router.GET("/updates", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, updates)
	})
	go router.Run("localhost:8181")
	time.Sleep(time.Second)

	time.Sleep(time.Second)
	open("localhost:8282")
}

func open(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		log.Fatal(err)
	}
}
