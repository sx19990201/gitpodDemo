package validator

import (
	"context"
	"fmt"
	"os"

	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type v struct {
}

func NewValidator() Validator {
	return &v{}
}

func (_v *v) ValidateOperations(ctx context.Context, operationsRootPath, schemaFilePath string) (*ast.QueryDocument, error) {
	schema, err := _v.loadSchema(schemaFilePath)
	if err != nil {
		return nil, err
	}
	return _v.loadQuery(schema, operationsRootPath)
}

func (_v *v) loadSchema(filePath string) (*ast.Schema, error) {
	schemaBytes, err := loadFile(filePath)
	if err != nil {
		return nil, err
	}
	return gqlparser.LoadSchema(&ast.Source{Input: string(schemaBytes), Name: filePath})
}

func (_v *v) loadQuery(schemaContent *ast.Schema, operationPath string) (*ast.QueryDocument, error) {
	operationBytes, err := loadFile(operationPath)
	if err != nil {
		return nil, err
	}
	return gqlparser.LoadQuery(schemaContent, string(operationBytes))
}

func loadFile(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", filePath)
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file=%s meet err=%v", filePath, err)
	}

	return bytes, nil
}
