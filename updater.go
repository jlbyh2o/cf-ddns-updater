package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// DDNSUpdater handles the DNS update process
type DDNSUpdater struct {
	config     *Config
	cfClient   *CloudflareClient
	ipDetector *IPDetector
	verbose    bool
}

// NewDDNSUpdater creates a new DDNS updater
func NewDDNSUpdater(config *Config, verbose bool) *DDNSUpdater {
	return &DDNSUpdater{
		config:     config,
		cfClient:   NewCloudflareClient(config.Cloudflare),
		ipDetector: NewIPDetector(),
		verbose:    verbose || config.Verbose,
	}
}

// Update performs the DNS update process
func (u *DDNSUpdater) Update() error {
	// Validate configuration
	if err := u.config.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	if u.verbose {
		log.Println("Starting DNS update process...")
	}

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
		} else if u.verbose {
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
		} else if u.verbose {
			log.Printf("Current IPv6 address: %s", ipv6)
		}
	}

	// Update each domain
	for _, domain := range u.config.Domains {
		if u.verbose {
			log.Printf("Processing domain: %s", domain.Name)
		}

		if err := u.updateDomain(domain, ipv4, ipv6); err != nil {
			log.Printf("Failed to update domain %s: %v", domain.Name, err)
			continue
		}

		if u.verbose {
			log.Printf("Successfully processed domain: %s", domain.Name)
		}
	}

	return nil
}

// updateDomain updates DNS records for a specific domain
func (u *DDNSUpdater) updateDomain(domain DomainConfig, ipv4, ipv6 string) error {
	// Get zone ID
	if u.verbose {
		log.Printf("Getting zone ID for domain: %s", extractRootDomain(domain.Name))
	}
	zoneID, err := u.cfClient.GetZoneID(extractRootDomain(domain.Name))
	if err != nil {
		return fmt.Errorf("failed to get zone ID: %w", err)
	}
	if u.verbose {
		log.Printf("Zone ID found: %s", zoneID)
	}

	// Update A record if needed
	if domain.ShouldUpdateA() && ipv4 != "" {
		if u.verbose {
			log.Printf("Checking A record for %s (current IP: %s)", domain.Name, ipv4)
		}
		if err := u.updateRecord(zoneID, domain, "A", ipv4); err != nil {
			return fmt.Errorf("failed to update A record: %w", err)
		}
	}

	// Update AAAA record if needed
	if domain.ShouldUpdateAAAA() && ipv6 != "" {
		if u.verbose {
			log.Printf("Checking AAAA record for %s (current IPv6: %s)", domain.Name, ipv6)
		}
		if err := u.updateRecord(zoneID, domain, "AAAA", ipv6); err != nil {
			return fmt.Errorf("failed to update AAAA record: %w", err)
		}
	}

	return nil
}

// updateRecord updates a specific DNS record
func (u *DDNSUpdater) updateRecord(zoneID string, domain DomainConfig, recordType, content string) error {
	if u.verbose {
		log.Printf("Checking %s record for %s (target IP: %s)", recordType, domain.Name, content)
	}

	// First, check what the domain currently resolves to via DNS
	if u.verbose {
		u.checkCurrentDNSResolution(domain.Name, recordType)
	}

	// Get existing DNS records from Cloudflare
	if u.verbose {
		log.Printf("Retrieving existing %s records for %s from Cloudflare API...", recordType, domain.Name)
	}
	existingRecords, err := u.cfClient.GetDNSRecords(zoneID, domain.Name, recordType)
	if err != nil {
		return fmt.Errorf("failed to get existing records: %w", err)
	}
	if u.verbose {
		log.Printf("Found %d existing %s record(s) for %s", len(existingRecords), recordType, domain.Name)
	}

	// Create new record data
	newRecord := DNSRecord{
		Type:    recordType,
		Name:    domain.Name,
		Content: content,
		TTL:     domain.TTL,
		Proxied: domain.Proxied,
	}

	// If record exists, compare and update if needed
	if len(existingRecords) > 0 {
		existingRecord := existingRecords[0]
		if u.verbose {
			log.Printf("Current %s record for %s: IP=%s, TTL=%d, Proxied=%t",
				recordType, domain.Name, existingRecord.Content, existingRecord.TTL, existingRecord.Proxied)
		}

		// Check if update is needed by comparing all relevant fields
		if existingRecord.Content == content && existingRecord.TTL == domain.TTL && existingRecord.Proxied == domain.Proxied {
			if u.verbose {
				log.Printf("%s record for %s is already up to date - no API call needed", recordType, domain.Name)
			}
			return nil
		}

		if u.verbose {
			log.Printf("DNS record needs update: Current IP (%s) != Target IP (%s) OR TTL/Proxy settings differ",
				existingRecord.Content, content)
		}

		// Update existing record via Cloudflare API
		log.Printf("Updating %s record for %s: %s to %s", recordType, domain.Name, existingRecord.Content, content)
		_, err = u.cfClient.UpdateDNSRecord(zoneID, existingRecord.ID, newRecord)
		if err != nil {
			return fmt.Errorf("failed to update existing record: %w", err)
		}

		log.Printf("Successfully updated %s record for %s", recordType, domain.Name)
	} else {
		// Create new record
		if u.verbose {
			log.Printf("No existing %s record found for %s, creating new record...", recordType, domain.Name)
		}
		log.Printf("Creating %s record for %s with IP %s", recordType, domain.Name, content)
		_, err = u.cfClient.CreateDNSRecord(zoneID, newRecord)
		if err != nil {
			return fmt.Errorf("failed to create new record: %w", err)
		}

		log.Printf("Successfully created %s record for %s", recordType, domain.Name)
	}

	return nil
}

// checkCurrentDNSResolution checks what the domain currently resolves to via DNS
func (u *DDNSUpdater) checkCurrentDNSResolution(domain, recordType string) {
	log.Printf("Performing DNS lookup to check current resolution for %s (%s record)...", domain, recordType)

	// Set a timeout for DNS resolution
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return d.DialContext(ctx, network, address)
		},
	}

	ctx := context.Background()

	if recordType == "A" {
		// Look up IPv4 addresses
		addrs, err := resolver.LookupIPAddr(ctx, domain)
		if err != nil {
			log.Printf("DNS lookup failed for %s (A record): %v", domain, err)
			return
		}

		var ipv4Addrs []string
		for _, addr := range addrs {
			if addr.IP.To4() != nil {
				ipv4Addrs = append(ipv4Addrs, addr.IP.String())
			}
		}

		if len(ipv4Addrs) > 0 {
			log.Printf("Current DNS resolution for %s (A): %v", domain, ipv4Addrs)
		} else {
			log.Printf("No A records found in DNS for %s", domain)
		}
	} else if recordType == "AAAA" {
		// Look up IPv6 addresses
		addrs, err := resolver.LookupIPAddr(ctx, domain)
		if err != nil {
			log.Printf("DNS lookup failed for %s (AAAA record): %v", domain, err)
			return
		}

		var ipv6Addrs []string
		for _, addr := range addrs {
			if addr.IP.To4() == nil && addr.IP.To16() != nil {
				ipv6Addrs = append(ipv6Addrs, addr.IP.String())
			}
		}

		if len(ipv6Addrs) > 0 {
			log.Printf("Current DNS resolution for %s (AAAA): %v", domain, ipv6Addrs)
		} else {
			log.Printf("No AAAA records found in DNS for %s", domain)
		}
	}
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
