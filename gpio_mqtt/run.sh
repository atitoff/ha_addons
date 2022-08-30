#!/usr/bin/env bashio

# #!/usr/bin/env bash
CONFIG_PATH=/data/options.json


MQTT_HOST=$(bashio::services mqtt "host")
MQTT_USER=$(bashio::services mqtt "username")
MQTT_PASSWORD=$(bashio::services mqtt "password")




if [[ -r "$CONFIG_PATH" ]]
then
  MqttPort="$(jq --raw-output '.MqttPort // empty' $CONFIG_PATH)"
  MqttClientId="$(jq --raw-output '.MqttClientId // empty' $CONFIG_PATH)"
  LogLevel="$(jq --raw-output '.LogLevel // empty' $CONFIG_PATH)"
fi


exec /gpio_mqtt $MQTT_HOST $MqttClientId $MQTT_USER $MQTT_PASSWORD $MqttPort $LogLevel
