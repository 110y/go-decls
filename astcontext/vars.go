package astcontext

import (
	"go/ast"
	"go/token"
)

type Var struct {
	Name   string    `json:"name" vim:"name"`
	VarPos *Position `json:"var" vim:"var"`

	node *ast.ValueSpec
}

type Vars []*Var

func (p *Parser) Vars() []*Decl {
	var files []*ast.File
	if p.file != nil {
		files = append(files, p.file)
	}

	if p.pkgs != nil {
		for _, pkg := range p.pkgs {
			for _, f := range pkg.Files {
				files = append(files, f)
			}
		}
	}

	// var vars []*Var
	// var typs []*Type
	var decls []*Decl

	for _, f := range files {
		for _, decl := range f.Decls {
			switch x := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range x.Specs {
					switch x.Tok {
					case token.VAR:
						s, ok := spec.(*ast.ValueSpec)
						if ok {
							v := &Var{
								Name:   s.Names[0].Name,
								VarPos: ToPosition(p.fset.Position(s.Names[0].Pos())),
								node:   s,
							}

							if v.Name != "_" {
								decls = append(decls, &Decl{
									Keyword:  "Var",
									Ident:    v.Name,
									Full:     v.Name,
									Filename: v.VarPos.Filename,
									Line:     v.VarPos.Line,
									Col:      v.VarPos.Column,
								})
							}
						}
					case token.CONST:
						s, ok := spec.(*ast.ValueSpec)
						if ok {
							v := &Var{
								Name:   s.Names[0].Name,
								VarPos: ToPosition(p.fset.Position(s.Names[0].Pos())),
								node:   s,
							}

							if v.Name != "_" {
								decls = append(decls, &Decl{
									Keyword:  "Const",
									Ident:    v.Name,
									Full:     v.Name,
									Filename: v.VarPos.Filename,
									Line:     v.VarPos.Line,
									Col:      v.VarPos.Column,
								})
							}
						}
					case token.TYPE:
						s, ok := spec.(*ast.TypeSpec)
						if ok {
							var keyword string
							_, ok := s.Type.(*ast.StructType)
							if ok {
								keyword = "Struct"
							} else if _, ok := s.Type.(*ast.InterfaceType); ok {
								keyword = "Interface"
							} else {
								keyword = "Type"
							}

							t := &Type{
								TypePos: ToPosition(p.fset.Position(s.Name.Pos())),
								node:    s,
							}

							if s.Doc != nil {
								t.Doc = ToPosition(p.fset.Position(x.Doc.Pos()))
							}

							t.Signature = NewTypeSignature(s)

							decls = append(decls, &Decl{
								Keyword:  keyword,
								Ident:    t.Signature.Name,
								Full:     t.Signature.Full,
								Filename: t.TypePos.Filename,
								Line:     t.TypePos.Line,
								Col:      t.TypePos.Column,
							})
						}
					}
				}
			case *ast.FuncDecl:
				fn := &Func{
					FuncPos: ToPosition(p.fset.Position(x.Type.Func)),
					node:    x,
				}

				// can be nil for forward declarations
				if x.Body != nil {
					fn.Lbrace = ToPosition(p.fset.Position(x.Body.Lbrace))
					fn.Rbrace = ToPosition(p.fset.Position(x.Body.Rbrace))
				}

				if x.Doc != nil {
					fn.Doc = ToPosition(p.fset.Position(x.Doc.Pos()))
				}

				fn.Signature = NewFuncSignature(x)

				var keyword string
				if x.Recv != nil {
					keyword = "Method"
				} else {
					keyword = "Function"
				}

				decls = append(decls, &Decl{
					Keyword:  keyword,
					Ident:    fn.Signature.Name,
					Full:     fn.Signature.Full,
					Filename: fn.FuncPos.Filename,
					Line:     fn.FuncPos.Line,
					Col:      fn.FuncPos.Column,
				})
			}
		}
	}

	return decls
}

func (v Vars) TopLevel() Vars {
	return Vars{}
}
