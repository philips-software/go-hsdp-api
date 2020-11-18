package pki_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCAs(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	muxPKI.HandleFunc("/core/pki/api/root/ca/pem", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pkix-cert")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `-----BEGIN CERTIFICATE-----
MIIHMzCCBRugAwIBAgIUB7awwVr04x/+xa1uFm7DAz85VdgwDQYJKoZIhvcNAQEL
BQAwgbMxFDASBgNVBAYTC05ldGhlcmxhbmRzMRYwFAYDVQQIEw1Ob29yZC1CcmFi
YW50MRIwEAYDVQQHEwlFaW5kaG92ZW4xKzApBgNVBAoTIlBoaWxpcHMgRWxlY3Ry
b25pY3MgTmVkZXJsYW5kIEIuVi4xFDASBgNVBAsTC0hlYWx0aFN1aXRlMSwwKgYD
VQQDEyNQaGlsaXBzIEhlYWx0aFN1aXRlIFByaXZhdGUgUm9vdCBDQTAeFw0yMDEx
MDYwOTE1NTVaFw0zMDExMDQwOTE2MTlaMIGzMRQwEgYDVQQGEwtOZXRoZXJsYW5k
czEWMBQGA1UECBMNTm9vcmQtQnJhYmFudDESMBAGA1UEBxMJRWluZGhvdmVuMSsw
KQYDVQQKEyJQaGlsaXBzIEVsZWN0cm9uaWNzIE5lZGVybGFuZCBCLlYuMRQwEgYD
VQQLEwtIZWFsdGhTdWl0ZTEsMCoGA1UEAxMjUGhpbGlwcyBIZWFsdGhTdWl0ZSBQ
cml2YXRlIFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDB
Dqi7GQ9oUMUU+6WMkVG+OeU90um2riOfCqICnDrXjYuGt6rjuyDO9Lk10OquVRaQ
f0gTAWirLfXbx0Ifrh0tPNB3XENTGjzf+K8zHxhHt2m18WoBWoCo8Bhc+v2UHqQy
CuKZhhZ9Wma4+kuowfmKJVJZD9zfGgkHRoqSSV+MyphggNukMnfjArSV0jHXOLc+
R8XMGJw9O++6kB1dOcxuj5Xmmv3bRyxRg1I9pWUBPovz400TZI1qz30jCj0TireS
sJyPD6SFH/4bSONEyAZ+n8U7m4JxwCrUlEnQ/zXSt7ZroKslYfAG/xRp1Jm5TqiJ
t1hJtzq7gRDCLkTGYRtaKoRUGpoZeES9GKVBlGKqNo5gCFbMpyulXf+IpTgQt9j2
scG6/l2qWhmPtdJ7atYzmL07/ooDk7SGV8fpetfRN0fdGw7Bn3NFk6wYsIfRZZKE
+68xizi6BZDhT54sHgZs1bTcYroAEWAijB6lMmBiK4kEDnp2ZEBSAwOh7ablepYf
Dbjht+Blfq7FzYeZUVCbILGl3CX2IVpXDTWzNjlo/hHGCaX9nh+5WaJVnNaihbwF
b86InqJ98oWfhl0PK+sh0ELJEt9dFQrSWxm6SSbEM4LT/Tr1A6dAqzz9Z65euXuA
sAHqZQGE25E/qQg9uWJxmiZ6bzsoh/zV3WkCUVasTwIDAQABo4IBOzCCATcwDgYD
VR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFP3me2KB210q
LZsoT8w5TMyHn5VrMB8GA1UdIwQYMBaAFP3me2KB210qLZsoT8w5TMyHn5VrMG4G
CCsGAQUFBwEBBGIwYDBeBggrBgEFBQcwAoZSaHR0cHM6Ly9wa2ktcHJveHktY2xp
ZW50LXRlc3QuZXUtd2VzdC5waGlsaXBzLWhlYWx0aHN1aXRlLmNvbS9jb3JlL3Br
aS9hcGkvcm9vdC9jYTBkBgNVHR8EXTBbMFmgV6BVhlNodHRwczovL3BraS1wcm94
eS1jbGllbnQtdGVzdC5ldS13ZXN0LnBoaWxpcHMtaGVhbHRoc3VpdGUuY29tL2Nv
cmUvcGtpL2FwaS9yb290L2NybDANBgkqhkiG9w0BAQsFAAOCAgEAU9/eE2dxC0Ph
7d/CXCFIlVqBEcqzFfJq8r9C2qQhAPM7cshDIq5m+wWod0zr3vA+wzbM7Wh/iu60
Ek8/p3E00CJPhiZeTcCejyYtJEqI1poVSnR/saQmuQ+Tahj6FkDRlo5a25NBBLKw
8Dx0g2NH6GEeHW4tW2yvV3c8/HcVMELGTcLuJW5cSe6W/Xt9WUyCJ852UDeGhvb4
RdxAzWDZrP48QRzc27YT7jg54MNPAKzf0mrNdmvtC5iOB19tSDK0FYTl1+Ab2/a5
gYAqnUSfVbKjudoaUaxU6vJKkcaTK1qOQhq0YelZZmoEPQSlrMtGIuExSy6UdMeo
ZxKBNQ9LXAu6r7elOv4pvcXk5v5q0S/uouAMxxxl4a5L5qyvJ3DPPVFHa3T4gJC0
mbToyqcYBi3HsqyryUXhpZmKlEX63U+tLmIsbxUJCqfCz97nNfgdUv6UUqD+INre
WSxEU8NbKQdbmM5WymUpTLZVx1JhPyl+DnAyWA9nfnlhU6IH+zpfkenphMsdmRSS
0TX2lxsQe5i1DXRw2o+XZxsJx2JWLIf74nuNourXdQs/taaibkTj6Y1dmrpMeDKx
GyUF/ncBCYaVXK+6DzS2kUCj9bPGVbGoXadaJxxGe9jGpOLcTR1xcQ/WSMAQuIZa
wo3KBVGxGCMPQZ8FeqGowJ0yDB8GxZ0=
-----END CERTIFICATE-----`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	cert, resp, err := pkiClient.Services.GetRootCA()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, cert) {
		return
	}
}
