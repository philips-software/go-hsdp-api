package iam

import (
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesCRUD(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	serviceName := "testservice"
	serviceDescription := "Service description"
	applicationID := "b0889958-3762-4427-af07-2d6268c26988"
	serviceID := serviceName + ".testapp.testprop@testdev.devorg.1e100.io"
	id := "2c266886-f918-4223-941d-437cb3cd09e8"
	orgID := "bda40124-54fa-4967-b2fb-23dcc4e0ad1a"

	muxIDM.HandleFunc("/authorize/identity/Service", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{
				"id": "`+id+`",
				"serviceId": "`+serviceID+`",
				"organizationId": "`+orgID+`",
				"privateKey": "-----BEGIN RSA PRIVATE KEY-----MIIEowIBAAKCAQEAh2kjyJRO/rJuwMOECMcR5ZoSCwIsq3205gABA8BNGAZk1CGxxbduIX3peZMdMC5nyQ2gVVHj97Bsx9sYHXs3ihUVzeCth9SEipmrIuktrIzMMbN1Hd0DBMYEVfgxER/F5yYvt3D2RwnydBOg9QSc7GoFa6zzJMfplyVQyCiDgCYN9qf+BHTPk7+x/7ev835ylBsYOFHebN0WsXUlEBWc8nCSJm0Z0hFcmUeLcVSsB2X7iMqGq4VpgSsRfP9sloCX8lNzDo4ujwvTIwsCpHwokQhxCV2avhmjH2hMQXHGlgUfyPm3RjQZ0waSGRRCEuWDhm7xLq9hRzad4cD2auVwhQIDAQABAoIBAD3YwrRBMNdZzgYTBsIvkjgJJ8aJZrepAa+vPsdk1JFtki3ledmxTwbTCIkzrTgtac/Ffn6ZmYKuvPCHXDtS5OoXeU8AGKIaabMYPrcCQ480+6qTqaFLKa7Ldn2Bj3+fwHcz1MV3PbTykR99O53NTpMYVYN5idA50rHrJDtXbcBgdc8KDA79keN6Fv6pggm//Vms7/E9/bgYYY1W8FmjngNHHDgH2EEnERMkyDp2Ng2/2jJIAa8uKPPFxjJAtApQjlLEv82MOcq+Yeq2/VXGRbvgxWlhdEUQcfCf2/ncFEhslau2llcgKTyvcqbaB+jI9Fj8AeXkVyAGVA4X3fuodp0CgYEA0p3oTonVovr+hXELd8hshElAQ+xP4KtPRFH00VPBkB5m3rAALVK+brDAzzZ7W/5X/QTpELcIfIbx5tbFOcQXjuakxVrbof0wZEDGT9KGPSqAAScp7GqYKMSNWUdZ3G5bLisnHT0q7jdWw9se5HGYPEBxYGW8NgbSX+ejsd84ydsCgYEApJax/f+h5RkTUx9TBSRoTz7jMmOls8bbC3IiAUSlVd3wyJX98etsiIsGFh0mrG16Z8Kg0C597woFoUEWTSv7vzN+SLH2CJAt7JzDdlbgmj7EtB0WJwOCPTRUCbQBeMvy3G34NwAlQxGcQqwbo9JMx/8haxqSgCpfrKaD0/CUrR8CgYBgbwqpwzR9Lj0RbkQY8Ty2iS+SqgWc0fM2TewxWA8dZL4nIiDCn8svtWBiwAhVg6xX3kK0c4nAMq1Zy2Z8X4uF05cIAeTkU6AvlvT2IWdzZB096eepJtlKeUxa32+GnUTEa9+55ILelZn1jUOkx1oz5DHFOG+nsRHr9Yye6Zz/1wKBgCAQF0aS6Rf3RZN407R5vjRJ3Pqw/NPD1mIpbsRuegL7RG/fAGSDZ1ZGNv5R2XnXrfPOr4M+u1u4yRX71vtbqSQ7RMuml3ZdmASzGUTRcdm6hplL3UfmYBXKPuDRB0Rf/sTAS41zYs7o/FbkrlHAoyKG6hyyRX3gQ1kf6yh7gosjAoGBAJv0O9x0oY7HZ5QF1PNlDLZUSF/8UtdmyKnm+6VLBmZBaOBJ7MWOJUzzJOHgMRZyzSbH1Z7aBfDXEdXuPihzel/m6TTt/sBv6P8UbVz8cA/uK5BC6nYVwhhu+/wQT/pWNxfqERDKcS7LJ+F3XWkngYA1fBVFmwwDeDWhzgcLQmU6-----END RSA PRIVATE KEY-----",
				"expiresOn": "2019-08-15T17:38:06.322Z",
				"name": "`+serviceName+`",
				"applicationId": "`+applicationID+`",
				"defaultScopes": [
					"openid"
				],
				"scopes": [
					"openid"
				]
			}`)
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
				"total": 1,
				"entry": [
					{
						"id": "`+id+`",
						"serviceId": "`+serviceID+`",
						"organizationId": "`+orgID+`",
						"expiresOn": "2019-08-15T17:38:06.322Z",
						"name": "testservice",
						"applicationId": "`+applicationID+`",
						"defaultScopes": [
							"openid"
						],
						"scopes": [
							"openid"
						]
					}
				]
			}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Service/"+id, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	})

	var r Service
	r.Name = serviceName
	r.Description = serviceDescription
	r.ApplicationID = applicationID

	service, resp, err := client.Services.CreateService(r)
	if ok := assert.Nil(t, err); !ok {
		return
	}
	if ok := assert.NotNil(t, resp); !ok {
		return
	}
	if ok := assert.NotNil(t, service); !ok {
		return
	}
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, serviceName, service.Name)
	assert.Equal(t, applicationID, service.ApplicationID)
	assert.True(t, len(service.PrivateKey) > 0)

	foundService, resp, err := client.Services.GetServiceByID(service.ID)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	if !assert.NotNil(t, foundService) {
		return
	}
	assert.Equal(t, service.ID, foundService.ID)

	services, resp, err := client.Services.GetServicesByApplicationID(applicationID)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, services) {
		return
	}
	assert.Equal(t, 1, len(*services))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ok, resp, err := client.Services.DeleteService(*foundService)
	if !assert.NotNil(t, resp) {
		return
	}
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.True(t, ok)
}

func TestScopes(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	serviceName := "testservice"
	serviceDescription := "Service description"
	applicationID := "b0889958-3762-4427-af07-2d6268c26988"
	id := "2c266886-f918-4223-941d-437cb3cd09e8"

	muxIDM.HandleFunc("/authorize/identity/Service/"+id+"/$scopes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "PUT":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	var r Service
	r.ID = id
	r.Name = serviceName
	r.Description = serviceDescription
	r.ApplicationID = applicationID

	ok, resp, err := client.Services.AddScopes(r, []string{"foo"}, []string{"foo"})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, ok) {
		return
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	ok, resp, err = client.Services.RemoveScopes(r, []string{"foo"}, []string{"foo"})
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, ok) {
		return
	}
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestServicesService_UpdateServiceCertificate(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	var privateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCbn6eqy4c9AyqHTNLUHD17yzSf47uiLdcl+Djo2acsIBO4sihI
XYrPTahSh+4MuwE0CkXN7R8wD5UdhBjjFPPJb8saWdz7wO5zIBpo0XaiBnzy/Bo2
/NfBt3jDgJW1KrzhO2iliNT05fkuV5qHaWqynjHrTX9XW02YWR+lhBwH7QIDAQAB
AoGAIo9QEtPkEp6TYTykTUqANO+pniBa2OdJMjtvl/ZN87Sj7I5T984eN29NhBZ/
vQACgmYNnm+cT3YGbWmahUNv5+FgrBzp+wP+vjioreKc3rbzEZ7Q1WSCgRYPWgjX
cnGbNoW9eIBU4AbIHjd4AZo6/8D+e/3lP2NID+sVw00+xMECQQDOwNd7qWmWixVR
ulNnPCpnvzucg0UMbxJKkM7JCSPItQ3v+YEBmUIS07JFWDAcM6eaNRNS8tEhjiKS
9f+S6nHZAkEAwLEXLrmtAgGrbBqGSF9h2OV1fAqucQbAy0jWkJ1ca2TxWzHSJEji
ba/Kz1gZIrUyRpz2b/Sq8YjOMdor2iZmNQJBAM61iPsvJs0HozoB8u8M/UyaqkvS
Bo/m4L/CRaxDQg6hzhx2NN/XTXMrSloAUKFxl4katYZCFIra44BRCMPnOAkCQB6V
fRp7egxRe4XjB+FA0mT8tpMDlwB1k/AwFW5Kq7qYAP76f2Im91l1h3k1gHfTzrEK
hMrsSJ2hmyiqej9bzOkCQQDGQzjdaCL/8c9DsR20vMwKUCmkUNx69Owb8UJ/BGwe
H4N1gfoHxQbN34YNKoAzQbode9Xv1p9CAPi2v9VlNen2
-----END RSA PRIVATE KEY-----`
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if !assert.NotNil(t, block) {
		return
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if !assert.Nil(t, err) {
		return
	}
	applicationID := "b0889958-3762-4427-af07-2d6268c26988"
	id := "2c266886-f918-4223-941d-437cb3cd09e8"
	orgID := "bda40124-54fa-4967-b2fb-23dcc4e0ad1a"

	muxIDM.HandleFunc("/authorize/identity/Service", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{
				"total": 1,
				"entry": [
					{
						"id": "`+id+`",
						"serviceId": "`+id+`",
						"organizationId": "`+orgID+`",
						"expiresOn": "2019-08-15T17:38:06.322Z",
						"name": "testservice",
						"applicationId": "`+applicationID+`",
						"defaultScopes": [
							"openid"
						],
						"scopes": [
							"openid"
						]
					}
				]
			}`)
		}
	})
	muxIDM.HandleFunc("/authorize/identity/Service/"+id+"/$update-certificate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case "POST":
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{}`)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	var r Service
	r.Name = "name"
	r.Description = "description"
	r.ApplicationID = applicationID
	r.ID = id

	service, resp, err := client.Services.UpdateServiceCertificate(r, privateKey)
	if ok := assert.Nil(t, err); !ok {
		return
	}
	if ok := assert.NotNil(t, resp); !ok {
		return
	}
	if ok := assert.NotNil(t, service); !ok {
		return
	}
}

func TestFixIAMPEM(t *testing.T) {
	var generatedIAMPEM = `-----BEGIN RSA PRIVATE KEY-----MIIEogIBAAKCAQEAg2ASjNf6XT/wPqjJlrjCCeD7aSC2JkFXkdlLIGUgrwYJdQcmP8cIfzsJmAa0Fba5TBBQIgEZiTutvyoUonso4Axw3itNyHlq/avRxMT7lQl/wp+xzYBGg0OHr8b1VaML031KQWJB7yn4OPeJebd2c16YgCSd1BWC0ysCB1bjbMj6f5XrHvkm1sbCLhXP4Nsz2oQAKvNYnz9B7qoI5ENlU6CLPAf22Ot885v/Ezg3QM8oSz7vPGTP5iEeGbHm3PppXEhUg47sWnc1zPHTvkEBkgSKuNKV/GPnrspynGCxfFaaFJBLxKKDiIgcn63Ui/WP6YJuk8C4DC0N8cT1fpcKQQIDAQABAoIBADe8XqJaTuGY6rfqq/vBzrCaPo9PTXfVaO17np/Tx5X5nbGEovmT20r/JQbWcTh1GQ/oK9cBWHAcXY4osDBux9yl6oyXGpzexkFwz8IHQCTursFHNwvUTBYMZFhytErsrWvAT+yFXyu0RpcEjdKRvzMbAt/6qOr9ZU9pg6bRehMAX/8XRVgbF4vAAGSIW9r2fDhoy/9w4NfI9YS9j3DSplCDEcKKt/L006d7egr1EkyK8jqGmUDHNRgusnlj1EkaeLY8KFxSaVFgzZsap2Q961d8GErHCktyQS3rdPR/R9VhyJ0fwnPZs3fdGHgz86NvByEqRXdu3S4fg09JUeOdgAECgYEA709J5YznPVfJQaXQGBug0yZnAJKxAiCZaGPu9rLzGhc+opFs1RPBXgZGVecLPTGUOASukYczt++A/AkNaugPrIPwR8odf6I8g9JUvh43ixM/AWQxRyxPpI6hmu07+OetTToTpUhNL7/49pAGIZmwL+q89fx9ZaQJk9ryblloMoECgYEAjImwMXoNpKLUSfJL+aoGyZp4g+uFn0zNInK6AFhoNDmuYdtKMfWpV35zJaIJkoSFFk6C7OfT5f3h/wC2CUx4spmdohAAySm3qopRop4cGjAu6C0eVU8GQuctybnupxSqM9tBw6h0DJP/qK63dqAw39PjpZSbtNedQgvN/LIPd8ECgYArGVGKuPywbI3/d7hpjxhAzurADKs6YdqfFVhlL8dejI8oQh6ipWMu5txl/i1lAPfbuuQhq4Lh/o1xVsQG2/JutHhixYDYisRxTUEb2F9Sec3gF0ZUtiNF1LZAK7egVUE95lbBbqsQDXfuF/QyuvYjOxh7CNUgqBFdgpIk049TgQKBgFCX2TnWY6KJWjJ5x8YKElBkzD3iv8Ml2PGrVy1BhaiFd5JYnkjeN1mcg04OiWzDoEq1hXF0VNks6gBBmSLl/sV1QKeig8DYsqITJlOaVkVSYmP1E9MaOV3+kwx/GvQn7p9EWodov/xClzdV0pIttTl2mrLm8RAZIe6xNFR4jFRBAoGAZO51gVNvHqfFlApJKKw6wJcFlQ0a4wzOY8ZEsRtdK+iLcZHtFdWDjiyzB/lx8DL0l6dlsVP/kH7HgEa1rPYfDNQW6V4QIdGxT8KcVOJXMUQabE/yhoVSpL05e/aX5IsAfJ7sbMe2qcHz+pivU17tsSjO+xTUyJf0Fv5kOfkQi2s=-----END RSA PRIVATE KEY-----`
	var fixedIAMPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAg2ASjNf6XT/wPqjJlrjCCeD7aSC2JkFXkdlLIGUgrwYJdQcmP8cIfzsJmAa0Fba5TBBQIgEZiTutvyoUonso4Axw3itNyHlq/avRxMT7lQl/wp+xzYBGg0OHr8b1VaML031KQWJB7yn4OPeJebd2c16YgCSd1BWC0ysCB1bjbMj6f5XrHvkm1sbCLhXP4Nsz2oQAKvNYnz9B7qoI5ENlU6CLPAf22Ot885v/Ezg3QM8oSz7vPGTP5iEeGbHm3PppXEhUg47sWnc1zPHTvkEBkgSKuNKV/GPnrspynGCxfFaaFJBLxKKDiIgcn63Ui/WP6YJuk8C4DC0N8cT1fpcKQQIDAQABAoIBADe8XqJaTuGY6rfqq/vBzrCaPo9PTXfVaO17np/Tx5X5nbGEovmT20r/JQbWcTh1GQ/oK9cBWHAcXY4osDBux9yl6oyXGpzexkFwz8IHQCTursFHNwvUTBYMZFhytErsrWvAT+yFXyu0RpcEjdKRvzMbAt/6qOr9ZU9pg6bRehMAX/8XRVgbF4vAAGSIW9r2fDhoy/9w4NfI9YS9j3DSplCDEcKKt/L006d7egr1EkyK8jqGmUDHNRgusnlj1EkaeLY8KFxSaVFgzZsap2Q961d8GErHCktyQS3rdPR/R9VhyJ0fwnPZs3fdGHgz86NvByEqRXdu3S4fg09JUeOdgAECgYEA709J5YznPVfJQaXQGBug0yZnAJKxAiCZaGPu9rLzGhc+opFs1RPBXgZGVecLPTGUOASukYczt++A/AkNaugPrIPwR8odf6I8g9JUvh43ixM/AWQxRyxPpI6hmu07+OetTToTpUhNL7/49pAGIZmwL+q89fx9ZaQJk9ryblloMoECgYEAjImwMXoNpKLUSfJL+aoGyZp4g+uFn0zNInK6AFhoNDmuYdtKMfWpV35zJaIJkoSFFk6C7OfT5f3h/wC2CUx4spmdohAAySm3qopRop4cGjAu6C0eVU8GQuctybnupxSqM9tBw6h0DJP/qK63dqAw39PjpZSbtNedQgvN/LIPd8ECgYArGVGKuPywbI3/d7hpjxhAzurADKs6YdqfFVhlL8dejI8oQh6ipWMu5txl/i1lAPfbuuQhq4Lh/o1xVsQG2/JutHhixYDYisRxTUEb2F9Sec3gF0ZUtiNF1LZAK7egVUE95lbBbqsQDXfuF/QyuvYjOxh7CNUgqBFdgpIk049TgQKBgFCX2TnWY6KJWjJ5x8YKElBkzD3iv8Ml2PGrVy1BhaiFd5JYnkjeN1mcg04OiWzDoEq1hXF0VNks6gBBmSLl/sV1QKeig8DYsqITJlOaVkVSYmP1E9MaOV3+kwx/GvQn7p9EWodov/xClzdV0pIttTl2mrLm8RAZIe6xNFR4jFRBAoGAZO51gVNvHqfFlApJKKw6wJcFlQ0a4wzOY8ZEsRtdK+iLcZHtFdWDjiyzB/lx8DL0l6dlsVP/kH7HgEa1rPYfDNQW6V4QIdGxT8KcVOJXMUQabE/yhoVSpL05e/aX5IsAfJ7sbMe2qcHz+pivU17tsSjO+xTUyJf0Fv5kOfkQi2s=
-----END RSA PRIVATE KEY-----
`

	fixed := FixPEM(generatedIAMPEM)
	assert.Equal(t, fixedIAMPEM, fixed)
	// A second run should be a NOOP
	assert.Equal(t, fixedIAMPEM, FixPEM(fixed))

	// Test decoding
	block, _ := pem.Decode([]byte(generatedIAMPEM))
	assert.Nil(t, block)

	block, _ = pem.Decode([]byte(fixed))
	assert.NotNil(t, block)
}
