package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/anyswap/CrossChain-Bridge/tools/crypto"
	"github.com/btcsuite/btcd/btcec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
	//core "github.com/terra-project/core/types"
)

var CDC = amino.NewCodec()

func init() {
	/*config := sdk.GetConfig()
	config.SetCoinType(core.CoinType)
	config.SetFullFundraiserPath(core.FullFundraiserPath)
	config.SetBech32PrefixForAccount(core.Bech32PrefixAccAddr, core.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(core.Bech32PrefixValAddr, core.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(core.Bech32PrefixConsAddr, core.Bech32PrefixConsPub)
	config.Seal()*/

	sdk.RegisterCodec(CDC)
	ctypes.RegisterAmino(CDC)

	//stargate-final
	bank.RegisterCodec(CDC)
	authtypes.RegisterCodec(CDC)

	// tequila-0004
	//CDC.RegisterConcrete(bank.MsgSend{}, "bank/MsgSend", nil)s
	//CDC.RegisterConcrete(authtypes.StdTx{}, "core/StdTx", nil)
}

//var ChainID = "tequila-0004"

var ChainID = "stargate-final"

func genKeyX() {
	priv := secp256k1.GenPrivKey()
	privkeyHex := hex.EncodeToString(priv.Bytes())
	pub := priv.PubKey().Bytes()
	pubkeyHex := hex.EncodeToString(pub)
	fmt.Printf("Private key: %v\nPublic key: %v\n", privkeyHex, pubkeyHex)
	pubkeyAddress := priv.PubKey().Address()
	fmt.Printf("Public key address: %v\n", pubkeyAddress)
	address, _ := sdk.AccAddressFromHex(pubkeyAddress.String())
	fmt.Printf("Address: %v\n", address.String())
}

func main() {
	//genKey()
	sendTx()
}

func genKey() {
	priv, _ := crypto.GenerateKey()
	privkeyHex := hex.EncodeToString(crypto.FromECDSA(priv))
	pub := priv.PublicKey
	pubkeyHex := hex.EncodeToString(crypto.FromECDSAPub(&pub))
	fmt.Printf("Private key: %v\nPublic key: %v\n", privkeyHex, pubkeyHex)
	address, err := PublicKeyToAddress(pubkeyHex)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Address: %v\n", address)

	var privBytes [32]byte
	privBytes1, _ := hex.DecodeString(privkeyHex)
	copy(privBytes[:], privBytes1[:33])
	priv1 := secp256k1.PrivKeySecp256k1(privBytes)
	pub1 := priv1.PubKey()
	pubkeyAddress := pub1.Address()
	address1, _ := sdk.AccAddressFromHex(pubkeyAddress.String())
	fmt.Printf("Address1: %v\n", address1.String())
}

/*
62625de7ed1d9ecaebb0dc8fe1425cfe994bb79c699c019d6b8c40f9e1ad8907
046693b7612ccd92f0ec57e62aa51c72d7a978c4871c482cbf2a896575bf67ac3041b8b0ee98f5c1433115c99bdba8548939d1a7eb232ec6df00af1d8d749ec23d
terra1qj05rkrpphd55dawh7qxxmd2c72g57j2r0nlp3
tequila-0004: 29325

1a05233ffa885bf369b5ff1ec829114975243fc7dbdbaabdee0cb9e4185dd678
04bfd55e4900a1de682907642843d16fd189ccac0656fcebacd22b3e10eecc6a374344c33799bca2831ac0cef6c739558ce879bb99ef2864a533e1e50e4d9dad6b
cosmos1s88t76ev084c6d35fahkslqseep5szgeggr3q0
stargate-final: 25361
*/

func sendTx() {
	//pubkeyHex := "046693b7612ccd92f0ec57e62aa51c72d7a978c4871c482cbf2a896575bf67ac3041b8b0ee98f5c1433115c99bdba8548939d1a7eb232ec6df00af1d8d749ec23d" // tequila-0004
	pubkeyHex := "04bfd55e4900a1de682907642843d16fd189ccac0656fcebacd22b3e10eecc6a374344c33799bca2831ac0cef6c739558ce879bb99ef2864a533e1e50e4d9dad6b" // stargate-final
	addr, err := PublicKeyToAddress(pubkeyHex)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("address:\n%v\n", addr) // terra1wr965945fxexk68mxne5en56pr6la006p7764e
	// cosmos1s88t76ev084c6d35fahkslqseep5szgeggr3q0
	address, _ := sdk.AccAddressFromBech32(addr)

	address2 := address

	// accountNumber := uint64(29325) // tequila-0004
	accountNumber := uint64(25361) // stargate-final
	sequence := uint64(2)

	msgs := []sdk.Msg{
		bank.MsgSend{
			FromAddress: address,
			ToAddress:   address2,
			//Amount:      sdk.Coins{sdk.NewCoin("uluna", sdk.NewInt(100000))},
			Amount: sdk.Coins{sdk.NewCoin("umuon", sdk.NewInt(10))},
		},
	}
	memo := ""

	//feeAmount := sdk.Coins{sdk.Coin{"uluna", sdk.NewInt(50000)}}
	feeAmount := sdk.Coins{sdk.Coin{"umuon", sdk.NewInt(50000)}}
	gas := uint64(300000)
	fee := authtypes.NewStdFee(gas, feeAmount)

	signBytes := authtypes.StdSignBytes(ChainID, accountNumber, sequence, fee, msgs, memo)
	signHash := fmt.Sprintf("%X", tmhash.Sum(signBytes))
	fmt.Printf("\nSign bytes hash:\n%v\n", signHash)

	var privBytes [32]byte
	//privBytes1, _ := hex.DecodeString("62625de7ed1d9ecaebb0dc8fe1425cfe994bb79c699c019d6b8c40f9e1ad8907") // tequila-0004
	privBytes1, _ := hex.DecodeString("1a05233ffa885bf369b5ff1ec829114975243fc7dbdbaabdee0cb9e4185dd678") // stargate-final
	copy(privBytes[:], privBytes1[:33])
	priv := secp256k1.PrivKeySecp256k1(privBytes)
	signature, err := priv.Sign(signBytes)

	stdsig := authtypes.StdSignature{
		PubKey:    priv.PubKey(),
		Signature: signature,
	}
	signatures := []authtypes.StdSignature{stdsig}
	stdtx := authtypes.StdTx{
		Msgs:       msgs,
		Fee:        fee,
		Signatures: signatures,
		Memo:       memo,
	}
	fmt.Printf("\nStd tx:\n%+v\n", stdtx)

	txBytes, err := CDC.MarshalBinaryLengthPrefixed(stdtx)
	if err != nil {
		log.Fatal(err)
	}
	txhash := ttypes.Tx(txBytes).Hash()
	fmt.Printf("\nTxhash:\n%X\n", txhash)

	// build post data
	bz, err := CDC.MarshalJSON(stdtx)
	if err != nil {
		log.Fatal(err)
	}
	// Take "value" from the json struct
	tempStr := make(map[string]interface{})
	err = json.Unmarshal(bz, &tempStr)
	if err != nil {
		log.Fatal(err)
	}
	value, ok := tempStr["value"].(map[string]interface{})
	if !ok {
		log.Fatal(err)
	}
	// repass account number and sequence
	signatures2, ok := value["signatures"].([]interface{})
	if !ok || len(signatures) < 1 {
		log.Fatal(err)
	}
	signatures2[0].(map[string]interface{})["account_number"] = fmt.Sprintf("%v", accountNumber)
	signatures2[0].(map[string]interface{})["sequence"] = fmt.Sprintf("%v", sequence)
	value["signatures"] = signatures2
	bz2, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}
	data := fmt.Sprintf(`{"tx":%v,"mode":"block"}`, string(bz2))
	fmt.Printf("\ndata:\n%+v\n", data)

	// broadcast
	// https://github.com/cosmos/cosmos-sdk/blob/v0.39.2/x/auth/client/rest/broadcast.go
	fmt.Printf("\n====================\n")
	client := &http.Client{}
	//req, err := http.NewRequest("POST", "https://tequila-lcd.terra.dev/txs", strings.NewReader(data)) // tequila-0004
	req, err := http.NewRequest("POST", "http://34.71.170.158:1317/txs", strings.NewReader(data)) // stargate-final
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nres:\n%v\n", string(bodyText))

	// simulate handle broadcast
	fmt.Printf("\n====================\n")
	r, err := http.NewRequest("POST", "", strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	var req1 rest.BroadcastReq
	body, err := ioutil.ReadAll(r.Body)
	err = CDC.UnmarshalJSON(body, &req1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nreq1.Tx:\n%+v\n", req1.Tx)
	txBytes2, err := CDC.MarshalBinaryLengthPrefixed(req1.Tx)
	if err != nil {
		log.Fatal(err)
	}
	txhash2 := ttypes.Tx(txBytes2).Hash()
	fmt.Printf("\nTxhash2:\n%X\n", txhash2)

	// verify tx
	// https://github.com/cosmos/cosmos-sdk/blob/v0.39.2/x/auth/ante/sigverify.go
	fmt.Printf("\n====================\n")
	sigTx := ante.SigVerifiableTx(req1.Tx)
	fmt.Printf("\nSigTx:\n%+v\n", sigTx)

	pubkeys := sigTx.GetPubKeys()
	signers := sigTx.GetSigners()
	sigs := sigTx.GetSignatures()
	fmt.Printf("\npubkeys:\n%+v\nsigners:\n%+v\nsigs:\n%+v\n", pubkeys, signers, sigs)

	fmt.Printf("\nsignBytes:\n%X\n", signBytes)
	valid1 := pubkeys[0].VerifyBytes(signBytes, sigs[0])
	fmt.Printf("\nvalid1:\n%v\n", valid1)

	signBytes2 := authtypes.StdSignBytes(
		ChainID, accountNumber, sequence, req1.Tx.Fee, req1.Tx.Msgs, req1.Tx.Memo,
	)
	fmt.Printf("\nsignBytes2:\n%X\n", signBytes2)
	valid2 := pubkeys[0].VerifyBytes(signBytes2, sigs[0])
	fmt.Printf("\nvalid2:\n%v\n", valid2)
}

func PublicKeyToAddress(pubKeyHex string) (address string, err error) {
	pubKeyHex = strings.TrimPrefix(pubKeyHex, "0x")
	bb, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return
	}
	pk, err := btcec.ParsePubKey(bb, btcec.S256())
	if err != nil {
		return
	}
	cpk := pk.SerializeCompressed()
	var pub [33]byte
	copy(pub[:], cpk[:33])
	pubkey := secp256k1.PubKeySecp256k1(pub)
	addr := pubkey.Address()
	accAddress, err := sdk.AccAddressFromHex(addr.String())
	if err != nil {
		return
	}
	address = accAddress.String()
	return
}
