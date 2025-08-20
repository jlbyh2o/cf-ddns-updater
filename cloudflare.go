/*
Cloudflare Dynamic DNS Updater
Copyright (C) 2025

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	cloudflareAPIBase = "https://api.cloudflare.com/client/v4"
)

// CloudflareClient handles Cloudflare API operations
type CloudflareClient struct {
	client *http.Client
	config CloudflareConfig
}

// NewCloudflareClient creates a new Cloudflare API client
func NewCloudflareClient(config CloudflareConfig) *CloudflareClient {
	return &CloudflareClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: config,
	}
}

// Zone represents a Cloudflare zone
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DNSRecord represents a Cloudflare DNS record
type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

// CloudflareResponse represents the standard Cloudflare API response
type CloudflareResponse struct {
	Success bool        `json:"success"`
	Errors  []CFError   `json:"errors"`
	Result  interface{} `json:"result"`
}

// CFError represents a Cloudflare API error
type CFError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetZoneID retrieves the zone ID for a domain
func (c *CloudflareClient) GetZoneID(domain string) (string, error) {
	if c.config.ZoneID != "" {
		return c.config.ZoneID, nil
	}

	url := fmt.Sprintf("%s/zones?name=%s", cloudflareAPIBase, domain)
	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	var zones []Zone
	if err := json.Unmarshal(resp, &zones); err != nil {
		return "", fmt.Errorf("failed to parse zones response: %w", err)
	}

	if len(zones) == 0 {
		return "", fmt.Errorf("zone not found for domain %s", domain)
	}

	return zones[0].ID, nil
}

// GetDNSRecords retrieves DNS records for a domain
func (c *CloudflareClient) GetDNSRecords(zoneID, name, recordType string) ([]DNSRecord, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records?name=%s&type=%s", cloudflareAPIBase, zoneID, name, recordType)
	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var records []DNSRecord
	if err := json.Unmarshal(resp, &records); err != nil {
		return nil, fmt.Errorf("failed to parse DNS records response: %w", err)
	}

	return records, nil
}

// CreateDNSRecord creates a new DNS record
func (c *CloudflareClient) CreateDNSRecord(zoneID string, record DNSRecord) (*DNSRecord, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records", cloudflareAPIBase, zoneID)

	payload, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DNS record: %w", err)
	}

	resp, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	var createdRecord DNSRecord
	if err := json.Unmarshal(resp, &createdRecord); err != nil {
		return nil, fmt.Errorf("failed to parse created DNS record response: %w", err)
	}

	return &createdRecord, nil
}

// UpdateDNSRecord updates an existing DNS record
func (c *CloudflareClient) UpdateDNSRecord(zoneID, recordID string, record DNSRecord) (*DNSRecord, error) {
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cloudflareAPIBase, zoneID, recordID)

	payload, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DNS record: %w", err)
	}

	resp, err := c.makeRequest("PUT", url, payload)
	if err != nil {
		return nil, err
	}

	var updatedRecord DNSRecord
	if err := json.Unmarshal(resp, &updatedRecord); err != nil {
		return nil, fmt.Errorf("failed to parse updated DNS record response: %w", err)
	}

	return &updatedRecord, nil
}

// makeRequest makes an HTTP request to the Cloudflare API
func (c *CloudflareClient) makeRequest(method, url string, body []byte) ([]byte, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Set authentication
	if c.config.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIToken)
	} else {
		req.Header.Set("X-Auth-Key", c.config.APIKey)
		req.Header.Set("X-Auth-Email", c.config.Email)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var cfResp CloudflareResponse
	unmarshalErr := json.Unmarshal(respBody, &cfResp)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !cfResp.Success {
		if len(cfResp.Errors) > 0 {
			return nil, fmt.Errorf("cloudflare API error: %s (code: %d)", cfResp.Errors[0].Message, cfResp.Errors[0].Code)
		}
		return nil, fmt.Errorf("cloudflare API request failed")
	}

	// Return the result as JSON
	result, err := json.Marshal(cfResp.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return result, nil
}
