build-local-container:
	docker buildx build \
			--load \
			--platform linux/arm64 \
			-t artel:local .
