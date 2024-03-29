package scheduler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
)

type UpdateLog interface {
	ATF() int
}

type CPUUpdate struct {
	ProcessName string `json:"Name"`
	Start       int    `json:"Start"`
	End         int    `json:"End"`
	QI          int    `json:"QI"`
	UpdateType  string `json:"Type"`
}

func (cu CPUUsage) toUpdate() *CPUUpdate {
	return &CPUUpdate{
		ProcessName: cu.ProcessName,
		Start:       cu.Start.Minute()*60 + cu.Start.Second(),
		End:         cu.End.Minute()*60 + cu.End.Second(),
		QI:          cu.QI,
		UpdateType:  "CPUUpdate",
	}
}

func (cu *CPUUpdate) ATF() int {
	return cu.Start
}

type AddProcess struct {
	CBT        int    `json:"CBT"`
	Name       string `json:"Name"`
	AT         int    `json:"AT"`
	QI         int    `json:"QI"`
	UpdateType string `json:"Type"`
}

func (proc *Process) toUpdate() *AddProcess {
	return &AddProcess{
		Name:       proc.Name,
		AT:         proc.AT.Minute()*60 + proc.AT.Second(),
		QI:         proc.QI,
		CBT:        int(proc.CBT.Seconds()),
		UpdateType: "AddProcess",
	}
}

func (proc *AddProcess) ATF() int {
	return proc.AT
}

func Display(updateChannel chan UpdateLog, wg *sync.WaitGroup) {
	updates := make([]UpdateLog, len(updateChannel))
	for i := 0; len(updateChannel) != 0; i++ {
		update, ok := <-updateChannel
		if !ok {
			panic("error in gui!")
		}
		updates[i] = update
	}

	router := gin.Default()
	router.GET("/updates", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, updates)
	})
	_ = router.Run("localhost:3131")
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
