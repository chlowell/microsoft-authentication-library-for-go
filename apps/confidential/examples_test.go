// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

// rename this package to get the user code experience here
package confidential_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

// users can do stuff like this today, dunno why they would want to
func MyCustomOption() confidential.AcquireTokenByAuthCodeOption {
	return func(a *confidential.AcquireTokenByAuthCodeOptions) {
		a.Challenge = "challenge"
		a.TenantID = "tenant"
	}
}

func ExampleNewCredFromCert_pem() {
	b, err := ioutil.ReadFile("key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// This extracts our public certificates and private key from the PEM file.
	// The private key must be in PKCS8 format. If it is encrypted, the second argument
	// must be password to decode.
	certs, priv, err := confidential.CertFromPEM(b, "")
	if err != nil {
		log.Fatal(err)
	}

	// PEM files can have multiple certs. This is usually for certificate chaining where roots
	// sign to leafs. Useful for TLS, not for this use case.
	if len(certs) > 1 {
		log.Fatal("too many certificates in PEM file")
	}

	cred := confidential.NewCredFromCert(certs[0], priv)
	fmt.Println(cred) // Simply here so cred is used, otherwise won't compile.

	client, err := confidential.New("clientId", cred)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	client.AcquireTokenByAuthCode(ctx, "", "", []string{},
		confidential.WithChallenge("challenge"),
		confidential.WithTenantID("tenant"),
		MyCustomOption(),
	)
}
