package lbry

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"math/big"

	"github.com/lbryio/commentron/util"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/btcsuite/btcd/btcec"
)

func ValidateSignature(channelClaimId, signature, signingTS, data string) error {
	channel, err := GetChannelClaim(channelClaimId)
	if err != nil {
		return errors.Err(err)
	}
	pk := channel.Value.GetChannel().GetPublicKey()
	return validateSignature(channelClaimId, signature, signingTS, data, pk)

}

func validateSignature(channelClaimId, signature, signingTS, data string, pubkey []byte) error {
	publicKey, err := getPublicKeyFromBytes(pubkey)
	if err != nil {
		return errors.Err(err)
	}
	injest := sha256.Sum256(
		util.CreateDigest(
			[]byte(signingTS),
			unhelixifyAndReverse(channelClaimId),
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
	asn1.Unmarshal(pubKeyBytes, &PKInfo)
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
