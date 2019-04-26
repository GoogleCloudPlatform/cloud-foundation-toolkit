FLY_CMD := "fly"

PIPELINES = \
	container-vm \
	kubernetes-engine \
	network \
	project-factory

.PHONY: $(PIPELINES)
$(PIPELINES):
	$(FLY_CMD) -t cft set-pipeline \
		-p terraform-google-$@ \
		-c pipelines/terraform-google-$@.yml \
		-l vars/phoogle.yml

POSTGRES_IP := 10.45.0.3

.PHONY: psql
psql:
	kubectl run -i --tty --attach psql --image=postgres:9.6-alpine -- \
		psql -h ${POSTGRES_IP} -U concourse -d atc
	kubectl delete deployment psql

include ./Makefile.BUILD
