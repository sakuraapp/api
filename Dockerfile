FROM node:14.6.0-alpine

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

ADD package.json /usr/src/app/
RUN npm i --production
ADD . /usr/src/app

ENTRYPOINT [ "npm", "start" ]
