package protoplugin

import (
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
)

type plugin struct {
	extension string
	generator Generator
	options   PluginOptions
}

func newPlugin(extension string, generator Generator, options PluginOptions) *plugin {
	return &plugin{extension, generator, options}
}

func (p *plugin) Handle() error {
	codeGeneratorRequest, err := p.getCodeGeneratorRequest()
	if err != nil {
		return err
	}
	if codeGeneratorRequest.Parameter != nil {
		if err := p.setParameters(codeGeneratorRequest.GetParameter()); err != nil {
			return err
		}
	}
	codeGeneratorResponseFiles, err := p.getCodeGeneratorResponseFiles(codeGeneratorRequest)
	if err != nil {
		return p.writeError(err)
	}
	return p.writeCodeGeneratorResponseFiles(codeGeneratorResponseFiles)
}

func (p *plugin) getCodeGeneratorRequest() (*google_protobuf_compiler.CodeGeneratorRequest, error) {
	data, err := p.readIn()
	if err != nil {
		return nil, err
	}
	codeGeneratorRequest := &google_protobuf_compiler.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, codeGeneratorRequest); err != nil {
		return nil, err
	}
	return codeGeneratorRequest, nil
}

func (p *plugin) readIn() ([]byte, error) {
	if p.options.In != nil {
		return ioutil.ReadAll(p.options.In)
	}
	return ioutil.ReadAll(os.Stdin)
}

func (p *plugin) setParameters(parameters string) error {
	for _, parameter := range strings.Split(parameters, ",") {
		split := strings.SplitN(parameter, "=", 2)
		if len(split) == 1 {
			if err := flag.CommandLine.Set(parameter, ""); err != nil {
				return err
			}
		} else {
			if err := flag.CommandLine.Set(split[0], split[1]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *plugin) getCodeGeneratorResponseFiles(codeGeneratorRequest *google_protobuf_compiler.CodeGeneratorRequest) ([]*google_protobuf_compiler.CodeGeneratorResponse_File, error) {
	var codeGeneratorResponseFiles []*google_protobuf_compiler.CodeGeneratorResponse_File
	for _, fileToGenerate := range codeGeneratorRequest.FileToGenerate {
		fileDescriptorProto, err := p.getProtoFile(fileToGenerate, codeGeneratorRequest.ProtoFile)
		if err != nil {
			return nil, err
		}
		codeGeneratorResponseFile, err := p.generate(fileDescriptorProto)
		if err != nil {
			return nil, err
		}
		if codeGeneratorResponseFile != nil {
			codeGeneratorResponseFiles = append(codeGeneratorResponseFiles, codeGeneratorResponseFile)
		}
	}
	return codeGeneratorResponseFiles, nil
}

func (p *plugin) getProtoFile(fileToGenerate string, fileDescriptorProtos []*descriptor.FileDescriptorProto) (*descriptor.FileDescriptorProto, error) {
	for _, fileDescriptorProto := range fileDescriptorProtos {
		if fileDescriptorProto.GetName() == fileToGenerate {
			return fileDescriptorProto, nil
		}
	}
	return nil, fmt.Errorf("protoplugin: no FileDescriptorProto for %s", fileToGenerate)
}

func (p *plugin) generate(fileDescriptorProto *descriptor.FileDescriptorProto) (*google_protobuf_compiler.CodeGeneratorResponse_File, error) {
	reader, err := p.generator.Generate(fileDescriptorProto)
	if err != nil {
		return nil, err
	}
	if reader == nil {
		return nil, nil
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if p.options.GoFmt {
		data, err = format.Source(data)
		if err != nil {
			return nil, err
		}
	}
	return &google_protobuf_compiler.CodeGeneratorResponse_File{
		Name:    proto.String(p.getOutFileName(fileDescriptorProto)),
		Content: proto.String(string(data)),
	}, nil
}

func (p *plugin) getOutFileName(fileDescriptorProto *descriptor.FileDescriptorProto) string {
	name := fileDescriptorProto.GetName()
	return fmt.Sprintf("%s.%s", strings.TrimSuffix(name, filepath.Ext(name)), p.extension)
}

func (p *plugin) writeCodeGeneratorResponseFiles(codeGeneratorResponseFiles []*google_protobuf_compiler.CodeGeneratorResponse_File) error {
	if len(codeGeneratorResponseFiles) > 0 {
		return p.writeCodeGeneratorResponse(&google_protobuf_compiler.CodeGeneratorResponse{File: codeGeneratorResponseFiles})
	}
	return nil
}

func (p *plugin) writeError(err error) error {
	return p.writeCodeGeneratorResponse(&google_protobuf_compiler.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func (p *plugin) writeCodeGeneratorResponse(codeGeneratorResponse *google_protobuf_compiler.CodeGeneratorResponse) error {
	data, err := proto.Marshal(codeGeneratorResponse)
	if err != nil {
		return err
	}
	var writer io.Writer = os.Stdout
	if p.options.Out != nil {
		writer = p.options.Out
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}
	return nil
}
