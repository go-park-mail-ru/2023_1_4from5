// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels(in *jlexer.Lexer, out *SubscriptionDetails) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "creator_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CreatorId).UnmarshalText(data))
			}
		case "id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.Id).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "month_count":
			out.MonthCount = int64(in.Int64())
		case "money":
			out.Money = int64(in.Int64())
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
func easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels(out *jwriter.Writer, in SubscriptionDetails) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"creator_id\":"
		out.RawString(prefix[1:])
		out.RawText((in.CreatorId).MarshalText())
	}
	if true {
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.RawText((in.Id).MarshalText())
	}
	if true {
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.RawText((in.UserID).MarshalText())
	}
	{
		const prefix string = ",\"month_count\":"
		out.RawString(prefix)
		out.Int64(int64(in.MonthCount))
	}
	{
		const prefix string = ",\"money\":"
		out.RawString(prefix)
		out.Int64(int64(in.Money))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SubscriptionDetails) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SubscriptionDetails) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SubscriptionDetails) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SubscriptionDetails) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels(l, v)
}
func easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels1(in *jlexer.Lexer, out *Subscription) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.Id).UnmarshalText(data))
			}
		case "creator":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.Creator).UnmarshalText(data))
			}
		case "creator_name":
			out.CreatorName = string(in.String())
		case "creator_photo":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CreatorPhoto).UnmarshalText(data))
			}
		case "month_cost":
			out.MonthCost = int64(in.Int64())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
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
func easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels1(out *jwriter.Writer, in Subscription) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	if true {
		const prefix string = ",\"creator\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.Creator).MarshalText())
	}
	if in.CreatorName != "" {
		const prefix string = ",\"creator_name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.CreatorName))
	}
	if true {
		const prefix string = ",\"creator_photo\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.CreatorPhoto).MarshalText())
	}
	{
		const prefix string = ",\"month_cost\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.MonthCost))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	if in.Description != "" {
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Subscription) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Subscription) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Subscription) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Subscription) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels1(l, v)
}
