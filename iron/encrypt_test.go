package iron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBrokenPubKey(t *testing.T) {
	pubkey := []byte("-----BEGIN PUBLIC KEY----- MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdS2oE9+dhexZc3/sEtI+a6ZKt 6FwBZaAgytdkQ7sX4FwbZAdJ7zFS1m0gDezyFTBJSPVjYOKYr0fu1ao/xkNkKnnz J2WkW6qsDNKwJgrHiCO1asnoW5XWtk8Yc4kKkg63REuV20x+QoD6onTCo3T2DfUI vZ8QOSJQ7NotGuO2wwIDAQAB -----END PUBLIC KEY-----")
	fixedPubkey := []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdS2oE9+dhexZc3/sEtI+a6ZKt
6FwBZaAgytdkQ7sX4FwbZAdJ7zFS1m0gDezyFTBJSPVjYOKYr0fu1ao/xkNkKnnz
J2WkW6qsDNKwJgrHiCO1asnoW5XWtk8Yc4kKkg63REuV20x+QoD6onTCo3T2DfUI
vZ8QOSJQ7NotGuO2wwIDAQAB
-----END PUBLIC KEY-----`)

	fixed := FormatBrokenPubkey(pubkey)
	assert.Equal(t, fixedPubkey, fixed)

	noop := FormatBrokenPubkey(fixed)
	assert.Equal(t, fixedPubkey, noop)
}

func TestBrokenPubkey(t *testing.T) {
	pubkey := []byte("broken!!!!")
	_, err := EncryptPayload([]byte(pubkey), []byte("Yo"))
	assert.NotNil(t, err)
}

func TestEncryptPayload(t *testing.T) {
	keys := []struct {
		PrivateKey []byte
		PublicKey  []byte
	}{
		{PrivateKey: []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDDLqIxJYtuJKFl4IlvJjaK2ZQVbEwgR5Daxch17rSyj41FmVC/
1ypsOWbiUnFkrnZThRyZxboKVLI8LWfIyruyBCX5oMnuk0nKbftpGy1WEf+ME7XU
FEZbYau5rQqXm2kNJJhGwYnm07rqIKNaL4bPQoRt2x1I+rUhzMi3WL+P3QIDAQAB
AoGBAICG8N8ULiC1lmKT3WyH6Vq9tDn3Opn3BnhJzZt7ORpsVUcDkp0BfzoNAqb+
SxVVnS2adh78iWnMJCJkc/dRKQ8FW86wknomvLKp3O11hGOwuSUlFK6HzKS92PxH
GS64yZiXUpdBMuTjnfwDLWV9kaiCqN4uC3HcXM8peKyNj+sBAkEA8Ofln7EPni/W
RF1IQnaB1BASNkRpc3FhMXGfmN+Asphv7FmSwvYRrYrcwzX5yrxZTF2M/fxmE2k9
cy5LHC+szQJBAM9pVXGJX1Fo3UYR4HtnvKZbWweAcEXLAiVrqCMVoPJN3YpfN/5s
H522MCSjWn3aQE+ZBzbns+ZU3Suw1Wixb1ECQQDE19dKyvTF/rSHm+klVYvz6UXY
TcIUcDpIml0cHtQcGm6pou9GmqYLNYH5iCsZOxmESpSgHBKUHdP2P4dj+pipAkEA
pJAwiNqz1AXduqCoeYE/PsaxHOydJ+MAmuwmBWA9yMJbClSuOqFTHHDXFdq+C6jE
6eLCxJ9mL1QZ/3ZYfK57YQJAZi+h0dWot/ARxES7HBXTnJQsBhwA6vf3VxOwr9YY
34BHxxQDi9t+5BmpUXs+nXFtLYmw2iGnc3ev1jAH85jUUQ==
-----END RSA PRIVATE KEY-----`),
			PublicKey: []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDLqIxJYtuJKFl4IlvJjaK2ZQV
bEwgR5Daxch17rSyj41FmVC/1ypsOWbiUnFkrnZThRyZxboKVLI8LWfIyruyBCX5
oMnuk0nKbftpGy1WEf+ME7XUFEZbYau5rQqXm2kNJJhGwYnm07rqIKNaL4bPQoRt
2x1I+rUhzMi3WL+P3QIDAQAB
-----END PUBLIC KEY-----`),
		},
		{
			PrivateKey: []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`),
			PublicKey: []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAziySK3sU+4XvHSosRNmZ
hbVIJdLgynPXsnlSgHhbHbSVlushZ8DHN4+xwJKkijqglv/VtAlsyzH1ppWePjNs
bwlqhd/vZxeXYHXWpfnH2vYkEd3+awkcmJB7t1Xb3iiAO6hIIWhaRsRhP19jhCH2
foNNxtezv2II5kiMPlnTqvxu9I3QjazhKQXF/bI13Yjw4cDqld2w0dOMdeb31XHx
RIOMOH4X8biPWWzAc0p0MMQ4W2M7X8dx92VojbUINwZEsSPM135gKL1PCULwnS6Q
CvYNkmzzjLTfT6iOna/Ze80dBEXwHFkLH5DolfSNXf7/3FjB2BHa64ejZlei3cFl
mgNXUGiHUYqfUtjE2P27Cqn2LkQnGMeRk1nniSUgkmFBcCwPmxsrBhsJ0+1ubrpl
sO2upE9KElR2F7HrKrs9ddBxPA8bJKgCYOlMkG+0kRPCgl8CbSO04qR+pvZimJUY
aftnXea/ylBWJ8SOtLgavIAT7eaUMFMQDhMs5IIgFhHe/zAqnk1vM730bKky47Sc
1HDuo3o47l/vRwo5ILDvP3zQepSZ4V3tsxn8LnnurOA9gM7A1cbzZX2g23Zgyfdd
0/MsbJtdrnBOO2d6kaGEMzXrLEMqY19nCthjkObtOhPsp/4SXO/ruwlUaLmdzD+w
yQH8MI8t9FirCX6dMx7OhhUCAwEAAQ==
-----END PUBLIC KEY-----
`),
		},
	}

	for _, key := range keys {
		encrypted, err := EncryptPayload(key.PublicKey, []byte("I’m gonna count to one. One!"))
		if !assert.Nil(t, err) {
			return
		}
		decrypted, err := DecryptPayload(key.PrivateKey, encrypted)
		if !assert.Nil(t, err) {
			return
		}
		assert.Equal(t, "I’m gonna count to one. One!", string(decrypted))

		_, err = DecryptPayload([]byte(""), "bogus")
		assert.NotNil(t, err)
	}
}
