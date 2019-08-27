FROM gcr.io/cft-test-workspace-221111/cft:latest
  
COPY cloud-foundation-tests.conf /etc/cloud-foundation-tests.conf
RUN cat /etc/cloud-foundation-tests.conf
RUN chmod 666 /etc/cloud-foundation-tests.conf
COPY docker-entrypoint.sh /root/
RUN chmod 777 /root/docker-entrypoint.sh


WORKDIR /cloud-foundation-toolkit/dm

ENTRYPOINT ["/root/docker-entrypoint.sh"]

CMD []
