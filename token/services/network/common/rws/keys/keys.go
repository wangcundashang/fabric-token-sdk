/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package keys

import (
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
)

const (
	minUnicodeRuneValue   = 0 // U+0000
	compositeKeyNamespace = "\x00"
	MaxUnicodeRuneValue   = utf8.MaxRune // U+10FFFF - maximum (and unallocated) code point

	numComponentsInKey = 2 // 2 components: txid, index, excluding TokenKeyPrefix

	TokenKeyPrefix          = "ztoken"
	TokenSetupKeyPrefix     = "setup"
	TokenSetupHashKeyPrefix = "setup.hash"

	TokenRequestKeyPrefix  = "token_request"
	SerialNumber           = "sn"
	IssueActionMetadata    = "iam"
	TransferActionMetadata = "tam"
)

func SplitCompositeKey(compositeKey string) (string, []string, error) {
	componentIndex := 1
	var components []string
	for i := 1; i < len(compositeKey); i++ {
		if compositeKey[i] == minUnicodeRuneValue {
			components = append(components, compositeKey[componentIndex:i])
			componentIndex = i + 1
		}
	}
	// there is an extra tokenIdPrefix component in the beginning, trim it off
	if len(components) < numComponentsInKey+1 {
		return "", nil, errors.Errorf("invalid composite key - not enough components found in key '%s', [%d][%v]", compositeKey, len(components), components)
	}
	return components[0], components[1:], nil
}

// CreateTokenKey Creates a rwset key for an individual output in a token transaction, as a function of
// the token owner, transaction ID, and index of the output
func CreateTokenKey(txID string, index uint64) (string, error) {
	return CreateCompositeKey(TokenKeyPrefix, []string{txID, strconv.FormatUint(index, 10)})
}

func CreateSNKey(sn string) (string, error) {
	return CreateCompositeKey(TokenKeyPrefix, []string{SerialNumber, sn})
}

func CreateSetupKey() (string, error) {
	return CreateCompositeKey(TokenKeyPrefix, []string{TokenSetupKeyPrefix})
}

func CreateSetupHashKey() (string, error) {
	return CreateCompositeKey(TokenKeyPrefix, []string{TokenSetupHashKeyPrefix})
}

func CreateTokenRequestKey(txID string) (string, error) {
	return CreateCompositeKey(TokenKeyPrefix, []string{TokenRequestKeyPrefix, txID})
}

func CreateIssueActionMetadataKey(hash string) (string, error) {
	return CreateCompositeKey(TokenKeyPrefix, []string{IssueActionMetadata, hash})
}

// CreateTransferActionMetadataKey returns the transfer action metadata key built from the passed
// transaction id, subkey, and index. Index is used to make sure the key is unique with the respect to the
// token request this key appears.
func CreateTransferActionMetadataKey(subKey string) (string, error) {
	return CreateCompositeKey(TokenKeyPrefix, []string{TransferActionMetadata, subKey})
}

func GetTransferMetadataSubKey(k string) (string, error) {
	prefix, components, err := SplitCompositeKey(k)
	if err != nil {
		return "", errors.Wrapf(err, "failed to split composite key [%s]", k)
	}
	if len(components) != 2 {
		return "", errors.Wrapf(err, "key [%s] should contain 2 components, got [%d]", k, len(components))
	}
	if prefix != TokenKeyPrefix {
		return "", errors.Errorf("key [%s] doesn not contain the token key prefix", k)
	}
	if components[0] != TransferActionMetadata {
		return "", errors.Errorf("key [%s] doesn not contain the token transfer action medatata prefix", k)
	}
	return components[1], nil
}

// CreateCompositeKey and its related functions and consts copied from core/chaincode/shim/chaincode.go
func CreateCompositeKey(objectType string, attributes []string) (string, error) {
	if err := ValidateCompositeKeyAttribute(objectType); err != nil {
		return "", err
	}
	ck := compositeKeyNamespace + objectType + string(rune(minUnicodeRuneValue))
	for _, att := range attributes {
		if err := ValidateCompositeKeyAttribute(att); err != nil {
			return "", err
		}
		ck += att + string(rune(minUnicodeRuneValue))
	}
	return ck, nil
}

func ValidateCompositeKeyAttribute(str string) error {
	if !utf8.ValidString(str) {
		return errors.Errorf("not a valid utf8 string: [%x]", str)
	}
	for index, runeValue := range str {
		if runeValue == minUnicodeRuneValue || runeValue == MaxUnicodeRuneValue {
			return errors.Errorf(`input contain unicode %#U starting at position [%d]. %#U and %#U are not allowed in the input attribute of a composite key`,
				runeValue, index, minUnicodeRuneValue, MaxUnicodeRuneValue)
		}
	}
	return nil
}

type Translator struct {
}

func (t *Translator) CreateTokenRequestKey(id string) (string, error) {
	return CreateTokenRequestKey(id)
}

func (t *Translator) CreateSetupKey() (string, error) {
	return CreateSetupKey()
}

func (t *Translator) CreateSetupHashKey() (string, error) {
	return CreateSetupHashKey()
}

func (t *Translator) CreateTokenKey(id string, index uint64) (string, error) {
	return CreateTokenKey(id, index)
}

func (t *Translator) GetTransferMetadataSubKey(k string) (string, error) {
	return GetTransferMetadataSubKey(k)
}

func (t *Translator) CreateSNKey(id string) (string, error) {
	return CreateSNKey(id)
}

func (t *Translator) CreateIssueActionMetadataKey(key string) (string, error) {
	return CreateIssueActionMetadataKey(key)
}

func (t *Translator) CreateTransferActionMetadataKey(key string) (string, error) {
	return CreateTransferActionMetadataKey(key)
}
