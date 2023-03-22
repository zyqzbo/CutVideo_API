package test

import (
	"CutVido_api/controller"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCutVideoController(t *testing.T) {
	// 构造测试用的gin.Context
	router := gin.Default()
	router.POST("/cutVideo", controller.CutVideoController)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cutVideo", nil)

	// 传递参数进行测试
	values := url.Values{}
	values.Set("inputVideoPath", "/Users/zyq/Desktop/test.mp4")
	values.Set("outputDir", "/Users/zyq/Desktop")
	values.Set("startCut", "00:00:01")
	values.Set("duration", "00:00:03")
	req.Form = values
	router.ServeHTTP(w, req)

	// 验证返回结果
	assert.Equal(t, 200, w.Code)
	var response map[string]interface{}
	// JSON 格式的响应解析为一个结构体类型 c.Writer.Body.Bytes()
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	// 检查响应是否符合预期 表达式返回false表示错误
	assert.Nil(t, err)
	assert.Equal(t, "剪辑成功！", response["msg"])
	assert.NotNil(t, response["data"])
	// 转为json串
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["resultVideoPath"])
	assert.NotNil(t, data["startTime"])
	assert.NotNil(t, data["endTime"])
}
