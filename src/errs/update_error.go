package errs

import "fmt"

// UpdateSource 更新操作的类型
type UpdateSource string

const (
	DB UpdateSource = "DB"
	ES UpdateSource = "ES"
)

type UpdateError struct {
	Source UpdateSource
	Err    error
}

func (e UpdateError) Error() string {
	return fmt.Sprintf("update error from %s: %v", e.Source, e.Err)
}

// NewUpdateError 创建UpdateError实例
func NewUpdateError(source UpdateSource, err error) UpdateError {
	return UpdateError{
		Source: source,
		Err:    err,
	}
}
