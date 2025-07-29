# Advanced DNS Resolver

A comprehensive educational DNS analysis tool built in Go, featuring both CLI and web interfaces for advanced DNS resolution, bulk domain analysis, reverse DNS lookups, and DNS server performance testing.

## ⚠️ Educational Use Only

**WARNING**: This tool is designed for educational purposes and authorized testing only. Use only on domains you own or have explicit permission to test. Unauthorized DNS scanning may be illegal in your jurisdiction.

## Features

### Core Capabilities
- **Multiple DNS Record Types**: A, AAAA, CNAME, MX, NS, TXT, SOA, PTR, SRV
- **Bulk Domain Processing**: Concurrent analysis of multiple domains
- **Reverse DNS Lookups**: IP address to hostname resolution
- **Server Performance Testing**: Compare DNS server response times
- **Query Tracing**: Debug DNS resolution paths
- **Multiple Output Formats**: Text, JSON, CSV

### Interfaces
- **Command Line Interface**: Full-featured CLI with extensive options
- **Web Interface**: Modern web GUI with dark theme
- **REST API**: Backend API for integration

## Installation

### Prerequisites
- Go 1.24+ installed
- Git for cloning the repository

### Building from Source

```bash
# Clone the repository
git clone https://github.com/sammtan/dns-resolver.git
cd dns-resolver

# Download dependencies
go mod download

# Build CLI tool
go build -o dns-resolver cmd/main.go

# Build web server
cd web
go build -o dns-resolver-web main.go
```

## Usage

### Command Line Interface

#### Basic DNS Resolution
```bash
# Resolve all common record types for a domain
./dns-resolver resolve google.com

# Resolve specific record types
./dns-resolver resolve google.com --types A,MX,NS

# Output as JSON
./dns-resolver resolve google.com --format json

# Save results to file
./dns-resolver resolve google.com --format csv --output results.csv
```

#### Bulk Domain Analysis
```bash
# Analyze multiple domains
./dns-resolver bulk google.com github.com stackoverflow.com

# Read domains from file
./dns-resolver bulk --input domains.txt --types A,MX

# Export bulk results
./dns-resolver bulk google.com github.com --format json --output bulk-results.json
```

#### Reverse DNS Lookups
```bash
# Reverse lookup for IPv4
./dns-resolver reverse 8.8.8.8

# Reverse lookup for IPv6
./dns-resolver reverse 2001:4860:4860::8888

# JSON output
./dns-resolver reverse 8.8.8.8 --format json
```

#### DNS Server Performance Testing
```bash
# Test default DNS servers
./dns-resolver test

# Test specific servers
./dns-resolver test --servers 8.8.8.8,1.1.1.1,9.9.9.9

# More iterations for accuracy
./dns-resolver test --domain google.com --iterations 10

# Export performance data
./dns-resolver test --format csv --output performance.csv
```

#### DNS Query Tracing
```bash
# Trace DNS resolution path
./dns-resolver trace google.com

# Trace specific record type
./dns-resolver trace google.com --type AAAA

# Export trace results
./dns-resolver trace subdomain.example.com --format json --output trace.json
```

#### Advanced Options
```bash
# Custom DNS servers
./dns-resolver resolve google.com --servers 8.8.8.8,1.1.1.1

# Custom timeout and retries
./dns-resolver resolve google.com --timeout 10s --retries 5

# Concurrent queries
./dns-resolver bulk --input domains.txt --concurrent 20

# Verbose output
./dns-resolver resolve google.com --verbose
```

### Web Interface

#### Starting the Web Server
```bash
# Start web interface on port 5002
cd web
./dns-resolver-web
```

#### Accessing the Interface
Open your browser and navigate to: `http://localhost:5002`

#### Web Features
- **Interactive Forms**: Easy-to-use interface for all DNS operations
- **Real-time Results**: Immediate display of DNS analysis results
- **Export Functionality**: Download results as JSON or CSV
- **Dark Theme**: Professional dark theme matching portfolio design
- **Responsive Design**: Works on desktop and mobile devices

### REST API

#### Base URL
```
http://localhost:5002/api
```

#### Endpoints

##### DNS Resolution
```bash
POST /api/resolve
Content-Type: application/json

{
  "domain": "google.com",
  "record_types": ["A", "MX", "NS"],
  "servers": ["8.8.8.8", "1.1.1.1"],
  "timeout": 5,
  "concurrent": 10
}
```

##### Bulk Analysis
```bash
POST /api/bulk
Content-Type: application/json

{
  "domains": ["google.com", "github.com"],
  "record_types": ["A", "MX"],
  "timeout": 5,
  "concurrent": 10
}
```

##### Reverse DNS
```bash
POST /api/reverse
Content-Type: application/json

{
  "ip": "8.8.8.8",
  "timeout": 5
}
```

##### Performance Testing
```bash
POST /api/test
Content-Type: application/json

{
  "test_domain": "google.com",
  "iterations": 5,
  "servers": ["8.8.8.8", "1.1.1.1", "9.9.9.9"],
  "timeout": 5
}
```

##### Health Check
```bash
GET /api/health
```

## Configuration

### Default DNS Servers
- Google DNS: 8.8.8.8, 8.8.4.4
- Cloudflare DNS: 1.1.1.1, 1.0.0.1
- Quad9 DNS: 9.9.9.9, 149.112.112.112

### Default Settings
- **Timeout**: 5 seconds
- **Retries**: 3 attempts
- **Concurrent Queries**: 10
- **Web Server Port**: 5002

## Output Formats

### Text Format
Human-readable output with detailed information:
```
============================================================
                DNS RESOLUTION RESULTS
============================================================

Domain: google.com
Record Type: A
DNS Server: 8.8.8.8:53
Response Time: 15.2ms
TTL: 300 seconds
Records:
  142.250.4.100
  142.250.4.101
  142.250.4.102
```

### JSON Format
Structured data perfect for automation:
```json
[
  {
    "domain": "google.com",
    "record_type": "A",
    "records": ["142.250.4.100", "142.250.4.101"],
    "ttl": 300,
    "response_time_ms": 15200000,
    "dns_server": "8.8.8.8:53",
    "timestamp": "2025-07-29T12:00:00Z"
  }
]
```

### CSV Format
Spreadsheet-compatible format:
```csv
Domain,RecordType,Records,TTL,ResponseTime,Server,Error,Timestamp
google.com,A,142.250.4.100; 142.250.4.101,300,15.2ms,8.8.8.8:53,,2025-07-29T12:00:00Z
```

## Architecture

### Project Structure
```
dns-resolver/
├── cmd/                    # CLI application
│   └── main.go            # CLI entry point
├── pkg/                   # Core packages
│   └── resolver/          # DNS resolution engine
│       └── resolver.go    # Core resolver implementation
├── web/                   # Web interface
│   ├── main.go           # Web server
│   ├── templates/        # HTML templates
│   └── static/          # CSS, JS, assets
├── go.mod               # Go module definition
└── README.md            # Documentation
```

### Core Components

#### DNS Resolver Engine (`pkg/resolver/`)
- **Concurrent Processing**: Goroutines for parallel DNS queries
- **Multiple Protocols**: Support for various DNS record types
- **Error Handling**: Comprehensive error reporting and recovery
- **Performance Metrics**: Detailed timing and success rate tracking

#### CLI Interface (`cmd/`)
- **Cobra Framework**: Professional command-line interface
- **Multiple Commands**: Resolve, bulk, reverse, test, trace
- **Flexible Output**: Text, JSON, CSV formats
- **File I/O**: Input from files, output to files

#### Web Interface (`web/`)
- **Gin Framework**: High-performance HTTP server
- **REST API**: RESTful endpoints for all operations
- **Modern UI**: Dark theme with responsive design
- **Real-time Processing**: AJAX-based interactions

## Security Considerations

### Ethical Usage
- Only test domains you own or have permission to test
- Respect rate limits and server resources
- Use for educational and research purposes only

### Technical Security
- Input validation for all domain names and IP addresses
- Timeout controls to prevent resource exhaustion
- Concurrent request limiting
- No sensitive data logging

## Educational Value

This tool demonstrates:
- **DNS Protocol Understanding**: How DNS resolution works
- **Network Programming**: Socket programming and network protocols
- **Concurrent Programming**: Goroutines and channels in Go
- **Web Development**: REST APIs and modern web interfaces
- **Performance Analysis**: Measuring and comparing DNS server performance
- **Data Formats**: Working with JSON, CSV, and structured data

## Troubleshooting

### Common Issues

#### DNS Resolution Failures
- Check network connectivity
- Verify DNS server accessibility
- Try alternative DNS servers
- Check domain name spelling

#### Timeout Errors
- Increase timeout duration with `--timeout` flag
- Check firewall settings
- Verify DNS server responsiveness

#### Permission Errors
- Ensure proper file permissions for output files
- Check directory write permissions
- Run with appropriate user privileges

### Performance Optimization
- Adjust concurrent query limits with `--concurrent`
- Use local DNS servers for better performance
- Implement DNS caching for repeated queries

## 🔗 Combined Workflow Integration

### Real-World Security Analysis with Packet Sniffer

The DNS Resolver works powerfully in combination with network monitoring tools like the Packet Sniffer for comprehensive security analysis:

#### **Complete Infrastructure Assessment Example**

```bash
# Step 1: DNS Intelligence Gathering
./dns-resolver resolve github.com --types A,AAAA,MX,NS,TXT --format json
# Discovers: github.com → 20.205.243.166, api.github.com → 20.205.243.168

# Step 2: Real-time Traffic Monitoring (with Packet Sniffer)
# Monitor actual traffic to discovered IPs
# Result: Traffic actually goes to 140.82.113.21 (load balancer)

# Step 3: Reverse DNS Validation
./dns-resolver reverse 140.82.113.21
# Result: lb-140-82-113-21-iad.github.com (GitHub Washington DC)
```

#### **Intelligence Correlation Benefits**

**DNS Discovery + Traffic Analysis**:
- **DNS shows potential targets**: All possible IPs and infrastructure
- **Packet monitoring shows reality**: Actual communication endpoints
- **Reverse DNS validates**: Confirms infrastructure ownership

**Real Example Results**:
```bash
DNS Resolution:    github.com → 20.205.243.166
Actual Traffic:    192.168.1.5 ↔ 140.82.113.21 (12 packets, 100% HTTPS)
Reverse Lookup:    140.82.113.21 → lb-140-82-113-21-iad.github.com
Infrastructure:    GitHub load balancer in Washington DC datacenter
```

#### **Automated Integration Potential**

```go
// Pseudo-code for combined platform
targets := dnsResolver.BulkResolve(domains)
packetSniffer.MonitorIPs(targets.GetAllIPs())
for packet := range packetSniffer.Packets() {
    if !targets.Contains(packet.DestIP) {
        // New IP discovered through traffic analysis
        hostname := dnsResolver.Reverse(packet.DestIP)
        targets.AddDiscovered(packet.DestIP, hostname)
    }
}
```

This integration provides **complete network intelligence** that neither tool could achieve independently.

## Contributing

This is an educational project. When extending functionality:
1. Maintain the educational focus
2. Add comprehensive error handling
3. Include usage examples
4. Update documentation
5. Follow Go best practices

## License

MIT License - Educational use encouraged. See LICENSE file for details.

## Disclaimer

This tool is provided for educational and authorized testing purposes only. Users are responsible for ensuring compliance with applicable laws and regulations. The authors assume no responsibility for misuse of this tool.

---

**Version**: 1.0.0  
**Author**: Samuel Tan  
**Purpose**: Educational DNS Analysis Tool  
**Last Updated**: July 2025