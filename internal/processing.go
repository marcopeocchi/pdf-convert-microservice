package internal

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/h2non/bimg"
)

type ProcessingOptions struct {
	Image   *bimg.Image
	ImgType bimg.ImageType
	Quality int
}

func Process(ctx context.Context, w io.Writer, opt *ProcessingOptions) error {
	select {
	case <-ctx.Done():
		return errors.New("context cancelled")

	default:
		return process(w, opt)
	}
}

func process(w io.Writer, opt *ProcessingOptions) error {
	start := time.Now()

	defer func() {
		TimePerOpGauge.Set(float64(time.Since(start)))
		OpsCounter.Inc()
	}()

	buf, err := opt.Image.Process(bimg.Options{
		Type:    opt.ImgType,
		Quality: opt.Quality,
	})
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
