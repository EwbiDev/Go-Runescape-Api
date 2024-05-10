package runescape

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestListGrandExchangeItems(t *testing.T) {
	// Mock HTTP server to simulate API responses
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request URL
		if !strings.Contains(r.URL.Path, "/m=itemdb_rs/api/catalogue/items.json") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Check query parameters
		q := r.URL.Query()
		if q.Get("category") != "1" || q.Get("alpha") != "a" || q.Get("page") != "1" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Prepare mock response
		resp := GeResponse{
			Total: 2,
			Items: []Item{
				{
					Id:          1,
					Name:        "Item 1",
					Description: "Description of item 1",
				},
				{
					Id:          2,
					Name:        "Item 2",
					Description: "Description of item 2",
				},
			},
		}
		respJSON, _ := json.Marshal(resp)

		// Write mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJSON)
	}))
	defer ts.Close()

	// Create a new Client with the test server's URL
	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(ts.URL)

	// Call the method being tested
	response, err := client.ListGrandExchangeItems("rs3", "a", 1, 1)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Check the response
	if response == nil {
		t.Error("Expected non-nil response")
		return
	}
	if response.Total != 2 {
		t.Errorf("Expected Total to be 2, got %d", response.Total)
	}
	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}
	if response.Items[0].Id != 1 || response.Items[1].Id != 2 {
		t.Error("Unexpected item IDs")
	}
}

func TestListGrandExchangeItems_InvalidGameType(t *testing.T) {
	client := NewClient(nil)
	_, err := client.ListGrandExchangeItems("invalidGameType", "a", 1, 1)
	if err == nil || !strings.Contains(err.Error(), "gameType must be") {
		t.Errorf("Expected error about invalid gameType, got: %v", err)
	}
}

func TestListGrandExchangeItems_InvalidItemAlpha(t *testing.T) {
	client := NewClient(nil)
	_, err := client.ListGrandExchangeItems("rs3", "invalidAlpha", 1, 1)
	if err == nil || !strings.Contains(err.Error(), "itemAlpha should be") {
		t.Errorf("Expected error about invalid itemAlpha, got: %v", err)
	}
}

func TestListGrandExchangeItems_NegativePage(t *testing.T) {
	client := NewClient(nil)
	_, err := client.ListGrandExchangeItems("rs3", "a", 1, -1)
	if err == nil || !strings.Contains(err.Error(), "page must be > 1") {
		t.Errorf("Expected error about negative page, got: %v", err)
	}
}

func TestListGrandExchangeItems_FailedRequest(t *testing.T) {
	// Mock HTTP server that always returns an error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	// Create a new Client with the test server's URL
	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(ts.URL)

	// Call the method being tested
	_, err := client.ListGrandExchangeItems("rs3", "a", 1, 1)

	// Check for errors
	if err == nil || !strings.Contains(err.Error(), "500 Internal Server Error") {
		t.Errorf("Expected error about failed request, got: %v", err)
	}
}

func TestListGrandExchangeItems_InvalidResponseJSON(t *testing.T) {
	// Mock HTTP server that returns invalid JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Write an invalid JSON response
		w.Write([]byte(`{"total": 2, "items": [{},"}`))
	}))
	defer ts.Close()

	// Create a new Client with the test server's URL
	client := NewClient(nil)
	client.BaseURL, _ = url.Parse(ts.URL)

	// Call the method being tested
	_, err := client.ListGrandExchangeItems("rs3", "a", 1, 1)

	// Check for errors
	if err == nil || !strings.Contains(err.Error(), "unexpected end of JSON input") {
		t.Errorf("Expected error about invalid JSON response, got: %v", err)
	}
}

func TestTrendPrice_UnmarshalJSON(t *testing.T) {
	// Valid JSON for TrendPrice
	validJSON1 := `{"trend":"neutral","price":"100"}`
	var tp TrendPrice

	err := json.Unmarshal([]byte(validJSON1), &tp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check unmarshaled fields
	if tp.Trend != "neutral" {
		t.Errorf("Expected Trend to be 'neutral', got '%s'", tp.Trend)
	}
	if tp.Price != "100" {
		t.Errorf("Expected Price to be '100', got '%s'", tp.Price)
	}

	validJSON2 := `{"trend":"neutral","price":100}`
	err = json.Unmarshal([]byte(validJSON2), &tp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check unmarshaled fields
	if tp.Trend != "neutral" {
		t.Errorf("Expected Trend to be 'neutral', got '%s'", tp.Trend)
	}
	if tp.Price != "100" {
		t.Errorf("Expected Price to be '100', got '%s'", tp.Price)
	}

	// Invalid JSON for TrendPrice
	invalidJSON := `{"trend":"up","price":[100]}`
	err = json.Unmarshal([]byte(invalidJSON), &tp)
	if err == nil || !strings.Contains(err.Error(), "unexpected type for price") {
		t.Errorf("Expected error about unexpected type for price, got: %v", err)
	}
}
