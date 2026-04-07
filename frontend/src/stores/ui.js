import { reactive, readonly } from 'vue'

const state = reactive({
  locale: localStorage.getItem('locale') || 'en',
})

const messages = {
  en: {
    navProblems: 'Problems',
    navAdmin: 'Admin',
    navLogin: 'Login',
    navSignup: 'Sign Up',
    navLogout: 'Logout',
    adminTitle: 'Admin Problem Management',
    adminSubtitle: 'Create, update, and delete algorithm problems.',
    createSection: 'Create Problem',
    manageSection: 'Manage Existing',
  },
  ru: {
    navProblems: 'Задачи',
    navAdmin: 'Админ',
    navLogin: 'Вход',
    navSignup: 'Регистрация',
    navLogout: 'Выход',
    adminTitle: 'Управление задачами',
    adminSubtitle: 'Создание, редактирование и удаление задач.',
    createSection: 'Создать задачу',
    manageSection: 'Управление существующими',
  },
  tm: {
    navProblems: 'Meseleler',
    navAdmin: 'Admin',
    navLogin: 'Giriş',
    navSignup: 'Hasaba al',
    navLogout: 'Çyk',
    adminTitle: 'Mesele Dolandyryşy',
    adminSubtitle: 'Mesele döretmek, täzelemek we pozmak.',
    createSection: 'Mesele döret',
    manageSection: 'Bar bolanlary dolandyrmak',
  },
}

export function useUIStore() {
  const setLocale = (locale) => {
    state.locale = locale
    localStorage.setItem('locale', locale)
  }

  const t = (key) => messages[state.locale]?.[key] || messages.en[key] || key

  return {
    state: readonly(state),
    setLocale,
    t,
  }
}
