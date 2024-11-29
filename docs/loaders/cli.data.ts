import fs from 'node:fs'
import {parse} from 'yaml'
import {defineLoader} from 'vitepress'

interface CLICommand {
    name: string;              // Название команды
    synopsis: string;          // Краткое описание
    description: string;       // Подробное описание
    usage: string;             // Инструкция по использованию
    options: CLIOption[];      // Список флагов/опций команды
    see_also: string[];        // Ссылки на связанные команды
}

interface CLIOption {
    name: string;              // Название флага
    shorthand: string;         // Короткая форма флага (однобуквенная)
    usage: string;             // Описание флага
    default_value?: string;    // Значение по умолчанию (необязательно)
}

type Data = Record<string, CLICommand>;

declare const data: Data
export {data}


export default defineLoader({
    watch: ['../data/cli/*.yaml'],
    load(watchedFiles): Data {
        return watchedFiles.map((file) => parse(fs.readFileSync(file, 'utf-8')) as CLICommand).reduce((acc, cur) => {
            const cmd = cur.name.split(' ')[1]
            if (!cmd) return acc

            acc[cmd] = cur

            return acc
        }, {})
    }
})
