// Package main is the command-line tool that does DNS lookups using
// dnsproxy/upstream.  See README.md for more information.
package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AdguardTeam/dnsproxy/upstream"
	"github.com/AdguardTeam/golibs/log"
	"github.com/miekg/dns"
)

// VersionString -- see the makefile
var VersionString = "master"

func main() {
	// parse env variables
	machineReadable := os.Getenv("JSON") == "1"
	shortTest := os.Getenv("TEST") == "1"
	insecureSkipVerify := os.Getenv("VERIFY") == "0"
	timeoutStr := os.Getenv("TIMEOUT")
	http3Enabled := os.Getenv("HTTP3") == "1"
	verbose := os.Getenv("VERBOSE") == "1"
	padding := os.Getenv("PAD") == "1"
	class := getClass()
	do := os.Getenv("DNSSEC") == "1"
	subnetOpt := getSubnet()
	ednsOpt := getEDNSOpt()
	rrType := getRRType()

	if verbose {
		log.SetLevel(log.DEBUG)
	}

	timeout := 10

	if !machineReadable {
		os.Stdout.WriteString(fmt.Sprintf("dnslookup %s\n", VersionString))

		if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
			os.Exit(0)
		}
	}

	if insecureSkipVerify {
		os.Stdout.WriteString("TLS verification has been disabled\n")
	}

	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		usage()
		os.Exit(0)
	}

	if len(os.Args) != 2 && len(os.Args) != 3 && len(os.Args) != 4 {
		log.Printf("Wrong number of arguments")
		usage()
		os.Exit(1)
	}

	if timeoutStr != "" {
		i, err := strconv.Atoi(timeoutStr)
		if err != nil {
			log.Printf("Wrong timeout value: %s", timeoutStr)
			usage()
			os.Exit(1)
		}

		timeout = i
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var n int
	n = 250
	rand.Seed(time.Now().UnixNano())
	iHost := Shuffle(NewSlice(1, n, 1))
	iDomain := Shuffle(NewSlice(0, n, 1))
	list := make([]string, 0, 1+(n*n))
	for _, d := range iDomain {
		for _, h := range iHost {
			url := "host" + strconv.Itoa(h) + ".domain" + strconv.Itoa(d) + ".rsx218-dox.cnam.fr\n"
			list = append(list, url)
		}
	}
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	now := time.Now().Unix()
	resultFileName := hostname + "_result_" + strconv.FormatInt(now, 10) + ".log"

	AppendToFile(resultFileName, "client,reqID,url,elapsedTime")
	var counter int = 1
	for _, url := range list {
		clean_url := CleanStr(url)

		domain := clean_url
		server := os.Args[1]

		var httpVersions []upstream.HTTPVersion
		if http3Enabled {
			httpVersions = []upstream.HTTPVersion{
				upstream.HTTPVersion3,
				upstream.HTTPVersion2,
				upstream.HTTPVersion11,
			}
		}

		opts := &upstream.Options{
			Timeout:            time.Duration(timeout) * time.Second,
			InsecureSkipVerify: insecureSkipVerify,
			HTTPVersions:       httpVersions,
		}

		if len(os.Args) == 4 {
			ip := net.ParseIP(os.Args[3])
			if ip == nil {
				log.Fatalf("invalid IP specified: %s", os.Args[3])
			}
			opts.ServerIPAddrs = []net.IP{ip}
		}

		u, err := upstream.AddressToUpstream(server, opts)
		if err != nil {
			log.Fatalf("Cannot create an upstream: %s", err)
		}

		req := &dns.Msg{}
		req.Id = dns.Id()
		req.RecursionDesired = true
		req.Question = []dns.Question{
			{Name: domain + ".", Qtype: rrType, Qclass: class},
		}

		if subnetOpt != nil {
			opt := getOrCreateOpt(req, do)
			opt.Option = append(opt.Option, subnetOpt)
		}

		if ednsOpt != nil {
			opt := getOrCreateOpt(req, do)
			opt.Option = append(opt.Option, ednsOpt)
		}

		if padding {
			opt := getOrCreateOpt(req, do)
			opt.Option = append(opt.Option, newEDNS0Padding(req))
		}

		startTime := time.Now()
		reply, err := u.Exchange(req)
		if err != nil {
			log.Fatalf("Cannot make the DNS request: %s", err)
		}

		if !machineReadable {
			str := fmt.Sprintf("%s,%d,%s,%s", hostname, counter, clean_url, time.Now().Sub(startTime))
			AppendToFile(resultFileName, str)
			fmt.Println(str)
		} else {
			var b []byte
			b, err = json.MarshalIndent(reply, "", "  ")
			if err != nil {
				log.Fatalf("Cannot marshal json: %s", err)
			}

			os.Stdout.WriteString(string(b) + "\n")
		}
		counter += 1
		if shortTest {
			time.Sleep(500 * time.Millisecond)
		}
		if counter == 5 && shortTest {
			os.Exit(0)
		}
	}
}

func getOrCreateOpt(req *dns.Msg, do bool) (opt *dns.OPT) {
	opt = req.IsEdns0()
	if opt == nil {
		req.SetEdns0(udpBufferSize, do)
		opt = req.IsEdns0()
	}

	return opt
}

func getEDNSOpt() (option *dns.EDNS0_LOCAL) {
	ednsOpt := os.Getenv("EDNSOPT")
	if ednsOpt == "" {
		return nil
	}

	parts := strings.Split(ednsOpt, ":")
	code, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("invalid EDNSOPT %s: %v", ednsOpt, err)
		usage()

		os.Exit(1)
	}

	var value []byte
	if len(parts) > 1 {
		value, err = hex.DecodeString(parts[1])
		if err != nil {
			log.Printf("invalid EDNSOPT %s: %v", ednsOpt, err)
			usage()

			os.Exit(1)
		}
	}

	return &dns.EDNS0_LOCAL{
		Code: uint16(code),
		Data: value,
	}
}

func getSubnet() (option *dns.EDNS0_SUBNET) {
	subnetStr := os.Getenv("SUBNET")
	if subnetStr == "" {
		return nil
	}

	_, ipNet, err := net.ParseCIDR(subnetStr)
	if err != nil {
		log.Printf("invalid SUBNET %s: %v", subnetStr, err)
		usage()

		os.Exit(1)
	}

	ones, _ := ipNet.Mask.Size()

	return &dns.EDNS0_SUBNET{
		Code:          dns.EDNS0SUBNET,
		Family:        1,
		SourceNetmask: uint8(ones),
		SourceScope:   0,
		Address:       ipNet.IP,
	}
}

func getClass() (class uint16) {
	classStr := os.Getenv("CLASS")
	var ok bool
	class, ok = dns.StringToClass[classStr]
	if !ok {
		if classStr != "" {
			log.Printf("Invalid CLASS: %q", classStr)
			usage()

			os.Exit(1)
		}

		class = dns.ClassINET
	}
	return class
}

func getRRType() (rrType uint16) {
	rrTypeStr := os.Getenv("RRTYPE")
	var ok bool
	rrType, ok = dns.StringToType[rrTypeStr]
	if !ok {
		if rrTypeStr != "" {
			log.Printf("Invalid RRTYPE: %q", rrTypeStr)
			usage()

			os.Exit(1)
		}

		rrType = dns.TypeA
	}
	return rrType
}

func usage() {
	os.Stdout.WriteString("Usage: dnslookup <domain> <server> [<providerName> <serverPk>]\n")
	os.Stdout.WriteString("<domain>: mandatory, domain name to lookup\n")
	os.Stdout.WriteString("<server>: mandatory, server address. Supported: plain, tls:// (DOT), https:// (DOH), sdns:// (DNSCrypt), quic:// (DOQ)\n")
	os.Stdout.WriteString("<providerName>: optional, DNSCrypt provider name\n")
	os.Stdout.WriteString("<serverPk>: optional, DNSCrypt server public key\n")
}

// requestPaddingBlockSize is used to pad responses over DoT and DoH according
// to RFC 8467.
const requestPaddingBlockSize = 128
const udpBufferSize = dns.DefaultMsgSize

// newEDNS0Padding constructs a new OPT RR EDNS0 Padding for the extra section.
func newEDNS0Padding(req *dns.Msg) (option *dns.EDNS0_PADDING) {
	msgLen := req.Len()
	padLen := requestPaddingBlockSize - msgLen%requestPaddingBlockSize

	// Truncate padding to fit in UDP buffer.
	if msgLen+padLen > udpBufferSize {
		padLen = udpBufferSize - msgLen
		if padLen < 0 {
			padLen = 0
		}
	}

	return &dns.EDNS0_PADDING{Padding: make([]byte, padLen)}
}

func NewSlice(start, end, step int) []int {
	if step <= 0 || end < start {
		return []int{}
	}
	s := make([]int, 0, 1+(end-start)/step)
	for start <= end {
		s = append(s, start)
		start += step
	}
	return s
}

func Shuffle(vals []int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func GetTimeMs() string {
	return time.Now().Format(time.StampMilli)
}

func CleanStr(str string) string {
	str = strings.ReplaceAll(str, "\n", "")
	return str
}

func AppendToFile(resultFileName string, str string) {
	f, err := os.OpenFile(resultFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(str + "\n"); err != nil {
		log.Println(err)
	}
}
