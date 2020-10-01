package query

import (
	"testing"

	"github.com/lbryio/lbrytv/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc"
)

func TestQueryParamsAsMap(t *testing.T) {
	q, err := NewQuery(jsonrpc.NewRequest("version"), "")
	require.NoError(t, err)
	assert.Nil(t, q.ParamsAsMap())

	q, err = NewQuery(jsonrpc.NewRequest("resolve", map[string]interface{}{"urls": "what"}), "")
	require.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"urls": "what"}, q.ParamsAsMap())

	q, err = NewQuery(jsonrpc.NewRequest("account_balance"), "123")
	require.NoError(t, err, errors.Unwrap(err))
	assert.Equal(t, map[string]interface{}{"wallet_id": "123"}, q.ParamsAsMap())

	searchParams := map[string]interface{}{
		"any_tags": []interface{}{
			"art", "automotive", "blockchain", "comedy", "economics", "education",
			"gaming", "music", "news", "science", "sports", "technology",
		},
	}
	q, err = NewQuery(jsonrpc.NewRequest("claim_search", searchParams), "")
	require.NoError(t, err)
	assert.Equal(t, searchParams, q.ParamsAsMap())
}

func TestQueryIsAuthenticated(t *testing.T) {
	q, err := NewQuery(jsonrpc.NewRequest("resolve"), "12345")
	require.NoError(t, err)
	assert.True(t, q.IsAuthenticated())

	q, err = NewQuery(jsonrpc.NewRequest("resolve"), "")
	require.NoError(t, err)
	assert.False(t, q.IsAuthenticated())
}

func TestMethodRequiresWallet(t *testing.T) {
	for _, m := range walletSpecificMethods {
		if methodInList(m, relaxedMethods) {
			assert.False(t, MethodRequiresWallet(m, nil), m)
		} else {
			assert.True(t, MethodRequiresWallet(m, nil), m)
		}
	}
	assert.True(t, MethodRequiresWallet(MethodCommentReactList, map[string]interface{}{
		ParamChannelID: "f77ff625ccf2401d6f47f9bfd7d17b8bb370bda7"}))
	assert.False(t, MethodRequiresWallet(MethodCommentReactList, nil))
}

func TestMethodAcceptsWallet(t *testing.T) {
	for _, m := range walletSpecificMethods {
		assert.True(t, MethodAcceptsWallet(m), m)
	}
}
