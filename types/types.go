// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2026 Jelius Basumatary

package types

type ErrorResp struct {
    Code    uint16 `json:"code"`
    Message string `json:"message"`
}

