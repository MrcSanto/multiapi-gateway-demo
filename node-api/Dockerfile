# Use the official lightweight Node.js image as base
FROM node:18-alpine

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+ automatically generates a package-lock.json)
COPY package*.json ./

# Install only production dependencies
RUN npm install --only=production

# Bundle app source
COPY . .

# Expose the port the app runs on
EXPOSE 8000

# Start the server
CMD [ "npm", "start" ]