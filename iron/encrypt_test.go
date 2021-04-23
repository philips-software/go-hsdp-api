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
	privkey := []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
	pubkey := []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDDLqIxJYtuJKFl4IlvJjaK2ZQV
bEwgR5Daxch17rSyj41FmVC/1ypsOWbiUnFkrnZThRyZxboKVLI8LWfIyruyBCX5
oMnuk0nKbftpGy1WEf+ME7XUFEZbYau5rQqXm2kNJJhGwYnm07rqIKNaL4bPQoRt
2x1I+rUhzMi3WL+P3QIDAQAB
-----END PUBLIC KEY-----`)
	fixed := FormatBrokenPubkey(pubkey)
	encrypted, err := EncryptPayload([]byte(fixed), []byte("I’m gonna count to one. One."))
	assert.Nil(t, err)

	decrypted, err := DecryptPayload(privkey, encrypted)
	assert.Nil(t, err)
	assert.Equal(t, "I’m gonna count to one. One.", string(decrypted))
}
