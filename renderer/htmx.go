package renderer

import "github.com/gofiber/fiber/v3"

type ViewManager struct{}

func New() *ViewManager {
    return &ViewManager{}
}

type TargetPart uint8

const (
    TPHead TargetPart = iota
    TPBody
    TPUndefined
)

// PageContext holds metadata about the current request lifecycle
type PageContext struct {
    Title      string
    IsPartial  bool
    TargetPart TargetPart
}

func (tp TargetPart) String() string {
    switch tp {
    case TPBody:
        return "body"
    case TPHead:
        return "head"
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
    default:
        return TPUndefined
    }
}

// DetermineRenderMode inspects HTTP headers to see what the frontend requested
func (pc *PageContext) DetermineRenderMode(c fiber.Ctx) {
    if c.Get("X-SPA-Request") != "" {
        pc.IsPartial = true
        pc.TargetPart = TargetPart.Into(TPUndefined, c.Get("X-SPA-Target"))
        return
    }

    // Default to full page load
    pc.IsPartial = false
    pc.TargetPart = TPUndefined
}

