package services

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ================== 领域对象定义 ==================
type QuestionBank struct {
	ID              int
	Topic           string
	TopicMaterialID int
	Answer          string
	TopicType       string
	Score           float64
	Difficulty      float64
	Chapter1        string
	Chapter2        string
	Label1          string
	Label2          string
}

type TestPaperGenHistory struct {
	TestPaperUID      string
	TestPaperName     string
	QuestionCount     int
	AverageDifficulty float64
	UpdateTime        time.Time
	Username          string
}

type QuestionGenHistory struct {
	TestPaperUID    string
	TestPaperName   string
	QuestionBankID  int
	Topic           string
	TopicMaterialID int
	Answer          string
	TopicType       string
	Score           float64
	Difficulty      float64
	Chapter1        string
	Chapter2        string
	Label1          string
	Label2          string
	UpdateTime      time.Time
}

// ================== 核心业务逻辑 ==================
type WordGenerator struct {
	DBRepository DBRepository // 数据库操作接口
}

func NewWordGenerator(repo DBRepository) *WordGenerator {
	return &WordGenerator{DBRepository: repo}
}

func (wg *WordGenerator) GenerateTestPaper(
	questions []QuestionBank,
	paperName string,
	username string,
) (string, error) {

	// 生成临时文件
	filePath, err := wg.generateWordFile(questions, paperName)
	if err != nil {
		return "", fmt.Errorf("生成文档失败: %w", err)
	}

	// 记录生成历史
	if err := wg.logGenerationHistory(questions, paperName, username); err != nil {
		os.Remove(filePath) // 清理文件
		return "", fmt.Errorf("记录历史失败: %w", err)
	}

	return filePath, nil
}

// ================== 模板处理逻辑 ==================
const templateName = "microcomputer_template.tmpl"

func (wg *WordGenerator) generateWordFile(questions []QuestionBank, paperName string) (string, error) {
	// 准备模板数据
	tmplData := struct {
		TotalScore  string
		TotalNumber int
		Questions   []template.HTML
	}{
		Questions: make([]template.HTML, 0, len(questions)),
	}

	totalScore := 0.0
	for i, q := range questions {
		totalScore += q.Score
		tmplData.Questions = append(tmplData.Questions,
			template.HTML(fmt.Sprintf("%d、%s<w:br />", i+1, q.Topic)))
	}

	tmplData.TotalScore = fmt.Sprintf("%.1f", totalScore)
	tmplData.TotalNumber = len(questions)

	// 创建临时文件
	tmpFile, err := os.CreateTemp(tempDir(), "testpaper-*.xml")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// 解析执行模板
	tmpl := template.Must(template.New(templateName).ParseFiles(
		filepath.Join(templateDir(), templateName)),
	)
	if err := tmpl.Execute(tmpFile, tmplData); err != nil {
		return "", fmt.Errorf("模板渲染错误: %w", err)
	}

	return tmpFile.Name(), nil
}

// ================== 历史记录处理 ==================
func (wg *WordGenerator) logGenerationHistory(
	questions []QuestionBank,
	paperName string,
	username string,
) error {
	now := time.Now()
	paperUID := generatePaperUID(now)

	// 试卷历史
	paperHistory := TestPaperGenHistory{
		TestPaperUID:      paperUID,
		TestPaperName:     paperName,
		QuestionCount:     len(questions),
		AverageDifficulty: calculateAverageDifficulty(questions),
		UpdateTime:        now,
		Username:          username,
	}

	// 题目历史
	questionHistories := make([]QuestionGenHistory, len(questions))
	for i, q := range questions {
		questionHistories[i] = QuestionGenHistory{
			TestPaperUID:    paperUID,
			TestPaperName:   paperName,
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
			UpdateTime:      now,
		}
	}

	// 数据库事务操作
	return wg.DBRepository.Transaction(func(repo DBRepository) error {
		if err := repo.CreateTestPaperHistory(paperHistory); err != nil {
			return err
		}
		return repo.BatchCreateQuestionHistories(questionHistories)
	})
}

// ================== 工具函数 ==================
func tempDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("TEMP")
	}
	return "/tmp"
}

func templateDir() string {
	return filepath.Join(os.Getenv("APP_HOME"), "templates")
}

func generatePaperUID(t time.Time) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d-%d", t.UnixNano(), rand.Intn(1000))))
	return fmt.Sprintf("%x", hash)
}

func calculateAverageDifficulty(questions []QuestionBank) float64 {
	sum := 0.0
	for _, q := range questions {
		sum += q.Difficulty
	}
	return sum / float64(len(questions))
}

// ================== 数据库接口定义 ==================
type DBRepository interface {
	Transaction(func(DBRepository) error) error
	CreateTestPaperHistory(TestPaperGenHistory) error
	BatchCreateQuestionHistories([]QuestionGenHistory) error
}

// ================== 使用示例 ==================
/*
func main() {
	// 初始化数据库仓库
	repo := NewSQLRepository() // 需要实现DBRepository接口

	// 创建生成器实例
	generator := NewWordGenerator(repo)

	// 准备测试数据
	questions := []QuestionBank{
		{
			ID:      1,
			Topic:   "8086处理器寻址方式",
			Score:   10,
			// ...其他字段
		},
	}

	// 执行生成
	filePath, err := generator.GenerateTestPaper(questions, "期中试卷", "teacher_zhang")
	if err != nil {
		panic(err)
	}

	fmt.Println("生成文件路径:", filePath)
}
*/
