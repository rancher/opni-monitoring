package tokens

import (
	"encoding/hex"

	"github.com/rancher/opni-monitoring/pkg/core"
)

func FromBootstrapToken(t *core.BootstrapToken) (*Token, error) {
	tokenID := t.GetTokenID()
	tokenSecret := t.GetSecret()
	token := &Token{
		ID:     make([]byte, hex.DecodedLen(len(tokenID))),
		Secret: make([]byte, hex.DecodedLen(len(tokenSecret))),
	}
	decodedID, err := hex.DecodeString(tokenID)
	if err != nil {
		return nil, err
	}
	decodedSecret, err := hex.DecodeString(tokenSecret)
	if err != nil {
		return nil, err
	}
	copy(token.ID, decodedID)
	copy(token.Secret, decodedSecret)
	return token, nil
}

func (t *Token) ToBootstrapToken() *core.BootstrapToken {
	return &core.BootstrapToken{
		TokenID: t.HexID(),
		Secret:  t.HexSecret(),
		Metadata: &core.BootstrapTokenMetadata{
			Labels: map[string]string{},
		},
	}
}
