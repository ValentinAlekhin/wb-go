import fs from 'fs'

export default {
    paths() {
        return fs
            .readdirSync('data/cli')
            .map((filename) => filename.split('_')[1])
            .filter((cmd) => cmd).map((cmd) => ({params: {cmd: cmd.replace('.yaml', '')}}))

    }
}