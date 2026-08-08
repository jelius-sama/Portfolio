package blogs

import (
    "errors"
    "io"
    "os"
    "path/filepath"

    "git.jelius.dev/jelius-sama/Portfolio/types"
    "github.com/gofiber/fiber/v3"
)

//	func exampleCaller() {
//	    var stream io.ReadCloser = io.NopCloser(strings.NewReader("/blog/abc1230"))
//
//	    if err := GetBlogMarkdown(nil, &stream); err != nil {
//	        return
//	    }
//	    // 'stream' has now been mutated! It is now an *os.File descriptor.
//	    // We MUST defer closing it to avoid resource leaks.
//	    defer stream.Close()
//
//	    // Manually step through the file chunk-by-chunk
//	    var chunkBuffer = make([]byte, 4096)
//	    for {
//	        n, err := stream.Read(chunkBuffer)
//	        if n > 0 {
//	            // Process chunkBuffer[:n] manually
//	        }
//	        if err == io.EOF {
//	            break // Finished chunk-by-chunk manual reading
//	        }
//	    }
//	}
func GetBlogMarkdown(c fiber.Ctx, buf ...*io.ReadCloser) error {
    var id string

    if has := len(buf) > 0; has && buf[0] == nil {
        return errors.New("length of variadic greater than 0 but nil buffer")
    } else if has {
        if idBytes, err := io.ReadAll(*buf[0]); err == nil {
            id = string(idBytes)
            (*buf[0]).Close()
        }
    } else {
        id = c.Params("id")
    }

    // Validate ID
    if len(id) == 0 {
        var err = errors.New("blog id is required")
        if len(buf) == 0 {
            return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResp{
                Code:    fiber.StatusBadRequest,
                Message: err.Error(),
            })
        }
        return err
    }

    var filePath = filepath.Join(types.EVDataDir.Get().Value, "blogs", id+".md")
    if len(buf) == 0 {
        return c.SendFile(filePath)
    }

    var file, err = os.Open(filePath)

    *buf[0] = file
    return err
}

