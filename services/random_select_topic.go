package services

import (
	"java2go/entity"
	"math"
	"math/rand"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子
}

// RandomSelectTopic 随机选题核心算法
func RandomSelectTopic(dataSource []entity.QuestionBank, targetDiff float64, selectCount int) []entity.QuestionBank {
	result := make([]entity.QuestionBank, 0, selectCount)

	// 参数校验
	if len(dataSource) == 0 ||
		targetDiff < 1.0 ||
		targetDiff > 5.0 ||
		selectCount <= 0 {
		return result
	}

	// 预排序（保持与Java版相同逻辑，虽然实际不影响算法正确性）
	sort.SliceStable(dataSource, func(i, j int) bool {
		return dataSource[i].Difficulty < dataSource[j].Difficulty
	})

	workingSet := make([]entity.QuestionBank, len(dataSource))
	copy(workingSet, dataSource) // 避免修改原始数据

	for remaining := selectCount; remaining > 0 && len(workingSet) > 0; remaining-- {
		// 单次选择流程
		if selected, ok := selectQuestion(workingSet, targetDiff); ok {
			result = append(result, selected)
			workingSet = removeQuestion(workingSet, selected.ID)
		} else {
			break // 无可用题目时提前终止
		}
	}

	return result
}

// selectQuestion 单次选题逻辑
func selectQuestion(candidates []entity.QuestionBank, targetDiff float64) (entity.QuestionBank, bool) {
	if len(candidates) == 0 {
		return entity.QuestionBank{}, false
	}

	// 第一阶段：寻找最小差值
	minDelta := math.MaxFloat64
	var minCandidates []entity.QuestionBank

	for _, q := range candidates {
		currentDelta := math.Abs(float64(q.Difficulty) - targetDiff)

		switch {
		case currentDelta < minDelta:
			minDelta = currentDelta
			minCandidates = []entity.QuestionBank{q}
		case currentDelta == minDelta:
			minCandidates = append(minCandidates, q)
		}
	}

	// 第二阶段：从候选中随机选择
	if len(minCandidates) == 0 {
		return entity.QuestionBank{}, false
	}
	return minCandidates[rand.Intn(len(minCandidates))], true
}

// removeQuestion 从切片中删除指定题目（保持顺序）
func removeQuestion(slice []entity.QuestionBank, id int) []entity.QuestionBank {
	for i, q := range slice {
		if q.ID == id {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

/* 使用示例
func main() {
    // 准备测试数据
    questions := []QuestionBank{
        {ID: 1, Difficulty: 3.2},
        {ID: 2, Difficulty: 4.5},
        {ID: 3, Difficulty: 2.8},
        {ID: 4, Difficulty: 4.5},
        {ID: 5, Difficulty: 3.9},
    }

    // 执行选题
    selected := RandomSelectTopic(questions, 4.0, 3)

    // 输出结果
    fmt.Println("Selected questions:")
    for _, q := range selected {
        fmt.Printf("ID: %d, Difficulty: %.1f\n", q.ID, q.Difficulty)
    }
}
*/
