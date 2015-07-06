echo "enum.proto"
protoc --proto_path=../../../../kit/prototype/src/common/protobuf/ --proto_path=../../../../kit/prototype/src/server/protobuf/ --go_out=./kit_ds ../../../../kit/prototype/src/server/protobuf/enum.proto
echo "msg_entity.proto"
protoc --proto_path=../../../../kit/prototype/src/common/protobuf/ --proto_path=../../../../kit/prototype/src/server/protobuf/ --go_out=./kit_ds ../../../../kit/prototype/src/server/protobuf/msg_entity.proto
echo "kit_ds.proto"
protoc --proto_path=../../../../kit/prototype/src/common/protobuf/ --proto_path=../../../../kit/prototype/src/server/protobuf/ --go_out=./kit_ds ../../../../kit/prototype/src/server/protobuf/kit_ds.proto
