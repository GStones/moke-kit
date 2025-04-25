package nosql

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 1. **基本类型转换速度最快**：仅处理基本类型时，函数性能最佳，每秒可处理约440万次操作，平均每次操作仅需257.6纳秒。
// 2. **JSON处理是性能瓶颈**：含有JSON序列化字符串的场景性能明显下降，处理速度仅为基本类型的21%左右。
// 3. **字段数量影响较小**：MultipleFields测试（10个基本类型字段）比基本测试（4个字段）慢约2.8倍，说明字段数量增加时性能下降并不是线性的。
// 4. **JSON数据大小显著影响性能**：处理大型JSON字符串时（LargeJSONString测试），性能下降最为明显，每秒仅能处理约3.1万次操作，是基本测试的0.7%。这表明当处理大型嵌套结构时，JSON解析成为严重瓶颈。
// 5. **混合类型结构体表现适中**：混合基本类型和JSON字符串的场景性能处于中间水平，每次操作约820.1纳秒，比纯基本类型慢约3.2倍。
// 用于基准测试的结构体
type benchStruct struct {
	ID       string
	Name     string
	Age      int
	IsActive bool
	Tags     []string
	SubData  *benchSubData
}

type benchSubData struct {
	SubField1 string
	SubField2 int
}

func BenchmarkMap2StructShallow(b *testing.B) {
	// 准备基础类型的字段
	testMap := map[string]any{
		"ID":       "12345",
		"Name":     "测试用户",
		"Age":      30,
		"IsActive": true,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var dst benchStruct
		err := map2StructShallow(testMap, &dst)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMap2StructShallow_WithJSONString(b *testing.B) {
	// 准备包含JSON字符串的测试数据
	tags, _ := json.Marshal([]string{"标签1", "标签2", "标签3"})
	subData, _ := json.Marshal(&benchSubData{
		SubField1: "子字段值",
		SubField2: 42,
	})

	testMap := map[string]any{
		"ID":       "12345",
		"Name":     "测试用户",
		"Age":      30,
		"IsActive": true,
		"Tags":     string(tags),    // JSON 序列化的字符串
		"SubData":  string(subData), // JSON 序列化的字符串
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var dst benchStruct
		err := map2StructShallow(testMap, &dst)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMap2StructShallow_MixedTypes(b *testing.B) {
	// 准备混合基本类型和JSON字符串的测试数据
	subData, _ := json.Marshal(&benchSubData{
		SubField1: "子字段值",
		SubField2: 42,
	})

	testMap := map[string]any{
		"ID":       "12345",
		"Name":     "测试用户",
		"Age":      30,
		"IsActive": true,
		"SubData":  string(subData), // JSON 序列化的字符串
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var dst benchStruct
		err := map2StructShallow(testMap, &dst)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMap2StructShallow_LargeJSONString(b *testing.B) {
	// 创建一个较大的JSON数据结构
	largeSubData := struct {
		Fields     map[string]string
		Values     []int
		Properties map[string][]float64
	}{
		Fields:     make(map[string]string),
		Values:     make([]int, 100),
		Properties: make(map[string][]float64),
	}

	// 填充大型数据结构
	for i := 0; i < 50; i++ {
		largeSubData.Fields[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("value_%d", i)
		largeSubData.Values[i] = i * 10
		largeSubData.Properties[fmt.Sprintf("prop_%d", i)] = []float64{float64(i), float64(i) + 0.5, float64(i) + 1.0}
	}

	largeDataJSON, _ := json.Marshal(largeSubData)

	testMap := map[string]any{
		"ID":       "12345",
		"Name":     "测试用户",
		"Age":      30,
		"IsActive": true,
		"SubData":  string(largeDataJSON), // 大型JSON字符串
	}

	type largeBenchStruct struct {
		ID       string
		Name     string
		Age      int
		IsActive bool
		SubData  struct {
			Fields     map[string]string
			Values     []int
			Properties map[string][]float64
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var dst largeBenchStruct
		err := map2StructShallow(testMap, &dst)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 测试多个不同大小的字段
func BenchmarkMap2StructShallow_MultipleFields(b *testing.B) {
	// 准备10个字段的测试数据
	testMap := map[string]any{
		"Field1":  "值1",
		"Field2":  "值2",
		"Field3":  "值3",
		"Field4":  100,
		"Field5":  200.5,
		"Field6":  true,
		"Field7":  int64(9999999999),
		"Field8":  uint(123),
		"Field9":  float32(3.14),
		"Field10": byte(65),
	}

	type tenFieldStruct struct {
		Field1  string
		Field2  string
		Field3  string
		Field4  int
		Field5  float64
		Field6  bool
		Field7  int64
		Field8  uint
		Field9  float32
		Field10 byte
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var dst tenFieldStruct
		err := map2StructShallow(testMap, &dst)
		if err != nil {
			b.Fatal(err)
		}
	}
}
