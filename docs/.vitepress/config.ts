import {defineConfig} from 'vitepress'

export default defineConfig({
    title: 'wb-go',
    lang: 'ru-RU',
    description: 'Библиотека для работы с устройствами Wirenboard.',
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