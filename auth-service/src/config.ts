import FS from 'fs';
import {GLOBAL_CONFIG_PATH, PROJECT_CONFIG_PATH} from './paths';
import {camelToSnakeCase} from './utils';

export interface Config {
  port: number;
  debug: boolean;
}

const DEFAULT_CONFIG: Config = {
  port: 3000,
  debug: true,
};

let configPath;

if (FS.existsSync(PROJECT_CONFIG_PATH)) {
  configPath = PROJECT_CONFIG_PATH;
} else if (FS.existsSync(GLOBAL_CONFIG_PATH)) {
  configPath = GLOBAL_CONFIG_PATH;
}

let config: Config = {...DEFAULT_CONFIG};

if (configPath) {
  console.log(`Config file found at ${configPath}.`);

  try {
    let configStr = FS.readFileSync(configPath, {encoding: 'utf8'});
    let configFromFile: Partial<Config> = JSON.parse(configStr);
    config = {...config, ...configFromFile};
    console.log(`- Config file successfully loaded.`);
  } catch (error) {
    console.error(
      `Failed to load config file (${configPath}): ${String(error)}`,
    );
  }
} else {
  console.log(`Config file not found, using default values...`);
}

for (let configKey of Object.keys(config)) {
  let envKey = camelToSnakeCase(configKey).toLocaleUpperCase();

  if (envKey in process.env) {
    let value;

    switch (envKey) {
      // number
      case 'port':
        value = Number(process.env[envKey]);
        if (isNaN(value)) {
          console.warn(`Env variable ${envKey} is not number, ignoring...`);
          continue;
        }
        break;
      // boolean
      case 'debug':
        let _value = process.env[envKey];
        value = _value && _value.toLocaleLowerCase() === 'true' ? true : false;
        break;
      // string
      default:
        value = process.env[envKey];
        if (!value) {
          console.warn(`Env variable ${envKey} is empty, ignoring...`);
          continue;
        }
        break;
    }

    (config as any)[configKey] = value;
    console.log(
      `Env variable ${envKey} found, overriding config: ${configKey}=${value}`,
    );
  }
}

console.log(`Final config: \n${JSON.stringify(config, undefined, 2)}`);

export {config};
