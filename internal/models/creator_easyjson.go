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

func easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels(in *jlexer.Lexer, out *UpdateCreatorInfo) {
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
		case "description":
			out.Description = string(in.String())
		case "creator_name":
			out.CreatorName = string(in.String())
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
func easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels(out *jwriter.Writer, in UpdateCreatorInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix[1:])
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"creator_name\":"
		out.RawString(prefix)
		out.String(string(in.CreatorName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UpdateCreatorInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UpdateCreatorInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UpdateCreatorInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UpdateCreatorInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels(l, v)
}
func easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels1(in *jlexer.Lexer, out *CreatorPage) {
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
		case "creator_info":
			(out.CreatorInfo).UnmarshalEasyJSON(in)
		case "aim":
			(out.Aim).UnmarshalEasyJSON(in)
		case "is_my_page":
			out.IsMyPage = bool(in.Bool())
		case "follows":
			out.Follows = bool(in.Bool())
		case "posts":
			if in.IsNull() {
				in.Skip()
				out.Posts = nil
			} else {
				in.Delim('[')
				if out.Posts == nil {
					if !in.IsDelim(']') {
						out.Posts = make([]Post, 0, 0)
					} else {
						out.Posts = []Post{}
					}
				} else {
					out.Posts = (out.Posts)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Post
					(v1).UnmarshalEasyJSON(in)
					out.Posts = append(out.Posts, v1)
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
					var v2 Subscription
					(v2).UnmarshalEasyJSON(in)
					out.Subscriptions = append(out.Subscriptions, v2)
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
func easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels1(out *jwriter.Writer, in CreatorPage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"creator_info\":"
		out.RawString(prefix[1:])
		(in.CreatorInfo).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"aim\":"
		out.RawString(prefix)
		(in.Aim).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"is_my_page\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsMyPage))
	}
	{
		const prefix string = ",\"follows\":"
		out.RawString(prefix)
		out.Bool(bool(in.Follows))
	}
	{
		const prefix string = ",\"posts\":"
		out.RawString(prefix)
		if in.Posts == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Posts {
				if v3 > 0 {
					out.RawByte(',')
				}
				(v4).MarshalEasyJSON(out)
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
			for v5, v6 := range in.Subscriptions {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreatorPage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreatorPage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreatorPage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreatorPage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels1(l, v)
}
func easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels2(in *jlexer.Lexer, out *Creator) {
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
				in.AddError((out.Id).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserId).UnmarshalText(data))
			}
		case "name":
			out.Name = string(in.String())
		case "cover_photo":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CoverPhoto).UnmarshalText(data))
			}
		case "profile_photo":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ProfilePhoto).UnmarshalText(data))
			}
		case "followers_count":
			out.FollowersCount = int64(in.Int64())
		case "description":
			out.Description = string(in.String())
		case "posts_count":
			out.PostsCount = int64(in.Int64())
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
func easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels2(out *jwriter.Writer, in Creator) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"creator_id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Id).MarshalText())
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.RawText((in.UserId).MarshalText())
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"cover_photo\":"
		out.RawString(prefix)
		out.RawText((in.CoverPhoto).MarshalText())
	}
	{
		const prefix string = ",\"profile_photo\":"
		out.RawString(prefix)
		out.RawText((in.ProfilePhoto).MarshalText())
	}
	{
		const prefix string = ",\"followers_count\":"
		out.RawString(prefix)
		out.Int64(int64(in.FollowersCount))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"posts_count\":"
		out.RawString(prefix)
		out.Int64(int64(in.PostsCount))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Creator) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Creator) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Creator) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Creator) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels2(l, v)
}
func easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels3(in *jlexer.Lexer, out *Aim) {
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
				in.AddError((out.Creator).UnmarshalText(data))
			}
		case "description":
			out.Description = string(in.String())
		case "money_needed":
			out.MoneyNeeded = int64(in.Int64())
		case "money_got":
			out.MoneyGot = int64(in.Int64())
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
func easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels3(out *jwriter.Writer, in Aim) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"creator_id\":"
		out.RawString(prefix[1:])
		out.RawText((in.Creator).MarshalText())
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"money_needed\":"
		out.RawString(prefix)
		out.Int64(int64(in.MoneyNeeded))
	}
	{
		const prefix string = ",\"money_got\":"
		out.RawString(prefix)
		out.Int64(int64(in.MoneyGot))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Aim) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Aim) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson7c25d2a6EncodeGithubComGoParkMailRu202314from5InternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Aim) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Aim) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson7c25d2a6DecodeGithubComGoParkMailRu202314from5InternalModels3(l, v)
}
