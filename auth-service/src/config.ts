import FS from 'fs';
import {GLOBAL_CONFIG_PATH, PROJECT_CONFIG_PATH} from './paths';
import {camelToSnakeCase} from './utils';

export interface Config {
  port: number;
  debug: boolean;
  accessTokens: string[];
}

const DEFAULT_CONFIG: Config = {
  port: 3000,
  debug: true,
  accessTokens: [],
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
    let _value = process.env[envKey];

    if (!_value) {
      console.warn(`Env variable ${envKey} is empty, ignoring...`);
      continue;
    }

    switch (envKey) {
      // number
      case 'port':
        value = Number(_value);
        if (isNaN(value)) {
          console.warn(`Env variable ${envKey} is not number, ignoring...`);
          continue;
        }
        break;
      // boolean
      case 'debug':
        value = _value && _value.toLocaleLowerCase() === 'true' ? true : false;
        break;
      // comma-split string array
      case 'accessTokens':
        {
          let items = _value
            .split(',')
            .map(item => item.trim())
            .filter(item => item.length > 0);
          if (items.length === 0) {
            console.warn(`Env variable ${envKey} is empty, ignoring...`);
            continue;
          }
          value = items;
        }
        break;
      // string
      default:
        value = _value;
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
