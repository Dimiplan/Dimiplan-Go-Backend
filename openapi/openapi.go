package openapi

import (
	"reflect"
	"regexp"

	"github.com/gofiber/fiber/v3"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
)

type Wrapper struct {
	router    fiber.Router
	reflector openapi3.Reflector
}

type Register struct {
	wrapper *Wrapper
	path    string
}

func NewWrapper(router fiber.Router) *Wrapper {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithTitle("Dimiplan Backend").
		WithVersion("2.0.0").
		WithDescription("Dimiplan Backend Rewritten in Go + Fiber")
	return &Wrapper{router: router, reflector: reflector}
}

func (w *Wrapper) Route(path string) *Register {
	return &Register{wrapper: w, path: path}
}

func (w *Wrapper) APIDocs() string {
	schema, err := w.reflector.Spec.MarshalYAML()
	if err != nil {
		return err.Error()
	}
	return string(schema)
}

func (r *Register) Route(path string) *Register {
	r.path += path
	return r
}

func (r *Register) Get(handler func(request interface{}, c fiber.Ctx) (interface{}, error), request interface{}, response interface{}, status int) *Register {
	// Convert Fiber path format (:param) to OpenAPI format ({param})
	re := regexp.MustCompile(`:([^/]+)`)
	openAPIPath := re.ReplaceAllString(r.path, "{$1}")

	op, err := r.wrapper.reflector.NewOperationContext(fiber.MethodGet, openAPIPath)
	if err != nil {
		panic(err)
	}
	op.AddReqStructure(request)
	op.AddRespStructure(response, func(cu *openapi.ContentUnit) { cu.HTTPStatus = status })
	r.wrapper.reflector.AddOperation(op)
	r.wrapper.router.Get(r.path, func(c fiber.Ctx) error {
		if request != nil {
			if err := c.Bind().All(request); err != nil {
				return fiber.ErrBadRequest
			}
		}
		value, err := handler(request, c)
		if value == nil || response == nil {
			c.SendStatus(status)
		}
		if err != nil || reflect.TypeOf(value) != reflect.TypeOf(response) {
			return err
		}
		return c.Status(status).JSON(value)
	})
	return r
}

func (r *Register) Post(handler func(request interface{}, c fiber.Ctx) (interface{}, error), request interface{}, response interface{}, status int) *Register {
	// Convert Fiber path format (:param) to OpenAPI format ({param})
	re := regexp.MustCompile(`:([^/]+)`)
	openAPIPath := re.ReplaceAllString(r.path, "{$1}")

	op, err := r.wrapper.reflector.NewOperationContext(fiber.MethodPost, openAPIPath)
	if err != nil {
		panic(err)
	}
	op.AddReqStructure(request)
	op.AddRespStructure(response, func(cu *openapi.ContentUnit) { cu.HTTPStatus = status })
	r.wrapper.reflector.AddOperation(op)
	r.wrapper.router.Post(r.path, func(c fiber.Ctx) error {
		if request != nil {
			if err := c.Bind().All(request); err != nil {
				return fiber.ErrBadRequest
			}
		}
		value, err := handler(request, c)
		if value == nil || response == nil {
			c.SendStatus(status)
		}
		if err != nil || reflect.TypeOf(value) != reflect.TypeOf(response) {
			return err
		}
		return c.Status(status).JSON(value)
	})
	return r
}

func (r *Register) Patch(handler func(request interface{}, c fiber.Ctx) (interface{}, error), request interface{}, response interface{}, status int) *Register {
	// Convert Fiber path format (:param) to OpenAPI format ({param})
	re := regexp.MustCompile(`:([^/]+)`)
	openAPIPath := re.ReplaceAllString(r.path, "{$1}")

	op, err := r.wrapper.reflector.NewOperationContext(fiber.MethodPatch, openAPIPath)
	if err != nil {
		panic(err)
	}
	op.AddReqStructure(request)
	op.AddRespStructure(response, func(cu *openapi.ContentUnit) { cu.HTTPStatus = status })
	r.wrapper.reflector.AddOperation(op)
	r.wrapper.router.Patch(r.path, func(c fiber.Ctx) error {
		if request != nil {
			if err := c.Bind().All(request); err != nil {
				return fiber.ErrBadRequest
			}
		}
		value, err := handler(request, c)
		if value == nil || response == nil {
			c.SendStatus(status)
		}
		if err != nil || reflect.TypeOf(value) != reflect.TypeOf(response) {
			return err
		}
		return c.Status(status).JSON(value)
	})
	return r
}

func (r *Register) Delete(handler func(request interface{}, c fiber.Ctx) (interface{}, error), request interface{}, response interface{}, status int) *Register {
	// Convert Fiber path format (:param) to OpenAPI format ({param})
	re := regexp.MustCompile(`:([^/]+)`)
	openAPIPath := re.ReplaceAllString(r.path, "{$1}")

	op, err := r.wrapper.reflector.NewOperationContext(fiber.MethodDelete, openAPIPath)
	if err != nil {
		panic(err)
	}
	op.AddReqStructure(request)
	op.AddRespStructure(response, func(cu *openapi.ContentUnit) { cu.HTTPStatus = status })
	r.wrapper.reflector.AddOperation(op)
	r.wrapper.router.Delete(r.path, func(c fiber.Ctx) error {
		if request != nil {
			if err := c.Bind().All(request); err != nil {
				return fiber.ErrBadRequest
			}
		}
		value, err := handler(request, c)
		if value == nil || response == nil {
			c.SendStatus(status)
		}
		if err != nil || reflect.TypeOf(value) != reflect.TypeOf(response) {
			return err
		}
		return c.Status(status).JSON(value)
	})
	return r
}
