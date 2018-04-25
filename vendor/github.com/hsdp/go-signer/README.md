# Go HSDP Signer

This package implements the HSDP API signing algorithm.
You can sign a http.Request instances 

## Usage

```go

import (
  "github.com/hsdp/go-signer"
  "net/http"
)

func signFilter(req *http.Request, sharedKey, secretKey string) (*http.Request, error) {
    s, err := signer.New(sharedKey, secretKey)
    if err != nil {
        return nil, err
    }
    s.SignRequest(req)
    return req, nil
}

```
## License

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at <http://www.apache.org/licenses/LICENSE-2.0>

