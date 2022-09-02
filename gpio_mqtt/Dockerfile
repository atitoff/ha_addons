ARG BUILD_FROM
FROM $BUILD_FROM

# Copy data for add-on
COPY gpio_mqtt /
COPY run.sh /
COPY web/index.html /
COPY web/bootstrap-dark.min.css /
COPY web/bootstrap.bundle.min.js /

RUN chmod a+x /gpio_mqtt
RUN chmod a+x /run.sh

ENTRYPOINT ["/run.sh"]