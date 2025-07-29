// Package resolver provides advanced DNS resolution capabilities
package resolver

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// RecordType represents different DNS record types
type RecordType string

const (
	A     RecordType = "A"
	AAAA  RecordType = "AAAA"
	CNAME RecordType = "CNAME"
	MX    RecordType = "MX"
	NS    RecordType = "NS"
	TXT   RecordType = "TXT"
	SOA   RecordType = "SOA"
	PTR   RecordType = "PTR"
	SRV   RecordType = "SRV"
)

// DNSResult represents the result of a DNS query
type DNSResult struct {
	Domain      string        `json:"domain"`
	RecordType  RecordType    `json:"record_type"`
	Records     []string      `json:"records"`
	TTL         uint32        `json:"ttl"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Server      string        `json:"dns_server"`
	Error       string        `json:"error,omitempty"`
	Timestamp   time.Time     `json:"timestamp"`
}

// BulkResult represents results for multiple domain queries
type BulkResult struct {
	Domain  string       `json:"domain"`
	Results []*DNSResult `json:"results"`
	Error   string       `json:"error,omitempty"`
}

// ServerPerformance represents DNS server performance metrics
type ServerPerformance struct {
	Server       string        `json:"server"`
	AvgResponse  time.Duration `json:"avg_response_time_ms"`
	MinResponse  time.Duration `json:"min_response_time_ms"`
	MaxResponse  time.Duration `json:"max_response_time_ms"`
	SuccessRate  float64       `json:"success_rate"`
	TotalQueries int           `json:"total_queries"`
	Failures     int           `json:"failures"`
}

// Resolver provides advanced DNS resolution functionality
type Resolver struct {
	servers    []string
	timeout    time.Duration
	retries    int
	concurrent int
	client     *dns.Client
}

// NewResolver creates a new DNS resolver with custom settings
func NewResolver(servers []string, timeout time.Duration, retries int, concurrent int) *Resolver {
	if len(servers) == 0 {
		// Default DNS servers (Google, Cloudflare, Quad9)
		servers = []string{
			"8.8.8.8:53",
			"1.1.1.1:53",
			"9.9.9.9:53",
		}
	}

	// Ensure servers have port numbers
	for i, server := range servers {
		if !strings.Contains(server, ":") {
			servers[i] = server + ":53"
		}
	}

	return &Resolver{
		servers:    servers,
		timeout:    timeout,
		retries:    retries,
		concurrent: concurrent,
		client: &dns.Client{
			Timeout: timeout,
		},
	}
}

// Resolve performs DNS resolution for a domain with specified record type
func (r *Resolver) Resolve(domain string, recordType RecordType) (*DNSResult, error) {
	domain = strings.TrimSpace(strings.ToLower(domain))
	if domain == "" {
		return nil, fmt.Errorf("domain cannot be empty")
	}

	// Remove trailing dot if present, then add it back for proper DNS query
	domain = strings.TrimSuffix(domain, ".")
	queryDomain := domain + "."

	var qtype uint16
	switch recordType {
	case A:
		qtype = dns.TypeA
	case AAAA:
		qtype = dns.TypeAAAA
	case CNAME:
		qtype = dns.TypeCNAME
	case MX:
		qtype = dns.TypeMX
	case NS:
		qtype = dns.TypeNS
	case TXT:
		qtype = dns.TypeTXT
	case SOA:
		qtype = dns.TypeSOA
	case PTR:
		qtype = dns.TypePTR
	case SRV:
		qtype = dns.TypeSRV
	default:
		return nil, fmt.Errorf("unsupported record type: %s", recordType)
	}

	result := &DNSResult{
		Domain:     domain,
		RecordType: recordType,
		Records:    []string{},
		Timestamp:  time.Now(),
	}

	// Try each DNS server until we get a successful response
	for _, server := range r.servers {
		start := time.Now()
		
		msg := new(dns.Msg)
		msg.SetQuestion(queryDomain, qtype)
		msg.RecursionDesired = true

		response, _, err := r.client.Exchange(msg, server)
		responseTime := time.Since(start)
		
		if err != nil {
			continue // Try next server
		}

		result.Server = server
		result.ResponseTime = responseTime

		if response.Rcode != dns.RcodeSuccess {
			result.Error = dns.RcodeToString[response.Rcode]
			return result, nil
		}

		// Parse the response
		records := []string{}
		var ttl uint32 = 0

		for _, answer := range response.Answer {
			if ttl == 0 {
				ttl = answer.Header().Ttl
			}

			switch rr := answer.(type) {
			case *dns.A:
				if recordType == A {
					records = append(records, rr.A.String())
				}
			case *dns.AAAA:
				if recordType == AAAA {
					records = append(records, rr.AAAA.String())
				}
			case *dns.CNAME:
				if recordType == CNAME {
					records = append(records, strings.TrimSuffix(rr.Target, "."))
				}
			case *dns.MX:
				if recordType == MX {
					records = append(records, fmt.Sprintf("%d %s", rr.Preference, strings.TrimSuffix(rr.Mx, ".")))
				}
			case *dns.NS:
				if recordType == NS {
					records = append(records, strings.TrimSuffix(rr.Ns, "."))
				}
			case *dns.TXT:
				if recordType == TXT {
					records = append(records, strings.Join(rr.Txt, " "))
				}
			case *dns.SOA:
				if recordType == SOA {
					records = append(records, fmt.Sprintf("%s %s %d %d %d %d %d",
						strings.TrimSuffix(rr.Ns, "."),
						strings.TrimSuffix(rr.Mbox, "."),
						rr.Serial, rr.Refresh, rr.Retry, rr.Expire, rr.Minttl))
				}
			case *dns.PTR:
				if recordType == PTR {
					records = append(records, strings.TrimSuffix(rr.Ptr, "."))
				}
			case *dns.SRV:
				if recordType == SRV {
					records = append(records, fmt.Sprintf("%d %d %d %s",
						rr.Priority, rr.Weight, rr.Port, strings.TrimSuffix(rr.Target, ".")))
				}
			}
		}

		result.Records = records
		result.TTL = ttl
		return result, nil
	}

	result.Error = "all DNS servers failed to respond"
	return result, nil
}

// ResolveAll performs resolution for multiple record types for a domain
func (r *Resolver) ResolveAll(domain string, recordTypes []RecordType) ([]*DNSResult, error) {
	if len(recordTypes) == 0 {
		recordTypes = []RecordType{A, AAAA, CNAME, MX, NS, TXT}
	}

	results := make([]*DNSResult, 0, len(recordTypes))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Use semaphore to limit concurrent queries
	sem := make(chan struct{}, r.concurrent)

	for _, recordType := range recordTypes {
		wg.Add(1)
		go func(rt RecordType) {
			defer wg.Done()
			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			result, err := r.Resolve(domain, rt)
			if err == nil {
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}
		}(recordType)
	}

	wg.Wait()

	// Sort results by record type for consistent output
	sort.Slice(results, func(i, j int) bool {
		return string(results[i].RecordType) < string(results[j].RecordType)
	})

	return results, nil
}

// BulkResolve performs DNS resolution for multiple domains
func (r *Resolver) BulkResolve(domains []string, recordTypes []RecordType) ([]*BulkResult, error) {
	results := make([]*BulkResult, 0, len(domains))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Use semaphore to limit concurrent domain processing
	sem := make(chan struct{}, r.concurrent)

	for _, domain := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			domainResults, err := r.ResolveAll(d, recordTypes)
			
			bulkResult := &BulkResult{
				Domain:  d,
				Results: domainResults,
			}
			
			if err != nil {
				bulkResult.Error = err.Error()
			}

			mu.Lock()
			results = append(results, bulkResult)
			mu.Unlock()
		}(domain)
	}

	wg.Wait()

	// Sort results by domain name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Domain < results[j].Domain
	})

	return results, nil
}

// ReverseDNS performs reverse DNS lookup for an IP address
func (r *Resolver) ReverseDNS(ip string) (*DNSResult, error) {
	// Validate IP address
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	// Convert IP to reverse DNS format
	var reverseDomain string
	if parsedIP.To4() != nil {
		// IPv4
		parts := strings.Split(parsedIP.String(), ".")
		reverseDomain = fmt.Sprintf("%s.%s.%s.%s.in-addr.arpa.", parts[3], parts[2], parts[1], parts[0])
	} else {
		// IPv6
		reverseDomain = reverseIPv6(parsedIP.String()) + ".ip6.arpa."
	}

	return r.Resolve(reverseDomain, PTR)
}

// TestServers tests the performance of configured DNS servers
func (r *Resolver) TestServers(testDomain string, iterations int) ([]*ServerPerformance, error) {
	if testDomain == "" {
		testDomain = "google.com"
	}
	if iterations <= 0 {
		iterations = 5
	}

	performances := make([]*ServerPerformance, 0, len(r.servers))

	for _, server := range r.servers {
		perf := &ServerPerformance{
			Server:      server,
			MinResponse: time.Hour, // Initialize with large value
		}

		var totalTime time.Duration
		var responses []time.Duration

		for i := 0; i < iterations; i++ {
			start := time.Now()
			
			msg := new(dns.Msg)
			msg.SetQuestion(testDomain+".", dns.TypeA)
			msg.RecursionDesired = true

			_, _, err := r.client.Exchange(msg, server)
			responseTime := time.Since(start)

			if err != nil {
				perf.Failures++
			} else {
				responses = append(responses, responseTime)
				totalTime += responseTime
				
				if responseTime < perf.MinResponse {
					perf.MinResponse = responseTime
				}
				if responseTime > perf.MaxResponse {
					perf.MaxResponse = responseTime
				}
			}
			
			perf.TotalQueries++
		}

		if len(responses) > 0 {
			perf.AvgResponse = totalTime / time.Duration(len(responses))
			perf.SuccessRate = float64(len(responses)) / float64(iterations) * 100
		} else {
			perf.MinResponse = 0
		}

		performances = append(performances, perf)
	}

	// Sort by average response time
	sort.Slice(performances, func(i, j int) bool {
		return performances[i].AvgResponse < performances[j].AvgResponse
	})

	return performances, nil
}

// Helper function to reverse IPv6 address for PTR lookup
func reverseIPv6(ip string) string {
	// Remove colons and pad with zeros
	ip = strings.ReplaceAll(ip, ":", "")
	
	// Pad to 32 characters
	for len(ip) < 32 {
		ip = "0" + ip
	}
	
	// Reverse and add dots
	var result []string
	for i := len(ip) - 1; i >= 0; i-- {
		result = append(result, string(ip[i]))
	}
	
	return strings.Join(result, ".")
}

// TraceQuery performs a DNS query trace showing the resolution path
func (r *Resolver) TraceQuery(domain string, recordType RecordType) ([]*DNSResult, error) {
	domain = strings.TrimSpace(strings.ToLower(domain))
	domain = strings.TrimSuffix(domain, ".")
	
	var trace []*DNSResult
	currentDomain := domain
	visited := make(map[string]bool)
	
	for i := 0; i < 10; i++ { // Limit trace depth
		if visited[currentDomain] {
			break // Avoid infinite loops
		}
		visited[currentDomain] = true
		
		result, err := r.Resolve(currentDomain, recordType)
		if err != nil {
			return trace, err
		}
		
		trace = append(trace, result)
		
		// If we got records or an error, stop tracing
		if len(result.Records) > 0 || result.Error != "" {
			break
		}
		
		// Try to follow CNAME if this was an A/AAAA query
		if recordType == A || recordType == AAAA {
			cnameResult, _ := r.Resolve(currentDomain, CNAME)
			if cnameResult != nil && len(cnameResult.Records) > 0 {
				currentDomain = cnameResult.Records[0]
				continue
			}
		}
		
		break
	}
	
	return trace, nil
}