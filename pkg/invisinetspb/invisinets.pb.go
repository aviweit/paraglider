//
//Copyright 2023 The Invisinets Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: invisinets.proto

package __

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PermitList_Direction int32

const (
	PermitList_INBOUND  PermitList_Direction = 0
	PermitList_OUTBOUND PermitList_Direction = 1
)

// Enum value maps for PermitList_Direction.
var (
	PermitList_Direction_name = map[int32]string{
		0: "INBOUND",
		1: "OUTBOUND",
	}
	PermitList_Direction_value = map[string]int32{
		"INBOUND":  0,
		"OUTBOUND": 1,
	}
)

func (x PermitList_Direction) Enum() *PermitList_Direction {
	p := new(PermitList_Direction)
	*p = x
	return p
}

func (x PermitList_Direction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PermitList_Direction) Descriptor() protoreflect.EnumDescriptor {
	return file_invisinets_proto_enumTypes[0].Descriptor()
}

func (PermitList_Direction) Type() protoreflect.EnumType {
	return &file_invisinets_proto_enumTypes[0]
}

func (x PermitList_Direction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PermitList_Direction.Descriptor instead.
func (PermitList_Direction) EnumDescriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{2, 0}
}

type BasicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *BasicResponse) Reset() {
	*x = BasicResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicResponse) ProtoMessage() {}

func (x *BasicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicResponse.ProtoReflect.Descriptor instead.
func (*BasicResponse) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{0}
}

func (x *BasicResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *BasicResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{1}
}

func (x *Resource) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type PermitList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string                           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Location   string                           `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	Id         string                           `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Properties *PermitList_PermitListProperties `protobuf:"bytes,4,opt,name=properties,proto3" json:"properties,omitempty"`
}

func (x *PermitList) Reset() {
	*x = PermitList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PermitList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermitList) ProtoMessage() {}

func (x *PermitList) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermitList.ProtoReflect.Descriptor instead.
func (*PermitList) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{2}
}

func (x *PermitList) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PermitList) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *PermitList) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PermitList) GetProperties() *PermitList_PermitListProperties {
	if x != nil {
		return x.Properties
	}
	return nil
}

type PermitList_PermitListRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag       []string             `protobuf:"bytes,1,rep,name=tag,proto3" json:"tag,omitempty"`
	Direction PermitList_Direction `protobuf:"varint,2,opt,name=direction,proto3,enum=invisinetspb.PermitList_Direction" json:"direction,omitempty"`
	SrcPort   []string             `protobuf:"bytes,3,rep,name=src_port,json=srcPort,proto3" json:"src_port,omitempty"`
	DstPort   []string             `protobuf:"bytes,4,rep,name=dst_port,json=dstPort,proto3" json:"dst_port,omitempty"`
	Protocol  int32                `protobuf:"varint,5,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *PermitList_PermitListRule) Reset() {
	*x = PermitList_PermitListRule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PermitList_PermitListRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermitList_PermitListRule) ProtoMessage() {}

func (x *PermitList_PermitListRule) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermitList_PermitListRule.ProtoReflect.Descriptor instead.
func (*PermitList_PermitListRule) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{2, 0}
}

func (x *PermitList_PermitListRule) GetTag() []string {
	if x != nil {
		return x.Tag
	}
	return nil
}

func (x *PermitList_PermitListRule) GetDirection() PermitList_Direction {
	if x != nil {
		return x.Direction
	}
	return PermitList_INBOUND
}

func (x *PermitList_PermitListRule) GetSrcPort() []string {
	if x != nil {
		return x.SrcPort
	}
	return nil
}

func (x *PermitList_PermitListRule) GetDstPort() []string {
	if x != nil {
		return x.DstPort
	}
	return nil
}

func (x *PermitList_PermitListRule) GetProtocol() int32 {
	if x != nil {
		return x.Protocol
	}
	return 0
}

type PermitList_PermitListProperties struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AssociatedResource string                       `protobuf:"bytes,1,opt,name=associated_resource,json=associatedResource,proto3" json:"associated_resource,omitempty"`
	Rules              []*PermitList_PermitListRule `protobuf:"bytes,2,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *PermitList_PermitListProperties) Reset() {
	*x = PermitList_PermitListProperties{}
	if protoimpl.UnsafeEnabled {
		mi := &file_invisinets_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PermitList_PermitListProperties) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PermitList_PermitListProperties) ProtoMessage() {}

func (x *PermitList_PermitListProperties) ProtoReflect() protoreflect.Message {
	mi := &file_invisinets_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PermitList_PermitListProperties.ProtoReflect.Descriptor instead.
func (*PermitList_PermitListProperties) Descriptor() ([]byte, []int) {
	return file_invisinets_proto_rawDescGZIP(), []int{2, 1}
}

func (x *PermitList_PermitListProperties) GetAssociatedResource() string {
	if x != nil {
		return x.AssociatedResource
	}
	return ""
}

func (x *PermitList_PermitListProperties) GetRules() []*PermitList_PermitListRule {
	if x != nil {
		return x.Rules
	}
	return nil
}

var File_invisinets_proto protoreflect.FileDescriptor

var file_invisinets_proto_rawDesc = []byte{
	0x0a, 0x10, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0c, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62,
	0x22, 0x43, 0x0a, 0x0d, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x85, 0x04, 0x0a, 0x0a, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x4d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74,
	0x73, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x2e, 0x50,
	0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74,
	0x69, 0x65, 0x73, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x1a,
	0xb6, 0x01, 0x0a, 0x0e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75,
	0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x03, 0x74, 0x61, 0x67, 0x12, 0x40, 0x0a, 0x09, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69,
	0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73,
	0x74, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x64, 0x69, 0x72,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x72, 0x63, 0x5f, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x72, 0x63, 0x50, 0x6f, 0x72,
	0x74, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x73, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x07, 0x64, 0x73, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x1a, 0x86, 0x01, 0x0a, 0x14, 0x50, 0x65, 0x72,
	0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65,
	0x73, 0x12, 0x2f, 0x0a, 0x13, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12,
	0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x12, 0x3d, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x27, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62,
	0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x2e, 0x50, 0x65, 0x72, 0x6d,
	0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65,
	0x73, 0x22, 0x26, 0x0a, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0b,
	0x0a, 0x07, 0x49, 0x4e, 0x42, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x4f,
	0x55, 0x54, 0x42, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x32, 0x9c, 0x01, 0x0a, 0x0b, 0x43, 0x6c,
	0x6f, 0x75, 0x64, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x48, 0x0a, 0x0d, 0x53, 0x65, 0x74,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x18, 0x2e, 0x69, 0x6e, 0x76,
	0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74,
	0x73, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x16, 0x2e, 0x69, 0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74,
	0x73, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x18, 0x2e, 0x69,
	0x6e, 0x76, 0x69, 0x73, 0x69, 0x6e, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x50, 0x65, 0x72, 0x6d,
	0x69, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_invisinets_proto_rawDescOnce sync.Once
	file_invisinets_proto_rawDescData = file_invisinets_proto_rawDesc
)

func file_invisinets_proto_rawDescGZIP() []byte {
	file_invisinets_proto_rawDescOnce.Do(func() {
		file_invisinets_proto_rawDescData = protoimpl.X.CompressGZIP(file_invisinets_proto_rawDescData)
	})
	return file_invisinets_proto_rawDescData
}

var file_invisinets_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_invisinets_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_invisinets_proto_goTypes = []interface{}{
	(PermitList_Direction)(0),               // 0: invisinetspb.PermitList.Direction
	(*BasicResponse)(nil),                   // 1: invisinetspb.BasicResponse
	(*Resource)(nil),                        // 2: invisinetspb.Resource
	(*PermitList)(nil),                      // 3: invisinetspb.PermitList
	(*PermitList_PermitListRule)(nil),       // 4: invisinetspb.PermitList.PermitListRule
	(*PermitList_PermitListProperties)(nil), // 5: invisinetspb.PermitList.PermitListProperties
}
var file_invisinets_proto_depIdxs = []int32{
	5, // 0: invisinetspb.PermitList.properties:type_name -> invisinetspb.PermitList.PermitListProperties
	0, // 1: invisinetspb.PermitList.PermitListRule.direction:type_name -> invisinetspb.PermitList.Direction
	4, // 2: invisinetspb.PermitList.PermitListProperties.rules:type_name -> invisinetspb.PermitList.PermitListRule
	3, // 3: invisinetspb.CloudPlugin.SetPermitList:input_type -> invisinetspb.PermitList
	2, // 4: invisinetspb.CloudPlugin.GetPermitList:input_type -> invisinetspb.Resource
	1, // 5: invisinetspb.CloudPlugin.SetPermitList:output_type -> invisinetspb.BasicResponse
	3, // 6: invisinetspb.CloudPlugin.GetPermitList:output_type -> invisinetspb.PermitList
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_invisinets_proto_init() }
func file_invisinets_proto_init() {
	if File_invisinets_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_invisinets_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BasicResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invisinets_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invisinets_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PermitList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invisinets_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PermitList_PermitListRule); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_invisinets_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PermitList_PermitListProperties); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_invisinets_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_invisinets_proto_goTypes,
		DependencyIndexes: file_invisinets_proto_depIdxs,
		EnumInfos:         file_invisinets_proto_enumTypes,
		MessageInfos:      file_invisinets_proto_msgTypes,
	}.Build()
	File_invisinets_proto = out.File
	file_invisinets_proto_rawDesc = nil
	file_invisinets_proto_goTypes = nil
	file_invisinets_proto_depIdxs = nil
}
