package sumup

import (
	"fmt"

	"github.com/sumup/sumup-go/readers"
)

func extractCreateCheckoutErrorDetails(err error) error {
	if e, ok := err.(*readers.CreateReaderCheckout422Response); ok {
		return extractSumup422ErrorDetail(*e)
	}

	return err
}

func extractSumup422ErrorDetail(err readers.CreateReaderCheckout422Response) error {
	if err.Errors == nil {
		return &err
	}

	detail := ""
	typ := ""

	if val, ok := (*err.Errors)["detail"]; ok {
		detail, _ = val.(string)
	}

	if val, ok := (*err.Errors)["type"]; ok {
		typ, _ = val.(string)
	}

	if detail != "" && typ != "" {
		return fmt.Errorf("%s (%s)", detail, typ)
	} else if detail != "" {
		return fmt.Errorf("%s", detail)
	}

	return &err
}
