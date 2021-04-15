package util

import (
    "context"
    "io"
)

type PipeStream struct {
    ErrorChannel chan error
    Blocked      bool
}

func (p *PipeStream) Block() {
    p.Blocked = true
}

func (p *PipeStream) Unblock() {
    p.Blocked = false
}

func Pipe(ctx context.Context, src io.ReadCloser, dst io.WriteCloser) *PipeStream {
    pipe := &PipeStream{
        ErrorChannel: make(chan error),
        Blocked:      false,
    }

    go func() {
        buffer := make([]byte, 4096)

        for {
            select {
            case <-ctx.Done():
                return
            default:
                if pipe.Blocked {
                    continue
                }

                n, err := src.Read(buffer)
                if err != nil {
                    pipe.ErrorChannel <- err

                    if err == io.EOF {
                        continue
                    }

                    return
                }

                _, err = dst.Write(buffer[0:n])
                if err != nil {
                    pipe.ErrorChannel <- err
                    return
                }
            }
        }
    }()

    return pipe
}
