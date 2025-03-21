package controller

import (
	"fmt"
	"java2go/entity"
	"java2go/mapper"
	"java2go/services"
	"java2go/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 处理 /getQuestionGenHistoriesByTestPaperUid 请求
func GetQuestionGenHistoriesByTestPaperUid(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	var questionGenHistories []entity.QuestionGenHistory
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Find(&questionGenHistories)
	resp := utils.Make200Resp("Success", questionGenHistories)
	c.String(http.StatusOK, resp)
}

// 处理 /deleteQuestionGenHistoryByTestPaperUid 请求
func DeleteQuestionGenHistoryByTestPaperUid(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	var delQuestionCount int64
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Delete(&entity.QuestionGenHistory{}).Count(&delQuestionCount)
	var delTestPaperCount int64
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Delete(&entity.TestPaperGenHistory{}).Count(&delTestPaperCount)
	response := map[string]interface{}{
		"delQuestionCount":  delQuestionCount,
		"delTestPaperCount": delTestPaperCount,
	}
	resp := utils.Make200Resp("Success", response)
	c.String(http.StatusOK, resp)
}

// 处理 /updateQuestionGenHistory 请求
func UpdateQuestionGenHistory(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	questionBankIdsStr := c.QueryArray("question_bank_id")
	var questionBankIds []int
	for _, idStr := range questionBankIdsStr {
		var id int
		fmt.Sscanf(idStr, "%d", &id)
		questionBankIds = append(questionBankIds, id)
	}

	var names []string
	mapper.DB.Table("test_paper_gen_histories").Where("test_paper_uid = ?", testPaperUid).Pluck("test_paper_name", &names)
	var testPaperName string
	if len(names) > 0 {
		testPaperName = names[0]
	}

	var questions []entity.QuestionBank
	for _, id := range questionBankIds {
		var question entity.QuestionBank
		if err := mapper.DB.Where("id = ?", id).First(&question).Error; err == nil {
			questions = append(questions, question)
		}
	}

	date := time.Now()
	var questionGenHistories []entity.QuestionGenHistory
	for _, q := range questions {
		questionGenHistory := entity.QuestionGenHistory{
			TestPaperUID:   testPaperUid,
			TestPaperName:  testPaperName,
			QuestionBankID: q.ID,
			Topic:          q.Topic,
			Answer:         q.Answer,
			TopicType:      q.TopicType,
			Score:          q.Score,
			Difficulty:     q.Difficulty,
			Chapter1:       q.Chapter1,
			Chapter2:       q.Chapter2,
			Label1:         q.Label1,
			Label2:         q.Label2,
			UpdateTime:     date,
		}
		questionGenHistories = append(questionGenHistories, questionGenHistory)
	}

	// 更新时间
	updateRes := mapper.DB.Model(&entity.TestPaperGenHistory{}).
		Where("test_paper_uid = ?", testPaperUid).
		Update("update_time", date)
	// 删除旧的
	deleteRes := mapper.DB.Where("test_paper_uid = ?", testPaperUid).Delete(&entity.QuestionGenHistory{})
	// 插入新的
	insertRes := mapper.DB.Create(&questionGenHistories)

	resp := utils.Make200Resp("Success", updateRes.RowsAffected+deleteRes.RowsAffected+insertRes.RowsAffected)
	c.String(http.StatusOK, resp)
}

// 处理 /reExportTestPaper 请求
func ReExportTestPaper(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	var questionGenHistories []entity.QuestionGenHistory
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Find(&questionGenHistories)

	var questionBanks []entity.QuestionBank
	for _, item := range questionGenHistories {
		var question entity.QuestionBank
		if err := mapper.DB.Where("id = ?", item.QuestionBankID).First(&question).Error; err == nil {
			questionBanks = append(questionBanks, question)
		}
	}
	var totalScore float64
	totalCount := 0
	contents := ""
	for _, q := range questionBanks {
		totalScore += q.Score
		totalCount++
		scoreStr := formatScore(q.Score)
		contents = fmt.Sprintf("%s%d、（本题%s分）%s\r\r", contents, totalCount, scoreStr, q.Topic)
	}

	mapData := map[string]string{
		"total_score": fmt.Sprintf("%s", formatScore(totalScore)),
		"total_count": fmt.Sprintf("%d", totalCount),
		"contents":    contents,
	}
	we := services.NewWordExporter(mapData)
	file, err := we.ExportTestPaper(1)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	downloadFile(c, file)
}

// 处理 /exportAnswer 请求
func ExportAnswer(c *gin.Context) {
	testPaperUid := c.Query("test_paper_uid")
	var questionGenHistories []entity.QuestionGenHistory
	mapper.DB.Where("test_paper_uid = ?", testPaperUid).Find(&questionGenHistories)

	var questionBanks []entity.QuestionBank
	for _, item := range questionGenHistories {
		var question entity.QuestionBank
		if err := mapper.DB.Where("id = ?", item.QuestionBankID).First(&question).Error; err == nil {
			questionBanks = append(questionBanks, question)
		}
	}

	var totalScore float64
	var totalCount int
	contents := ""
	for _, q := range questionBanks {
		totalScore += q.Score
		totalCount++
		contents = fmt.Sprintf("%s%d、（本题%d分）%s\r\r", contents, totalCount, q.Score, q.Answer)
	}

	mapData := map[string]string{
		"total_score": fmt.Sprintf("%d", totalScore),
		"total_count": fmt.Sprintf("%d", totalCount),
		"contents":    contents,
	}

	we := services.NewWordExporter(mapData)
	file, _ := we.ExportTestPaper(2)
	downloadFile(c, file)
}

// 格式化分数输出
func formatScore(score float64) string {
	if score == float64(int(score)) {
		return fmt.Sprintf("%d", int(score))
	}
	return fmt.Sprintf("%.1f", score)
}
