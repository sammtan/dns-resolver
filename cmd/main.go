// Advanced DNS Resolver - Command Line Interface
// Educational Security Tool for DNS Analysis and Testing
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sammtan/dns-resolver/pkg/resolver"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	servers    []string
	timeout    time.Duration
	retries    int
	concurrent int
	format     string
	output     string
	verbose    bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dns-resolver",
		Short: "Advanced DNS Resolver - Educational Security Tool",
		Long: `
╔══════════════════════════════════════════════════════════════╗
║                   ADVANCED DNS RESOLVER                      ║
║                     Educational Tool v1.0.0                 ║
║                                                              ║
║  WARNING: Use only on domains you own or have permission    ║
║          to test. Unauthorized scanning may be illegal!     ║
╚══════════════════════════════════════════════════════════════╝

Advanced DNS Resolver provides comprehensive DNS analysis capabilities
including multiple record type queries, bulk domain processing, reverse
DNS lookups, and DNS server performance testing.

Features:
  • Multiple DNS record types (A, AAAA, CNAME, MX, NS, TXT, SOA, PTR, SRV)
  • Bulk domain processing with concurrent queries
  • Reverse DNS lookups for IP addresses
  • DNS server performance testing and comparison
  • Query tracing for debugging DNS resolution paths
  • Multiple output formats (text, JSON, CSV)
  • Educational focus with detailed explanations

Educational Use Only - Authorized Testing Required`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(servers) == 0 {
				servers = []string{"8.8.8.8", "1.1.1.1", "9.9.9.9"}
			}
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringSliceVarP(&servers, "servers", "s", []string{}, "DNS servers to query (default: 8.8.8.8,1.1.1.1,9.9.9.9)")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 5*time.Second, "Query timeout duration")
	rootCmd.PersistentFlags().IntVarP(&retries, "retries", "r", 3, "Number of retries per query")
	rootCmd.PersistentFlags().IntVarP(&concurrent, "concurrent", "c", 10, "Maximum concurrent queries")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "text", "Output format (text, json, csv)")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Add subcommands
	rootCmd.AddCommand(createResolveCommand())
	rootCmd.AddCommand(createBulkCommand())
	rootCmd.AddCommand(createReverseCommand())
	rootCmd.AddCommand(createTestCommand())
	rootCmd.AddCommand(createTraceCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func createResolveCommand() *cobra.Command {
	var recordTypes []string
	
	cmd := &cobra.Command{
		Use:   "resolve [domain]",
		Short: "Resolve DNS records for a domain",
		Long: `Resolve DNS records for a single domain. Supports multiple record types
and provides detailed information about each record including TTL values
and response times.

Examples:
  dns-resolver resolve google.com
  dns-resolver resolve google.com --types A,AAAA,MX
  dns-resolver resolve example.com --format json --output results.json`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			
			// Parse record types
			var types []resolver.RecordType
			if len(recordTypes) == 0 {
				types = []resolver.RecordType{resolver.A, resolver.AAAA, resolver.CNAME, resolver.MX, resolver.NS, resolver.TXT}
			} else {
				for _, rt := range recordTypes {
					types = append(types, resolver.RecordType(strings.ToUpper(rt)))
				}
			}
			
			// Create resolver
			r := resolver.NewResolver(servers, timeout, retries, concurrent)
			
			// Perform resolution
			results, err := r.ResolveAll(domain, types)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error resolving domain: %v\n", err)
				os.Exit(1)
			}
			
			// Output results
			outputResults(results, format, output)
		},
	}
	
	cmd.Flags().StringSliceVar(&recordTypes, "types", []string{}, "Record types to query (A,AAAA,CNAME,MX,NS,TXT,SOA,PTR,SRV)")
	
	return cmd
}

func createBulkCommand() *cobra.Command {
	var inputFile string
	var recordTypes []string
	
	cmd := &cobra.Command{
		Use:   "bulk [domains...]",
		Short: "Perform bulk DNS resolution for multiple domains",
		Long: `Perform DNS resolution for multiple domains simultaneously with
concurrent processing for improved performance.

Examples:
  dns-resolver bulk google.com facebook.com twitter.com
  dns-resolver bulk --input domains.txt --types A,MX
  dns-resolver bulk --input domains.txt --format csv --output results.csv`,
		Run: func(cmd *cobra.Command, args []string) {
			var domains []string
			
			// Get domains from arguments or file
			if inputFile != "" {
				data, err := os.ReadFile(inputFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
					os.Exit(1)
				}
				domains = strings.Fields(string(data))
			} else if len(args) > 0 {
				domains = args
			} else {
				fmt.Fprintf(os.Stderr, "Error: No domains provided. Use arguments or --input file\n")
				os.Exit(1)
			}
			
			// Parse record types
			var types []resolver.RecordType
			if len(recordTypes) == 0 {
				types = []resolver.RecordType{resolver.A, resolver.AAAA, resolver.MX}
			} else {
				for _, rt := range recordTypes {
					types = append(types, resolver.RecordType(strings.ToUpper(rt)))
				}
			}
			
			// Create resolver
			r := resolver.NewResolver(servers, timeout, retries, concurrent)
			
			if verbose {
				fmt.Printf("[INFO] Processing %d domains with %d concurrent workers\n", len(domains), concurrent)
			}
			
			// Perform bulk resolution
			results, err := r.BulkResolve(domains, types)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error performing bulk resolution: %v\n", err)
				os.Exit(1)
			}
			
			// Output results
			outputBulkResults(results, format, output)
		},
	}
	
	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file containing domains (one per line)")
	cmd.Flags().StringSliceVar(&recordTypes, "types", []string{}, "Record types to query")
	
	return cmd
}

func createReverseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reverse [ip]",
		Short: "Perform reverse DNS lookup for an IP address",
		Long: `Perform reverse DNS lookup to find the hostname associated with
an IP address. Supports both IPv4 and IPv6 addresses.

Examples:
  dns-resolver reverse 8.8.8.8
  dns-resolver reverse 2001:4860:4860::8888
  dns-resolver reverse 192.168.1.1 --format json`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ip := args[0]
			
			// Create resolver
			r := resolver.NewResolver(servers, timeout, retries, concurrent)
			
			// Perform reverse DNS lookup
			result, err := r.ReverseDNS(ip)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error performing reverse DNS lookup: %v\n", err)
				os.Exit(1)
			}
			
			// Output result
			outputResults([]*resolver.DNSResult{result}, format, output)
		},
	}
	
	return cmd
}

func createTestCommand() *cobra.Command {
	var testDomain string
	var iterations int
	
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test DNS server performance",
		Long: `Test the performance of configured DNS servers by measuring
response times and success rates across multiple queries.

Examples:
  dns-resolver test
  dns-resolver test --domain example.com --iterations 10
  dns-resolver test --servers 8.8.8.8,1.1.1.1 --format json`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create resolver
			r := resolver.NewResolver(servers, timeout, retries, concurrent)
			
			if verbose {
				fmt.Printf("[INFO] Testing %d DNS servers with %d iterations each\n", len(servers), iterations)
			}
			
			// Test server performance
			results, err := r.TestServers(testDomain, iterations)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error testing DNS servers: %v\n", err)
				os.Exit(1)
			}
			
			// Output results
			outputServerPerformance(results, format, output)
		},
	}
	
	cmd.Flags().StringVar(&testDomain, "domain", "google.com", "Domain to use for testing")
	cmd.Flags().IntVar(&iterations, "iterations", 5, "Number of test iterations per server")
	
	return cmd
}

func createTraceCommand() *cobra.Command {
	var recordType string
	
	cmd := &cobra.Command{
		Use:   "trace [domain]",
		Short: "Trace DNS query resolution path",
		Long: `Trace the DNS resolution path for a domain to debug DNS issues
and understand how queries are resolved through the DNS hierarchy.

Examples:
  dns-resolver trace google.com
  dns-resolver trace example.com --type AAAA
  dns-resolver trace subdomain.example.com --format json`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]
			rt := resolver.RecordType(strings.ToUpper(recordType))
			
			// Create resolver
			r := resolver.NewResolver(servers, timeout, retries, concurrent)
			
			if verbose {
				fmt.Printf("[INFO] Tracing DNS resolution path for %s (%s)\n", domain, rt)
			}
			
			// Perform trace
			results, err := r.TraceQuery(domain, rt)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error tracing DNS query: %v\n", err)
				os.Exit(1)
			}
			
			// Output results
			outputTraceResults(results, format, output)
		},
	}
	
	cmd.Flags().StringVar(&recordType, "type", "A", "Record type to trace")
	
	return cmd
}

// Output formatting functions
func outputResults(results []*resolver.DNSResult, format, output string) {
	var data []byte
	var err error
	
	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(results, "", "  ")
	case "csv":
		data, err = formatCSV(results)
	default:
		data = []byte(formatText(results))
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}
	
	writeOutput(data, output)
}

func outputBulkResults(results []*resolver.BulkResult, format, output string) {
	var data []byte
	var err error
	
	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(results, "", "  ")
	case "csv":
		// Flatten bulk results for CSV
		var flatResults []*resolver.DNSResult
		for _, bulk := range results {
			flatResults = append(flatResults, bulk.Results...)
		}
		data, err = formatCSV(flatResults)
	default:
		data = []byte(formatBulkText(results))
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}
	
	writeOutput(data, output)
}

func outputServerPerformance(results []*resolver.ServerPerformance, format, output string) {
	var data []byte
	var err error
	
	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(results, "", "  ")
	case "csv":
		data, err = formatPerformanceCSV(results)
	default:
		data = []byte(formatPerformanceText(results))
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}
	
	writeOutput(data, output)
}

func outputTraceResults(results []*resolver.DNSResult, format, output string) {
	var data []byte
	var err error
	
	switch strings.ToLower(format) {
	case "json":
		data, err = json.MarshalIndent(results, "", "  ")
	case "csv":
		data, err = formatCSV(results)
	default:
		data = []byte(formatTraceText(results))
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}
	
	writeOutput(data, output)
}

func writeOutput(data []byte, output string) {
	if output != "" {
		// Ensure output path is in the tools directory, not portfolio directory
		if !strings.HasPrefix(output, "/") && !strings.Contains(output, ":") {
			// Relative path - save to tools directory
			toolsDir := "../sammtan.github.io-tools/dns-resolver"
			output = toolsDir + "/" + output
		}
		
		err := os.WriteFile(output, data, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[INFO] Results saved to %s\n", output)
	} else {
		fmt.Print(string(data))
	}
}

// Text formatting functions
func formatText(results []*resolver.DNSResult) string {
	var output strings.Builder
	
	output.WriteString("============================================================\n")
	output.WriteString("                DNS RESOLUTION RESULTS\n")
	output.WriteString("============================================================\n\n")
	
	for _, result := range results {
		output.WriteString(fmt.Sprintf("Domain: %s\n", result.Domain))
		output.WriteString(fmt.Sprintf("Record Type: %s\n", result.RecordType))
		output.WriteString(fmt.Sprintf("DNS Server: %s\n", result.Server))
		output.WriteString(fmt.Sprintf("Response Time: %v\n", result.ResponseTime))
		output.WriteString(fmt.Sprintf("TTL: %d seconds\n", result.TTL))
		
		if result.Error != "" {
			output.WriteString(fmt.Sprintf("Error: %s\n", result.Error))
		} else if len(result.Records) > 0 {
			output.WriteString("Records:\n")
			for _, record := range result.Records {
				output.WriteString(fmt.Sprintf("  %s\n", record))
			}
		} else {
			output.WriteString("No records found\n")
		}
		
		output.WriteString(fmt.Sprintf("Timestamp: %s\n", result.Timestamp.Format(time.RFC3339)))
		output.WriteString("\n" + strings.Repeat("-", 60) + "\n\n")
	}
	
	output.WriteString("DISCLAIMER: This tool is for educational and authorized testing only.\n")
	return output.String()
}

func formatBulkText(results []*resolver.BulkResult) string {
	var output strings.Builder
	
	output.WriteString("============================================================\n")
	output.WriteString("               BULK DNS RESOLUTION RESULTS\n")
	output.WriteString("============================================================\n\n")
	
	for _, bulk := range results {
		output.WriteString(fmt.Sprintf("Domain: %s\n", bulk.Domain))
		
		if bulk.Error != "" {
			output.WriteString(fmt.Sprintf("Error: %s\n", bulk.Error))
		} else {
			for _, result := range bulk.Results {
				output.WriteString(fmt.Sprintf("  %s: ", result.RecordType))
				if result.Error != "" {
					output.WriteString(fmt.Sprintf("Error - %s\n", result.Error))
				} else if len(result.Records) > 0 {
					output.WriteString(fmt.Sprintf("%s (TTL: %ds, %v)\n", 
						strings.Join(result.Records, ", "), result.TTL, result.ResponseTime))
				} else {
					output.WriteString("No records\n")
				}
			}
		}
		
		output.WriteString("\n")
	}
	
	output.WriteString("DISCLAIMER: This tool is for educational and authorized testing only.\n")
	return output.String()
}

func formatPerformanceText(results []*resolver.ServerPerformance) string {
	var output strings.Builder
	
	output.WriteString("============================================================\n")
	output.WriteString("              DNS SERVER PERFORMANCE TEST\n")
	output.WriteString("============================================================\n\n")
	
	output.WriteString(fmt.Sprintf("%-20s %-12s %-12s %-12s %-12s %-8s\n", 
		"SERVER", "AVG_TIME", "MIN_TIME", "MAX_TIME", "SUCCESS%", "QUERIES"))
	output.WriteString(strings.Repeat("-", 80) + "\n")
	
	for _, perf := range results {
		output.WriteString(fmt.Sprintf("%-20s %-12v %-12v %-12v %-12.1f %-8d\n",
			perf.Server,
			perf.AvgResponse.Truncate(time.Millisecond),
			perf.MinResponse.Truncate(time.Millisecond),
			perf.MaxResponse.Truncate(time.Millisecond),
			perf.SuccessRate,
			perf.TotalQueries))
	}
	
	output.WriteString("\nDISCLAIMER: This tool is for educational and authorized testing only.\n")
	return output.String()
}

func formatTraceText(results []*resolver.DNSResult) string {
	var output strings.Builder
	
	output.WriteString("============================================================\n")
	output.WriteString("                DNS QUERY TRACE RESULTS\n")
	output.WriteString("============================================================\n\n")
	
	for i, result := range results {
		output.WriteString(fmt.Sprintf("Step %d: %s (%s)\n", i+1, result.Domain, result.RecordType))
		output.WriteString(fmt.Sprintf("  Server: %s\n", result.Server))
		output.WriteString(fmt.Sprintf("  Response Time: %v\n", result.ResponseTime))
		
		if result.Error != "" {
			output.WriteString(fmt.Sprintf("  Error: %s\n", result.Error))
		} else if len(result.Records) > 0 {
			output.WriteString("  Records:\n")
			for _, record := range result.Records {
				output.WriteString(fmt.Sprintf("    %s\n", record))
			}
		}
		output.WriteString("\n")
	}
	
	output.WriteString("DISCLAIMER: This tool is for educational and authorized testing only.\n")
	return output.String()
}

// CSV formatting functions
func formatCSV(results []*resolver.DNSResult) ([]byte, error) {
	var output strings.Builder
	writer := csv.NewWriter(&output)
	
	// Write header
	writer.Write([]string{"Domain", "RecordType", "Records", "TTL", "ResponseTime", "Server", "Error", "Timestamp"})
	
	// Write data
	for _, result := range results {
		writer.Write([]string{
			result.Domain,
			string(result.RecordType),
			strings.Join(result.Records, "; "),
			fmt.Sprintf("%d", result.TTL),
			result.ResponseTime.String(),
			result.Server,
			result.Error,
			result.Timestamp.Format(time.RFC3339),
		})
	}
	
	writer.Flush()
	return []byte(output.String()), writer.Error()
}

func formatPerformanceCSV(results []*resolver.ServerPerformance) ([]byte, error) {
	var output strings.Builder
	writer := csv.NewWriter(&output)
	
	// Write header
	writer.Write([]string{"Server", "AvgResponseTime", "MinResponseTime", "MaxResponseTime", "SuccessRate", "TotalQueries", "Failures"})
	
	// Write data
	for _, perf := range results {
		writer.Write([]string{
			perf.Server,
			perf.AvgResponse.String(),
			perf.MinResponse.String(),
			perf.MaxResponse.String(),
			fmt.Sprintf("%.2f", perf.SuccessRate),
			fmt.Sprintf("%d", perf.TotalQueries),
			fmt.Sprintf("%d", perf.Failures),
		})
	}
	
	writer.Flush()
	return []byte(output.String()), writer.Error()
}