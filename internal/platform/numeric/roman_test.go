package numeric

import (
	"fmt"
	"testing"
)

func TestToRomanSystem(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{1, "I"},
		{2, "II"},
		{3, "III"},
		{4, "IV"},
		{5, "V"},
		{6, "VI"},
		{7, "VII"},
		{8, "VIII"},
		{9, "IX"},
		{10, "X"},
		{11, "XI"},
		{14, "XIV"},
		{15, "XV"},
		{19, "XIX"},
		{20, "XX"},
		{39, "XXXIX"},
		{40, "XL"},
		{49, "XLIX"},
		{50, "L"},
		{89, "LXXXIX"},
		{90, "XC"},
		{99, "XCIX"},
		{100, "C"},
		{399, "CCCXCIX"},
		{400, "CD"},
		{499, "CDXCIX"},
		{500, "D"},
		{899, "DCCCXCIX"},
		{900, "CM"},
		{999, "CMXCIX"},
		{1000, "M"},
		{3999, "MMMCMXCIX"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input %d", tc.input), func(t *testing.T) {
			output := ToRomanSystem(tc.input)
			if output != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, output)
			}
		})
	}
}
