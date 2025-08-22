FROM node:22.16.0 as build

# Create app directory
WORKDIR /usr/src/app
# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# Copying this separately prevents re-running npm install on every code change
COPY . .

RUN yarn

RUN yarn build

FROM node:22.16.0 as production

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/.nuxt ./.nuxt
COPY --from=build /usr/src/app/.output ./.output

CMD node .output/server/index.mjs
EXPOSE 3000