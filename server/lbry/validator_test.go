package lbry

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"strconv"
	"testing"
	"time"

	"github.com/lbryio/commentron/util"

	"github.com/btcsuite/btcd/btcec"
)

func TestValidateSignature1(t *testing.T) {
	channelClaimID := "7fadfe1d0dce928350137a13497b6fc36627cf45"
	signature := "fe35046bd949fc89037d64ac3558fea859022a166558b459b6883acafa15ca9ec567ca23e7b4ae19e4dbc3f92aac30a132315db7abcb03c15c61662fb9f49458"
	signingTS := "1582846386"
	data := "nicee"
	err := ValidateSignature(channelClaimID, signature, signingTS, data)
	if err != nil {
		t.Error(err)
	}
}

func TestValidateSignature2(t *testing.T) {
	channelClaimID := "6dab3a207b6551b9c4a0c782e22963d2b444d609"
	signature := "cae3b6ca34c141bd0a3b20355c5ed7c5f718c45a764194629ab612d48448061dd42ae3ccf49848d529421265c9ee348c60233d0c76feafbb9ad4221aee9c9072"
	signingTS := "1591846880"
	data := "thank you"
	err := ValidateSignature(channelClaimID, signature, signingTS, data)
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
	digest := sha256.Sum256(util.CreateDigest(
		unhelixifyAndReverse(channelClaimID),
		[]byte(comment),
		[]byte(strconv.FormatInt(time.Now().Unix(), 10))))
	sig, err := private.Sign(digest[:])
	if err != nil {
		t.Error(err)
	}

	valid := ecdsa.Verify(private.PubKey().ToECDSA(), digest[:], sig.R, sig.S)
	if !valid {
		t.Error("sig not valid")
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
