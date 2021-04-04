package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/agilepathway/gauge-confluence/gauge_messages"
	"github.com/agilepathway/gauge-confluence/internal/confluence"
	"github.com/agilepathway/gauge-confluence/internal/env"
	"github.com/agilepathway/gauge-confluence/internal/git"
	"github.com/agilepathway/gauge-confluence/util"
	"google.golang.org/grpc"
)

const (
	gaugeSpecsDir = "GAUGE_SPEC_DIRS"
	fileSeparator = "||"
)

var projectRoot = util.GetProjectRoot() //nolint:gochecknoglobals

type handler struct {
	server *grpc.Server
}

func (h *handler) GenerateDocs(c context.Context, m *gauge_messages.SpecDetails) (*gauge_messages.Empty, error) {
	var ( //nolint:prealloc
		specsPaths []string // the absolute paths for all the specs
		specs      confluence.Specs
	)

	for _, providedSpecPath := range strings.Split(providedSpecsPaths(), fileSeparator) {
		specsPaths = append(specsPaths, util.GetFiles(providedSpecPath)...)
	}

	for _, specPath := range specsPaths {
		specs = append(specs, confluence.NewSpec(specPath, git.SpecGitURL(specPath, projectRoot)))
	}

	specs.PublishToConfluence()

	return &gauge_messages.Empty{}, nil
}

func (h *handler) Kill(c context.Context, m *gauge_messages.KillProcessRequest) (*gauge_messages.Empty, error) {
	defer h.stopServer()
	return &gauge_messages.Empty{}, nil
}

func (h *handler) stopServer() {
	h.server.Stop()
}

func main() {
	checkRequiredConfigVars()

	err := os.Chdir(projectRoot)
	util.Fatal("failed to change directory to project root.", err)

	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	util.Fatal("failed to start server.", err)

	l, err := net.ListenTCP("tcp", address)
	util.Fatal("TCP listening failed.", err)

	server := grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 10)) //nolint:gomnd
	h := &handler{server: server}
	gauge_messages.RegisterDocumenterServer(server, h)
	fmt.Printf("Listening on port:%d /n", l.Addr().(*net.TCPAddr).Port)
	server.Serve(l) //nolint:errcheck,gosec
}

func checkRequiredConfigVars() {
	env.GetRequired("CONFLUENCE_BASE_URL")
	env.GetRequired("CONFLUENCE_USERNAME")
	env.GetRequired("CONFLUENCE_TOKEN")
}

// providedSpecsPaths returns the list of specs paths passed in
// by the user of the plugin (converted to absolute paths by the
// core Gauge engine).
// Each spec path can be a directory or a spec file, as the `gauge docs`
// command accepts arguments in the same way as `gauge run`:
// https://docs.gauge.org/execution.html#multiple-arguments-passed-to-gauge-run
func providedSpecsPaths() string {
	return os.Getenv(gaugeSpecsDir)
}
