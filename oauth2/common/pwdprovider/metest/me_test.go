package metest

import (
	"testing"

	"github.com/bernishen/lion-go/oauth2/common/pwdprovider"
)

func TestProviderV1(t *testing.T) {
	t.Log(pwdprovider.Compute(1, "123"))
	t.Log(pwdprovider.Compute(2, "123"))
}
