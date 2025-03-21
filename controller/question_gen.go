package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"java2go/entity"
	"java2go/mapper"
	"java2go/services"
	"java2go/utils"

	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 处理 /randomSelect 请求
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
	query := mapper.DB

	if len(selectedTopicIds) > 0 {
		query = query.Where("id NOT IN ?", selectedTopicIds)
	}
	if len(generateRange) > 0 {
		query = query.Where("label_1 IN ?", generateRange)
	}

	result := query.Find(&randomQuestionBankList)
	if result.Error != nil {
		fmt.Printf("查询出错: %v\n", result.Error)
		return
	}

	//var randomQuestionBankList []entity.QuestionBank
	//mapper.DB.Where("id NOT IN ? AND topic_type IN ?", selectedTopicIds, generateRange).Find(&randomQuestionBankList)

	TKTRandomList := filterQuestions(randomQuestionBankList, "填空题")
	XZTRandomList := filterQuestions(randomQuestionBankList, "选择题")
	PDTRandomList := filterQuestions(randomQuestionBankList, "判断题")
	JDTRandomList := filterQuestions(randomQuestionBankList, "程序设计题", "程序阅读题")

	TKTList := services.RandomSelectTopic(TKTRandomList, averageDifficulty, TKTCount)
	XZTList := services.RandomSelectTopic(XZTRandomList, averageDifficulty, XZTCount)
	PDTList := services.RandomSelectTopic(PDTRandomList, averageDifficulty, PDTCount)
	JDTList := services.RandomSelectTopic(JDTRandomList, averageDifficulty, JDTCount)

	response := map[string]interface{}{
		"TKTList": TKTList,
		"XZTList": XZTList,
		"PDTList": PDTList,
		"JDTList": JDTList,
	}
	resp := utils.Make200Resp("Success", response)
	c.String(http.StatusOK, resp)
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
	query := mapper.DB

	if len(selectedTopicIds) > 0 {
		query = query.Where("id NOT IN ?", selectedTopicIds)
	}
	if len(generateRange) > 0 {
		query = query.Where("label_1 IN ?", generateRange)
	}

	result := query.Find(&randomQuestionBankList)
	if result.Error != nil {
		fmt.Printf("查询出错: %v\n", result.Error)
		return
	}

	genIter := services.NewGeneticIteration(iterationsNum, randomQuestionBankList, targetDifficulty, TKTCount, XZTCount, PDTCount, JDTCount)
	genIter.Run()

	response := map[string]interface{}{
		"TKTList":  genIter.TKTCurrent,
		"XZTList":  genIter.XZTCurrent,
		"PDTList":  genIter.PDTCurrent,
		"JDTList":  genIter.JDTCurrent,
		"variance": genIter.Variance,
	}
	resp := utils.Make200Resp("Success", response)
	c.String(http.StatusOK, resp)
}

// 处理 /questionGen 请求
func QuestionGen(c *gin.Context) {
	session := sessions.Default(c)
	username, ok := session.Get("username").(string)
	fmt.Println("username: ", username)
	if !ok || username == "" {
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

	wE := services.NewWordGenerator()
	str, _ := wE.GenerateTestPaper(questionBanks, testPaperName, username)
	fmt.Println(str)

	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
	}
	c.String(http.StatusOK, utils.Make200Resp("Success", response))
}

// 处理 /questionGen2 请求
func QuestionGen2(c *gin.Context) {
	session := sessions.Default(c)
	username, ok := session.Get("username").(string)
	fmt.Println("username: ", username)
	if !ok || username == "" {
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
		scoreStr := formatScore(q.Score)
		contents = fmt.Sprintf("%s%d、（本题%s分）%s\r\r", contents, totalCount, scoreStr, q.Topic)
	}

	mapData := map[string]string{
		"total_score": fmt.Sprintf("%s", formatScore(totalScore)),
		"total_count": fmt.Sprintf("%d", totalCount),
		"contents":    contents,
	}
	wE := services.NewWordExporter(mapData)
	file, _ := wE.ExportTestPaper(1)
	logHistory(questionBanks, testPaperName, username, file)
	downloadFile(c, file)
}

// 处理 /getFile 请求
func GetFile(c *gin.Context) {
	genWord := services.NewWordGenerator()
	file := genWord.GetFile()
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
