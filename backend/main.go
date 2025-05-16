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

// save the riddles in a map
// var riddles []Riddle
var riddles map[int]Riddle

// init the riddles
func init() {
	godotenv.Load() // loding .env
	// 打开CSV文件
	riddles = make(map[int]Riddle)
	file, err := os.Open("riddles.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// create a csv reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// skip the header, start from the second row
	for i, record := range records {
		if i == 0 {
			continue // skip the header
		}

		// ensure the row data is complete
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
	// create a gin router
	r := gin.Default()

	// config the cors
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	// get all the riddles
	r.GET("/api/riddles", func(c *gin.Context) {
		lang := c.DefaultQuery("lang", "CH") // default is chinese

		// return the data according to the language
		var responseData []map[string]interface{}
		for _, riddle := range riddles {
			item := map[string]interface{}{
				"id":         riddle.ID,
				"difficulty": riddle.Difficulty,
			}

			// if lang == "CH" {
			// 	item["title"] = riddle.TitleCH
			// 	item["content_ch"] = riddle.ContentCH
			// 	item["content_en"] = riddle.ContentEN
			// 	item["answer_ch"] = riddle.AnswerCH
			// 	item["answer_en"] = riddle.AnswerEN
			// } else {
			// 	item["title"] = riddle.TitleEN
			// 	item["content"] = riddle.ContentEN
			// 	item["content_ch"] = riddle.ContentCH
			// 	item["content_en"] = riddle.ContentEN
			// 	item["answer_ch"] = riddle.AnswerCH
			// 	item["answer_en"] = riddle.AnswerEN
			// }

			if lang == "CH" {
				item["title"] = riddle.TitleCH
				item["content"] = riddle.ContentCH
				item["answer_ch"] = riddle.AnswerCH
				item["answer_en"] = riddle.AnswerEN
			} else {
				item["title"] = riddle.TitleEN
				item["content"] = riddle.ContentEN
				item["answer_ch"] = riddle.AnswerCH
				item["answer_en"] = riddle.AnswerEN
			}

			responseData = append(responseData, item)
		}

		c.JSON(http.StatusOK, responseData)
	})

	// check the answer
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

		// get the corresponding riddle
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

		// check the answer
		status, err := checkAnswerWithChatGPT(riddle, request.Answer, request.Lang)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": status})
	})

	// start the server
	r.Run(":8080")
}
