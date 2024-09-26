import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import zh from './locales/zh.json'
import en from './locales/en.json'


i18n
    .use(LanguageDetector)
    .use(initReactI18next)
    .init({
        resources: {
            zh: {
                name: '简体中文',
                translation: zh
            },
            en: {
                name: 'English',
                translation: en
            }
        },
        fallbackLng: 'zh',
        debug: true,
        interpolation: {
            escapeValue: false,
        },
        backend: {
            loadPath: '/locales/{{lng}}.json',
        }
    });

export default i18n;
