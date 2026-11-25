package password_test

import (
	"strings"
	"testing"

	"github.com/ozaitsev92/tododdd/internal/infrastructure/password"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher(t *testing.T) {
	t.Parallel()

	h, err := password.NewBcryptHasher(bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := h.Hash("Secret123")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "Secret123" || !strings.HasPrefix(hash, "$2a$04$") {
		t.Fatalf("unexpected hash %q", hash)
	}
	if err := h.Compare(hash, "Secret123"); err != nil {
		t.Fatalf("Compare(correct) = %v", err)
	}
	if err := h.Compare(hash, "wrong"); err == nil {
		t.Fatal("Compare(wrong) succeeded")
	}
	if again, _ := h.Hash("Secret123"); again == hash {
		t.Fatal("hashes are not salted")
	}
}

func TestBcryptHasher_UsesConfiguredCost(t *testing.T) {
	t.Parallel()
	h, err := password.NewBcryptHasher(10)
	if err != nil {
		t.Fatal(err)
	}
	hash, _ := h.Hash("Secret123")
	if cost, _ := bcrypt.Cost([]byte(hash)); cost != 10 {
		t.Fatalf("cost = %d, want 10", cost)
	}
}

func TestNewBcryptHasher_RejectsInvalidCost(t *testing.T) {
	t.Parallel()
	for _, c := range []int{0, 3, 32} {
		if _, err := password.NewBcryptHasher(c); err == nil {
			t.Errorf("cost %d accepted", c)
		}
	}
}
