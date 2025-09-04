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
	if err := u.validateConfig(); err != nil {
		return err
	}

	u.logStart()

	ipv4, ipv6, err := u.getRequiredIPs()
	if err != nil {
		return err
	}

	return u.updateAllDomains(ipv4, ipv6)
}

// validateConfig validates the configuration
func (u *DDNSUpdater) validateConfig() error {
	if err := u.config.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	return nil
}

// logStart logs the start of the update process
func (u *DDNSUpdater) logStart() {
	if u.verbose {
		log.Println("Starting DNS update process...")
	}
}

// getRequiredIPs determines which IP addresses are needed and fetches them
func (u *DDNSUpdater) getRequiredIPs() (ipv4, ipv6 string, err error) {
	needsIPv4 := u.needsIPv4()
	needsIPv6 := u.needsIPv6()

	if needsIPv4 {
		ipv4, err = u.getIPv4WithLogging()
		if err != nil {
			return "", "", err
		}
	}

	if needsIPv6 {
		ipv6, err = u.getIPv6WithLogging()
		if err != nil {
			return ipv4, "", err
		}
	}

	return ipv4, ipv6, nil
}

// needsIPv4 checks if any domain needs IPv4 updates
func (u *DDNSUpdater) needsIPv4() bool {
	for _, domain := range u.config.Domains {
		if domain.ShouldUpdateA() {
			return true
		}
	}
	return false
}

// needsIPv6 checks if any domain needs IPv6 updates
func (u *DDNSUpdater) needsIPv6() bool {
	for _, domain := range u.config.Domains {
		if domain.ShouldUpdateAAAA() {
			return true
		}
	}
	return false
}

// getIPv4WithLogging gets IPv4 address with appropriate logging
func (u *DDNSUpdater) getIPv4WithLogging() (string, error) {
	ipv4, err := u.ipDetector.GetIPv4()
	if err != nil {
		log.Printf("Warning: Failed to get IPv4 address: %v", err)
		return "", nil // Return empty string but no error to continue processing
	}
	if u.verbose {
		log.Printf("Current IPv4 address: %s", ipv4)
	}
	return ipv4, nil
}

// getIPv6WithLogging gets IPv6 address with appropriate logging
func (u *DDNSUpdater) getIPv6WithLogging() (string, error) {
	ipv6, err := u.ipDetector.GetIPv6()
	if err != nil {
		log.Printf("Warning: Failed to get IPv6 address: %v", err)
		return "", nil // Return empty string but no error to continue processing
	}
	if u.verbose {
		log.Printf("Current IPv6 address: %s", ipv6)
	}
	return ipv6, nil
}

// updateAllDomains updates all configured domains
func (u *DDNSUpdater) updateAllDomains(ipv4, ipv6 string) error {
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
	u.logRecordCheck(recordType, domain.Name, content)

	if u.verbose {
		u.checkCurrentDNSResolution(domain.Name, recordType)
	}

	existingRecords, err := u.getExistingRecords(zoneID, domain.Name, recordType)
	if err != nil {
		return err
	}

	newRecord := u.createNewRecord(domain, recordType, content)

	if len(existingRecords) > 0 {
		return u.handleExistingRecord(zoneID, existingRecords[0], newRecord, recordType, domain.Name, content)
	}

	return u.createRecord(zoneID, newRecord, recordType, domain.Name, content)
}

// logRecordCheck logs the initial record check
func (u *DDNSUpdater) logRecordCheck(recordType, domainName, content string) {
	if u.verbose {
		log.Printf("Checking %s record for %s (target IP: %s)", recordType, domainName, content)
	}
}

// getExistingRecords retrieves existing DNS records from Cloudflare
func (u *DDNSUpdater) getExistingRecords(zoneID, domainName, recordType string) ([]DNSRecord, error) {
	if u.verbose {
		log.Printf("Retrieving existing %s records for %s from Cloudflare API...", recordType, domainName)
	}

	existingRecords, err := u.cfClient.GetDNSRecords(zoneID, domainName, recordType)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing records: %w", err)
	}

	if u.verbose {
		log.Printf("Found %d existing %s record(s) for %s", len(existingRecords), recordType, domainName)
	}

	return existingRecords, nil
}

// createNewRecord creates a new DNS record structure
func (u *DDNSUpdater) createNewRecord(domain DomainConfig, recordType, content string) DNSRecord {
	return DNSRecord{
		Type:    recordType,
		Name:    domain.Name,
		Content: content,
		TTL:     domain.TTL,
		Proxied: domain.Proxied,
	}
}

// handleExistingRecord handles updating an existing DNS record
func (u *DDNSUpdater) handleExistingRecord(zoneID string, existingRecord DNSRecord, newRecord DNSRecord, recordType, domainName, content string) error {
	if u.verbose {
		log.Printf("Current %s record for %s: IP=%s, TTL=%d, Proxied=%t",
			recordType, domainName, existingRecord.Content, existingRecord.TTL, existingRecord.Proxied)
	}

	if u.recordNeedsUpdate(existingRecord, newRecord) {
		return u.updateExistingRecord(zoneID, existingRecord, newRecord, recordType, domainName, content)
	}

	if u.verbose {
		log.Printf("%s record for %s is already up to date - no API call needed", recordType, domainName)
	}
	return nil
}

// recordNeedsUpdate checks if a record needs to be updated
func (u *DDNSUpdater) recordNeedsUpdate(existing, new DNSRecord) bool {
	return existing.Content != new.Content ||
		existing.TTL != new.TTL ||
		existing.Proxied != new.Proxied
}

// updateExistingRecord updates an existing DNS record
func (u *DDNSUpdater) updateExistingRecord(zoneID string, existingRecord DNSRecord, newRecord DNSRecord, recordType, domainName, content string) error {
	if u.verbose {
		log.Printf("DNS record needs update: Current IP (%s) != Target IP (%s) OR TTL/Proxy settings differ",
			existingRecord.Content, content)
	}

	log.Printf("Updating %s record for %s: %s to %s", recordType, domainName, existingRecord.Content, content)
	_, err := u.cfClient.UpdateDNSRecord(zoneID, existingRecord.ID, newRecord)
	if err != nil {
		return fmt.Errorf("failed to update existing record: %w", err)
	}

	log.Printf("Successfully updated %s record for %s", recordType, domainName)
	return nil
}

// createRecord creates a new DNS record
func (u *DDNSUpdater) createRecord(zoneID string, newRecord DNSRecord, recordType, domainName, content string) error {
	if u.verbose {
		log.Printf("No existing %s record found for %s, creating new record...", recordType, domainName)
	}

	log.Printf("Creating %s record for %s with IP %s", recordType, domainName, content)
	_, err := u.cfClient.CreateDNSRecord(zoneID, newRecord)
	if err != nil {
		return fmt.Errorf("failed to create new record: %w", err)
	}

	log.Printf("Successfully created %s record for %s", recordType, domainName)
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
