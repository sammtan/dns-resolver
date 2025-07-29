# DNS Resolver Usage Examples

This document provides practical examples of using the Advanced DNS Resolver tool for various educational scenarios.

## Basic DNS Analysis

### Simple Domain Resolution
```bash
# Analyze a single domain with default settings
./dns-resolver resolve google.com
```

### Specific Record Types
```bash
# Check only A and MX records
./dns-resolver resolve example.com --types A,MX

# Check email configuration
./dns-resolver resolve company.com --types MX,TXT

# Check nameservers
./dns-resolver resolve domain.com --types NS,SOA
```

## Security Analysis Examples

### Email Security Analysis
```bash
# Check SPF, DKIM, and DMARC records
./dns-resolver resolve company.com --types TXT --format json | grep -E "spf|dkim|dmarc"

# Analyze mail exchange servers
./dns-resolver resolve company.com --types MX --format csv --output mail-servers.csv
```

### Infrastructure Discovery
```bash
# Map out domain infrastructure
./dns-resolver resolve target-domain.com --types A,AAAA,CNAME,NS --format json --output infrastructure.json

# Check for IPv6 support
./dns-resolver resolve modernsite.com --types AAAA
```

## Bulk Analysis Scenarios

### Competitor Analysis
```bash
# Analyze multiple competitor domains
./dns-resolver bulk competitor1.com competitor2.com competitor3.com --types A,MX --format csv --output competitors.csv
```

### Domain Portfolio Management
```bash
# Check all company domains from file
./dns-resolver bulk --input company-domains.txt --types A,MX,NS --format json --output portfolio-status.json
```

### Subdomain Analysis
```bash
# Analyze common subdomains
./dns-resolver bulk www.example.com mail.example.com ftp.example.com api.example.com --types A,CNAME
```

## Performance Testing Examples

### DNS Server Comparison
```bash
# Compare different DNS providers
./dns-resolver test --servers 8.8.8.8,1.1.1.1,9.9.9.9,208.67.222.222 --iterations 10 --format csv --output dns-comparison.csv
```

### Geographic Performance Testing
```bash
# Test regional DNS servers
./dns-resolver test --servers 8.8.8.8,1.1.1.1,9.9.9.9 --domain local-site.com --iterations 20
```

### Load Testing Simulation
```bash
# High-concurrency testing
./dns-resolver test --concurrent 50 --iterations 100 --domain test-site.com
```

## Reverse DNS Examples

### IP Range Analysis
```bash
# Analyze IP block ownership
./dns-resolver reverse 8.8.8.8
./dns-resolver reverse 8.8.4.4
./dns-resolver reverse 1.1.1.1
```

### Server Identification
```bash
# Identify servers behind IP addresses
./dns-resolver reverse 192.168.1.1 --format json
./dns-resolver reverse 10.0.0.1 --format csv --output internal-servers.csv
```

## Advanced Query Tracing

### DNS Resolution Path Analysis
```bash
# Trace how a domain resolves
./dns-resolver trace complex-domain.com --format json --output trace-analysis.json

# Debug CNAME chains
./dns-resolver trace www.redirected-site.com --type CNAME
```

### Troubleshooting DNS Issues
```bash
# Debug slow resolution
./dns-resolver trace slow-site.com --verbose

# Analyze resolution failures
./dns-resolver trace broken-domain.com --type A --format json
```

## Educational Scenarios

### DNS Learning Lab
```bash
# Compare different record types for the same domain
./dns-resolver resolve university.edu --types A,AAAA,MX,NS,TXT,SOA --format json --output learning-lab.json
```

### Network Administration Practice
```bash
# Simulate network monitoring
./dns-resolver bulk --input critical-services.txt --types A,MX --concurrent 5 --output monitoring-results.csv
```

### Security Research
```bash
# Analyze DNS infrastructure patterns
./dns-resolver bulk --input research-domains.txt --types NS,SOA --format json --output research-data.json
```

## Automation Examples

### Bash Script Integration
```bash
#!/bin/bash
# Monitor domain changes
DOMAIN="important-site.com"
TODAY=$(date +%Y%m%d)
./dns-resolver resolve $DOMAIN --format json --output "monitoring-$TODAY.json"
```

### PowerShell Integration
```powershell
# Windows automation example
$domains = @("site1.com", "site2.com", "site3.com")
foreach ($domain in $domains) {
    .\dns-resolver.exe resolve $domain --format csv --output "results-$domain.csv"
}
```

### Python Integration
```python
import subprocess
import json

# Run DNS analysis from Python
result = subprocess.run([
    './dns-resolver', 'resolve', 'example.com', 
    '--format', 'json'
], capture_output=True, text=True)

data = json.loads(result.stdout)
print(f"Found {len(data)} DNS records")
```

## Data Analysis Examples

### CSV Analysis with Excel/Google Sheets
```bash
# Generate data suitable for spreadsheet analysis
./dns-resolver bulk --input domains.txt --types A,MX --format csv --output analysis.csv
```

### JSON Processing with jq
```bash
# Extract only A records from JSON output
./dns-resolver resolve example.com --format json | jq '.[] | select(.record_type == "A")'

# Count records by type
./dns-resolver resolve example.com --format json | jq 'group_by(.record_type) | map({type: .[0].record_type, count: length})'
```

## Sample Domains for Testing

### Safe Testing Domains
- google.com (reliable, well-configured)
- cloudflare.com (modern DNS setup)
- github.com (developer-focused)
- stackoverflow.com (community site)

### Educational Domains
- example.com (RFC reserved)
- example.org (RFC reserved)
- example.net (RFC reserved)

### DNS Feature Testing
- ipv6.google.com (IPv6 testing)
- mx.google.com (mail exchange testing)
- ns.google.com (nameserver testing)

## Performance Benchmarks

### Typical Response Times
- Local DNS: 1-5ms
- Public DNS: 10-50ms
- International DNS: 50-200ms

### Bulk Analysis Rates
- Single domain: < 1 second
- 10 domains: 2-5 seconds
- 100 domains: 20-60 seconds

## Best Practices

### Rate Limiting
- Use reasonable concurrent limits (default: 10)
- Add delays between bulk operations
- Respect DNS server resources

### Data Management
- Save results with timestamps
- Use descriptive filenames
- Archive old results regularly

### Security Considerations
- Only test authorized domains
- Be aware of network monitoring
- Use VPN for sensitive research

## Troubleshooting Common Issues

### No Results Returned
```bash
# Increase timeout for slow networks
./dns-resolver resolve slow-site.com --timeout 30s

# Try different DNS servers
./dns-resolver resolve problematic-site.com --servers 1.1.1.1,8.8.8.8
```

### Permission Errors
```bash
# Ensure output directory exists and is writable
mkdir -p results
./dns-resolver resolve example.com --output results/test.json
```

### Network Connectivity Issues
```bash
# Test with verbose output
./dns-resolver resolve google.com --verbose

# Check specific DNS server connectivity
./dns-resolver resolve google.com --servers 8.8.8.8 --timeout 10s
```

---

Remember: This tool is for educational purposes. Always ensure you have permission to analyze the domains you're testing, and use the knowledge gained responsibly.