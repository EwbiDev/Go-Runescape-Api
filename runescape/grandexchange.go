// Package runescape provides a Go wrapper for the Runescape API.
package runescape

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"strconv"
)

// GeResponse represents the JSON response from the Grand Exchange endpoint.
type GeResponse struct {
	Total int    `json:"total"`
	Items []Item `json:"items"`
}

// Item represents an item in the Grand Exchange.
type Item struct {
	Icon        string     `json:"icon"`
	Icon_large  string     `json:"icon_large"`
	Id          int        `json:"id"`
	TypeItem    string     `json:"type"`
	TypeIcon    string     `json:"typeIcon"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Current     TrendPrice `json:"current"`
	Today       TrendPrice `json:"today"`
	Members     string     `json:"members"`
}

// TrendPrice represents the Current and Today properties in Item
type TrendPrice struct {
	Trend string `json:"trend"`
	Price string `json:"price"`
}

// UnmarshalJSON unmarshals JSON data into a TrendPrice struct.
//
// The Runescape api returns price as both a string and an int so therefore needs to be coerced into a string for consistency.
func (tp *TrendPrice) UnmarshalJSON(data []byte) error {
	var raw struct {
		Trend string      `json:"trend"`
		Price interface{} `json:"price"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	tp.Trend = raw.Trend

	switch v := raw.Price.(type) {
	case string:
		tp.Price = v
	case int, float64:
		tp.Price = strconv.Itoa(int(v.(float64)))
	default:
		return errors.New("unexpected type for price")
	}

	return nil
}

// ListGrandExchangeItems retrieves a list of items from the Grand Exchange.
//
// gameType specifies the game type ("rs3" or "osrs").
//
// itemAlpha specifies the first character of the item name to filter by.
//
// itemCategory specifies the category of items to retrieve. Reference for the Item Category Ids can be found on the
// [RuneScape Wiki]. itemCategory will default to 1 if gameType is set to "osrs"
//
// page specifies the page number of the results.
//
// [RuneScape Wiki]: https://runescape.wiki/w/Application_programming_interface#items
func (c *Client) ListGrandExchangeItems(gameType string, itemAlpha string, itemCategory int64, page int64) (*GeResponse, error) {
	validGameTypes := map[string]bool{
		"rs3":  true,
		"osrs": true,
	}

	if !validGameTypes[gameType] {
		return nil, errors.New("gameType must be \"rs3\" or \"osrs\"")
	}

	itemByte := []byte(itemAlpha)
	if len(itemByte) > 1 {
		return nil, errors.New("itemAlpha should be one character a-z or #")
	}

	charCode := int(itemByte[0])
	if (charCode < 97 && charCode != 35) || charCode > 122 {
		return nil, errors.New("invalid character")
	}

	if page < 1 {
		return nil, errors.New("page must be > 1")
	}

	var gameTypeString string
	var category int64
	if gameType == "rs3" {
		gameTypeString = "rs"
		category = itemCategory
	} else {
		gameTypeString = "oldschool"
		category = 1
	}
	path := c.BaseURL.JoinPath("/m=itemdb_"+gameTypeString, "/api/catalogue/items.json")

	v := url.Values{}
	v.Set("category", strconv.FormatInt(category, 10))
	v.Set("alpha", itemAlpha)
	v.Set("page", strconv.FormatInt(page, 10))

	req := path.String() + "?" + v.Encode()

	resp, err := c.client.Get(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response GeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
