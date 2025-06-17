import { createLogger, format, transports } from 'winston'
import { config } from '../config'

const logger = createLogger({
  level: config.log.level,
  format:
    config.log.format === 'json'
      ? format.combine(format.timestamp(), format.json())
      : format.combine(
          format.colorize(),
          format.timestamp(),
          format.printf(({ level, message, timestamp }) => `${timestamp} [${level}]: ${message}`)
        ),
  transports: [new transports.Console()]
})

export { logger }
