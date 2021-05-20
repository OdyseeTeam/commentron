package lbry

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"math/big"
	"strconv"

	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/lbryio/lbry.go/v2/schema/keys"

	"github.com/btcsuite/btcd/btcec"
)

// ValidateSignatures determines if signatures should be validated or not ( not used yet)
var ValidateSignatures bool

// ValidateSignature validates the signature was signed by the channel reference.
func ValidateSignature(channelClaimID, signature, signingTS, data string) error {
	channel, err := SDK.GetClaim(channelClaimID)
	if err != nil {
		return errors.Err(err)
	}
	amount, err := strconv.ParseFloat(channel.Amount, 64)
	if err != nil {
		return errors.Err(err)
	}
	supports, err := strconv.ParseFloat(channel.Meta.SupportAmount, 64)
	if err != nil {
		return errors.Err(err)
	}
	if amount+supports < 0.001 {
		return errors.Err("validation is disallowed for non controlling channels")
	}
	pk := channel.Value.GetChannel().GetPublicKey()
	return validateSignature(channelClaimID, signature, signingTS, data, pk)

}

// ValidateSignatureFromClaim validates the signature was signed by the channel reference.
func ValidateSignatureFromClaim(channel *jsonrpc.Claim, signature, signingTS, data string) error {
	if channel == nil {
		return errors.Err("no channel to validate")
	}
	if channel.Value.GetChannel() == nil {
		return errors.Err("no channel for public key")
	}
	pk := channel.Value.GetChannel().GetPublicKey()
	return validateSignature(channel.ClaimID, signature, signingTS, data, pk)

}

// encodePrivateKey encodes an ECDSA private key to PEM format.
func encodePrivateKey(key *btcec.PrivateKey) ([]byte, error) {
	derPrivKey, err := keys.PrivateKeyToDER(key)
	if err != nil {
		return nil, err
	}

	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derPrivKey,
	}

	return pem.EncodeToMemory(keyBlock), nil
}

func validateSignature(channelClaimID, signature, signingTS, data string, pubkey []byte) error {
	publicKey, err := getPublicKeyFromBytes(pubkey)
	if err != nil {
		return errors.Err(err)
	}
	injest := sha256.Sum256(
		helper.CreateDigest(
			[]byte(signingTS),
			unhelixifyAndReverse(channelClaimID),
			[]byte(data),
		))
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return errors.Err(err)
	}
	signatureBytes := [64]byte{}
	for i, b := range sig {
		signatureBytes[i] = b
	}
	sigValid := isSignatureValid(signatureBytes, publicKey, injest[:])
	if !sigValid {
		return errors.Err("could not validate the signature")
	}
	return nil
}

func isSignatureValid(signature [64]byte, publicKey *btcec.PublicKey, injest []byte) bool {

	R := &big.Int{}
	S := &big.Int{}
	R.SetBytes(signature[:32])
	S.SetBytes(signature[32:])
	return ecdsa.Verify(publicKey.ToECDSA(), injest, R, S)
}

// rev reverses a byte slice. useful for switching endian-ness
func reverseBytes(b []byte) []byte {
	r := make([]byte, len(b))
	for left, right := 0, len(b)-1; left < right; left, right = left+1, right-1 {
		r[left], r[right] = b[right], b[left]
	}
	return r
}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

func getPublicKeyFromBytes(pubKeyBytes []byte) (*btcec.PublicKey, error) {
	PKInfo := publicKeyInfo{}
	_, err := asn1.Unmarshal(pubKeyBytes, &PKInfo)
	if err != nil {
		return nil, errors.Err(err)
	}
	pubkeyBytes1 := PKInfo.PublicKey.Bytes
	return btcec.ParsePubKey(pubkeyBytes1, btcec.S256())
}

func unhelixifyAndReverse(claimID string) []byte {
	b, err := hex.DecodeString(claimID)
	if err != nil {
		return nil
	}
	return reverseBytes(b)
}
