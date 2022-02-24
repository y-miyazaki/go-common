package main

import "testing"

type TestStruct struct {
	Title       string
	Name        string
	Age         string
	TestStruct2 *TestStruct2
}
type TestStruct2 struct {
	Title2      string
	Name2       string
	Age2        string
	TestStruct3 *TestStruct3
}
type TestStruct3 struct {
	Title3 string
	Name3  string
	Age3   string
}

func CreateTestData() []TestStruct {
	data := []TestStruct{}
	for i := 0; i < 100000; i++ {
		data = append(data, TestStruct{
			Title: "testtitle",
			Name:  "testname",
			Age:   "1",
			TestStruct2: &TestStruct2{
				Title2: "testtitle2",
				Name2:  "testname2",
				Age2:   "2",
				TestStruct3: &TestStruct3{
					Title3: "testtitle3",
					Name3:  "testname3",
					Age3:   "3",
				},
			},
		})
	}
	return data
}

func BenchmarkStruct1(b *testing.B) {
	data := CreateTestData()
	var base []string = []string{}

	for _, v := range data {
		base = append(base, v.TestStruct2.TestStruct3.Title3)
		base = append(base, v.TestStruct2.TestStruct3.Name3)
		base = append(base, v.TestStruct2.TestStruct3.Age3)
		base = append(base, v.TestStruct2.TestStruct3.Title3)
		base = append(base, v.TestStruct2.TestStruct3.Name3)
		base = append(base, v.TestStruct2.TestStruct3.Age3)
	}
}
func BenchmarkStruct2(b *testing.B) {
	data := CreateTestData()
	var base []string = []string{}

	for _, v := range data {
		t3 := v.TestStruct2.TestStruct3
		base = append(base, t3.Title3)
		base = append(base, t3.Name3)
		base = append(base, t3.Age3)
		base = append(base, t3.Title3)
		base = append(base, t3.Name3)
		base = append(base, t3.Age3)
	}
}

func BenchmarkMake(b *testing.B) {
	data := []string{}

	for i := 0; i < b.N; i++ {
		line := make([]string, 0, 20)
		line = []string{
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
		}
		data = append(data, line...)
	}
}

func BenchmarkMake2(b *testing.B) {
	data := []string{}

	for i := 0; i < b.N; i++ {
		line := make([]string, 0, 10)
		line = []string{
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
		}
		data = append(data, line...)
	}
}
func BenchmarkMakeNo(b *testing.B) {
	data := []string{}

	for i := 0; i < b.N; i++ {
		line := []string{
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
			"1",
		}
		data = append(data, line...)
	}
}
