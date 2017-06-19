IMAGE=mazzolino/wilddns

default: image

wilddns: main.go
	dapper -m bind

image: wilddns
	docker build -t ${IMAGE} .

push: image
	docker push ${IMAGE}
