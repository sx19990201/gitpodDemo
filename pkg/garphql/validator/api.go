package validator

import (
	"context"
	"github.com/vektah/gqlparser/v2/ast"
)

// type ValidatorRes struct {
// 	Error    error
// 	FileName string
// }

type Validator interface {
	ValidateOperations(ctx context.Context, operationsRootPath, schemaFilePath string) (*ast.QueryDocument, error)
}
