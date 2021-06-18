import Koa from 'koa';
import bodyParser from 'koa-bodyparser';

import {config} from './config';
import {router} from './router';
import {loggingMiddleware} from './utils';

const app = new Koa();

app.use(bodyParser());

if (config.debug) {
  app.use(loggingMiddleware);
}

app.use(router.routes()).use(router.allowedMethods());

app.listen(config.port, () => {
  console.info(`BFF mock service listening on port: ${config.port}...`);
});
