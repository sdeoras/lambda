// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// EmailRequest is the request to send email. This data is streamed from client
// to the server.
type EmailRequest struct {
	ToEmail              string   `protobuf:"bytes,1,opt,name=to_email,json=toEmail,proto3" json:"to_email,omitempty"`
	ToName               string   `protobuf:"bytes,2,opt,name=to_name,json=toName,proto3" json:"to_name,omitempty"`
	FromEmail            string   `protobuf:"bytes,3,opt,name=from_email,json=fromEmail,proto3" json:"from_email,omitempty"`
	FromName             string   `protobuf:"bytes,4,opt,name=from_name,json=fromName,proto3" json:"from_name,omitempty"`
	Subject              string   `protobuf:"bytes,5,opt,name=subject,proto3" json:"subject,omitempty"`
	Body                 []byte   `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmailRequest) Reset()         { *m = EmailRequest{} }
func (m *EmailRequest) String() string { return proto.CompactTextString(m) }
func (*EmailRequest) ProtoMessage()    {}
func (*EmailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_6548425e868eea19, []int{0}
}
func (m *EmailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmailRequest.Unmarshal(m, b)
}
func (m *EmailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmailRequest.Marshal(b, m, deterministic)
}
func (dst *EmailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmailRequest.Merge(dst, src)
}
func (m *EmailRequest) XXX_Size() int {
	return xxx_messageInfo_EmailRequest.Size(m)
}
func (m *EmailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EmailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EmailRequest proto.InternalMessageInfo

func (m *EmailRequest) GetToEmail() string {
	if m != nil {
		return m.ToEmail
	}
	return ""
}

func (m *EmailRequest) GetToName() string {
	if m != nil {
		return m.ToName
	}
	return ""
}

func (m *EmailRequest) GetFromEmail() string {
	if m != nil {
		return m.FromEmail
	}
	return ""
}

func (m *EmailRequest) GetFromName() string {
	if m != nil {
		return m.FromName
	}
	return ""
}

func (m *EmailRequest) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *EmailRequest) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

// SendResponse is the response from the server after email is sent.
type EmailResponse struct {
	StatusCode           int64                    `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Body                 string                   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Headers              map[string]*ListOfString `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *EmailResponse) Reset()         { *m = EmailResponse{} }
func (m *EmailResponse) String() string { return proto.CompactTextString(m) }
func (*EmailResponse) ProtoMessage()    {}
func (*EmailResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_6548425e868eea19, []int{1}
}
func (m *EmailResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmailResponse.Unmarshal(m, b)
}
func (m *EmailResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmailResponse.Marshal(b, m, deterministic)
}
func (dst *EmailResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmailResponse.Merge(dst, src)
}
func (m *EmailResponse) XXX_Size() int {
	return xxx_messageInfo_EmailResponse.Size(m)
}
func (m *EmailResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EmailResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EmailResponse proto.InternalMessageInfo

func (m *EmailResponse) GetStatusCode() int64 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *EmailResponse) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *EmailResponse) GetHeaders() map[string]*ListOfString {
	if m != nil {
		return m.Headers
	}
	return nil
}

// ListOfString is a list of strings
type ListOfString struct {
	Value                []string `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListOfString) Reset()         { *m = ListOfString{} }
func (m *ListOfString) String() string { return proto.CompactTextString(m) }
func (*ListOfString) ProtoMessage()    {}
func (*ListOfString) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_6548425e868eea19, []int{2}
}
func (m *ListOfString) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListOfString.Unmarshal(m, b)
}
func (m *ListOfString) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListOfString.Marshal(b, m, deterministic)
}
func (dst *ListOfString) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListOfString.Merge(dst, src)
}
func (m *ListOfString) XXX_Size() int {
	return xxx_messageInfo_ListOfString.Size(m)
}
func (m *ListOfString) XXX_DiscardUnknown() {
	xxx_messageInfo_ListOfString.DiscardUnknown(m)
}

var xxx_messageInfo_ListOfString proto.InternalMessageInfo

func (m *ListOfString) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

// InferImageRequest takes a list of images (as bytes) and model and label file path URIs.
type InferImageRequest struct {
	// List of images
	Images []*Image `protobuf:"bytes,1,rep,name=images,proto3" json:"images,omitempty"`
	// Path to model pb file
	ModelPath string `protobuf:"bytes,2,opt,name=model_path,json=modelPath,proto3" json:"model_path,omitempty"`
	// Path to label file
	LabelPath            string   `protobuf:"bytes,3,opt,name=label_path,json=labelPath,proto3" json:"label_path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InferImageRequest) Reset()         { *m = InferImageRequest{} }
func (m *InferImageRequest) String() string { return proto.CompactTextString(m) }
func (*InferImageRequest) ProtoMessage()    {}
func (*InferImageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_6548425e868eea19, []int{3}
}
func (m *InferImageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InferImageRequest.Unmarshal(m, b)
}
func (m *InferImageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InferImageRequest.Marshal(b, m, deterministic)
}
func (dst *InferImageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InferImageRequest.Merge(dst, src)
}
func (m *InferImageRequest) XXX_Size() int {
	return xxx_messageInfo_InferImageRequest.Size(m)
}
func (m *InferImageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InferImageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InferImageRequest proto.InternalMessageInfo

func (m *InferImageRequest) GetImages() []*Image {
	if m != nil {
		return m.Images
	}
	return nil
}

func (m *InferImageRequest) GetModelPath() string {
	if m != nil {
		return m.ModelPath
	}
	return ""
}

func (m *InferImageRequest) GetLabelPath() string {
	if m != nil {
		return m.LabelPath
	}
	return ""
}

// Image consists of a name and bytes
type Image struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Image) Reset()         { *m = Image{} }
func (m *Image) String() string { return proto.CompactTextString(m) }
func (*Image) ProtoMessage()    {}
func (*Image) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_6548425e868eea19, []int{4}
}
func (m *Image) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Image.Unmarshal(m, b)
}
func (m *Image) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Image.Marshal(b, m, deterministic)
}
func (dst *Image) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Image.Merge(dst, src)
}
func (m *Image) XXX_Size() int {
	return xxx_messageInfo_Image.Size(m)
}
func (m *Image) XXX_DiscardUnknown() {
	xxx_messageInfo_Image.DiscardUnknown(m)
}

var xxx_messageInfo_Image proto.InternalMessageInfo

func (m *Image) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Image) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// InferImageResponse is a list of output, one per input image.
type InferImageResponse struct {
	Outputs              []*InferOutput `protobuf:"bytes,1,rep,name=outputs,proto3" json:"outputs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *InferImageResponse) Reset()         { *m = InferImageResponse{} }
func (m *InferImageResponse) String() string { return proto.CompactTextString(m) }
func (*InferImageResponse) ProtoMessage()    {}
func (*InferImageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_6548425e868eea19, []int{5}
}
func (m *InferImageResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InferImageResponse.Unmarshal(m, b)
}
func (m *InferImageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InferImageResponse.Marshal(b, m, deterministic)
}
func (dst *InferImageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InferImageResponse.Merge(dst, src)
}
func (m *InferImageResponse) XXX_Size() int {
	return xxx_messageInfo_InferImageResponse.Size(m)
}
func (m *InferImageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InferImageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InferImageResponse proto.InternalMessageInfo

func (m *InferImageResponse) GetOutputs() []*InferOutput {
	if m != nil {
		return m.Outputs
	}
	return nil
}

// InferOutput consists of a name, inferred label and probability for that label.
type InferOutput struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Label                string   `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	Probability          int64    `protobuf:"varint,3,opt,name=probability,proto3" json:"probability,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InferOutput) Reset()         { *m = InferOutput{} }
func (m *InferOutput) String() string { return proto.CompactTextString(m) }
func (*InferOutput) ProtoMessage()    {}
func (*InferOutput) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_6548425e868eea19, []int{6}
}
func (m *InferOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InferOutput.Unmarshal(m, b)
}
func (m *InferOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InferOutput.Marshal(b, m, deterministic)
}
func (dst *InferOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InferOutput.Merge(dst, src)
}
func (m *InferOutput) XXX_Size() int {
	return xxx_messageInfo_InferOutput.Size(m)
}
func (m *InferOutput) XXX_DiscardUnknown() {
	xxx_messageInfo_InferOutput.DiscardUnknown(m)
}

var xxx_messageInfo_InferOutput proto.InternalMessageInfo

func (m *InferOutput) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *InferOutput) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *InferOutput) GetProbability() int64 {
	if m != nil {
		return m.Probability
	}
	return 0
}

func init() {
	proto.RegisterType((*EmailRequest)(nil), "api.EmailRequest")
	proto.RegisterType((*EmailResponse)(nil), "api.EmailResponse")
	proto.RegisterMapType((map[string]*ListOfString)(nil), "api.EmailResponse.HeadersEntry")
	proto.RegisterType((*ListOfString)(nil), "api.ListOfString")
	proto.RegisterType((*InferImageRequest)(nil), "api.InferImageRequest")
	proto.RegisterType((*Image)(nil), "api.Image")
	proto.RegisterType((*InferImageResponse)(nil), "api.InferImageResponse")
	proto.RegisterType((*InferOutput)(nil), "api.InferOutput")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ApiClient is the client API for Api service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ApiClient interface {
	// Email service sends email
	Email(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error)
	// InferImage applies trained model on input image for inferring labels
	InferImage(ctx context.Context, opts ...grpc.CallOption) (Api_InferImageClient, error)
}

type apiClient struct {
	cc *grpc.ClientConn
}

func NewApiClient(cc *grpc.ClientConn) ApiClient {
	return &apiClient{cc}
}

func (c *apiClient) Email(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, "/api.Api/Email", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) InferImage(ctx context.Context, opts ...grpc.CallOption) (Api_InferImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Api_serviceDesc.Streams[0], "/api.Api/InferImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &apiInferImageClient{stream}
	return x, nil
}

type Api_InferImageClient interface {
	Send(*InferImageRequest) error
	CloseAndRecv() (*InferImageResponse, error)
	grpc.ClientStream
}

type apiInferImageClient struct {
	grpc.ClientStream
}

func (x *apiInferImageClient) Send(m *InferImageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *apiInferImageClient) CloseAndRecv() (*InferImageResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(InferImageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ApiServer is the server API for Api service.
type ApiServer interface {
	// Email service sends email
	Email(context.Context, *EmailRequest) (*EmailResponse, error)
	// InferImage applies trained model on input image for inferring labels
	InferImage(Api_InferImageServer) error
}

func RegisterApiServer(s *grpc.Server, srv ApiServer) {
	s.RegisterService(&_Api_serviceDesc, srv)
}

func _Api_Email_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Email(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Email",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Email(ctx, req.(*EmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_InferImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ApiServer).InferImage(&apiInferImageServer{stream})
}

type Api_InferImageServer interface {
	SendAndClose(*InferImageResponse) error
	Recv() (*InferImageRequest, error)
	grpc.ServerStream
}

type apiInferImageServer struct {
	grpc.ServerStream
}

func (x *apiInferImageServer) SendAndClose(m *InferImageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *apiInferImageServer) Recv() (*InferImageRequest, error) {
	m := new(InferImageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Api_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Api",
	HandlerType: (*ApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Email",
			Handler:    _Api_Email_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "InferImage",
			Handler:       _Api_InferImage_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "api.proto",
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_api_6548425e868eea19) }

var fileDescriptor_api_6548425e868eea19 = []byte{
	// 482 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x53, 0x51, 0x6f, 0x94, 0x40,
	0x10, 0x96, 0x52, 0x8e, 0x32, 0x60, 0xd2, 0xdb, 0x34, 0x16, 0xcf, 0x98, 0x12, 0x62, 0x22, 0xf1,
	0xe1, 0x34, 0xe7, 0x8b, 0xfa, 0x64, 0x63, 0x9a, 0xd8, 0x44, 0xad, 0xc1, 0x27, 0x9f, 0x2e, 0x4b,
	0xd9, 0xeb, 0xa1, 0xc0, 0x22, 0x3b, 0x98, 0x9c, 0xbf, 0xc9, 0xdf, 0xe2, 0x6f, 0x32, 0x3b, 0xbb,
	0x58, 0xcc, 0xf5, 0x6d, 0xe6, 0xfb, 0xbe, 0x99, 0xfd, 0x66, 0x06, 0x20, 0xe0, 0x5d, 0xb5, 0xec,
	0x7a, 0x89, 0x92, 0xb9, 0xbc, 0xab, 0xd2, 0xdf, 0x0e, 0x44, 0x17, 0x0d, 0xaf, 0xea, 0x5c, 0xfc,
	0x18, 0x84, 0x42, 0xf6, 0x10, 0x8e, 0x50, 0xae, 0x85, 0x86, 0x62, 0x27, 0x71, 0xb2, 0x20, 0xf7,
	0x51, 0x92, 0x82, 0x9d, 0x82, 0x8f, 0x72, 0xdd, 0xf2, 0x46, 0xc4, 0x07, 0xc4, 0xcc, 0x50, 0x7e,
	0xe2, 0x8d, 0x60, 0x8f, 0x01, 0x36, 0xbd, 0x6c, 0x6c, 0x95, 0x4b, 0x5c, 0xa0, 0x11, 0x53, 0xf7,
	0x08, 0x28, 0x31, 0x95, 0x87, 0xc4, 0x1e, 0x69, 0x80, 0x6a, 0x63, 0xf0, 0xd5, 0x50, 0x7c, 0x13,
	0xd7, 0x18, 0x7b, 0xe6, 0x39, 0x9b, 0x32, 0x06, 0x87, 0x85, 0x2c, 0x77, 0xf1, 0x2c, 0x71, 0xb2,
	0x28, 0xa7, 0x38, 0xfd, 0xe3, 0xc0, 0x7d, 0x6b, 0x57, 0x75, 0xb2, 0x55, 0x82, 0x9d, 0x41, 0xa8,
	0x90, 0xe3, 0xa0, 0xd6, 0xd7, 0xb2, 0x14, 0x64, 0xd9, 0xcd, 0xc1, 0x40, 0xef, 0x64, 0x29, 0xfe,
	0xb5, 0x31, 0x96, 0x29, 0x66, 0xaf, 0xc1, 0xdf, 0x0a, 0x5e, 0x8a, 0x5e, 0xc5, 0x6e, 0xe2, 0x66,
	0xe1, 0xea, 0x6c, 0xa9, 0xf7, 0xf2, 0x5f, 0xe7, 0xe5, 0x7b, 0xa3, 0xb8, 0x68, 0xb1, 0xdf, 0xe5,
	0xa3, 0x7e, 0xf1, 0x11, 0xa2, 0x29, 0xc1, 0x8e, 0xc1, 0xfd, 0x2e, 0x76, 0x76, 0x55, 0x3a, 0x64,
	0x4f, 0xc1, 0xfb, 0xc9, 0xeb, 0xc1, 0x2c, 0x29, 0x5c, 0xcd, 0xa9, 0xf5, 0x87, 0x4a, 0xe1, 0xd5,
	0xe6, 0x0b, 0xf6, 0x55, 0x7b, 0x93, 0x1b, 0xfe, 0xcd, 0xc1, 0x2b, 0x27, 0x7d, 0x02, 0xd1, 0x94,
	0x62, 0x27, 0x63, 0xb1, 0x93, 0xb8, 0x59, 0x60, 0x95, 0xe9, 0x00, 0xf3, 0xcb, 0x76, 0x23, 0xfa,
	0xcb, 0x86, 0xdf, 0x88, 0xf1, 0x52, 0x29, 0xcc, 0x2a, 0x9d, 0x2b, 0xd2, 0x86, 0x2b, 0xa0, 0x87,
	0x8c, 0xc4, 0x32, 0xfa, 0x32, 0x8d, 0x2c, 0x45, 0xbd, 0xee, 0x38, 0x6e, 0xed, 0x0a, 0x02, 0x42,
	0x3e, 0x73, 0xdc, 0x6a, 0xba, 0xe6, 0xc5, 0x48, 0xdb, 0xc3, 0x11, 0xa2, 0xe9, 0xf4, 0x39, 0x78,
	0xd4, 0x4e, 0xef, 0x90, 0x8e, 0x67, 0xa6, 0xa4, 0x58, 0x63, 0x25, 0x47, 0x4e, 0x4d, 0xa3, 0x9c,
	0xe2, 0xf4, 0x2d, 0xb0, 0xa9, 0x4f, 0x7b, 0xa2, 0x67, 0xe0, 0xcb, 0x01, 0xbb, 0x01, 0x47, 0xa7,
	0xc7, 0xc6, 0xa9, 0x56, 0x5e, 0x11, 0x91, 0x8f, 0x82, 0xf4, 0x2b, 0x84, 0x13, 0xfc, 0xce, 0x87,
	0x4f, 0xc0, 0x23, 0x8b, 0x76, 0x1c, 0x93, 0xb0, 0x04, 0xc2, 0xae, 0x97, 0x05, 0x2f, 0xaa, 0xba,
	0xc2, 0x1d, 0xcd, 0xe2, 0xe6, 0x53, 0x68, 0xf5, 0x0b, 0xdc, 0xf3, 0xae, 0x62, 0x2f, 0xc0, 0x33,
	0x9f, 0xe5, 0x7c, 0x7a, 0x73, 0x5a, 0xe9, 0x82, 0xed, 0x7f, 0x06, 0xe9, 0x3d, 0x76, 0x0e, 0x70,
	0x3b, 0x15, 0x7b, 0x70, 0x6b, 0x7e, 0x7a, 0x8e, 0xc5, 0xe9, 0x1e, 0x3e, 0x36, 0xc8, 0x9c, 0x62,
	0x46, 0xbf, 0xdc, 0xcb, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x14, 0xac, 0x3e, 0x96, 0x7f, 0x03,
	0x00, 0x00,
}
