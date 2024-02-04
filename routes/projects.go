package routes

import (
	"log"
	"net"
	"strings"
	"github.com/gin-gonic/gin"
)

type DomainCheckResult struct {
	Domain       string
	HasMX        bool
	HasSPF       bool
	SPFRecord    string
	HasDMARC     bool
	DMARCRecord  string
}

func checkDomain(domain string) DomainCheckResult {
	var result DomainCheckResult

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	result.Domain = domain
	if len(mxRecords) > 0 {
		result.HasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			result.HasSPF = true
			result.SPFRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("ErrorL%v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			result.HasDMARC = true
			result.DMARCRecord = record
			break
		}
	}

	return result
}

func VerifyEmail(c *gin.Context) {
	domain := c.Query("domain")
	log.Printf("Domain: %v\n", domain) 

	result := checkDomain(domain)

	c.JSON(200, result)
}
