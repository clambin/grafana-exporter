FROM python:3.9-slim

WORKDIR /app

RUN groupadd -g 1000 abc && \
    useradd -u 1000 -g abc abc && \
    pip install --upgrade pip && \
    pip install pipenv

COPY Pip* ./
RUN pipenv install --system --deploy
COPY grafana_exporter.py ./

USER abc
ENTRYPOINT ["/usr/local/bin/python3", "grafana_exporter.py"]
CMD []