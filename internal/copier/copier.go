package copier

import (
	"fmt"
	"io"
)

const BufferSize = 20

//CopyParams copy parameters structure
type CopyParams struct {
	Limit  int64
	Offset int64
}

//CopyBuffered copy from reader to dst using buffer
func CopyBuffered(src io.Reader, dst io.Writer, param CopyParams) (int64, error) {
	var bytesCopied int64
	buf := make([]byte, BufferSize)
	for {
		if bytesCopied >= param.Limit {
			break
		}
		n, err := src.Read(buf)
		if err != nil {
			switch err {
			case io.EOF:
				return bytesCopied, nil
			default:
				return bytesCopied, err
			}
		}
		if n == 0 {
			break
		}
		var writeCount = n
		if int64(writeCount)+bytesCopied > param.Limit {
			writeCount = int(param.Limit - bytesCopied)
		}
		if _, err := dst.Write(buf[:writeCount]); err != nil {
			return bytesCopied, err
		}
		bytesCopied = bytesCopied + int64(writeCount)
		fmt.Printf("bytes copied %d\n", writeCount)
	}
	return bytesCopied, nil
}
