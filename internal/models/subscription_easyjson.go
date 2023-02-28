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

func easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels(in *jlexer.Lexer, out *Subscription) {
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
		case "month_const":
			out.MonthConst = int(in.Int())
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
func easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels(out *jwriter.Writer, in Subscription) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	{
		const prefix string = ",\"creator\":"
		out.RawString(prefix)
		out.RawText((in.Creator).MarshalText())
	}
	{
		const prefix string = ",\"month_const\":"
		out.RawString(prefix)
		out.Int(int(in.MonthConst))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Subscription) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Subscription) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFfbd3743EncodeGithubComGoParkMailRu202314from5InternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Subscription) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Subscription) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFfbd3743DecodeGithubComGoParkMailRu202314from5InternalModels(l, v)
}
