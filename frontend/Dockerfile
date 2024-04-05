# Dockerfile

# Start from the latest Node base image
FROM node:20.2-slim

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy package.json and package-lock.json files
COPY package.json ./
COPY yarn.lock ./

# Install all dependencies. Dependencies will be cached if the package files are not changed
RUN yarn install

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the React app
RUN yarn build

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to start the server
CMD ["yarn", "start"]
