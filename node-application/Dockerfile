FROM node:16 as builder

# Set up args
ARG ssh_key

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


FROM node:16

# Copy the application
RUN mkdir /app
COPY --from=builder /app /app

# Change the directory
WORKDIR /app

# Run the application
CMD npm start

