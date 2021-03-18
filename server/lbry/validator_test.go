package lbry

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"github.com/lbryio/commentron/helper"

	"github.com/lbryio/lbry.go/v2/schema/keys"

	"github.com/btcsuite/btcd/btcec"
)

func TestValidateSignature1(t *testing.T) {
	channelClaimID := "7fadfe1d0dce928350137a13497b6fc36627cf45"
	pubkeyHex := "3056301006072a8648ce3d020106052b8104000a03420004e0743cfa62857d1d7bda9ca6ba0ec3325902866e6442f51a9da2b143bc0ba40cda532e483e1a8a48c84b4b9dc16a117b2f9763d518db50d8fed2b818937ef8b1"
	signature := "fe35046bd949fc89037d64ac3558fea859022a166558b459b6883acafa15ca9ec567ca23e7b4ae19e4dbc3f92aac30a132315db7abcb03c15c61662fb9f49458"
	signingTS := "1582846386"
	data := "nicee"
	pk, err := hex.DecodeString(pubkeyHex)
	if err != nil {
		t.Fatal(err)
	}
	err = validateSignature(channelClaimID, signature, signingTS, data, pk)
	if err != nil {
		t.Error(err)
	}
}

func TestValidateSignature2(t *testing.T) {
	channelClaimID := "6dab3a207b6551b9c4a0c782e22963d2b444d609"
	pubkeyHex := "3056301006072a8648ce3d020106052b8104000a0342000428f5f61f7e051aa7c9f6f1e9802773ac4d77a0ffcc4f282252c8c889e9c225cbb5afa5bc12f4c2c5017513a767a138123cf0e3919b7927c9f1249750e7f688f2"
	signature := "cae3b6ca34c141bd0a3b20355c5ed7c5f718c45a764194629ab612d48448061dd42ae3ccf49848d529421265c9ee348c60233d0c76feafbb9ad4221aee9c9072"
	signingTS := "1591846880"
	data := "thank you"
	pk, err := hex.DecodeString(pubkeyHex)
	if err != nil {
		t.Fatal(err)
	}
	err = validateSignature(channelClaimID, signature, signingTS, data, pk)
	if err != nil {
		t.Error(err)
	}
}

func TestCommentSignAndVerify(t *testing.T) {
	channelClaimID := "9cb713f01bf247a0e03170b5ed00d5161340c486"
	private, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		t.Error(err)
	}

	comment := "sign this shit"
	strconv.FormatInt(time.Now().Unix(), 10)
	digest := sha256.Sum256(helper.CreateDigest(
		[]byte(strconv.FormatInt(time.Now().Unix(), 10)),
		unhelixifyAndReverse(channelClaimID),
		[]byte(comment)))
	sig, err := private.Sign(digest[:])
	valid := ecdsa.Verify(private.PubKey().ToECDSA(), digest[:], sig.R, sig.S)
	if !valid {
		t.Error("sig not valid")
	}
}

func TestCommentSignAndVerifyNew(t *testing.T) {
	private, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		t.Fatal(err)
	}
	pubkeyBytes, err := keys.PublicKeyToDER(private.PubKey())
	if err != nil {
		t.Fatal(err)
	}
	channel, err := newChannel("@MyTestChannel", "9cb713f01bf247a0e03170b5ed00d5161340c486", private)
	if err != nil {
		t.Fatal(err)
	}

	comment := "sign this shit"
	signature, timestamp, err := channel.Sign([]byte(comment))
	err = validateSignature(channel.ChannelID, signature, timestamp, comment, pubkeyBytes)
	if err != nil {
		t.Error(err)
	}
}

func TestSignatures(t *testing.T) {
	private, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		t.Error(err)
	}

	stuff := "sign this shit"
	digest := sha256.Sum256([]byte(stuff))
	sig, err := private.Sign(digest[:])
	if err != nil {
		t.Error(err)
	}

	valid := ecdsa.Verify(private.PubKey().ToECDSA(), digest[:], sig.R, sig.S)
	if !valid {
		t.Error("sig not valid")
	}
}
