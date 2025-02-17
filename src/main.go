package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/evilsocket/islazy/zip"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ConvertReq struct {
	Nuts    string                `form:"nuts" binding:"required"`
	Vat     string                `form:"vat" binding:"required"`
	Version string                `form:"version" binding:"required"`
	Az      string                `form:"az" binding:"required"`
	File    *multipart.FileHeader `form:"file" binding:"required"`
}

func main() {
	g := gin.Default()

	g.Use(cors.Default())

	g.GET("/health", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })

	g.POST("/", func(ctx *gin.Context) {
		var recObj ConvertReq
		log.Println("Incoming request")

		if err := ctx.ShouldBind(&recObj); err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
			return
		}

		workdir, err := os.MkdirTemp("", "gtfs2netex.")
		if err != nil {
			log.Panic("Could not create temp directory")
		}
		defer os.RemoveAll(workdir)

		gtfs := filepath.Join(workdir, "in.gtfs")
		if err := ctx.SaveUploadedFile(recObj.File, gtfs); err != nil {
			log.Panic(err)
		}

		_, err = zip.Unzip(gtfs, workdir)
		if err != nil {
			log.Panic(err)
		}

		// launch the external converter tool
		cmd := exec.Command("python3", os.Args[1], "--folder", workdir, "--db", "tmp", "--NUTS", recObj.Nuts, "--vat", recObj.Vat, "--version", recObj.Version, "--az", recObj.Az)
		output, err := cmd.Output()
		if err != nil {
			log.Panic(err)
		}
		fmt.Print(string(output))

		e, _ := os.ReadDir(workdir)
		for _, e := range e {
			fmt.Println(e)
		}

		matches, err := filepath.Glob(filepath.Join(workdir, "*.xml"))
		if err != nil {
			log.Panic(err)
		}

		if len(matches) != 1 {
			log.Panic("Could not find output file. Check stdout output,  maybe the converter failed without setting exit code", matches)
		}

		f, err := os.Open(matches[0])
		if err != nil {
			log.Panic(err)
		}

		netex, err := io.ReadAll(f)
		if err != nil {
			log.Panic(err)
		}

		ctx.Data(http.StatusOK, "application/xml", netex)
		log.Println("Incoming request")
	})

	log.Println("gtfs2netex-api starting up an ready to serve...")

	g.Run()
}
