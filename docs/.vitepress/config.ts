import {defineConfig} from 'vitepress'

export default defineConfig({
    title: 'wb-go',
    lang: 'ru-RU',
    base: '/wb-go',
    description: 'Библиотека для работы с устройствами Wiren Board.',
    themeConfig: {
        search: {
            provider: 'local'
        },
        nav: [
            {text: 'GitHub', link: 'https://github.com/ValentinAlekhin/wb-go'},
        ],
        sidebar: [
            {
                text: "Введение",
                items: [
                    {
                        text: 'О библиотеке',
                        link: '/guide/about'
                    },
                    {
                        text: 'Установка',
                        link: '/guide/install'
                    }
                ]
            },
            {
                text: 'Команды',
                items: [
                    {text: 'Generate', link: '/guide/commands/generate'},
                ]
            },
        ]
    }
})