package excel

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/util"
	"github.com/gophab/gophrame/core/util/array"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/xuri/excelize/v2"
)

type ExcelColumn struct {
	Title  string /* 标题 */
	Path   string /* 路径 */
	Column string /* 列 */
}

type Exporter struct {
	*excelize.File
	Columns    []ExcelColumn
	Sheet      string
	Index      int
	CurrentRow int
	hasHeader  bool
	Sheets     []string
}

func NewExporter(columns []ExcelColumn) *Exporter {
	return &Exporter{
		File:       excelize.NewFile(),
		Columns:    columns,
		Sheet:      "Sheet1",
		Index:      1,
		CurrentRow: 1,
		hasHeader:  false,
		Sheets:     []string{"Sheet1"},
	}
}

func (e *Exporter) WithColumns(columns []ExcelColumn) *Exporter {
	e.Columns = columns
	return e
}

func (e *Exporter) WithIndex(index int) *Exporter {
	e.Index = index
	return e
}

func (e *Exporter) WithCurrentRow(currentRow int) *Exporter {
	e.CurrentRow = currentRow
	return e
}

func (e *Exporter) WithSheet(sheet string) *Exporter {
	if sheet != e.Sheet {
		if b, _ := array.Contains[string](e.Sheets, sheet); !b {
			if _, err := e.NewSheet(sheet); err == nil {
				e.Sheets = append(e.Sheets, sheet)
			}
		}
		e.Sheet = sheet
	}
	return e
}

func (e *Exporter) PrepareHeader() *Exporter {
	// Head
	if !e.hasHeader {
		for _, column := range e.Columns {
			e.SetCellStr("sheet1", column.Column+"1", column.Title)
		}
		e.CurrentRow++

		e.hasHeader = true
	}
	return e
}

func (e *Exporter) AppendRow(record interface{}) {
	// 导出的Column
	row := util.Struct2Map(record)

	for _, column := range e.Columns {
		if column.Path == "_index_" {
			e.SetCellValue(e.Sheet, fmt.Sprintf("%s%d", column.Column, e.CurrentRow), e.Index)
		} else if v, b := util.GetRecordField(row, column.Path); b {
			e.SetCellValue(e.Sheet, fmt.Sprintf("%s%d", column.Column, e.CurrentRow), v)
		}
	}

	e.CurrentRow++
	e.Index++
}

// Write provides a function to write to an io.Writer.
func (e *Exporter) Write(ctx *gin.Context, prefix string, opts ...excelize.Options) {
	// 设置 HTTP 响应的头信息
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	filename := fmt.Sprintf("%s_%v.xlsx", prefix, time.Now().Format("20060102_150405"))
	ctx.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(filename))
	// 将 Excel 文件写入 HTTP 响应
	// buffer, _ := e.WriteToBuffer()
	if _, err := e.WriteTo(ctx.Writer, opts...); err != nil {
		response.FailMessage(ctx, 500, err.Error())
		return
	}
	response.OK(ctx, nil)
}

func (e *Exporter) Save() {
	if err := e.SaveAs("/tmp/temp.xlsx"); err != nil {
		logger.Error("Error: ", err.Error())
	}
}
