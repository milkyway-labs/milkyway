// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: stride/stakeibc/ica_account.proto

package types

import (
	fmt "fmt"
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

type ICAAccountType int32

const (
	ICAAccountType_DELEGATION             ICAAccountType = 0
	ICAAccountType_FEE                    ICAAccountType = 1
	ICAAccountType_WITHDRAWAL             ICAAccountType = 2
	ICAAccountType_REDEMPTION             ICAAccountType = 3
	ICAAccountType_COMMUNITY_POOL_DEPOSIT ICAAccountType = 4
	ICAAccountType_COMMUNITY_POOL_RETURN  ICAAccountType = 5
	ICAAccountType_CONVERTER_UNWIND       ICAAccountType = 6
	ICAAccountType_CONVERTER_TRADE        ICAAccountType = 7
)

var ICAAccountType_name = map[int32]string{
	0: "DELEGATION",
	1: "FEE",
	2: "WITHDRAWAL",
	3: "REDEMPTION",
	4: "COMMUNITY_POOL_DEPOSIT",
	5: "COMMUNITY_POOL_RETURN",
	6: "CONVERTER_UNWIND",
	7: "CONVERTER_TRADE",
}

var ICAAccountType_value = map[string]int32{
	"DELEGATION":             0,
	"FEE":                    1,
	"WITHDRAWAL":             2,
	"REDEMPTION":             3,
	"COMMUNITY_POOL_DEPOSIT": 4,
	"COMMUNITY_POOL_RETURN":  5,
	"CONVERTER_UNWIND":       6,
	"CONVERTER_TRADE":        7,
}

func (x ICAAccountType) String() string {
	return proto.EnumName(ICAAccountType_name, int32(x))
}

func (ICAAccountType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2976ae6e7f6ce824, []int{0}
}

type ICAAccount struct {
	ChainId      string         `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Type         ICAAccountType `protobuf:"varint,2,opt,name=type,proto3,enum=stride.stakeibc.ICAAccountType" json:"type,omitempty"`
	ConnectionId string         `protobuf:"bytes,3,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	Address      string         `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *ICAAccount) Reset()         { *m = ICAAccount{} }
func (m *ICAAccount) String() string { return proto.CompactTextString(m) }
func (*ICAAccount) ProtoMessage()    {}
func (*ICAAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_2976ae6e7f6ce824, []int{0}
}
func (m *ICAAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ICAAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ICAAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ICAAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ICAAccount.Merge(m, src)
}
func (m *ICAAccount) XXX_Size() int {
	return m.Size()
}
func (m *ICAAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_ICAAccount.DiscardUnknown(m)
}

var xxx_messageInfo_ICAAccount proto.InternalMessageInfo

func (m *ICAAccount) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

func (m *ICAAccount) GetType() ICAAccountType {
	if m != nil {
		return m.Type
	}
	return ICAAccountType_DELEGATION
}

func (m *ICAAccount) GetConnectionId() string {
	if m != nil {
		return m.ConnectionId
	}
	return ""
}

func (m *ICAAccount) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func init() {
	proto.RegisterEnum("stride.stakeibc.ICAAccountType", ICAAccountType_name, ICAAccountType_value)
	proto.RegisterType((*ICAAccount)(nil), "stride.stakeibc.ICAAccount")
}

func init() { proto.RegisterFile("stride/stakeibc/ica_account.proto", fileDescriptor_2976ae6e7f6ce824) }

var fileDescriptor_2976ae6e7f6ce824 = []byte{
	// 366 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xcf, 0x6e, 0xda, 0x30,
	0x18, 0xc0, 0x63, 0x60, 0x64, 0xb3, 0x36, 0xb0, 0xbc, 0x3f, 0x0a, 0x3b, 0x64, 0x6c, 0xbb, 0xa0,
	0x49, 0x4b, 0x24, 0x78, 0x82, 0x8c, 0x78, 0x9b, 0xa5, 0x90, 0x20, 0xcf, 0x0c, 0x6d, 0x97, 0x28,
	0x71, 0xa2, 0x62, 0x01, 0x09, 0x22, 0x41, 0x6d, 0xde, 0xa2, 0xf7, 0x3e, 0x42, 0x5f, 0xa4, 0x47,
	0x8e, 0x3d, 0x56, 0xf0, 0x22, 0x15, 0xa1, 0x6d, 0x54, 0x8e, 0xbf, 0xef, 0xfb, 0xe9, 0x67, 0xc9,
	0x1f, 0xfc, 0x9c, 0xe5, 0x6b, 0x19, 0xc5, 0x66, 0x96, 0x07, 0xf3, 0x58, 0x86, 0xc2, 0x94, 0x22,
	0xf0, 0x03, 0x21, 0xd2, 0x4d, 0x92, 0x1b, 0xab, 0x75, 0x9a, 0xa7, 0xb8, 0x7d, 0x54, 0x8c, 0x47,
	0xe5, 0xcb, 0x15, 0x80, 0x90, 0x0e, 0x2d, 0xeb, 0x68, 0xe1, 0x0e, 0x7c, 0x29, 0x66, 0x81, 0x4c,
	0x7c, 0x19, 0x69, 0xa0, 0x0b, 0x7a, 0xaf, 0x98, 0x5a, 0x32, 0x8d, 0xf0, 0x00, 0x36, 0xf2, 0x62,
	0x15, 0x6b, 0xb5, 0x2e, 0xe8, 0xb5, 0xfa, 0x9f, 0x8c, 0x93, 0x92, 0x51, 0x55, 0x78, 0xb1, 0x8a,
	0x59, 0x29, 0xe3, 0xaf, 0xf0, 0x8d, 0x48, 0x93, 0x24, 0x16, 0xb9, 0x4c, 0xcb, 0x68, 0xbd, 0x8c,
	0xbe, 0xae, 0x86, 0x34, 0xc2, 0x1a, 0x54, 0x83, 0x28, 0x5a, 0xc7, 0x59, 0xa6, 0x35, 0x8e, 0x6f,
	0x3e, 0xe0, 0xb7, 0x6b, 0x00, 0x5b, 0xcf, 0xbb, 0xb8, 0x05, 0xa1, 0x4d, 0x1c, 0xf2, 0xcb, 0xe2,
	0xd4, 0x73, 0x91, 0x82, 0x55, 0x58, 0xff, 0x49, 0x08, 0x02, 0x87, 0xc5, 0x94, 0xf2, 0xdf, 0x36,
	0xb3, 0xa6, 0x96, 0x83, 0x6a, 0x07, 0x66, 0xc4, 0x26, 0xa3, 0x71, 0x29, 0xd6, 0xf1, 0x47, 0xf8,
	0x61, 0xe8, 0x8d, 0x46, 0x13, 0x97, 0xf2, 0x7f, 0xfe, 0xd8, 0xf3, 0x1c, 0xdf, 0x26, 0x63, 0xef,
	0x0f, 0xe5, 0xa8, 0x81, 0x3b, 0xf0, 0xfd, 0xc9, 0x8e, 0x11, 0x3e, 0x61, 0x2e, 0x7a, 0x81, 0xdf,
	0x41, 0x34, 0xf4, 0xdc, 0xbf, 0x84, 0x71, 0xc2, 0xfc, 0x89, 0x3b, 0xa5, 0xae, 0x8d, 0x9a, 0xf8,
	0x2d, 0x6c, 0x57, 0x53, 0xce, 0x2c, 0x9b, 0x20, 0xf5, 0x87, 0x73, 0xb3, 0xd3, 0xc1, 0x76, 0xa7,
	0x83, 0xbb, 0x9d, 0x0e, 0x2e, 0xf7, 0xba, 0xb2, 0xdd, 0xeb, 0xca, 0xed, 0x5e, 0x57, 0xfe, 0xf7,
	0xcf, 0x64, 0x3e, 0xdb, 0x84, 0x86, 0x48, 0x97, 0xe6, 0x52, 0x2e, 0xe6, 0xc5, 0x79, 0x50, 0x7c,
	0x5f, 0x04, 0x61, 0xf6, 0x44, 0xe6, 0x45, 0x75, 0xb6, 0xc3, 0xcf, 0x65, 0x61, 0xb3, 0xbc, 0xd8,
	0xe0, 0x3e, 0x00, 0x00, 0xff, 0xff, 0x53, 0xbb, 0xc0, 0x8f, 0xd6, 0x01, 0x00, 0x00,
}

func (m *ICAAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ICAAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ICAAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintIcaAccount(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.ConnectionId) > 0 {
		i -= len(m.ConnectionId)
		copy(dAtA[i:], m.ConnectionId)
		i = encodeVarintIcaAccount(dAtA, i, uint64(len(m.ConnectionId)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Type != 0 {
		i = encodeVarintIcaAccount(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x10
	}
	if len(m.ChainId) > 0 {
		i -= len(m.ChainId)
		copy(dAtA[i:], m.ChainId)
		i = encodeVarintIcaAccount(dAtA, i, uint64(len(m.ChainId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintIcaAccount(dAtA []byte, offset int, v uint64) int {
	offset -= sovIcaAccount(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ICAAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChainId)
	if l > 0 {
		n += 1 + l + sovIcaAccount(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovIcaAccount(uint64(m.Type))
	}
	l = len(m.ConnectionId)
	if l > 0 {
		n += 1 + l + sovIcaAccount(uint64(l))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovIcaAccount(uint64(l))
	}
	return n
}

func sovIcaAccount(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozIcaAccount(x uint64) (n int) {
	return sovIcaAccount(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ICAAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIcaAccount
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
			return fmt.Errorf("proto: ICAAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ICAAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIcaAccount
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
				return ErrInvalidLengthIcaAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIcaAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIcaAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= ICAAccountType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIcaAccount
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
				return ErrInvalidLengthIcaAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIcaAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIcaAccount
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
				return ErrInvalidLengthIcaAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIcaAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIcaAccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIcaAccount
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
func skipIcaAccount(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIcaAccount
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
					return 0, ErrIntOverflowIcaAccount
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
					return 0, ErrIntOverflowIcaAccount
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
				return 0, ErrInvalidLengthIcaAccount
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupIcaAccount
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthIcaAccount
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthIcaAccount        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIcaAccount          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupIcaAccount = fmt.Errorf("proto: unexpected end of group")
)