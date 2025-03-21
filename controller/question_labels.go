package controller

import (
	"github.com/gin-gonic/gin"
	"java2go/entity"
	"java2go/mapper"
	"java2go/utils"
	"net/http"
)

// 处理 /getAllQuestionLabels 请求
func GetAllQuestionLabels(c *gin.Context) {
	var questionLabels []entity.QuestionLabels
	mapper.DB.Find(&questionLabels)
	resp := utils.Make200Resp("successfully get all question labels", questionLabels)
	c.String(http.StatusOK, resp)
}

// 处理 /getDistinctChapter1 请求
func GetDistinctChapter1(c *gin.Context) {
	var distinctChapter1 []entity.QuestionLabels
	mapper.DB.Distinct("chapter_1").Find(&distinctChapter1)
	var chapter1List []string
	for _, label := range distinctChapter1 {
		chapter1List = append(chapter1List, label.Chapter1)
	}
	resp := utils.Make200Resp("successfully get chapter1", chapter1List)
	c.String(http.StatusOK, resp)
}

// 处理 /getDistinctChapter2 请求
func GetDistinctChapter2(c *gin.Context) {
	var distinctChapter2 []entity.QuestionLabels
	mapper.DB.Distinct("chapter_2").Find(&distinctChapter2)
	var chapter2List []string
	for _, label := range distinctChapter2 {
		chapter2List = append(chapter2List, label.Chapter2)
	}
	resp := utils.Make200Resp("successfully get chapter2", chapter2List)
	c.String(http.StatusOK, resp)
}

// 处理 /getChapter2ByChapter1 请求
func GetChapter2ByChapter1(c *gin.Context) {
	chapter1 := c.Query("chapter1")
	var chapter2ByChapter1 []entity.QuestionLabels
	mapper.DB.Where("chapter_1 = ?", chapter1).Find(&chapter2ByChapter1)
	var chapter2List []string
	for _, label := range chapter2ByChapter1 {
		chapter2List = append(chapter2List, label.Chapter2)
	}
	resp := utils.Make200Resp("successfully get chapter2 by chapter1", chapter2List)
	c.String(http.StatusOK, resp)
}

// 处理 /getDistinctLabel1 请求
func GetDistinctLabel1(c *gin.Context) {
	var distinctLabel1 []entity.QuestionLabels
	mapper.DB.Distinct("label_1").Find(&distinctLabel1)
	var label1List []string
	for _, label := range distinctLabel1 {
		label1List = append(label1List, label.Label1)
	}
	resp := utils.Make200Resp("successfully get label1", label1List)
	c.String(http.StatusOK, resp)
}

// 处理 /getDistinctLabel2 请求
func GetDistinctLabel2(c *gin.Context) {
	var distinctLabel2 []entity.QuestionLabels
	mapper.DB.Distinct("label_2").Find(&distinctLabel2)
	var label2List []string
	for _, label := range distinctLabel2 {
		label2List = append(label2List, label.Label2)
	}
	resp := utils.Make200Resp("successfully get label2", label2List)
	c.String(http.StatusOK, resp)
}
