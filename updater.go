package main

import (
	"fmt"
	"log"
	"strings"
)

// DDNSUpdater handles the DNS update process
type DDNSUpdater struct {
	config     *Config
	cfClient   *CloudflareClient
	ipDetector *IPDetector
}

// NewDDNSUpdater creates a new DDNS updater
func NewDDNSUpdater(config *Config) *DDNSUpdater {
	return &DDNSUpdater{
		config:     config,
		cfClient:   NewCloudflareClient(config.Cloudflare),
		ipDetector: NewIPDetector(),
	}
}

// Update performs the DNS update process
func (u *DDNSUpdater) Update() error {
	// Validate configuration
	if err := u.config.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	
	log.Println("Starting DNS update process...")
	
	// Get current IP addresses
	var ipv4, ipv6 string
	var err error
	
	// Check if any domain needs IPv4 updates
	needsIPv4 := false
	for _, domain := range u.config.Domains {
		if domain.ShouldUpdateA() {
			needsIPv4 = true
			break
		}
	}
	
	if needsIPv4 {
		ipv4, err = u.ipDetector.GetIPv4()
		if err != nil {
			log.Printf("Warning: Failed to get IPv4 address: %v", err)
		} else {
			log.Printf("Current IPv4 address: %s", ipv4)
		}
	}
	
	// Check if any domain needs IPv6 updates
	needsIPv6 := false
	for _, domain := range u.config.Domains {
		if domain.ShouldUpdateAAAA() {
			needsIPv6 = true
			break
		}
	}
	
	if needsIPv6 {
		ipv6, err = u.ipDetector.GetIPv6()
		if err != nil {
			log.Printf("Warning: Failed to get IPv6 address: %v", err)
		} else {
			log.Printf("Current IPv6 address: %s", ipv6)
		}
	}
	
	// Update each domain
	for _, domain := range u.config.Domains {
		log.Printf("Processing domain: %s", domain.Name)
		
		if err := u.updateDomain(domain, ipv4, ipv6); err != nil {
			log.Printf("Failed to update domain %s: %v", domain.Name, err)
			continue
		}
		
		log.Printf("Successfully updated domain: %s", domain.Name)
	}
	
	return nil
}

// updateDomain updates DNS records for a specific domain
func (u *DDNSUpdater) updateDomain(domain DomainConfig, ipv4, ipv6 string) error {
	// Get zone ID
	zoneID, err := u.cfClient.GetZoneID(extractRootDomain(domain.Name))
	if err != nil {
		return fmt.Errorf("failed to get zone ID: %w", err)
	}
	
	// Update A record if needed
	if domain.ShouldUpdateA() && ipv4 != "" {
		if err := u.updateRecord(zoneID, domain, "A", ipv4); err != nil {
			return fmt.Errorf("failed to update A record: %w", err)
		}
	}
	
	// Update AAAA record if needed
	if domain.ShouldUpdateAAAA() && ipv6 != "" {
		if err := u.updateRecord(zoneID, domain, "AAAA", ipv6); err != nil {
			return fmt.Errorf("failed to update AAAA record: %w", err)
		}
	}
	
	return nil
}

// updateRecord updates a specific DNS record
func (u *DDNSUpdater) updateRecord(zoneID string, domain DomainConfig, recordType, content string) error {
	log.Printf("Updating %s record for %s to %s", recordType, domain.Name, content)
	
	// Get existing records
	existingRecords, err := u.cfClient.GetDNSRecords(zoneID, domain.Name, recordType)
	if err != nil {
		return fmt.Errorf("failed to get existing records: %w", err)
	}
	
	// Create new record data
	newRecord := DNSRecord{
		Type:    recordType,
		Name:    domain.Name,
		Content: content,
		TTL:     domain.TTL,
		Proxied: domain.Proxied,
	}
	
	// If record exists, update it
	if len(existingRecords) > 0 {
		existingRecord := existingRecords[0]
		
		// Check if update is needed
		if existingRecord.Content == content && existingRecord.TTL == domain.TTL && existingRecord.Proxied == domain.Proxied {
			log.Printf("%s record for %s is already up to date", recordType, domain.Name)
			return nil
		}
		
		// Update existing record
		_, err = u.cfClient.UpdateDNSRecord(zoneID, existingRecord.ID, newRecord)
		if err != nil {
			return fmt.Errorf("failed to update existing record: %w", err)
		}
		
		log.Printf("Updated %s record for %s from %s to %s", recordType, domain.Name, existingRecord.Content, content)
	} else {
		// Create new record
		_, err = u.cfClient.CreateDNSRecord(zoneID, newRecord)
		if err != nil {
			return fmt.Errorf("failed to create new record: %w", err)
		}
		
		log.Printf("Created new %s record for %s with value %s", recordType, domain.Name, content)
	}
	
	return nil
}

// extractRootDomain extracts the root domain from a subdomain
func extractRootDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) <= 2 {
		return domain
	}
	
	// Return the last two parts (root domain)
	return strings.Join(parts[len(parts)-2:], ".")
}