package pki_test

import (
	"crypto/ecdsa"
	"crypto/x509"
	"io"
	"net/http"
	"testing"

	"github.com/philips-software/go-hsdp-api/pki"

	"github.com/stretchr/testify/assert"
)

func TestGetCAs(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	certificate := `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`

	returnCA := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pkix-cert")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, certificate)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	muxPKI.HandleFunc("/core/pki/api/root/ca/pem", returnCA)
	muxPKI.HandleFunc("/core/pki/api/policy/ca/pem", returnCA)

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
	cert, resp, err = pkiClient.Services.GetPolicyCA()
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

func TestGetCRLs(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	getCrl := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pkix-crl")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `-----BEGIN X509 CRL-----
MIIDJTCCAQ0CAQEwDQYJKoZIhvcNAQELBQAwgbMxFDASBgNVBAYTC05ldGhlcmxh
bmRzMRYwFAYDVQQIEw1Ob29yZC1CcmFiYW50MRIwEAYDVQQHEwlFaW5kaG92ZW4x
KzApBgNVBAoTIlBoaWxpcHMgRWxlY3Ryb25pY3MgTmVkZXJsYW5kIEIuVi4xFDAS
BgNVBAsTC0hlYWx0aFN1aXRlMSwwKgYDVQQDEyNQaGlsaXBzIEhlYWx0aFN1aXRl
IFByaXZhdGUgUm9vdCBDQRcNMjAxMTA2MDkxNjI1WhcNMzAxMTA0MDkxNjI1WjAA
oCMwITAfBgNVHSMEGDAWgBT95ntigdtdKi2bKE/MOUzMh5+VazANBgkqhkiG9w0B
AQsFAAOCAgEAfhyfD3l3pjEqp8ALA+d9D/mO8ZaEZXrt4kRLYpkCfTKgdIFFlNQP
kUgkMGwvBpYiRkO7TZkM8doEEnVzcsGfsYDmw37knM7+YIIEtqFp5n/nBUcx2z4H
lGy17K84t93yRDu50qfMs5OIJ/1zWX/3HBthuM0CfsoB2nN8TEyKDwYp7yTlxxaf
RngFC6Hkxn6cRwcOLOv64rn6ej+xoH+A4C3e61HQ8SZdbDvNrHEyxaUj51ztl1RP
axGEAx+hAQ98GZ7FCCrkuj45nfayFTl+5B+wgd8rN5Dx4PzyEShfARNqUwmXkuHn
ZDrbVXb1PlW2ngzpWca1L0VIP5TD3nJak/7je+HzhKKITBuALKiMF0muelmWH8wt
2ElsIB0wpTw02o+jrDeEvoHFFg8do6yWGp9vFvvt9gz9GgvI7gVXec/NwTuphh+z
SeINakInr13EWHeOhoTaf6hxyysS9g5hx2lV9X5wH0K1xoL5wqK8DCv1SL1+Vz1X
HBqISrpKF7YNyXbmjIUhWqglS1fJvgdSJHlH9utxhnd1K4kyQT4VjPtZfU9dM3GJ
yK0P70Ty2WxPuug1wgSccAdOCArYC5te7KOVTQJy7lG/YHPZLkN5hDO6rmtTQA+O
kljJ1cnVriYSyGoStCTCep8b4zDjl3KTdu2cGU4tUZIif6E2DruBZJ8=
-----END X509 CRL-----`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

	muxPKI.HandleFunc("/core/pki/api/root/crl/pem", getCrl)
	muxPKI.HandleFunc("/core/pki/api/policy/crl/pem", getCrl)

	crl, resp, err := pkiClient.Services.GetRootCRL()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, crl) {
		return
	}
	crl, resp, err = pkiClient.Services.GetPolicyCRL()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, crl) {
		return
	}
}

func TestIssueAndSignCertificates(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// NOTE: The generated certificate, including the Private Key are throwaway value
	response := `{
  "request_id": "8f74aacf-09f0-b6d9-31e0-1c3e7a9d7c22",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "ca_chain": [
      "-----BEGIN CERTIFICATE-----\nMIIFZDCCA0ygAwIBAgIUJkQQ/86O/lBUV7SXfC2LcyLs9eAwDQYJKoZIhvcNAQEL\nBQAwgbUxFDASBgNVBAYTC05ldGhlcmxhbmRzMRYwFAYDVQQIEw1Ob29yZC1CcmFi\nYW50MRIwEAYDVQQHEwlFaW5kaG92ZW4xKzApBgNVBAoTIlBoaWxpcHMgRWxlY3Ry\nb25pY3MgTmVkZXJsYW5kIEIuVi4xFDASBgNVBAsTC0hlYWx0aFN1aXRlMS4wLAYD\nVQQDEyVQaGlsaXBzIEhlYWx0aFN1aXRlIFByaXZhdGUgUG9saWN5IENBMB4XDTIw\nMTExODE5NTQ0N1oXDTIxMTExMzE5NTUxN1owdDELMAkGA1UEBhMCTkwxFzAVBgNV\nBAgTDkhhYXJsZW1tZXJtZWVyMRIwEAYDVQQHEwlIb29mZGRvcnAxDzANBgNVBAoT\nBlBhd25lZTETMBEGA1UECxMKcm9uc3dhbnNvbjESMBAGA1UEAxMJYW5keS10ZXN0\nMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEKv1wg5OkV0pJ9JyBgMEtqBVWeGdVSHED\niOjc3dfarWuD+qUKWaK3XExntORo/Ke2mxGIOyOYaDHeoc4pMZJPivpGL7UL9RvD\nYf4Jw6O+ZdbAOLvT9BPmmZG841Qy554Go4IBWDCCAVQwDgYDVR0PAQH/BAQDAgEG\nMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFHv8fKsXUSqN/wzsK1CybDEh\nq8daMB8GA1UdIwQYMBaAFMe8sLFh7iqofvjzQsaQDs98ZXJpMHAGCCsGAQUFBwEB\nBGQwYjBgBggrBgEFBQcwAoZUaHR0cHM6Ly9wa2ktcHJveHktY2xpZW50LXRlc3Qu\nZXUtd2VzdC5waGlsaXBzLWhlYWx0aHN1aXRlLmNvbS9jb3JlL3BraS9hcGkvcG9s\naWN5L2NhMBQGA1UdEQQNMAuCCWFuZHktdGVzdDBmBgNVHR8EXzBdMFugWaBXhlVo\ndHRwczovL3BraS1wcm94eS1jbGllbnQtdGVzdC5ldS13ZXN0LnBoaWxpcHMtaGVh\nbHRoc3VpdGUuY29tL2NvcmUvcGtpL2FwaS9wb2xpY3kvY3JsMA0GCSqGSIb3DQEB\nCwUAA4ICAQCY6xWRyMWoWJfaekh7v3mlIx14yCff9ub2kZTdKo9fWyyxLmnsiqRy\npZ2kVoF0vG9SSDBhCoP7nBPkRt8j7+HzixGcErVpOTi1D8XJtd1bdREvL3nHHlcA\nmSrbJMTDE+8NsDzs7s/MmAzqpkS4p0Vw+NZ5VdJqcTWFYg3KdWcaxSBkzAi8biUp\nIipK1rA7bnQSpAUvIi3jEUhmNrhv4Bt0HRDdrCR53R5isZVMAZpb0+0uTmSUrr0d\nxQpCggr0SOGi192Vx8Qz9hYDnSpBjf2d57VUw2/IJRfxRhYTzQMYrq0fXYu798MX\nxtk1oy20smAivUSWr9pFS5MB+H/oH4IXLFtFLkh81v8lypIkWCyLuAvmxVi1HpWx\na58KCUIVmavb0UU7qcHZKFdFn/HVtEOKiCScBNbYItsIiCXjQvmLNWBiLOFVLiny\nklt017KSbXU/5jHh24c4O6hbxiiss1MllFTlQvrr8a2kbS8/7gLTu3vGE7dzn5X3\nHI2ghP25mFWk4u74jio9hE/0gmq2r75HY4lRsW0jn4HRYerYyXVcbNxJQquXHsdp\nJXRhMUrJuPaZ/pHe8z3ZexmffPQOY7cmC5/yNsr1fxI3lL4VqOEocuO5PTrsPRXv\nh4V6Hqe1CWsYP7jdFZ4ztpOfbRCbXuc2oOOu/U9YXHs+aGNe+FmS6g==\n-----END CERTIFICATE-----",
      "-----BEGIN CERTIFICATE-----\nMIIHNTCCBR2gAwIBAgIUKT/4+P0M6aQJ612sKh4AMalwT90wDQYJKoZIhvcNAQEL\nBQAwgbMxFDASBgNVBAYTC05ldGhlcmxhbmRzMRYwFAYDVQQIEw1Ob29yZC1CcmFi\nYW50MRIwEAYDVQQHEwlFaW5kaG92ZW4xKzApBgNVBAoTIlBoaWxpcHMgRWxlY3Ry\nb25pY3MgTmVkZXJsYW5kIEIuVi4xFDASBgNVBAsTC0hlYWx0aFN1aXRlMSwwKgYD\nVQQDEyNQaGlsaXBzIEhlYWx0aFN1aXRlIFByaXZhdGUgUm9vdCBDQTAeFw0yMDEx\nMDYwOTE1NThaFw0yNTExMDUwOTE2MjhaMIG1MRQwEgYDVQQGEwtOZXRoZXJsYW5k\nczEWMBQGA1UECBMNTm9vcmQtQnJhYmFudDESMBAGA1UEBxMJRWluZGhvdmVuMSsw\nKQYDVQQKEyJQaGlsaXBzIEVsZWN0cm9uaWNzIE5lZGVybGFuZCBCLlYuMRQwEgYD\nVQQLEwtIZWFsdGhTdWl0ZTEuMCwGA1UEAxMlUGhpbGlwcyBIZWFsdGhTdWl0ZSBQ\ncml2YXRlIFBvbGljeSBDQTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIB\nAMcD0r7nvQFTDOPBgGFgA9bnVvRixjusB8iAkr6alTQh/RwUwXtc+il2YtFHFurb\nlVl/w1nBqCwYqLB2cf6mLujCZjUTrq6EU7atBvrf0qgEbW/IWM5R0jAvI7AoTNZD\n1YTS2ZcSajBmVaVRcPyIAtZiA88sb7zSwjPBE1S+UMni28v0ARJhVq5cmAvB2QJq\nSekjUsbwhTt3aJvprGaUCjQa3Yg3xy1Lzi/EEl7Mb8hpEgLolGxaQ73ITOwnmAFi\n+yNlBuFtHC9VL00CdZOqpnyEoLVNP8NidxUR4LvbnIuTxPWmFEgLe4riBJig0m9A\nr/inR3yiIG3RdH9lxxZhV9Du0r0g1iWxjPUElBgFGRamHH6HUr1SLk5sxopT9i5n\nYKE/IawrGS+3uk/WoVPO9wwIK2ZqbRTrJZhXz76VMpjY7ZBwjhTEevr0yIrDvBaw\nmxTIERbL0tEvWJPaFyJgcdjFvZL8Gciqp7vtnL9Kl4zGBUPNaDKEXZkDk7Z0PX5L\nGuLVeQKWjrOJY6cW6eUl/CKzlv7I0Bl/Ql/DyhicjpchZohUTv5GLeAh0U/vrldB\nUiZrDweitVmPd4z1zQZBbiDcwc600nM+Jj6dCWiwszRDtImnSLoq4R5PuaYvGJpR\nREmcn8ujf0okYe70WG9cgo6W0DvFL27yfenutOyLd/hHAgMBAAGjggE7MIIBNzAO\nBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUx7ywsWHu\nKqh++PNCxpAOz3xlcmkwHwYDVR0jBBgwFoAU/eZ7YoHbXSotmyhPzDlMzIeflWsw\nbgYIKwYBBQUHAQEEYjBgMF4GCCsGAQUFBzAChlJodHRwczovL3BraS1wcm94eS1j\nbGllbnQtdGVzdC5ldS13ZXN0LnBoaWxpcHMtaGVhbHRoc3VpdGUuY29tL2NvcmUv\ncGtpL2FwaS9yb290L2NhMGQGA1UdHwRdMFswWaBXoFWGU2h0dHBzOi8vcGtpLXBy\nb3h5LWNsaWVudC10ZXN0LmV1LXdlc3QucGhpbGlwcy1oZWFsdGhzdWl0ZS5jb20v\nY29yZS9wa2kvYXBpL3Jvb3QvY3JsMA0GCSqGSIb3DQEBCwUAA4ICAQBX3Se18iZ9\nQ5q5Su9x9SE0+NyMBOvCEIEusWFKodr7G9CAKN8HYa4NkLFfw2yZoRQ2wtFUPBo0\nDtYXnsDjNjdEQMx2KKQnCRigFQSkQErPQ03zvGxxULlreljym9N0weGnWZoFT+uI\nBdteXAbLWWk3alkGJhirfFd/mNKx4UzpqowUzzzJl1cXThKSD6J7irJZYUspOeqN\nwD+ELev1c312ZP4+GgOH36drp3edKq7HxrfTh66OAOS3JOY4CrlV/pa/4lAUj/tv\nGvh4zANg1PREBj2GrSBvknRsFTupIDC+IuPaDt9zqnA9CNZgHpcD1nLoDtA9QkOU\nP3nIJ999Kep5pj57kKTK7AhDnR8d7fMSUkF0n1ZDYsjocSnNX2VdUp45WrBvCola\nm7m9Zzp7EkWhU/DqQnCgMFfm4aXxpLpz5OGuBlL1/6j6EVW5WXguyD8Ozh7F6tuo\njW55Q+4N5C/aJmZwYNM++YnRWlb2O6coAONtjz52QViHcB/D7fWNxmTzpQcShMT3\nFDQjdC3Y6CWmQRRY4CTpfE/XxN9hjUeO7wUGiTPJ4V1/iSbI+/m8/f/8KS5opVjJ\nQ3AsiRh4Q77Cc5Y7mKE16bjwUNM6QOm32UROoqIOAFHh3cNLrSa2N7vIeGB8dVPC\n6RhjuwvAE1jojDfWvO7+PFC3fz4GVmvf3g==\n-----END CERTIFICATE-----",
      "-----BEGIN CERTIFICATE-----\nMIIHMzCCBRugAwIBAgIUB7awwVr04x/+xa1uFm7DAz85VdgwDQYJKoZIhvcNAQEL\nBQAwgbMxFDASBgNVBAYTC05ldGhlcmxhbmRzMRYwFAYDVQQIEw1Ob29yZC1CcmFi\nYW50MRIwEAYDVQQHEwlFaW5kaG92ZW4xKzApBgNVBAoTIlBoaWxpcHMgRWxlY3Ry\nb25pY3MgTmVkZXJsYW5kIEIuVi4xFDASBgNVBAsTC0hlYWx0aFN1aXRlMSwwKgYD\nVQQDEyNQaGlsaXBzIEhlYWx0aFN1aXRlIFByaXZhdGUgUm9vdCBDQTAeFw0yMDEx\nMDYwOTE1NTVaFw0zMDExMDQwOTE2MTlaMIGzMRQwEgYDVQQGEwtOZXRoZXJsYW5k\nczEWMBQGA1UECBMNTm9vcmQtQnJhYmFudDESMBAGA1UEBxMJRWluZGhvdmVuMSsw\nKQYDVQQKEyJQaGlsaXBzIEVsZWN0cm9uaWNzIE5lZGVybGFuZCBCLlYuMRQwEgYD\nVQQLEwtIZWFsdGhTdWl0ZTEsMCoGA1UEAxMjUGhpbGlwcyBIZWFsdGhTdWl0ZSBQ\ncml2YXRlIFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDB\nDqi7GQ9oUMUU+6WMkVG+OeU90um2riOfCqICnDrXjYuGt6rjuyDO9Lk10OquVRaQ\nf0gTAWirLfXbx0Ifrh0tPNB3XENTGjzf+K8zHxhHt2m18WoBWoCo8Bhc+v2UHqQy\nCuKZhhZ9Wma4+kuowfmKJVJZD9zfGgkHRoqSSV+MyphggNukMnfjArSV0jHXOLc+\nR8XMGJw9O++6kB1dOcxuj5Xmmv3bRyxRg1I9pWUBPovz400TZI1qz30jCj0TireS\nsJyPD6SFH/4bSONEyAZ+n8U7m4JxwCrUlEnQ/zXSt7ZroKslYfAG/xRp1Jm5TqiJ\nt1hJtzq7gRDCLkTGYRtaKoRUGpoZeES9GKVBlGKqNo5gCFbMpyulXf+IpTgQt9j2\nscG6/l2qWhmPtdJ7atYzmL07/ooDk7SGV8fpetfRN0fdGw7Bn3NFk6wYsIfRZZKE\n+68xizi6BZDhT54sHgZs1bTcYroAEWAijB6lMmBiK4kEDnp2ZEBSAwOh7ablepYf\nDbjht+Blfq7FzYeZUVCbILGl3CX2IVpXDTWzNjlo/hHGCaX9nh+5WaJVnNaihbwF\nb86InqJ98oWfhl0PK+sh0ELJEt9dFQrSWxm6SSbEM4LT/Tr1A6dAqzz9Z65euXuA\nsAHqZQGE25E/qQg9uWJxmiZ6bzsoh/zV3WkCUVasTwIDAQABo4IBOzCCATcwDgYD\nVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFP3me2KB210q\nLZsoT8w5TMyHn5VrMB8GA1UdIwQYMBaAFP3me2KB210qLZsoT8w5TMyHn5VrMG4G\nCCsGAQUFBwEBBGIwYDBeBggrBgEFBQcwAoZSaHR0cHM6Ly9wa2ktcHJveHktY2xp\nZW50LXRlc3QuZXUtd2VzdC5waGlsaXBzLWhlYWx0aHN1aXRlLmNvbS9jb3JlL3Br\naS9hcGkvcm9vdC9jYTBkBgNVHR8EXTBbMFmgV6BVhlNodHRwczovL3BraS1wcm94\neS1jbGllbnQtdGVzdC5ldS13ZXN0LnBoaWxpcHMtaGVhbHRoc3VpdGUuY29tL2Nv\ncmUvcGtpL2FwaS9yb290L2NybDANBgkqhkiG9w0BAQsFAAOCAgEAU9/eE2dxC0Ph\n7d/CXCFIlVqBEcqzFfJq8r9C2qQhAPM7cshDIq5m+wWod0zr3vA+wzbM7Wh/iu60\nEk8/p3E00CJPhiZeTcCejyYtJEqI1poVSnR/saQmuQ+Tahj6FkDRlo5a25NBBLKw\n8Dx0g2NH6GEeHW4tW2yvV3c8/HcVMELGTcLuJW5cSe6W/Xt9WUyCJ852UDeGhvb4\nRdxAzWDZrP48QRzc27YT7jg54MNPAKzf0mrNdmvtC5iOB19tSDK0FYTl1+Ab2/a5\ngYAqnUSfVbKjudoaUaxU6vJKkcaTK1qOQhq0YelZZmoEPQSlrMtGIuExSy6UdMeo\nZxKBNQ9LXAu6r7elOv4pvcXk5v5q0S/uouAMxxxl4a5L5qyvJ3DPPVFHa3T4gJC0\nmbToyqcYBi3HsqyryUXhpZmKlEX63U+tLmIsbxUJCqfCz97nNfgdUv6UUqD+INre\nWSxEU8NbKQdbmM5WymUpTLZVx1JhPyl+DnAyWA9nfnlhU6IH+zpfkenphMsdmRSS\n0TX2lxsQe5i1DXRw2o+XZxsJx2JWLIf74nuNourXdQs/taaibkTj6Y1dmrpMeDKx\nGyUF/ncBCYaVXK+6DzS2kUCj9bPGVbGoXadaJxxGe9jGpOLcTR1xcQ/WSMAQuIZa\nwo3KBVGxGCMPQZ8FeqGowJ0yDB8GxZ0=\n-----END CERTIFICATE-----"
    ],
    "certificate": "-----BEGIN CERTIFICATE-----\nMIIDSDCCAs2gAwIBAgIUIVPQZChvUoCsfv+e++1zquZIEQ0wCgYIKoZIzj0EAwIw\ndDELMAkGA1UEBhMCTkwxFzAVBgNVBAgTDkhhYXJsZW1tZXJtZWVyMRIwEAYDVQQH\nEwlIb29mZGRvcnAxDzANBgNVBAoTBlBhd25lZTETMBEGA1UECxMKcm9uc3dhbnNv\nbjESMBAGA1UEAxMJYW5keS10ZXN0MB4XDTIwMTExOTA2NDIxMFoXDTIxMDUxOTE3\nNDI0MFowJTELMAkGA1UEBhMCTkwxFjAUBgNVBAMTDXRlc3QuMWUxMDAuaW8wdjAQ\nBgcqhkjOPQIBBgUrgQQAIgNiAAQqzZ7/6OonroRCfU4EvJ+fbkFGkN1GjycWfnNT\nLR/X8FhQ0kciK+Kv4u0VuBp70PceKllqcXXHRuPvTLHpbKflTBMilE5o+t8S9jfn\naV44jxNCgDNVkdUnpN7JwyxqGs2jggFtMIIBaTAOBgNVHQ8BAf8EBAMCA6gwHQYD\nVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQWBBRSgmUpo7JeW0NR\nV/+PkhofGG3SCDAfBgNVHSMEGDAWgBR7/HyrF1Eqjf8M7CtQsmwxIavHWjBzBggr\nBgEFBQcBAQRnMGUwYwYIKwYBBQUHMAKGV2h0dHBzOi8vcGtpLXByb3h5LWNsaWVu\ndC10ZXN0LmV1LXdlc3QucGhpbGlwcy1oZWFsdGhzdWl0ZS5jb20vY29yZS9wa2kv\nYXBpL2FuZHktdGVzdC9jYTAYBgNVHREEETAPgg10ZXN0LjFlMTAwLmlvMGkGA1Ud\nHwRiMGAwXqBcoFqGWGh0dHBzOi8vcGtpLXByb3h5LWNsaWVudC10ZXN0LmV1LXdl\nc3QucGhpbGlwcy1oZWFsdGhzdWl0ZS5jb20vY29yZS9wa2kvYXBpL2FuZHktdGVz\ndC9jcmwwCgYIKoZIzj0EAwIDaQAwZgIxALo/Jog4roubu/R6iOJrjqnD9n1tGlZG\n+Oh5736fAD7KIlsA3XPbRf3/IiNLHmvpsQIxAILrzL6hHUmtAmXSsT/OGcfk+RbX\nq0KoJbXtM/VZI3IDqsBb3ywXRQut4F5TJRPD2g==\n-----END CERTIFICATE-----",
    "expiration": 1621446160,
    "issuing_ca": "-----BEGIN CERTIFICATE-----\nMIIFZDCCA0ygAwIBAgIUJkQQ/86O/lBUV7SXfC2LcyLs9eAwDQYJKoZIhvcNAQEL\nBQAwgbUxFDASBgNVBAYTC05ldGhlcmxhbmRzMRYwFAYDVQQIEw1Ob29yZC1CcmFi\nYW50MRIwEAYDVQQHEwlFaW5kaG92ZW4xKzApBgNVBAoTIlBoaWxpcHMgRWxlY3Ry\nb25pY3MgTmVkZXJsYW5kIEIuVi4xFDASBgNVBAsTC0hlYWx0aFN1aXRlMS4wLAYD\nVQQDEyVQaGlsaXBzIEhlYWx0aFN1aXRlIFByaXZhdGUgUG9saWN5IENBMB4XDTIw\nMTExODE5NTQ0N1oXDTIxMTExMzE5NTUxN1owdDELMAkGA1UEBhMCTkwxFzAVBgNV\nBAgTDkhhYXJsZW1tZXJtZWVyMRIwEAYDVQQHEwlIb29mZGRvcnAxDzANBgNVBAoT\nBlBhd25lZTETMBEGA1UECxMKcm9uc3dhbnNvbjESMBAGA1UEAxMJYW5keS10ZXN0\nMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEKv1wg5OkV0pJ9JyBgMEtqBVWeGdVSHED\niOjc3dfarWuD+qUKWaK3XExntORo/Ke2mxGIOyOYaDHeoc4pMZJPivpGL7UL9RvD\nYf4Jw6O+ZdbAOLvT9BPmmZG841Qy554Go4IBWDCCAVQwDgYDVR0PAQH/BAQDAgEG\nMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFHv8fKsXUSqN/wzsK1CybDEh\nq8daMB8GA1UdIwQYMBaAFMe8sLFh7iqofvjzQsaQDs98ZXJpMHAGCCsGAQUFBwEB\nBGQwYjBgBggrBgEFBQcwAoZUaHR0cHM6Ly9wa2ktcHJveHktY2xpZW50LXRlc3Qu\nZXUtd2VzdC5waGlsaXBzLWhlYWx0aHN1aXRlLmNvbS9jb3JlL3BraS9hcGkvcG9s\naWN5L2NhMBQGA1UdEQQNMAuCCWFuZHktdGVzdDBmBgNVHR8EXzBdMFugWaBXhlVo\ndHRwczovL3BraS1wcm94eS1jbGllbnQtdGVzdC5ldS13ZXN0LnBoaWxpcHMtaGVh\nbHRoc3VpdGUuY29tL2NvcmUvcGtpL2FwaS9wb2xpY3kvY3JsMA0GCSqGSIb3DQEB\nCwUAA4ICAQCY6xWRyMWoWJfaekh7v3mlIx14yCff9ub2kZTdKo9fWyyxLmnsiqRy\npZ2kVoF0vG9SSDBhCoP7nBPkRt8j7+HzixGcErVpOTi1D8XJtd1bdREvL3nHHlcA\nmSrbJMTDE+8NsDzs7s/MmAzqpkS4p0Vw+NZ5VdJqcTWFYg3KdWcaxSBkzAi8biUp\nIipK1rA7bnQSpAUvIi3jEUhmNrhv4Bt0HRDdrCR53R5isZVMAZpb0+0uTmSUrr0d\nxQpCggr0SOGi192Vx8Qz9hYDnSpBjf2d57VUw2/IJRfxRhYTzQMYrq0fXYu798MX\nxtk1oy20smAivUSWr9pFS5MB+H/oH4IXLFtFLkh81v8lypIkWCyLuAvmxVi1HpWx\na58KCUIVmavb0UU7qcHZKFdFn/HVtEOKiCScBNbYItsIiCXjQvmLNWBiLOFVLiny\nklt017KSbXU/5jHh24c4O6hbxiiss1MllFTlQvrr8a2kbS8/7gLTu3vGE7dzn5X3\nHI2ghP25mFWk4u74jio9hE/0gmq2r75HY4lRsW0jn4HRYerYyXVcbNxJQquXHsdp\nJXRhMUrJuPaZ/pHe8z3ZexmffPQOY7cmC5/yNsr1fxI3lL4VqOEocuO5PTrsPRXv\nh4V6Hqe1CWsYP7jdFZ4ztpOfbRCbXuc2oOOu/U9YXHs+aGNe+FmS6g==\n-----END CERTIFICATE-----",
    "private_key": "-----BEGIN EC PRIVATE KEY-----\nMIGkAgEBBDB1AH5vdB5dWuQPoHJTLq1h9O1y0VQkfn292Cdne5sVYXPCHC9Y6pSb\njtEwmuYdYgWgBwYFK4EEACKhZANiAAQqzZ7/6OonroRCfU4EvJ+fbkFGkN1GjycW\nfnNTLR/X8FhQ0kciK+Kv4u0VuBp70PceKllqcXXHRuPvTLHpbKflTBMilE5o+t8S\n9jfnaV44jxNCgDNVkdUnpN7JwyxqGs0=\n-----END EC PRIVATE KEY-----",
    "private_key_type": "ec",
    "serial_number": "21:53:d0:64:28:6f:52:80:ac:7e:ff:9e:fb:ed:73:aa:e6:48:11:0d"
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}`
	returnCert := func(logicalPath, role string) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case "POST", "GET":
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, response)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
	}
	serial := "21:53:d0:64:28:6f:52:80:ac:7e:ff:9e:fb:ed:73:aa:e6:48:11:0d"
	logicalPath := "ron-swanson"
	role := "ec384"
	muxPKI.HandleFunc("/core/pki/api/"+logicalPath+"/cert/"+serial, returnCert(logicalPath, role))
	muxPKI.HandleFunc("/core/pki/api/"+logicalPath+"/issue/"+role, returnCert(logicalPath, role))
	muxPKI.HandleFunc("/core/pki/api/"+logicalPath+"/sign/"+role, returnCert(logicalPath, role))
	cert, resp, err := pkiClient.Services.IssueCertificate(logicalPath, role, pki.CertificateRequest{
		CommonName: "test",
		TTL:        "4355h",
	})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, cert) {
		return
	}
	key, err := cert.Data.GetPrivateKey()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, key) {
		return
	}
	assert.IsType(t, &ecdsa.PrivateKey{}, key)

	certificate, err := cert.Data.GetCertificate()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, certificate) {
		return
	}
	assert.IsType(t, &x509.Certificate{}, certificate)

	cert, resp, err = pkiClient.Services.GetCertificateBySerial(logicalPath, serial)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, cert) {
		return
	}

	cert, resp, err = pkiClient.Services.Sign(logicalPath, role, pki.SignRequest{
		CSR:        "",
		CommonName: "1e100.io",
		Format:     "pem",
	})
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

func TestServicesErrors(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	_, _, err := pkiClient.Services.GetCertificateBySerial("logicalPath", "serial")
	assert.NotNil(t, err)
	_, _, err = pkiClient.Services.IssueCertificate("logicalPath", "role", pki.CertificateRequest{})
	assert.NotNil(t, err)
	_, _, err = pkiClient.Services.GetPolicyCRL()
	assert.NotNil(t, err)
	_, _, err = pkiClient.Services.GetRootCRL()
	assert.NotNil(t, err)
	_, _, err = pkiClient.Services.GetRootCA()
	assert.NotNil(t, err)
}
