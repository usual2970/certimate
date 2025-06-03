package udnr

type DomainDNSRecord struct {
	DnsType    string
	RecordName string
	Content    string
	TTL        int
	Prio       int
}
