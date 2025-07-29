/**
 * Advanced DNS Resolver Web Interface JavaScript
 * ==============================================
 * 
 * Interactive frontend for the DNS Resolver tool providing
 * real-time DNS analysis with multiple query types.
 * 
 * Author: Samuel Tan
 * License: MIT
 * Version: 1.0.0
 */

class DNSResolverApp {
    constructor() {
        this.currentResults = null;
        this.apiBase = '';
        
        this.initializeEventListeners();
        this.initializeTabs();
    }
    
    initializeEventListeners() {
        // Form submissions
        document.getElementById('resolveForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.performDNSLookup();
        });
        
        document.getElementById('bulkForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.performBulkAnalysis();
        });
        
        document.getElementById('reverseForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.performReverseDNS();
        });
        
        document.getElementById('testForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.performServerTest();
        });
    }
    
    initializeTabs() {
        const tabs = document.querySelectorAll('.nav-tab');
        const contents = document.querySelectorAll('.tab-content');
        
        tabs.forEach(tab => {
            tab.addEventListener('click', () => {
                // Remove active class from all tabs and contents
                tabs.forEach(t => t.classList.remove('active'));
                contents.forEach(c => c.classList.remove('active'));
                
                // Add active class to clicked tab
                tab.classList.add('active');
                
                // Show corresponding content
                const targetTab = tab.getAttribute('data-tab');
                document.getElementById(`${targetTab}-tab`).classList.add('active');
            });
        });
    }
    
    async performDNSLookup() {
        const formData = new FormData(document.getElementById('resolveForm'));
        const domain = formData.get('domain').trim();
        
        if (!domain) {
            this.showError('Please enter a domain name');
            return;
        }
        
        const recordTypes = [];
        formData.getAll('record_types').forEach(type => {
            recordTypes.push(type);
        });
        
        if (recordTypes.length === 0) {
            this.showError('Please select at least one record type');
            return;
        }
        
        const requestData = {
            domain: domain,
            record_types: recordTypes,
            timeout: parseInt(formData.get('timeout')) || 5,
            concurrent: parseInt(formData.get('concurrent')) || 10
        };
        
        // Parse servers if provided
        const servers = formData.get('servers').trim();
        if (servers) {
            requestData.servers = servers.split(',').map(s => s.trim());
        }
        
        this.showLoading('Analyzing DNS records...');
        
        try {
            const response = await fetch('/api/resolve', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestData)
            });
            
            const result = await response.json();
            
            if (response.ok) {
                this.currentResults = result;
                this.displayDNSResults(result);
            } else {
                this.showError(result.error || 'DNS lookup failed');
            }
            
        } catch (error) {
            this.showError('Network error: ' + error.message);
        } finally {
            this.hideLoading();
        }
    }
    
    async performBulkAnalysis() {
        const formData = new FormData(document.getElementById('bulkForm'));
        const domainsText = formData.get('domains').trim();
        
        if (!domainsText) {
            this.showError('Please enter domain names');
            return;
        }
        
        const domains = domainsText.split('\n').map(d => d.trim()).filter(d => d);
        
        if (domains.length === 0) {
            this.showError('Please enter valid domain names');
            return;
        }
        
        if (domains.length > 50) {
            this.showError('Maximum 50 domains allowed');
            return;
        }
        
        const recordTypes = [];
        formData.getAll('bulk_record_types').forEach(type => {
            recordTypes.push(type);
        });
        
        if (recordTypes.length === 0) {
            this.showError('Please select at least one record type');
            return;
        }
        
        const requestData = {
            domains: domains,
            record_types: recordTypes,
            timeout: 5,
            concurrent: 10
        };
        
        this.showLoading(`Analyzing ${domains.length} domains...`);
        
        try {
            const response = await fetch('/api/bulk', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestData)
            });
            
            const result = await response.json();
            
            if (response.ok) {
                this.currentResults = result;
                this.displayBulkResults(result);
            } else {
                this.showError(result.error || 'Bulk analysis failed');
            }
            
        } catch (error) {
            this.showError('Network error: ' + error.message);
        } finally {
            this.hideLoading();
        }
    }
    
    async performReverseDNS() {
        const formData = new FormData(document.getElementById('reverseForm'));
        const ip = formData.get('ip').trim();
        
        if (!ip) {
            this.showError('Please enter an IP address');
            return;
        }
        
        const requestData = {
            ip: ip,
            timeout: 5,
            concurrent: 10
        };
        
        this.showLoading('Performing reverse DNS lookup...');
        
        try {
            const response = await fetch('/api/reverse', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestData)
            });
            
            const result = await response.json();
            
            if (response.ok) {
                this.currentResults = result;
                this.displayReverseResults(result);
            } else {
                this.showError(result.error || 'Reverse DNS lookup failed');
            }
            
        } catch (error) {
            this.showError('Network error: ' + error.message);
        } finally {
            this.hideLoading();
        }
    }
    
    async performServerTest() {
        const formData = new FormData(document.getElementById('testForm'));
        const testDomain = formData.get('test_domain').trim() || 'google.com';
        const iterations = parseInt(formData.get('iterations')) || 5;
        const servers = formData.get('servers').trim();
        
        const requestData = {
            test_domain: testDomain,
            iterations: iterations,
            timeout: 5
        };
        
        if (servers) {
            requestData.servers = servers.split(',').map(s => s.trim());
        }
        
        this.showLoading('Testing DNS server performance...');
        
        try {
            const response = await fetch('/api/test', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(requestData)
            });
            
            const result = await response.json();
            
            if (response.ok) {
                this.currentResults = result;
                this.displayTestResults(result);
            } else {
                this.showError(result.error || 'Server test failed');
            }
            
        } catch (error) {
            this.showError('Network error: ' + error.message);
        } finally {
            this.hideLoading();
        }
    }
    
    displayDNSResults(result) {
        const container = document.getElementById('resolveResultsContainer');
        const panel = document.getElementById('resolveResults');
        
        document.getElementById('resultDomain').textContent = result.domain;
        document.getElementById('recordCount').textContent = result.count;
        
        container.innerHTML = '';
        
        if (result.results && result.results.length > 0) {
            result.results.forEach(record => {
                const card = this.createRecordCard(record);
                container.appendChild(card);
            });
        } else {
            container.innerHTML = '<div class="error-message">No DNS records found</div>';
        }
        
        panel.classList.remove('hidden');
        panel.scrollIntoView({ behavior: 'smooth' });
    }
    
    displayBulkResults(result) {
        const container = document.getElementById('bulkResultsContainer');
        const panel = document.getElementById('bulkResults');
        
        document.getElementById('bulkDomainCount').textContent = result.count;
        
        container.innerHTML = '';
        
        if (result.results && result.results.length > 0) {
            result.results.forEach(bulk => {
                const domainSection = document.createElement('div');
                domainSection.className = 'domain-section';
                
                const domainHeader = document.createElement('h3');
                domainHeader.textContent = bulk.domain;
                domainHeader.style.color = 'var(--accent-primary)';
                domainHeader.style.marginBottom = '15px';
                domainSection.appendChild(domainHeader);
                
                if (bulk.error) {
                    const errorDiv = document.createElement('div');
                    errorDiv.className = 'error-message';
                    errorDiv.textContent = bulk.error;
                    domainSection.appendChild(errorDiv);
                } else if (bulk.results) {
                    bulk.results.forEach(record => {
                        const card = this.createRecordCard(record);
                        domainSection.appendChild(card);
                    });
                }
                
                container.appendChild(domainSection);
            });
        } else {
            container.innerHTML = '<div class="error-message">No results found</div>';
        }
        
        panel.classList.remove('hidden');
        panel.scrollIntoView({ behavior: 'smooth' });
    }
    
    displayReverseResults(result) {
        const container = document.getElementById('reverseResultsContainer');
        const panel = document.getElementById('reverseResults');
        
        container.innerHTML = '';
        
        if (result.result) {
            const card = this.createRecordCard(result.result);
            container.appendChild(card);
        } else {
            container.innerHTML = '<div class="error-message">Reverse DNS lookup failed</div>';
        }
        
        panel.classList.remove('hidden');
        panel.scrollIntoView({ behavior: 'smooth' });
    }
    
    displayTestResults(result) {
        const container = document.getElementById('testResultsContainer');
        const panel = document.getElementById('testResults');
        
        container.innerHTML = '';
        
        if (result.results && result.results.length > 0) {
            const table = document.createElement('table');
            table.className = 'performance-table';
            
            // Create header
            const thead = document.createElement('thead');
            const headerRow = document.createElement('tr');
            ['Server', 'Avg Response', 'Min Response', 'Max Response', 'Success Rate', 'Queries'].forEach(text => {
                const th = document.createElement('th');
                th.textContent = text;
                headerRow.appendChild(th);
            });
            thead.appendChild(headerRow);
            table.appendChild(thead);
            
            // Create body
            const tbody = document.createElement('tbody');
            result.results.forEach(perf => {
                const row = document.createElement('tr');
                
                const serverCell = document.createElement('td');
                serverCell.textContent = perf.server;
                row.appendChild(serverCell);
                
                const avgCell = document.createElement('td');
                avgCell.textContent = this.formatDuration(perf.avg_response_time_ms);
                row.appendChild(avgCell);
                
                const minCell = document.createElement('td');
                minCell.textContent = this.formatDuration(perf.min_response_time_ms);
                row.appendChild(minCell);
                
                const maxCell = document.createElement('td');
                maxCell.textContent = this.formatDuration(perf.max_response_time_ms);
                row.appendChild(maxCell);
                
                const successCell = document.createElement('td');
                successCell.textContent = perf.success_rate.toFixed(1) + '%';
                row.appendChild(successCell);
                
                const queriesCell = document.createElement('td');
                queriesCell.textContent = perf.total_queries;
                row.appendChild(queriesCell);
                
                tbody.appendChild(row);
            });
            table.appendChild(tbody);
            
            container.appendChild(table);
        } else {
            container.innerHTML = '<div class="error-message">No performance data available</div>';
        }
        
        panel.classList.remove('hidden');
        panel.scrollIntoView({ behavior: 'smooth' });
    }
    
    createRecordCard(record) {
        const card = document.createElement('div');
        card.className = 'record-card';
        
        const header = document.createElement('div');
        header.className = 'record-header';
        
        const typeSpan = document.createElement('span');
        typeSpan.className = 'record-type';
        typeSpan.textContent = record.record_type;
        header.appendChild(typeSpan);
        
        const meta = document.createElement('div');
        meta.className = 'record-meta';
        meta.innerHTML = `TTL: ${record.ttl}s | Response: ${this.formatDuration(record.response_time_ms)} | Server: ${record.dns_server}`;
        header.appendChild(meta);
        
        card.appendChild(header);
        
        if (record.error) {
            const errorDiv = document.createElement('div');
            errorDiv.className = 'error-message';
            errorDiv.textContent = record.error;
            card.appendChild(errorDiv);
        } else if (record.records && record.records.length > 0) {
            const valuesDiv = document.createElement('div');
            valuesDiv.className = 'record-values';
            
            record.records.forEach(value => {
                const valueDiv = document.createElement('div');
                valueDiv.className = 'record-value';
                valueDiv.textContent = value;
                valuesDiv.appendChild(valueDiv);
            });
            
            card.appendChild(valuesDiv);
        } else {
            const noDataDiv = document.createElement('div');
            noDataDiv.textContent = 'No records found';
            noDataDiv.style.color = 'var(--text-muted)';
            noDataDiv.style.fontStyle = 'italic';
            card.appendChild(noDataDiv);
        }
        
        return card;
    }
    
    formatDuration(nanoseconds) {
        if (!nanoseconds) return '0ms';
        const ms = nanoseconds / 1000000;
        if (ms < 1) {
            return (nanoseconds / 1000).toFixed(0) + 'μs';
        }
        return ms.toFixed(1) + 'ms';
    }
    
    showLoading(message) {
        const overlay = document.getElementById('loadingOverlay');
        if (message) {
            overlay.querySelector('p').textContent = message;
        }
        overlay.classList.remove('hidden');
    }
    
    hideLoading() {
        document.getElementById('loadingOverlay').classList.add('hidden');
    }
    
    showError(message) {
        alert('Error: ' + message);
    }
}

// Export functionality
window.exportResults = function(format, type) {
    const app = window.dnsResolverApp;
    if (!app.currentResults) {
        alert('No results to export');
        return;
    }
    
    let data;
    let filename;
    let mimeType;
    
    if (format === 'json') {
        data = JSON.stringify(app.currentResults, null, 2);
        filename = `dns_${type}_results.json`;
        mimeType = 'application/json';
    } else if (format === 'csv') {
        data = app.convertToCSV(app.currentResults, type);
        filename = `dns_${type}_results.csv`;
        mimeType = 'text/csv';
    }
    
    if (data) {
        app.downloadFile(data, filename, mimeType);
    }
};

DNSResolverApp.prototype.convertToCSV = function(results, type) {
    const csvRows = [];
    
    if (type === 'resolve' && results.results) {
        csvRows.push(['Domain', 'RecordType', 'Records', 'TTL', 'ResponseTime', 'DNSServer', 'Error']);
        
        results.results.forEach(record => {
            csvRows.push([
                record.domain,
                record.record_type,
                record.records ? record.records.join('; ') : '',
                record.ttl,
                record.response_time_ms,
                record.dns_server,
                record.error || ''
            ]);
        });
    } else if (type === 'bulk' && results.results) {
        csvRows.push(['Domain', 'RecordType', 'Records', 'TTL', 'ResponseTime', 'DNSServer', 'Error']);
        
        results.results.forEach(bulk => {
            if (bulk.results) {
                bulk.results.forEach(record => {
                    csvRows.push([
                        bulk.domain,
                        record.record_type,
                        record.records ? record.records.join('; ') : '',
                        record.ttl,
                        record.response_time_ms,
                        record.dns_server,
                        record.error || bulk.error || ''
                    ]);
                });
            }
        });
    }
    
    return csvRows.map(row => 
        row.map(cell => 
            typeof cell === 'string' && cell.includes(',') ? `"${cell}"` : cell
        ).join(',')
    ).join('\n');
};

DNSResolverApp.prototype.downloadFile = function(data, filename, mimeType) {
    const blob = new Blob([data], { type: mimeType });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
};

// Initialize the application when the page loads
document.addEventListener('DOMContentLoaded', () => {
    window.dnsResolverApp = new DNSResolverApp();
});