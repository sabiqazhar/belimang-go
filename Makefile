# This is a phony target, meaning it's a command name and not a file.
.PHONY: proto

# This rule runs the protoc command when you type "make proto".
proto:
	@echo "🔥 Generating Go code from protobuf..."
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/merchant/merchant.proto
	@echo "✅ Done."
