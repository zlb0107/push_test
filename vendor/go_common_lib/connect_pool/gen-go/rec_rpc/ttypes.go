// Autogenerated by Thrift Compiler (0.9.3)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package rec_rpc

import (
	"bytes"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"go_common_lib/connect_pool/gen-go/response"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

var _ = response.GoUnusedProtection__
var GoUnusedProtection__ int

// Attributes:
//  - ID
//  - RecTag
//  - Token
type RecPublisher struct {
	ID     int64  `thrift:"id,1,required" json:"id"`
	RecTag string `thrift:"rec_tag,2,required" json:"rec_tag"`
	Token  string `thrift:"token,3" json:"token,omitempty"`
}

func NewRecPublisher() *RecPublisher {
	return &RecPublisher{}
}

func (p *RecPublisher) GetID() int64 {
	return p.ID
}

func (p *RecPublisher) GetRecTag() string {
	return p.RecTag
}

var RecPublisher_Token_DEFAULT string = ""

func (p *RecPublisher) GetToken() string {
	return p.Token
}
func (p *RecPublisher) IsSetToken() bool {
	return p.Token != RecPublisher_Token_DEFAULT
}

func (p *RecPublisher) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetID bool = false
	var issetRecTag bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
			issetID = true
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
			issetRecTag = true
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetID {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field ID is not set"))
	}
	if !issetRecTag {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field RecTag is not set"))
	}
	return nil
}

func (p *RecPublisher) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.ID = v
	}
	return nil
}

func (p *RecPublisher) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.RecTag = v
	}
	return nil
}

func (p *RecPublisher) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 3: ", err)
	} else {
		p.Token = v
	}
	return nil
}

func (p *RecPublisher) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("RecPublisher"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *RecPublisher) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("id", thrift.I64, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:id: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.ID)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.id (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:id: ", p), err)
	}
	return err
}

func (p *RecPublisher) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("rec_tag", thrift.STRING, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:rec_tag: ", p), err)
	}
	if err := oprot.WriteString(string(p.RecTag)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.rec_tag (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:rec_tag: ", p), err)
	}
	return err
}

func (p *RecPublisher) writeField3(oprot thrift.TProtocol) (err error) {
	if p.IsSetToken() {
		if err := oprot.WriteFieldBegin("token", thrift.STRING, 3); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:token: ", p), err)
		}
		if err := oprot.WriteString(string(p.Token)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T.token (3) field write error: ", p), err)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 3:token: ", p), err)
		}
	}
	return err
}

func (p *RecPublisher) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RecPublisher(%+v)", *p)
}

// Attributes:
//  - LogID
//  - UID
//  - Publisher
//  - Extension
//  - IsNew
//  - PageIdx
type RecTagSearchRequest struct {
	LogID     string            `thrift:"log_id,1,required" json:"log_id"`
	UID       int64             `thrift:"uid,2,required" json:"uid"`
	Publisher []*RecPublisher   `thrift:"publisher,3,required" json:"publisher"`
	Extension map[string]string `thrift:"extension,4,required" json:"extension"`
	IsNew     int32             `thrift:"is_new,5,required" json:"is_new"`
	PageIdx   int32             `thrift:"page_idx,6,required" json:"page_idx"`
}

func NewRecTagSearchRequest() *RecTagSearchRequest {
	return &RecTagSearchRequest{}
}

func (p *RecTagSearchRequest) GetLogID() string {
	return p.LogID
}

func (p *RecTagSearchRequest) GetUID() int64 {
	return p.UID
}

func (p *RecTagSearchRequest) GetPublisher() []*RecPublisher {
	return p.Publisher
}

func (p *RecTagSearchRequest) GetExtension() map[string]string {
	return p.Extension
}

func (p *RecTagSearchRequest) GetIsNew() int32 {
	return p.IsNew
}

func (p *RecTagSearchRequest) GetPageIdx() int32 {
	return p.PageIdx
}
func (p *RecTagSearchRequest) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetLogID bool = false
	var issetUID bool = false
	var issetPublisher bool = false
	var issetExtension bool = false
	var issetIsNew bool = false
	var issetPageIdx bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
			issetLogID = true
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
			issetUID = true
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
			issetPublisher = true
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
			issetExtension = true
		case 5:
			if err := p.readField5(iprot); err != nil {
				return err
			}
			issetIsNew = true
		case 6:
			if err := p.readField6(iprot); err != nil {
				return err
			}
			issetPageIdx = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetLogID {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field LogID is not set"))
	}
	if !issetUID {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field UID is not set"))
	}
	if !issetPublisher {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Publisher is not set"))
	}
	if !issetExtension {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Extension is not set"))
	}
	if !issetIsNew {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field IsNew is not set"))
	}
	if !issetPageIdx {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field PageIdx is not set"))
	}
	return nil
}

func (p *RecTagSearchRequest) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.LogID = v
	}
	return nil
}

func (p *RecTagSearchRequest) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI64(); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.UID = v
	}
	return nil
}

func (p *RecTagSearchRequest) readField3(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*RecPublisher, 0, size)
	p.Publisher = tSlice
	for i := 0; i < size; i++ {
		_elem0 := &RecPublisher{}
		if err := _elem0.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem0), err)
		}
		p.Publisher = append(p.Publisher, _elem0)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *RecTagSearchRequest) readField4(iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Extension = tMap
	for i := 0; i < size; i++ {
		var _key1 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_key1 = v
		}
		var _val2 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_val2 = v
		}
		p.Extension[_key1] = _val2
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *RecTagSearchRequest) readField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 5: ", err)
	} else {
		p.IsNew = v
	}
	return nil
}

func (p *RecTagSearchRequest) readField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 6: ", err)
	} else {
		p.PageIdx = v
	}
	return nil
}

func (p *RecTagSearchRequest) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("RecTagSearchRequest"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := p.writeField5(oprot); err != nil {
		return err
	}
	if err := p.writeField6(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *RecTagSearchRequest) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("log_id", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:log_id: ", p), err)
	}
	if err := oprot.WriteString(string(p.LogID)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.log_id (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:log_id: ", p), err)
	}
	return err
}

func (p *RecTagSearchRequest) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("uid", thrift.I64, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:uid: ", p), err)
	}
	if err := oprot.WriteI64(int64(p.UID)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.uid (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:uid: ", p), err)
	}
	return err
}

func (p *RecTagSearchRequest) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("publisher", thrift.LIST, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:publisher: ", p), err)
	}
	if err := oprot.WriteListBegin(thrift.STRUCT, len(p.Publisher)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Publisher {
		if err := v.Write(oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:publisher: ", p), err)
	}
	return err
}

func (p *RecTagSearchRequest) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("extension", thrift.MAP, 4); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 4:extension: ", p), err)
	}
	if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Extension)); err != nil {
		return thrift.PrependError("error writing map begin: ", err)
	}
	for k, v := range p.Extension {
		if err := oprot.WriteString(string(k)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
		if err := oprot.WriteString(string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteMapEnd(); err != nil {
		return thrift.PrependError("error writing map end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 4:extension: ", p), err)
	}
	return err
}

func (p *RecTagSearchRequest) writeField5(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("is_new", thrift.I32, 5); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 5:is_new: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.IsNew)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.is_new (5) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 5:is_new: ", p), err)
	}
	return err
}

func (p *RecTagSearchRequest) writeField6(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("page_idx", thrift.I32, 6); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 6:page_idx: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.PageIdx)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.page_idx (6) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 6:page_idx: ", p), err)
	}
	return err
}

func (p *RecTagSearchRequest) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RecTagSearchRequest(%+v)", *p)
}

// Attributes:
//  - Code
//  - GroupRes
//  - Extension
type RecTagSearchResult_ struct {
	Code      int32                `thrift:"code,1,required" json:"code"`
	GroupRes  *response.MatchGroup `thrift:"group_res,2,required" json:"group_res"`
	Extension map[string]string    `thrift:"extension,3,required" json:"extension"`
}

func NewRecTagSearchResult_() *RecTagSearchResult_ {
	return &RecTagSearchResult_{}
}

func (p *RecTagSearchResult_) GetCode() int32 {
	return p.Code
}

var RecTagSearchResult__GroupRes_DEFAULT *response.MatchGroup

func (p *RecTagSearchResult_) GetGroupRes() *response.MatchGroup {
	if !p.IsSetGroupRes() {
		return RecTagSearchResult__GroupRes_DEFAULT
	}
	return p.GroupRes
}

func (p *RecTagSearchResult_) GetExtension() map[string]string {
	return p.Extension
}
func (p *RecTagSearchResult_) IsSetGroupRes() bool {
	return p.GroupRes != nil
}

func (p *RecTagSearchResult_) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetCode bool = false
	var issetGroupRes bool = false
	var issetExtension bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
			issetCode = true
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
			issetGroupRes = true
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
			issetExtension = true
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetCode {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Code is not set"))
	}
	if !issetGroupRes {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field GroupRes is not set"))
	}
	if !issetExtension {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field Extension is not set"))
	}
	return nil
}

func (p *RecTagSearchResult_) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.Code = v
	}
	return nil
}

func (p *RecTagSearchResult_) readField2(iprot thrift.TProtocol) error {
	p.GroupRes = &response.MatchGroup{}
	if err := p.GroupRes.Read(iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.GroupRes), err)
	}
	return nil
}

func (p *RecTagSearchResult_) readField3(iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	tMap := make(map[string]string, size)
	p.Extension = tMap
	for i := 0; i < size; i++ {
		var _key3 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_key3 = v
		}
		var _val4 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			_val4 = v
		}
		p.Extension[_key3] = _val4
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func (p *RecTagSearchResult_) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("RecTagSearchResult"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *RecTagSearchResult_) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("code", thrift.I32, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:code: ", p), err)
	}
	if err := oprot.WriteI32(int32(p.Code)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.code (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:code: ", p), err)
	}
	return err
}

func (p *RecTagSearchResult_) writeField2(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("group_res", thrift.STRUCT, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:group_res: ", p), err)
	}
	if err := p.GroupRes.Write(oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.GroupRes), err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:group_res: ", p), err)
	}
	return err
}

func (p *RecTagSearchResult_) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("extension", thrift.MAP, 3); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 3:extension: ", p), err)
	}
	if err := oprot.WriteMapBegin(thrift.STRING, thrift.STRING, len(p.Extension)); err != nil {
		return thrift.PrependError("error writing map begin: ", err)
	}
	for k, v := range p.Extension {
		if err := oprot.WriteString(string(k)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
		if err := oprot.WriteString(string(v)); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T. (0) field write error: ", p), err)
		}
	}
	if err := oprot.WriteMapEnd(); err != nil {
		return thrift.PrependError("error writing map end: ", err)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 3:extension: ", p), err)
	}
	return err
}

func (p *RecTagSearchResult_) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RecTagSearchResult_(%+v)", *p)
}
