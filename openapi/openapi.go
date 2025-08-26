package openapi

import (
	"reflect"
	"regexp"

	"github.com/gofiber/fiber/v3"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
)

type Wrapper struct {
	router    fiber.Router
	reflector openapi31.Reflector
}

type Register struct {
	wrapper *Wrapper
	path    string
}

func HasBody(request interface{}) bool {
	if request == nil {
		return false
	}

	requestType := reflect.TypeOf(request)
	if requestType == nil {
		return false
	}

	if requestType.Kind() == reflect.Pointer {
		requestType = requestType.Elem()
	}

	if requestType.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < requestType.NumField(); i++ {
		field := requestType.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" {
			return true
		}
	}
	return false
}

func NewWrapper(router fiber.Router) *Wrapper {
	var reflector openapi31.Reflector
	if !fiber.IsChild() {
		reflector := openapi31.Reflector{}
		reflector.Spec = &openapi31.Spec{Openapi: "3.1.0"}
		reflector.Spec.Info.
			WithTitle("Dimiplan Backend").
			WithVersion("2.0.0").
			WithDescription("Dimiplan Backend Rewritten in Go + Fiber")
	}
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

func (r *Register) Get(handler func(request interface{}, c fiber.Ctx) (interface{}, error), requestFactory func() interface{}, response interface{}, status int) *Register {
	if !fiber.IsChild() {
		re := regexp.MustCompile(`:([^/]+)`)
		openAPIPath := re.ReplaceAllString(r.path, "{$1}")

		op, err := r.wrapper.reflector.NewOperationContext(fiber.MethodGet, openAPIPath)
		if err != nil {
			panic(err)
		}
		if requestFactory != nil {
			op.AddReqStructure(requestFactory())
		}
		op.AddRespStructure(response, func(cu *openapi.ContentUnit) { cu.HTTPStatus = status })
		r.wrapper.reflector.AddOperation(op)
	}
	r.wrapper.router.Get(r.path, func(c fiber.Ctx) error {
		var request interface{}
		if requestFactory != nil {
			request = requestFactory()
		}
		value, err := handler(request, c)
		if err != nil {
			return err
		}
		if value == nil {
			return c.SendStatus(status)
		}
		return c.Status(status).JSON(value)
	})
	return r
}

func (r *Register) Post(handler func(request interface{}, c fiber.Ctx) (interface{}, error), requestFactory func() interface{}, response interface{}, status int) *Register {
	if !fiber.IsChild() {
		re := regexp.MustCompile(`:([^/]+)`)
		openAPIPath := re.ReplaceAllString(r.path, "{$1}")

		op, err := r.wrapper.reflector.NewOperationContext(fiber.MethodPost, openAPIPath)
		if err != nil {
			panic(err)
		}
		if requestFactory != nil {
			sampleRequest := requestFactory()
			op.AddReqStructure(sampleRequest)
		}
		op.AddRespStructure(response, func(cu *openapi.ContentUnit) { cu.HTTPStatus = status })
		r.wrapper.reflector.AddOperation(op)
	}
	r.wrapper.router.Post(r.path, func(c fiber.Ctx) error {
		var request interface{}
		if requestFactory != nil {
			request = requestFactory()
			if !HasBody(&request) {
				request = nil
			}
			if request != nil {
				if err := c.Bind().Body(request); err != nil {
					return fiber.NewError(fiber.StatusBadRequest, err.Error())
				}
			}
		}
		value, err := handler(request, c)
		if err != nil {
			return err
		}
		if value == nil {
			return c.SendStatus(status)
		}
		return c.Status(status).JSON(value)
	})
	return r
}

func (r *Register) Patch(handler func(request interface{}, c fiber.Ctx) (interface{}, error), requestFactory func() interface{}, response interface{}, status int) *Register {
	if !fiber.IsChild() {
		re := regexp.MustCompile(`:([^/]+)`)
		openAPIPath := re.ReplaceAllString(r.path, "{$1}")

		op, err := r.wrapper.reflector.NewOperationContext(fiber.MethodPatch, openAPIPath)
		if err != nil {
			panic(err)
		}
		if requestFactory != nil {
			sampleRequest := requestFactory()
			op.AddReqStructure(sampleRequest)
		}
		op.AddRespStructure(response, func(cu *openapi.ContentUnit) { cu.HTTPStatus = status })
		r.wrapper.reflector.AddOperation(op)
	}
	r.wrapper.router.Patch(r.path, func(c fiber.Ctx) error {
		var request interface{}
		if requestFactory != nil {
			request = requestFactory()
			if !HasBody(&request) {
				request = nil
			}
			if request != nil {
				if err := c.Bind().Body(request); err != nil {
					return fiber.NewError(fiber.StatusBadRequest, err.Error())
				}
			}
		}
		value, err := handler(request, c)
		if err != nil {
			return err
		}
		if value == nil {
			return c.SendStatus(status)
		}
		return c.Status(status).JSON(value)
	})
	return r
}

func (r *Register) Delete(handler func(request interface{}, c fiber.Ctx) (interface{}, error), requestFactory func() interface{}, response interface{}, status int) *Register {
	if !fiber.IsChild() {
		re := regexp.MustCompile(`:([^/]+)`)
		openAPIPath := re.ReplaceAllString(r.path, "{$1}")

		op, err := r.wrapper.reflector.NewOperationContext(fiber.MethodDelete, openAPIPath)
		if err != nil {
			panic(err)
		}
		if requestFactory != nil {
			sampleRequest := requestFactory()
			op.AddReqStructure(sampleRequest)
		}
		op.AddRespStructure(response, func(cu *openapi.ContentUnit) { cu.HTTPStatus = status })
		r.wrapper.reflector.AddOperation(op)
	}
	r.wrapper.router.Delete(r.path, func(c fiber.Ctx) error {
		var request interface{}
		if requestFactory != nil {
			request = requestFactory()
			if !HasBody(&request) {
				request = nil
			}
			if request != nil {
				if err := c.Bind().Body(request); err != nil {
					return fiber.NewError(fiber.StatusBadRequest, err.Error())
				}
			}
		}
		value, err := handler(request, c)
		if err != nil {
			return err
		}
		if value == nil {
			return c.SendStatus(status)
		}
		return c.Status(status).JSON(value)
	})
	return r
}
