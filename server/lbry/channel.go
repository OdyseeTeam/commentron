package lbry

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"strconv"
	"time"

	"github.com/lbryio/commentron/helper"

	"github.com/sirupsen/logrus"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/lbryio/lbry.go/v2/schema/keys"
)

var testChannel *Channel

func setSerializedTestChannel(serializedChannel string) error {
	channelBytes := base58.Decode(serializedChannel)
	channel := &Channel{}
	err := json.Unmarshal(channelBytes, channel)
	if err != nil {
		return errors.Err(err)
	}
	testChannel = channel
	testChannel.Keys()
	return nil
}

// Channel structure representing an exported LBRY channel
type Channel struct {
	Name              string `json:"name"`
	ChannelID         string `json:"channel_id"`
	HoldingAddress    string `json:"holding_address"`
	HoldingPublicKey  string `json:"holding_public_key"`
	SigningPrivateKey string `json:"signing_private_key"`
}

// ImportChannel creates a Channel representation from the serialization
func ImportChannel(serializedChannel string) *Channel {
	channelBytes := base58.Decode(serializedChannel)
	channel := &Channel{}
	err := json.Unmarshal(channelBytes, channel)
	if err != nil {
		logrus.Panic(err)
	}
	return channel
}

func newChannel(name, claimID string, private *btcec.PrivateKey) (*Channel, error) {
	b := bytes.NewBuffer(nil)
	derBytes, err := keys.PrivateKeyToDER(private)
	if err != nil {
		return nil, errors.Err(err)
	}
	err = pem.Encode(b, &pem.Block{Type: "PRIVATE KEY", Bytes: derBytes})
	if err != nil {
		return nil, errors.Err(err)
	}
	channel := Channel{
		Name:              name,
		ChannelID:         claimID,
		HoldingAddress:    "",
		HoldingPublicKey:  "",
		SigningPrivateKey: string(b.Bytes()),
	}

	private2, _ := channel.Keys()
	if !private2.ToECDSA().Equal(private.ToECDSA()) {
		return nil, errors.Err("private keys don't match")
	}
	return &channel, nil
}

// Keys will return the private and public key of a channel
func (l *Channel) Keys() (*btcec.PrivateKey, *btcec.PublicKey) {
	return keys.ExtractKeyFromPem(l.SigningPrivateKey)
}

// Sign will sign the data with the channels private key returning the signature and a timestamp it was signed. This is
// used in the signature itself so it is required input to verify the signature. This ensures someone cannot fudge a
// timestamp of a signature and allows for expiration times of signatures.
func (l *Channel) Sign(data []byte) (string, string, error) {
	private, _ := l.Keys()
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	hash := sha256.Sum256(helper.CreateDigest(
		[]byte(timestamp),
		unhelixifyAndReverse(l.ChannelID),
		data))
	hashBytes := make([]byte, len(hash))
	for i, b := range hash {
		hashBytes[i] = b
	}
	sig, err := private.Sign(hashBytes[:])
	if err != nil {
		return "", "", errors.Err(err)
	}

	valid := ecdsa.Verify(private.PubKey().ToECDSA(), hashBytes[:], sig.R, sig.S)
	if !valid {
		return "", "", errors.Err("sig not valid")
	}
	keysSig := keys.Signature{Signature: *sig}
	sigBytes, err := keysSig.LBRYSDKEncode()
	if err != nil {
		return "", "", errors.Err(err)
	}
	return hex.EncodeToString(sigBytes), timestamp, nil
}

// Claim will return the metadata for the channel.
func (l *Channel) Claim() (*jsonrpc.Claim, error) {
	return nil, nil
}
