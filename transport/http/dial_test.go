// Copyright 2023 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"context"
	"fmt"
	"testing"

	"go.opencensus.io/plugin/ochttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func TestNewClient(t *testing.T) {
	client, endpoint, err := NewClient(context.Background())

	if err != nil {
		t.Fatalf("unable to create client: %v", err)
	}
	if client == nil {
		t.Fatalf("client is nil")
	}
	if endpoint != "" {
		t.Errorf("got: %s, want: ''", endpoint)
	}
	if got, want := fmt.Sprintf("%T", client.Transport), "*oauth2.Transport"; got != want {
		t.Fatalf("got %s, want: %s", got, want)
	}
	t1 := client.Transport.(*oauth2.Transport)
	if got, want := fmt.Sprintf("%T", t1.Base), "*ochttp.Transport"; got != want {
		t.Fatalf("got %s, want: %s", got, want)
	}
	t2 := t1.Base.(*ochttp.Transport)
	if got, want := fmt.Sprintf("%T", t2.Base), "*otelhttp.Transport"; got != want {
		t.Fatalf("got %s, want: %s", got, want)
	}
}

func TestNewClient_MismatchedUniverseDomainCreds(t *testing.T) {
	rootTokenScope := "https://www.googleapis.com/auth/cloud-platform"
	universeDomain := "example.com"
	universeDomainDefault := "googleapis.com"
	creds := &google.Credentials{} // universeDomainDefault
	wantErr := errUniverseNotMatch(universeDomain, universeDomainDefault)
	_, _, err := NewClient(context.Background(), option.WithUniverseDomain(universeDomain),
		option.WithCredentials(creds), option.WithScopes(rootTokenScope))

	if err.Error() != wantErr.Error() {
		t.Fatalf("got: %v, want: %v", err, wantErr)
	}
}
