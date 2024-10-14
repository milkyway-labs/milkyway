package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	icacallbacktypes "github.com/milkyway-labs/milkyway/x/icacallbacks/types"
	recordstypes "github.com/milkyway-labs/milkyway/x/records/types"
)

func FilterDepositRecords(arr []recordstypes.DepositRecord, condition func(recordstypes.DepositRecord) bool) (ret []recordstypes.DepositRecord) {
	for _, elem := range arr {
		if condition(elem) {
			ret = append(ret, elem)
		}
	}
	return ret
}

func Int64ToCoinString(amount int64, denom string) string {
	return strconv.FormatInt(amount, 10) + denom
}

func ValidateAdminAddress(address string) error {
	if !Admins[address] {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "address (%s) is not an admin", address)
	}
	return nil
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func StringMapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func Int32MapKeys[V any](m map[int32]V) []int32 {
	keys := make([]int32, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

// ==============================  ADDRESS VERIFICATION UTILS  ================================
// ref: https://github.com/cosmos/cosmos-sdk/blob/b75c2ebcfab1a6b535723f1ac2889a2fc2509520/types/address.go#L177

var errBech32EmptyAddress = errors.New("decoding Bech32 address failed: must provide a non empty address")

// GetFromBech32 decodes a bytestring from a Bech32 encoded string.
func GetFromBech32(bech32str, prefix string) ([]byte, error) {
	if len(bech32str) == 0 {
		return nil, errBech32EmptyAddress
	}

	hrp, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return nil, err
	}

	if hrp != prefix {
		return nil, fmt.Errorf("invalid Bech32 prefix; expected %s, got %s", prefix, hrp)
	}

	return bz, nil
}

// VerifyAddressFormat verifies that the provided bytes form a valid address
// according to the default address rules or a custom address verifier set by
// GetConfig().SetAddressVerifier().
// TODO make an issue to get rid of global Config
// ref: https://github.com/cosmos/cosmos-sdk/issues/9690
func VerifyAddressFormat(bz []byte) error {
	verifier := func(bz []byte) error {
		n := len(bz)
		// Base accounts are length 20, module/ICA accounts are length 32
		if n == 20 || n == 32 {
			return nil
		}
		return fmt.Errorf("incorrect address length %d", n)
	}
	if verifier(bz) != nil {
		return verifier(bz)
	}

	if len(bz) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrUnknownAddress, "addresses cannot be empty")
	}

	if len(bz) > address.MaxAddrLen {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "address max length is %d, got %d", address.MaxAddrLen, len(bz))
	}

	return nil
}

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress []byte

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromBech32(address string, bech32prefix string) (addr AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return AccAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := GetFromBech32(address, bech32prefix)
	if err != nil {
		return nil, err
	}

	err = VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return AccAddress(bz), nil
}

// LogWithHostZone returns a log string with a chainId and tab as the prefix
// Ex:
//
//	| COSMOSHUB-4   |   string
func LogWithHostZone(chainID string, s string, a ...any) string {
	msg := fmt.Sprintf(s, a...)
	return fmt.Sprintf("|   %-13s |  %s", strings.ToUpper(chainID), msg)
}

// Returns a log string with a chain Id and callback as a prefix
// callbackType is either ICACALLBACK or ICQCALLBACK
// Format:
//
//	|   CHAIN-ID    |  {CALLBACK_ID} {CALLBACK_TYPE}  |  string
func logCallbackWithHostZone(chainID string, callbackID string, callbackType string, s string, a ...any) string {
	msg := fmt.Sprintf(s, a...)
	return fmt.Sprintf("|   %-13s |  %s %s  |  %s", strings.ToUpper(chainID), strings.ToUpper(callbackID), callbackType, msg)
}

// LogICACallbackWithHostZone returns a log string with a chain Id and icacallback as a prefix
// Ex:
//
//	| COSMOSHUB-4   |  DELEGATE ICACALLBACK  |  string
func LogICACallbackWithHostZone(chainID string, callbackID string, s string, a ...any) string {
	return logCallbackWithHostZone(chainID, callbackID, "ICACALLBACK", s, a...)
}

// LogICACallbackStatusWithHostZone returns a log string with a chain Id and icacallback as a prefix, and status of the callback
// Ex:
//
//	| COSMOSHUB-4   |  DELEGATE ICACALLBACK  |  ICA SUCCESS, Packet: ...
func LogICACallbackStatusWithHostZone(chainID string, callbackID string, status icacallbacktypes.AckResponseStatus, packet channeltypes.Packet) string {
	var statusMsg string
	switch status {
	case icacallbacktypes.AckResponseStatus_SUCCESS:
		statusMsg = "ICA SUCCESSFUL"
	case icacallbacktypes.AckResponseStatus_TIMEOUT:
		statusMsg = "ICA TIMEOUT"
	default:
		statusMsg = "ICA FAILED (ack error)"
	}
	return logCallbackWithHostZone(chainID, callbackID, "ICACALLBACK", "%s, Packet: %+v", statusMsg, packet)
}

// LogICQCallbackWithHostZone returns a log string with a chain Id and icqcallback as a prefix
// Ex:
//
//	| COSMOSHUB-4   |  WITHDRAWALHOSTBALANCE ICQCALLBACK  |  string
func LogICQCallbackWithHostZone(chainID string, callbackID string, s string, a ...any) string {
	return logCallbackWithHostZone(chainID, callbackID, "ICQCALLBACK", s, a...)
}

// LogHeader returns a log header string with a dash padding on either side
// Ex:
//
//	------------------------------ string ------------------------------
func LogHeader(s string, a ...any) string {
	lineLength := 120
	header := fmt.Sprintf(s, a...)
	pad := strings.Repeat("-", (lineLength-len(header))/2)
	return fmt.Sprintf("%s %s %s", pad, header, pad)
}

// StAssetDenomFromHostZoneDenom returns the stDenom from a native denom by appending a st prefix
func StAssetDenomFromHostZoneDenom(hostZoneDenom string) string {
	return "milk" + hostZoneDenom
}

// VerifyTxHash verifies a tx hash is valid
func VerifyTxHash(txHash string) (err error) {
	if txHash == "" {
		return errorsmod.Wrapf(sdkerrors.ErrTxDecode, "tx hash is empty")
	}
	_, err = hex.DecodeString(txHash)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrTxDecode, "tx hash is invalid %s", txHash)
	}
	return nil
}

func ParseUint32Slice(s string) ([]uint32, error) {
	if s == "" {
		return nil, nil
	}
	ss := strings.Split(s, ",")
	out := make([]uint32, len(ss))
	for i, d := range ss {
		u, err := strconv.ParseUint(d, 10, 32)
		if err != nil {
			return nil, err
		}
		out[i] = uint32(u)
	}
	return out, nil
}

func FormatUint32Slice(s []uint32) string {
	ss := make([]string, len(s))
	for i, u := range s {
		ss[i] = fmt.Sprint(u)
	}
	return strings.Join(ss, ",")
}

func MustParseCoin(s string) sdk.Coin {
	c, err := sdk.ParseCoinNormalized(strings.ReplaceAll(s, "_", ""))
	if err != nil {
		panic(err)
	}
	return c
}

func MustParseCoins(s string) sdk.Coins {
	c, err := sdk.ParseCoinsNormalized(strings.ReplaceAll(s, "_", ""))
	if err != nil {
		panic(err)
	}
	return c
}

func MustParseDecCoins(s string) sdk.DecCoins {
	d, err := sdk.ParseDecCoins(strings.ReplaceAll(s, "_", ""))
	if err != nil {
		panic(err)
	}
	return d
}

func MustParseDec(s string) sdkmath.LegacyDec {
	return sdkmath.LegacyMustNewDecFromStr(strings.ReplaceAll(s, "_", ""))
}
