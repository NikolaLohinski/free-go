package client

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("parseSSDPResponse", func() {
	var (
		data   []byte
		result *types.SSDPDiscovery
	)

	JustBeforeEach(func() {
		result = parseSSDPResponse(data)
	})

	Context("with a valid Freebox M-SEARCH response", func() {
		BeforeEach(func() {
			data = []byte("HTTP/1.1 200 OK\r\n" +
				"CACHE-CONTROL: max-age=1800\r\n" +
				"ST: urn:schemas-freebox-fr:device:Freebox:1\r\n" +
				"USN: uuid:abc::urn:schemas-freebox-fr:device:Freebox:1\r\n" +
				"LOCATION: http://192.168.1.254/index.html\r\n" +
				"SERVER: Linux/3.10 UPnP/1.0 Freebox/1.0\r\n" +
				"\r\n")
		})

		It("returns a populated SSDPDiscovery", func() {
			Expect(result).NotTo(BeNil())
			Expect(result.Location).To(Equal("http://192.168.1.254/index.html"))
			Expect(result.USN).To(Equal("uuid:abc::urn:schemas-freebox-fr:device:Freebox:1"))
			Expect(result.Server).To(Equal("Linux/3.10 UPnP/1.0 Freebox/1.0"))
		})
	})

	Context("with a response from a different device type", func() {
		BeforeEach(func() {
			data = []byte("HTTP/1.1 200 OK\r\n" +
				"ST: urn:schemas-upnp-org:device:Basic:1\r\n" +
				"LOCATION: http://192.168.1.1/desc.xml\r\n" +
				"\r\n")
		})

		It("returns nil", func() {
			Expect(result).To(BeNil())
		})
	})

	Context("with a non-200 status line", func() {
		BeforeEach(func() {
			data = []byte("HTTP/1.1 404 Not Found\r\n" +
				"ST: urn:schemas-freebox-fr:device:Freebox:1\r\n" +
				"\r\n")
		})

		It("returns nil", func() {
			Expect(result).To(BeNil())
		})
	})

	Context("with a Freebox response missing the LOCATION header", func() {
		BeforeEach(func() {
			data = []byte("HTTP/1.1 200 OK\r\n" +
				"ST: urn:schemas-freebox-fr:device:Freebox:1\r\n" +
				"USN: uuid:abc\r\n" +
				"\r\n")
		})

		It("returns nil", func() {
			Expect(result).To(BeNil())
		})
	})

	Context("with an empty payload", func() {
		BeforeEach(func() {
			data = []byte{}
		})

		It("returns nil", func() {
			Expect(result).To(BeNil())
		})
	})
})
