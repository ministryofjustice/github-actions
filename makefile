IMAGE := ministryofjustice/github-action-code-formatter
TAG := 1.0

build: .built-image

run: .built-image
	docker run --rm -it $(IMAGE) bash

push: .built-image
	docker tag $(IMAGE):$(TAG) $(IMAGE):$(TAG)
	# docker push $(IMAGE):$(TAG)

.built-image: Dockerfile makefile format-code.rb github.rb
	docker build -t $(IMAGE) .
	touch .built-image
