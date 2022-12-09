FROM rabbitmq:3-management
COPY conf-rabbitmq/10-defaults.conf /etc/rabbitmq/conf.d/10-defaults.conf
COPY conf-rabbitmq/definitions.json /etc/rabbitmq/conf.d/definitions.json
RUN chown rabbitmq:rabbitmq /etc/rabbitmq/conf.d/10-defaults.conf
RUN chown rabbitmq:rabbitmq /etc/rabbitmq/conf.d/definitions.json
#RUN rabbitmqctl add_user test test
#RUN rabbitmqctl set_user_tags test administrator
#RUN rabbitmqctl set_permissions -p / test ".*" ".*" ".*"
#ENV RABBITMQ_DEFAULT_USER=test1
#ENV RABBITMQ_DEFAULT_PASS=test1
RUN apt-get update
RUN apt-get install nano

EXPOSE 15672
#ENTRYPOINT ["/bin/bash"]
