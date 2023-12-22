package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// defining two structures DomainURL and DomainVar

// DomainURL will decode the input

// DomainVar will encode the input

type DomainURL struct {
    DomainURL string `string:"domainurl"`
}

type DomainVar struct {

    // the atrributes will be 
    // mapped as shown in
    // the JSON object

    Domain string `json:"domain"`
    HasMX bool `json:"hasmx"`
    HasSPF bool `json:haspf`
    SpfRecord string `json:"spfrecord"`
    HasDMARC bool `json:"hasdmarc"`
    DmarcRecord string `json:"dmarcRecord"`
}

// definition domainVars slice
var domainVars []DomainVar

func isValidDomain (domain string) {

	var hasMX, hasSPF, hasDMARC bool

    // SPF --> Sender Policy Framework
	var spfRecord string

    // DMARC --> Domain-based Message Authentication
    //           Reporting and Conformance

	var dmarcRecord string

    // Look for Mail Exchange records
    // usign the net/http import
    mxRecord, err := net.LookupMX(domain)

    if err != nil {
        log.Printf("Error: %v\n", err)
    }

    if len(mxRecord) > 0 {
        hasMX = true
    }

    txtRecords, err := net.LookupTXT(domain)

    if err != nil {
        log.Printf("Error: %v\n", err)
    }

    for _, record := range txtRecords {
        if strings.HasPrefix(record, "v=spf1") {
            hasSPF = true
            spfRecord = record
            break
        }
    }

    dmarcRecords, err := net.LookupTXT("_dmarc." + domain)

    if err != nil {

        log.Printf("Error: %v\n", err)
    }

    for _, record := range dmarcRecords {

        if strings.HasPrefix(records, "v=DMARC1") {
            hasDMARC = true
            dmarcRecord = record
            break
        }
        
    }

    // printing our boolean values

    fmt.Printf("domain=%v\n,hasMX=%v\n,hasSPF=%v\n,
                spfRecord=%v\n, hasDMARC=%v\n,
                dmarcRecord=%v\n",domain,hasMX,hasSPF,
                spfRecord,hasDMARC,dmarcRecord)



    // Assigning into each attribute 
    // of the domainVar the email data
    // collected

    var domainVar DomainVar
    domainVar.Domain      = domain
    domainVar.HasMX       = hasMX
    domainVar.HasSPF      = hasSPF
    domainVar.SpfRecord   = spfRecord 
    domainVar.HasDMARC    = hasDMARC
    domainVar.DmarcRecord = dmarcRecord

    return domainVar
}


// Handler function

// w --> interface for the server's response
//       to the client.

// r --> pointer to the client's request to the server
//

func formHandler(w http.ResponseWriter, r *http.Request) {

    // setting the content of the response to
    // an application/json. 

    // With the Header() method we acces the header
    // map that will be sent by WriteHeader

               // setting the header key 
    w.Header().set("Content-Type", "application/json")


    // domainURL will hold data
    // decoded from the JSON

    var domainURL DomainURL
    json.NewDecoder(r.Body).Decode(&domainURL)

    // Calling isValidDomain function
    // and passing the attribute DomainURL 
    domainVar  := isValidDomain(domainURL.DomainURL)

    // adding into domainVars the new value in
    // domainVar
    domainVars =  append(domainVars, domainVar)

    json.NewEncoder(w).Encode(domain.Vars)
}

func main () {

    r:= mux.NewRouter()
    r.HandleFunc("/from", formHandler).Methods("POST")
   
    fmt.Print("Starting server at port 8000\n")
    
    // if the ListenAnServe function returns an error,
    // log.Fatal will log an error an stop the program.
    log.Fatal(http.ListenAndServe(":8000",r))
}
