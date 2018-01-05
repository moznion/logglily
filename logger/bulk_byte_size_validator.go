package logger

import "fmt"

func validateBulkByteSizeThreshold(given int) error {
	const maximumBulkByteSizeThreshold = 5 * 1024 * 1024 // nearly 5MB (doc: https://www.loggly.com/docs/http-bulk-endpoint/)

	if given <= 0 {
		return fmt.Errorf(
			"bulk byte size threshold must be natural number [given: %d]",
			given,
		)
	}

	if given > maximumBulkByteSizeThreshold {
		return fmt.Errorf(
			"bulk byte size threshold is exceeded [maximum: %d, given: %d]",
			maximumBulkByteSizeThreshold,
			given,
		)
	}

	return nil
}
