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
	"time"
	//"strconv"
)

// CutVideoController 剪辑视频接口
func CutVideoController(c *gin.Context) {
	db := common.GetDB()
	inputVideoPath := c.PostForm("inputVideoPath") // 输入视频路径
	outputDir := c.PostForm("outputDir")           // 输出视频目录
	startCut := c.PostForm("startCut")             // 启始剪辑时间 （00:00:01）
	duration := c.PostForm("duration")             // 持续剪辑时间 （00:00:01）

	// 参数不为空的简单验证
	if inputVideoPath == "" || outputDir == "" || startCut == "" || duration == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters"})
		return
	}

	// 视频默认格式
	formatArr := []string{"mp4", "flv"}
	// 切割视频传入的格式
	_, file := filepath.Split(inputVideoPath)
	tmps := strings.Split(file, ".")
	ext := tmps[len(tmps)-1]
	// 格式验证
	if !in(ext, formatArr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "格式不支持",
		})
		return
	}

	// 最终视频名使用uuid避免重复
	name, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	startTime := time.Now().Format("2006-01-02 15:04:05")
	// 获取路径
	resultVideoPath := filepath.Join(outputDir, fmt.Sprintf("%s.%s", name.String(), ext))
	// 剪切视频
	err = ffmpeg.Input(inputVideoPath).
		Output(resultVideoPath, ffmpeg.KwArgs{"ss": startCut, "t": duration, "c:v": "copy", "c:a": "copy"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	endTime := time.Now().Format("2006-01-02 15:04:05")
	// 把视频信息保存到数据库
	video := models.Video{
		Name:           name.String(),
		StartTime:      startTime,
		EndTime:        endTime,
		InputVideoPath: inputVideoPath,
		OutputDir:      outputDir,
		StartCut:       startCut,
		Duration:       duration,
	}
	db.Create(&video)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"resultVideoPath": resultVideoPath, // 视频剪辑后的url
			"startTime":       startTime,       // 开始剪辑时间
			"endTime":         endTime,         // 结束剪辑时间
		},
		"msg": "剪辑成功！",
	})

}

// 格式验证
func in(target string, array []string) bool {
	for _, element := range array {
		if target == element {
			return true
		}
	}
	return false
}
