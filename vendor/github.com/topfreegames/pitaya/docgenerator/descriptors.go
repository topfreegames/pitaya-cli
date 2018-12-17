package docgenerator

import (
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/topfreegames/pitaya/constants"
)

// ProtoDescriptors returns the descriptor for a given message name or .proto file
func ProtoDescriptors(protoName string) ([]byte, error) {
	if strings.HasSuffix(protoName, ".proto") {
		descriptor := proto.FileDescriptor(protoName)
		if descriptor == nil {
			return nil, constants.ErrProtodescriptor
		}
		return descriptor, nil
	}

	protoReflectTypePointer := proto.MessageType(protoName)

	if protoReflectTypePointer == nil {
		return nil, constants.ErrProtodescriptor
	}

	protoReflectType := protoReflectTypePointer.Elem()
	protoValue := reflect.New(protoReflectType)
	descriptorMethod, ok := protoReflectTypePointer.MethodByName("Descriptor")

	if !ok {
		return nil, constants.ErrProtodescriptor
	}

	descriptorValue := descriptorMethod.Func.Call([]reflect.Value{protoValue})
	protoDescriptor := descriptorValue[0].Bytes()

	return protoDescriptor, nil
}
