// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: nexus/v1beta1/query.proto

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

type QueryChainMaintainersResponse struct {
	Maintainers []github_com_cosmos_cosmos_sdk_types.ValAddress `protobuf:"bytes,1,rep,name=maintainers,proto3,casttype=github.com/cosmos/cosmos-sdk/types.ValAddress" json:"maintainers,omitempty"`
}

func (m *QueryChainMaintainersResponse) Reset()         { *m = QueryChainMaintainersResponse{} }
func (m *QueryChainMaintainersResponse) String() string { return proto.CompactTextString(m) }
func (*QueryChainMaintainersResponse) ProtoMessage()    {}
func (*QueryChainMaintainersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_18ecb24985e280bf, []int{0}
}
func (m *QueryChainMaintainersResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryChainMaintainersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryChainMaintainersResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryChainMaintainersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryChainMaintainersResponse.Merge(m, src)
}
func (m *QueryChainMaintainersResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryChainMaintainersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryChainMaintainersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryChainMaintainersResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryChainMaintainersResponse)(nil), "nexus.v1beta1.QueryChainMaintainersResponse")
}

func init() { proto.RegisterFile("nexus/v1beta1/query.proto", fileDescriptor_18ecb24985e280bf) }

var fileDescriptor_18ecb24985e280bf = []byte{
	// 232 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcc, 0x4b, 0xad, 0x28,
	0x2d, 0xd6, 0x2f, 0x33, 0x4c, 0x4a, 0x2d, 0x49, 0x34, 0xd4, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x05, 0x4b, 0xe9, 0x41, 0xa5, 0xa4, 0x44, 0xd2, 0xf3,
	0xd3, 0xf3, 0xc1, 0x32, 0xfa, 0x20, 0x16, 0x44, 0x91, 0x52, 0x09, 0x97, 0x6c, 0x20, 0x48, 0x8f,
	0x73, 0x46, 0x62, 0x66, 0x9e, 0x6f, 0x62, 0x66, 0x5e, 0x49, 0x62, 0x66, 0x5e, 0x6a, 0x51, 0x71,
	0x50, 0x6a, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0x50, 0x30, 0x17, 0x77, 0x2e, 0x42, 0x58, 0x82,
	0x51, 0x81, 0x59, 0x83, 0xc7, 0xc9, 0xf0, 0xd7, 0x3d, 0x79, 0xdd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2,
	0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xe4, 0xfc, 0xe2, 0xdc, 0xfc, 0x62, 0x28, 0xa5, 0x5b, 0x9c,
	0x92, 0xad, 0x5f, 0x52, 0x59, 0x90, 0x5a, 0xac, 0x17, 0x96, 0x98, 0xe3, 0x98, 0x92, 0x52, 0x94,
	0x5a, 0x5c, 0x1c, 0x84, 0x6c, 0x8a, 0x53, 0xc0, 0x89, 0x87, 0x72, 0x0c, 0x27, 0x1e, 0xc9, 0x31,
	0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb,
	0x31, 0xdc, 0x78, 0x2c, 0xc7, 0x10, 0x65, 0x84, 0x64, 0x72, 0x62, 0x45, 0x6a, 0x4e, 0x62, 0x51,
	0x5e, 0x6a, 0x49, 0x79, 0x7e, 0x51, 0x36, 0x94, 0xa7, 0x9b, 0x9c, 0x5f, 0x94, 0xaa, 0x5f, 0xa1,
	0x0f, 0xf1, 0x3a, 0xd8, 0xa6, 0x24, 0x36, 0xb0, 0x77, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x20, 0x25, 0xb6, 0x5c, 0x10, 0x01, 0x00, 0x00,
}

func (m *QueryChainMaintainersResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryChainMaintainersResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryChainMaintainersResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Maintainers) > 0 {
		for iNdEx := len(m.Maintainers) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Maintainers[iNdEx])
			copy(dAtA[i:], m.Maintainers[iNdEx])
			i = encodeVarintQuery(dAtA, i, uint64(len(m.Maintainers[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryChainMaintainersResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Maintainers) > 0 {
		for _, b := range m.Maintainers {
			l = len(b)
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryChainMaintainersResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryChainMaintainersResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryChainMaintainersResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Maintainers", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Maintainers = append(m.Maintainers, make([]byte, postIndex-iNdEx))
			copy(m.Maintainers[len(m.Maintainers)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
