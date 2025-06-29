package imgproc

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// ResizeImage resizes and compresses an image buffer according to parameters.
// Returns output bytes, format string, error
func ResizeImage(
	imageData []byte,
	width, height int,
	maintainAspect bool,
	quality int,
	maxSizeKB int,
) ([]byte, string, error) {
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", err
	}

	// Resize
	var dst *image.NRGBA
	if maintainAspect || width == 0 || height == 0 {
		dst = imaging.Resize(img, width, height, imaging.Lanczos)
	} else {
		dst = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	}

	var buf bytes.Buffer
	save := func(q int) error {
		switch strings.ToLower(format) {
		case "jpeg", "jpg":
			return jpeg.Encode(&buf, dst, &jpeg.Options{Quality: q})
		case "png":
			return png.Encode(&buf, dst)
		case "gif":
			return gif.Encode(&buf, dst, nil)
		case "webp":
			return webp.Encode(&buf, dst, &webp.Options{Quality: float32(q)})
		default:
			return jpeg.Encode(&buf, dst, &jpeg.Options{Quality: q})
		}
	}

	// Try with user quality
	if err := save(quality); err != nil {
		return nil, "", err
	}

	// If maxSizeKB is set, iteratively reduce quality to fit
	if maxSizeKB > 0 && (format == "jpeg" || format == "jpg" || format == "webp") {
		for q := quality; q >= 10; q -= 5 {
			buf.Reset()
			if err := save(q); err != nil {
				return nil, "", err
			}
			if buf.Len() <= maxSizeKB*1024 {
				break
			}
		}
		// If can't fit, still return the smallest
		if buf.Len() > maxSizeKB*1024 {
			return buf.Bytes(), format, errors.New("could not fit image into specified max size; returned lowest quality")
		}
	}

	return buf.Bytes(), format, nil
}
