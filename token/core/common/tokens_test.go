/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokensService_NewUpgradeChallenge(t *testing.T) {
	ts := NewTokensService()
	challange, err := ts.NewUpgradeChallenge()
	assert.NoError(t, err)
	assert.Len(t, challange, ChallengeSize)
}

func TestTokensService_GenUpgradeProof(t *testing.T) {
	ts := NewTokensService()
	res, err := ts.GenUpgradeProof(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "not supported")
}

func TestTokensService_CheckUpgradeProof(t *testing.T) {
	ts := NewTokensService()
	res, err := ts.CheckUpgradeProof(nil, nil, nil)
	assert.Error(t, err)
	assert.False(t, res)
	assert.EqualError(t, err, "not supported")
}
