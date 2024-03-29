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

func easyjsonE9abebc9DecodeGithubComGoParkMailRu202314from5InternalModels(in *jlexer.Lexer, out *Comment) {
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
		case "comment_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.CommentID).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "username":
			out.Username = string(in.String())
		case "user_photo":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserPhoto).UnmarshalText(data))
			}
		case "post_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.PostID).UnmarshalText(data))
			}
		case "text":
			out.Text = string(in.String())
		case "creation":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Creation).UnmarshalJSON(data))
			}
		case "likes_count":
			out.LikesCount = int64(in.Int64())
		case "is_liked":
			out.IsLiked = bool(in.Bool())
		case "is_owner":
			out.IsOwner = bool(in.Bool())
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
func easyjsonE9abebc9EncodeGithubComGoParkMailRu202314from5InternalModels(out *jwriter.Writer, in Comment) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		const prefix string = ",\"comment_id\":"
		first = false
		out.RawString(prefix[1:])
		out.RawText((in.CommentID).MarshalText())
	}
	if true {
		const prefix string = ",\"user_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.UserID).MarshalText())
	}
	if in.Username != "" {
		const prefix string = ",\"username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	if true {
		const prefix string = ",\"user_photo\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.UserPhoto).MarshalText())
	}
	if true {
		const prefix string = ",\"post_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.RawText((in.PostID).MarshalText())
	}
	if in.Text != "" {
		const prefix string = ",\"text\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Text))
	}
	if true {
		const prefix string = ",\"creation\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Creation).MarshalJSON())
	}
	if in.LikesCount != 0 {
		const prefix string = ",\"likes_count\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.LikesCount))
	}
	{
		const prefix string = ",\"is_liked\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.IsLiked))
	}
	{
		const prefix string = ",\"is_owner\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsOwner))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (comment Comment) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE9abebc9EncodeGithubComGoParkMailRu202314from5InternalModels(&w, comment)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (comment Comment) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE9abebc9EncodeGithubComGoParkMailRu202314from5InternalModels(w, comment)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (comment *Comment) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE9abebc9DecodeGithubComGoParkMailRu202314from5InternalModels(&r, comment)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (comment *Comment) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE9abebc9DecodeGithubComGoParkMailRu202314from5InternalModels(l, comment)
}
