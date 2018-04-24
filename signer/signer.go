package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"
)

var (
	LOG_TIME_FORMAT      = "2006-01-02T15:04:05.000Z07:00"
	TIME_FORMAT          = time.RFC3339
	AUTHORIZATION_HEADER = "hsdp-api-signature"
	SIGNED_DATE_HEADER   = "SignedDate"
	DEFAULT_PREFIX_64    = "REhQV1M="
	ALGORITHM_NAME       = "HmacSHA256"
	errSignatureExpired  = errors.New("Signture expired")
	errInvalidSignature  = errors.New("Invalid signature")
	errInvalidAlgorithm  = errors.New("Invalid algorithm")
	errInvalidHeaderList = errors.New("Invalid header list")
	errInvalidCredential = errors.New("Invalid credential")
	errNotSupportedYet   = errors.New("Missing implementation. Please contact the author(s)")
)

type Signer struct {
	sharedKey    string
	sharedSecret string
	prefix       string
	nowFunc      NowFunc
}

type NowFunc func() time.Time

func New(sharedKey, sharedSecret, prefix string, nowFunc NowFunc) (*Signer, error) {
	signer := &Signer{
		sharedKey:    sharedKey,
		sharedSecret: sharedSecret,
		prefix:       prefix,
	}
	if signer.prefix == "" {
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(DEFAULT_PREFIX_64)))
		l, _ := base64.StdEncoding.Decode(decoded, []byte(DEFAULT_PREFIX_64))
		signer.prefix = string(decoded[:l])
	}
	if nowFunc != nil {
		signer.nowFunc = nowFunc
	} else {
		signer.nowFunc = func() time.Time {
			return time.Now()
		}
	}
	return signer, nil
}

func (s *Signer) SignRequest(request *http.Request) error {
	signTime := s.nowFunc().UTC().Format(TIME_FORMAT)

	seed1 := base64.StdEncoding.EncodeToString([]byte(signTime))

	hashedSeed := hash([]byte(seed1), []byte(s.prefix+s.sharedSecret))

	signature := base64.StdEncoding.EncodeToString(hashedSeed)

	authorization := ALGORITHM_NAME + ";" +
		"Credential:" + s.sharedKey + ";" +
		"SignedHeaders:SignedDate" + ";" +
		"Signature:" + signature

	request.Header.Set(AUTHORIZATION_HEADER, authorization)
	request.Header.Set(SIGNED_DATE_HEADER, signTime)
	return nil
}

func (s *Signer) ValidateRequest(request *http.Request) (bool, error) {
	signature := request.Header.Get(AUTHORIZATION_HEADER)
	signedDate := request.Header.Get(SIGNED_DATE_HEADER)

	comps := strings.Split(signature, ";")
	if len(comps) != 4 {
		return false, errInvalidSignature
	}
	if comps[0] != ALGORITHM_NAME {
		return false, errInvalidAlgorithm
	}
	credential := strings.TrimPrefix(comps[1], "Credential:")
	if credential != s.sharedKey {
		return false, errInvalidCredential
	}

	headers := strings.Split(strings.TrimPrefix(comps[2], "SignedHeaders:"), ",")
	if len(headers) < 1 {
		return false, errInvalidHeaderList
	}
	currentSeed := []byte("")
	currentKey := []byte("")
	for _, h := range headers {
		if len(currentKey) == 0 {
			currentKey = []byte(request.Header.Get(h)) // SignedDate!
			continue
		}
		switch h {
		case "body":
			return false, errNotSupportedYet
		case "method":
			return false, errNotSupportedYet
		case "URI":
			return false, errNotSupportedYet
		default:
			currentSeed = []byte(request.Header.Get(h))
		}
		currentKey = hash(currentSeed, currentKey)
	}

	finalHMAC := base64.StdEncoding.EncodeToString([]byte(currentKey))

	hashedSeed := hash([]byte(finalHMAC), []byte(s.prefix+s.sharedSecret))

	signature = base64.StdEncoding.EncodeToString(hashedSeed)
	receivedSignature := strings.TrimPrefix(comps[3], "Signature:")

	if signature != receivedSignature {
		return false, errInvalidSignature
	}

	signed, err := time.Parse(TIME_FORMAT, signedDate)
	if err != nil {
		return false, err
	}
	now := s.nowFunc()
	if now.Sub(signed).Seconds() > 900 {
		return false, errSignatureExpired
	}
	return true, nil
}

func hash(data []byte, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
