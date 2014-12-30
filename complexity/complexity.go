package complexity

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

func CheckFiles(filepath string, limit int) {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, filepath, nil, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	a, err := BuildAst("", f)

	decl := a.GetChildByString("Type", "Decl")
	score := 0
	for _, ch := range decl.Children {
		if strings.Contains(ch.Attrs["Type"], "FuncDecl") {
			score = CyclomaticComplexity(ch)
			if score > limit {
				Position := fset.Position(ch.Pos)
				fmt.Printf("%s :Line %d func %s is too complex %d > %d\n", filepath, Position.Line, ch.Attrs["Name"], score, limit)
			}
		}
	}
}

func CyclomaticComplexity(a *Ast) (score int) {
	score = 0
	body := a.GetChildByString("Prefix", "Body")
	if body != nil {
		for _, child := range body.Children {
			score += CyclomaticComplexity(child)
		}
	}
	if strings.Contains(a.Attrs["Type"], "List") || strings.Contains(a.Attrs["Prefix"], "List") {
		for _, child := range a.Children {
			score += CyclomaticComplexity(child)
		}
	}
	switch {
	case strings.Contains(a.Attrs["Type"], "IfStmt"):
		// count how many conds
		conds := a.GetChildByString("Prefix", "Cond")
		if conds != nil {
			score += CountConds(conds)
		}
	case strings.Contains(a.Attrs["Type"], "ForStmt"):
		score += 1
	case strings.Contains(a.Attrs["Type"], "CaseClause"):
		// count how many cases
		score += 1
	}
	return score
}

func (a *Ast) GetChildByString(key, name string) *Ast {
	for _, child := range a.Children {
		if strings.Contains(child.Attrs[key], name) {
			return child
		}
	}
	return nil
}

func CountConds(a *Ast) int {
	count := 0
	if strings.Contains(a.Attrs["Type"], "BinaryExpr") {
		for _, child := range a.Children {
			count += CountConds(child)
		}
		return count
	} else {
		return 1
	}
}
