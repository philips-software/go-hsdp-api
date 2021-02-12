package stl_test

import (
	"context"
	"github.com/philips-software/go-hsdp-api/stl"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestGetCustomCertByID(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "appCustomCert": {
      "id": 53,
      "deviceId": 53615,
      "name": "terrakube.com",
      "key": "-----BEGIN EC PRIVATE KEY-----\nFAKE\n-----END EC PRIVATE KEY-----",
      "cert": "-----BEGIN CERTIFICATE-----\nMIIEtDCCApygAwIBAgIUNuPCrKttZ2Wrf12rRa/dDb3kGiQwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLY29tbW9uLm5hbWUwHhcNMjEwMjA0MTE0NzM1WhcNMjEw\nMzA2MjE0ODA1WjAYMRYwFAYDVQQDEw10ZXJyYWt1YmUuY29tMHYwEAYHKoZIzj0C\nAQYFK4EEACIDYgAEEFiEVayyclxPAe3MN4u3oBj3YO9L2UR3k19qBw6SPjUkneig\nYrAW12fHeZgZZ7awpStsZy6cdwapLAa+0grTauPs7nKV12cLfOB0hQG4t7MJquR/\nXS+PmfdMrTjThnVco4IBpDCCAaAwDgYDVR0PAQH/BAQDAgOoMB0GA1UdJQQWMBQG\nCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQUyHiRrI5kTA9Q+l7Jmk1iBKju\nsKgwHwYDVR0jBBgwFoAUy0DPAuyXCjxGzweAPrXtvVZBhO0wgY4GCCsGAQUFBwEB\nBIGBMH8wfQYIKwYBBQUHMAKGcWh0dHA6Ly9wa2ktcHJveHktY2xpZW50LXRlc3Qu\nZXUtd2VzdC5waGlsaXBzLWhlYWx0aHN1aXRlLmNvbS9jb3JlL3BraS9hcGkvYWNm\nMWFlMDgtNTU5OS00ZTc3LWIyMGItMThlYmVhZmI2MjVjL2NhMBgGA1UdEQQRMA+C\nDXRlcnJha3ViZS5jb20wgYMGA1UdHwR8MHoweKB2oHSGcmh0dHA6Ly9wa2ktcHJv\neHktY2xpZW50LXRlc3QuZXUtd2VzdC5waGlsaXBzLWhlYWx0aHN1aXRlLmNvbS9j\nb3JlL3BraS9hcGkvYWNmMWFlMDgtNTU5OS00ZTc3LWIyMGItMThlYmVhZmI2MjVj\nL2NybDANBgkqhkiG9w0BAQsFAAOCAgEAM7xB7ymibfENe05wFIE9LjL06CJi9Hip\nuy1m9kLOfbysfhJLcRUUKXFw1v8lUjhw6IuiOEY8WJDX7F35XLy29lvseh/rYMtm\nrHE3w2p3nzOmTdVUL4JrB0ZNk91FuyrK3G5X3A6P5HOz8DMeWnpJVsIkt5AaP5Tn\nt0PjzBxvlbzVgpXRdUQI5u1YrxMU7v9dKRcXT067oHHN7mR4hqT+JrOINqTIXkFf\n1Gzd2XZ6VKLGQ8OjA2g11ShI4SnTm2uLWmzUuj8ARSDuqSzsZY+7+Rou4YYezREU\n1OZPWIJb55vi6frcufQU3Nf5gHDmSCMrYlpqHqLmyojOUXaA047Bwjjjzu6Rxky/\nGbbjoBMxKBy/YuLj0stLX5JOICWFHFN3rASxCVx9M6stPQ4RnTPe7xb6zkBVGaw4\nfg6DoPkdVGCSxJeKFdxLSuPXpPDj6J1YIQKWiwjcflGLWYPo996AeChDXA4tlx54\n2L7M2J5t/1oedR6Y3F3RtWtAFDJdIY+N2Hgf4cNMxHUfk62o+LrrXAaRJKxE7um6\nTLIUUzwtGG7QDZMutiv3f2d1/7MtjbTUEYkCIySUO+vzJZwDfXPA1TpTcAElM0FI\nPKgEFHzJpe4qL1/O+NGeyZv43/UkDvQyvEnLLJ/a2rBnmZb8MwzqTJqyV8+fWoNh\nYVmko1j2+jg=\n-----END CERTIFICATE-----"
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	cert, err := client.Certs.GetCustomCertByID(ctx, 1)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, cert) {
		return
	}
}

func TestCertsService_GetCustomCertsBySerial(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	muxSTL.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
  "data": {
    "appCustomCerts": {
      "edges": [
        {
          "node": {
            "id": 47,
            "deviceId": 53615,
            "name": "test",
            "key": "-----BEGIN EC PRIVATE KEY-----\nFAKE\n-----END EC PRIVATE KEY-----",
            "cert": "-----BEGIN CERTIFICATE-----\nMIIEtDCCApygAwIBAgIUNuPCrKttZ2Wrf12rRa/dDb3kGiQwDQYJKoZIhvcNAQEL\nBQAwFjEUMBIGA1UEAxMLY29tbW9uLm5hbWUwHhcNMjEwMjA0MTE0NzM1WhcNMjEw\nMzA2MjE0ODA1WjAYMRYwFAYDVQQDEw10ZXJyYWt1YmUuY29tMHYwEAYHKoZIzj0C\nAQYFK4EEACIDYgAEEFiEVayyclxPAe3MN4u3oBj3YO9L2UR3k19qBw6SPjUkneig\nYrAW12fHeZgZZ7awpStsZy6cdwapLAa+0grTauPs7nKV12cLfOB0hQG4t7MJquR/\nXS+PmfdMrTjThnVco4IBpDCCAaAwDgYDVR0PAQH/BAQDAgOoMB0GA1UdJQQWMBQG\nCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQUyHiRrI5kTA9Q+l7Jmk1iBKju\nsKgwHwYDVR0jBBgwFoAUy0DPAuyXCjxGzweAPrXtvVZBhO0wgY4GCCsGAQUFBwEB\nBIGBMH8wfQYIKwYBBQUHMAKGcWh0dHA6Ly9wa2ktcHJveHktY2xpZW50LXRlc3Qu\nZXUtd2VzdC5waGlsaXBzLWhlYWx0aHN1aXRlLmNvbS9jb3JlL3BraS9hcGkvYWNm\nMWFlMDgtNTU5OS00ZTc3LWIyMGItMThlYmVhZmI2MjVjL2NhMBgGA1UdEQQRMA+C\nDXRlcnJha3ViZS5jb20wgYMGA1UdHwR8MHoweKB2oHSGcmh0dHA6Ly9wa2ktcHJv\neHktY2xpZW50LXRlc3QuZXUtd2VzdC5waGlsaXBzLWhlYWx0aHN1aXRlLmNvbS9j\nb3JlL3BraS9hcGkvYWNmMWFlMDgtNTU5OS00ZTc3LWIyMGItMThlYmVhZmI2MjVj\nL2NybDANBgkqhkiG9w0BAQsFAAOCAgEAM7xB7ymibfENe05wFIE9LjL06CJi9Hip\nuy1m9kLOfbysfhJLcRUUKXFw1v8lUjhw6IuiOEY8WJDX7F35XLy29lvseh/rYMtm\nrHE3w2p3nzOmTdVUL4JrB0ZNk91FuyrK3G5X3A6P5HOz8DMeWnpJVsIkt5AaP5Tn\nt0PjzBxvlbzVgpXRdUQI5u1YrxMU7v9dKRcXT067oHHN7mR4hqT+JrOINqTIXkFf\n1Gzd2XZ6VKLGQ8OjA2g11ShI4SnTm2uLWmzUuj8ARSDuqSzsZY+7+Rou4YYezREU\n1OZPWIJb55vi6frcufQU3Nf5gHDmSCMrYlpqHqLmyojOUXaA047Bwjjjzu6Rxky/\nGbbjoBMxKBy/YuLj0stLX5JOICWFHFN3rASxCVx9M6stPQ4RnTPe7xb6zkBVGaw4\nfg6DoPkdVGCSxJeKFdxLSuPXpPDj6J1YIQKWiwjcflGLWYPo996AeChDXA4tlx54\n2L7M2J5t/1oedR6Y3F3RtWtAFDJdIY+N2Hgf4cNMxHUfk62o+LrrXAaRJKxE7um6\nTLIUUzwtGG7QDZMutiv3f2d1/7MtjbTUEYkCIySUO+vzJZwDfXPA1TpTcAElM0FI\nPKgEFHzJpe4qL1/O+NGeyZv43/UkDvQyvEnLLJ/a2rBnmZb8MwzqTJqyV8+fWoNh\nYVmko1j2+jg=\n-----END CERTIFICATE-----"
          }
        }
      ]
    }
  }
}`)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	ctx := context.Background()
	cert, err := client.Certs.GetCustomCertsBySerial(ctx, "serial")
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, cert) {
		return
	}
}

func TestCertsServiceErrors(t *testing.T) {
	teardown, err := setup(t)
	if !assert.Nil(t, err) {
		return
	}
	defer teardown()

	serial := "A444900Z0822111"
	ctx := context.Background()
	c, err := client.Certs.CreateCustomCert(ctx, stl.CreateAppCustomCertInput{})
	assert.NotNil(t, err)
	assert.Nil(t, c)
	l, err := client.Certs.GetCustomCertsBySerial(ctx, serial)
	assert.NotNil(t, err)
	assert.Nil(t, l)
	c, err = client.Certs.UpdateCustomCert(ctx, stl.UpdateAppCustomCertInput{})
	assert.NotNil(t, err)
	assert.Nil(t, c)
	c, err = client.Certs.GetCustomCertByID(ctx, 1)
	assert.NotNil(t, err)
	assert.Nil(t, c)
	ok, err := client.Certs.DeleteCustomCert(ctx, stl.DeleteAppCustomCertInput{})
	assert.NotNil(t, err)
	assert.False(t, ok)
}
