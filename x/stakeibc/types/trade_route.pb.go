// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: stride/stakeibc/trade_route.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Stores pool information needed to execute the swap along a trade route
type TradeConfig struct {
	// Currently Osmosis is the only trade chain so this is an osmosis pool id
	PoolId uint64 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	// Spot price in the pool to convert the reward denom to the host denom
	// output_tokens = swap_price * input tokens
	// This value may be slightly stale as it is updated by an ICQ
	SwapPrice cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=swap_price,json=swapPrice,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"swap_price"`
	// unix time in seconds that the price was last updated
	PriceUpdateTimestamp uint64 `protobuf:"varint,3,opt,name=price_update_timestamp,json=priceUpdateTimestamp,proto3" json:"price_update_timestamp,omitempty"`
	// Threshold defining the percentage of tokens that could be lost in the trade
	// This captures both the loss from slippage and from a stale price on stride
	// 0.05 means the output from the trade can be no less than a 5% deviation
	// from the current value
	MaxAllowedSwapLossRate cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=max_allowed_swap_loss_rate,json=maxAllowedSwapLossRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"max_allowed_swap_loss_rate"`
	// min and max set boundaries of reward denom on trade chain we will swap
	// min also decides when reward token transfers are worth it (transfer fees)
	MinSwapAmount cosmossdk_io_math.Int `protobuf:"bytes,5,opt,name=min_swap_amount,json=minSwapAmount,proto3,customtype=cosmossdk.io/math.Int" json:"min_swap_amount"`
	MaxSwapAmount cosmossdk_io_math.Int `protobuf:"bytes,6,opt,name=max_swap_amount,json=maxSwapAmount,proto3,customtype=cosmossdk.io/math.Int" json:"max_swap_amount"`
}

func (m *TradeConfig) Reset()         { *m = TradeConfig{} }
func (m *TradeConfig) String() string { return proto.CompactTextString(m) }
func (*TradeConfig) ProtoMessage()    {}
func (*TradeConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_c252b142ecf88017, []int{0}
}
func (m *TradeConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TradeConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TradeConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TradeConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TradeConfig.Merge(m, src)
}
func (m *TradeConfig) XXX_Size() int {
	return m.Size()
}
func (m *TradeConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_TradeConfig.DiscardUnknown(m)
}

var xxx_messageInfo_TradeConfig proto.InternalMessageInfo

func (m *TradeConfig) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *TradeConfig) GetPriceUpdateTimestamp() uint64 {
	if m != nil {
		return m.PriceUpdateTimestamp
	}
	return 0
}

// TradeRoute represents a round trip including info on transfer and how to do
// the swap. It makes the assumption that the reward token is always foreign to
// the host so therefore the first two hops are to unwind the ibc denom enroute
// to the trade chain and the last hop is the return so funds start/end in the
// withdrawl ICA on hostZone
// The structure is key'd on reward denom and host denom in their native forms
// (i.e. reward_denom_on_reward_zone and host_denom_on_host_zone)
type TradeRoute struct {
	// ibc denom for the reward on the host zone
	RewardDenomOnHostZone string `protobuf:"bytes,1,opt,name=reward_denom_on_host_zone,json=rewardDenomOnHostZone,proto3" json:"reward_denom_on_host_zone,omitempty"`
	// should be the native denom for the reward chain
	RewardDenomOnRewardZone string `protobuf:"bytes,2,opt,name=reward_denom_on_reward_zone,json=rewardDenomOnRewardZone,proto3" json:"reward_denom_on_reward_zone,omitempty"`
	// ibc denom of the reward on the trade chain, input to the swap
	RewardDenomOnTradeZone string `protobuf:"bytes,3,opt,name=reward_denom_on_trade_zone,json=rewardDenomOnTradeZone,proto3" json:"reward_denom_on_trade_zone,omitempty"`
	// ibc of the host denom on the trade chain, output from the swap
	HostDenomOnTradeZone string `protobuf:"bytes,4,opt,name=host_denom_on_trade_zone,json=hostDenomOnTradeZone,proto3" json:"host_denom_on_trade_zone,omitempty"`
	// should be the same as the native host denom on the host chain
	HostDenomOnHostZone string `protobuf:"bytes,5,opt,name=host_denom_on_host_zone,json=hostDenomOnHostZone,proto3" json:"host_denom_on_host_zone,omitempty"`
	// ICAAccount on the host zone with the reward tokens
	// This is the same as the host zone withdrawal ICA account
	HostAccount ICAAccount `protobuf:"bytes,6,opt,name=host_account,json=hostAccount,proto3" json:"host_account"`
	// ICAAccount on the reward zone that is acts as the intermediate
	// receiver of the transfer from host zone to trade zone
	RewardAccount ICAAccount `protobuf:"bytes,7,opt,name=reward_account,json=rewardAccount,proto3" json:"reward_account"`
	// ICAAccount responsible for executing the swap of reward
	// tokens for host tokens
	TradeAccount ICAAccount `protobuf:"bytes,8,opt,name=trade_account,json=tradeAccount,proto3" json:"trade_account"`
	// Channel responsible for the transfer of reward tokens from the host
	// zone to the reward zone. This is the channel ID on the host zone side
	HostToRewardChannelId string `protobuf:"bytes,9,opt,name=host_to_reward_channel_id,json=hostToRewardChannelId,proto3" json:"host_to_reward_channel_id,omitempty"`
	// Channel responsible for the transfer of reward tokens from the reward
	// zone to the trade zone. This is the channel ID on the reward zone side
	RewardToTradeChannelId string `protobuf:"bytes,10,opt,name=reward_to_trade_channel_id,json=rewardToTradeChannelId,proto3" json:"reward_to_trade_channel_id,omitempty"`
	// Channel responsible for the transfer of host tokens from the trade
	// zone, back to the host zone. This is the channel ID on the trade zone side
	TradeToHostChannelId string `protobuf:"bytes,11,opt,name=trade_to_host_channel_id,json=tradeToHostChannelId,proto3" json:"trade_to_host_channel_id,omitempty"`
	// specifies the configuration needed to execute the swap
	// such as pool_id, slippage, min trade amount, etc.
	TradeConfig TradeConfig `protobuf:"bytes,12,opt,name=trade_config,json=tradeConfig,proto3" json:"trade_config"`
}

func (m *TradeRoute) Reset()         { *m = TradeRoute{} }
func (m *TradeRoute) String() string { return proto.CompactTextString(m) }
func (*TradeRoute) ProtoMessage()    {}
func (*TradeRoute) Descriptor() ([]byte, []int) {
	return fileDescriptor_c252b142ecf88017, []int{1}
}
func (m *TradeRoute) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TradeRoute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TradeRoute.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TradeRoute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TradeRoute.Merge(m, src)
}
func (m *TradeRoute) XXX_Size() int {
	return m.Size()
}
func (m *TradeRoute) XXX_DiscardUnknown() {
	xxx_messageInfo_TradeRoute.DiscardUnknown(m)
}

var xxx_messageInfo_TradeRoute proto.InternalMessageInfo

func (m *TradeRoute) GetRewardDenomOnHostZone() string {
	if m != nil {
		return m.RewardDenomOnHostZone
	}
	return ""
}

func (m *TradeRoute) GetRewardDenomOnRewardZone() string {
	if m != nil {
		return m.RewardDenomOnRewardZone
	}
	return ""
}

func (m *TradeRoute) GetRewardDenomOnTradeZone() string {
	if m != nil {
		return m.RewardDenomOnTradeZone
	}
	return ""
}

func (m *TradeRoute) GetHostDenomOnTradeZone() string {
	if m != nil {
		return m.HostDenomOnTradeZone
	}
	return ""
}

func (m *TradeRoute) GetHostDenomOnHostZone() string {
	if m != nil {
		return m.HostDenomOnHostZone
	}
	return ""
}

func (m *TradeRoute) GetHostAccount() ICAAccount {
	if m != nil {
		return m.HostAccount
	}
	return ICAAccount{}
}

func (m *TradeRoute) GetRewardAccount() ICAAccount {
	if m != nil {
		return m.RewardAccount
	}
	return ICAAccount{}
}

func (m *TradeRoute) GetTradeAccount() ICAAccount {
	if m != nil {
		return m.TradeAccount
	}
	return ICAAccount{}
}

func (m *TradeRoute) GetHostToRewardChannelId() string {
	if m != nil {
		return m.HostToRewardChannelId
	}
	return ""
}

func (m *TradeRoute) GetRewardToTradeChannelId() string {
	if m != nil {
		return m.RewardToTradeChannelId
	}
	return ""
}

func (m *TradeRoute) GetTradeToHostChannelId() string {
	if m != nil {
		return m.TradeToHostChannelId
	}
	return ""
}

func (m *TradeRoute) GetTradeConfig() TradeConfig {
	if m != nil {
		return m.TradeConfig
	}
	return TradeConfig{}
}

func init() {
	proto.RegisterType((*TradeConfig)(nil), "stride.stakeibc.TradeConfig")
	proto.RegisterType((*TradeRoute)(nil), "stride.stakeibc.TradeRoute")
}

func init() { proto.RegisterFile("stride/stakeibc/trade_route.proto", fileDescriptor_c252b142ecf88017) }

var fileDescriptor_c252b142ecf88017 = []byte{
	// 663 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xcd, 0x4e, 0xdb, 0x40,
	0x10, 0xc7, 0x93, 0xf2, 0xd5, 0x6c, 0x42, 0x91, 0x5c, 0x3e, 0x4c, 0x68, 0x03, 0xe5, 0xc4, 0x05,
	0x47, 0xa5, 0x08, 0xa1, 0xaa, 0x97, 0x40, 0xa8, 0x88, 0x84, 0x54, 0xe4, 0xa6, 0x17, 0x2e, 0xab,
	0x8d, 0xbd, 0x4d, 0x2c, 0xbc, 0x1e, 0xcb, 0xbb, 0x51, 0x92, 0x3e, 0x45, 0xdf, 0xa3, 0xd7, 0x3e,
	0x04, 0x47, 0xd4, 0x53, 0xd5, 0x03, 0xaa, 0xe0, 0x31, 0x7a, 0xa9, 0x76, 0xd6, 0x06, 0x87, 0xf4,
	0x10, 0xf5, 0x96, 0xc9, 0xcc, 0xef, 0x3f, 0x3b, 0xfb, 0xf7, 0x0e, 0x79, 0x25, 0x55, 0x12, 0xf8,
	0xbc, 0x2e, 0x15, 0xbb, 0xe4, 0x41, 0xc7, 0xab, 0xab, 0x84, 0xf9, 0x9c, 0x26, 0xd0, 0x57, 0xdc,
	0x89, 0x13, 0x50, 0x60, 0x2d, 0x99, 0x12, 0x27, 0x2b, 0xa9, 0x2e, 0x77, 0xa1, 0x0b, 0x98, 0xab,
	0xeb, 0x5f, 0xa6, 0xac, 0x3a, 0xa1, 0x14, 0x78, 0x8c, 0x32, 0xcf, 0x83, 0x7e, 0xa4, 0xd2, 0x92,
	0x75, 0x0f, 0xa4, 0x00, 0x49, 0x0d, 0x6b, 0x02, 0x93, 0xda, 0xfe, 0x36, 0x43, 0xca, 0x6d, 0xdd,
	0xfa, 0x18, 0xa2, 0xcf, 0x41, 0xd7, 0x5a, 0x23, 0x0b, 0x31, 0x40, 0x48, 0x03, 0xdf, 0x2e, 0x6e,
	0x15, 0x77, 0x66, 0xdd, 0x79, 0x1d, 0xb6, 0x7c, 0xeb, 0x9c, 0x10, 0x39, 0x60, 0x31, 0x8d, 0x93,
	0xc0, 0xe3, 0xf6, 0x93, 0xad, 0xe2, 0x4e, 0xe9, 0xe8, 0xf5, 0xd5, 0xcd, 0x66, 0xe1, 0xd7, 0xcd,
	0xe6, 0x86, 0x91, 0x94, 0xfe, 0xa5, 0x13, 0x40, 0x5d, 0x30, 0xd5, 0x73, 0xce, 0x78, 0x97, 0x79,
	0xa3, 0x26, 0xf7, 0x7e, 0x7c, 0xdf, 0x25, 0x69, 0xc7, 0x26, 0xf7, 0xdc, 0x92, 0x16, 0x39, 0xd7,
	0x1a, 0xd6, 0x3e, 0x59, 0x45, 0x31, 0xda, 0x8f, 0x7d, 0xa6, 0x38, 0x55, 0x81, 0xe0, 0x52, 0x31,
	0x11, 0xdb, 0x33, 0xd8, 0x79, 0x19, 0xb3, 0x9f, 0x30, 0xd9, 0xce, 0x72, 0x96, 0x20, 0x55, 0xc1,
	0x86, 0x94, 0x85, 0x21, 0x0c, 0xb8, 0x4f, 0xf1, 0x4c, 0x21, 0x48, 0x49, 0x13, 0xa6, 0xb8, 0x3d,
	0xfb, 0xbf, 0xe7, 0x5a, 0x15, 0x6c, 0xd8, 0x30, 0x9a, 0x1f, 0x07, 0x2c, 0x3e, 0x03, 0x29, 0x5d,
	0xa6, 0xb8, 0x75, 0x42, 0x96, 0x44, 0x10, 0x99, 0x36, 0x4c, 0xe8, 0x3b, 0xb5, 0xe7, 0xb0, 0xc7,
	0xcb, 0xb4, 0xc7, 0xca, 0x64, 0x8f, 0x56, 0xa4, 0xdc, 0x45, 0x11, 0x44, 0x5a, 0xa8, 0x81, 0x0c,
	0xca, 0xb0, 0xe1, 0x98, 0xcc, 0xfc, 0x74, 0x32, 0x6c, 0xf8, 0x20, 0xb3, 0xfd, 0x67, 0x8e, 0x10,
	0x74, 0xcb, 0xd5, 0xdf, 0x89, 0x75, 0x48, 0xd6, 0x13, 0x3e, 0x60, 0x89, 0x4f, 0x7d, 0x1e, 0x81,
	0xa0, 0x10, 0xd1, 0x1e, 0x48, 0x45, 0xbf, 0x40, 0xc4, 0xd1, 0xbe, 0x92, 0xbb, 0x62, 0x0a, 0x9a,
	0x3a, 0xff, 0x21, 0x3a, 0x05, 0xa9, 0x2e, 0x20, 0xe2, 0xd6, 0x3b, 0xb2, 0xf1, 0x98, 0x4c, 0x63,
	0x64, 0xd1, 0x5e, 0x77, 0x6d, 0x8c, 0x75, 0x31, 0x40, 0xfa, 0x2d, 0xa9, 0x3e, 0xa6, 0xcd, 0xe7,
	0x8b, 0xf0, 0x0c, 0xc2, 0xab, 0x63, 0x30, 0x1e, 0x1a, 0xd9, 0x03, 0x62, 0xe3, 0x19, 0xff, 0x45,
	0xa2, 0x7b, 0xee, 0xb2, 0xce, 0x4f, 0x70, 0xfb, 0x64, 0x6d, 0x9c, 0x7b, 0x98, 0x14, 0x0d, 0x71,
	0x9f, 0xe7, 0xb0, 0xfb, 0x39, 0x9b, 0xa4, 0x82, 0x75, 0xe9, 0x7b, 0xc0, 0x4b, 0x2f, 0xef, 0x6d,
	0x38, 0x8f, 0x9e, 0x96, 0xd3, 0x3a, 0x6e, 0x34, 0x4c, 0xc9, 0xd1, 0xac, 0x76, 0xc4, 0x2d, 0x6b,
	0x2c, 0xfd, 0xcb, 0x3a, 0x25, 0xcf, 0xd2, 0x79, 0x33, 0x9d, 0x85, 0x69, 0x75, 0x16, 0x0d, 0x98,
	0x29, 0xbd, 0x27, 0x8b, 0x66, 0xde, 0x4c, 0xe8, 0xe9, 0xb4, 0x42, 0x15, 0xe4, 0x32, 0x9d, 0x43,
	0xb2, 0x8e, 0x73, 0x29, 0xc8, 0x7c, 0xf3, 0x7a, 0x2c, 0x8a, 0x38, 0x3e, 0xdc, 0x92, 0x71, 0x5e,
	0x17, 0xb4, 0xc1, 0xd8, 0x76, 0x6c, 0xb2, 0x2d, 0x3f, 0xe7, 0x9d, 0x82, 0xf4, 0xee, 0x73, 0x28,
	0xc9, 0x7b, 0xd7, 0x06, 0xb3, 0x19, 0xee, 0xd9, 0x03, 0x62, 0x1b, 0x42, 0x81, 0xb9, 0xfe, 0x1c,
	0x59, 0x36, 0xde, 0x61, 0xbe, 0x0d, 0xda, 0x80, 0x07, 0xee, 0x84, 0x54, 0xd2, 0x4e, 0xb8, 0x64,
	0xec, 0x0a, 0x0e, 0xfd, 0x62, 0x62, 0xe8, 0xdc, 0x22, 0xca, 0x6c, 0x50, 0xb9, 0xbf, 0xce, 0xae,
	0x6e, 0x6b, 0xc5, 0xeb, 0xdb, 0x5a, 0xf1, 0xf7, 0x6d, 0xad, 0xf8, 0xf5, 0xae, 0x56, 0xb8, 0xbe,
	0xab, 0x15, 0x7e, 0xde, 0xd5, 0x0a, 0x17, 0x7b, 0xdd, 0x40, 0xf5, 0xfa, 0x1d, 0xc7, 0x03, 0x51,
	0x17, 0x41, 0x78, 0x39, 0x1a, 0xb0, 0xd1, 0x6e, 0xc8, 0x3a, 0xf2, 0x3e, 0xaa, 0x0f, 0x73, 0xab,
	0x76, 0x14, 0x73, 0xd9, 0x99, 0xc7, 0x05, 0xf8, 0xe6, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x25,
	0x60, 0xac, 0xff, 0x8a, 0x05, 0x00, 0x00,
}

func (m *TradeConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TradeConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TradeConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.MaxSwapAmount.Size()
		i -= size
		if _, err := m.MaxSwapAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTradeRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.MinSwapAmount.Size()
		i -= size
		if _, err := m.MinSwapAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTradeRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.MaxAllowedSwapLossRate.Size()
		i -= size
		if _, err := m.MaxAllowedSwapLossRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTradeRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.PriceUpdateTimestamp != 0 {
		i = encodeVarintTradeRoute(dAtA, i, uint64(m.PriceUpdateTimestamp))
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.SwapPrice.Size()
		i -= size
		if _, err := m.SwapPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTradeRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.PoolId != 0 {
		i = encodeVarintTradeRoute(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TradeRoute) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TradeRoute) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TradeRoute) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.TradeConfig.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTradeRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	if len(m.TradeToHostChannelId) > 0 {
		i -= len(m.TradeToHostChannelId)
		copy(dAtA[i:], m.TradeToHostChannelId)
		i = encodeVarintTradeRoute(dAtA, i, uint64(len(m.TradeToHostChannelId)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.RewardToTradeChannelId) > 0 {
		i -= len(m.RewardToTradeChannelId)
		copy(dAtA[i:], m.RewardToTradeChannelId)
		i = encodeVarintTradeRoute(dAtA, i, uint64(len(m.RewardToTradeChannelId)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.HostToRewardChannelId) > 0 {
		i -= len(m.HostToRewardChannelId)
		copy(dAtA[i:], m.HostToRewardChannelId)
		i = encodeVarintTradeRoute(dAtA, i, uint64(len(m.HostToRewardChannelId)))
		i--
		dAtA[i] = 0x4a
	}
	{
		size, err := m.TradeAccount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTradeRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size, err := m.RewardAccount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTradeRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size, err := m.HostAccount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTradeRoute(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.HostDenomOnHostZone) > 0 {
		i -= len(m.HostDenomOnHostZone)
		copy(dAtA[i:], m.HostDenomOnHostZone)
		i = encodeVarintTradeRoute(dAtA, i, uint64(len(m.HostDenomOnHostZone)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.HostDenomOnTradeZone) > 0 {
		i -= len(m.HostDenomOnTradeZone)
		copy(dAtA[i:], m.HostDenomOnTradeZone)
		i = encodeVarintTradeRoute(dAtA, i, uint64(len(m.HostDenomOnTradeZone)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.RewardDenomOnTradeZone) > 0 {
		i -= len(m.RewardDenomOnTradeZone)
		copy(dAtA[i:], m.RewardDenomOnTradeZone)
		i = encodeVarintTradeRoute(dAtA, i, uint64(len(m.RewardDenomOnTradeZone)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.RewardDenomOnRewardZone) > 0 {
		i -= len(m.RewardDenomOnRewardZone)
		copy(dAtA[i:], m.RewardDenomOnRewardZone)
		i = encodeVarintTradeRoute(dAtA, i, uint64(len(m.RewardDenomOnRewardZone)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.RewardDenomOnHostZone) > 0 {
		i -= len(m.RewardDenomOnHostZone)
		copy(dAtA[i:], m.RewardDenomOnHostZone)
		i = encodeVarintTradeRoute(dAtA, i, uint64(len(m.RewardDenomOnHostZone)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTradeRoute(dAtA []byte, offset int, v uint64) int {
	offset -= sovTradeRoute(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TradeConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovTradeRoute(uint64(m.PoolId))
	}
	l = m.SwapPrice.Size()
	n += 1 + l + sovTradeRoute(uint64(l))
	if m.PriceUpdateTimestamp != 0 {
		n += 1 + sovTradeRoute(uint64(m.PriceUpdateTimestamp))
	}
	l = m.MaxAllowedSwapLossRate.Size()
	n += 1 + l + sovTradeRoute(uint64(l))
	l = m.MinSwapAmount.Size()
	n += 1 + l + sovTradeRoute(uint64(l))
	l = m.MaxSwapAmount.Size()
	n += 1 + l + sovTradeRoute(uint64(l))
	return n
}

func (m *TradeRoute) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RewardDenomOnHostZone)
	if l > 0 {
		n += 1 + l + sovTradeRoute(uint64(l))
	}
	l = len(m.RewardDenomOnRewardZone)
	if l > 0 {
		n += 1 + l + sovTradeRoute(uint64(l))
	}
	l = len(m.RewardDenomOnTradeZone)
	if l > 0 {
		n += 1 + l + sovTradeRoute(uint64(l))
	}
	l = len(m.HostDenomOnTradeZone)
	if l > 0 {
		n += 1 + l + sovTradeRoute(uint64(l))
	}
	l = len(m.HostDenomOnHostZone)
	if l > 0 {
		n += 1 + l + sovTradeRoute(uint64(l))
	}
	l = m.HostAccount.Size()
	n += 1 + l + sovTradeRoute(uint64(l))
	l = m.RewardAccount.Size()
	n += 1 + l + sovTradeRoute(uint64(l))
	l = m.TradeAccount.Size()
	n += 1 + l + sovTradeRoute(uint64(l))
	l = len(m.HostToRewardChannelId)
	if l > 0 {
		n += 1 + l + sovTradeRoute(uint64(l))
	}
	l = len(m.RewardToTradeChannelId)
	if l > 0 {
		n += 1 + l + sovTradeRoute(uint64(l))
	}
	l = len(m.TradeToHostChannelId)
	if l > 0 {
		n += 1 + l + sovTradeRoute(uint64(l))
	}
	l = m.TradeConfig.Size()
	n += 1 + l + sovTradeRoute(uint64(l))
	return n
}

func sovTradeRoute(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTradeRoute(x uint64) (n int) {
	return sovTradeRoute(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TradeConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTradeRoute
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TradeConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TradeConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SwapPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SwapPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PriceUpdateTimestamp", wireType)
			}
			m.PriceUpdateTimestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PriceUpdateTimestamp |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxAllowedSwapLossRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxAllowedSwapLossRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinSwapAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinSwapAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxSwapAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxSwapAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTradeRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TradeRoute) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTradeRoute
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TradeRoute: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TradeRoute: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardDenomOnHostZone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardDenomOnHostZone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardDenomOnRewardZone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardDenomOnRewardZone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardDenomOnTradeZone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardDenomOnTradeZone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostDenomOnTradeZone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HostDenomOnTradeZone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostDenomOnHostZone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HostDenomOnHostZone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostAccount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.HostAccount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardAccount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardAccount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeAccount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TradeAccount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostToRewardChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HostToRewardChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardToTradeChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardToTradeChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeToHostChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TradeToHostChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradeConfig", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTradeRoute
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TradeConfig.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTradeRoute(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTradeRoute
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTradeRoute(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTradeRoute
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTradeRoute
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTradeRoute
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTradeRoute
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTradeRoute
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTradeRoute        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTradeRoute          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTradeRoute = fmt.Errorf("proto: unexpected end of group")
)