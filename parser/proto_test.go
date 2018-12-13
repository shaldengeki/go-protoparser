package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/go-protoparser/internal/lexer"
	"github.com/yoheimuta/go-protoparser/internal/util_test"
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
)

func TestParser_ParseProto(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		filename  string
		wantProto *parser.Proto
		wantErr   bool
	}{
		{
			name:    "parsing an empty",
			wantErr: true,
		},
		{
			name: "parsing an excerpt from the official reference",
			input: `
syntax = "proto3";
import public "other.proto";
option java_package = "com.example.foo";
enum EnumAllowingAlias {
  option allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 2 [(custom_option) = "hello world"];
}
message outer {
  option (my_option).a = true;
  message inner {
    int64 ival = 1;
  }
  repeated inner inner_message = 2;
  EnumAllowingAlias enum_field =3;
  map<int32, string> my_map = 4;
}
`,
			filename: "official.proto",
			wantProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
					Meta: meta.Meta{
						Pos: meta.Position{
							Filename: "official.proto",
							Offset:   1,
							Line:     2,
							Column:   1,
						},
					},
				},
				ProtoBody: []parser.Visitee{
					&parser.Import{
						Modifier: parser.ImportModifierPublic,
						Location: `"other.proto"`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "official.proto",
								Offset:   20,
								Line:     3,
								Column:   1,
							},
						},
					},
					&parser.Option{
						OptionName: "java_package",
						Constant:   `"com.example.foo"`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "official.proto",
								Offset:   49,
								Line:     4,
								Column:   1,
							},
						},
					},
					&parser.Enum{
						EnumName: "EnumAllowingAlias",
						EnumBody: []parser.Visitee{
							&parser.Option{
								OptionName: "allow_alias",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   117,
										Line:     6,
										Column:   3,
									},
								},
							},
							&parser.EnumField{
								Ident:  "UNKNOWN",
								Number: "0",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   146,
										Line:     7,
										Column:   3,
									},
								},
							},
							&parser.EnumField{
								Ident:  "STARTED",
								Number: "1",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   161,
										Line:     8,
										Column:   3,
									},
								},
							},
							&parser.EnumField{
								Ident:  "RUNNING",
								Number: "2",
								EnumValueOptions: []*parser.EnumValueOption{
									{
										OptionName: "(custom_option)",
										Constant:   `"hello world"`,
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   176,
										Line:     9,
										Column:   3,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "official.proto",
								Offset:   90,
								Line:     5,
								Column:   1,
							},
						},
					},
					&parser.Message{
						MessageName: "outer",
						MessageBody: []parser.Visitee{
							&parser.Option{
								OptionName: "(my_option).a",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   243,
										Line:     12,
										Column:   3,
									},
								},
							},
							&parser.Message{
								MessageName: "inner",
								MessageBody: []parser.Visitee{
									&parser.Field{
										Type:        "int64",
										FieldName:   "ival",
										FieldNumber: "1",
										Meta: meta.Meta{
											Pos: meta.Position{
												Filename: "official.proto",
												Offset:   294,
												Line:     14,
												Column:   5,
											},
										},
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   274,
										Line:     13,
										Column:   3,
									},
								},
							},
							&parser.Field{
								IsRepeated:  true,
								Type:        "inner",
								FieldName:   "inner_message",
								FieldNumber: "2",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   316,
										Line:     16,
										Column:   3,
									},
								},
							},
							&parser.Field{
								Type:        "EnumAllowingAlias",
								FieldName:   "enum_field",
								FieldNumber: "3",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   352,
										Line:     17,
										Column:   3,
									},
								},
							},
							&parser.MapField{
								KeyType:     "int32",
								Type:        "string",
								MapName:     "my_map",
								FieldNumber: "4",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "official.proto",
										Offset:   387,
										Line:     18,
										Column:   3,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "official.proto",
								Offset:   225,
								Line:     11,
								Column:   1,
							},
						},
					},
				},
				Meta: &parser.ProtoMeta{
					Filename: "official.proto",
				},
			},
		},
		{
			name: "parsing a service",
			input: `
syntax = "proto3";
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
}
`,
			filename: "service.proto",
			wantProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
					Meta: meta.Meta{
						Pos: meta.Position{
							Filename: "service.proto",
							Offset:   1,
							Line:     2,
							Column:   1,
						},
					},
				},
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceName: "SearchService",
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "Search",
								RPCRequest: &parser.RPCRequest{
									MessageType: "SearchRequest",
									Meta: meta.Meta{
										Pos: meta.Position{
											Filename: "service.proto",
											Offset:   57,
											Line:     4,
											Column:   14,
										},
									},
								},
								RPCResponse: &parser.RPCResponse{
									MessageType: "SearchResponse",
									Meta: meta.Meta{
										Pos: meta.Position{
											Filename: "service.proto",
											Offset:   81,
											Line:     4,
											Column:   38,
										},
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "service.proto",
										Offset:   46,
										Line:     4,
										Column:   3,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "service.proto",
								Offset:   20,
								Line:     3,
								Column:   1,
							},
						},
					},
				},
				Meta: &parser.ProtoMeta{
					Filename: "service.proto",
				},
			},
		},
		{
			name: "parsing comments",
			input: `
// syntax
/*
syntax2
*/
syntax = "proto3";
// import
import public "other.proto";
/* package */
package foo.bar;
// option
option java_package = "com.example.foo";
// message
message outer {
}
// enum
enum EnumAllowingAlias {
  option allow_alias = true;
}
// service
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
}
`,
			filename: "comments.proto",
			wantProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
					Comments: []*parser.Comment{
						{
							Raw: `// syntax`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Filename: "comments.proto",
									Offset:   1,
									Line:     2,
									Column:   1,
								},
							},
						},
						{
							Raw: `/*
syntax2
*/`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Filename: "comments.proto",
									Offset:   11,
									Line:     3,
									Column:   1,
								},
							},
						},
					},
					Meta: meta.Meta{
						Pos: meta.Position{
							Filename: "comments.proto",
							Offset:   25,
							Line:     6,
							Column:   1,
						},
					},
				},
				ProtoBody: []parser.Visitee{
					&parser.Import{
						Modifier: parser.ImportModifierPublic,
						Location: `"other.proto"`,
						Comments: []*parser.Comment{
							{
								Raw: `// import`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "comments.proto",
										Offset:   44,
										Line:     7,
										Column:   1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "comments.proto",
								Offset:   54,
								Line:     8,
								Column:   1,
							},
						},
					},
					&parser.Package{
						Name: `foo.bar`,
						Comments: []*parser.Comment{
							{
								Raw: `/* package */`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "comments.proto",
										Offset:   83,
										Line:     9,
										Column:   1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "comments.proto",
								Offset:   97,
								Line:     10,
								Column:   1,
							},
						},
					},
					&parser.Option{
						OptionName: "java_package",
						Constant:   `"com.example.foo"`,
						Comments: []*parser.Comment{
							{
								Raw: `// option`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "comments.proto",
										Offset:   114,
										Line:     11,
										Column:   1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "comments.proto",
								Offset:   124,
								Line:     12,
								Column:   1,
							},
						},
					},
					&parser.Message{
						MessageName: "outer",
						Comments: []*parser.Comment{
							{
								Raw: `// message`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "comments.proto",
										Offset:   165,
										Line:     13,
										Column:   1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "comments.proto",
								Offset:   176,
								Line:     14,
								Column:   1,
							},
						},
					},
					&parser.Enum{
						EnumName: "EnumAllowingAlias",
						EnumBody: []parser.Visitee{
							&parser.Option{
								OptionName: "allow_alias",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "comments.proto",
										Offset:   229,
										Line:     18,
										Column:   3,
									},
								},
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// enum`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "comments.proto",
										Offset:   194,
										Line:     16,
										Column:   1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "comments.proto",
								Offset:   202,
								Line:     17,
								Column:   1,
							},
						},
					},
					&parser.Service{
						ServiceName: "SearchService",
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "Search",
								RPCRequest: &parser.RPCRequest{
									MessageType: "SearchRequest",
									Meta: meta.Meta{
										Pos: meta.Position{
											Filename: "comments.proto",
											Offset:   306,
											Line:     22,
											Column:   14,
										},
									},
								},
								RPCResponse: &parser.RPCResponse{
									MessageType: "SearchResponse",
									Meta: meta.Meta{
										Pos: meta.Position{
											Filename: "comments.proto",
											Offset:   330,
											Line:     22,
											Column:   38,
										},
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "comments.proto",
										Offset:   295,
										Line:     22,
										Column:   3,
									},
								},
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// service`,
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "comments.proto",
										Offset:   258,
										Line:     20,
										Column:   1,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "comments.proto",
								Offset:   269,
								Line:     21,
								Column:   1,
							},
						},
					},
				},
				Meta: &parser.ProtoMeta{
					Filename: "comments.proto",
				},
			},
		},
		{
			name: "parsing inline comments",
			input: `
syntax = "proto3"; // syntax
import public "other.proto"; // import
package foo.bar; /* package */
option java_package = "com.example.foo"; // option
message outer {
} // message
enum EnumAllowingAlias {
  option allow_alias = true;
} // enum
service SearchService {
  rpc Search (SearchRequest) returns (SearchResponse);
} // service
`,
			filename: "inlineComments.proto",
			wantProto: &parser.Proto{
				Syntax: &parser.Syntax{
					ProtobufVersion: "proto3",
					InlineComment: &parser.Comment{
						Raw: `// syntax`,
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "inlineComments.proto",
								Offset:   20,
								Line:     2,
								Column:   20,
							},
						},
					},
					Meta: meta.Meta{
						Pos: meta.Position{
							Filename: "inlineComments.proto",
							Offset:   1,
							Line:     2,
							Column:   1,
						},
					},
				},
				ProtoBody: []parser.Visitee{
					&parser.Import{
						Modifier: parser.ImportModifierPublic,
						Location: `"other.proto"`,
						InlineComment: &parser.Comment{
							Raw: `// import`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Filename: "inlineComments.proto",
									Offset:   59,
									Line:     3,
									Column:   30,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "inlineComments.proto",
								Offset:   30,
								Line:     3,
								Column:   1,
							},
						},
					},
					&parser.Package{
						Name: `foo.bar`,
						InlineComment: &parser.Comment{
							Raw: `/* package */`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Filename: "inlineComments.proto",
									Offset:   86,
									Line:     4,
									Column:   18,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "inlineComments.proto",
								Offset:   69,
								Line:     4,
								Column:   1,
							},
						},
					},
					&parser.Option{
						OptionName: "java_package",
						Constant:   `"com.example.foo"`,
						InlineComment: &parser.Comment{
							Raw: `// option`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Filename: "inlineComments.proto",
									Offset:   141,
									Line:     5,
									Column:   42,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "inlineComments.proto",
								Offset:   100,
								Line:     5,
								Column:   1,
							},
						},
					},
					&parser.Message{
						MessageName: "outer",
						InlineComment: &parser.Comment{
							Raw: `// message`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Filename: "inlineComments.proto",
									Offset:   169,
									Line:     7,
									Column:   3,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "inlineComments.proto",
								Offset:   151,
								Line:     6,
								Column:   1,
							},
						},
					},
					&parser.Enum{
						EnumName: "EnumAllowingAlias",
						EnumBody: []parser.Visitee{
							&parser.Option{
								OptionName: "allow_alias",
								Constant:   "true",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "inlineComments.proto",
										Offset:   207,
										Line:     9,
										Column:   3,
									},
								},
							},
						},
						InlineComment: &parser.Comment{
							Raw: `// enum`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Filename: "inlineComments.proto",
									Offset:   236,
									Line:     10,
									Column:   3,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "inlineComments.proto",
								Offset:   180,
								Line:     8,
								Column:   1,
							},
						},
					},
					&parser.Service{
						ServiceName: "SearchService",
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "Search",
								RPCRequest: &parser.RPCRequest{
									MessageType: "SearchRequest",
									Meta: meta.Meta{
										Pos: meta.Position{
											Filename: "inlineComments.proto",
											Offset:   281,
											Line:     12,
											Column:   14,
										},
									},
								},
								RPCResponse: &parser.RPCResponse{
									MessageType: "SearchResponse",
									Meta: meta.Meta{
										Pos: meta.Position{
											Filename: "inlineComments.proto",
											Offset:   305,
											Line:     12,
											Column:   38,
										},
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "inlineComments.proto",
										Offset:   270,
										Line:     12,
										Column:   3,
									},
								},
							},
						},
						InlineComment: &parser.Comment{
							Raw: `// service`,
							Meta: meta.Meta{
								Pos: meta.Position{
									Filename: "inlineComments.proto",
									Offset:   325,
									Line:     13,
									Column:   3,
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "inlineComments.proto",
								Offset:   244,
								Line:     11,
								Column:   1,
							},
						},
					},
				},
				Meta: &parser.ProtoMeta{
					Filename: "inlineComments.proto",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			p := parser.NewParser(
				lexer.NewLexer(
					strings.NewReader(test.input),
					lexer.WithFilename(test.filename),
				),
			)
			got, err := p.ParseProto()
			switch {
			case test.wantErr:
				if err == nil {
					t.Errorf("got err nil, but want err, parsed=%v", got)
				}
				return
			case !test.wantErr && err != nil:
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantProto) {
				t.Errorf("got %v, but want %v", util_test.PrettyFormat(got), util_test.PrettyFormat(test.wantProto))
			}

			if !p.IsEOF() {
				t.Errorf("got not eof, but want eof")
			}
		})
	}
}
