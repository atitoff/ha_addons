ARG BUILD_FROM
FROM $BUILD_FROM

# Copy data for add-on
COPY gpio_mqtt /
RUN chmod a+x /gpio_mqtt

CMD [ "/gpio_mqtt" ]