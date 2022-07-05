FROM node:16 as builder

# Set up args
ARG ssh_key
ARG redis_host
ARG redis_port
ARG redis_db

# Set up env vars
ENV REDIS_HOST=$redis_host
ENV REDIS_PORT=$redis_port
ENV REDIS_DB=$redis_db

# Install packages
RUN apt-get update && apt-get install -y --no-install-recommends git openssh-client

# Set up SSH
RUN mkdir /root/.ssh && echo "${ssh_key}" > /root/.ssh/id_rsa && chmod 400 /root/.ssh/id_rsa && touch /root/.ssh/known_hosts && ssh-keyscan github.com >> /root/.ssh/known_hosts

# Copy the code
RUN mkdir /app
COPY . /app

# Change the directory
WORKDIR /app

# Install dependencies
RUN git clone git@github.com:getcohesive/platform.git
RUN npm install

# Replace variables in config file
RUN python -c 'import os,sys; sys.stdout.write(os.path.expandvars(sys.stdin.read()))' < conf.cohesive.yaml > conf.yaml


FROM node:16

# Copy the application
RUN mkdir /app
COPY --from=builder /app /app

# Change the directory
WORKDIR /app

# Run the application
CMD npm start

