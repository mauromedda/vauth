package token

import (
	"testing"

	vt "github.com/hashicorp/vault/command/token"
)

// TestCommand re-uses the existing Test function to ensure proper behavior of
// the internal token helper
func TestCommand(t *testing.T) {
	vt.Test(t, &InternalTokenHelper{})
}
