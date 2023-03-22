package controller

import (
	"CutVido_api/common"
	"CutVido_api/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"net/http"
	"path/filepath"
	"strings"
	//"strconv"
)

func CutVideoController(c *gin.Context) {
	db := common.GetDB()
	inputVideoPath := c.Query("inputVideoPath")
	outputDir := c.Query("outputDir")
	startCut := c.Query("startCut")
	duration := c.Query("duration")

	if inputVideoPath == "" || outputDir == "" || startCut == "" || duration == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})
		return
	}

	formatArr := []string{"mp4", "flv"}
	_, file := filepath.Split(inputVideoPath)
	tmps := strings.Split(file, ".")
	ext := tmps[len(tmps)-1]
	if !in(ext, formatArr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "格式不支持",
		})
		return
	}

	name, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	//startTime := time.Now().Unix()
	//剪切视频
	resultVideoPath := filepath.Join(outputDir, fmt.Sprintf("%s.%s", name.String(), ext))
	err = ffmpeg.Input(inputVideoPath).
		Output(resultVideoPath, ffmpeg.KwArgs{"ss": startCut, "t": duration, "c:v": "copy", "c:a": "copy"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	//endTime := time.Now().Unix()
	video := models.Video{
		Name: name.String(),
		//StartTime:      startTime,
		//EndTime:        endTime,
		InputVideoPath: inputVideoPath,
		OutputDir:      outputDir,
		StartCut:       startCut,
		Duration:       duration,
	}
	db.Create(&video)
	c.JSON(http.StatusOK, gin.H{
		"data": resultVideoPath,
		"msg":  "剪辑成功！",
	})

}

func in(target string, array []string) bool {
	for _, element := range array {
		if target == element {
			return true
		}
	}
	return false
}
