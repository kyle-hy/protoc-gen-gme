package gme

import (
	"path"
	"strconv"
	"strings"

	generator "github.com/kyle-hy/protoc-gen-gme/generator"

	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Paths for packages used by code generated in this file,
// relative to the import_prefix of the generator.Generator.
const (
	logPath   = "git.giitllm.cn/platform/zlog"
	zapPath   = "go.uber.org/zap"
	codecPath = "git.giitllm.cn/platform/sdk/core/codec"
	corePath  = "git.giitllm.cn/platform/sdk/core"
	utilsPath = "git.giitllm.cn/platform/sdk/core/utils"
)

func init() {
	generator.RegisterPlugin(new(gme))
}

// gme is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for go-gme support.
type gme struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "gme".
func (g *gme) Name() string {
	return "gme"
}

// The names for packages imported in the generated code.
// They may vary from the final path component of the import path
// if the name is used by other packages.
var (
	pkgImports map[generator.GoPackageName]bool
)

// Init initializes the plugin.
func (g *gme) Init(gen *generator.Generator) {
	g.gen = gen
}

// P forwards to g.gen.P.
func (g *gme) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *gme) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.MessageType) == 0 {
		return
	}
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	// g.P("var _ ", contextPkg, ".Context")

	for i, messageType := range file.FileDescriptorProto.MessageType {
		g.generateMessageType(file, messageType, i)
	}
}

func (g *gme) generateMessageType(file *generator.FileDescriptor, messageType *pb.DescriptorProto, index int) {
	messageName := strings.Title(*messageType.Name)
	g.P("func RegHandler", messageName, "(register core.Register, handler func(ctx *core.Context, message *", messageName, ") error, opt ...core.RegHandleOption) {")
	g.P("if register != nil && handler != nil {")
	g.P("options := core.NewRegHandleOptions()")
	g.P("for _, o := range opt {")
	g.P("o(options)")
	g.P("}")
	g.P("messageName := proto.MessageName((*", messageName, ")(nil))")
	g.P("decodeFunc := func(unmarshaller codec.Unmarshaller, msg []byte) (interface{}, error) {")
	g.P("data := new(", messageName, ")")
	g.P("err := unmarshaller.Unmarshal(msg, data)")
	g.P("return data, err")
	g.P("}")
	g.P("handleFunc := func(ctx *core.CtxParam, data interface{}) error {")
	g.P("msg := data.(*", messageName, ")")
	g.P("if utils.RecvNeedLogging(ctx, options.LogRecvBody) {")
	g.P("fields := utils.LogFieldsFromCtx(ctx)")
	g.P("zlog.Info(\"recv msg\",append(fields, zap.Reflect(\"bady\", msg))...)")
	g.P("}")
	g.P("ctx.SetLogSendBody(utils.SendNeedLogging(ctx, options.LogSendBody))")
	g.P("gwctx := core.NewContext(ctx)")
	g.P("return handler(gwctx, msg)")
	g.P("}")
	g.P("register.RegHandleFunc(utils.MessageName(messageName), handleFunc, decodeFunc, opt...)")
	g.P("}")
	g.P("}")
	g.P()
}

// GenerateImports generates the import declaration for this file.
func (g *gme) GenerateImports(file *generator.FileDescriptor, imports map[generator.GoImportPath]generator.GoPackageName) {
	g.P("import (")
	g.P(strconv.Quote(path.Join(g.gen.ImportPrefix, logPath)))
	g.P(strconv.Quote(path.Join(g.gen.ImportPrefix, zapPath)))
	g.P(strconv.Quote(path.Join(g.gen.ImportPrefix, codecPath)))
	g.P(strconv.Quote(path.Join(g.gen.ImportPrefix, corePath)))
	g.P(strconv.Quote(path.Join(g.gen.ImportPrefix, utilsPath)))
	g.P(")")
	g.P()
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	// We need to keep track of imported packages to make sure we don't produce
	// a name collision when generating types.
	pkgImports = make(map[generator.GoPackageName]bool)
	for _, name := range imports {
		pkgImports[name] = true
	}
}
