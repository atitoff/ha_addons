ARG BUILD_ARCH
ARG BUILD_FROM
FROM $BUILD_FROM

# Copy data for add-on
COPY ${BUILD_ARCH}_gpio_mqtt /gpio_mqtt
COPY run.sh /
COPY web/ /web/
COPY index.html /

RUN chmod a+x /gpio_mqtt
RUN chmod a+x /run.sh

ENTRYPOINT ["/run.sh"]