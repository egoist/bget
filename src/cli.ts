#!/usr/bin/env node
import { version } from '../package.json'
import { cac } from 'cac'

const cli = cac(`bget`)

cli
  .command(`<source>`, `Download and install binaries from GitHub Releases`)
  .action(async (source) => {
    const { bget } = await import('./')
    await bget(source).catch((error) => {
      if (error instanceof Error) {
        console.error(error.stack)
      }
      process.exit(1)
    })
  })

cli.version(version)
cli.help()
cli.parse()
