import fs from 'fs'
import path from 'path'
import axios from 'axios'
import { prompt } from 'enquirer'

const installDir = `/usr/local/bin`

export async function bget(source: string) {
  const { repo, tag } = parseSource(source)
  const res = await axios
    .get<any>(`https://api.github.com/repos/${repo}/releases/${tag}`)
    .catch((error) => {
      if (error.response) {
        if (error.response.status === 404) {
          throw new Error(`this release does not exist`)
        } else {
          throw new Error(
            `unable to fetch release info: ${error.response.status}`,
          )
        }
      } else {
        throw error
      }
    })

  const assets = res.data.assets
  if (assets.length === 0) {
    throw new Error(`no assets in this release`)
  }

  const answers = await prompt<{ downloadURL: string; outputPath: string }>([
    {
      type: 'select',
      name: 'downloadURL',
      message: 'Choose an asset to download and install',
      choices: assets.map((asset: any) => {
        return {
          message: asset.name,
          name: asset.browser_download_url,
        }
      }),
    },
    {
      type: 'text',
      name: 'outputPath',
      message: 'Install the binary to',
      initial: path.join(installDir, source.split('/')[1]),
    },
  ])

  console.log(`Downloading ${answers.downloadURL}`)
  await downloadFile(answers.downloadURL, answers.outputPath)
  fs.promises.chmod(answers.outputPath, 0o755)
  console.log(`Installed to ${answers.outputPath}`)
}

async function downloadFile(url: string, outputLocationPath: string) {
  const writer = fs.createWriteStream(outputLocationPath)

  return axios
    .get<any>(url, {
      responseType: 'stream',
    })
    .then((response) => {
      // ensure that the user can call `then()` only when the file has
      // been downloaded entirely.

      return new Promise((resolve, reject) => {
        response.data.pipe(writer)
        let error: Error | null = null
        writer.on('error', (err) => {
          error = err
          writer.close()
          reject(err)
        })
        writer.on('close', () => {
          if (!error) {
            resolve(true)
          }
          // no need to call the reject here, as it will have been called in the
          // 'error' stream;
        })
      })
    })
}

function parseSource(source: string) {
  const [repo, tag] = source.split('@')
  return {
    repo,
    tag: tag || 'latest',
  }
}
