package utils

import (
	"math"
	"strconv"
)

// Default values
var (
	// DefaultPage default page
	DefaultPage = 1
	// DefaultLimit default limit
	DefaultLimit = 10
)

// Pagination struct
type Pagination struct {
	Page  int `json:"page"`  // Page number
	Limit int `json:"limit"` // Limit per page
	Total struct {
		ItemInPage int `json:"item_in_page"` // ItemInPage total items in a page
		Items      int `json:"items"`        // Items total items
		Pages      int `json:"pages"`        // Pages total pages
	} `json:"total"`
	HasPreviousPage bool `json:"has_previous_page"` // HasPreviousPage check if it has previous page
	HasNextPage     bool `json:"has_next_page"`     // HasNextPage check if it has next page
	Links           struct {
		Previous string `json:"previous_url"`
		Next     string `json:"next_url"`
	} `json:"links"`
}

// CalculateTotalPages method to calculate total pages
func (p *Pagination) CalculateTotalPages() {
	p.Total.Pages = int(math.Ceil(float64(p.Total.Items) / float64(p.Limit)))
}

// CalculatePrev method to calculate previous page
func (p *Pagination) CalculatePrev(url string) {
	p.HasPreviousPage = p.Page > 1 // Check if current page is greater than 1

	// If current page is greater than 1, then set a previous page.
	if p.HasPreviousPage {
		p.Links.Previous = url + "?page=" + strconv.Itoa(p.Page-1) + "&limit=" + strconv.Itoa(p.Limit)
	} else { // If current page is 1, then set a previous page to empty string.
		p.Links.Previous = ""
	}
}

// CalculateNext method to calculate next page
func (p *Pagination) CalculateNext(url string) {
	p.HasNextPage = p.Page < p.Total.Pages // Check if current page is less than total pages

	// If current page is less than total pages, then set a next page.
	if p.HasNextPage {
		p.Links.Next = url + "?page=" + strconv.Itoa(p.Page+1) + "&limit=" + strconv.Itoa(p.Limit)
	} else { // If current page is equal to total pages, then set a next page to empty string.
		p.Links.Next = ""
	}
}

// NewPagination function to create new pagination
func NewPagination(url string, page, limit, pageItems, totalItems int) *Pagination {
	// Create a new pagination instance
	pagination := &Pagination{
		Page:  page,
		Limit: limit,
		Total: struct {
			ItemInPage int `json:"item_in_page"`
			Items      int `json:"items"`
			Pages      int `json:"pages"`
		}{ItemInPage: pageItems, Items: totalItems},
		Links: struct {
			Previous string `json:"previous_url"`
			Next     string `json:"next_url"`
		}{},
	}

	// Calculate total pages, previous page, and next page
	pagination.CalculateTotalPages()
	pagination.CalculatePrev(url)
	pagination.CalculateNext(url)

	// Return the pagination
	return pagination
}
