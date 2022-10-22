FROM python:3.10
# Install dapr CLI
RUN wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash

# Install daprd
ARG DAPR_BUILD_DIR
COPY $DAPR_BUILD_DIR /opt/dapr
ENV PATH="/opt/dapr/:${PATH}"
RUN dapr init --slim
RUN pip install flask
RUN pip install grpcio
RUN pip install grpcio-tools

# Install your app
WORKDIR /app
COPY . .
EXPOSE 3030
EXPOSE 3000
#RUN pip install requests
ENTRYPOINT ["dapr"]
CMD ["run", "--app-id", "client", "--dapr-grpc-port", "3030", "python", "main.py"]