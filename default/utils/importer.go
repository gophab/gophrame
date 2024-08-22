package utils

import (
	"fmt"
	"io"

	"github.com/gophab/gophrame/core/util/array"
	"github.com/xuri/excelize/v2"
)

type Importer struct {
	*excelize.File
	Columns    []string
	Sheet      string
	Index      int
	CurrentRow int
	hasHeader  bool
	Sheets     []string
}

func NewImporter() *Importer {
	return &Importer{
		Columns:    nil,
		Sheet:      "Sheet1",
		Index:      1,
		CurrentRow: 1,
		hasHeader:  false,
		Sheets:     []string{"Sheet1"},
	}
}

func (e *Importer) WithColumns(columns []string) *Importer {
	e.Columns = columns
	return e
}

func (e *Importer) WithIndex(index int) *Importer {
	e.Index = index
	return e
}

func (e *Importer) WithCurrentRow(currentRow int) *Importer {
	e.CurrentRow = currentRow
	return e
}

func (e *Importer) WithSheet(sheet string) *Importer {
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

func (e *Importer) PrepareHeader() *Importer {
	// Head
	// 1. read header
	e.Columns = make([]string, 0)
	for i := 1; ; i++ {
		if cell, err := e.GetCellValue(e.Sheet, fmt.Sprintf("%s%d", getColumnName(i-1), e.CurrentRow)); err != nil || cell == "" {
			break
		} else {
			e.Columns = append(e.Columns, cell)
		}
	}
	e.CurrentRow++
	e.hasHeader = true
	return e
}

func getColumnName(index int) string {
	var prefix = index / 26
	if prefix == 0 {
		var c = 'A'
		var name = rune(int(c) + index)
		return string(name)
	} else {
		index -= prefix * 26
		return getColumnName(prefix-1) + getColumnName(index)
	}
}

func (e *Importer) ReadFile(fileName string, callback func([]map[string]string), batchInSize int) (err error) {
	e.File, err = excelize.OpenFile(fileName)
	if err != nil {
		return err
	}
	defer func() {
		err = e.Close()
	}()
	e.readInBatch(callback, batchInSize)
	return
}

func (e *Importer) ReadReader(reader io.Reader, callback func([]map[string]string), batchInSize int) (err error) {
	e.File, err = excelize.OpenReader(reader)
	if err != nil {
		return err
	}
	defer func() {
		err = e.Close()
	}()

	e.readInBatch(callback, batchInSize)
	return
}

func (e *Importer) readInBatch(callback func([]map[string]string), batchInSize int) {
	// 1. 准备Header
	if e.Columns == nil {
		e.PrepareHeader()
	}

	// 2. read rows
	var rows = make([]map[string]string, 0)
	for r := e.CurrentRow; ; r++ {
		var blank = true
		var row = make(map[string]string)
		for c := 1; c <= len(e.Columns); c++ {
			if cell, err := e.GetCellValue("Sheet1", fmt.Sprintf("%s%d", getColumnName(c-1), r)); err == nil {
				row[e.Columns[c-1]] = cell
				if cell != "" {
					blank = false
				}
			}
		}
		if blank {
			break
		}

		rows = append(rows, row)
		if len(rows) == batchInSize {
			callback(rows)
			rows = make([]map[string]string, 0)
		}
	}
	callback(rows)
}
