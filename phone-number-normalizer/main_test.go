package main

import "testing"

type testPhoneNumberCase struct {
	input string
	output string
}

func TestNormalizer(t *testing.T) {
	testCases := []testPhoneNumberCase {
		{"0128787989", "0128787989"},
		{"012-8787989", "0128787989"},
		{"012-8787-989", "0128787989"},
		{"+60-12-8787989", "0128787989"},
		{"(012)8787989", "0128787989"},
		{"016-90908989", "01690908989"},
		{"+6016-9090-8989", "01690908989"},
		{"01790908989", "01790908989"},
		{"6019802080208", "false"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			givenNumber := normalizer(tc.input)
			if givenNumber != tc.output {
				t.Errorf("intput: %s; output: %s", givenNumber, tc.output)
			}
		})
	}
}