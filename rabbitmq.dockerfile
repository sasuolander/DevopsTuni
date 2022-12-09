FROM rabbitmq:3-management
COPY conf-rabbitmq/10-defaults.conf /etc/rabbitmq/conf.d/10-defaults.conf
COPY conf-rabbitmq/definitions.json /etc/rabbitmq/conf.d/definitions.json
RUN chown rabbitmq:rabbitmq /etc/rabbitmq/conf.d/10-defaults.conf
RUN chown rabbitmq:rabbitmq /etc/rabbitmq/conf.d/definitions.json
EXPOSE 15672
#ENTRYPOINT ["/bin/bash"]
