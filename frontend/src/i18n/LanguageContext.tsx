import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import { EventsOn } from 'wailsjs/runtime/runtime'
import { GetLanguage, SetLanguage } from 'wailsjs/go/main/App'
import translations, { type LangCode, type Translations } from './translations'

interface LanguageContextValue {
  lang: LangCode
  t: Translations
  setLang: (code: LangCode) => void
}

const LanguageContext = createContext<LanguageContextValue>({
  lang: 'en',
  t: translations.en,
  setLang: () => {},
})

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [lang, setLangState] = useState<LangCode>('en')

  useEffect(() => {
    GetLanguage()
      .then(saved => {
        if (saved && saved in translations) {
          setLangState(saved as LangCode)
        }
      })
      .catch(() => {})

    const off = EventsOn('language-changed', (newLang: string) => {
      if (newLang in translations) {
        setLangState(newLang as LangCode)
      }
    })
    return () => { off() }
  }, [])

  const setLang = (code: LangCode) => {
    setLangState(code)
    SetLanguage(code).catch(() => {})
  }

  return (
    <LanguageContext.Provider value={{ lang, t: translations[lang], setLang }}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useTranslation() {
  return useContext(LanguageContext)
}
