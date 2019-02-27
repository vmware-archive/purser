FROM node:9.6.1 as builder

LABEL maintainer = "VMware <hkatyal@vmware.com>"
LABEL author = "Hemani Katyal <hkatyal@vmware.com>"

# set working directory
RUN mkdir /usr/src/app
WORKDIR /usr/src/app

# add `/usr/src/app/node_modules/.bin` to $PATH
ENV PATH /usr/src/app/node_modules/.bin:$PATH

# install and cache app dependencies
COPY package.json package-lock.json ./
RUN npm install
RUN npm install -g @angular/cli@6.2.1

# add purser application to the working directory
COPY . .

# start purser application
RUN npm run build

# Build a small nginx image
FROM nginx:latest
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /usr/src/app/dist /usr/share/nginx/html