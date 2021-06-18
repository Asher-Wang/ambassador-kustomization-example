FROM node:14-alpine

WORKDIR /opt/node_app

COPY package.json yarn.lock ./

RUN yarn install

COPY . .

RUN yarn build

CMD ["yarn", "start-nodemon"]
