package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"java2go/mapper"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// StreamResponse 定义流式响应的完整结构
type StreamResponse struct {
	Model      string `json:"model"`
	CreatedAt  string `json:"created_at"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason,omitempty"`
	Context    []int  `json:"context,omitempty"`
}

// Question 定义题目结构体
type Question struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
}

func main() {
	// 连接数据库
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/java_test")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	QB := mapper.NewQuestionBankMapper()
	data, err := QB.GetAllQuestionBank()
	if err != nil {
		panic(err)
	}
	spew.Dump(data)
	// 组织信息
	questionJSON, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// 定义试卷组装的需求
	prompt := fmt.Sprintf(`请从以下题目中选择部分题目生成一份试卷，包含 15 道选择题、5 道填空题和 3 道解答题。试卷总分要求100分，可以动态调整各类题目的数量，题目难度适中。题目信息：%s`, string(questionJSON))

	// 定义 Ollama API 的地址
	url := "http://localhost:11434/api/generate"

	// 构造流式请求
	reqBody := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}{
		Model:  "deepseek-r1:7b",
		Prompt: prompt,
		Stream: true,
	}

	// 发送请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(
		must(json.Marshal(reqBody)),
	))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 流式处理
	var fullResponse bytes.Buffer
	decoder := json.NewDecoder(resp.Body)

	for {
		var chunk StreamResponse
		err := decoder.Decode(&chunk)
		if err == io.EOF {
			break // 流结束
		}
		if err != nil {
			panic(fmt.Errorf("解析 JSON 块失败: %w", err))
		}

		// 打印每次接收到的数据
		fmt.Print(chunk.Response)

		fullResponse.WriteString(chunk.Response)

		// 处理结束标志
		if chunk.Done {
			fmt.Printf("\n完整响应: \n%s\n", fullResponse.String())
			fmt.Printf("元信息：结束原因 =%s, 上下文长度 =%d\n",
				chunk.DoneReason, len(chunk.Context))
			break
		}
	}
}

// getQuestionsFromDB 从数据库获取题目信息
func getQuestionsFromDB(db *sql.DB) ([]Question, error) {
	rows, err := db.Query("SELECT id, type, content, answer, difficulty FROM questions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var question Question
		err := rows.Scan(&question.ID, &question.Type, &question.Content, &question.Answer, &question.Difficulty)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
