# CutVideo_API
十行笔记后端开发笔试题-视频剪辑接口

## 1、项目所需的框架和第三方库

下载gin框架依赖

```go
go get -u github.com/gin-gonic/gin
```



gorm框架下载包命令

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```



mysql

```go
go get -u github.com/go-sql-driver/mysql
```



viper

```go
go get github.com/spf13/viper
```



mac准备ffmpeg本机环境

```bash
brew install ffmpeg
```



下载ffmpeg第三方库

```go
 go get -u github.com/u2takey/ffmpeg-go
```



下载jwt包

```go
go get -u github.com/dgrijalva/jwt-go
```



测试工具包assert

```go
go get -u github.com/stretchr/testify/
```



执行单元测试用例

```go
go test -v ./test/cut_video_controller_test.go
```

![image](https://user-images.githubusercontent.com/59222555/227130573-94d07235-cb71-4d53-a950-9ccbdd1c212b.png)



程序正常运行

![image](https://user-images.githubusercontent.com/59222555/227130906-93199e53-7e8b-45f4-84d9-f567a6dd70d1.png)

postMen测试

<img width="1440" alt="image" src="https://user-images.githubusercontent.com/59222555/227130968-e955b463-910f-40c3-a64c-ab91f0525f00.png">

## 2、接口文档说明

### 剪辑视频接口

剪辑指定的视频和剪辑视频存放路径，根据输入的启始剪辑时间和持续剪辑时间，生成一个新的视频文件。并且返回开始剪辑时间和结束剪辑时间等信息、剪辑完成后的视频URL。

**请求URL：**

- `POST /cutVideo`

**请求参数：**

| 参数名         | 必选 | 类型   | 说明                     |
| -------------- | ---- | ------ | ------------------------ |
| inputVideoPath | 是   | string | 输入视频的绝对路径       |
| outputDir      | 是   | string | 输出视频文件的目录       |
| startCut       | 是   | string | 启始剪辑时间（00:00:01） |
| duration       | 是   | string | 持续剪辑时间（00:00:01） |

**请求示例：**

```go
jsonCopy code
{
    "inputVideoPath": "/Users/zyq/Desktop/test.mp4",
    "outputDir": "/Users/zyq/Desktop",
    "startCut": "0:01",
    "duration": "4"
}
```

**返回参数：**

| 参数名               | 类型   | 说明             |
| -------------------- | ------ | ---------------- |
| data.resultVideoPath | string | 剪辑后的视频路径 |
| data.startTime       | string | 剪辑开始时间     |
| data.endTime         | string | 剪辑结束时间     |
| msg                  | string | 剪辑结果提示信息 |

**返回示例：**

```go
jsonCopy code
{
    "data": {
        "endTime": "2023-03-23 12:21:22",
        "resultVideoPath": "/Users/zyq/Desktop/1d73c603-6d3c-4d7e-8d7a-8ad97be3e305.mp4",
        "startTime": "2023-03-23 12:21:21"
    },
    "msg": "剪辑成功！"
}
```

**错误返回：**

| HTTP状态码 | 返回码                      | 说明                                              |
| ---------- | --------------------------- | ------------------------------------------------- |
| 400        | missing required parameters | 缺少必要的参数                                    |
| 400        | 格式不支持                  | 输入的视频格式不支持                              |
| 500        |                             | 服务器错误，具体信息在返回结果中的 `error` 字段中 |

**错误返回示例：**

```go
jsonCopy code
{
    "error": "missing required parameters"
}
```



## 计分说明第4

应该是梯子不稳定等问题依赖下不下来这边提供代码参考

加上如下第三方依赖

```
"github.com/rylio/ytdl"
"github.com/vansante/go-ffmpeg"
"github.com/vansante/go-ffprobe"
```

接口方法进行改进

1、剪辑的那部分代码使用goroutine的方式

在剪辑视频接口中异步执行剪辑任务，同时返回任务的唯一标识

```go
go func() {
		// 剪辑视频
		err = ffmpeg.Input(inputVideoPath).
			Output(resultVideoPath, ffmpeg.KwArgs{"ss": startCut, "t": duration, "c:v": "copy", "c:a": "copy"}).
			OverWriteOutput().
			ErrorToStdOut().
			WithProgressCallback(func(progress ffmpeg.Progress) {
				// 更新任务进度
				updateTaskProgress(taskID.String(), progress)
			}).
			Run()
		if err != nil {
			updateTaskStatus(taskID.String(), "failed")
		} else {
			updateTaskStatus(taskID.String(), "completed")
		}
	}()
	// 返回任务ID
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"taskID": taskID.String(),
		},
		"msg": "剪辑任务已经开始执行！",
	})
```

2、提供一个查询任务进度的接口，通过传入任务ID来查询任务进度。任务进度可以保存在数据库或者缓存中。

```go
func GetTaskProgressController(c *gin.Context) {
    taskID := c.Query("taskID")
    progress := getTaskProgress(taskID)
    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "progress": progress,
        },
        "msg": "获取任务进度成功！",
    })
}
```

3、剪辑任务执行的进度可以通过FFmpeg提供的进度回调函数获取，进度信息可以保存在数据库或者缓存中。

```
func updateTaskProgress(taskID string, progress ffmpeg.Progress) {
    db := common.GetDB()
    task := models.Task{}
    db.First(&task, "task_id = ?", taskID)
    task.Progress = int(progress.Percent)
    db.Save(&task)
}

func getTaskProgress(taskID string) int {
    db := common.GetDB()
    task := models.Task{}
    db.First(&task, "task_id = ?", taskID)
    return task.Progress
}

```

## 计分项第5

要在剪辑视频接口完成之后主动推送剪辑完成事件给用户，可以采用消息队列的方式。具体实现步骤如下：

1. 配置消息队列服务器，比如 RabbitMQ 或者 Kafka。
2. 在剪辑视频接口中，剪辑完成后，将剪辑完成事件发送到消息队列。
3. 用户客户端通过消息队列订阅剪辑完成事件，收到剪辑完成事件后，即可进行相应处理，比如提示用户剪辑完成、展示剪辑后的视频等。

注：用户客户端订阅消息队列时，需要提供一个唯一标识符，比如用户的 ID，以便服务器将消息发送到正确的客户端。另外，在发送消息时，可以添加一些额外的信息，比如剪辑后的视频路径等，以便客户端能够方便地处理剪辑完成事件。
