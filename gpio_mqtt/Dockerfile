ARG BUILD_FROM
FROM $BUILD_FROM

# Copy data for add-on
COPY gpio_mqtt /
COPY run.sh /
COPY index.html /
COPY static /


RUN chmod a+x /gpio_mqtt
RUN chmod a+x /run.sh

ENTRYPOINT ["/run.sh"]