package controller

import (
	"fmt"
	"java2go/entity"
	"java2go/mapper"
	"java2go/services"
	"java2go/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// QuestionBankMapper 定义数据库操作接口
type QuestionBankMapper interface {
	GetAllQuestionBank() ([]entity.QuestionBank, error)
	GetQuestionBankById(id int) ([]entity.QuestionBank, error)
	GetDistinctTopicType() ([]string, error)
	SearchQuestionByTopic(topicType, keyword string) ([]entity.QuestionBank, error)
	InsertSingleQuestionBank(questionBank *entity.QuestionBank) (int64, error)
	DeleteSingleQuestionBank(id int) (int64, error)
	UpdateSingleQuestionBank(questionBank *entity.QuestionBank) (int64, error)
	GetAvgDifficultyByIds(ids []int) (float64, error)
	GetDistinctLabel1FromQuestionBank() ([]string, error)
	GetQuestionBankCountByLabel1(label1 string) (int, error)
	GetDistinctScoreFromQuestionBank() ([]float64, error)
	GetQuestionBankCountByScore(score float64) (int, error)
	GetQuestionBankByIds(ids []int, generateRange []string) ([]entity.QuestionBank, error)
	GetAll() ([]entity.QuestionBank, error)
}

// QuestionBankController 定义问题银行控制器结构体
type QuestionBankController struct {
	mapper         QuestionBankMapper
	default200Resp string
}

// NewQuestionBankController 创建新的问题银行控制器
func NewQuestionBankController() *QuestionBankController {
	return &QuestionBankController{
		mapper:         mapper.NewQuestionBankMapper(),
		default200Resp: "default 200 response",
	}
}

// GetAllQuestionBank 获取所有问题银行记录
func (c *QuestionBankController) GetAllQuestionBank(ctx *gin.Context) {
	allQuestionBank, _ := c.mapper.GetAllQuestionBank()
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, allQuestionBank))
}

// GetQuestionBank 获取指定范围的问题银行记录
func (c *QuestionBankController) GetQuestionBank(ctx *gin.Context) {
	//startItemStr := ctx.Query("startItem")
	//endItemStr := ctx.Query("endItem")
	//startItem, _ := strconv.Atoi(startItemStr)
	//endItem, _ := strconv.Atoi(endItemStr)
	questionBankList, _ := c.mapper.GetAllQuestionBank()
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, questionBankList))
}

// GetTopicType 获取不同的主题类型
func (c *QuestionBankController) GetTopicType(ctx *gin.Context) {
	topicType, _ := c.mapper.GetDistinctTopicType()
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, topicType))
}

// SearchQuestionByTopic 根据主题类型和关键字搜索问题
func (c *QuestionBankController) SearchQuestionByTopic(ctx *gin.Context) {
	topicType := ctx.Query("topicType")
	keyword := ctx.Query("keyword")
	questions, _ := c.mapper.SearchQuestionByTopic(topicType, keyword)
	ctx.JSON(http.StatusOK, utils.Make200Resp(c.default200Resp, questions))
}

type QuestionBank struct {
	ID              string    `json:"id"`
	Topic           string    `json:"topic"`
	TopicMaterialID int       `json:"topic_material_id"`
	Answer          string    `json:"answer"`
	TopicType       string    `json:"topic_type"`
	Score           float64   `json:"score"`
	Difficulty      string    `json:"difficulty"`
	Chapter1        string    `json:"chapter_1"`
	Chapter2        string    `json:"chapter_2"`
	Label1          string    `json:"label_1"`
	Label2          string    `json:"label_2"`
	UpdateTime      time.Time `json:"update_time"`
}

// InsertSingleQuestionBank 插入单条问题银行记录
func (c *QuestionBankController) InsertSingleQuestionBank(ctx *gin.Context) {
	var questionBank QuestionBank
	if err := ctx.BindJSON(&questionBank); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	questionBank.UpdateTime = time.Now()
	diff, _ := strconv.Atoi(questionBank.Difficulty)
	commitStatus, _ := c.mapper.InsertSingleQuestionBank(&entity.QuestionBank{
		Topic:           questionBank.Topic,
		TopicMaterialID: questionBank.TopicMaterialID,
		Answer:          questionBank.Answer,
		TopicType:       questionBank.TopicType,
		Score:           questionBank.Score,
		Difficulty:      diff,
		Chapter1:        questionBank.Chapter1,
		Chapter2:        questionBank.Chapter2,
		Label1:          questionBank.Label1,
		Label2:          questionBank.Label2,
		UpdateTime:      questionBank.UpdateTime,
	})
	retJson := map[string]interface{}{
		"insertStatus": commitStatus,
		"insertObject": questionBank,
	}
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, retJson))
}

// DeleteSingleQuestionBank 删除单条问题银行记录
func (c *QuestionBankController) DeleteSingleQuestionBank(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, _ := strconv.Atoi(idStr)
	commitStatus, _ := c.mapper.DeleteSingleQuestionBank(id)
	retJson := map[string]interface{}{
		"deleteStatus": commitStatus,
		"deleteObject": id,
	}
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, retJson))
}

// GetQuestionBankById 根据 ID 获取问题银行记录
func (c *QuestionBankController) GetQuestionBankById(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, _ := strconv.Atoi(idStr)
	questionBankByIdList, _ := c.mapper.GetQuestionBankById(id)
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, questionBankByIdList))
}

// UpdateQuestionBankById 根据 ID 更新问题银行记录
func (c *QuestionBankController) UpdateQuestionBankById(ctx *gin.Context) {
	var questionBank QuestionBank
	if err := ctx.BindJSON(&questionBank); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	questionBank.UpdateTime = time.Now()
	diff, _ := strconv.Atoi(questionBank.Difficulty)
	id, _ := strconv.Atoi(questionBank.ID)
	updateStatus, _ := c.mapper.UpdateSingleQuestionBank(&entity.QuestionBank{
		ID:              id,
		Topic:           questionBank.Topic,
		TopicMaterialID: questionBank.TopicMaterialID,
		Answer:          questionBank.Answer,
		TopicType:       questionBank.TopicType,
		Score:           questionBank.Score,
		Difficulty:      diff,
		Chapter1:        questionBank.Chapter1,
		Chapter2:        questionBank.Chapter2,
		Label1:          questionBank.Label1,
		Label2:          questionBank.Label2,
		UpdateTime:      questionBank.UpdateTime,
	})
	retJson := map[string]interface{}{
		"updateStatus": updateStatus,
		"updateObject": questionBank,
	}
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, retJson))
}

// UploadFile 上传 Excel 文件到数据库
func (c *QuestionBankController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	isDeleteAllStr := ctx.PostForm("isDeleteAll")
	isDeleteAll, _ := strconv.ParseBool(isDeleteAllStr)

	deleteCount := 0
	insertCount := 0
	if isDeleteAll {
		allQuestionBank, _ := c.mapper.GetAllQuestionBank()
		for _, questionBank := range allQuestionBank {
			num, _ := c.mapper.DeleteSingleQuestionBank(questionBank.ID)
			deleteCount += int(num)
		}
	}
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	eR := services.NewExcelReader(src)
	questionBanMap, err := eR.ReadExcel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, v := range questionBanMap {
		questionBank := &entity.QuestionBank{
			Topic:           v["topic"].(string),
			TopicMaterialID: toInt(v["topic_material_id"].(string)),
			Answer:          v["answer"].(string),
			TopicType:       v["topic_type"].(string),
			Score:           v["score"].(float64),
			Difficulty:      int(v["difficulty"].(int64)),
			Chapter1:        v["chapter_1"].(string),
			Chapter2:        v["chapter_2"].(string),
			Label1:          v["label_1"].(string),
			Label2:          v["label_2"].(string),
			UpdateTime:      time.Now(),
		}
		num, _ := c.mapper.InsertSingleQuestionBank(questionBank)
		insertCount += int(num)
	}

	rs := map[string]interface{}{
		"deleteCount": deleteCount,
		"insertCount": insertCount,
	}
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, rs))
}

// GetEachChapterCount 获取各 Label1 下的统计数量
func (c *QuestionBankController) GetEachChapterCount(ctx *gin.Context) {
	distinctLabel1FromQuestionBank, _ := c.mapper.GetDistinctLabel1FromQuestionBank()
	var ret []map[string]interface{}
	for _, eachLabel1 := range distinctLabel1FromQuestionBank {
		num, _ := c.mapper.GetQuestionBankCountByLabel1(eachLabel1)
		count := num
		tmp := map[string]interface{}{
			"label_1": eachLabel1,
			"count":   count,
		}
		ret = append(ret, tmp)
	}
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, ret))
}

// GetEachScoreCount 获取各 Score 下的统计数量
func (c *QuestionBankController) GetEachScoreCount(ctx *gin.Context) {
	distinctScoreFromQuestionBank, _ := c.mapper.GetDistinctScoreFromQuestionBank()
	var ret []map[string]interface{}
	for _, eachScore := range distinctScoreFromQuestionBank {
		num, _ := c.mapper.GetQuestionBankCountByScore(eachScore)
		count := num
		tmp := map[string]interface{}{
			"score": eachScore,
			"count": count,
		}
		ret = append(ret, tmp)
	}
	ctx.String(http.StatusOK, utils.Make200Resp(c.default200Resp, ret))
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func toFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
