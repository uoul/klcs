package api

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (a *Api) getPrintJobs(ctx *gin.Context) {
	printerId := ctx.Param("printerId")
	jobs := a.printService.Register(printerId)
	defer a.printService.UnRegister(printerId, jobs)
	ticker := time.NewTicker(10 * time.Second)
	ctx.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			ctx.SSEvent("heartbeat", nil)
			return true
		case job := <-jobs:
			if job.Error == nil {
				ctx.SSEvent("message", job.Value)
			}
			return true
		}
	})
}

func (a *Api) acknowledgePrintJob(ctx *gin.Context) {
	printerId := ctx.Param("printerId")
	transactionId := ctx.Param("transactionId")
	if err := a.logic.AcknowledgePrintJob(ctx, printerId, transactionId); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
