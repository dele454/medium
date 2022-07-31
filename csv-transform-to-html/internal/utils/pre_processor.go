package utils

import (
	"fmt"
	"html"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dele454/medium/csv-transform-to-html/internal/errs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// SalesRecord details data about a sale record
type SalesRecord struct {
	Region        string `csv:"Region"`
	Country       string `csv:"Country" processor:"required"`
	ItemType      string `csv:"ItemType" processor:"required"`
	SalesChannel  string `csv:"SalesChannel"`
	OrderPriority string `csv:"OrderPriority"`
	OrderDate     string `csv:"OrderDate" processor:"date,required"`
	OrderID       string `csv:"OrderID" processor:"numeric,required"`
	ShipDate      string `csv:"ShipDate" processor:"date,required"`
	UnitsSold     string `csv:"UnitsSold" processor:"numeric,required"`
	UnitPrice     string `csv:"UnitPrice" processor:"amount,required"`
	UnitCost      string `csv:"UnitCost" processor:"amount,required"`
	TotalRevenue  string `csv:"TotalRevenue" processor:"amount"`
	TotalCost     string `csv:"TotalCost" processor:"amount,required"`
	TotalProfit   string `csv:"TotalProfit" processor:"amount,required"`
}

// Preprocessor operations a transformer must perform
type PreProcessor interface {
	EscapeHTML(val string) string
	NotEmpty(val, field string) error
	SanitizeString(str string) string
	ParseDate(val, field string) error
	ParseFloat(val, field string) error
	ParseInteger(val, field string) error

	Unmarshal(record []string, sr SalesRecord) (SalesRecord, error)
}

// Processor
type Processor struct{}

// NewProcessor creates a new processor
func NewProcessor() PreProcessor {
	return &Processor{}
}

// EscapeHTML escapes html tags from a value
func (p *Processor) EscapeHTML(val string) string {
	return html.EscapeString(val)
}

// SanitizeString strips out tabs, trailing and leading spaces and quotes from a string
func (p *Processor) SanitizeString(str string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(str, "\"", ""), "\t", ""))
}

// NotEmpty checks if a value is empty
func (p *Processor) NotEmpty(val, field string) error {
	if val == "" {
		return fmt.Errorf(errs.ErrorFieldIsEmpty.Error(), field)
	}

	return nil
}

// ParseDate parses a date string in a DD/MM/YYYY format
func (p *Processor) ParseDate(val, field string) error {
	var err error

	if err := p.NotEmpty(val, field); err != nil {
		return err
	}

	// parse date field
	_, err = time.Parse("1/2/2006", val)
	if err != nil {
		return fmt.Errorf(errs.ErrorFieldNotValid.Error(), field)
	}

	return nil
}

// ParseFloat parse a float value
func (p *Processor) ParseFloat(val, field string) error {
	if err := p.NotEmpty(val, field); err != nil {
		return err
	}

	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fmt.Errorf(errs.ErrorFieldNotValid.Error(), field)
	}

	return nil
}

// ParseInteger parse a int value
func (p *Processor) ParseInteger(val, field string) error {
	if err := p.NotEmpty(val, field); err != nil {
		return err
	}

	_, err := strconv.ParseInt(val, 0, 36)
	if err != nil {
		return fmt.Errorf(errs.ErrorFieldNotValid.Error(), field)
	}

	return nil
}

// Unmarshal unmarshals records found into the SalesRecord struct
//
// Cycles through all the fields for the SalesRecord struct in order
// to decipher which field(s) needs a pre-processor and apply
// as record is unmarshalled.
func (p *Processor) Unmarshal(record []string, sr SalesRecord) (SalesRecord, error) {
	s := reflect.ValueOf(sr).Type()
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		tags := strings.Split(field.Tag.Get("processor"), ",")

		for _, t := range tags {
			switch t {
			case "date":
				err := p.ParseDate(record[i], field.Name)
				if err != nil {
					return sr, err
				}
				reflect.ValueOf(&sr).Elem().Field(i).SetString(record[i])
			case "amount":
				err := p.ParseFloat(record[i], field.Name)
				if err != nil {
					return sr, err
				}
				reflect.ValueOf(&sr).Elem().Field(i).SetString(record[i])
			case "numeric":
				err := p.ParseInteger(record[i], field.Name)
				if err != nil {
					return sr, err
				}
				reflect.ValueOf(&sr).Elem().Field(i).SetString(record[i])
			default:
				err := p.NotEmpty(record[i], field.Name)
				if err != nil {
					return sr, nil
				}
				reflect.ValueOf(&sr).Elem().Field(i).SetString(p.EscapeHTML(p.SanitizeString(record[i])))
			}
		}
	}

	return sr, nil
}

// UnsupportedType
type UnsupportedType struct {
	Type string
}

// Error error message for UnsupportedType type
func (e *UnsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}

// GetHeaders gets all expected headers for the SalesOrder struct
//
// Introspects all csv tags for the SalesOrder struct fields
// and builds a list of headers from that thus making it dynamic
func GetHeaders() []string {
	var (
		headers []string
		sr      SalesRecord
	)

	s := reflect.ValueOf(sr).Type()
	for i := 0; i < s.NumField(); i++ {
		headers = append(headers, cases.Title(language.Und, cases.NoLower).
			String(s.Field(i).Tag.Get("csv")))
	}

	return headers
}

// RootDir returns the root directory for the application
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Join(filepath.Dir(d), "..")
}
