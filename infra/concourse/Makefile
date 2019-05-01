FLY_CMD := "fly"
PIPELINE_CONFS := $(wildcard ./pipelines/*.yml)
PIPELINES := $(foreach CONF, $(PIPELINE_CONFS), $(subst terraform-google-,,$(basename $(notdir $(CONF)))))

.PHONY: $(PIPELINES)
$(PIPELINES):
	$(FLY_CMD) -t cft set-pipeline \
		-p terraform-google-$@ \
		-c pipelines/terraform-google-$@.yml \
		-l vars/phoogle.yml

check-%:
	$(FLY_CMD) -t cft validate-pipeline -s \
		-c pipelines/terraform-google-$*.yml \
		-l vars/phoogle.yml

POSTGRES_IP := 10.45.0.3

.PHONY: psql
psql:
	kubectl run -i --tty --attach psql --image=postgres:9.6-alpine -- \
		psql -h ${POSTGRES_IP} -U concourse -d atc
	kubectl delete deployment psql

include ./Makefile.BUILD
