FROM fluent/fluentd:v1.17.1-debian-1.0

USER root

RUN apt-get update && apt-get install -y \
    build-essential \
    ruby-dev \
    libsystemd-dev \
    && gem install fluent-plugin-elasticsearch fluent-plugin-systemd --no-document \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

USER fluent

CMD fluentd -c /fluentd/etc/fluent.conf