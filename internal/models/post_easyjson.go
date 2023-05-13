// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	uuid "github.com/google/uuid"
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

func easyjson5a72dc82DecodeGithubComGoParkMailRu202314from5InternalModels(in *jlexer.Lexer, out *PostEditData) {
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
		case "title":
			out.Title = string(in.String())
		case "text":
			out.Text = string(in.String())
		case "available_subscriptions":
			if in.IsNull() {
				in.Skip()
				out.AvailableSubscriptions = nil
			} else {
				in.Delim('[')
				if out.AvailableSubscriptions == nil {
					if !in.IsDelim(']') {
						out.AvailableSubscriptions = make([]uuid.UUID, 0, 4)
					} else {
						out.AvailableSubscriptions = []uuid.UUID{}
					}
				} else {
					out.AvailableSubscriptions = (out.AvailableSubscriptions)[:0]
				}
				for !in.IsDelim(']') {
					var v1 uuid.UUID
					if data := in.UnsafeBytes(); in.Ok() {
						in.AddError((v1).UnmarshalText(data))
					}
					out.AvailableSubscriptions = append(out.AvailableSubscriptions, v1)
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
func easyjson5a72dc82EncodeGithubComGoParkMailRu202314from5InternalModels(out *jwriter.Writer, in PostEditData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"available_subscriptions\":"
		out.RawString(prefix)
		if in.AvailableSubscriptions == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.AvailableSubscriptions {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.RawText((v3).MarshalText())
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostEditData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComGoParkMailRu202314from5InternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostEditData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComGoParkMailRu202314from5InternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostEditData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComGoParkMailRu202314from5InternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostEditData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComGoParkMailRu202314from5InternalModels(l, v)
}
func easyjson5a72dc82DecodeGithubComGoParkMailRu202314from5InternalModels1(in *jlexer.Lexer, out *Post) {
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
		case "creator_photo":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CreatorPhoto).UnmarshalText(data))
			}
		case "creator_name":
			out.CreatorName = string(in.String())
		case "creation_date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Creation).UnmarshalJSON(data))
			}
		case "likes_count":
			out.LikesCount = int64(in.Int64())
		case "comments_count":
			out.CommentsCount = int64(in.Int64())
		case "title":
			out.Title = string(in.String())
		case "text":
			out.Text = string(in.String())
		case "is_available":
			out.IsAvailable = bool(in.Bool())
		case "is_liked":
			out.IsLiked = bool(in.Bool())
		case "attachments":
			if in.IsNull() {
				in.Skip()
				out.Attachments = nil
			} else {
				in.Delim('[')
				if out.Attachments == nil {
					if !in.IsDelim(']') {
						out.Attachments = make([]Attachment, 0, 2)
					} else {
						out.Attachments = []Attachment{}
					}
				} else {
					out.Attachments = (out.Attachments)[:0]
				}
				for !in.IsDelim(']') {
					var v4 Attachment
					(v4).UnmarshalEasyJSON(in)
					out.Attachments = append(out.Attachments, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "subscriptions":
			if in.IsNull() {
				in.Skip()
				out.Subscriptions = nil
			} else {
				in.Delim('[')
				if out.Subscriptions == nil {
					if !in.IsDelim(']') {
						out.Subscriptions = make([]Subscription, 0, 0)
					} else {
						out.Subscriptions = []Subscription{}
					}
				} else {
					out.Subscriptions = (out.Subscriptions)[:0]
				}
				for !in.IsDelim(']') {
					var v5 Subscription
					(v5).UnmarshalEasyJSON(in)
					out.Subscriptions = append(out.Subscriptions, v5)
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
func easyjson5a72dc82EncodeGithubComGoParkMailRu202314from5InternalModels1(out *jwriter.Writer, in Post) {
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
	if true {
		const prefix string = ",\"creator_photo\":"
		out.RawString(prefix)
		out.RawText((in.CreatorPhoto).MarshalText())
	}
	if in.CreatorName != "" {
		const prefix string = ",\"creator_name\":"
		out.RawString(prefix)
		out.String(string(in.CreatorName))
	}
	{
		const prefix string = ",\"creation_date\":"
		out.RawString(prefix)
		out.Raw((in.Creation).MarshalJSON())
	}
	{
		const prefix string = ",\"likes_count\":"
		out.RawString(prefix)
		out.Int64(int64(in.LikesCount))
	}
	{
		const prefix string = ",\"comments_count\":"
		out.RawString(prefix)
		out.Int64(int64(in.CommentsCount))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"is_available\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsAvailable))
	}
	{
		const prefix string = ",\"is_liked\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsLiked))
	}
	{
		const prefix string = ",\"attachments\":"
		out.RawString(prefix)
		if in.Attachments == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v6, v7 := range in.Attachments {
				if v6 > 0 {
					out.RawByte(',')
				}
				(v7).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"subscriptions\":"
		out.RawString(prefix)
		if in.Subscriptions == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Subscriptions {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComGoParkMailRu202314from5InternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComGoParkMailRu202314from5InternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComGoParkMailRu202314from5InternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComGoParkMailRu202314from5InternalModels1(l, v)
}
