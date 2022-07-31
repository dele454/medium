package utils

import (
	"testing"
)

func TestParseDate(t *testing.T) {
	var p Processor

	err := p.ParseDate("6/27/2010", "ShipDate")
	if err != nil {
		t.Fatalf("Date should be valid")
	}

	err = p.ParseDate("27/01/2010", "ShipDate")
	if err == nil {
		t.Fatal("Date should not be valid")
	}

	err = p.ParseDate("", "ShipDate")
	if err == nil {
		t.Fatal("Date should not be valid")
	}
}

func TestParseFloat(t *testing.T) {
	var p Processor

	err := p.ParseFloat("ParseBirthday", "Credit Limit")
	if err == nil {
		t.Fatal("Credit limit should not be valid")
	}

	err = p.ParseFloat("121.99", "Credit Limit")
	if err != nil {
		t.Fatal("Credit limit should be valid")
	}

	err = p.ParseFloat("1000000", "Credit Limit")
	if err != nil {
		t.Fatal("Credit limit should be valid")
	}
}

func TestUnmarshal(t *testing.T) {
	var err error
	var p Processor

	type sample struct {
		record  []string
		errored bool
	}

	expectations := []sample{
		{
			// all data required and fields are valid
			[]string{"Australia and Oceania", "Tuvalu", "Baby Food", "Offline", "H", "5/28/2010", "669165933", "6/27/2010", "9925", "255.28", "159.42", "2533654.00", "1582243.50", "951410.50"},
			false,
		},
		{
			// country is required but empty
			[]string{"Australia and Oceania", "", "Baby Food", "Offline", "H", "5/28/2010", "669165933", "6/27/2010", "9925", "255.28", "159.42", "2533654.00", "1582243.50", "951410.50"},
			true,
		},
		{
			// TotalProfit is expected to be numeric but string is passed
			[]string{"Central America and the Caribbean", "Grenada", "Cereal", "Online", "C", "8/22/2012", "963881480", "9/15/2012", "2804", "205.70", "117.11", "576782.80", "328376.44", "SAY WHAT"},
			true,
		},
	}

	for _, tc := range expectations {
		t.Run("", func(t *testing.T) {
			_, err = p.Unmarshal(tc.record, SalesRecord{})
			if tc.errored != (err != nil) {
				t.Fatalf("\nUnmarshal Mismatch:\nExpected: %v\nGot: %v", tc.errored, err != nil)
			}
		})
	}
}

func TestGetHeaders(t *testing.T) {
	if len(GetHeaders()) == 0 {
		t.Fatal("Headers should be detected")
	}
}
