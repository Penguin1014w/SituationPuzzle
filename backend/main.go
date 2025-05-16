package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Riddle 结构体定义
type Riddle struct {
	ID         int    `json:"id"`
	TitleCH    string `json:"title_ch"`
	TitleEN    string `json:"title_en"`
	ContentCH  string `json:"content_ch"`
	ContentEN  string `json:"content_en"`
	AnswerCH   string `json:"answer_ch"`
	AnswerEN   string `json:"answer_en"`
	Difficulty int    `json:"difficulty"`
}

// 全局变量存储谜面数据
// var riddles []Riddle
var riddles map[int]Riddle

// 初始化函数，在程序启动时加载CSV数据
func init() {
	godotenv.Load() // 加载.env文件
	// 打开CSV文件
	riddles = make(map[int]Riddle)
	file, err := os.Open("riddles.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建CSV reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// 跳过表头，从第二行开始读取数据
	for i, record := range records {
		if i == 0 {
			continue // 跳过表头
		}

		// 确保行数据完整
		if len(record) >= 8 {
			id, _ := strconv.Atoi(record[0])
			difficulty, _ := strconv.Atoi(record[7])
			riddle := Riddle{
				ID:         id,
				TitleCH:    record[1],
				TitleEN:    record[2],
				ContentCH:  record[3],
				ContentEN:  record[4],
				AnswerCH:   record[5],
				AnswerEN:   record[6],
				Difficulty: difficulty,
			}
			// riddles = append(riddles, riddle)
			riddles[id] = riddle
		}
	}
}

func main() {
	// 创建Gin路由
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	// 获取所有谜面
	r.GET("/api/riddles", func(c *gin.Context) {
		lang := c.DefaultQuery("lang", "CH") // 默认中文

		// 根据语言返回对应的数据
		var responseData []map[string]interface{}
		for _, riddle := range riddles {
			item := map[string]interface{}{
				"id":         riddle.ID,
				"difficulty": riddle.Difficulty,
			}

			if lang == "CH" {
				item["title"] = riddle.TitleCH
				item["content_ch"] = riddle.ContentCH
				item["content_en"] = riddle.ContentEN
				item["answer_ch"] = riddle.AnswerCH
				item["answer_en"] = riddle.AnswerEN
			} else {
				item["title"] = riddle.TitleEN
				item["content"] = riddle.ContentEN
				item["content_ch"] = riddle.ContentCH
				item["content_en"] = riddle.ContentEN
				item["answer_ch"] = riddle.AnswerCH
				item["answer_en"] = riddle.AnswerEN
			}

			responseData = append(responseData, item)
		}

		c.JSON(http.StatusOK, responseData)
	})

	// 检查答案
	r.POST("/api/check-answer", func(c *gin.Context) {
		var request struct {
			RiddleID int    `json:"riddleId"`
			Answer   string `json:"answer"`
			Lang     string `json:"lang"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 获取对应的谜面
		// var riddle Riddle
		// for _, r := range riddles {
		// 	if r.ID == request.RiddleID {
		// 		riddle = r
		// 		break
		// 	}
		// }
		riddle, exists := riddles[request.RiddleID]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Riddle not found"})
			return
		}

		// 检查答案
		status, err := checkAnswerWithChatGPT(riddle, request.Answer, request.Lang)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": status})
	})

	// 启动服务器
	r.Run(":8080")
}
