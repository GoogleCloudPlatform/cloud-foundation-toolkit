FROM gcr.io/cloud-builders/gcloud
  
RUN apt-get update
RUN apt-get install python-setuptools -y
RUN apt-get install npm -y
RUN apt-get install jq -y
RUN pip install yq
RUN npm install -g ajv-cli
RUN ln -s /usr/bin/nodejs /usr/bin/node

COPY docker-entrypoint.sh /root/
RUN chmod 777 /root/docker-entrypoint.sh

ENTRYPOINT ["/root/docker-entrypoint.sh"]

CMD []
