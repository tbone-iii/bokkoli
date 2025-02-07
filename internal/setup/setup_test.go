package setup

import (
	"testing"
)

func TestValidatePortNumberInRange(t *testing.T) {
	expected := true
	result := validatePort("8080")

	if expected != result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}
func TestValidatePortNumberOutOfRange(t *testing.T) {
	expected := false
	result := validatePort("20")

	if expected != result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}

func TestValidatePortNumberNotANumber(t *testing.T) {
	expected := false
	result := validatePort("veggie-town")

	if expected != result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}

func TestValidatePortNumberEmpty(t *testing.T) {
	expected := false
	result := validatePort("")

	if expected != result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}

func TestValidateUsernameExists(t *testing.T) {
	expected := true
	result := validateUsername("Pickle132")

	if expected != result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}

func TestValidateUsernameEmpty(t *testing.T) {
	expected := false
	result := validateUsername("")

	if expected != result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}
