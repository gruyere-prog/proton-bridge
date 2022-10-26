// Copyright (c) 2022 Proton AG
//
// This file is part of Proton Mail Bridge.
//
// Proton Mail Bridge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Proton Mail Bridge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Proton Mail Bridge.  If not, see <https://www.gnu.org/licenses/>.

package user

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"gitlab.protontech.ch/go/liteapi"
)

func withAddrKR(apiUser liteapi.User, apiAddr liteapi.Address, keyPass []byte, fn func(userKR, addrKR *crypto.KeyRing) error) error {
	userKR, err := apiUser.Keys.Unlock(keyPass, nil)
	if err != nil {
		return fmt.Errorf("failed to unlock user keys: %w", err)
	}
	defer userKR.ClearPrivateParams()

	addrKR, err := apiAddr.Keys.Unlock(keyPass, userKR)
	if err != nil {
		return fmt.Errorf("failed to unlock address keys: %w", err)
	}
	defer addrKR.ClearPrivateParams()

	return fn(userKR, addrKR)
}

func withAddrKRs(apiUser liteapi.User, apiAddr map[string]liteapi.Address, keyPass []byte, fn func(*crypto.KeyRing, map[string]*crypto.KeyRing) error) error {
	userKR, err := apiUser.Keys.Unlock(keyPass, nil)
	if err != nil {
		return fmt.Errorf("failed to unlock user keys: %w", err)
	}
	defer userKR.ClearPrivateParams()

	addrKRs := make(map[string]*crypto.KeyRing, len(apiAddr))

	for addrID, apiAddr := range apiAddr {
		addrKR, err := apiAddr.Keys.Unlock(keyPass, userKR)
		if err != nil {
			return fmt.Errorf("failed to unlock address keys: %w", err)
		}
		defer addrKR.ClearPrivateParams()

		addrKRs[addrID] = addrKR
	}

	return fn(userKR, addrKRs)
}
