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

package status

import (
	"fmt"
	ht "net/http"
)

// HTTPError wraps an error from an HTTP request.
func HTTPError(resp *ht.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode >= ht.StatusInternalServerError {
		return fmt.Errorf("%w: %s", ErrServerError, resp.Status)
	}
	if resp.StatusCode >= ht.StatusBadRequest {
		return fmt.Errorf("%w: %s", ErrClientError, resp.Status)
	}
	return nil
}
