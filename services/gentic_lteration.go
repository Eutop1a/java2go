package services

import (
	"java2go/entity"
	"math"
	"math/rand"
	"time"
)

// GeneticIteration 遗传迭代算法结构体
type GeneticIteration struct {
	iterationsNum   int
	Variance        []float64
	targetDifficult float64
	TKTCurrent      []entity.QuestionBank
	XZTCurrent      []entity.QuestionBank
	PDTCurrent      []entity.QuestionBank
	JDTCurrent      []entity.QuestionBank
	TKTLibrary      []entity.QuestionBank
	XZTLibrary      []entity.QuestionBank
	PDTLibrary      []entity.QuestionBank
	JDTLibrary      []entity.QuestionBank
}

// NewGeneticIteration 构造函数
func NewGeneticIteration(iterationsNum int, questionList []entity.QuestionBank, targetDifficult float64, TKTCount, XZTCount, PDTCount, JDTCount int) *GeneticIteration {
	rand.Seed(time.Now().UnixNano())
	gi := &GeneticIteration{
		iterationsNum:   iterationsNum,
		targetDifficult: targetDifficult,
		Variance:        make([]float64, 0),
		TKTLibrary:      make([]entity.QuestionBank, 0),
		XZTLibrary:      make([]entity.QuestionBank, 0),
		PDTLibrary:      make([]entity.QuestionBank, 0),
		JDTLibrary:      make([]entity.QuestionBank, 0),
		TKTCurrent:      make([]entity.QuestionBank, 0),
		XZTCurrent:      make([]entity.QuestionBank, 0),
		PDTCurrent:      make([]entity.QuestionBank, 0),
		JDTCurrent:      make([]entity.QuestionBank, 0),
	}

	// 分门别类，先放到未选列表里
	for _, q := range questionList {
		switch q.TopicType {
		case "填空题":
			gi.TKTLibrary = append(gi.TKTLibrary, q)
		case "选择题":
			gi.XZTLibrary = append(gi.XZTLibrary, q)
		case "判断题":
			gi.PDTLibrary = append(gi.PDTLibrary, q)
		case "程序设计题", "程序阅读题":
			gi.JDTLibrary = append(gi.JDTLibrary, q)
		}
	}

	// 初始化已选列表
	gi.initCurrentList(&gi.TKTLibrary, &gi.TKTCurrent, TKTCount)
	gi.initCurrentList(&gi.XZTLibrary, &gi.XZTCurrent, XZTCount)
	gi.initCurrentList(&gi.PDTLibrary, &gi.PDTCurrent, PDTCount)
	gi.initCurrentList(&gi.JDTLibrary, &gi.JDTCurrent, JDTCount)

	return gi
}

// initCurrentList 初始化已选列表
func (gi *GeneticIteration) initCurrentList(library *[]entity.QuestionBank, current *[]entity.QuestionBank, count int) {
	for count > 0 && len(*library) > 0 {
		index := rand.Intn(len(*library))
		*current = append(*current, (*library)[index])
		*library = append((*library)[:index], (*library)[index+1:]...)
		count--
	}
}

// Run 运行迭代算法
func (gi *GeneticIteration) Run() {
	n := len(gi.TKTCurrent) + len(gi.XZTCurrent) + len(gi.PDTCurrent) + len(gi.JDTCurrent)
	for gi.iterationsNum > 0 {
		i := rand.Intn(n)
		if i < len(gi.TKTCurrent) {
			gi.singleIteration(&gi.TKTLibrary, &gi.TKTCurrent)
		} else if i < len(gi.TKTCurrent)+len(gi.XZTCurrent) {
			gi.singleIteration(&gi.XZTLibrary, &gi.XZTCurrent)
		} else if i < len(gi.TKTCurrent)+len(gi.XZTCurrent)+len(gi.PDTCurrent) {
			gi.singleIteration(&gi.PDTLibrary, &gi.PDTCurrent)
		} else {
			gi.singleIteration(&gi.JDTLibrary, &gi.JDTCurrent)
		}
		res := gi.calcVariance()
		gi.Variance = append(gi.Variance, res)
		gi.iterationsNum--
	}
}

// calcVariance 计算当前已选列表与预设难度的方差
func (gi *GeneticIteration) calcVariance() float64 {
	sum := 0.0
	for _, q := range gi.TKTCurrent {
		sum += math.Pow(float64(q.Difficulty)-gi.targetDifficult, 2)
	}
	for _, q := range gi.XZTCurrent {
		sum += math.Pow(float64(q.Difficulty)-gi.targetDifficult, 2)
	}
	for _, q := range gi.PDTCurrent {
		sum += math.Pow(float64(q.Difficulty)-gi.targetDifficult, 2)
	}
	for _, q := range gi.JDTCurrent {
		sum += math.Pow(float64(q.Difficulty)-gi.targetDifficult, 2)
	}
	n := len(gi.TKTCurrent) + len(gi.XZTCurrent) + len(gi.PDTCurrent) + len(gi.JDTCurrent)
	return sum / float64(n)
}

// singleIteration 单次迭代
func (gi *GeneticIteration) singleIteration(library *[]entity.QuestionBank, current *[]entity.QuestionBank) {
	if len(*library) > 0 && len(*current) > 0 {
		index1 := rand.Intn(len(*library))
		index2 := rand.Intn(len(*current))
		(*library)[index1], (*current)[index2] = (*current)[index2], (*library)[index1]
		v := gi.calcVariance()
		if len(gi.Variance) > 0 && v > gi.Variance[len(gi.Variance)-1] {
			(*library)[index1], (*current)[index2] = (*current)[index2], (*library)[index1]
		}
	}
}

//
//func main() {
//	// 示例使用
//	questionList := []QuestionBank{
//		{TopicType: "填空题", Difficulty: 0.5},
//		{TopicType: "选择题", Difficulty: 0.6},
//		{TopicType: "判断题", Difficulty: 0.4},
//		{TopicType: "程序设计题", Difficulty: 0.7},
//	}
//	iterationsNum := 10
//	targetDifficult := 0.5
//	TKTCount := 1
//	XZTCount := 1
//	PDTCount := 1
//	JDTCount := 1
//
//	gi := NewGeneticIteration(iterationsNum, questionList, targetDifficult, TKTCount, XZTCount, PDTCount, JDTCount)
//	gi.Run()
//
//	fmt.Println("Variance:", gi.variance)
//}
