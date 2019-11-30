package tribecrypto

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/aliras1/FileTribe/tribecrypto/bgls"
	. "github.com/aliras1/FileTribe/tribecrypto/curves"
	"github.com/aliras1/FileTribe/tribecrypto/dkg"
)

// Usage examples:
// ./dkgmain -func=cgen

var cmd string

const INTERACTIVE = false
const BigintAsStrBase = 16 // All bigints as strings are in hex

const G1_COORDS_LEN = 2
const G2_COORDS_LEN = 4

type KeyPair struct {
	SK *big.Int
	PK Point
}

type Member struct {
	Key    *KeyPair
	Commit *DataForCommit
	Index  int
}

type DataForCommit struct {
	Coefficients 	[]*big.Int
	PubCommitG1  	[]Point
	PubCommitG2  	[]Point
	PrvCommit    	[]*big.Int
	PrvCommitEnc 	[]*big.Int
	MerkleTreePubCommitG1 *MerkleTree
	MerkleTreePubCommitG2 *MerkleTree
	MerkleTreePrvCommit   *MerkleTree
	SK              *big.Int
	PK              Point
}

type JsonDataForCommit struct {
	Coefficients []string
	PubCommitG1  []string
	PubCommitG2  []string
	PrvCommit    []string
	PrvCommitEnc []string
	SK           string
	PK           []string
}

// Conversions between array of numbers and G1/G2 points:
//func (g1Point *altbn128Point1) ToAffineCoords() []*big.Int
// func (g1Point *altbn128Point2) ToAffineCoords() []*big.Int
// func (curve *altbn128) MakeG2Point(coords []*big.Int, check bool) (Point, bool)

// Data for commitment:
// Generate t+1 random coefficients from (mod p) field for the polynomial
// Generate public commitments
// Generate private commitments

// index is 1-based
func GetCommitDataForSingleParticipant(curve CurveSystem, index int, t int, n int, mySK *big.Int, pks []Point) (*DataForCommit, error) {
	data := DataForCommit{
		Coefficients: make([]*big.Int, t+1),
		PubCommitG1:  make([]Point, t+1),
		PubCommitG2:  make([]Point, t+1),
		PrvCommit:    make([]*big.Int, n),
		PrvCommitEnc: make([]*big.Int, n),		
	}
	pubCommitG1MerkleLeaves := make([][]byte, t+1)
	pubCommitG2MerkleLeaves := make([][]byte, t+1)

	for i := 0; i < t+1; i++ {
		var err error
		data.Coefficients[i], data.PubCommitG1[i], data.PubCommitG2[i], err = dkg.CoefficientGen(curve)
		if err != nil {
			return nil, err
		}

		pubCommitG1MerkleLeaves[i] = ethcrypto.Keccak256(pointTo32ByteArray(data.PubCommitG1[i])...)
		pubCommitG2MerkleLeaves[i] = ethcrypto.Keccak256(pointTo32ByteArray(data.PubCommitG2[i])...)

		// verifyResult := dkg.VerifyPublicCommitment(curve, data.PubCommitG1[i], data.PubCommitG2[i])
		// if !verifyResult {
		// 	return nil, fmt.Errorf("VerifyPublicCommitment() failed for (participant=%v i=%v)", index, i)
		// }
		//fmt.Printf("PASSED VerifyPublicCommitment() (index=%v i=%v)\n", index, i)
	}

	data.MerkleTreePubCommitG1 = NewMerkleTree(pubCommitG1MerkleLeaves)
	data.MerkleTreePubCommitG2 = NewMerkleTree(pubCommitG2MerkleLeaves)

	prvCommitLeaves := make([][]byte, n)
	j := big.NewInt(1)
	for i := 0; i < n; i++ {
		//if i == index-1 {
		//  data.PrvCommit[i] = big.NewInt(0) // Don't calculate private commitment from me to myself
		//} else
		//{
		plainPrvCommit := dkg.GetPrivateCommitment(curve, j, data.Coefficients)
		//fmt.Printf("Calling Encrypt() with sk=%v pks[%v]=%v\n", mySK, i, pks[i].ToAffineCoords(), )
		data.PrvCommit[i] = plainPrvCommit
		data.PrvCommitEnc[i] = dkg.Encrypt(curve, mySK, pks[i], plainPrvCommit)

		prvCommitLeaves[i] = ethcrypto.Keccak256(bigIntTo32ByteArray(data.PrvCommitEnc[i]))
		//fmt.Printf("Encrypt() result: %v\n", data.PrvCommit[i])
		//}
		j.Add(j, big.NewInt(1))
	}

	data.MerkleTreePrvCommit = NewMerkleTree(prvCommitLeaves)

	return &data, nil
}

func SignAndVerify(curve CurveSystem, threshold int, n int, data []*DataForCommit) (bool, error) {

	// == Calculate SK, Pks and group PK ==
	// Should be happen only once, after DKG flow is done, and not for every SignAndVerify()

	//fmt.Println()
	//fmt.Printf("Starting SignAndVerify with threshold=%v n=%v\n", threshold, n)

	fmt.Println("Calculating SK, PK and Commitments - this is done just once, before signing & verifying messages.")

	// == Verify phase ==

	commitPrvAllDec := make([][]*big.Int, n)

	// First decrypt
	for committedParticipant := 0; committedParticipant < n; committedParticipant++ {
		pk := data[committedParticipant].PK // this is the encrypted pk
		//fmt.Printf("PK: %v\n", pk)
		commitPrvDec := make([]*big.Int, n)
		for participant := 0; participant < n; participant++ {
			if committedParticipant != participant {
				sk := data[participant].SK
				//fmt.Printf("SK[%v]: %v\n", participant, sk)
				enc := big.NewInt(0).Set(data[committedParticipant].PrvCommitEnc[participant]) // PrvCommit is encrypted
				//fmt.Printf("Enc prv commit[%v][%v]: %v\n", committedParticipant, participant, enc)
				commitPrvDec[participant] = dkg.Decrypt(curve, sk, pk, enc)
				//fmt.Printf("Dec prv commit[%v]: %v\n", participant, commitPrvDec[participant])
				if commitPrvDec[participant].Cmp(data[committedParticipant].PrvCommit[participant]) != 0 {
					panic("commitment is not the same after decryption")
				}
			} else {
				commitPrvDec[participant] = data[committedParticipant].PrvCommit[participant] // personal data
			}
		}
		commitPrvAllDec[committedParticipant] = commitPrvDec
	}
	fmt.Println("PASSED First Decrypt")

	// j := big.NewInt(1)
	// for participant := 0; participant < n; participant++ {
	// 	for commitParticipant := 0; commitParticipant < n; commitParticipant++ {
	// 		if participant != commitParticipant {
	// 			prv := commitPrvAllDec[commitParticipant][participant]
	// 			pub := data[commitParticipant].PubCommitG1
	// 			if !dkg.VerifyPrivateCommitment(curve, j, prv, pub) {
	// 				panic("private commit doesnt match public commit")
	// 			}
	// 		}
	// 	}
	// 	j.Add(j, big.NewInt(1))
	// }
	//fmt.Println("PASSED VerifyPrivateCommitment")

	// END OF DKG

	// == Calculate SK, Pks and group PK ==
	skAll := make([]*big.Int, n)
	//pkAll := make([][]Point, n)
	pubCommitG2Zero := make([]Point, n)
	pubCommitG2All := make([][]Point, n)
	for participant := 0; participant < n; participant++ {
		pubCommitG2All[participant] = data[participant].PubCommitG2
		//fmt.Printf("pubCommitG2All[%v]: %v\n", participant, pubCommitG2All[participant])
	}

	start := time.Now()
	pkAll := dkg.GetAllPublicKey(curve, threshold, pubCommitG2All)
	end := time.Now()
	fmt.Println(end.Sub(start).String())

	for participant := 0; participant < n; participant++ {
		pubCommitG2Zero[participant] = pubCommitG2All[participant][0]
		prvCommit := make([]*big.Int, n)
		for commitParticipant := 0; commitParticipant < n; commitParticipant++ {
			prvCommit[commitParticipant] = commitPrvAllDec[commitParticipant][participant]
		}

		start = time.Now()		
		skAll[participant] = dkg.GetSecretKey(prvCommit)		
		end = time.Now()
		fmt.Printf("SK: %s\n", end.Sub(start).String())
	}
	fmt.Println("Calculated SKs, PKs")

	//Verify pkAll are the same for all
	//fmt.Println("Public Key shares - same values are calculated by each client")
	// for participant := 0; participant < n; participant++ {
	// 	pks := pkAll[participant]
	// 	for otherParticipant := 0; otherParticipant < n; otherParticipant++ {
	// 		if !pks[participant].Equals(pkAll[otherParticipant][participant]) {
	// 			panic("pk for the same participant is different among other paricipants")
	// 		}
	// 	}
	// }

	//fmt.Println("PASSED Verification that same PKs are shared between all participants")

	//fmt.Println("Completed one-time calculation of SK, PK and Commitments")
	// fmt.Println("** SECRET KEYS [DEBUG ONLY] **")
	// for _, sk := range skAll {
	// 	fmt.Printf("** SK: %x\n", sk)
	// }
	//fmt.Println()

	start = time.Now()		
	groupPk := dkg.GetGroupPublicKey(curve, pubCommitG2Zero)	
	end = time.Now()		
	fmt.Printf("Group PK: %v\n", pointToHexCoords(groupPk))
	fmt.Printf("GROUP PK: %s\n", end.Sub(start).String())

	// coefsZero := make([]*big.Int, n)
	// for participant := 0; participant < n; participant++ {
	// 	coefsZero[participant] = data[participant].Coefficients[0]
	// }
	// groupSk := dkg.GetPrivateCommitment(curve, big.NewInt(1), coefsZero)
	// if !groupPk.Equals(bgls.LoadPublicKey(curve, groupSk)) {
	// 	panic("groupPK doesnt match to groupSK")
	// }

	// == Sign and reconstruct ==

	var msg string
	if INTERACTIVE {
		msg = readFromStdin("*** Enter message: ")
	} else {
		msg = "Hello Orbs"
	}	

	//fmt.Println()
	//fmt.Printf("Message for signature verification: %v\n", msg)
	msgBytes := []byte(msg)
	fmt.Println("message:")
	fmt.Println(hex.EncodeToString(msgBytes))
	fmt.Println(msgBytes)
	//fmt.Printf("Message bytes: %v\n", msgBytes)
	//fmt.Printf("Message hex: %v\n", hex.EncodeToString(msgBytes))
	sigs := make([]Point, n)

	// For each participant, generate signature with its SK
	for participant := 0; participant < n; participant++ {
		sigs[participant] = bgls.Sign(curve, skAll[participant], msgBytes)

		if !bgls.VerifySingleSignature(curve, sigs[participant], pkAll[participant], msgBytes) {
			return false, fmt.Errorf("signature invalid")
		}
		//fmt.Printf("PASSED VerifySingleSignature() sig share for client ID #%v: %v\n", participant+1, pointToHexCoords(sigs[participant]))
	}

	// Generates indices [0..n)
	indices := make([]*big.Int, n)
	index := big.NewInt(0)
	for participant := 0; participant < n; participant++ {
	index.Add(index, big.NewInt(1))
		indices[participant] = big.NewInt(0).Set(index)
	}

	
	// These are 1-based (not 0-based)
	subIndices := [][]int{}
	idxs := rand.Perm(n)
	idxs = idxs[:threshold+1]
	for i := 0; i <= threshold; i++ {
		idxs[i]++
	}

	subIndices = append(subIndices, idxs)

	for i := 0; i < len(subIndices); i++ {
		//fmt.Println()
		//fmt.Printf("=====> verifySigOnSubset() subIndices #%v <=====\n", subIndices[i])

		_, err := verifySigOnSubset(curve, indices, sigs, groupPk, msgBytes, subIndices[i])
		if err != nil {
			//fmt.Printf("Error in subgroup %v: %v", subIndices[i], err)
			return false, err
		}
		//fmt.Printf("PASSED verifySigOnSubset() subIndices #%v\n", subIndices[i])
		//fmt.Printf("Verify signature completed successfully for subgroup %v\n", subIndices[i])
		//fmt.Println("======================================================")
	}

	//fmt.Println()

	return true, nil
}

func verifySigOnSubset(curve CurveSystem, indices []*big.Int, sigs []Point, groupPk Point, msgBytes []byte, subIndices []int) (bool, error) {

	subSigs := make([]Point, len(subIndices))
	subIndicesBigInt := make([]*big.Int, len(subIndices))

	for i, idx := range subIndices {
		subSigs[i] = sigs[idx-1]
		subIndicesBigInt[i] = big.NewInt(int64(idx))
		//subIndices[i] = indices[idx]
	}

	//fmt.Printf("Sending to SignatureReconstruction(): indices=%v\n", subIndices)
	//for i, subSig := range subSigs {
	//fmt.Printf("Signature Share %v: %v\n", subIndicesBigInt[i], pointToHexCoords(subSig))
	//}
	start := time.Now()
	groupSig1, err := dkg.SignatureReconstruction(curve, subSigs, subIndicesBigInt)
	end := time.Now()
	fmt.Printf("Sig reconstruction: %v\n", end.Sub(start).String())
	if err != nil {
		return false, fmt.Errorf("group signature reconstruction failed")
	}

	//fmt.Printf("* Created group signature: %v *\n", pointToHexCoords(groupSig1))

	if !bgls.VerifySingleSignature(curve, groupSig1, groupPk, msgBytes) {
		return false, fmt.Errorf("group signature invalid")
	}
	//fmt.Printf("* PASSED VerifySingleSignature for subgroup signature: %v\n", pointToHexCoords(groupSig1))
	//fmt.Printf("Group PK: %v\n", pointToHexCoords(groupPk))

	return true, nil 
}
       
func TestBenchmark(t *testing.T) {
	var testCases = []struct{
		threshold int
		n		  int
	}{
		{2, 5},
		{5, 10},
		{25, 50},
		{50, 100},
		// {128, 256},
		// {500, 1000},
		// {1000, 2000}, 
	}
	curve := Altbn128
	
	for _, testCase := range(testCases) {
		n := testCase.n
		threshold := testCase.threshold
		group := make([]*Member, n)
		pks := make([]Point, n)
		commits := make([]*DataForCommit, n)

		fmt.Printf("N: %v, T: %v", n, threshold)

		// gen member
		for i := 0; i < n; i++ {
			group[i] = newMember(i+1, curve)
			pks[i] = group[i].Key.PK
		}

		// gen commits
		for i := 0; i < n; i++ {
			var err error
			member := group[i]
			start := time.Now()
			member.Commit, err = GetCommitDataForSingleParticipant(curve, member.Index, threshold, n, member.Key.SK, pks)
			end := time.Now()
			fmt.Printf("CommitData time: %v\n", end.Sub(start).String())
			if err != nil {
				panic(err)
			}

			member.Commit.SK = member.Key.SK
			member.Commit.PK = member.Key.PK
			commits[i] = member.Commit
		}

		// sign and verify
		isOk, err := SignAndVerify(curve, threshold, n, commits)
		if err != nil {
			panic(err)
		}
		if !isOk {
			panic("verification failed")
		}

		// case "VerifyPrivateCommitment":
		// 	// func VerifyPrivateCommitment(curve CurveSystem, myIndex *big.Int, prvCommit *big.Int, pubCommitG1 []Point) bool {
		// 	complainerIndex := toInt(flag.Arg(0)) // 1-based
		// 	accusedIndex := toInt(flag.Arg(1))    // 1-based

		// 	//myIndex, _ := strconv.Atoi(flag.Args()[0]) // 2   1-based
		// 	//prvCommit, _ := new(big.Int).SetString(flag.Arg(1), 0)
		// 	//pubCommitG1 := strToG1s(curve, flag.Arg(2))
		// 	// prvCommit is prvCommitAll[0][1] - this is what client 0 has commited to client 1
		// 	// pubCommitG1 [0] - this is all of client 0 public commitments over G1
		// 	myIndex := big.NewInt(int64(complainerIndex))
		// 	dataFile := flag.Arg(2)
		// 	jsonData := readDataFile(dataFile, curve)
		// 	data := make([]*DataForCommit, len(jsonData))
		// 	for i := 0; i < len(jsonData); i++ {
		// 		data[i] = jsonData[i].toData(curve)
		// 	}
		// 	prvCommit := data[accusedIndex-1].PrvCommit[complainerIndex-1]
		// 	pubCommitG1 := data[accusedIndex-1].PubCommitG1

		// 	fmt.Printf("Calling VerifyPrivateCommitment(): myIndex: %v prvCommit: %v pubCommitG1: %v\n", myIndex, bigIntToHexStr(prvCommit), pointsToStr(pubCommitG1))
		// 	passedVerification := dkg.VerifyPrivateCommitment(curve, myIndex, prvCommit, pubCommitG1)
		// 	res := fmt.Sprintf("%v\n", boolToStr(passedVerification))
		// 	fmt.Println(res) 
	}          
}     
  
type truffleCommit struct {
	SenderIndex int `json:"senderIndex"`
	RootPubCommitG1 string `json:"rootPubCommitG1"`
	RootPubCommitG2 string `json:"rootPubCommitG2"`
	RootEncPrvCommit string `json:"rootEncPrvCommit"`
	CommitIpfsHash string `json:"commitIpfsHash"`
	YG1 []string `json:"yG1"`
}   

type complaintPubCommit struct {
	ComplainerIndex int `json:"complainerIndex"`
    AccusedIndex int `json:"accusedIndex"`
	PubCommitG1 []string `json:"pubCommitG1"`
    PubCommitG2 []string `json:"pubCommitG2"`
    ProofG1 []string `json:"proofG1"`
    ProofG2 []string `json:"proofG2"`
}

type complaintPrvCommit struct {
	ComplainerIndex int `json:"complainerIndex"`
    AccusedIndex int `json:"accusedIndex"`	
	ComplainerSk string `json:"complainerSk"`
	EncPrvCommit string `json:"encPrvCommit"`
	ProofPrvCommit []string `json:"proofPrvCommit"`
	PubCommitsG1 [][]string `json:"pubCommitsG1"`
}
         
type TruffleTestData struct {
	T int `json:"t"`
	N int `json:"n"`
	PKs [][]string `json:"pks"`
	Commits []truffleCommit  `json:"commits"`
	ComplaintPubCommit complaintPubCommit `json:"complaintPubCommit"`
	ComplaintPrvCommit complaintPrvCommit `json:"complaintPrvCommit"`
}

func TestGenTestDataForComplaining(t *testing.T) {
	var testCases = []struct{
		threshold int
		n		  int
	}{
		{1, 2},
		{2, 5},
		{5, 10},
		{25, 50},
		{50, 100},
		// {128, 256},
		// {500, 1000},
		// {1000, 2000},
	} 
	curve := Altbn128
	genTests := make([]*TruffleTestData, len(testCases))

	for j, testCase := range(testCases) {
		genData := new(TruffleTestData)
		n := testCase.n
		threshold := testCase.threshold
				
		// gen members & commits
		group, pks := genMembers(n, curve)
		genCommits(group, pks, threshold, curve)

		genData.T = threshold 
		genData.N = n
		genData.PKs = make([][]string, n)
		genData.Commits = make([]truffleCommit, 1)
		for i, member := range(group) {
			genData.PKs[i] = pointToStrArray(member.Key.PK)			
		}
		member := group[0]
		genData.Commits[0].SenderIndex = member.Index
		genData.Commits[0].RootPubCommitG1 = fmt.Sprintf("0x%s", hex.EncodeToString(member.Commit.MerkleTreePubCommitG1.Root()))
		genData.Commits[0].RootPubCommitG2 = fmt.Sprintf("0x%s", hex.EncodeToString(member.Commit.MerkleTreePubCommitG2.Root()))
		genData.Commits[0].RootEncPrvCommit = fmt.Sprintf("0x%s", hex.EncodeToString(member.Commit.MerkleTreePrvCommit.Root()))
		genData.Commits[0].CommitIpfsHash = "ipfshash"
		genData.Commits[0].YG1 = pointToStrArray(member.Commit.PubCommitG1[0])

		leafG1 := ethcrypto.Keccak256(pointTo32ByteArray(member.Commit.PubCommitG1[1])...)
		leafG2 := ethcrypto.Keccak256(pointTo32ByteArray(member.Commit.PubCommitG2[1])...)
		leafPrv := ethcrypto.Keccak256(bigIntTo32ByteArray(member.Commit.PrvCommitEnc[1]))		
		fmt.Printf("leafG1: 0x%s\n", hex.EncodeToString(leafG1))
		fmt.Printf("leafG2: 0x%s\n", hex.EncodeToString(leafG2))
		fmt.Printf("leafPrv: 0x%s\n", hex.EncodeToString(leafPrv))

		proofG1, err := member.Commit.MerkleTreePubCommitG1.Prove(leafG1)
		if err != nil {
			t.Fatal(err)
		}
		if !VerifyMerkleProof(proofG1, member.Commit.MerkleTreePubCommitG1.Root(), leafG1) {
			t.Fatal("invalid proof G1")
		}
		proofG2, err := member.Commit.MerkleTreePubCommitG2.Prove(leafG2)
		if err != nil {
			t.Fatal(err)
		}
		if !VerifyMerkleProof(proofG2, member.Commit.MerkleTreePubCommitG2.Root(), leafG2) {
			t.Fatal("invalid proof G1")
		}
		proofPrv, err := member.Commit.MerkleTreePrvCommit.Prove(leafPrv)
		if err != nil {
			t.Fatal(err)
		}
		if !VerifyMerkleProof(proofPrv, member.Commit.MerkleTreePrvCommit.Root(), leafPrv) {
			t.Fatal("invalid proof Prv commit")
		}

		genData.ComplaintPubCommit = complaintPubCommit{
			ComplainerIndex: 2,
			AccusedIndex: 1,
			PubCommitG1: pointToStrArray(member.Commit.PubCommitG1[1]),
			PubCommitG2: pointToStrArray(member.Commit.PubCommitG2[1]),		
			ProofG1: bytesArrayToStringArray(proofG1),
			ProofG2: bytesArrayToStringArray(proofG2),
		}

		pubCommitsG1 := make([][]string, len(member.Commit.PubCommitG1))
		for i, commit := range(member.Commit.PubCommitG1) {
			pubCommitsG1[i] = pointToStrArray(commit)
		}

		genData.ComplaintPrvCommit = complaintPrvCommit {
			ComplainerIndex: 2,
			AccusedIndex: 1,
			ComplainerSk: fmt.Sprintf("0x%s", hex.EncodeToString(bigIntTo32ByteArray(group[1].Key.SK))),
			EncPrvCommit: fmt.Sprintf("0x%s", hex.EncodeToString(bigIntTo32ByteArray(member.Commit.PrvCommitEnc[1]))),
			ProofPrvCommit: bytesArrayToStringArray(proofPrv),
			PubCommitsG1:  pubCommitsG1,
		}

		// rand participant --> make him prove everything

		genTests[j] = genData		
	}

	jsonBytes, err := json.Marshal(genTests)
	if err != nil {
		t.Fatal(err)
	}
	
	if err := ioutil.WriteFile("truffle_test_data_complaints.json", jsonBytes, 0644); err != nil {
		t.Fatal(err)
	}
}
       
func TestGenerateTestDataForTruffle(t *testing.T) {
	var testCases = []struct{
		threshold int
		n		  int
	}{
		{1, 2},
		{2, 5},
		{5, 10},
		{25, 50},
		{50, 100},
		// {128, 256},
		// {500, 1000},
		// {1000, 2000},
	} 
	curve := Altbn128
	genTests := make([]*TruffleTestData, len(testCases))

	for j, testCase := range(testCases) {
		genData := new(TruffleTestData)
		n := testCase.n
		threshold := testCase.threshold
				
		// gen members & commits
		group, pks := genMembers(n, curve)
		genCommits(group, pks, threshold, curve)

		genData.T = threshold 
		genData.N = n
		genData.PKs = make([][]string, n)
		genData.Commits = make([]truffleCommit, n)
		for i, member := range(group) {
			genData.PKs[i] = pointToStrArray(member.Key.PK)
			genData.Commits[i].SenderIndex = member.Index
			genData.Commits[i].RootPubCommitG1 = fmt.Sprintf("0x%s", hex.EncodeToString(member.Commit.MerkleTreePubCommitG1.Root()))
			genData.Commits[i].RootPubCommitG2 = fmt.Sprintf("0x%s", hex.EncodeToString(member.Commit.MerkleTreePubCommitG2.Root()))
			genData.Commits[i].RootEncPrvCommit = fmt.Sprintf("0x%s", hex.EncodeToString(member.Commit.MerkleTreePrvCommit.Root()))
			genData.Commits[i].CommitIpfsHash = "ipfshash"
			genData.Commits[i].YG1 = pointToStrArray(member.Commit.PubCommitG1[0])
		}

		genTests[j] = genData		
	}

	jsonBytes, err := json.Marshal(genTests)
	if err != nil {
		t.Fatal(err)
	}
	
	if err := ioutil.WriteFile("truffle_test_data.json", jsonBytes, 0644); err != nil {
		t.Fatal(err)
	}
}

func genMembers(n int, curve CurveSystem) ([]*Member, []Point) {
	group := make([]*Member, n)
	pks := make([]Point, n)
	for i := 0; i < n; i++ {
		group[i] = newMember(i+1, curve)
		pks[i] = group[i].Key.PK
	}
	return group, pks
}

func genCommits(group []*Member, pks []Point, threshold int, curve CurveSystem) []*DataForCommit {
	n := len(group)
	commits := make([]*DataForCommit, n)	
	for i := 0; i < n; i++ {
		var err error
		member := group[i]
		start := time.Now()
		member.Commit, err = GetCommitDataForSingleParticipant(curve, member.Index, threshold, n, member.Key.SK, pks)
		end := time.Now()
		fmt.Printf("CommitData time: %v\n", end.Sub(start).String())
		if err != nil {
			panic(err)
		}

		member.Commit.SK = member.Key.SK
		member.Commit.PK = member.Key.PK
		commits[i] = member.Commit
	}

	return commits
}

func newMember(idx int, curve CurveSystem) *Member {
	sk, pk, _, _ := dkg.CoefficientGen(curve)
	keyPair := KeyPair{sk, pk}

	return &Member{
	Key:   &keyPair,
		Index: idx,
	}
}

// Gets array of the form p0[0], p0[1], p1[0], p1[1], p2[0], p2[1], etc.
// Each pair is a G1 point so an array of Points is returned.
// This is not
func strToG1s(curve CurveSystem, pointStr string) []Point {
	//fmt.Printf("pointStr=%v\n", pointStr)
	pointStrCoords := strings.Split(pointStr, ",")
	points := make([]Point, len(pointStrCoords)/2)
	for i := 0; i < len(pointStrCoords); i += 2 {
		//fmt.Printf("Reading pointsStrCoords i=%v of %v", i, len(pointStrCoords))
		coord0, ok := new(big.Int).SetString(pointStrCoords[i], 0)
		if !ok {
			panic(fmt.Errorf("failed parsing coord0 to big.Int: %v (big.Int value: %v)", pointStrCoords[i], coord0))
		}
		coord1, ok := new(big.Int).SetString(pointStrCoords[i+1], 0)
		if !ok {
			panic(fmt.Errorf("failed parsing coord1 to big.Int: %v (big.Int value: %v)", pointStrCoords[i], coord1))
		}

		bigintCoords := []*big.Int{coord0, coord1}
		//fmt.Printf("strToG1: coord0=%v coord1=%v\n", coord0, coord1)
		point, _ := curve.MakeG1Point(bigintCoords, true)
		points[i/2] = point
	}
	return points
}

func strToG1(curve CurveSystem, pointStr string) Point {
	pointStrCoords := strings.Split(pointStr, ",")
	bigintCoords := make([]*big.Int, len(pointStrCoords))
	for i := 0; i < len(pointStr); i++ {
		bigintCoords[i], _ = new(big.Int).SetString(pointStrCoords[i], 0)
	}
	point, _ := curve.MakeG1Point(bigintCoords, true)
	return point
}

func readDataFile(dataFile string, curve CurveSystem) []*JsonDataForCommit {
	var inBuf []byte
	var resJson []*JsonDataForCommit
	var err error
	inBuf, err = ioutil.ReadFile(dataFile)
	//err = readGob("./data.gob", data)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(inBuf, &resJson); err != nil {
		panic(err)
	}
	return resJson
}

func (keyPair KeyPair) Marshal() ([]byte, error) {

	return json.Marshal(keyPair)
}

func bytesArrayToStringArray(bytesArray [][]byte) ([]string) {
	strArr := make([]string, len(bytesArray))
	for i, b := range(bytesArray) {
		strArr[i] = fmt.Sprintf("0x%s", hex.EncodeToString(b))
	}
	return strArr
}

func bigIntTo32ByteArray(num *big.Int) []byte {
	numBytes := num.Bytes()
	b := make([]byte, 32-len(numBytes))
	return append(b, numBytes...)
}

func pointTo32ByteArray(point Point) [][]byte {
	coords := point.ToAffineCoords()
	coordsBytesArr := make([][]byte, len(coords))
	for k := 0; k < len(coords); k++ {
		coordsBytesArr[k] = bigIntTo32ByteArray(coords[k])
	}
	return coordsBytesArr
}

func pointToStrArray(point Point) []string {
	coords := point.ToAffineCoords()
	coordsStr := make([]string, len(coords))
	for k := 0; k < len(coords); k++ {
		coordsStr[k] = fmt.Sprintf("0x%s", hex.EncodeToString(bigIntTo32ByteArray(coords[k])))
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

func readFromStdin(caption string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println()
	fmt.Print(caption)
	text, _ := reader.ReadString('\n')
return text
}

func (jd *JsonDataForCommit) toData(curve CurveSystem) *DataForCommit {

	res := new(DataForCommit)
	res.Coefficients = make([]*big.Int, len(jd.Coefficients))
	res.PubCommitG1 = make([]Point, len(jd.PubCommitG1)/G1_COORDS_LEN)
	res.PubCommitG2 = make([]Point, len(jd.PubCommitG2)/G2_COORDS_LEN)
	res.PrvCommit = make([]*big.Int, len(jd.PrvCommit))
	res.PrvCommitEnc = make([]*big.Int, len(jd.PrvCommitEnc))

	for i := 0; i < len(jd.Coefficients); i++ {
		res.Coefficients[i] = toBigInt(jd.Coefficients[i])
	}

	res.SK = toBigInt(jd.SK)

	coords := make([]*big.Int, G1_COORDS_LEN)
	for k := 0; k < G1_COORDS_LEN; k++ {
		coords[k] = toBigInt(jd.PK[k])
	}
	pk, _ := curve.MakeG1Point(coords, true)
	res.PK = pk

	// PubCommitG1 is a string array of flattened x,y points in this format: x0,y0,x1,y1,x2,y2,..
	for i := 0; i < len(jd.PubCommitG1); i += G1_COORDS_LEN {
		coords := make([]*big.Int, G1_COORDS_LEN)
		for k := 0; k < G1_COORDS_LEN; k++ {
			coords[k] = toBigInt(jd.PubCommitG1[i+k])
		}
		var isOk bool
		res.PubCommitG1[i/G1_COORDS_LEN], isOk = curve.MakeG1Point(coords, true)
		//fmt.Printf("G1 Point: %v\n", res.PubCommitG1[i/G1_COORDS_LEN].ToAffineCoords())
		if !isOk {
			panic(fmt.Errorf("failed to make G1 point"))
		}
	}

	// PubCommitG2 is a string array of flattened x,y,z,t points in this format: x0,y0,z0,t0,x1,y1,z1,t1,x2,y2,z2,t2,..
	for i := 0; i < len(jd.PubCommitG2); i += G2_COORDS_LEN {
		coords := make([]*big.Int, G2_COORDS_LEN)
		for k := 0; k < G2_COORDS_LEN; k++ {
			coords[k] = toBigInt(jd.PubCommitG2[i+k])
		}
		var isOk bool
		//fmt.Printf("Make G2 point from coords: %v\n", coords)
		res.PubCommitG2[i/G2_COORDS_LEN], isOk = curve.MakeG2Point(coords, true)
		//fmt.Printf("G2 Point: %v\n", res.PubCommitG2[i/G2_COORDS_LEN].ToAffineCoords())
		if !isOk {
			panic(fmt.Errorf("failed to make G2 point"))
		}
	}

	for i := 0; i < len(jd.PrvCommit); i++ {
		res.PrvCommit[i] = toBigInt(jd.PrvCommit[i])
	}
	for i := 0; i < len(jd.PrvCommitEnc); i++ {
		res.PrvCommitEnc[i] = toBigInt(jd.PrvCommitEnc[i])
	}

	return res

}

func (data DataForCommit) MarshalJSON() ([]byte, error) {

	res := new(JsonDataForCommit)

	res.Coefficients = bigIntArrayToHexStrArray(data.Coefficients)
	res.PubCommitG1 = pointsToStrArray(data.PubCommitG1)
	res.PubCommitG2 = pointsToStrArray(data.PubCommitG2)
	//fmt.Printf("G1=%v G2=%v\n", res.PubCommitG1, res.PubCommitG2)
	res.PrvCommit = bigIntArrayToHexStrArray(data.PrvCommit)
	res.PrvCommitEnc = bigIntArrayToHexStrArray(data.PrvCommitEnc)

	return json.Marshal(res)
}
