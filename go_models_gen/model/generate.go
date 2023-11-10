package model

import "github.com/spf13/cast"

// interval.间隔
var _interval = " "

// IGenerate Generate Printing Interface.生成打印接口
type IGenerate interface {
	// Get the generate data .获取结果数据
	Generate() string
}

// PrintAtom . atom print .原始打印
type PrintAtom struct {
	lines []string
}

// Add add one to print.打印
func (p *PrintAtom) Add(str ...interface{}) {
	var tmp string
	for _, v := range str {
		tmp += cast.ToString(v) + _interval
	}
	p.lines = append(p.lines, tmp)
}

// Generates Get the generated list.获取生成列表
func (p *PrintAtom) Generates() []string {
	return p.lines
}
