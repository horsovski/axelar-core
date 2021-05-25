// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ethereum/v1beta1/types.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// BurnerInfo describes information required to burn token at an burner address
// that is deposited by an user
type BurnerInfo struct {
	TokenAddress Address `protobuf:"bytes,1,opt,name=token_address,json=tokenAddress,proto3,customtype=Address" json:"token_address"`
	Symbol       string  `protobuf:"bytes,2,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Salt         Hash    `protobuf:"bytes,3,opt,name=salt,proto3,customtype=Hash" json:"salt"`
}

func (m *BurnerInfo) Reset()         { *m = BurnerInfo{} }
func (m *BurnerInfo) String() string { return proto.CompactTextString(m) }
func (*BurnerInfo) ProtoMessage()    {}
func (*BurnerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c65a343c55145f9, []int{0}
}
func (m *BurnerInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BurnerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BurnerInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BurnerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BurnerInfo.Merge(m, src)
}
func (m *BurnerInfo) XXX_Size() int {
	return m.Size()
}
func (m *BurnerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_BurnerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_BurnerInfo proto.InternalMessageInfo

// ERC20Deposit contains information for an ERC20 deposit
type ERC20Deposit struct {
	TxID          Hash                                    `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3,customtype=Hash" json:"tx_id"`
	Amount        github_com_cosmos_cosmos_sdk_types.Uint `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Uint" json:"amount"`
	Symbol        string                                  `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	BurnerAddress Address                                 `protobuf:"bytes,4,opt,name=burner_address,json=burnerAddress,proto3,customtype=Address" json:"burner_address"`
}

func (m *ERC20Deposit) Reset()         { *m = ERC20Deposit{} }
func (m *ERC20Deposit) String() string { return proto.CompactTextString(m) }
func (*ERC20Deposit) ProtoMessage()    {}
func (*ERC20Deposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c65a343c55145f9, []int{1}
}
func (m *ERC20Deposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ERC20Deposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ERC20Deposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ERC20Deposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ERC20Deposit.Merge(m, src)
}
func (m *ERC20Deposit) XXX_Size() int {
	return m.Size()
}
func (m *ERC20Deposit) XXX_DiscardUnknown() {
	xxx_messageInfo_ERC20Deposit.DiscardUnknown(m)
}

var xxx_messageInfo_ERC20Deposit proto.InternalMessageInfo

// ERC20TokenDeployment describes information about an ERC20 token
type ERC20TokenDeployment struct {
	Symbol       string  `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	TokenAddress Address `protobuf:"bytes,2,opt,name=token_address,json=tokenAddress,proto3,customtype=Address" json:"token_address"`
}

func (m *ERC20TokenDeployment) Reset()         { *m = ERC20TokenDeployment{} }
func (m *ERC20TokenDeployment) String() string { return proto.CompactTextString(m) }
func (*ERC20TokenDeployment) ProtoMessage()    {}
func (*ERC20TokenDeployment) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c65a343c55145f9, []int{2}
}
func (m *ERC20TokenDeployment) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ERC20TokenDeployment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ERC20TokenDeployment.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ERC20TokenDeployment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ERC20TokenDeployment.Merge(m, src)
}
func (m *ERC20TokenDeployment) XXX_Size() int {
	return m.Size()
}
func (m *ERC20TokenDeployment) XXX_DiscardUnknown() {
	xxx_messageInfo_ERC20TokenDeployment.DiscardUnknown(m)
}

var xxx_messageInfo_ERC20TokenDeployment proto.InternalMessageInfo

func init() {
	proto.RegisterType((*BurnerInfo)(nil), "ethereum.v1beta1.BurnerInfo")
	proto.RegisterType((*ERC20Deposit)(nil), "ethereum.v1beta1.ERC20Deposit")
	proto.RegisterType((*ERC20TokenDeployment)(nil), "ethereum.v1beta1.ERC20TokenDeployment")
}

func init() { proto.RegisterFile("ethereum/v1beta1/types.proto", fileDescriptor_9c65a343c55145f9) }

var fileDescriptor_9c65a343c55145f9 = []byte{
	// 379 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcf, 0x4e, 0xf2, 0x40,
	0x14, 0xc5, 0x3b, 0x1f, 0xfd, 0x30, 0x4e, 0x8a, 0x9a, 0x86, 0x18, 0x62, 0x4c, 0x21, 0x6c, 0xc4,
	0x05, 0x54, 0xfc, 0xb7, 0xb7, 0x62, 0x94, 0x6d, 0xc5, 0x8d, 0x1b, 0xd2, 0xd2, 0x11, 0x1a, 0xda,
	0xde, 0x66, 0x66, 0xaa, 0x25, 0xf1, 0x21, 0x7c, 0x2c, 0x96, 0xc4, 0x95, 0x71, 0x41, 0xb4, 0xbc,
	0x88, 0x61, 0x18, 0x09, 0x44, 0x13, 0x57, 0xed, 0xe9, 0x39, 0x33, 0xe7, 0xf6, 0x97, 0x8b, 0xf7,
	0x09, 0x1f, 0x10, 0x4a, 0x92, 0xd0, 0x7c, 0x6c, 0xba, 0x84, 0x3b, 0x4d, 0x93, 0x8f, 0x62, 0xc2,
	0x1a, 0x31, 0x05, 0x0e, 0xfa, 0xce, 0xb7, 0xdb, 0x90, 0xee, 0x5e, 0xb1, 0x0f, 0x7d, 0x10, 0xa6,
	0x39, 0x7f, 0x5b, 0xe4, 0xaa, 0xcf, 0x18, 0x5b, 0x09, 0x8d, 0x08, 0x6d, 0x47, 0x0f, 0xa0, 0x9f,
	0xe2, 0x02, 0x87, 0x21, 0x89, 0xba, 0x8e, 0xe7, 0x51, 0xc2, 0x58, 0x09, 0x55, 0x50, 0x4d, 0xb3,
	0xb6, 0xc7, 0xd3, 0xb2, 0xf2, 0x3e, 0x2d, 0x6f, 0x5c, 0x2c, 0x3e, 0xdb, 0x9a, 0x48, 0x49, 0xa5,
	0xef, 0xe2, 0x3c, 0x1b, 0x85, 0x2e, 0x04, 0xa5, 0x7f, 0x15, 0x54, 0xdb, 0xb4, 0xa5, 0xd2, 0x2b,
	0x58, 0x65, 0x4e, 0xc0, 0x4b, 0x39, 0x71, 0x89, 0x26, 0x2f, 0x51, 0x6f, 0x1c, 0x36, 0xb0, 0x85,
	0x53, 0x7d, 0x45, 0x58, 0xbb, 0xb2, 0x2f, 0x8f, 0x8f, 0x5a, 0x24, 0x06, 0xe6, 0x73, 0xfd, 0x10,
	0xff, 0xe7, 0x69, 0xd7, 0xf7, 0x64, 0x71, 0x71, 0xf5, 0x4c, 0x36, 0x2d, 0xab, 0x9d, 0xb4, 0xdd,
	0xb2, 0x55, 0x9e, 0xb6, 0x3d, 0xfd, 0x1a, 0xe7, 0x9d, 0x10, 0x92, 0x88, 0x8b, 0x56, 0xcd, 0x32,
	0x65, 0xf6, 0xa0, 0xef, 0xf3, 0x41, 0xe2, 0x36, 0x7a, 0x10, 0x9a, 0x3d, 0x60, 0x21, 0x30, 0xf9,
	0xa8, 0x33, 0x6f, 0x28, 0x19, 0xdd, 0xf9, 0x11, 0xb7, 0xe5, 0xf1, 0x95, 0xf1, 0x73, 0x6b, 0xe3,
	0x9f, 0xe3, 0x2d, 0x57, 0xa0, 0x59, 0xd2, 0x50, 0x7f, 0xa7, 0x51, 0x58, 0xc4, 0xa4, 0xac, 0x7a,
	0xb8, 0x28, 0xfe, 0xa9, 0x33, 0x67, 0xd4, 0x22, 0x71, 0x00, 0xa3, 0x90, 0xac, 0xf5, 0xa0, 0xb5,
	0x9e, 0x1f, 0xd0, 0x05, 0xc5, 0x3f, 0xa0, 0x5b, 0xb7, 0xe3, 0x4f, 0x43, 0x19, 0x67, 0x06, 0x9a,
	0x64, 0x06, 0xfa, 0xc8, 0x0c, 0xf4, 0x32, 0x33, 0x94, 0xc9, 0xcc, 0x50, 0xde, 0x66, 0x86, 0x72,
	0x7f, 0xb6, 0x02, 0xc1, 0x49, 0x49, 0xe0, 0xd0, 0x88, 0xf0, 0x27, 0xa0, 0x43, 0xa9, 0xea, 0x3d,
	0xa0, 0xc4, 0x4c, 0xcd, 0xe5, 0x0e, 0x09, 0x2e, 0x6e, 0x5e, 0x2c, 0xc5, 0xc9, 0x57, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x48, 0x1c, 0x4d, 0xa2, 0x5c, 0x02, 0x00, 0x00,
}

func (m *BurnerInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BurnerInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BurnerInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Salt.Size()
		i -= size
		if _, err := m.Salt.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.TokenAddress.Size()
		i -= size
		if _, err := m.TokenAddress.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *ERC20Deposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ERC20Deposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ERC20Deposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.BurnerAddress.Size()
		i -= size
		if _, err := m.BurnerAddress.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0x1a
	}
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.TxID.Size()
		i -= size
		if _, err := m.TxID.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *ERC20TokenDeployment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ERC20TokenDeployment) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ERC20TokenDeployment) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.TokenAddress.Size()
		i -= size
		if _, err := m.TokenAddress.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BurnerInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.TokenAddress.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.Salt.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *ERC20Deposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.TxID.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Amount.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.BurnerAddress.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *ERC20TokenDeployment) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.TokenAddress.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BurnerInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: BurnerInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BurnerInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenAddress.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Salt", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Salt.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *ERC20Deposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: ERC20Deposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ERC20Deposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TxID.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BurnerAddress", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BurnerAddress.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *ERC20TokenDeployment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: ERC20TokenDeployment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ERC20TokenDeployment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenAddress.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
