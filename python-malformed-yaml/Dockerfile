FROM python:3.12.0-alpine3.17

# We are installing a dependency here directly into our app source dir
RUN pip install pathlib
RUN pip install yaml

COPY main.py /main.py

ENTRYPOINT ["/main.py"]
