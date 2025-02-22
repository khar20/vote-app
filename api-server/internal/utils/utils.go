package utils

import (
	"crypto/rand"
	"math/big"
)

// PublicKey represents the Paillier public key
type PublicKey struct {
	N *big.Int // modulus
	G *big.Int // generator
}

// PrivateKey represents the Paillier private key
type PrivateKey struct {
	PublicKey
	Lambda *big.Int // λ = lcm(p-1, q-1)
	Mu     *big.Int // μ = L(g^λ mod N²)⁻¹ mod N
}

// Decrypt decrypts a ciphertext
func (priv *PrivateKey) Decrypt(c *big.Int) *big.Int {
	nSquared := new(big.Int).Mul(priv.N, priv.N)
	// c^λ mod N²
	cLambda := new(big.Int).Exp(c, priv.Lambda, nSquared)
	// L(c^λ mod N²)
	l := new(big.Int).Div(new(big.Int).Sub(cLambda, big.NewInt(1)), priv.N)
	// m = L(c^λ mod N²) * μ mod N
	m := new(big.Int).Mod(new(big.Int).Mul(l, priv.Mu), priv.N)
	return m
}

func GeneratePaillierKeys(bitSize int) (*PublicKey, *PrivateKey, error) {
	// Generate two large prime numbers p and q
	p, err := rand.Prime(rand.Reader, bitSize/2)
	if err != nil {
		return nil, nil, err
	}
	q, err := rand.Prime(rand.Reader, bitSize/2)
	if err != nil {
		return nil, nil, err
	}

	// Compute n = p * q
	n := new(big.Int).Mul(p, q)

	// Compute λ = lcm(p-1, q-1)
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	qMinus1 := new(big.Int).Sub(q, big.NewInt(1))
	lambda := new(big.Int).Mul(pMinus1, qMinus1)
	gcd := new(big.Int).GCD(nil, nil, pMinus1, qMinus1)
	lambda.Div(lambda, gcd)

	// Choose g = n + 1 (a common choice for g)
	g := new(big.Int).Add(n, big.NewInt(1))

	// Compute μ = (L(g^λ mod n^2))⁻¹ mod n
	nSquared := new(big.Int).Mul(n, n)
	gLambda := new(big.Int).Exp(g, lambda, nSquared)
	l := new(big.Int).Div(new(big.Int).Sub(gLambda, big.NewInt(1)), n)
	mu := new(big.Int).ModInverse(l, n)

	// Create keys
	publicKey := &PublicKey{N: n, G: g}
	privateKey := &PrivateKey{PublicKey: *publicKey, Lambda: lambda, Mu: mu}

	return publicKey, privateKey, nil
}
