// Copyright (c) 2017-2019 The Elastos Foundation
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.
//

package api

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/core/contract"
	"github.com/elastos/Elastos.ELA/core/types/payload"
	"github.com/elastos/Elastos.ELA/crypto"

	lua "github.com/yuin/gopher-lua"
)

const (
	luaCoinBaseTypeName      = "coinbase"
	luaTransferAssetTypeName = "transferasset"
	luaRegisterProducerName  = "registerproducer"
	luaUpdateProducerName    = "updateproducer"
	luaCancelProducerName    = "cancelproducer"
	luaActivateProducerName  = "activateproducer"
	luaReturnDepositCoinName = "returndepositcoin"
	luaSideChainPowName      = "sidechainpow"
	luaRegisterCRName        = "registercr"
	luaUpdateCRName          = "updatecr"
	luaUnregisterCRName      = "unregistercr"
)

func RegisterCoinBaseType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaCoinBaseTypeName)
	L.SetGlobal("coinbase", mt)
	L.SetField(mt, "new", L.NewFunction(newCoinBase))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), coinbaseMethods))
}

// Constructor
func newCoinBase(L *lua.LState) int {
	data, _ := hex.DecodeString(L.ToString(1))
	cb := &payload.CoinBase{
		Content: data,
	}
	ud := L.NewUserData()
	ud.Value = cb
	L.SetMetatable(ud, L.GetTypeMetatable(luaCoinBaseTypeName))
	L.Push(ud)

	return 1
}

// Checks whether the first lua argument is a *LUserData with *CoinBase and
// returns this *CoinBase.
func checkCoinBase(L *lua.LState, idx int) *payload.CoinBase {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.CoinBase); ok {
		return v
	}
	L.ArgError(1, "CoinBase expected")
	return nil
}

var coinbaseMethods = map[string]lua.LGFunction{
	"get": coinbaseGet,
}

// Getter and setter for the Person#Name
func coinbaseGet(L *lua.LState) int {
	p := checkCoinBase(L, 1)
	fmt.Println(p)

	return 0
}

// Registers my person type to given L.
func RegisterTransferAssetType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaTransferAssetTypeName)
	L.SetGlobal("transferasset", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newTransferAsset))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), transferassetMethods))
}

// Constructor
func newTransferAsset(L *lua.LState) int {
	ta := &payload.TransferAsset{}
	ud := L.NewUserData()
	ud.Value = ta
	L.SetMetatable(ud, L.GetTypeMetatable(luaTransferAssetTypeName))
	L.Push(ud)

	return 1
}

// Checks whether the first lua argument is a *LUserData with *TransferAsset and
// returns this *TransferAsset.
func checkTransferAsset(L *lua.LState, idx int) *payload.TransferAsset {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.TransferAsset); ok {
		return v
	}
	L.ArgError(1, "TransferAsset expected")
	return nil
}

var transferassetMethods = map[string]lua.LGFunction{
	"get": transferassetGet,
}

// Getter and setter for the Person#Name
func transferassetGet(L *lua.LState) int {
	p := checkTransferAsset(L, 1)
	fmt.Println(p)

	return 0
}

func RegisterUpdateProducerType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaUpdateProducerName)
	L.SetGlobal("updateproducer", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newUpdateProducer))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), updateProducerMethods))
}

// Constructor
func newUpdateProducer(L *lua.LState) int {
	ownerPublicKeyStr := L.ToString(1)
	nodePublicKeyStr := L.ToString(2)
	nickName := L.ToString(3)
	url := L.ToString(4)
	location := L.ToInt64(5)
	address := L.ToString(6)
	needSign := true
	client, err := checkClient(L, 7)
	if err != nil {
		needSign = false
	}

	ownerPublicKey, err := common.HexStringToBytes(ownerPublicKeyStr)
	if err != nil {
		fmt.Println("wrong producer public key")
		os.Exit(1)
	}
	nodePublicKey, err := common.HexStringToBytes(nodePublicKeyStr)
	if err != nil {
		fmt.Println("wrong producer public key")
		os.Exit(1)
	}
	updateProducer := &payload.ProducerInfo{
		OwnerPublicKey: []byte(ownerPublicKey),
		NodePublicKey:  []byte(nodePublicKey),
		NickName:       nickName,
		Url:            url,
		Location:       uint64(location),
		NetAddress:     address,
	}

	if needSign {
		upSignBuf := new(bytes.Buffer)
		err = updateProducer.SerializeUnsigned(upSignBuf, payload.ProducerInfoVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		codeHash, err := contract.PublicKeyToStandardCodeHash(ownerPublicKey)
		acc := client.GetAccountByCodeHash(*codeHash)
		if acc == nil {
			fmt.Println("no available account in wallet")
			os.Exit(1)
		}
		rpSig, err := crypto.Sign(acc.PrivKey(), upSignBuf.Bytes())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		updateProducer.Signature = rpSig
	}

	ud := L.NewUserData()
	ud.Value = updateProducer
	L.SetMetatable(ud, L.GetTypeMetatable(luaUpdateProducerName))
	L.Push(ud)

	return 1
}

func checkUpdateProducer(L *lua.LState, idx int) *payload.ProducerInfo {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.ProducerInfo); ok {
		return v
	}
	L.ArgError(1, "ProducerInfo expected")
	return nil
}

var updateProducerMethods = map[string]lua.LGFunction{
	"get": updateProducerGet,
}

// Getter and setter for the Person#Name
func updateProducerGet(L *lua.LState) int {
	p := checkUpdateProducer(L, 1)
	fmt.Println(p)

	return 0
}

// Registers my person type to given L.
func RegisterRegisterProducerType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaRegisterProducerName)
	L.SetGlobal("registerproducer", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newRegisterProducer))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), registerProducerMethods))
}

// Constructor
func newRegisterProducer(L *lua.LState) int {
	ownerPublicKeyStr := L.ToString(1)
	nodePublicKeyStr := L.ToString(2)
	nickName := L.ToString(3)
	url := L.ToString(4)
	location := L.ToInt64(5)
	address := L.ToString(6)
	needSign := true
	client, err := checkClient(L, 7)
	if err != nil {
		needSign = false
	}

	ownerPublicKey, err := common.HexStringToBytes(ownerPublicKeyStr)
	if err != nil {
		fmt.Println("wrong producer public key")
		os.Exit(1)
	}
	nodePublicKey, err := common.HexStringToBytes(nodePublicKeyStr)
	if err != nil {
		fmt.Println("wrong producer public key")
		os.Exit(1)
	}

	registerProducer := &payload.ProducerInfo{
		OwnerPublicKey: []byte(ownerPublicKey),
		NodePublicKey:  []byte(nodePublicKey),
		NickName:       nickName,
		Url:            url,
		Location:       uint64(location),
		NetAddress:     address,
	}

	if needSign {
		rpSignBuf := new(bytes.Buffer)
		err = registerProducer.SerializeUnsigned(rpSignBuf, payload.ProducerInfoVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		codeHash, err := contract.PublicKeyToStandardCodeHash(ownerPublicKey)
		acc := client.GetAccountByCodeHash(*codeHash)
		if acc == nil {
			fmt.Println("no available account in wallet")
			os.Exit(1)
		}
		rpSig, err := crypto.Sign(acc.PrivKey(), rpSignBuf.Bytes())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		registerProducer.Signature = rpSig
	}

	ud := L.NewUserData()
	ud.Value = registerProducer
	L.SetMetatable(ud, L.GetTypeMetatable(luaRegisterProducerName))
	L.Push(ud)

	return 1
}

// Checks whether the first lua argument is a *LUserData with *ProducerInfo and
// returns this *ProducerInfo.
func checkRegisterProducer(L *lua.LState, idx int) *payload.ProducerInfo {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.ProducerInfo); ok {
		return v
	}
	L.ArgError(1, "ProducerInfo expected")
	return nil
}

var registerProducerMethods = map[string]lua.LGFunction{
	"get": registerProducerGet,
}

// Getter and setter for the Person#Name
func registerProducerGet(L *lua.LState) int {
	p := checkRegisterProducer(L, 1)
	fmt.Println(p)

	return 0
}

func RegisterCancelProducerType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaCancelProducerName)
	L.SetGlobal("cancelproducer", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newProcessProducer))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), cancelProducerMethods))
}

// Constructor
func newProcessProducer(L *lua.LState) int {
	publicKeyStr := L.ToString(1)
	client, err := checkClient(L, 2)
	if err != nil {
		fmt.Println(err)
	}

	publicKey, err := common.HexStringToBytes(publicKeyStr)
	if err != nil {
		fmt.Println("wrong producer public key")
		os.Exit(1)
	}
	processProducer := &payload.ProcessProducer{
		OwnerPublicKey: []byte(publicKey),
	}

	cpSignBuf := new(bytes.Buffer)
	err = processProducer.SerializeUnsigned(cpSignBuf, payload.ProcessProducerVersion)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	codeHash, err := contract.PublicKeyToStandardCodeHash(publicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	acc := client.GetAccountByCodeHash(*codeHash)
	if acc == nil {
		fmt.Println("no available account in wallet")
		os.Exit(1)
	}
	rpSig, err := crypto.Sign(acc.PrivKey(), cpSignBuf.Bytes())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	processProducer.Signature = rpSig

	ud := L.NewUserData()
	ud.Value = processProducer
	L.SetMetatable(ud, L.GetTypeMetatable(luaCancelProducerName))
	L.Push(ud)

	return 1
}

func checkCancelProducer(L *lua.LState, idx int) *payload.ProcessProducer {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.ProcessProducer); ok {
		return v
	}
	L.ArgError(1, "CancelProducer expected")
	return nil
}

var cancelProducerMethods = map[string]lua.LGFunction{
	"get": cancelProducerGet,
}

// Getter and setter for the Person#Name
func cancelProducerGet(L *lua.LState) int {
	p := checkCancelProducer(L, 1)
	fmt.Println(p)

	return 0
}

func RegisterReturnDepositCoinType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaReturnDepositCoinName)
	L.SetGlobal("returndepositcoin", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newReturnDepositCoin))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), returnDepositCoinMethods))
}

// Constructor
func newReturnDepositCoin(L *lua.LState) int {
	returnDeposit := &payload.ReturnDepositCoin{}
	ud := L.NewUserData()
	ud.Value = returnDeposit
	L.SetMetatable(ud, L.GetTypeMetatable(luaReturnDepositCoinName))
	L.Push(ud)

	return 1
}

// Checks whether the first lua argument is a *LUserData with *ReturnDepositCoin and
// returns this *ReturnDepositCoin.
func checkReturnDepositCoin(L *lua.LState, idx int) *payload.ReturnDepositCoin {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.ReturnDepositCoin); ok {
		return v
	}
	L.ArgError(1, "ReturnDepositCoin expected")
	return nil
}

var returnDepositCoinMethods = map[string]lua.LGFunction{
	"get": returnDepositCoinGet,
}

// Getter and setter for the Person#Name
func returnDepositCoinGet(L *lua.LState) int {
	p := checkReturnDepositCoin(L, 1)
	fmt.Println(p)

	return 0
}

func RegisterActivateProducerType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaActivateProducerName)
	L.SetGlobal("activateproducer", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newActivateProducer))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), activateProducerMethods))
}

func newActivateProducer(L *lua.LState) int {
	publicKeyStr := L.ToString(1)
	client, err := checkClient(L, 2)
	if err != nil {
		fmt.Println(err)
	}

	publicKey, err := common.HexStringToBytes(publicKeyStr)
	if err != nil {
		fmt.Println("wrong producer node public key")
		os.Exit(1)
	}
	activateProducer := &payload.ActivateProducer{
		NodePublicKey: []byte(publicKey),
	}

	apSignBuf := new(bytes.Buffer)
	err = activateProducer.SerializeUnsigned(apSignBuf, payload.ActivateProducerVersion)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	codeHash, err := contract.PublicKeyToStandardCodeHash(publicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	acc := client.GetAccountByCodeHash(*codeHash)
	if err != nil {
		fmt.Println(err)
	}
	if acc == nil {
		fmt.Println("no available account in wallet")
		os.Exit(1)
	}
	rpSig, err := crypto.Sign(acc.PrivKey(), apSignBuf.Bytes())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	activateProducer.Signature = rpSig

	ud := L.NewUserData()
	ud.Value = activateProducer
	L.SetMetatable(ud, L.GetTypeMetatable(luaActivateProducerName))
	L.Push(ud)

	return 1
}

func checkActivateProducer(L *lua.LState, idx int) *payload.ActivateProducer {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.ActivateProducer); ok {
		return v
	}
	L.ArgError(1, "ActivateProducer expected")
	return nil
}

var activateProducerMethods = map[string]lua.LGFunction{
	"get": activateProducerGet,
}

// Getter and setter for the Person#Name
func activateProducerGet(L *lua.LState) int {
	p := checkActivateProducer(L, 1)
	fmt.Println(p)

	return 0
}

func RegisterSidechainPowType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaSideChainPowName)
	L.SetGlobal("sidechainpow", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newSideChainPow))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), returnSideChainPowMethods))
}

// Constructor
func newSideChainPow(L *lua.LState) int {
	sideBlockHashStr := L.ToString(1)
	sideGenesisHashStr := L.ToString(2)
	blockHeight := L.ToInt(3)
	client, err := checkClient(L, 4)
	if err != nil {
		fmt.Println(err)
	}

	sideBlockHash, err := common.Uint256FromHexString(sideBlockHashStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sideGenesisHash, err := common.Uint256FromHexString(sideGenesisHashStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sideChainPow := &payload.SideChainPow{
		SideBlockHash:   *sideBlockHash,
		SideGenesisHash: *sideGenesisHash,
		BlockHeight:     uint32(blockHeight),
	}

	spSignBuf := new(bytes.Buffer)
	err = sideChainPow.SerializeUnsigned(spSignBuf, payload.SideChainPowVersion)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	acc := client.GetMainAccount()
	spSig, err := crypto.Sign(acc.PrivKey(), spSignBuf.Bytes())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sideChainPow.Signature = spSig

	ud := L.NewUserData()
	ud.Value = sideChainPow
	L.SetMetatable(ud, L.GetTypeMetatable(luaSideChainPowName))
	L.Push(ud)

	return 1
}

func checkSideChainPow(L *lua.LState, idx int) *payload.SideChainPow {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.SideChainPow); ok {
		return v
	}
	L.ArgError(1, "SideChainPow expected")
	return nil
}

var returnSideChainPowMethods = map[string]lua.LGFunction{
	"get": returnSideChainPowGet,
}

// Getter and setter for the Person#Name
func returnSideChainPowGet(L *lua.LState) int {
	p := checkSideChainPow(L, 1)
	fmt.Println(p)

	return 0
}

// Registers my person type to given L.
func RegisterRegisterCRType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaRegisterCRName)
	L.SetGlobal("registercr", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newRegisterCR))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), registerCRMethods))
}

// Constructor
func newRegisterCR(L *lua.LState) int {
	publicKeyStr := L.ToString(1)
	nickName := L.ToString(2)
	url := L.ToString(3)
	location := L.ToInt64(4)
	payloadVersion := byte(L.ToInt(5))
	needSign := true
	client, err := checkClient(L, 6)
	if err != nil {
		needSign = false
	}
	publicKey, err := common.HexStringToBytes(publicKeyStr)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	pk, err := crypto.DecodePoint(publicKey)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	code, err := contract.CreateStandardRedeemScript(pk)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	ct, err := contract.CreateCRIDContractByCode(code)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	didCode := make([]byte, len(code))
	copy(didCode, code)
	didCode = append(didCode[:len(code)-1], common.DID)
	didCT, err := contract.CreateCRIDContractByCode(didCode)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	registerCR := &payload.CRInfo{
		Code:     code,
		CID:      *ct.ToProgramHash(),
		DID:      *didCT.ToProgramHash(),
		NickName: nickName,
		Url:      url,
		Location: uint64(location),
	}

	if needSign {
		rpSignBuf := new(bytes.Buffer)
		err = registerCR.SerializeUnsigned(rpSignBuf, payloadVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		codeHash, err := contract.PublicKeyToStandardCodeHash(publicKey)
		acc := client.GetAccountByCodeHash(*codeHash)
		if acc == nil {
			fmt.Println("no available account in wallet")
			os.Exit(1)
		}
		rpSig, err := crypto.Sign(acc.PrivKey(), rpSignBuf.Bytes())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		registerCR.Signature = rpSig
	}

	ud := L.NewUserData()
	ud.Value = registerCR
	L.SetMetatable(ud, L.GetTypeMetatable(luaRegisterCRName))
	L.Push(ud)

	return 1
}

// Checks whether the first lua argument is a *LUserData with *CRInfo and
// returns this *CRInfo.
func checkRegisterCR(L *lua.LState, idx int) *payload.CRInfo {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.CRInfo); ok {
		return v
	}
	L.ArgError(1, "ProducerInfo expected")
	return nil
}

var registerCRMethods = map[string]lua.LGFunction{
	"get": registerCRGet,
}

// Getter and setter for the Person#Name
func registerCRGet(L *lua.LState) int {
	p := checkRegisterCR(L, 1)
	fmt.Println(p)

	return 0
}

//
// Registers my person type to given L.
func RegisterUpdateCRType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaUpdateCRName)
	L.SetGlobal("updatecr", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newUpdateCR))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), updateCRMethods))
}

// Constructor
func newUpdateCR(L *lua.LState) int {
	publicKeyStr := L.ToString(1)
	nickName := L.ToString(2)
	url := L.ToString(3)
	location := L.ToInt64(4)
	payloadVersion := byte(L.ToInt(5))
	needSign := true
	client, err := checkClient(L, 6)
	if err != nil {
		needSign = false
	}
	publicKey, err := common.HexStringToBytes(publicKeyStr)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	pk, err := crypto.DecodePoint(publicKey)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	code, err := contract.CreateStandardRedeemScript(pk)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	ct, err := contract.CreateCRIDContractByCode(code)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	didCode := make([]byte, len(code))
	copy(didCode, code)
	didCode = append(didCode[:len(code)-1], common.DID)
	didCT, err := contract.CreateCRIDContractByCode(didCode)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	updateCR := &payload.CRInfo{
		Code:     ct.Code,
		CID:      *ct.ToProgramHash(),
		DID:      *didCT.ToProgramHash(),
		NickName: nickName,
		Url:      url,
		Location: uint64(location),
	}

	if needSign {
		rpSignBuf := new(bytes.Buffer)
		err = updateCR.SerializeUnsigned(rpSignBuf, payloadVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		codeHash, err := contract.PublicKeyToStandardCodeHash(publicKey)
		acc := client.GetAccountByCodeHash(*codeHash)
		if acc == nil {
			fmt.Println("no available account in wallet")
			os.Exit(1)
		}
		rpSig, err := crypto.Sign(acc.PrivKey(), rpSignBuf.Bytes())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		updateCR.Signature = rpSig
	}

	ud := L.NewUserData()
	ud.Value = updateCR
	L.SetMetatable(ud, L.GetTypeMetatable(luaUpdateCRName))
	L.Push(ud)

	return 1
}

// Checks whether the first lua argument is a *LUserData with *CRInfo and
// returns this *CRInfo.
func checkUpdateCR(L *lua.LState, idx int) *payload.CRInfo {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.CRInfo); ok {
		return v
	}
	L.ArgError(1, "CRInfo expected")
	return nil
}

var updateCRMethods = map[string]lua.LGFunction{
	"get": updateCRGet,
}

// Getter and setter for the Person#Name
func updateCRGet(L *lua.LState) int {
	p := checkUpdateCR(L, 1)
	fmt.Println(p)

	return 0
}

func RegisterUnregisterCRType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaUnregisterCRName)
	L.SetGlobal("unregistercr", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newUnregisterCR))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), unregisterCRMethods))
}
func getIDProgramHash(code []byte) *common.Uint168 {
	ct, _ := contract.CreateCRIDContractByCode(code)
	return ct.ToProgramHash()
}

// Constructor
func newUnregisterCR(L *lua.LState) int {
	publicKeyStr := L.ToString(1)
	needSign := true
	client, err := checkClient(L, 2)
	if err != nil {
		needSign = false
	}
	publicKey, err := common.HexStringToBytes(publicKeyStr)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	pk, err := crypto.DecodePoint(publicKey)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}

	ct, err := contract.CreateStandardContract(pk)
	if err != nil {
		fmt.Println("wrong cr public key")
		os.Exit(1)
	}
	cid := getIDProgramHash(ct.Code)
	unregisterCR := &payload.UnregisterCR{
		CID: *cid,
	}

	if needSign {
		rpSignBuf := new(bytes.Buffer)
		err = unregisterCR.SerializeUnsigned(rpSignBuf, payload.UnregisterCRVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		codeHash, err := contract.PublicKeyToStandardCodeHash(publicKey)
		acc := client.GetAccountByCodeHash(*codeHash)
		if acc == nil {
			fmt.Println("no available account in wallet")
			os.Exit(1)
		}
		rpSig, err := crypto.Sign(acc.PrivKey(), rpSignBuf.Bytes())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		unregisterCR.Signature = rpSig
	}

	ud := L.NewUserData()
	ud.Value = unregisterCR
	L.SetMetatable(ud, L.GetTypeMetatable(luaUnregisterCRName))
	L.Push(ud)

	return 1
}

// Checks whether the first lua argument is a *LUserData with *CRInfo and
// returns this *CRInfo.
func checkUnregisterCR(L *lua.LState, idx int) *payload.UnregisterCR {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*payload.UnregisterCR); ok {
		return v
	}
	L.ArgError(1, "UnregisterCR expected")
	return nil
}

var unregisterCRMethods = map[string]lua.LGFunction{
	"get": unregisterCRGet,
}

// Getter and setter for the Person#Name
func unregisterCRGet(L *lua.LState) int {
	p := checkUnregisterCR(L, 1)
	fmt.Println(p)

	return 0
}
