package controller

import (
	"fmt"
	"java2go/entity"
	"java2go/mapper"

	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenWord 表示生成 Word 文件的服务
type GenWord struct{}

// GeneticIteration 表示遗传算法迭代服务
type GeneticIteration struct {
	IterationsNum    int
	QuestionBankList []entity.QuestionBank
	TargetDifficulty float64
	TKTCount         int
	XZTCount         int
	PDTCount         int
	JDTCount         int
	TKTCurrent       []entity.QuestionBank
	XZTCurrent       []entity.QuestionBank
	PDTCurrent       []entity.QuestionBank
	JDTCurrent       []entity.QuestionBank
	Variance         float64
}

// RandomSelectTopic 表示随机选题服务
type RandomSelectTopic struct{}

// WordExport 表示 Word 导出服务
type WordExport struct {
	Map map[string]string
}

var db *gorm.DB

// 处理 /RandomSelect 请求
func RandomSelect(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	selectedTopicIds := getIntList(payload, "selectedTopicIds")
	generateRange := getStringList(payload, "generateRange")
	TKTCount := int(getFloat64(payload, "TKTCount"))
	XZTCount := int(getFloat64(payload, "XZTCount"))
	PDTCount := int(getFloat64(payload, "PDTCount"))
	JDTCount := int(getFloat64(payload, "JDTCount"))
	averageDifficulty := getFloat64(payload, "averageDifficulty")

	var randomQuestionBankList []entity.QuestionBank
	mapper.DB.Where("id NOT IN ? AND topic_type IN ?", selectedTopicIds, generateRange).Find(&randomQuestionBankList)

	TKTRandomList := filterQuestions(randomQuestionBankList, "填空题")
	XZTRandomList := filterQuestions(randomQuestionBankList, "选择题")
	PDTRandomList := filterQuestions(randomQuestionBankList, "判断题")
	JDTRandomList := filterQuestions(randomQuestionBankList, "程序设计题", "程序阅读题")

	randomSelectTopic := RandomSelectTopic{}
	TKTList := randomSelectTopic.randomSelectTopic(TKTRandomList, averageDifficulty, TKTCount)
	XZTList := randomSelectTopic.randomSelectTopic(XZTRandomList, averageDifficulty, XZTCount)
	PDTList := randomSelectTopic.randomSelectTopic(PDTRandomList, averageDifficulty, PDTCount)
	JDTList := randomSelectTopic.randomSelectTopic(JDTRandomList, averageDifficulty, JDTCount)

	response := map[string]interface{}{
		"TKTList": TKTList,
		"XZTList": XZTList,
		"PDTList": PDTList,
		"JDTList": JDTList,
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /geneticSelect 请求
func GeneticSelect(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	selectedTopicIds := getIntList(payload, "selectedTopicIds")
	generateRange := getStringList(payload, "generateRange")
	TKTCount := int(getFloat64(payload, "TKTCount"))
	XZTCount := int(getFloat64(payload, "XZTCount"))
	PDTCount := int(getFloat64(payload, "PDTCount"))
	JDTCount := int(getFloat64(payload, "JDTCount"))
	targetDifficulty := getFloat64(payload, "averageDifficulty")
	iterationsNum := int(getFloat64(payload, "iterationsNum"))

	var randomQuestionBankList []entity.QuestionBank
	mapper.DB.Where("id NOT IN ? AND topic_type IN ?", selectedTopicIds, generateRange).Find(&randomQuestionBankList)

	gi := GeneticIteration{
		IterationsNum:    iterationsNum,
		QuestionBankList: randomQuestionBankList,
		TargetDifficulty: targetDifficulty,
		TKTCount:         TKTCount,
		XZTCount:         XZTCount,
		PDTCount:         PDTCount,
		JDTCount:         JDTCount,
	}
	gi.run()

	response := map[string]interface{}{
		"TKTList":  gi.TKTCurrent,
		"XZTList":  gi.XZTCurrent,
		"PDTList":  gi.PDTCurrent,
		"JDTList":  gi.JDTCurrent,
		"variance": gi.Variance,
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /questionGen 请求
func QuestionGen(c *gin.Context) {
	username := c.GetHeader("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is missing"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	questionIdList := getIntList(payload, "questionIdList")
	TKTIdList := getIntList(payload, "TKTIdList")
	XZTIdList := getIntList(payload, "XZTIdList")
	PDTIdList := getIntList(payload, "PDTIdList")
	JDTIdList := getIntList(payload, "JDTIdList")
	testPaperName := getString(payload, "testPaperName")

	questionBanks := getQuestionBanks(questionIdList, TKTIdList, XZTIdList, PDTIdList, JDTIdList)
	if len(questionBanks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid questions found"})
		return
	}

	genWord := GenWord{}
	genWord.genWordTest(questionBanks, testPaperName, username)

	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
	}
	c.JSON(http.StatusOK, response)
}

// 处理 /questionGen2 请求
func QuestionGen2(c *gin.Context) {
	username := c.GetHeader("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is missing"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	questionIdList := getIntList(payload, "questionIdList")
	TKTIdList := getIntList(payload, "TKTIdList")
	XZTIdList := getIntList(payload, "XZTIdList")
	PDTIdList := getIntList(payload, "PDTIdList")
	JDTIdList := getIntList(payload, "JDTIdList")
	testPaperName := getString(payload, "testPaperName")

	questionBanks := getQuestionBanks(questionIdList, TKTIdList, XZTIdList, PDTIdList, JDTIdList)
	if len(questionBanks) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid questions found"})
		return
	}

	var totalScore float64
	totalCount := 0
	contents := ""
	for _, q := range questionBanks {
		totalScore += q.Score
		totalCount++
		contents = fmt.Sprintf("%s%d、（本题%d分）%s\r\r", contents, totalCount, q.Score, q.Topic)
	}

	mapData := map[string]string{
		"total_score": fmt.Sprintf("%d", totalScore),
		"total_count": fmt.Sprintf("%d", totalCount),
		"contents":    contents,
	}

	we := WordExport{Map: mapData}
	file := we.exportTestPaper(1)
	logHistory(questionBanks, testPaperName, username, file)
	downloadFile(c, file)
}

// 处理 /getFile 请求
func GetFile(c *gin.Context) {
	genWord := GenWord{}
	file := genWord.getFile()
	downloadFile(c, file)
}

// 从数据库中获取题目列表
func getQuestionBanks(ids ...[]int) []entity.QuestionBank {
	var allIds []int
	for _, idList := range ids {
		allIds = append(allIds, idList...)
	}

	var questionBanks []entity.QuestionBank
	mapper.DB.Where("id IN ?", allIds).Find(&questionBanks)
	return questionBanks
}

// 过滤题目列表
func filterQuestions(questions []entity.QuestionBank, topicTypes ...string) []entity.QuestionBank {
	var filtered []entity.QuestionBank
	for _, q := range questions {
		for _, t := range topicTypes {
			if q.TopicType == t {
				filtered = append(filtered, q)
				break
			}
		}
	}
	return filtered
}

// 记录历史记录
func logHistory(questionBanks []entity.QuestionBank, testPaperName, username string, file *os.File) {
	date := time.Now()
	uid := fmt.Sprintf("%s_%s_%d", file.Name(), uuid.New().String(), date.Unix())

	testPaperGenHistory := entity.TestPaperGenHistory{
		TestPaperUID:      uid,
		TestPaperName:     testPaperName,
		QuestionCount:     len(questionBanks),
		AverageDifficulty: calculateAverageDifficulty(questionBanks),
		UpdateTime:        date,
		Username:          username,
	}
	mapper.DB.Create(&testPaperGenHistory)

	var questionGenHistoryList []entity.QuestionGenHistory
	for _, q := range questionBanks {
		questionGenHistory := entity.QuestionGenHistory{
			TestPaperUID:    uid,
			TestPaperName:   testPaperName,
			QuestionBankID:  q.ID,
			Topic:           q.Topic,
			TopicMaterialID: q.TopicMaterialID,
			Answer:          q.Answer,
			TopicType:       q.TopicType,
			Score:           q.Score,
			Difficulty:      q.Difficulty,
			Chapter1:        q.Chapter1,
			Chapter2:        q.Chapter2,
			Label1:          q.Label1,
			Label2:          q.Label2,
			UpdateTime:      date,
		}
		questionGenHistoryList = append(questionGenHistoryList, questionGenHistory)
	}
	mapper.DB.Create(&questionGenHistoryList)
}

// 计算平均难度
func calculateAverageDifficulty(questionBanks []entity.QuestionBank) float64 {
	if len(questionBanks) == 0 {
		return 0
	}
	var totalDifficulty int
	for _, q := range questionBanks {
		totalDifficulty += q.Difficulty
	}
	return float64(totalDifficulty) / float64(len(questionBanks))
}

// 下载文件
func downloadFile(c *gin.Context, file *os.File) {
	if file == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting file stats: %v", err)})
		return
	}

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=TestPaperExport_%d.docx", time.Now().Unix()))
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Last-Modified", time.Now().String())
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))
	c.Header("Content-Type", "application/octet-stream")

	c.File(file.Name())
}

// 从 map 中获取 int 列表
func getIntList(m map[string]interface{}, key string) []int {
	var list []int
	if val, ok := m[key].([]interface{}); ok {
		for _, v := range val {
			if num, ok := v.(float64); ok {
				list = append(list, int(num))
			}
		}
	}
	return list
}

// 从 map 中获取 string 列表
func getStringList(m map[string]interface{}, key string) []string {
	var list []string
	if val, ok := m[key].([]interface{}); ok {
		for _, v := range val {
			if str, ok := v.(string); ok {
				list = append(list, str)
			}
		}
	}
	return list
}

// 从 map 中获取 float64 类型的值
func getFloat64(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0
}

// 从 map 中获取 string 类型的值
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// 随机选题方法
func (r *RandomSelectTopic) randomSelectTopic(questions []entity.QuestionBank, difficulty float64, count int) []entity.QuestionBank {
	// 这里需要实现具体的随机选题逻辑
	return []entity.QuestionBank{}
}

// 遗传算法迭代方法
func (gi *GeneticIteration) run() {
	// 这里需要实现具体的遗传算法迭代逻辑
}

// 生成 Word 文件方法
func (gw *GenWord) genWordTest(questionBanks []entity.QuestionBank, testPaperName, username string) {
	// 这里需要实现具体的生成 Word 文件逻辑
}

// 获取文件方法
func (gw *GenWord) getFile() *os.File {
	// 这里需要实现具体的获取文件逻辑
	return nil
}

// 导出试卷方法
func (we *WordExport) exportTestPaper(num int) *os.File {
	// 这里需要实现具体的导出试卷逻辑
	return nil
}
