FROM python:3

WORKDIR "/code"

RUN pip3 install bandit

COPY tmp/src/*.py ./

CMD ["bandit", "-r", "/code"]