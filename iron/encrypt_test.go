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
	pubkey := []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdS2oE9+dhexZc3/sEtI+a6ZKt
6FwBZaAgytdkQ7sX4FwbZAdJ7zFS1m0gDezyFTBJSPVjYOKYr0fu1ao/xkNkKnnz
J2WkW6qsDNKwJgrHiCO1asnoW5XWtk8Yc4kKkg63REuV20x+QoD6onTCo3T2DfUI
vZ8QOSJQ7NotGuO2wwIDAQAB
-----END PUBLIC KEY-----`)
	fixed := FormatBrokenPubkey(pubkey)
	encrypted, err := EncryptPayload([]byte(fixed), []byte("Yo"))
	assert.Nil(t, err)
	assert.Equal(t, 212, len(encrypted))
}
