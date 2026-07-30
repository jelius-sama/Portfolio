package renderer

import (
    "net/url"

    "git.jelius.dev/jelius-sama/Portfolio/types"

    "github.com/gofiber/fiber/v3"
)

type ViewManager struct{}

func New() *ViewManager {
    return &ViewManager{}
}

type TargetPart uint8

const (
    TPHead TargetPart = iota
    TPBody
    TPBoth
    TPUndefined
)

// PageContext holds metadata about the current request lifecycle
type PageContext struct {
    Metadata   types.Metadata
    IsPartial  bool
    TargetPart TargetPart
}

func (tp TargetPart) String() string {
    switch tp {
    case TPBody:
        return "body"
    case TPHead:
        return "head"
    case TPBoth:
        return "both"
    default:
        return ""
    }
}

func (tp TargetPart) Into(target string) TargetPart {
    switch target {
    case "body":
        return TPBody
    case "head":
        return TPHead
    case "both":
        return TPBoth
    default:
        return TPUndefined
    }
}

func (pc *PageContext) DetermineRenderMode(c fiber.Ctx) {
    if c.Get("HX-Request") == "true" {
        pc.IsPartial = true

        var parsedURL, _ = url.Parse(c.Get("HX-Current-URL"))

        if parsedURL.Path != c.Path() {
            pc.TargetPart = TPUndefined.Into("both")
        } else {
            pc.TargetPart = TPUndefined.Into("body")
        }
        return
    }

    // Default to full page load
    pc.IsPartial = false
    pc.TargetPart = TPUndefined
}

