// Copyright (C) 2018 Authors
// distributed under Apache 2.0 license

package bgls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	. "github.com/aliras1/FileTribe/tribecrypto/curves"
	"github.com/stretchr/testify/assert"
)

var curves = []CurveSystem{Altbn128}
var benchmarkCurve = Altbn128

func toHexBigInt(n *big.Int) string {
	return fmt.Sprintf("0x%x", n) // or %X or upper case
}

func pointToHexCoordsArray(p Point) []string {

	coords := p.ToAffineCoords()
	res := make([]string, len(coords))
	for i, coord := range coords {
		res[i] = toHexBigInt(coord)
	}
	return res
}

func pointsToStrArray(points []Point) []string {
	pointStrs := make([]string, 0)
	for i := 0; i < len(points); i++ {
		pointStrs = append(pointStrs, pointToHexCoordsArray(points[i])...)
	}
	return pointStrs
}

func TestSingleSigner(t *testing.T) {
	for _, curve := range curves {
		sk, vk, err := KeyGen(curve)
		assert.Nil(t, err, "Key generation failed")
		// 0x76699d3ff42826c5446413fd0eb7e0322416b1c9dddd62769d1628a084f822791bf3b990a29d3c70f87cfbc6b91b6ae972a5c914c96d5d74df6d92c55cdd38b1
		d := []byte{0x76, 0x69, 0x9d, 0x3f, 0xf4, 0x28, 0x26, 0xc5, 0x44, 0x64, 0x13, 0xfd, 0x0e, 0xb7, 0xe0, 0x32, 0x24, 0x16, 0xb1, 0xc9, 0xdd, 0xdd, 0x62, 0x76, 0x9d, 0x16, 0x28, 0xa0, 0x84, 0xf8, 0x22, 0x79, 0x1b, 0xf3, 0xb9, 0x90, 0xa2, 0x9d, 0x3c, 0x70, 0xf8, 0x7c, 0xfb, 0xc6, 0xb9, 0x1b, 0x6a, 0xe9, 0x72, 0xa5, 0xc9, 0x14, 0xc9, 0x6d, 0x5d, 0x74, 0xdf, 0x6d, 0x92, 0xc5, 0x5c, 0xdd, 0x38, 0xb1}
		// _, err = rand.Read(d)
		// assert.Nil(t, err, "test data generation failed")
		fmt.Println(hex.EncodeToString(d))
		sig := Sign(curve, sk, d)
		assert.True(t, VerifySingleSignature(curve, sig, vk, d), "Standard BLS "+
			"signature verification failed")

		//h := curve.HashToG1(d).Mul(new(big.Int).SetInt64(-1))
		h := curve.HashToG1(d)
		fmt.Println(h.MarshalUncompressed())
		fmt.Println(h.Marshal())
		hps := h.ToAffineCoords()
		for _, hq := range hps {
			fmt.Println(hq.String())
		}

		hp := []Point{h}
		sigp := []Point{sig}
		vkp := []Point{vk}
		g2p := []Point{curve.GetG2()}

		hpstr := pointsToStrArray(hp)
		sigpstr := pointsToStrArray(sigp)
		vkpstr := pointsToStrArray(vkp)
		g2pstr := pointsToStrArray(g2p)

		s, _ := json.Marshal(hpstr)
		fmt.Println(string(s))
		s, _ = json.Marshal(sigpstr)
		fmt.Println(string(s))
		s, _ = json.Marshal(vkpstr)
		fmt.Println(string(s))
		s, _ = json.Marshal(g2pstr)
		fmt.Println(string(s))

		sigTmp := sig.Copy()
		sigTmp, _ = sigTmp.Add(curve.GetG1())
		sig2 := sigTmp
		assert.False(t, VerifySingleSignature(curve, sig2, vk, d), "Standard BLS "+
			"signature verification succeeding when it shouldn't")

		// TODO Add tests to show that this doesn't succeed if d or vk is altered
	}
}

func TestAggregation(t *testing.T) {
	for _, curve := range curves {
		N, Size := 6, 32
		msgs := make([][]byte, N+1)
		sigs := make([]Point, N+1)
		pubkeys := make([]Point, N+1)
		for i := 0; i < N; i++ {
			msgs[i] = make([]byte, Size)
			rand.Read(msgs[i])

			sk, vk, _ := KeyGen(curve)
			sig := Sign(curve, sk, msgs[i])
			pubkeys[i] = vk
			sigs[i] = sig
		}
		aggSig := AggregateSignatures(sigs[:N])
		assert.True(t, VerifyAggregateSignature(curve, aggSig, pubkeys[:N], msgs[:N]),
			"Aggregate Point1 verification failed")
		assert.False(t, VerifyAggregateSignature(curve, aggSig, pubkeys[:N-1], msgs[:N]),
			"Aggregate Point1 verification succeeding without enough pubkeys")
		skf, vkf, _ := KeyGen(curve)
		pubkeys[N] = vkf
		sigs[N] = Sign(curve, skf, msgs[0])
		msgs[N] = msgs[0]
		aggSig = AggregateSignatures(sigs)
		assert.False(t, VerifyAggregateSignature(curve, aggSig, pubkeys, msgs),
			"Aggregate Signature succeeding with duplicate messages")
		assert.False(t, VerifyAggregateSignature(curve, aggSig, pubkeys[:N], msgs[:N]),
			"Aggregate Point1 succeeding with invalid signature")
		msgs[0] = msgs[1]
		msgs[1] = msgs[N]
		aggSig = AggregateSignatures(sigs[:N])
		assert.False(t, VerifyAggregateSignature(curve, aggSig, pubkeys[:N], msgs[:N]),
			"Aggregate Point1 succeeded with messages 0 and 1 switched")

		// TODO Add tests to make sure there is no mutation
	}
}

func BenchmarkKeygen(b *testing.B) {
	b.ResetTimer()
	curve := Altbn128
	for i := 0; i < b.N; i++ {
		_, _, res := KeyGen(curve)
		if res != nil {
			b.Error("key gen failure")
		}
	}
}

func BenchmarkAltBnHashToCurve(b *testing.B) {
	curve := Altbn128
	ms := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		ms[i] = make([]byte, 64)
		rand.Read(ms[i])
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		curve.HashToG1(ms[i])
	}
}

func BenchmarkSigning(b *testing.B) {
	curve := Altbn128
	sks := make([]*big.Int, b.N)
	ms := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		ms[i] = make([]byte, 64)
		rand.Read(ms[i])
		sk, _, _ := KeyGen(curve)
		sks[i] = sk
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Sign(curve, sks[i], ms[i])
	}
}

func BenchmarkVerification(b *testing.B) {
	curve := Altbn128
	message := make([]byte, 64)
	rand.Read(message)
	sk, vk, _ := KeyGen(curve)
	sig := Sign(curve, sk, message)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !VerifySingleSignature(curve, sig, vk, message) {
			b.Error("verification failed")
		}
	}
}

var vks []Point
var sgs []Point
var msg []byte

func TestMain(m *testing.M) {
	vks = make([]Point, 2048)
	sgs = make([]Point, 2048)
	msg = make([]byte, 64)
	rand.Read(msg)
	for i := 0; i < 2048; i++ {
		sk, vk, _ := KeyGen(benchmarkCurve)
		vks[i] = vk
		sgs[i] = KoskSign(benchmarkCurve, sk, msg)
	}
	os.Exit(m.Run())
}

func benchmulti(b *testing.B, k int) {
	//multisig := MultiSig{vks[:k], AggregateSignatures(sgs[:k]), msg}
	aggsig := AggregateSignatures(sgs[:k])
	keys := vks[:k]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !KoskVerifyMultiSignature(benchmarkCurve, aggsig, keys, msg) {
			b.Error("MultiSig verification failed")
		}
	}
}

func BenchmarkMultiVerification64(b *testing.B) {
	benchmulti(b, 64)
}

func BenchmarkMultiVerification128(b *testing.B) {
	benchmulti(b, 128)
}

func BenchmarkMultiVerification256(b *testing.B) {
	benchmulti(b, 256)
}

func BenchmarkMultiVerification512(b *testing.B) {
	benchmulti(b, 512)
}

func BenchmarkMultiVerification1024(b *testing.B) {
	benchmulti(b, 1024)
}

func BenchmarkMultiVerification2048(b *testing.B) {
	benchmulti(b, 2048)
}

func BenchmarkAggregateVerification(b *testing.B) {
	verifkeys := make([]Point, b.N)
	sigs := make([]Point, b.N)
	messages := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		messages[i] = make([]byte, 64)
		rand.Read(messages[i])
		sk, vk, _ := KeyGen(benchmarkCurve)
		verifkeys[i] = vk
		sigs[i] = Sign(benchmarkCurve, sk, messages[i])
	}
	aggsig := AggregateSignatures(sigs)
	b.ResetTimer()
	if !VerifyAggregateSignature(benchmarkCurve, aggsig, verifkeys, messages) {
		b.Error("Aggregate verificaton failed")
	}
}
