// Advanced DNS Resolver - Web Interface
// Educational Security Tool with Web GUI
package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sammtan/dns-resolver/pkg/resolver"
)

type QueryRequest struct {
	Domain      string   `json:"domain" binding:"required"`
	RecordTypes []string `json:"record_types"`
	Servers     []string `json:"servers"`
	Timeout     int      `json:"timeout"`
	Concurrent  int      `json:"concurrent"`
}

type BulkQueryRequest struct {
	Domains     []string `json:"domains" binding:"required"`
	RecordTypes []string `json:"record_types"`
	Servers     []string `json:"servers"`
	Timeout     int      `json:"timeout"`
	Concurrent  int      `json:"concurrent"`
}

type ReverseQueryRequest struct {
	IP         string   `json:"ip" binding:"required"`
	Servers    []string `json:"servers"`
	Timeout    int      `json:"timeout"`
	Concurrent int      `json:"concurrent"`
}

type TestRequest struct {
	TestDomain string   `json:"test_domain"`
	Iterations int      `json:"iterations"`
	Servers    []string `json:"servers"`
	Timeout    int      `json:"timeout"`
}

func main() {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.Default()
	
	// Load HTML templates
	r.LoadHTMLGlob("templates/*")
	
	// Serve static files
	r.Static("/static", "./static")
	
	// Routes
	r.GET("/", homePage)
	r.POST("/api/resolve", resolveHandler)
	r.POST("/api/bulk", bulkHandler)
	r.POST("/api/reverse", reverseHandler)
	r.POST("/api/test", testHandler)
	r.GET("/api/health", healthHandler)
	
	// CORS middleware for API endpoints
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// Start server
	println(`
==============================================================
               DNS RESOLVER WEB INTERFACE                    
                     Starting Server...                      
                                                              
  WARNING: Use only on domains you own or have permission   
          to test. Unauthorized scanning may be illegal!    
==============================================================
`)
	
	r.Run(":5002") // Run on port 5002 to avoid conflicts
}

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Advanced DNS Resolver",
	})
}

func resolveHandler(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set defaults
	if len(req.RecordTypes) == 0 {
		req.RecordTypes = []string{"A", "AAAA", "CNAME", "MX", "NS", "TXT"}
	}
	if len(req.Servers) == 0 {
		req.Servers = []string{"8.8.8.8", "1.1.1.1", "9.9.9.9"}
	}
	if req.Timeout == 0 {
		req.Timeout = 5
	}
	if req.Concurrent == 0 {
		req.Concurrent = 10
	}
	
	// Convert record types
	var recordTypes []resolver.RecordType
	for _, rt := range req.RecordTypes {
		recordTypes = append(recordTypes, resolver.RecordType(strings.ToUpper(rt)))
	}
	
	// Create resolver
	r := resolver.NewResolver(req.Servers, time.Duration(req.Timeout)*time.Second, 3, req.Concurrent)
	
	// Perform resolution
	results, err := r.ResolveAll(req.Domain, recordTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"domain":  req.Domain,
		"results": results,
		"count":   len(results),
	})
}

func bulkHandler(c *gin.Context) {
	var req BulkQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set defaults
	if len(req.RecordTypes) == 0 {
		req.RecordTypes = []string{"A", "AAAA", "MX"}
	}
	if len(req.Servers) == 0 {
		req.Servers = []string{"8.8.8.8", "1.1.1.1", "9.9.9.9"}
	}
	if req.Timeout == 0 {
		req.Timeout = 5
	}
	if req.Concurrent == 0 {
		req.Concurrent = 10
	}
	
	// Limit bulk queries for web interface
	if len(req.Domains) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 50 domains allowed for bulk queries"})
		return
	}
	
	// Convert record types
	var recordTypes []resolver.RecordType
	for _, rt := range req.RecordTypes {
		recordTypes = append(recordTypes, resolver.RecordType(strings.ToUpper(rt)))
	}
	
	// Create resolver
	r := resolver.NewResolver(req.Servers, time.Duration(req.Timeout)*time.Second, 3, req.Concurrent)
	
	// Perform bulk resolution
	results, err := r.BulkResolve(req.Domains, recordTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"domains": req.Domains,
		"results": results,
		"count":   len(results),
	})
}

func reverseHandler(c *gin.Context) {
	var req ReverseQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set defaults
	if len(req.Servers) == 0 {
		req.Servers = []string{"8.8.8.8", "1.1.1.1", "9.9.9.9"}
	}
	if req.Timeout == 0 {
		req.Timeout = 5
	}
	if req.Concurrent == 0 {
		req.Concurrent = 10
	}
	
	// Create resolver
	r := resolver.NewResolver(req.Servers, time.Duration(req.Timeout)*time.Second, 3, req.Concurrent)
	
	// Perform reverse DNS lookup
	result, err := r.ReverseDNS(req.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"ip":     req.IP,
		"result": result,
	})
}

func testHandler(c *gin.Context) {
	var req TestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set defaults
	if req.TestDomain == "" {
		req.TestDomain = "google.com"
	}
	if req.Iterations == 0 {
		req.Iterations = 5
	}
	if len(req.Servers) == 0 {
		req.Servers = []string{"8.8.8.8", "1.1.1.1", "9.9.9.9"}
	}
	if req.Timeout == 0 {
		req.Timeout = 5
	}
	
	// Create resolver
	r := resolver.NewResolver(req.Servers, time.Duration(req.Timeout)*time.Second, 3, 5)
	
	// Test server performance
	results, err := r.TestServers(req.TestDomain, req.Iterations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"test_domain": req.TestDomain,
		"iterations":  req.Iterations,
		"results":     results,
		"count":       len(results),
	})
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "dns-resolver",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}