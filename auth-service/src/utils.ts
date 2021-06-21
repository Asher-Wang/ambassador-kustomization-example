import {Middleware} from 'koa';

export function camelToSnakeCase(text: string): string {
  return text.replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`);
}

export function fatalError(message: string): never {
  console.error(message);
  process.exit(1);
}

export const loggingMiddleware: Middleware = async (context, next) => {
  let date = new Date().toUTCString();

  console.log(
    `[Request] ${date} - ${context.method} - ${context.path} - ${
      context.querystring
    } - ${JSON.stringify(context.request.body)}`,
  );

  await next();

  console.log(`[Response] ${context.status} ${JSON.stringify(context.body)}`);
};
