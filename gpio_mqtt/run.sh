#!/usr/bin/env bash

CONFIG_PATH=/data/options.json

if [[ -r "$CONFIG_PATH" ]]
then
  MqttPort="$(jq --raw-output '.MqttPort // empty' "$CONFIG_PATH")"
fi

export MqttPort="${MqttPort:-pi}"

exec /gpio_mqtt

