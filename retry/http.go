// Copyright 2023 Scott M. Long
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import (
	"context"
	ht "net/http"

	"github.com/neuralnorthwest/mu/status"
)

// HTTPGet is a convenience function that performs an HTTP GET request with
// the given URL and returns the response body. It retries the request according
// to the provided strategy, until the context is canceled or the request
// succeeds. "Success" is defined as any response with a non-5xx status code.
// The caller is responsible for closing the response body.
func HTTPGet(ctx context.Context, url string, strategy Strategy, opts ...DoOption) (*ht.Response, error) {
	var resp *ht.Response
	err := Do(ctx, strategy, func() error {
		resp, err := ht.Get(url)
		return status.HTTPError(resp, err)
	}, opts...)
	return resp, err
}
