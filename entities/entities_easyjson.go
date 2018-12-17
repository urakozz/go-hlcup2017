// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entities

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	"sync/atomic"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities(in *jlexer.Lexer, out *VisitDiff) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "LocationID":
			easyjson3e8ab7adDecode(in, &out.LocationID)
		case "UserID":
			easyjson3e8ab7adDecode(in, &out.UserID)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities(out *jwriter.Writer, in VisitDiff) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"LocationID\":")
	easyjson3e8ab7adEncode(out, in.LocationID)
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"UserID\":")
	easyjson3e8ab7adEncode(out, in.UserID)
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v VisitDiff) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VisitDiff) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *VisitDiff) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VisitDiff) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities(l, v)
}
func easyjson3e8ab7adDecode(in *jlexer.Lexer, out *struct {
	HasDiff bool
	Old     int64
	New     int64
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "HasDiff":
			out.HasDiff = bool(in.Bool())
		case "Old":
			out.Old = int64(in.Int64())
		case "New":
			out.New = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncode(out *jwriter.Writer, in struct {
	HasDiff bool
	Old     int64
	New     int64
}) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"HasDiff\":")
	out.Bool(bool(in.HasDiff))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Old\":")
	out.Int64(int64(in.Old))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"New\":")
	out.Int64(int64(in.New))
	out.RawByte('}')
}
func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities1(in *jlexer.Lexer, out *VisitContainer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "visits":
			if in.IsNull() {
				in.Skip()
				out.Visits = nil
			} else {
				in.Delim('[')
				if out.Visits == nil {
					if !in.IsDelim(']') {
						out.Visits = make([]*Visit, 0, 8)
					} else {
						out.Visits = []*Visit{}
					}
				} else {
					out.Visits = (out.Visits)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *Visit
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(Visit)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Visits = append(out.Visits, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities1(out *jwriter.Writer, in VisitContainer) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"visits\":")
	if in.Visits == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in.Visits {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				(*v3).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v VisitContainer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v VisitContainer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *VisitContainer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *VisitContainer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities1(l, v)
}
func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities2(in *jlexer.Lexer, out *Visit) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "location":
			if in.IsNull() {
				in.Skip()
				out.LocationID = nil
			} else {
				if out.LocationID == nil {
					out.LocationID = new(int64)
				}
				*out.LocationID = int64(in.Int64())
			}
		case "user":
			if in.IsNull() {
				in.Skip()
				out.UserID = nil
			} else {
				if out.UserID == nil {
					out.UserID = new(int64)
				}
				*out.UserID = int64(in.Int64())
			}
		case "visited_at":
			if in.IsNull() {
				in.Skip()
				out.VisitedAt = nil
			} else {
				if out.VisitedAt == nil {
					out.VisitedAt = new(int64)
				}
				*out.VisitedAt = int64(in.Int64())
			}
		case "mark":
			if in.IsNull() {
				in.Skip()
				out.Mark = nil
			} else {
				if out.Mark == nil {
					out.Mark = new(uint8)
				}
				*out.Mark = uint8(in.Uint8())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities2(out *jwriter.Writer, in Visit) {
	if atomic.LoadInt32(&in.hasJSON) == 1 {
		out.Raw(in.json, nil)
		return
	}
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.Int64(int64(in.ID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"location\":")
	if in.LocationID == nil {
		out.RawString("null")
	} else {
		out.Int64(int64(*in.LocationID))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"user\":")
	if in.UserID == nil {
		out.RawString("null")
	} else {
		out.Int64(int64(*in.UserID))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"visited_at\":")
	if in.VisitedAt == nil {
		out.RawString("null")
	} else {
		out.Int64(int64(*in.VisitedAt))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"mark\":")
	if in.Mark == nil {
		out.RawString("null")
	} else {
		out.Uint8(uint8(*in.Mark))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Visit) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Visit) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Visit) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Visit) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities2(l, v)
}
func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities3(in *jlexer.Lexer, out *UserContainer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "users":
			if in.IsNull() {
				in.Skip()
				out.Users = nil
			} else {
				in.Delim('[')
				if out.Users == nil {
					if !in.IsDelim(']') {
						out.Users = make([]*User, 0, 8)
					} else {
						out.Users = []*User{}
					}
				} else {
					out.Users = (out.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v4 *User
					if in.IsNull() {
						in.Skip()
						v4 = nil
					} else {
						if v4 == nil {
							v4 = new(User)
						}
						(*v4).UnmarshalEasyJSON(in)
					}
					out.Users = append(out.Users, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities3(out *jwriter.Writer, in UserContainer) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"users\":")
	if in.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in.Users {
			if v5 > 0 {
				out.RawByte(',')
			}
			if v6 == nil {
				out.RawString("null")
			} else {
				(*v6).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserContainer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserContainer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserContainer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserContainer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities3(l, v)
}
func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities4(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "email":
			if in.IsNull() {
				in.Skip()
				out.Email = nil
			} else {
				if out.Email == nil {
					out.Email = new(string)
				}
				*out.Email = string(in.String())
			}
		case "first_name":
			if in.IsNull() {
				in.Skip()
				out.FirstName = nil
			} else {
				if out.FirstName == nil {
					out.FirstName = new(string)
				}
				*out.FirstName = string(in.String())
			}
		case "last_name":
			if in.IsNull() {
				in.Skip()
				out.LastName = nil
			} else {
				if out.LastName == nil {
					out.LastName = new(string)
				}
				*out.LastName = string(in.String())
			}
		case "gender":
			if in.IsNull() {
				in.Skip()
				out.Gender = nil
			} else {
				if out.Gender == nil {
					out.Gender = new(string)
				}
				*out.Gender = string(in.String())
			}
		case "birth_date":
			if in.IsNull() {
				in.Skip()
				out.Birthdate = nil
			} else {
				if out.Birthdate == nil {
					out.Birthdate = new(int64)
				}
				*out.Birthdate = int64(in.Int64())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities4(out *jwriter.Writer, in User) {
	if atomic.LoadInt32(&in.hasJSON) == 1 {
		out.Raw(in.json, nil)
		return
	}
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.Int64(int64(in.ID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"email\":")
	if in.Email == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Email))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"first_name\":")
	if in.FirstName == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.FirstName))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"last_name\":")
	if in.LastName == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.LastName))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"gender\":")
	if in.Gender == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Gender))
	}
	if in.Birthdate != nil {
		if !first {
			out.RawByte(',')
		}
		first = false
		out.RawString("\"birth_date\":")
		if in.Birthdate == nil {
			out.RawString("null")
		} else {
			out.Int64(int64(*in.Birthdate))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities4(l, v)
}
func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities5(in *jlexer.Lexer, out *ShortVisitContainer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "visits":
			if in.IsNull() {
				in.Skip()
				out.Visits = nil
			} else {
				in.Delim('[')
				if out.Visits == nil {
					if !in.IsDelim(']') {
						out.Visits = make([]*ShortVisit, 0, 8)
					} else {
						out.Visits = []*ShortVisit{}
					}
				} else {
					out.Visits = (out.Visits)[:0]
				}
				for !in.IsDelim(']') {
					var v7 *ShortVisit
					if in.IsNull() {
						in.Skip()
						v7 = nil
					} else {
						if v7 == nil {
							v7 = new(ShortVisit)
						}
						(*v7).UnmarshalEasyJSON(in)
					}
					out.Visits = append(out.Visits, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities5(out *jwriter.Writer, in ShortVisitContainer) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"visits\":")
	if in.Visits == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in.Visits {
			if v8 > 0 {
				out.RawByte(',')
			}
			if v9 == nil {
				out.RawString("null")
			} else {
				(*v9).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ShortVisitContainer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortVisitContainer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortVisitContainer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortVisitContainer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities5(l, v)
}
func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities6(in *jlexer.Lexer, out *ShortVisit) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "mark":
			out.Mark = uint8(in.Uint8())
		case "visited_at":
			out.VisitedAt = int64(in.Int64())
		case "place":
			out.Place = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities6(out *jwriter.Writer, in ShortVisit) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"mark\":")
	out.Uint8(uint8(in.Mark))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"visited_at\":")
	out.Int64(int64(in.VisitedAt))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"place\":")
	out.String(string(in.Place))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ShortVisit) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ShortVisit) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ShortVisit) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ShortVisit) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities6(l, v)
}
func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities7(in *jlexer.Lexer, out *LocationContainer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "locations":
			if in.IsNull() {
				in.Skip()
				out.Locations = nil
			} else {
				in.Delim('[')
				if out.Locations == nil {
					if !in.IsDelim(']') {
						out.Locations = make([]*Location, 0, 8)
					} else {
						out.Locations = []*Location{}
					}
				} else {
					out.Locations = (out.Locations)[:0]
				}
				for !in.IsDelim(']') {
					var v10 *Location
					if in.IsNull() {
						in.Skip()
						v10 = nil
					} else {
						if v10 == nil {
							v10 = new(Location)
						}
						(*v10).UnmarshalEasyJSON(in)
					}
					out.Locations = append(out.Locations, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities7(out *jwriter.Writer, in LocationContainer) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"locations\":")
	if in.Locations == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v11, v12 := range in.Locations {
			if v11 > 0 {
				out.RawByte(',')
			}
			if v12 == nil {
				out.RawString("null")
			} else {
				(*v12).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LocationContainer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LocationContainer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LocationContainer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LocationContainer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities7(l, v)
}
func easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities8(in *jlexer.Lexer, out *Location) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "distance":
			if in.IsNull() {
				in.Skip()
				out.Distance = nil
			} else {
				if out.Distance == nil {
					out.Distance = new(int64)
				}
				*out.Distance = int64(in.Int64())
			}
		case "city":
			if in.IsNull() {
				in.Skip()
				out.City = nil
			} else {
				if out.City == nil {
					out.City = new(string)
				}
				*out.City = string(in.String())
			}
		case "place":
			if in.IsNull() {
				in.Skip()
				out.Place = nil
			} else {
				if out.Place == nil {
					out.Place = new(string)
				}
				*out.Place = string(in.String())
			}
		case "country":
			if in.IsNull() {
				in.Skip()
				out.Country = nil
			} else {
				if out.Country == nil {
					out.Country = new(string)
				}
				*out.Country = string(in.String())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities8(out *jwriter.Writer, in Location) {
	if atomic.LoadInt32(&in.hasJSON) == 1 {
		out.Raw(in.json, nil)
		return
	}
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"id\":")
	out.Int64(int64(in.ID))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"distance\":")
	if in.Distance == nil {
		out.RawString("null")
	} else {
		out.Int64(int64(*in.Distance))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"city\":")
	if in.City == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.City))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"place\":")
	if in.Place == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Place))
	}
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"country\":")
	if in.Country == nil {
		out.RawString("null")
	} else {
		out.String(string(*in.Country))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Location) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Location) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson3e8ab7adEncodeGithubComUrakozzHighloadcampEntities8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Location) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Location) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson3e8ab7adDecodeGithubComUrakozzHighloadcampEntities8(l, v)
}
