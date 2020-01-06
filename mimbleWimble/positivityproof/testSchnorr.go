package main

import (
	"github.com/rynobey/bn256"
	"github.com/viteshan/simple-test/mimbleWimble/positivityproof/util"
	//"golang.org/x/crypto/bn256"

	"fmt"
	"math/big"
)

func main() {
	m := "veritaserum"

	sqrtb := util.CryptoRandBigInt()
	sqrtv := new(big.Int).SetInt64(10000)

	kb := util.CryptoRandBigInt()
	kv := util.CryptoRandBigInt()

	G := new(bn256.G1).ScalarBaseMult(new(big.Int).SetInt64(1))

	secret := util.CryptoRandBigInt()

	H := new(bn256.G1).ScalarMult(G, secret)

	for count := 0; count < 8; count++ {
		sb, sv, e, K, Cp, C := util.GenSignature(m, sqrtb, sqrtv, kb, kv, G, H)

		isValid := util.VerifySignature(m, sb, sv, e, K, Cp, C, G, H)
		fmt.Printf("%t\n", isValid)
	}

}
