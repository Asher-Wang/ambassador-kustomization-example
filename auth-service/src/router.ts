import KoaRouter from '@koa/router';

const router = new KoaRouter();

router.get('/', context => {
  context.body = `
<html>
  <head>
    <title>Auth Service</title>
  </head>
  <body>
    <h1>Congratulations!</h1>
    <h4>Auth service is up and running.</h4>
  </body>
</html>`;
});

router.all('/auth', context => {
  let authorization = context.get('authorization');
});

export {router};
