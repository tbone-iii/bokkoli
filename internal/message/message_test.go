package message

import (
	"testing"
)

func TestValidatePortNumberInRange(t *testing.T) {
	max_line_length := 50
	test_string := "This is a sample string with a superlongwordthatexceedsthemaximumlinelengthlimitandneedstobebroken into parts."
	expected :=
		"This is a sample string with a \n" +
			"superlongwordthatexceedsthemaximumlinelengthlimi\n" +
			"tandneedstobebroken into parts."

	result := wrapText(test_string, max_line_length)

	if expected != result {
		t.Errorf("Expected \n'%s'\n, got \n'%s'", expected, result)
	}
}
