package transform

import (
	"github.com/dele454/medium/csv-transform-to-html/internal/report"
)

// Output details of the transformation sent to an io.Writer
type Output struct {
	FileName     string
	FileLocation string
	TotalHeaders int
	Headers      []string
	Data         [][]string
}

// Transformer ops every transformer should conform to
type Transformer interface {
	WriteOutputToFile(file Output)
}

// HTMLTransformer handles the processing of customer data
// with the aid of preprocessors and writing the output
// as an HTML document.
type HTMLTransformer struct {
	reporter report.Reporter
}

// XMLTransformer handles the processing of customer data
// with the aid of preprocessors and writing the output
// as an XML document.
type XMLTransformer struct{}
