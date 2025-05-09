package services

import (
	"fmt"
	"io"
	"net/http"
)

func FetchExternalAPI(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Close Request at the end of the function

	// Check if the HTTP Request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API return status %d", resp.StatusCode)
	}

	// Read the Data within the Response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
