// Copyright (C) 2018 Authors
// distributed under Apache 2.0 license

package bgls

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

	. "github.com/aliras1/FileTribe/tribecrypto/curves" // nolint: golint
)

func pointToStrArray(point Point) []string {
	coords := point.ToAffineCoords()
	coordsStr := make([]string, len(coords))
	for k := 0; k < len(coords); k++ {
		coordsStr[k] = toHexBigInt(coords[k])
	}
	return coordsStr

}
func toHexBigInt(n *big.Int) string {
	return fmt.Sprintf("0x%x", n) // or %X or upper case
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func toBigInt(s string) *big.Int {
	bigInt := new(big.Int)
	bigInt, ok := bigInt.SetString(s, 0)
	if !ok {
		panic(fmt.Errorf("toBigInt() failed on string %v", s))
	}
	return bigInt
}

func boolToStr(boolRes bool) string {
	return fmt.Sprintf("%v", boolRes)
}

func bigIntToHexStr(bigInt *big.Int) string {
	return fmt.Sprintf("0x%x", bigInt)
}

func bigIntArrayToHexStrArray(bigInts []*big.Int) []string {

	arr := make([]string, len(bigInts))
	for i := 0; i < len(bigInts); i++ {
		arr[i] = bigIntToHexStr(bigInts[i])
	}
	return arr
}

func pointToHexCoords(p Point) string {

	return strings.Join(pointToHexCoordsArray(p), ",")
}

func pointToHexCoordsArray(p Point) []string {

	coords := p.ToAffineCoords()
	res := make([]string, len(coords))
	for i, coord := range coords {
		res[i] = toHexBigInt(coord)
	}
	return res
}

func pointsToStr(points []Point) string {
	return strings.Join(pointsToStrArray(points), ",")
}

func pointsToStrArray(points []Point) []string {
	pointStrs := make([]string, 0)
	for i := 0; i < len(points); i++ {
		pointStrs = append(pointStrs, pointToHexCoordsArray(points[i])...)
	}
	return pointStrs
}

//MultiSig holds set of keys and one message plus signature
type MultiSig struct {
	keys []Point
	sig  Point
	msg  []byte
}

//AggSig holds paired sequences of keys and messages, and one signature
type AggSig struct {
	keys []Point
	msgs [][]byte
	sig  Point
}

//KeyGen generates a *big.Int and Point2
func KeyGen(curve CurveSystem) (*big.Int, Point, error) {
	x, err := rand.Int(rand.Reader, curve.GetG1Order())
	if err != nil {
		return nil, nil, err
	}
	pubKey := LoadPublicKey(curve, x)
	return x, pubKey, nil
}

//LoadPublicKey turns secret key into a public key of type Point2
func LoadPublicKey(curve CurveSystem, sk *big.Int) Point {
	pubKey := curve.GetG2().Mul(sk)
	return pubKey
}

// Sign creates a standard BLS signature on a message with a private key
func Sign(curve CurveSystem, sk *big.Int, msg []byte) Point {
	return SignCustHash(sk, msg, curve.HashToG1)
}

// SignCustHash creates a standard BLS signature on a message with a private key,
// using a supplied function to hash onto the curve where signatures lie.
func SignCustHash(sk *big.Int, msg []byte, hash func([]byte) Point) Point {
	h := hash(msg)
	i := h.Mul(sk)
	return i
}

// VerifySingleSignature checks that a single standard BLS signature is valid
func VerifySingleSignature(curve CurveSystem, sig Point, pubKey Point, msg []byte) bool {
	return VerifySingleSignatureCustHash(curve, sig, pubKey, msg, curve.HashToG1)
}

// VerifySingleSignatureCustHash checks that a single standard BLS signature is
// valid, using the supplied hash function to hash onto the curve where signatures lie.
func VerifySingleSignatureCustHash(
	curve CurveSystem,
	sig Point,
	pubkey Point,
	msg []byte,
	hash func([]byte) Point) bool {

	h1 := hash(msg)
	coords1 := h1.ToAffineCoords()
	fmt.Println(coords1[0].String())
	fmt.Println(coords1[1].String())
	h := hash(msg).Mul(new(big.Int).SetInt64(-1))
	coords := h.ToAffineCoords()
	fmt.Println(coords[0].String())
	fmt.Println(coords[1].String())
	paired, _ := curve.PairingProduct([]Point{h, sig}, []Point{pubkey, curve.GetG2()})
	fmt.Println("h")
	fmt.Println(h.ToAffineCoords())
	fmt.Println("sig")
	fmt.Println(sig.ToAffineCoords())
	fmt.Println("pk")
	fmt.Println(pubkey.ToAffineCoords())
	fmt.Println("g")
	fmt.Println(curve.GetG2().ToAffineCoords())
	fmt.Println()
	return curve.GetGTIdentity().Equals(paired)
}

// Verify verifies an aggregate signature type.
func (a *AggSig) Verify(curve CurveSystem) bool {
	return VerifyAggregateSignature(curve, a.sig, a.keys, a.msgs)
}

// VerifyAggregateSignature verifies that the aggregated signature proves that
// all messages were signed by the associated keys. This will fail if there are
// duplicate messages, due to the possibility of the rogue public-key attack.
// If duplicate messages should be allowed, one of the protections against the
// rogue public-key attack should be used. See doc.go for more details.
func VerifyAggregateSignature(curve CurveSystem, aggsig Point, keys []Point, msgs [][]byte) bool {
	return verifyAggSig(curve, aggsig, keys, msgs, false)
}

// verifyMultiSignature checks that the aggregate signature correctly proves
// that a single message has been signed by a set of keys. This is
// vulnerable to the rogue public attack, so one of the defense mechanisms should be used.
func verifyMultiSignature(curve CurveSystem, aggsig Point, keys []Point, msg []byte) bool {
	vs := AggregatePoints(keys)
	return VerifySingleSignature(curve, aggsig, vs, msg)
}

func verifyAggSig(curve CurveSystem, aggsig Point, keys []Point, msgs [][]byte, allowDuplicates bool) bool {
	if len(keys) != len(msgs) {
		return false
	}
	if !allowDuplicates {
		if containsDuplicateMessage(msgs) {
			return false
		}
	}
	pts1 := make([]Point, len(keys)+1)
	pts2 := make([]Point, len(keys)+1)
	var wg sync.WaitGroup
	wg.Add(len(msgs))
	for i := 0; i < len(msgs); i++ {
		go concurrentHash(curve, i, pts1, msgs[i], &wg)
		pts2[i] = keys[i]
	}
	wg.Wait()
	pts1[len(keys)] = aggsig.Mul(new(big.Int).SetInt64(-1))
	pts2[len(keys)] = curve.GetG2()
	aggPt, ok := curve.PairingProduct(pts1, pts2)
	if ok {
		return aggPt.Equals(curve.GetGTIdentity())
	}
	return ok
}

// AggregateSignatures aggregates an array of signatures into one aggsig.
// This wrapper only exists so end-users don't have to use the method from curves
func AggregateSignatures(sigs []Point) Point {
	return AggregatePoints(sigs)
}

// AggregateKeys sums an array of public keys into one key.
// This wrapper only exists so end-users don't have to use the method from curve
func AggregateKeys(keys []Point) Point {
	return AggregatePoints(keys)
}

// concurrentHash hashes the message and sends the result down the channel.
func concurrentHash(curve CurveSystem, i int, pts []Point, msg []byte, wg *sync.WaitGroup) {
	pts[i] = curve.HashToG1(msg)
	wg.Done()
}

func containsDuplicateMessage(msgs [][]byte) bool {
	hashmap := make(map[string]bool)
	for i := 0; i < len(msgs); i++ {
		msg := string(msgs[i])
		if _, ok := hashmap[msg]; !ok {
			hashmap[msg] = true
		} else {
			return true
		}
	}
	return false
}
