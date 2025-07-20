package sumup

import (
	"strings"
	"testing"

	"github.com/sumup/sumup-go/readers"
)

func TestExtractSumup422ErrorDetail(t *testing.T) {
	err := &readers.CreateReaderCheckout422Response{
		Errors: &readers.CreateReaderCheckout422ResponseErrors{
			"detail": "The device is offline",
			"type":   "READER_OFFLINE",
		},
	}
	expected := "The device is offline (READER_OFFLINE)"
	actual := extractSumup422ErrorDetail(*err)

	if actual.Error() != expected {
		t.Errorf("Expected %q but got %q", expected, actual.Error())
	}

	err = &readers.CreateReaderCheckout422Response{
		Errors: &readers.CreateReaderCheckout422ResponseErrors{
			"detail": "The device is offline",
			"foo":    "bar",
		},
	}
	expected = "The device is offline"
	actual = extractSumup422ErrorDetail(*err)

	if actual.Error() != expected {
		t.Errorf("Expected %q but got %q", expected, actual.Error())
	}

	err = &readers.CreateReaderCheckout422Response{
		Errors: &readers.CreateReaderCheckout422ResponseErrors{
			"foo": "bar",
		},
	}
	expected = "errors=&map[foo:bar]"
	actual = extractSumup422ErrorDetail(*err)

	if actual.Error() != expected {
		t.Errorf("Expected %q but got %q", expected, actual.Error())
	}
}

func TestExtractCreateCheckoutErrorDetails(t *testing.T) {
	err := &readers.CreateReaderCheckout422Response{
		Errors: &readers.CreateReaderCheckout422ResponseErrors{
			"detail": "The device is offline",
			"type":   "READER_OFFLINE",
		},
	}
	expected := "The device is offline (READER_OFFLINE)"
	actual := extractCreateCheckoutErrorDetails(err)

	if actual.Error() != expected {
		t.Errorf("Expected %q but got %q", expected, actual.Error())
	}

	internalServerError := "Internal server error"
	err2 := &readers.CreateReaderCheckout500Response{
		Errors: &readers.CreateReaderCheckout500ResponseErrors{
			Detail: &internalServerError,
		},
	}
	actual = extractCreateCheckoutErrorDetails(err2)
	expected = "errors="

	if !strings.HasPrefix(actual.Error(), expected) {
		t.Errorf("Expected %q but got %q", expected, actual.Error())
	}
}
