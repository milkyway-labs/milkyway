package keeper

import (
	icacallbackstypes "github.com/milkyway-labs/milk/x/icacallbacks/types"
)

const IBCCallbacksID_NativeTransfer = "transfer"
const IBCCallbacksID_LSMTransfer = "lsm-transfer"

func (k Keeper) Callbacks() icacallbackstypes.ModuleCallbacks {
	return []icacallbackstypes.ICACallback{
		{CallbackId: IBCCallbacksID_NativeTransfer, CallbackFunc: icacallbackstypes.ICACallbackFunction(k.TransferCallback)},
		{CallbackId: IBCCallbacksID_LSMTransfer, CallbackFunc: icacallbackstypes.ICACallbackFunction(k.LSMTransferCallback)},
	}
}
