package client

import (
	"net"

	"github.com/miekg/dns"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("parseMDNSResponse", func() {
	var (
		entries  map[string]*types.MDNSDiscovery
		response *dns.Msg
	)

	BeforeEach(func() {
		entries = map[string]*types.MDNSDiscovery{}
		response = new(dns.Msg)
	})

	JustBeforeEach(func() {
		parseMDNSResponse(response, entries)
	})

	Context("with a complete PTR + SRV + TXT + A response", func() {
		BeforeEach(func() {
			response.Answer = []dns.RR{
				&dns.PTR{
					Hdr: dns.RR_Header{Name: freeboxMDNSService, Rrtype: dns.TypePTR},
					Ptr: "Freebox Server._fbx-api._tcp.local.",
				},
			}
			response.Extra = []dns.RR{
				&dns.SRV{
					Hdr:    dns.RR_Header{Name: "Freebox Server._fbx-api._tcp.local.", Rrtype: dns.TypeSRV},
					Target: "mafreebox.freebox.fr.",
					Port:   80,
				},
				&dns.TXT{
					Hdr: dns.RR_Header{Name: "Freebox Server._fbx-api._tcp.local.", Rrtype: dns.TypeTXT},
					Txt: []string{
						"uid=23b86ec8091013d668829fe12791fdab",
						"api_version=4.0",
						"device_type=FreeboxServer1,2",
						"api_base_url=/api/",
						"api_domain=abcdefgh.fbxos.fr",
						"https_available=1",
						"https_port=3615",
					},
				},
				&dns.A{
					Hdr: dns.RR_Header{Name: "mafreebox.freebox.fr.", Rrtype: dns.TypeA},
					A:   net.ParseIP("192.168.1.254"),
				},
			}
		})

		It("returns a fully populated MDNSDiscovery", func() {
			Expect(entries).To(HaveLen(1))
			e := entries["Freebox Server._fbx-api._tcp.local."]
			Expect(e).NotTo(BeNil())
			Expect(e.Name).To(Equal("Freebox Server"))
			Expect(e.Hostname).To(Equal("mafreebox.freebox.fr"))
			Expect(e.Port).To(Equal(uint16(80)))
			Expect(e.IPv4).To(HaveLen(1))
			Expect(e.IPv4[0].String()).To(Equal("192.168.1.254"))
			Expect(e.UID).To(Equal("23b86ec8091013d668829fe12791fdab"))
			Expect(e.APIVersion).To(Equal("4.0"))
			Expect(e.DeviceType).To(Equal("FreeboxServer1,2"))
			Expect(e.APIBaseURL).To(Equal("/api/"))
			Expect(e.APIDomain).To(Equal("abcdefgh.fbxos.fr"))
			Expect(e.HTTPSAvailable).To(BeTrue())
			Expect(e.HTTPSPort).To(Equal(3615))
		})
	})

	Context("with only a PTR record", func() {
		BeforeEach(func() {
			response.Answer = []dns.RR{
				&dns.PTR{
					Hdr: dns.RR_Header{Name: freeboxMDNSService, Rrtype: dns.TypePTR},
					Ptr: "Freebox Server._fbx-api._tcp.local.",
				},
			}
		})

		It("creates an entry with name but empty hostname and IPs", func() {
			Expect(entries).To(HaveLen(1))
			e := entries["Freebox Server._fbx-api._tcp.local."]
			Expect(e.Name).To(Equal("Freebox Server"))
			Expect(e.Hostname).To(BeEmpty())
			Expect(e.Port).To(BeZero())
			Expect(e.IPv4).To(BeEmpty())
			Expect(e.IPv6).To(BeEmpty())
		})

		Context("when SRV and A records arrive in a subsequent packet", func() {
			JustBeforeEach(func() {
				second := new(dns.Msg)
				second.Answer = []dns.RR{
					&dns.SRV{
						Hdr:    dns.RR_Header{Name: "Freebox Server._fbx-api._tcp.local.", Rrtype: dns.TypeSRV},
						Target: "mafreebox.freebox.fr.",
						Port:   80,
					},
				}
				second.Extra = []dns.RR{
					&dns.A{
						Hdr: dns.RR_Header{Name: "mafreebox.freebox.fr.", Rrtype: dns.TypeA},
						A:   net.ParseIP("192.168.1.254"),
					},
				}
				parseMDNSResponse(second, entries)
			})

			It("assembles the full entry across both packets", func() {
				e := entries["Freebox Server._fbx-api._tcp.local."]
				Expect(e.Hostname).To(Equal("mafreebox.freebox.fr"))
				Expect(e.Port).To(Equal(uint16(80)))
				Expect(e.IPv4).To(HaveLen(1))
				Expect(e.IPv4[0].String()).To(Equal("192.168.1.254"))
			})
		})
	})

	Context("with a PTR record for a different service type", func() {
		BeforeEach(func() {
			response.Answer = []dns.RR{
				&dns.PTR{
					Hdr: dns.RR_Header{Name: "_other._tcp.local.", Rrtype: dns.TypePTR},
					Ptr: "Some Device._other._tcp.local.",
				},
			}
		})

		It("ignores unrelated records", func() {
			Expect(entries).To(BeEmpty())
		})
	})

	Context("with an AAAA record", func() {
		BeforeEach(func() {
			entries["Freebox Server._fbx-api._tcp.local."] = &types.MDNSDiscovery{
				Name:     "Freebox Server",
				Hostname: "mafreebox.freebox.fr",
			}
			response.Extra = []dns.RR{
				&dns.AAAA{
					Hdr:  dns.RR_Header{Name: "mafreebox.freebox.fr.", Rrtype: dns.TypeAAAA},
					AAAA: net.ParseIP("2a01:cb00::1"),
				},
			}
		})

		It("populates IPv6", func() {
			e := entries["Freebox Server._fbx-api._tcp.local."]
			Expect(e.IPv6).To(HaveLen(1))
			Expect(e.IPv6[0].String()).To(Equal("2a01:cb00::1"))
		})
	})

	Context("when the same IP appears twice", func() {
		BeforeEach(func() {
			entries["Freebox Server._fbx-api._tcp.local."] = &types.MDNSDiscovery{
				Name:     "Freebox Server",
				Hostname: "mafreebox.freebox.fr",
				IPv4:     []net.IP{net.ParseIP("192.168.1.254")},
			}
			response.Extra = []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{Name: "mafreebox.freebox.fr.", Rrtype: dns.TypeA},
					A:   net.ParseIP("192.168.1.254"),
				},
			}
		})

		It("deduplicates IP addresses", func() {
			e := entries["Freebox Server._fbx-api._tcp.local."]
			Expect(e.IPv4).To(HaveLen(1))
		})
	})

	Context("with two Freeboxes in one response", func() {
		BeforeEach(func() {
			response.Answer = []dns.RR{
				&dns.PTR{
					Hdr: dns.RR_Header{Name: freeboxMDNSService, Rrtype: dns.TypePTR},
					Ptr: "Freebox One._fbx-api._tcp.local.",
				},
				&dns.PTR{
					Hdr: dns.RR_Header{Name: freeboxMDNSService, Rrtype: dns.TypePTR},
					Ptr: "Freebox Two._fbx-api._tcp.local.",
				},
			}
		})

		It("creates two separate entries", func() {
			Expect(entries).To(HaveLen(2))
			Expect(entries).To(HaveKey("Freebox One._fbx-api._tcp.local."))
			Expect(entries).To(HaveKey("Freebox Two._fbx-api._tcp.local."))
		})
	})
})
